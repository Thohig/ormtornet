// Klienten: visar det servern skickar och skickar tangenttryck. Ingen spellogik här.
"use strict";

const $ = (id) => document.getElementById(id);

let ws = null;
let state = null;
let myName = "";
let joined = false;
let mapReady = false;
let mapSize = 0; // tornets storlek när kartan ritades upp
let towerClosed = false; // spelledaren har stängt tornet, eller den här datorn leder det
let closedText = "tornet är stängt";
let hostingHere = false; // spelledarläget slogs på och den här datorn leder tornet

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
    if (towerClosed) {
      $("conn").textContent = closedText;
      return;
    }
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
      trackSeen(state);
      if (!joined) {
        joined = true;
        try { localStorage.setItem("ormtornet-namn", myName); } catch {}
        show("game");
      }
      if (!mapReady || mapSize !== msg.size) {
        MapView.setup($("map"), msg.size);
        mapReady = true;
        mapSize = msg.size;
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
    case "spelledare":
      // Spelledarläget slogs på, och den här datorn leder tornet.
      towerClosed = true;
      hostingHere = true;
      closedText = "du leder tornet";
      $("conn").textContent = closedText;
      Log.push(msg.text, "varn");
      if (ws) ws.close();
      showHosting(msg.text);
      break;
    case "stangt":
      towerClosed = true;
      $("conn").textContent = "tornet är stängt";
      Sound.stop();
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

// Strid mot andra vandrare: Tab väljer vem när flera står i rummet.
let selectedFoe = 0;
function foesHere() {
  const me = myPlayer();
  if (!me || me.status !== "inne") return [];
  return state.players.filter((p) => p.id !== me.id && !p.hidden && p.status === "inne" && p.x === me.x && p.y === me.y);
}
function currentFoe() {
  const foes = foesHere();
  if (!foes.length) return null;
  return foes.find((p) => p.id === selectedFoe) || foes[0];
}
function cycleFoe() {
  const foes = foesHere();
  if (!foes.length) return;
  const i = foes.findIndex((p) => p.id === (currentFoe() || {}).id);
  selectedFoe = foes[(i + 1) % foes.length].id;
  renderSide();
}

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

// Med synligheten "moten" eller "egen" syns andra vandrare bara i samma rum. Klienten
// minns var och när den såg var och en senast, så att kartan kan visa det.
let lastSeen = {};
function hidesOthers(st) {
  return st.sight === "moten" || st.sight === "egen";
}
function trackSeen(st) {
  if (st.phase !== "spel" || !hidesOthers(st)) {
    lastSeen = {};
  } else {
    st.players.forEach((p) => {
      if (p.id === st.you) return;
      if (p.status !== "inne") delete lastSeen[p.id];
      else if (!p.hidden) lastSeen[p.id] = { x: p.x, y: p.y, slot: p.slot, round: st.round };
    });
  }
  st.lastSeen = lastSeen;
}

function seenTag(p) {
  if (state.phase !== "spel" || !hidesOthers(state) || p.id === state.you) return "";
  if (p.status !== "inne" || !p.hidden) return "";
  const s = lastSeen[p.id];
  return s ? "sågs runda " + s.round : "inte sedd";
}

function playerTag(p) {
  if (!p.online) return "borta";
  if (p.status === "ute") return "ute";
  if (p.status === "forlorad") return "förlorad";
  if (p.status === "dod") return state.phase === "spel" && p.chosen ? "död, klar" : "död";
  if (state.phase === "spel") {
    const seen = seenTag(p);
    return (p.chosen ? "klar" : "väljer") + (seen ? ", " + seen : "");
  }
  return p.ready ? "redo" : "";
}

// Klassen: i lobbyn väljer man en av två, sedan står klassen överst.
function renderClass(box, line) {
  const offer = state.erbjudande || [];
  if (offer.length) {
    line("Välj klass" + (touch ? "" : " (1 eller 2)"), "plan");
    offer.forEach((c, i) => {
      const d = document.createElement("div");
      d.className = "classpick";
      d.tabIndex = 0;
      const title = document.createElement("div");
      title.className = "classname";
      title.textContent = (i + 1) + " " + c.namn;
      const st = document.createElement("div");
      st.className = "dimline";
      st.textContent = "+ " + c.styrka;
      const sv = document.createElement("div");
      sv.className = "dimline";
      sv.textContent = "− " + c.svaghet;
      d.append(title, st, sv);
      d.title = c.beskrivning;
      d.addEventListener("click", () => send({ t: "klass", name: c.id }));
      box.appendChild(d);
    });
    return;
  }
  const c = state.klass;
  if (!c) return;
  line(c.namn, "classname");
  if (state.phase !== "spel") {
    line(c.beskrivning, "dimline");
    line("+ " + c.styrka, "dimline");
    line("− " + c.svaghet, "dimline");
  }
}

function renderSide() {
  setEnterLabel(state.phase === "slut" ? "Nytt torn" : "Redo");
  const ul = $("players");
  ul.innerHTML = "";
  state.players.forEach((p) => {
    const li = document.createElement("li");
    const mark = document.createElement("span");
    mark.textContent = "@ ";
    mark.style.color = MapView.PLAYER_COLORS[p.slot];
    li.appendChild(mark);
    li.appendChild(document.createTextNode(p.name + (p.id === state.you ? " (du)" : p.ai ? " (AI)" : p.bot ? " (dator)" : "") + (p.klass ? ", " + p.klass : "")));
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
  renderClass(box, line);
  box.appendChild(hpMeter(state.hp, state.maxHp));
  if (me.status === "dod") line("Död", "cond");
  line("Skatt: " + me.treasure + " guld");
  if (state.potions > 0) line("Drycker: " + state.potions);
  (state.foremal || []).forEach((it) => {
    const d = document.createElement("div");
    d.className = "item";
    d.textContent = "· " + it.namn + (it.kvar > 1 ? " (" + it.kvar + ")" : "");
    d.title = it.beskrivning;
    box.appendChild(d);
  });
  if (state.kompass) line(state.kompass, "dimline");
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
      if (!touch) line("A anfall · P parera · U undan", "dimline");
    }
    const others = foe || state.cellar > 0 ? [] : foesHere();
    if (others.length) {
      const names = others.map((p) => p.name).join(" och ");
      const f = currentFoe();
      if (state.pvp === "ja") {
        line("Strid: " + f.name + (others.length > 1 ? " (" + keyName("Tab") + " byter)" : ""), "foe");
        if (!touch) line("A anfall · P parera · U undan", "dimline");
      } else {
        line("Här står också " + names, "dimline");
        if (state.pvp === "fredat") line("Porten är fredad", "dimline");
        if (state.pvp === "for_tidigt") line("Ingen strid förrän solen sjunkit halvvägs", "dimline");
      }
    }
  }
  if (state.controller) {
    line("Du styr tornets varelser", "cond");
    const m = controlledMonster();
    if (m) {
      line("Vald: " + m.name + " (" + keyName("Tab") + " byter)", "foe");
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
    line(me.ready ? "Du är redo" : "Tryck " + keyName("Enter") + " när du är redo", "plan");
  } else if (state.phase === "slut") {
    line(me.ready ? "Väntar på de andra" : touch ? "Tryck Nytt torn" : "Enter: nytt torn", "plan");
  }
}

// Namnet på en tangent så som spelaren ser den: tangenten på en dator,
// knappen på en pekskärm.
const TOUCH_NAMES = { Tab: "Byt", Enter: "Redo" };
function keyName(k) {
  return touch ? TOUCH_NAMES[k] || k : k;
}

// Pekknappen för Enter heter olika saker beroende på vad den gör.
function setEnterLabel(text) {
  const b = $("padEnter");
  if (b.textContent !== text) b.textContent = text;
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
  const el = $("solmatare");
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
    el.textContent = (n < state.minPlayers ? "väntar på vandrare" : "väntar på att alla är redo") + " · tornet " + state.size + " × " + state.size;
    el.classList.remove("urgent");
    return;
  }
  if (state.phase === "slut") {
    el.textContent = state.sun > 0 ? "vandringen är slut" : "solen har gått ner";
    el.classList.remove("urgent");
    return;
  }
  if (state.paused) {
    el.textContent = "runda " + state.round + "  pausad";
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

// Rutan som visas när den här datorn har blivit spelledare: härifrån går
// det bara tillbaka till startsidan.
function showHosting(text) {
  const el = $("overlay");
  const words = text.split(" ");
  const lines = [{ text: "DU LEDER TORNET" }, { text: "" }];
  let row = "";
  words.forEach((w) => {
    if ((row + " " + w).trim().length > 30) {
      lines.push({ text: row });
      row = w;
    } else {
      row = (row + " " + w).trim();
    }
  });
  if (row) lines.push({ text: row });
  lines.push({ text: "" }, { text: touch ? "Tryck Tillbaka" : "Enter: tillbaka till början" });
  setEnterLabel("Tillbaka");
  el.innerHTML = "";
  frame(lines, 32).forEach((l) => {
    const d = document.createElement("div");
    d.textContent = l.text;
    el.appendChild(d);
  });
  el.hidden = false;
}

// backToStart går tillbaka till startsidan med tornet i solnedgången, som
// när sidan just har öppnats.
function backToStart() {
  if (ws) {
    ws.onclose = null;
    ws.close();
    ws = null;
  }
  joined = false;
  state = null;
  towerClosed = false;
  hostingHere = false;
  closedText = "tornet är stängt";
  $("conn").textContent = "";
  $("overlay").hidden = true;
  $("log").textContent = "";
  $("nameerr").textContent = "";
  Intro.reset();
  show("start");
  $("enter").focus();
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
    const name = (r.winner ? "*" : " ") + (r.ai ? "~" : " ") + r.name;
    const row = name.padEnd(20, " ") + what.padStart(10, " ");
    lines.push({ text: row, cls: r.winner ? "win" : r.status !== "ute" ? "lost" : "" });
  });
  if (state.results.some((r) => r.ai)) lines.push({ text: "" }, { text: "~ AI spelade en del av vandringen" });
  lines.push({ text: "" }, { text: touch ? "Tryck Nytt torn" : "Enter: nytt torn" });
  el.innerHTML = "";
  frame(lines, 32).forEach((l) => {
    const d = document.createElement("div");
    d.textContent = l.text;
    if (l.cls) d.className = l.cls;
    el.appendChild(d);
  });
  el.hidden = false;
}

// Introt: tornet i en storm vid solnedgången, med rykten som tonar fram.
// Ljudet får bara starta efter ett klick eller en tangent. "Gå in i tornet"
// visas först efter en stund, utom för den som redan har sett introt i
// fliken och på tornets egen dator.
const Intro = (() => {
  const reduced = window.matchMedia && matchMedia("(prefers-reduced-motion: reduce)").matches;
  let seen = false; // har sett hela övergången förut i den här fliken
  try { seen = sessionStorage.getItem("ormtornet-intro") === "1"; } catch {}
  let rumors = [];
  // På tornets egen dator leder introt till värdsidan, utom när man
  // har valt att spela själv (adressen slutar med ?spela).
  let toHostPage = false;
  const wantsToPlay = new URLSearchParams(location.search).has("spela");
  let reveal = true;
  let next = 0;
  let entering = false;
  let done = false;
  let timer = null;
  let ready = false; // knappen syns och går att använda
  let heard = false; // besökaren har startat stormens ljud
  let flashTimer = null;
  let readyTimer = null;

  function showRumor() {
    if (!rumors.length || entering) return;
    const el = $("rumor");
    el.classList.add("fade");
    setTimeout(() => {
      el.textContent = rumors[next++ % rumors.length];
      el.classList.remove("fade");
    }, el.textContent ? 900 : 0);
  }

  // Regnet ritas på en egen duk: sneda streck som faller i vinden.
  const Rain = (() => {
    const canvas = $("rain");
    const c = canvas.getContext("2d");
    let drops = [];
    let w = 0;
    let h = 0;
    let running = false;
    let last = 0;

    function resize() {
      const dpr = window.devicePixelRatio || 1;
      w = canvas.clientWidth;
      h = canvas.clientHeight;
      canvas.width = Math.round(w * dpr);
      canvas.height = Math.round(h * dpr);
      c.setTransform(dpr, 0, 0, dpr, 0, 0);
      const n = Math.min(220, Math.round((w * h) / (reduced ? 30000 : 9000)));
      drops = Array.from({ length: n }, () => drop(Math.random() * h));
    }

    function drop(y) {
      const near = Math.random();
      return {
        x: Math.random() * (w + 200) - 100,
        y,
        len: 8 + near * 16,
        speed: 420 + near * 520, // bildpunkter per sekund
        alpha: 0.12 + near * 0.28,
      };
    }

    function frame(t) {
      if (!running) return;
      if (canvas.clientWidth !== w || canvas.clientHeight !== h) resize();
      const dt = Math.min(0.05, (t - last) / 1000 || 0);
      last = t;
      draw(dt);
      requestAnimationFrame(frame);
    }

    function draw(dt) {
      c.clearRect(0, 0, w, h);
      c.lineWidth = 1;
      drops.forEach((d) => {
        d.y += d.speed * dt;
        d.x -= d.speed * dt * 0.28;
        if (d.y > h + d.len) Object.assign(d, drop(-d.len));
        c.strokeStyle = "rgba(160, 175, 210, " + d.alpha + ")";
        c.beginPath();
        c.moveTo(d.x, d.y);
        c.lineTo(d.x + d.len * 0.28, d.y - d.len);
        c.stroke();
      });
    }

    return {
      start() {
        resize();
        if (reduced) {
          draw(0); // en stilla bild av regnet
          return;
        }
        if (running) return;
        running = true;
        last = performance.now();
        requestAnimationFrame(frame);
      },
      stop() { running = false; },
      resize() { if (drops.length) resize(); if (reduced) draw(0); },
    };
  })();
  window.addEventListener("resize", () => Rain.resize());

  // Blixten: himlen blir vit, tornet en svart siluett och sedan mullrar det.
  function flash(after) {
    if (entering) return;
    const scene = $("scene");
    const fl = $("flash");
    const steps = reduced
      ? [[0, 0.22, true], [900, 0, false]]
      : [[0, 0.7, true], [70, 0.1, true], [140, 0.45, true], [230, 0, false], [420, 0.25, true], [480, 0, false]];
    fl.style.transition = reduced ? "opacity 0.8s" : "none";
    steps.forEach(([ms, op, lit]) => setTimeout(() => {
      fl.style.opacity = op;
      scene.classList.toggle("lit", lit);
    }, ms));
    const near = Math.random();
    if (heard) setTimeout(() => Sound.thunder(near), 300 + (1 - near) * 2600);
    // Ibland glöder ögat i tornet en stund efter blixten.
    if (reveal && Math.random() < 0.35) {
      setTimeout(() => $("eye").classList.add("glint"), 600);
      setTimeout(() => $("eye").classList.remove("glint"), 1900);
    }
    if (after) setTimeout(after, reduced ? 700 : 480); // när blixten har slocknat
  }

  function scheduleFlash() {
    clearTimeout(flashTimer);
    flashTimer = setTimeout(() => {
      flash();
      scheduleFlash();
    }, 8000 + Math.random() * 12000);
  }

  // Knappen tonar fram, helst precis efter en blixt, och porten glimmar.
  function showEnter(withFlash) {
    if (ready) return;
    const appear = () => {
      ready = true;
      const b = $("enter");
      b.classList.remove("waiting");
      b.classList.add("appear");
      $("gate").classList.add("glow");
      if (Intro.active()) b.focus();
    };
    if (withFlash) {
      // Blixten kommer precis innan knappen, och nästa blixt dröjer en stund.
      flash(appear);
      scheduleFlash();
    } else {
      appear();
    }
  }

  function listen() {
    if (heard) return;
    heard = true;
    Sound.start();
    Sound.storm(reveal);
    $("listen").classList.add("gone");
  }

  function startStorm() {
    Rain.start();
    scheduleFlash();
    if (heard) Sound.storm(reveal);
  }

  fetch("/intro" + (document.visibilityState === "visible" ? "?synlig" : ""), { cache: "no-store" })
    .then((r) => r.json())
    .then((info) => {
      reveal = !!info.avslojande;
      if (!reveal) $("intro").classList.add("invigd");
      toHostPage = info.vard && !wantsToPlay;
      $("listen").textContent = info.ljud || "";
      rumors = (info.rykten || []).slice().sort(() => Math.random() - 0.5);
      showRumor();
      setInterval(showRumor, 7500);
      if (seen || toHostPage) showEnter(false);
      else readyTimer = setTimeout(() => showEnter(true), (info.vanta || 15) * 1000);
    })
    .catch(() => showEnter(false));
  startStorm();

  // Fönstret blinkar, en enda gång, en stund efter att sidan har öppnats.
  if (!reduced) {
    setTimeout(() => {
      const eye = $("eye");
      eye.classList.add("blink");
      setTimeout(() => eye.classList.remove("blink"), 170);
    }, 7000 + Math.random() * 9000);
  }

  function finish() {
    if (done) return;
    done = true;
    clearTimeout(timer);
    Rain.stop();
    if (toHostPage) {
      location.href = "/spelledare";
      return;
    }
    show("namescreen");
    try { $("name").value = localStorage.getItem("ormtornet-namn") || ""; } catch {}
    $("name").focus();
  }

  function enter() {
    if (entering) return;
    if (!ready) {
      listen(); // innan knappen syns startar ett klick bara ljudet
      return;
    }
    entering = true;
    clearTimeout(flashTimer);
    clearTimeout(readyTimer);
    Sound.start();
    Sound.calm(3);
    Sound.play("intro");
    try { sessionStorage.setItem("ormtornet-intro", "1"); } catch {}
    const quick = seen;
    seen = true;
    $("start").classList.add("leaving");
    $("intro").classList.add(quick ? "quick" : "entering");
    timer = setTimeout(finish, quick ? 600 : 3000);
  }

  return {
    enter,
    // poke: ett klick någonstans på startsidan. Under övergången hoppar det
    // fram, annars startar det stormens ljud.
    poke() {
      if (entering) finish();
      else listen();
    },
    reset() {
      clearTimeout(timer);
      entering = false;
      done = false;
      $("intro").classList.remove("entering", "quick");
      $("start").classList.remove("leaving");
      showEnter(false);
      startStorm();
    },
    active() { return $("start").classList.contains("active"); },
    entering() { return entering; },
  };
})();

// Ett tryck på kartan växlar mellan utsnittet och hela tornet, när tornet
// är för stort för skärmen.
$("map").addEventListener("click", () => MapView.toggleWhole());

$("enter").addEventListener("click", (e) => {
  e.stopPropagation();
  Intro.enter();
});
$("start").addEventListener("click", () => Intro.poke());
document.addEventListener("keydown", (e) => {
  if (!Intro.active() || e.metaKey || e.ctrlKey || e.altKey) return;
  if (Intro.entering()) {
    e.preventDefault();
    Intro.poke();
  } else if (e.key === "Enter" || e.key === " ") {
    e.preventDefault();
    Intro.enter();
  }
});

$("back").addEventListener("click", () => pressKey("Escape"));
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
  if (e.metaKey || e.ctrlKey || e.altKey) return;
  if (pressKey(e.key)) e.preventDefault();
});

