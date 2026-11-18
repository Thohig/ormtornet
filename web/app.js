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
      if (state && state.phase !== msg.phase && msg.phase === "spel") lastTick = -1;
      state = msg;
      state.receivedAt = performance.now();
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
      renderTop();
      renderOverlay();
      break;
    case "log":
      Log.push(msg.text, msg.c);
      break;
    case "sound":
      Sound.play(msg.s);
      break;
    case "ping":
      MapView.ping(msg.x, msg.y, msg.slot);
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

const ACTION_WORDS = {
  N: "gå norrut", E: "gå österut", S: "gå söderut", W: "gå västerut",
  vanta: "vänta", ta: "ta skatt", sok: "sök i rummet", drick: "drick",
};

function playerTag(p) {
  if (!p.online) return "borta";
  if (p.status === "ute") return "ute";
  if (p.status === "forlorad") return "förlorad";
  if (p.status === "dod") return "död";
  if (state.phase === "spel") return p.chosen ? "klar" : "väljer";
  return p.ready ? "redo" : "";
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
    li.appendChild(document.createTextNode(p.name + (p.id === state.you ? " (du)" : "")));
    const tag = playerTag(p);
    if (tag) {
      const t = document.createElement("span");
      t.className = "tag";
      t.textContent = " " + tag;
      li.appendChild(t);
    }
    if (p.id === state.you) li.classList.add("me");
    if (!p.online) li.classList.add("off");
    ul.appendChild(li);
  });
  const me = myPlayer();
  const room = me && state.rooms.find((r) => r.x === me.x && r.y === me.y);
  let here = me && me.status === "inne" && room ? room.name : "";
  if (state.cellar > 0) here = "Källarvalven";
  else if (state.dark > 0 && room) here = "Mörker";
  $("here").textContent = here;

  const box = $("me");
  box.innerHTML = "";
  const line = (text, cls) => {
    const d = document.createElement("div");
    d.textContent = text;
    if (cls) d.className = cls;
    box.appendChild(d);
  };
  if (!me) return;
  box.appendChild(hpMeter(state.hp, state.maxHp));
  line("Skatt: " + me.treasure + " guld");
  if (state.potions > 0) line("Drycker: " + state.potions);
  if (me.status === "inne") {
    if (state.cellar > 0) line("I källarvalven", "cond");
    if (state.stuck > 0) line("Fast i " + state.stuck + (state.stuck === 1 ? " runda" : " rundor"), "cond");
    if (state.poison > 0) line("Förgiftad", "cond");
    if (state.dark > 0) line("Lyktan är släckt", "cond");
  }
  if (state.phase === "spel" && me.status === "inne") {
    line(state.action ? "Drag: " + ACTION_WORDS[state.action] : "Välj ditt drag", "plan");
  } else if (state.phase === "lobby") {
    line(me.ready ? "Du är redo" : "Tryck Enter när du är redo", "plan");
  } else if (state.phase === "slut") {
    line(me.ready ? "Väntar på de andra" : "Enter: nytt torn", "plan");
  }
}

// Livmätaren byggs med fast bredd i kod: ett hjärta per liv.
function hpMeter(hp, maxHp) {
  const d = document.createElement("div");
  d.className = "hp";
  const full = document.createElement("span");
  full.textContent = "♥".repeat(Math.max(0, hp));
  const empty = document.createElement("span");
  empty.className = "empty";
  empty.textContent = "♡".repeat(Math.max(0, maxHp - hp));
  d.append("Liv ", full, empty);
  return d;
}

function myPlayer() {
  return state && state.players.find((p) => p.id === state.you);
}

// Solmätaren byggs med fast bredd i kod.
const SUN_WIDTH = 24;
function renderTop() {
  const el = $("sun");
  el.innerHTML = "";
  const filled = Math.max(0, Math.round((state.sun / state.sunMax) * SUN_WIDTH));
  const label = document.createElement("span");
  label.textContent = "SOL ";
  const bar = document.createElement("span");
  bar.className = "bar";
  if (state.sun <= 3) bar.classList.add("last");
  else if (state.sun <= state.sunMax / 4) bar.classList.add("low");
  bar.textContent = "⣿".repeat(filled);
  const gone = document.createElement("span");
  gone.className = "gone";
  gone.textContent = "⣀".repeat(SUN_WIDTH - filled);
  const num = document.createElement("span");
  num.textContent = " " + String(state.sun).padStart(2, " ");
  el.append(label, bar, gone, num);
  updateTimer();
}

