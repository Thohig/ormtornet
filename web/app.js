// Klienten: visar det servern skickar och skickar tangenttryck. Ingen spellogik här.
"use strict";

const $ = (id) => document.getElementById(id);

let ws = null;
let state = null;
let myName = "";
let joined = false;
let mapReady = false;

function show(id) {
  document.querySelectorAll(".screen").forEach((s) => s.classList.toggle("active", s.id === id));
}

// Loggen skriver fram texten tecken för tecken.
const Log = (() => {
  const queue = [];
  let busy = false;
  const MAX_LINES = 80;

  function push(text, cls) {
    if (!text) return;
    queue.push({ text, cls });
    if (!busy) next();
  }

  function next() {
    const item = queue.shift();
    if (!item) {
      busy = false;
      return;
    }
    busy = true;
    const box = $("log");
    const line = document.createElement("div");
    if (item.cls) line.className = item.cls;
    line.classList.add("cursor");
    box.appendChild(line);
    while (box.children.length > MAX_LINES) box.removeChild(box.firstChild);
    const chars = [...item.text];
    let i = 0;
    function tick() {
      // Skriv snabbare om det står mycket i kö.
      const step = queue.length > 2 ? 4 : 1;
      i = Math.min(chars.length, i + step);
      line.textContent = chars.slice(0, i).join("");
      box.scrollTop = box.scrollHeight;
      if (i < chars.length) {
        const ch = chars[i - 1];
        const pause = ".!?".includes(ch) ? 180 : ",".includes(ch) ? 70 : 22;
        setTimeout(tick, queue.length > 2 ? 6 : pause);
      } else {
        line.classList.remove("cursor");
        setTimeout(next, 120);
      }
    }
    tick();
  }

  return { push };
})();

function send(obj) {
  if (ws && ws.readyState === WebSocket.OPEN) ws.send(JSON.stringify(obj));
}

function connect() {
  const proto = location.protocol === "https:" ? "wss:" : "ws:";
  ws = new WebSocket(proto + "//" + location.host + "/ws");
  ws.onopen = () => {
    $("conn").textContent = "";
    send({ t: "join", name: myName });
  };
  ws.onmessage = (ev) => {
    let msg;
    try {
      msg = JSON.parse(ev.data);
    } catch {
      return;
    }
    handle(msg);
  };
  ws.onclose = () => {
    if (!joined) return;
    $("conn").textContent = "förbindelsen bröts, försöker igen";
    setTimeout(connect, 2000);
  };
}

function handle(msg) {
  switch (msg.t) {
    case "state":
      state = msg;
      if (!joined) {
        joined = true;
        try { localStorage.setItem("ormtornet-namn", myName); } catch {}
        show("game");
      }
      if (!mapReady) {
        MapView.setup($("map"), msg.size);
        mapReady = true;
      }
      MapView.draw(state);
      renderSide();
      break;
    case "log":
      Log.push(msg.text, msg.c);
      break;
    case "sound":
      Sound.play(msg.s);
      break;
    case "error":
      if (!joined) {
        $("nameerr").textContent = msg.text;
        ws.onclose = null;
        ws.close();
      } else {
        Log.push(msg.text, "varn");
      }
      break;
  }
}

function renderSide() {
  const ul = $("players");
  ul.innerHTML = "";
  state.players.forEach((p) => {
    const li = document.createElement("li");
    const mark = document.createElement("span");
    mark.textContent = "@ ";
    mark.style.color = MapView.PLAYER_COLORS[p.slot];
    li.appendChild(mark);
    li.appendChild(document.createTextNode(p.name + (p.id === state.you ? " (du)" : "") + (p.online ? "" : " (borta)")));
    if (p.id === state.you) li.classList.add("me");
    if (!p.online) li.classList.add("off");
    ul.appendChild(li);
  });
  const me = state.players.find((p) => p.id === state.you);
  const room = me && state.rooms.find((r) => r.x === me.x && r.y === me.y);
  $("here").textContent = room ? room.name : "";
}

$("enter").addEventListener("click", () => {
  Sound.start();
  show("namescreen");
  try { $("name").value = localStorage.getItem("ormtornet-namn") || ""; } catch {}
  $("name").focus();
});

$("nameform").addEventListener("submit", (e) => {
  e.preventDefault();
  myName = $("name").value.trim();
  if (!myName) {
    $("nameerr").textContent = "Tornet vill veta vad du heter.";
    return;
  }
  $("nameerr").textContent = "";
  connect();
});

const KEYS = { ArrowUp: "N", ArrowRight: "E", ArrowDown: "S", ArrowLeft: "W" };
document.addEventListener("keydown", (e) => {
  if (!joined) return;
  const dir = KEYS[e.key];
  if (!dir) return;
  e.preventDefault();
  send({ t: "move", dir });
});

// Rita om då och då så att den egna markören blinkar.
setInterval(() => {
  if (state && mapReady) MapView.draw(state);
}, 275);