// pressKey gör det en tangent gör. Pekknapparna använder samma väg och
// skickar namnet på den tangent de står för. Svarar true om tangenten
// betydde något.
function pressKey(key) {
  // Esc i namnrutan, eller Enter när den här datorn har blivit spelledare,
  // går tillbaka till startsidan.
  if ((key === "Escape" && $("namescreen").classList.contains("active")) || (key === "Enter" && hostingHere)) {
    backToStart();
    return true;
  }
  if (!joined) return false;
  const offer = (state && state.erbjudande) || [];
  if ((key === "1" || key === "2") && offer[+key - 1]) {
    send({ t: "klass", name: offer[+key - 1].id });
    return true;
  }
  if (key === "Enter") {
    send({ t: "ready" });
    return true;
  }
  if (state && state.controller) {
    // Den döda spelaren: Tab väljer varelse, pilarna flyttar den.
    if (key === "Tab") {
      cycleMonster();
      return true;
    }
    const a = KEYS[key];
    if (a !== "vanta" && !DIR_WORDS[a]) return false;
    const m = controlledMonster();
    if (a === "vanta") send({ t: "act", a });
    else if (m) send({ t: "act", a, m: m.id });
    return true;
  }
  if (key === "Tab" && foesHere().length) {
    cycleFoe();
    return true;
  }
  const a = KEYS[key];
  if (!a) return false;
  if (!state || state.phase !== "spel") return true;
  const me = myPlayer();
  const monsterHere = me && (state.monsters || []).some((m) => m.x === me.x && m.y === me.y);
  const foe = currentFoe();
  if (["anfall", "parera", "undan"].includes(a) && !monsterHere && foe) {
    send({ t: "act", a, m: -foe.id }); // negativt id: en annan vandrare
  } else {
    send({ t: "act", a });
  }
  return true;
}

// Pekknapparna: varje knapp står för en tangent (data-key).
$("pad").addEventListener("click", (e) => {
  const b = e.target.closest("button[data-key]");
  if (!b) return;
  pressKey(b.dataset.key);
});

// Pekskärm: knapparna visas från början på en telefon, och annars så fort
// någon rör vid skärmen.
let touch = !!(window.matchMedia && matchMedia("(pointer: coarse)").matches);
function setTouch() {
  document.body.classList.toggle("touch", touch);
  if (state && !hostingHere) {
    renderSide();
    renderOverlay();
  }
}
window.addEventListener("touchstart", () => {
  if (touch) return;
  touch = true;
  setTouch();
}, { passive: true });
setTouch();

// Säg till tornet då och då att sidan syns, så att det inte stänger av sig
// självt medan någon sitter och tittar. En gömd flik säger ingenting.
setInterval(() => {
  if (document.visibilityState !== "visible") return;
  if (joined) send({ t: "puls" });
  else fetch("/intro?synlig", { cache: "no-store" }).catch(() => {});
}, 60000);

// Rita om då och då så att den egna markören blinkar.
setInterval(() => {
  if (state && mapReady) MapView.draw(state);
  if (state) updateTimer();
}, 100);