function remainingMs() {
  if (!state || state.phase !== "spel") return 0;
  return Math.max(0, state.deadlineMs - (performance.now() - state.receivedAt));
}

let lastTick = -1;
function updateTimer() {
  const el = $("turn");
  if (!state) return;
  if (state.phase === "lobby") {
    const n = state.players.filter((p) => p.online).length;
    el.textContent = n < state.minPlayers ? "väntar på vandrare" : "väntar på att alla är redo";
    el.classList.remove("urgent");
    return;
  }
  if (state.phase === "slut") {
    el.textContent = "solen har gått ner";
    el.classList.remove("urgent");
    return;
  }
  const secs = Math.ceil(remainingMs() / 1000);
  el.textContent = "runda " + state.round + "  " + String(secs).padStart(2, " ") + " s";
  const urgent = secs <= 5;
  el.classList.toggle("urgent", urgent);
  const me = myPlayer();
  if (urgent && secs > 0 && secs !== lastTick && me && me.status === "inne" && !state.action) {
    Sound.play("tick");
  }
  lastTick = secs;
}

// En ram av fast bredd runt raderna, byggd i kod.
function frame(lines, width) {
  const pad = (t) => {
    const chars = [...t];
    if (chars.length > width) return chars.slice(0, width).join("");
    const left = Math.floor((width - chars.length) / 2);
    return " ".repeat(left) + t + " ".repeat(width - chars.length - left);
  };
  const out = ["╭" + "─".repeat(width + 2) + "╮"];
  lines.forEach((l) => {
    const text = typeof l === "string" ? l : l.text;
    out.push({ text: "│ " + pad(text) + " │", cls: l.cls });
  });
  out.push("╰" + "─".repeat(width + 2) + "╯");
  return out.map((l) => (typeof l === "string" ? { text: l } : l));
}

function renderOverlay() {
  const el = $("overlay");
  if (state.phase !== "slut" || !state.results) {
    el.hidden = true;
    return;
  }
  const lines = [{ text: "SOLEN HAR GÅTT NER" }, { text: "" }];
  state.results.forEach((r, i) => {
    let what;
    if (r.status === "ute") what = r.treasure > 0 ? r.treasure + " guld" : "tomhänt";
    else if (r.status === "dod") what = "död";
    else what = "förlorad";
    const name = (r.winner ? "* " : "  ") + r.name;
    const row = name.padEnd(20, " ") + what.padStart(10, " ");
    lines.push({ text: row, cls: r.winner ? "win" : r.status !== "ute" ? "lost" : "" });
  });
  lines.push({ text: "" }, { text: "Enter: nytt torn" });
  el.innerHTML = "";
  frame(lines, 32).forEach((l) => {
    const d = document.createElement("div");
    d.textContent = l.text;
    if (l.cls) d.className = l.cls;
    el.appendChild(d);
  });
  el.hidden = false;
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

const KEYS = {
  ArrowUp: "N", ArrowRight: "E", ArrowDown: "S", ArrowLeft: "W",
  " ": "vanta", t: "ta", T: "ta", s: "sok", S: "sok", d: "drick", D: "drick",
};
document.addEventListener("keydown", (e) => {
  if (!joined || e.metaKey || e.ctrlKey || e.altKey) return;
  if (e.key === "Enter") {
    e.preventDefault();
    send({ t: "ready" });
    return;
  }
  const a = KEYS[e.key];
  if (!a) return;
  e.preventDefault();
  if (state && state.phase === "spel") send({ t: "act", a });
});

// Rita om då och då så att den egna markören blinkar.
setInterval(() => {
  if (state && mapReady) MapView.draw(state);
  if (state) updateTimer();
}, 250);
