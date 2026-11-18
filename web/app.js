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
      Sound.breath(state.phase === "spel" ? state.sleep : -1);
      Sound.level(state.reveal ? state.level : 0);
      document.body.dataset.niva = state.reveal ? String(state.level) : "1";
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
    case "vaknar":
      MapView.wake();
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
  anfall: "anfall", parera: "parera", undan: "undanmanöver",
};
const DIR_WORDS = { N: "norrut", E: "österut", S: "söderut", W: "västerut" };

// Den döda spelaren som styr varelserna väljer en av dem med Tab.
let selectedMonster = 0;
function controlledMonster() {
  if (!state || !state.monsters || state.monsters.length === 0) return null;
  let m = state.monsters.find((x) => x.id === selectedMonster);
  if (!m) {
    m = state.monsters[0];
    selectedMonster = m.id;
  }
  return m;
}
function cycleMonster() {
  const list = (state && state.monsters) || [];
  if (list.length === 0) return;
  const i = list.findIndex((x) => x.id === selectedMonster);
  selectedMonster = list[(i + 1) % list.length].id;
  MapView.draw(state);
  renderSide();
}

// Mätare med fast bredd i kod, ett tecken per liv.
function meter(cls, label, value, maxValue, full, empty) {
  const d = document.createElement("div");
  d.className = cls;
  const a = document.createElement("span");
  a.textContent = full.repeat(Math.max(0, value));
  const b = document.createElement("span");
  b.className = "empty";
  b.textContent = empty.repeat(Math.max(0, maxValue - value));
  d.append(label, a, b);
  return d;
}

function playerTag(p) {
  if (!p.online) return "borta";
  if (p.status === "ute") return "ute";
  if (p.status === "forlorad") return "förlorad";
  if (p.status === "dod") return state.phase === "spel" && p.chosen ? "död, klar" : "död";
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
  if (me.status === "dod") line("Död", "cond");
  line("Skatt: " + me.treasure + " guld");
  if (state.potions > 0) line("Drycker: " + state.potions);
  if (me.status === "inne") {
    if (state.cellar > 0) line("I källarvalven", "cond");
    if (state.stuck > 0) line("Fast i " + state.stuck + (state.stuck === 1 ? " runda" : " rundor"), "cond");
    if (state.poison > 0) line("Förgiftad", "cond");
    if (state.dark > 0) line("Lyktan är släckt", "cond");
    if (state.sleep >= 0 && state.sleepText) line(state.sleepText, state.sleep >= 2 ? "cond" : "sleep");
    const foe = (state.monsters || []).find((m) => m.x === me.x && m.y === me.y);
    if (foe && state.cellar === 0) {
      line("Strid: " + foe.name, "foe");
      box.appendChild(meter("foehp", "", foe.hp, foe.maxHp, "■", "□"));
      line("A anfall · P parera · U undan", "dimline");
    }
  }
  if (state.controller) {
    line("Du styr tornets varelser", "cond");
    const m = controlledMonster();
    if (m) {
      line("Vald: " + m.name + " (Tab byter)", "foe");
      const a = state.action;
      if (a && a[0] === "m") line("Drag: " + a.replace(/^m\d+/, "").split("").map((d) => DIR_WORDS[d]).join(""), "plan");
      else if (a === "vanta") line("Drag: vila", "plan");
      else line("Välj ett drag", "plan");
    } else {
      line("Inga varelser att styra", "dimline");
    }
    return;
  }
  if (state.phase === "spel" && me.status === "inne") {
    line(state.action ? "Drag: " + ACTION_WORDS[state.action] : "Välj ditt drag", "plan");
  } else if (state.phase === "lobby") {
    line(me.ready ? "Du är redo" : "Tryck Enter när du är redo", "plan");
  } else if (state.phase === "slut") {
    line(me.ready ? "Väntar på de andra" : "Enter: nytt torn", "plan");
  }
}

function hpMeter(hp, maxHp) {
  return meter("hp", "Liv ", hp, maxHp, "♥", "♡");
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
    el.textContent = state.sun > 0 ? "vandringen är slut" : "solen har gått ner";
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

function endTitle() {
  return state.sun > 0 ? "VANDRINGEN ÄR SLUT" : "SOLEN HAR GÅTT NER";
}

function renderOverlay() {
  const el = $("overlay");
  if (state.phase !== "slut" || !state.results) {
    el.hidden = true;
    return;
  }
  const lines = [{ text: endTitle() }, { text: "" }];
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
  a: "anfall", A: "anfall", p: "parera", P: "parera", u: "undan", U: "undan",
};
document.addEventListener("keydown", (e) => {
  if (!joined || e.metaKey || e.ctrlKey || e.altKey) return;
  if (e.key === "Enter") {
    e.preventDefault();
    send({ t: "ready" });
    return;
  }
  if (state && state.controller) {
    // Den döda spelaren: Tab väljer varelse, pilarna flyttar den.
    if (e.key === "Tab") {
      e.preventDefault();
      cycleMonster();
      return;
    }
    const a = KEYS[e.key];
    if (a !== "vanta" && !DIR_WORDS[a]) return;
    e.preventDefault();
    const m = controlledMonster();
    if (a === "vanta") send({ t: "act", a });
    else if (m) send({ t: "act", a, m: m.id });
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
}, 100);
