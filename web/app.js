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

// Loggen. Raderna samlas i grupper, en per runda (och en för tiden före
// och efter vandringen). Den första raden i en grupp skrivs fram tecken för
// tecken, resten kommer en hel rad i taget i lugn takt. Ligger loggen
// efter skrivs rader från förbigångna rundor ikapp fort. Har man rullat upp
// för att läsa stannar loggen där, och "↓ nya rader" tar en tillbaka. När
// en ny runda börjar fälls rundorna före den förra ihop till en rad som går
// att klicka upp.
const Log = (() => {
  const queue = [];
  let busy = false;
  const MAX_LINES = 300;
  const LINE_GAP = 320; // ms mellan hela rader
  const CHAR_MS = 24;
  const CATCH_UP = 60; // ms mellan rader från en runda som redan är förbi
  let stick = true; // loggen följer med nedåt
  let current = null; // gruppen som raderna hamnar i

  function box() {
    return $("log");
  }

  function keyNow() {
    if (!state || state.phase === "lobby" || state.phase === "klassval") return "lobby";
    if (state.phase === "slut") return "slut";
    return "r" + state.round;
  }

  function label(key) {
    if (key === "lobby") return "Före vandringen";
    if (key === "slut") return "Slutet";
    return "Runda " + key.slice(1);
  }

  function push(text, cls) {
    if (!text) return;
    queue.push({ text, cls, key: keyNow() });
    if (!busy) next();
  }

  // group ger gruppen för key, och gör en ny när rundan har bytts.
  function group(key) {
    if (current && current.dataset.key === key) return current;
    const g = document.createElement("section");
    g.className = "round";
    g.dataset.key = key;
    const head = document.createElement("button");
    head.type = "button";
    head.className = "roundhead";
    head.addEventListener("click", () => {
      g.classList.toggle("folded");
      summarize(g);
    });
    const lines = document.createElement("div");
    lines.className = "lines";
    g.append(head, lines);
    box().appendChild(g);
    current = g;
    summarize(g);
    fold();
    return g;
  }

  // summarize skriver gruppens rubrik: rundan, och när den är hopfälld
  // rummen man kom till och hur många rader den har.
  function summarize(g) {
    const head = g.querySelector(".roundhead");
    const lines = g.querySelectorAll(".lines > div");
    const folded = g.classList.contains("folded");
    let text = (folded ? "▸ " : "▾ ") + label(g.dataset.key);
    if (folded) {
      const rooms = [...g.querySelectorAll(".lines > .rubrik")].map((d) => d.textContent.replace(/\.$/, ""));
      if (rooms.length) text += " · " + rooms.join(", ");
      text += " · " + lines.length + (lines.length === 1 ? " rad" : " rader");
    }
    head.textContent = text;
    head.setAttribute("aria-expanded", String(!folded));
  }

  // fold fäller ihop alla grupper utom de två senaste. Det väntar tills man
  // är längst ner, så att texten man läser inte flyttar sig.
  function fold() {
    if (!stick) return;
    const groups = [...box().querySelectorAll(".round")];
    groups.slice(0, -2).forEach((g) => {
      if (g.dataset.auto) return;
      g.dataset.auto = "1";
      g.classList.add("folded");
      summarize(g);
    });
    trim();
    follow();
  }

  function trim() {
    let n = box().querySelectorAll(".lines > div").length;
    while (n > MAX_LINES && box().children.length > 2) {
      const first = box().firstElementChild;
      n -= first.querySelectorAll(".lines > div").length;
      first.remove();
    }
  }

  // paint skriver raden och färgar vandrarnas namn i deras färg, så att
  // det syns vem raden handlar om. Ett s efter namnet hör till namnet.
  function paint(line, text) {
    const players = (state && state.players) || [];
    const names = players.map((p) => p.name).filter(Boolean).sort((a, b) => b.length - a.length);
    if (!names.length) {
      line.textContent = text;
      return;
    }
    const esc = names.map((n) => n.replace(/[.*+?^${}()|[\]\\]/g, "\\$&"));
    const re = new RegExp("(?<![\\p{L}\\p{N}])(" + esc.join("|") + ")(?=s?(?![\\p{L}\\p{N}]))", "gu");
    line.textContent = "";
    let at = 0;
    for (const m of text.matchAll(re)) {
      if (m.index > at) line.appendChild(document.createTextNode(text.slice(at, m.index)));
      const p = players.find((x) => x.name === m[1]);
      const s = document.createElement("span");
      s.className = "who";
      s.textContent = m[1];
      if (p) s.style.color = MapView.PLAYER_COLORS[p.slot];
      line.appendChild(s);
      at = m.index + m[1].length;
    }
    if (at < text.length) line.appendChild(document.createTextNode(text.slice(at)));
  }

  function atBottom() {
    const b = box();
    return b.scrollTop + b.clientHeight >= b.scrollHeight - 6;
  }

  // follow håller loggen längst ner, eller visar att det finns nya rader.
  function follow() {
    const b = box();
    if (stick) {
      b.scrollTop = b.scrollHeight;
      $("logDown").hidden = true;
    } else {
      $("logDown").hidden = false;
    }
  }

  function next() {
    const item = queue.shift();
    if (!item) {
      busy = false;
      return;
    }
    busy = true;
    const g = group(item.key);
    const list = g.querySelector(".lines");
    const first = list.children.length === 0;
    const line = document.createElement("div");
    if (item.cls) line.className = item.cls;
    list.appendChild(line);
    if (g.classList.contains("folded")) summarize(g);
    // Rader från en runda som redan är förbi kommer ikapp fort, så att
    // loggen inte hamnar efter. Den aktuella rundan kommer i lugn takt.
    const behind = queue.some((q) => q.key !== item.key);
    if (!first || behind) {
      // Hela raden på en gång, med en mjuk övergång.
      paint(line, item.text);
      line.classList.add("appear");
      follow();
      setTimeout(next, behind ? CATCH_UP : LINE_GAP);
      return;
    }
    line.classList.add("cursor");
    const chars = [...item.text];
    let i = 0;
    function tick() {
      i = Math.min(chars.length, i + 1);
      line.textContent = chars.slice(0, i).join("");
      follow();
      if (i < chars.length) {
        const ch = chars[i - 1];
        setTimeout(tick, ".!?".includes(ch) ? 180 : ",".includes(ch) ? 70 : CHAR_MS);
      } else {
        line.classList.remove("cursor");
        paint(line, item.text);
        setTimeout(next, LINE_GAP);
      }
    }
    tick();
  }

  function init() {
    box().addEventListener("scroll", () => {
      const was = stick;
      stick = atBottom();
      if (stick) $("logDown").hidden = true;
      if (stick && !was) fold();
    });
    $("logDown").addEventListener("click", () => {
      stick = true;
      fold();
      follow();
    });
  }

  return { push, init };
})();
Log.init();

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
      }
      showPhase();
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
      renderChoose();
      firstHelp();
      checkRoundCard();
      restoreTip();
      break;
    case "log":
      Log.push(msg.text, msg.c);
      break;
    case "runda":
      showRoundCard(msg);
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
function actionWord(a) {
  if (/^f\d/.test(a)) {
    const it = (state.foremal || [])[+a[1] - 1];
    return "använd " + (it ? it.namn.toLowerCase() : "föremål") + (a[2] ? " " + DIR_WORDS[a[2]] : "");
  }
  return ACTION_WORDS[a] || a;
}
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
// En mätare: en etikett, en rad rutor (en per enhet) och ett tal. Rutorna
// byggs i kod, så att mätaren alltid har samma bredd.
function meter(cls, label, value, maxValue, num) {
  const d = document.createElement("div");
  d.className = "meter " + cls;
  const l = document.createElement("span");
  l.className = "mlabel";
  l.textContent = label;
  const segs = document.createElement("span");
  segs.className = "segs";
  for (let i = 0; i < maxValue; i++) {
    const seg = document.createElement("span");
    seg.className = i < value ? "seg on" : "seg";
    segs.appendChild(seg);
  }
  const n = document.createElement("span");
  n.className = "num";
  n.textContent = num !== undefined ? num : Math.max(0, value) + "/" + maxValue;
  d.append(l, segs, n);
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

// Klassvalet har en egen skärm mellan lobbyn och vandringen. Kartan ritas
// om när spelskärmen syns igen, eftersom den inte har någon storlek medan
// den är dold.
function showPhase() {
  const want = state.phase === "klassval" ? "classscreen" : "game";
  if ($(want).classList.contains("active")) return;
  show(want);
  if (want === "game") MapView.resize();
}

// Korten: ett för varje klass. Ett kort som någon annan har tagit går inte
// att välja, och när man själv har valt ligger de andra kvar som de är.
function renderChoose() {
  if (state.phase !== "klassval") return;
  const cards = state.erbjudande || [];
  const box = $("cards");
  const mine = cards.find((c) => c.din);
  box.innerHTML = "";
  cards.forEach((c, i) => {
    const b = document.createElement("button");
    b.type = "button";
    b.className = "card" + (c.din ? " mine" : "");
    b.setAttribute("role", "listitem");
    b.disabled = !!c.tagen || !!mine;
    const add = (cls, text) => {
      const d = document.createElement("div");
      d.className = cls;
      d.textContent = text;
      b.appendChild(d);
      return d;
    };
    add("num", String(i + 1));
    add("num low", String(i + 1));
    add("cname", (c.tecken ? c.tecken + " " : "") + c.namn);
    add("desc", c.beskrivning);
    const stats = [];
    if (c.liv) stats.push((c.liv > 0 ? "+" : "−") + Math.abs(c.liv) + " liv");
    if (c.drycker) stats.push(c.drycker + (c.drycker === 1 ? " dryck" : " drycker") + " från början");
    if (stats.length) add("stats", stats.join(" · "));
    add("good", "+ " + c.styrka);
    add("bad", "− " + c.svaghet);
    add("owner", c.din ? "ditt kort" : c.tagen ? "tagen av " + c.tagen : "");
    b.addEventListener("click", () => chooseCard(i));
    box.appendChild(b);
  });
  updateChooseStatus();
}

function chooseCard(i) {
  const c = state && state.phase === "klassval" && (state.erbjudande || [])[i];
  if (!c || c.tagen || state.erbjudande.some((x) => x.din)) return false;
  send({ t: "klass", name: c.id });
  return true;
}

function updateChooseStatus() {
  if (!state || state.phase !== "klassval") return;
  const el = $("chooseStatus");
  const secs = Math.ceil(remainingMs() / 1000);
  const mine = (state.erbjudande || []).find((c) => c.din);
  const waiting = state.players.filter((p) => p.online && !p.klass && p.id !== state.you).map((p) => p.name);
  let text = mine ? "Du är " + mine.namn.toLowerCase() + ". " : touch ? "Tryck på ett kort. " : "Tryck 1–" + state.erbjudande.length + " eller klicka på ett kort. ";
  if (mine && waiting.length) text += "Väntar på " + waiting.join(", ") + ". ";
  el.textContent = text;
  const s = document.createElement("span");
  s.className = "secs" + (secs <= 5 ? " urgent" : "");
  s.textContent = secs + " s";
  el.appendChild(s);
}

function renderSide() {
  setEnterLabel(state.phase === "slut" ? "Nytt torn" : "Redo");
  // Under vandringen står vandrarna som prickar i översta raden.
  $("playersBlock").hidden = state.phase === "spel";
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
  $("roomname").textContent = here;
  renderHud(me);
  renderNow(me);
}

// Statusen står som tecken uppe till höger på kartan: klassen, liv, skatt,
// drycker, tillstånd och föremålen. Varje tecken har en förklaring som syns
// när musen är över det, eller efter ett tryck på en pekskärm.
function renderHud(me) {
  const box = $("status");
  const bag = $("bag");
  box.innerHTML = "";
  bag.innerHTML = "";
  if (!me) return;
  const icon = (glyph, value, tip, cls) => {
    const d = document.createElement("span");
    d.className = "hudicon" + (cls ? " " + cls : "");
    d.dataset.tip = tip;
    d.tabIndex = 0;
    const g = document.createElement("span");
    g.className = "glyph";
    g.textContent = glyph;
    d.appendChild(g);
    if (value !== "") {
      const v = document.createElement("span");
      v.className = "val";
      v.textContent = value;
      d.appendChild(v);
    }
    d.setAttribute("aria-label", tip);
    box.appendChild(d);
  };
  const c = state.klass;
  if (c) icon(c.tecken || c.namn[0], "", c.namn + "\n" + c.beskrivning + "\n+ " + c.styrka + "\n− " + c.svaghet, "klass");
  if (me.status === "dod") {
    icon("✕", "", "Du är död.", "warn");
  } else {
    const low = state.hp <= Math.max(1, Math.floor(state.maxHp / 3));
    icon("♥", state.hp + "/" + state.maxHp, "Liv: " + state.hp + " av " + state.maxHp + ". Vid noll är du död.", low ? "warn" : "hp");
  }
  icon("¤", String(me.treasure), "Skatt: " + me.treasure + " guld. Den räknas bara om du bär ut den genom en port.", "gold");
  if (state.potions > 0) {
    icon("⚱", String(state.potions), (state.potions === 1 ? "En dryck. " : state.potions + " drycker. ") +
      (touch ? "Tryck Drick" : "Tryck D") + " för att dricka en. Den ger tillbaka liv och tar bort gift.", "");
  }
  if (me.status === "inne") {
    if (state.cellar > 0) icon("▼", "", "I källarvalven. Ingen ser var du är, och du kommer upp någon annanstans.", "warn");
    if (state.stuck > 0) icon("⛓", String(state.stuck), "Fast i " + state.stuck + (state.stuck === 1 ? " runda" : " rundor") + ". Du kan inte gå härifrån.", "warn");
    if (state.poison > 0) icon("☠", "", "Förgiftad. Du tappar liv varje runda tills giftet går ur. En dryck tar bort det.", "warn");
    if (state.dark > 0) icon("◐", "", "Lyktan är släckt. Du ser bara rummet du står i.", "warn");
  }
  if (state.kompass) icon("✧", "", state.kompass, "");
  if (state.hjartslag) icon("♡", "", state.hjartslag, "");
  if (state.invigd) icon("❝", "", "Tornet har hört dig.", "rost");
  if (state.kanViska) icon("V", "", (touch ? "Tryck Viska" : "Tryck V") + " för att viska till tornet.", "rost");
  const pw = $("padWhisper");
  if (pw) pw.hidden = !state.kanViska;

  (state.foremal || []).forEach((it, i) => {
    const b = document.createElement("button");
    b.dataset.slot = "f" + i;
    b.type = "button";
    b.className = "hudicon item" + (it.aktiv ? " active" : "") + (aimSlot === i ? " aiming" : "");
    b.textContent = (it.aktiv && !touch ? i + 1 + " " : "") + it.namn + (it.kvar > 1 ? " ×" + it.kvar : "");
    b.dataset.tip = it.beskrivning + (it.aktiv ? (touch ? "" : " Tryck " + (i + 1) + " för att använda det.") : "");
    if (it.aktiv) b.addEventListener("click", () => useItemSlot(i));
    bag.appendChild(b);
  });
}

// Ett tecken känns igen på sitt glyftecken, eller på platsen i väskan.
let tipKey = "";
function keyOf(el) {
  return el.dataset.slot || (el.querySelector(".glyph") || el).textContent;
}

// Panelen Nu: det som gäller just nu, som draget, striden och vad som väntar.
function renderNow(me) {
  const box = $("me");
  box.innerHTML = "";
  const line = (text, cls) => {
    const d = document.createElement("div");
    d.textContent = text;
    if (cls) d.className = cls;
    box.appendChild(d);
  };
  if (!me) return;
  if (aimSlot !== null && state.foremal && state.foremal[aimSlot]) {
    line(state.foremal[aimSlot].namn + ": åt vilket håll? (pilarna" + (touch ? "" : ", esc ångrar") + ")", "plan");
  }
  if (me.status === "dod" && !state.controller) line("Du är död", "cond");
  if (me.status === "inne") {
    if (state.sleep >= 0 && state.sleepText) line(state.sleepText, state.sleep >= 2 ? "cond" : "sleep");
    const foe = (state.monsters || []).find((m) => m.x === me.x && m.y === me.y);
    if (foe && state.cellar === 0) {
      line("Strid: " + foe.name, "foe");
      box.appendChild(meter("foehp", "Dess liv", foe.hp, foe.maxHp));
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
    line(state.action ? "Drag: " + actionWord(state.action) : "Välj ditt drag", "plan");
  } else if (state.phase === "lobby") {
    line(me.ready ? "Du är redo. Väntar på de andra" : "Tryck " + keyName("Enter") + " när du är redo", "plan");
  } else if (state.phase === "slut") {
    line(me.ready ? "Väntar på de andra" : touch ? "Tryck Nytt torn" : "Enter: nytt torn", "plan");
  }
}

// Rundans sammanfattning: ett kort över kartan med vad som hände dig, det
// farliga först, och en kort rad i loggen. Kortet försvinner när du väljer
// nästa drag, efter en stund, eller vid ett klick.
const CARD_MS = 7000;
let cardTimer = null;
let cardAction = "";
function showRoundCard(msg) {
  const lines = msg.lines || [];
  if (!lines.length) return;
  const short = lines.filter((l) => l.kort).map((l) => l.kort);
  if (short.length) Log.push("Runda " + msg.round + ": " + short.join(", ") + ".", "sam");
  const card = $("roundcard");
  card.innerHTML = "";
  const h = document.createElement("div");
  h.className = "cardhead";
  h.textContent = "Runda " + msg.round;
  card.appendChild(h);
  lines.forEach((l) => {
    const d = document.createElement("div");
    d.textContent = l.text;
    if (l.c) d.className = l.c;
    card.appendChild(d);
  });
  card.hidden = false;
  card.classList.remove("leaving");
  cardAction = "";
  clearTimeout(cardTimer);
  cardTimer = setTimeout(hideRoundCard, CARD_MS);
}
function hideRoundCard() {
  const card = $("roundcard");
  if (card.hidden || card.classList.contains("leaving")) return;
  clearTimeout(cardTimer);
  card.classList.add("leaving");
  setTimeout(() => {
    if (card.classList.contains("leaving")) card.hidden = true;
  }, 400);
}
// checkRoundCard döljer kortet när ett nytt drag är valt eller vandringen
// är slut.
function checkRoundCard() {
  if ($("roundcard").hidden) return;
  if (state.phase !== "spel") {
    $("roundcard").hidden = true;
    return;
  }
  if (state.action && state.action !== cardAction) hideRoundCard();
  cardAction = state.action || "";
}
$("roundcard").addEventListener("click", hideRoundCard);

// Teckenförklaringen uppe till vänster på kartan: de tecken som står i
// rummen just nu. Texterna kommer från tecken_* i data/texter.json, och
// fällor och varelser förklaras med sina namn. Den går att fälla ihop.
let signTexts = {};
let legendOpen = !(window.matchMedia && matchMedia("(max-width: 700px)").matches);
try {
  const saved = localStorage.getItem("ormtornet-tecken");
  if (saved) legendOpen = saved === "1";
} catch {}
function renderLegend(list) {
  const box = $("legendList");
  box.innerHTML = "";
  list.forEach((e) => {
    const text = e.name || signTexts[e.key.split(":")[0]];
    if (!text) return;
    const row = document.createElement("div");
    row.className = "sign";
    const g = document.createElement("span");
    g.className = "glyph";
    g.textContent = e.glyph;
    g.style.color = e.color;
    const w = document.createElement("span");
    w.className = "what";
    w.textContent = text;
    row.append(g, w);
    box.appendChild(row);
  });
  $("legend").hidden = box.childElementCount === 0;
  showLegendOpen();
}
function showLegendOpen() {
  $("legendList").hidden = !legendOpen;
  $("legendBtn").setAttribute("aria-expanded", String(legendOpen));
}
$("legendBtn").addEventListener("click", () => {
  legendOpen = !legendOpen;
  try { localStorage.setItem("ormtornet-tecken", legendOpen ? "1" : "0"); } catch {}
  showLegendOpen();
});
MapView.setLegend(renderLegend);

// Förklaringen till ett tecken på kartan. Den läggs under tecknet, inom
// kartans ruta.
function showTip(el) {
  const tip = $("tip");
  const r = el.getBoundingClientRect();
  tip.textContent = el.dataset.tip;
  tip.hidden = false;
  tipKey = keyOf(el);
  const w = tip.offsetWidth;
  // Under tecknet, inom fönstret. Tecknen på kartan står till höger, så
  // förklaringen slutar där tecknet slutar; prickarna börjar där de börjar.
  const left = el.closest("#crew") ? r.left : r.right - w;
  tip.style.top = r.bottom + 6 + "px";
  tip.style.left = Math.max(4, Math.min(window.innerWidth - w - 4, left)) + "px";
}
function hideTip() {
  $("tip").hidden = true;
  tipKey = "";
}
// Tecknen på kartan och prickarna i översta raden har samma förklaringar.
for (const id of ["hud", "crew"]) {
  const area = $(id);
  area.addEventListener("mouseover", (e) => {
    if (touch) return; // på en pekskärm är det trycket som räknas
    const el = e.target.closest("[data-tip]");
    if (el) showTip(el);
  });
  area.addEventListener("mouseout", (e) => {
    if (touch) return;
    if (!e.relatedTarget || !e.relatedTarget.closest || !e.relatedTarget.closest("[data-tip]")) hideTip();
  });
  area.addEventListener("focusin", (e) => {
    if (touch) return;
    const el = e.target.closest("[data-tip]");
    if (el) showTip(el);
  });
  area.addEventListener("focusout", () => {
    if (!touch) hideTip();
  });
  area.addEventListener("click", (e) => {
    // På en pekskärm visar ett tryck förklaringen, och ett tryck till döljer den.
    const el = e.target.closest("[data-tip]:not(button)");
    if (!el) return;
    e.stopPropagation();
    if (!$("tip").hidden && tipKey === keyOf(el)) hideTip();
    else showTip(el);
  });
}
// Ett tryck någon annanstans döljer förklaringen.
document.addEventListener("click", (e) => {
  if (touch && tipKey && !e.target.closest("#hud, #crew")) hideTip();
});

// restoreTip visar förklaringen igen efter att läget har ritats om.
function restoreTip() {
  if (!tipKey) return;
  const el = [...document.querySelectorAll("#hud [data-tip], #crew [data-tip]")].find((x) => keyOf(x) === tipKey);
  if (el) showTip(el);
  else hideTip();
}

// Styrningen ligger bakom frågetecknet. Första vandringen visas den en gång.
function toggleHelp(open) {
  const box = $("helpbox");
  const show = open === undefined ? box.hidden : open;
  box.hidden = !show;
  $("helpBtn").setAttribute("aria-expanded", String(show));
  if (show) {
    try { localStorage.setItem("ormtornet-styrning", "1"); } catch {}
  }
}
$("helpBtn").addEventListener("click", (e) => {
  e.stopPropagation();
  toggleHelp();
});
$("helpbox").addEventListener("click", (e) => {
  e.stopPropagation();
  toggleHelp(false);
});
let helpOffered = false;
function firstHelp() {
  if (touch || helpOffered || state.phase !== "spel") return;
  helpOffered = true;
  let seen = false;
  try { seen = localStorage.getItem("ormtornet-styrning") === "1"; } catch {}
  if (!seen) toggleHelp(true);
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
  return meter("hp", "Liv", hp, maxHp);
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
  const m = meter("sun", "Sol", filled, SUN_WIDTH, String(state.sun));
  if (state.sun <= 3) m.classList.add("last");
  else if (state.sun <= state.sunMax / 4) m.classList.add("low");
  el.appendChild(m);
  renderCrew();
  updateTimer();
}

// Vandrarna som prickar i sin färg i översta raden. Under vandringen är
// pricken fylld när vandraren har valt drag, i lobbyn när den är redo.
// Förklaringen säger vem det är, klassen och läget.
function renderCrew() {
  const box = $("crew");
  box.innerHTML = "";
  state.players.forEach((p) => {
    const d = document.createElement("span");
    d.className = "dot" + (p.id === state.you ? " me" : "") + (p.online ? "" : " off");
    d.dataset.slot = "p" + p.id;
    d.tabIndex = 0;
    let glyph;
    if (p.status === "dod") glyph = "✕";
    else if (p.status === "ute") glyph = "◌";
    else if (p.status === "forlorad") glyph = "✕";
    else glyph = (state.phase === "spel" ? p.chosen : p.ready) ? "●" : "○";
    d.textContent = glyph;
    d.style.color = MapView.PLAYER_COLORS[p.slot];
    const who = p.name + (p.id === state.you ? " (du)" : p.ai ? " (AI spelar)" : p.bot ? " (dator)" : "");
    const tag = playerTag(p);
    d.dataset.tip = who + (p.klass ? ", " + p.klass : "") + (tag ? "\n" + tag : "");
    d.setAttribute("aria-label", d.dataset.tip);
    box.appendChild(d);
  });
}

function remainingMs() {
  if (!state || (state.phase !== "spel" && state.phase !== "klassval")) return 0;
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
  if (state.phase === "klassval") {
    el.textContent = "klassval";
    updateChooseStatus();
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
      signTexts = info.tecken || {};
      $("legendBtn").textContent = signTexts.rubrik || "";
      MapView.setLegend(renderLegend); // texterna kom efter kartan
      if (state && mapReady) MapView.draw(state);
      rumors = (info.rykten || []).slice().sort(() => Math.random() - 0.5);
      showRumor();
      setInterval(showRumor, 7500);
      if (seen || toHostPage) showEnter(false);
      else readyTimer = setTimeout(() => showEnter(true), (info.vanta || 5) * 1000);
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
  if (e.target && e.target.id === "whisperText") return; // viskningen skrivs i sin egen ruta
  if (pressKey(e.key)) e.preventDefault();
});

// Föremålen: 1–3 använder föremålet på den platsen i väskan. Ett föremål som
// pekar åt ett håll väntar på en pil.
let aimSlot = null;
function useItemSlot(i) {
  const it = state && (state.foremal || [])[i];
  if (!it || !it.aktiv || state.phase !== "spel") return;
  if (it.riktning) {
    aimSlot = i;
    renderSide();
    return;
  }
  send({ t: "act", a: "f" + (i + 1) });
}

// Viskningen: en egen ruta under panelen, så att bokstäverna inte blir drag.
function openWhisper() {
  const f = $("whisper");
  f.hidden = false;
  $("whisperText").value = "";
  $("whisperText").focus();
}
function closeWhisper() {
  $("whisper").hidden = true;
  $("whisperText").blur();
}
$("whisper").addEventListener("submit", (e) => {
  e.preventDefault();
  const text = $("whisperText").value.trim();
  if (text) send({ t: "viska", name: text });
  closeWhisper();
});
$("whisperText").addEventListener("keydown", (e) => {
  if (e.key === "Escape") {
    e.preventDefault();
    closeWhisper();
  }
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
  if (key === "?" || ((key === "h" || key === "H") && state && state.phase !== "klassval")) {
    toggleHelp();
    return true;
  }
  if (key === "Escape" && !$("helpbox").hidden) {
    toggleHelp(false);
    return true;
  }
  if (aimSlot !== null) {
    // Ett föremål som pekar åt ett håll: nästa pil väljer riktningen.
    const dir = { ArrowUp: "N", ArrowRight: "E", ArrowDown: "S", ArrowLeft: "W" }[key];
    if (dir) {
      send({ t: "act", a: "f" + (aimSlot + 1) + dir });
      aimSlot = null;
      renderSide();
      return true;
    }
    if (key === "Escape") {
      aimSlot = null;
      renderSide();
      return true;
    }
  }
  if (state && state.phase === "spel" && /^[1-9]$/.test(key)) {
    useItemSlot(+key - 1);
    return true;
  }
  if (state && state.phase === "klassval") {
    // Siffrorna väljer kort. Tab och Enter får göra som vanligt bland knapparna.
    if (/^[1-9]$/.test(key)) {
      chooseCard(+key - 1);
      return true;
    }
    return false;
  }
  if ((key === "v" || key === "V") && state && state.kanViska) {
    openWhisper();
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
