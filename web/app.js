// Klienten: visar det servern skickar och skickar tangenttryck. Ingen spellogik här.
"use strict";

const $ = (id) => document.getElementById(id);

let ws = null;
let state = null;
let myName = "";
let savedName = ""; // namnet som senast sparades i fliken och webbläsaren

// Namnet tornet gav sparas med tiden det senast användes, i fliken och i
// webbläsaren. Efter 12 timmar glöms det, så att varje kväll ger ett nytt
// namn men samma kväll behåller man sitt.
const NAME_KEY = "ormtornet-namn";
const NAME_MS = 12 * 60 * 60 * 1000;
function loadName() {
  for (const store of [sessionStorage, localStorage]) {
    try {
      const v = JSON.parse(store.getItem(NAME_KEY) || "null");
      if (v && typeof v.namn === "string" && Date.now() - v.tid < NAME_MS) return v.namn;
    } catch {}
  }
  return "";
}
function saveName(name) {
  const v = JSON.stringify({ namn: name, tid: Date.now() });
  try {
    sessionStorage.setItem(NAME_KEY, v);
    localStorage.setItem(NAME_KEY, v);
  } catch {}
}
let joined = false;
let mapReady = false;
let mapSize = 0; // tornets storlek när kartan ritades upp
let towerClosed = false; // spelledaren har stängt tornet

function show(id) {
  document.querySelectorAll(".screen").forEach((s) => s.classList.toggle("active", s.id === id));
}

// Loggen. Raderna samlas i grupper, en per runda (och en för tiden före
// och efter vandringen). Den senaste rundan står överst, och inom en runda
// står raderna i den ordning de hände. Den första raden i en grupp skrivs
// fram tecken för tecken, resten kommer en hel rad i taget i lugn takt.
// Ligger loggen efter skrivs rader från förbigångna rundor ikapp fort. Har
// man rullat ner för att läsa ligger texten kvar där, och "↑ nya rader" tar
// en tillbaka upp. Bara den senaste rundan står i bild; de äldre ligger
// under strecket "tidigare" och syns när man rullar ner.
const Log = (() => {
  const queue = [];
  let busy = false;
  const MAX_LINES = 300;
  const LINE_GAP = 320; // ms mellan hela rader
  const CHAR_MS = 24;
  const CATCH_UP = 60; // ms mellan rader från en runda som redan är förbi
  let stick = true; // loggen står överst, där det nya kommer

  function box() {
    return $("log");
  }

  function keyNow() {
    if (!state || state.phase === "lobby" || state.phase === "klassval") return "lobby";
    if (state.phase === "slut") return "slut";
    return "r" + state.round;
  }

  function push(text, cls) {
    if (!text) return;
    queue.push({ text, cls, key: keyNow() });
    if (!busy) next();
  }

  // room lägger rumsblocket (here.go) överst i rundan: rummets namn,
  // beskrivningen när rummet är nytt och vad som finns där. Ett nytt block
  // i samma runda ersätter det förra.
  function room(msg) {
    queue.push({ room: msg, key: keyNow() });
    if (!busy) next();
  }

  function roomBlock(g, msg) {
    const list = g.querySelector(".lines");
    const block = document.createElement("div");
    block.className = "roomblock" + (msg.kvar ? " kvar" : "");
    const name = document.createElement("div");
    name.className = "roomname";
    name.textContent = msg.namn;
    block.appendChild(name);
    let desc = null;
    if (msg.text) {
      desc = document.createElement("div");
      desc.className = "roomdesc";
      desc.textContent = msg.text;
      block.appendChild(desc);
      // På en telefon visar beskrivningen tre rader; "mer" visar resten.
      const more = document.createElement("button");
      more.type = "button";
      more.className = "roommore";
      more.textContent = t("logg_mer");
      more.addEventListener("click", () => block.classList.add("open"));
      block.appendChild(more);
    }
    if (msg.har) {
      const h = document.createElement("div");
      h.className = "roomhere";
      paint(h, msg.har);
      block.appendChild(h);
    }
    block.classList.add("appear");
    keep(() => {
      const old = list.querySelector(":scope > .roomblock");
      if (old) old.replaceWith(block);
      else list.prepend(block);
    });
    if (desc && desc.scrollHeight > desc.clientHeight + 1) block.classList.add("lang");
    follow();
  }

  // keep gör en ändring utan att texten man läser flyttar sig när man har
  // rullat ner: raden överst i rutan står kvar på samma plats.
  function keep(change) {
    const b = box();
    if (stick || b.scrollTop <= 0) {
      change();
      return;
    }
    const top = b.getBoundingClientRect().top;
    const anchor = [...b.querySelectorAll(".roundhead, .lines > div")].find(
      (el) => el.offsetParent && el.getBoundingClientRect().bottom > top
    );
    const before = anchor ? anchor.getBoundingClientRect().top : 0;
    change();
    if (anchor && anchor.isConnected) b.scrollTop += anchor.getBoundingClientRect().top - before;
  }

  // group ger gruppen för key, och gör en ny överst när rundan har bytts.
  // Den översta rundan fyller hela loggens höjd, så att de äldre hamnar
  // under strecket "tidigare" och bara syns när man rullar ner.
  function group(key) {
    const first = box().firstElementChild;
    if (first && first.dataset.key === key) return first;
    const g = document.createElement("section");
    g.className = "round";
    g.dataset.key = key;
    const head = document.createElement("div");
    head.className = "roundhead";
    const lines = document.createElement("div");
    lines.className = "lines";
    g.append(head, lines);
    label(g);
    keep(() => box().prepend(g));
    trim();
    follow();
    return g;
  }

  // label skriver gruppens rubrik och strecket som står ovanför den när den
  // inte längre är den översta.
  function label(g) {
    const key = g.dataset.key;
    let text;
    if (key === "lobby") text = t("logg_fore");
    else if (key === "slut") text = t("logg_slut");
    else text = t("logg_runda", { runda: key.slice(1) });
    g.querySelector(".roundhead").textContent = text;
    g.dataset.earlier = t("logg_tidigare");
  }

  function trim() {
    let n = box().querySelectorAll(".lines > div").length;
    while (n > MAX_LINES && box().children.length > 2) {
      const last = box().lastElementChild;
      n -= last.querySelectorAll(".lines > div").length;
      last.remove();
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

  function atTop() {
    return box().scrollTop <= 6;
  }

  // follow håller loggen överst, eller visar att det finns nya rader.
  function follow() {
    if (stick) {
      box().scrollTop = 0;
      $("logNew").hidden = true;
    } else {
      $("logNew").hidden = false;
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
    if (item.room) {
      roomBlock(g, item.room);
      setTimeout(next, queue.some((q) => q.key !== item.key) ? CATCH_UP : LINE_GAP);
      return;
    }
    const list = g.querySelector(".lines");
    const first = list.children.length === 0;
    const line = document.createElement("div");
    if (item.cls) line.className = item.cls;
    keep(() => list.appendChild(line));
    // Rader från en runda som redan är förbi kommer ikapp fort, så att
    // loggen inte hamnar efter. Den aktuella rundan kommer i lugn takt.
    const behind = queue.some((q) => q.key !== item.key);
    if (!first || behind) {
      // Hela raden på en gång, med en mjuk övergång.
      keep(() => paint(line, item.text));
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
      keep(() => (line.textContent = chars.slice(0, i).join("")));
      follow();
      if (i < chars.length) {
        const ch = chars[i - 1];
        setTimeout(tick, ".!?".includes(ch) ? 180 : ",".includes(ch) ? 70 : CHAR_MS);
      } else {
        line.classList.remove("cursor");
        keep(() => paint(line, item.text));
        setTimeout(next, LINE_GAP);
      }
    }
    tick();
  }

  function init() {
    box().addEventListener("scroll", () => {
      stick = atTop();
      if (stick) $("logNew").hidden = true;
    });
    $("logNew").addEventListener("click", () => {
      stick = true;
      follow();
    });
  }

  // relabel skriver om gruppernas rubriker, när språket har bytts.
  function relabel() {
    box().querySelectorAll(".round").forEach(label);
  }

  return { push, room, init, relabel };
})();
Log.init();

function send(obj) {
  if (ws && ws.readyState === WebSocket.OPEN) ws.send(JSON.stringify(obj));
}

function connect() {
  const proto = location.protocol === "https:" ? "wss:" : "ws:";
  ws = new WebSocket(proto + "//" + location.host + "/ws");
  savedName = ""; // spara namnet igen, med ny tid, när tornet svarar
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
    if (!joined) {
      // Tornet svarade inte: säg det, med en knapp tillbaka.
      $("nameerr").textContent = t("anslutning_ingen");
      show("namescreen");
      return;
    }
    if (towerClosed) {
      $("conn").textContent = t("anslutning_stangt");
      return;
    }
    $("conn").textContent = t("anslutning_bruten");
    setTimeout(connect, 2000);
  };
}

function handle(msg) {
  switch (msg.t) {
    case "state":
      if (state && state.phase !== msg.phase && msg.phase === "spel") lastTick = -1;
      // Styrningen hör till vandringen och ska inte ligga över slutkortet.
      if (state && state.phase === "spel" && msg.phase !== "spel") toggleHelp(false);
      state = msg;
      state.receivedAt = performance.now();
      if (msg.sprak && msg.sprak !== UI.sprak) changeLanguage(msg.sprak, true);
      trackSeen(state);
      joined = true;
      {
        // Namnet tornet gav, så att sidan kommer tillbaka till samma vandrare.
        const me = msg.players.find((p) => p.id === msg.you);
        if (me && me.name !== savedName) {
          myName = savedName = me.name;
          saveName(myName);
        }
      }
      showPhase();
      if (!mapReady || mapSize !== msg.size) {
        MapView.setup($("map"), msg.size);
        mapReady = true;
        mapSize = msg.size;
      }
      MapView.draw(state);
      padArrows();
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
    case "rum":
      Log.room(msg);
      break;
    case "runda":
      showRoundCard(msg);
      break;
    case "sound":
      Sound.play(msg.s);
      break;
    case "stangt":
      towerClosed = true;
      $("conn").textContent = t("anslutning_stangt");
      Sound.stop();
      break;
    case "vaknar":
      MapView.wake();
      showWakeCard(msg);
      break;
    case "ormkort":
      showSerpentCard(msg);
      break;
    case "ping":
      MapView.ping(msg.x, msg.y, msg.slot);
      break;
    case "error":
      if (!joined) {
        $("nameerr").textContent = msg.text;
        show("namescreen");
        $("back").focus();
        ws.onclose = null;
        ws.close();
      } else {
        Log.push(msg.text, "varn");
      }
      break;
  }
}

// Byte av språk: spelledaren har valt ett annat språk, och servern säger
// det i läget och i svaret från /intro. Sidan hämtar de nya texterna och
// ritar om det den själv har skrivit, utan att laddas om.
async function changeLanguage(lang, againIntro) {
  if (!(await refreshTexts("spelarsidan", lang))) return;
  Log.relabel();
  if (zoomKey) $("zoomName").textContent = t(zoomKey);
  const head = document.querySelector("#roundcard .cardhead");
  if (head) head.textContent = t("rundkort_rubrik", { runda: head.dataset.round });
  if (state) {
    renderSide();
    renderTop();
    renderOverlay();
    renderChoose();
    restoreTip();
  }
  if (!againIntro) return;
  try {
    const info = await (await fetch("/intro", { cache: "no-store" })).json();
    if (state && mapReady) MapView.draw(state);
  } catch {}
}

const ACTION_WORDS = {
  N: "drag_norrut", E: "drag_osterut", S: "drag_soderut", W: "drag_vasterut",
  vanta: "drag_vanta", ta: "drag_ta", sok: "drag_sok", drick: "drag_drick",
  anfall: "drag_anfall", parera: "drag_parera", undan: "drag_undan",
};
function actionWord(a) {
  return ACTION_WORDS[a] ? t(ACTION_WORDS[a]) : a;
}

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

// sleepMeter visar i skattkammaren hur djupt ormen sover, i tre steg som
// andas: lugnt när den sover djupt och fort när den är nära att vakna.
// Ormtjusaren ser också vakna-korten som är kvar i leken, som små kort.
function sleepMeter(level, text, wakeCards) {
  const d = document.createElement("div");
  d.className = "somn somn" + Math.min(2, level);
  const row = document.createElement("div");
  row.className = "somnrad";
  const l = document.createElement("span");
  l.className = "mlabel";
  l.textContent = t("nu_ormen");
  const segs = document.createElement("span");
  segs.className = "andetag";
  segs.setAttribute("aria-hidden", "true");
  for (let i = 0; i < 3; i++) {
    const s = document.createElement("span");
    s.className = i <= level ? "seg on" : "seg";
    segs.appendChild(s);
  }
  row.append(l, segs);
  if (typeof wakeCards === "number") {
    const k = document.createElement("span");
    k.className = "vaknakort";
    k.title = t("nu_vaknakort", { antal: wakeCards });
    k.setAttribute("aria-label", k.title);
    for (let i = 0; i < wakeCards; i++) k.appendChild(document.createElement("i"));
    row.appendChild(k);
  }
  const p = document.createElement("div");
  p.className = level >= 2 ? "cond" : "sleep";
  p.textContent = text;
  d.append(row, p);
  return d;
}

// Andra vandrare syns bara i samma rum. Klienten minns var och när den såg
// var och en senast, så att kartan kan visa det.
let lastSeen = {};
function trackSeen(st) {
  if (st.phase !== "spel") {
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
  if (state.phase !== "spel" || p.id === state.you) return "";
  if (p.status !== "inne" || !p.hidden) return "";
  const s = lastSeen[p.id];
  return s ? t("vandrare_sags", { runda: s.round }) : t("vandrare_inte_sedd");
}

function playerTag(p) {
  if (!p.online) return t("vandrare_borta");
  if (p.status === "ute") return t("vandrare_ute");
  if (p.status === "ny") return t("vandrare_ny");
  if (p.status === "forlorad") return t("vandrare_forlorad");
  if (p.status === "dod") return t(state.phase === "spel" && p.chosen ? "vandrare_dod_klar" : "vandrare_dod");
  if (state.phase === "spel") {
    const seen = seenTag(p);
    const lage = t(p.chosen ? "vandrare_klar" : "vandrare_valjer");
    return seen ? t("vandrare_lage_sedd", { lage, sedd: seen }) : lage;
  }
  return p.ready ? t("vandrare_redo") : "";
}

// Klassvalet har en egen skärm mellan lobbyn och vandringen. Kartan ritas
// om när spelskärmen syns igen, eftersom den inte har någon storlek medan
// den är dold.
function showPhase() {
  const chron = wantsChron();
  if (chron && chron !== chronWas) chronMap = false; // en ny anledning att visa krönikan
  chronWas = chron;
  renderChron();
  const want = state.phase === "klassval" ? "classscreen" : chron && !chronMap ? "chronscreen" : "game";
  if ($(want).classList.contains("active")) return;
  show(want);
  if (want === "game") MapView.resize();
}

// Krönikan på egen skärm. Den som har dött ser sin egen vandring: rundans
// kort runda för runda, hur hen dog och sina siffror, men inget om var de
// andra är. Knappen Visa kartan går tillbaka till spelet, och Din vandring
// i panelen Nu hit igen.
// När vandringen är slut får alla skärmen, med resultaten, epilogen, hela
// krönikan och sin egen vandring under.
let chronMap = false; // kartan visas i stället för krönikan
let chronWas = "";
function wantsChron() {
  const me = myPlayer();
  if (!state || !state.vandring || !me) return "";
  if (state.phase === "slut") return "slut";
  return state.phase === "spel" && me.status === "dod" ? "dod" : "";
}
function renderChron() {
  const v = state && state.vandring;
  if (!v) return;
  const end = state.phase === "slut";
  const rows = (id, list) => {
    const box = $(id);
    box.innerHTML = "";
    (list || []).forEach((l) => {
      const d = document.createElement("div");
      d.textContent = l.text;
      if (l.c) d.className = l.c;
      if (l.kort && l.kort.length) d.appendChild(cardRow(l.kort));
      box.appendChild(d);
    });
    box.hidden = !box.childElementCount;
  };
  $("chronTitle").textContent = end ? endTitle() : t("vandring_rubrik");
  $("chronLead").textContent = end ? "" : v.dod || "";
  $("chronStatus").textContent = end ? "" : t("vandring_pagar", { runda: state.round, sol: state.sun, inne: v.inne, ute: v.ute });
  // Resultaten: vem som kom ut med hur mycket, vinnaren först.
  rows("chronResults", end ? (state.results || []).map((r) => ({ text: (r.winner ? "★ " : "") + r.name + (r.ai ? " ~" : "") + " · " + resultWord(r), c: r.winner ? "win" : r.status !== "ute" ? "lost" : "" })) : []);
  rows("chronEpilog", end ? state.epilog : []);
  rows("chronAll", end ? (state.kronika || []).filter((l) => l.c !== "rubrik").map((l) => ({ text: l.text, c: l.c, kort: l.kort })) : []);
  if (end && state.kronika && state.kronika[0]) {
    const h = document.createElement("h2");
    h.textContent = state.kronika[0].text;
    $("chronAll").prepend(h);
  }
  $("chronOwn").textContent = end ? t("vandring_rubrik") : "";
  $("chronOwn").hidden = !end;
  const next = $("chronNext");
  const me = myPlayer();
  next.hidden = !end;
  next.textContent = me && me.ready ? t("nu_vantar") : t("knapp_nytt_torn");
  next.disabled = !!(me && me.ready);
  const box = $("chronRounds");
  box.innerHTML = "";
  if (!v.rundor.length) {
    const p = document.createElement("p");
    p.className = "dimline";
    p.textContent = t("vandring_tom");
    box.appendChild(p);
  }
  v.rundor.forEach((r) => {
    const d = document.createElement("div");
    d.className = "chronRound";
    const h = document.createElement("h3");
    h.textContent = r.rum ? t("vandring_runda_rum", { runda: r.runda, rum: r.rum }) : t("vandring_runda", { runda: r.runda });
    d.appendChild(h);
    if (r.har) {
      const here = document.createElement("div");
      here.className = "here";
      here.textContent = r.har;
      d.appendChild(here);
    }
    r.rader.forEach((l) => {
      const p = document.createElement("div");
      p.textContent = l.text;
      if (l.ormkort) p.prepend(miniCard(l.ormkort));
      if (l.c) p.className = l.c;
      d.appendChild(p);
    });
    box.appendChild(d);
  });
  $("chronStats").textContent = end ? "" : v.siffror;
  if (!end && v.kort && v.kort.length) $("chronStats").appendChild(cardRow(v.kort));
  $("chronWait").textContent = end ? (touch ? "" : t("slut_nytt_torn")) : t("vandring_vantar");
}
$("chronNext").addEventListener("click", () => send({ t: "ready" }));
$("chronMap").addEventListener("click", () => {
  chronMap = true;
  showPhase();
});
function showChron() {
  chronMap = false;
  showPhase();
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
    if (c.liv) stats.push(t(c.liv > 0 ? "klasskort_liv_mer" : "klasskort_liv_mindre", { n: Math.abs(c.liv) }));
    if (c.drycker) stats.push(t(c.drycker === 1 ? "klasskort_dryck_en" : "klasskort_dryck_flera", { n: c.drycker }));
    if (stats.length) add("stats", stats.join(" · "));
    add("good", t("klasskort_styrka", { styrka: c.styrka }));
    add("bad", t("klasskort_svaghet", { svaghet: c.svaghet }));
    add("owner", c.din ? t("klasskort_ditt") : c.tagen ? t("klasskort_tagen", { namn: c.tagen }) : "");
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
  const parts = [];
  if (mine) parts.push(t("klassval_du_ar", { klass: mine.namn.toLowerCase() }));
  else if (touch) parts.push(t("klassval_tryck_kort"));
  else parts.push(t("klassval_tryck_siffra", { antal: state.erbjudande.length }));
  if (mine && waiting.length) parts.push(t("klassval_vantar", { namn: waiting.join(", ") }));
  el.textContent = parts.join(" ") + " ";
  const s = document.createElement("span");
  s.className = "secs" + (secs <= 5 ? " urgent" : "");
  s.textContent = t("klassval_sekunder", { s: secs });
  el.appendChild(s);
}

function renderSide() {
  renderReady();
  const me = myPlayer();
  const room = me && state.rooms.find((r) => r.x === me.x && r.y === me.y);
  let here = me && me.status === "inne" && room ? room.name : "";
  if (state.cellar > 0) here = t("rum_kallarvalven");
  else if (state.dark > 0 && room) here = t("rum_morker");
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
  if (c) icon(c.tecken || c.namn[0], "", t("tecken_klass", { namn: c.namn, beskrivning: c.beskrivning, styrka: c.styrka, svaghet: c.svaghet }), "klass");
  if (me.status === "dod") {
    icon("✕", "", t("tecken_dod"), "warn");
  } else {
    const low = state.hp <= Math.max(1, Math.floor(state.maxHp / 3));
    icon("♥", state.hp + "/" + state.maxHp, t("tecken_liv", { liv: state.hp, max: state.maxHp }), low ? "warn" : "hp");
  }
  icon("¤", String(me.treasure), t("tecken_skatt", { guld: me.treasure }), "gold");
  if (state.potions > 0) {
    const tryck = pressWord("knapp_drick", "D");
    icon("⚱", String(state.potions), state.potions === 1 ? t("tecken_dryck_en", { tryck }) : t("tecken_dryck_flera", { n: state.potions, tryck }), "");
  }
  if (me.status === "inne") {
    if (state.cellar > 0) icon("▼", "", t("tecken_kallare"), "warn");
    if (state.stuck > 0) icon("⛓", String(state.stuck), t(state.stuck === 1 ? "tecken_fast_en" : "tecken_fast_flera", { n: state.stuck }), "warn");
    if (state.poison > 0) icon("☠", "", t("tecken_gift"), "warn");
    if (state.dark > 0) icon("◐", "", t("tecken_morker"), "warn");
  }
  if (state.kompass) icon("✧", "", state.kompass, "");
  if (state.hjartslag) icon("♡", "", state.hjartslag, "");
  if (state.invigd) icon("❝", "", t("tecken_invigd"), "rost");
  if (state.kanViska) icon("V", "", t("tecken_viska", { tryck: pressWord("knapp_viska", "V") }), "rost");
  // Delarna av ormens namn som vandraren har fått, med prickar för resten.
  if (state.namndelar) icon("⁂", state.namndelar.map((p) => p || "·").join(" "), t("tecken_namndelar"), "rost");
  // Viska har en fast plats bland knapparna under vandringen och är osynlig
  // tills den går att använda, så att inga andra knappar flyttar sig när den
  // dyker upp. Redo behövs bara i lobbyn och i slutet och tar då platsen.
  const playing = state.phase === "spel";
  const pw = $("padWhisper");
  if (pw) {
    pw.hidden = !playing;
    pw.classList.toggle("off", !state.kanViska);
  }
  const pe = $("padEnter");
  if (pe) pe.hidden = playing;

  // Väskan: en plats för varje föremål man kan bära, med föremålets tecken
  // och namn. Under vandringen står de tomma platserna som en prick, så att
  // man ser hur mycket det finns plats för. Ett nytt föremål lyser en stund.
  const items = state.foremal || [];
  const walking = state.phase === "spel" && me.status === "inne";
  const slots = walking ? Math.max(items.length, state.maxForemal || 3) : items.length;
  const now = performance.now();
  for (let i = 0; i < slots; i++) {
    const it = items[i];
    const id = it ? it.id : "";
    if (bagIds[i] !== id) {
      if (id && bagIds[i] !== undefined) bagNewAt[i] = now;
      bagIds[i] = id;
    }
    const b = document.createElement("span");
    b.dataset.slot = "f" + i;
    if (!it) {
      b.className = "hudicon item empty";
      b.textContent = "·";
      b.tabIndex = 0;
      b.dataset.tip = t("tecken_tom_plats", { max: state.maxForemal || 3 });
      b.setAttribute("aria-label", b.dataset.tip);
      bag.appendChild(b);
      continue;
    }
    b.tabIndex = 0;
    const fresh = now - (bagNewAt[i] || -Infinity);
    b.className = "hudicon item" + (fresh < BAG_NEW_MS ? " fresh" : "");
    if (fresh < BAG_NEW_MS) b.style.animationDelay = -fresh + "ms";
    if (it.tecken) {
      const g = document.createElement("span");
      g.className = "glyph";
      g.textContent = it.tecken;
      b.appendChild(g);
    }
    b.appendChild(document.createTextNode(it.namn + (it.kvar > 1 ? " ×" + it.kvar : "")));
    b.dataset.tip = it.beskrivning;
    b.setAttribute("aria-label", it.namn + ". " + b.dataset.tip);
    bag.appendChild(b);
  }
  if (!walking && !items.length) bagIds.length = 0;
}

// Vilket föremål som låg på varje plats i väskan förra gången, och när ett
// nytt kom dit, så att det nya kan lysa en stund.
const BAG_NEW_MS = 2600;
const bagIds = [];
const bagNewAt = [];

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
  if (!me) {
    $("nowBlock").hidden = $("whisper").hidden;
    return;
  }
  if (me.status === "dod" && state.phase === "spel") line(t("nu_dod"), "cond");
  // Tillbaka till krönikans skärm: för den som har dött, och för alla i
  // slutet när de har gått till kartan.
  if (wantsChron()) {
    const b = document.createElement("button");
    b.type = "button";
    b.className = "chronBack";
    b.textContent = t(state.phase === "slut" ? "kronika_visa" : "vandring_visa");
    b.addEventListener("click", showChron);
    box.appendChild(b);
  }
  if (me.status === "inne") {
    if (state.sleep >= 0 && state.sleepText) box.appendChild(sleepMeter(state.sleep, state.sleepText, state.vaknaKort));
    if (state.stuck > 0) line(state.stuck === 1 ? t("nu_fast_en") : t("nu_fast", { n: state.stuck }), "cond");
    const foe = (state.monsters || []).find((m) => m.x === me.x && m.y === me.y);
    const here = (state.rooms || []).find((r) => r.x === me.x && r.y === me.y);
    if (!foe && state.cellar === 0 && here && here.k === "kammare" && state.phase === "spel") {
      line(t("nu_kammare", { tryck: pressWord("knapp_ta", "T") }), "dimline");
    }
    if (foe && state.cellar === 0) {
      line(t("nu_strid", { namn: foe.name }), "foe");
      box.appendChild(meter("foehp", t("nu_dess_liv"), foe.hp, foe.maxHp));
      if (!touch) line(t("nu_strid_tangenter"), "dimline");
    }
    const others = foe || state.cellar > 0 ? [] : foesHere();
    if (others.length) {
      const names = others.map((p) => p.name).join(" " + t("nu_och") + " ");
      const f = currentFoe();
      if (state.pvp === "ja") {
        if (others.length > 1) line(t("nu_strid_byt", { namn: f.name, tangent: keyName("Tab") }), "foe");
        else line(t("nu_strid", { namn: f.name }), "foe");
        if (!touch) line(t("nu_strid_tangenter"), "dimline");
      } else {
        line(t("nu_har_star", { namn: names }), "dimline");
        if (state.pvp === "fredat") line(t("nu_fredad"), "dimline");
      }
    }
  }
  if (state.phase === "spel" && me.status === "inne") {
    line(state.action ? t("nu_drag", { drag: actionWord(state.action) }) : t("nu_valj_ditt"), "plan");
  } else if (state.phase === "spel" && me.status === "ute") {
    line(me.treasure > 0 ? t("nu_ute", { guld: me.treasure }) : t("nu_ute_tom"), "plan");
  }
  // Panelen Nu syns bara när den har något att säga.
  $("nowBlock").hidden = box.childElementCount === 0 && $("whisper").hidden;
}

// Rundans sammanfattning: ett kort över kartan med vad som hände dig, det
// farliga först. Kortet försvinner när du väljer
// nästa drag, efter en stund, eller vid ett klick.
const CARD_MS = 7000;
let cardTimer = null;
let cardAction = "";
function showRoundCard(msg) {
  const lines = msg.lines || [];
  if (!lines.length) return;
  const card = $("roundcard");
  card.innerHTML = "";
  const h = document.createElement("div");
  h.className = "cardhead";
  h.dataset.round = msg.round;
  h.textContent = t("rundkort_rubrik", { runda: msg.round });
  card.appendChild(h);
  lines.forEach((l) => {
    const d = document.createElement("div");
    d.textContent = l.text;
    if (l.ormkort) d.prepend(miniCard(l.ormkort));
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

// Ormkortet: kortet du drog när du tog en näve i skattkammaren. Det vänds
// upp på kartan med en bild av ormen, kortets text och guldet, och glider
// undan efter en stund eller vid ett klick. Bilderna ritas här som SVG.
const SERPENT_MS = 6500;
let serpentTimer = null;

// serpentSpiral ger slingorna som en bana: en spiral från mitten och ut,
// med huvudet i änden.
function serpentSpiral(cx, cy, turns, r0, grow) {
  const pts = [];
  for (let a = 0; a <= turns * Math.PI * 2; a += 0.15) {
    const r = r0 + grow * a;
    pts.push([cx + Math.cos(a) * r, cy + Math.sin(a) * r * 0.62]);
  }
  return pts;
}

// serpentArt ritar bilden på kortet: ormen som sover hoprullad, som rör
// sig med ett öga på glänt, eller två gula ögon i mörkret.
function serpentArt(id, wakes) {
  const kind = wakes || id === "vaknar" ? "vaknar" : id === "ror_sig" ? "ror_sig" : "sover";
  const f = (n) => n.toFixed(1);
  const pts = kind === "vaknar"
    ? serpentSpiral(50, 50, 2.6, 3, 2.5) // slingorna i mörkret, innanför rutan
    : serpentSpiral(50, 60, kind === "ror_sig" ? 2.3 : 2.6, 4, 4.6);
  const d = "M" + pts.map(([x, y]) => f(x) + " " + f(y)).join("L");
  const [hx, hy] = pts[pts.length - 1];
  const [px, py] = pts[pts.length - 4];
  const ang = (Math.atan2(hy - py, hx - px) * 180) / Math.PI;
  const head = `<g transform="translate(${f(hx)} ${f(hy)}) rotate(${f(ang)})">
<ellipse rx="8" ry="5.5" fill="#4f5e2f" stroke="#1d2410" stroke-width="1.2"/>
${kind === "sover"
    ? `<path d="M0 -2.2q2.5 1.6 5 0" fill="none" stroke="#1d2410" stroke-width="1.1"/>`
    : `<ellipse cx="2.4" cy="-1.6" rx="1.8" ry="1.1" fill="#e8c040"/><path d="M0.6 -2.6q1.8 -0.8 3.6 0" fill="none" stroke="#1d2410" stroke-width="1"/>`}
</g>`;
  if (kind === "vaknar") {
    return `<svg viewBox="0 0 100 100" aria-hidden="true">
<rect x="4" y="4" width="92" height="92" rx="6" fill="#1a0b08"/>
<path d="${d}" fill="none" stroke="#3a1410" stroke-width="11" stroke-linecap="round" opacity="0.8"/>
${[34, 66].map((x) => `<ellipse cx="${x}" cy="46" rx="11" ry="7" fill="#ffd23a"/><ellipse cx="${x}" cy="46" rx="2" ry="6.5" fill="#1a0b08"/>`).join("")}
<circle cx="50" cy="46" r="40" fill="none" stroke="#a8402c" stroke-width="1.5" opacity="0.6"/>
</svg>`;
  }
  const moves = kind === "ror_sig"
    ? `<path d="M10 34q6 -6 12 0M78 30q6 -6 12 0M84 80q4 -5 8 0" fill="none" stroke="#6b4a24" stroke-width="1.6" stroke-linecap="round"/>`
    : `<text x="78" y="28" font-size="9" fill="#6b4a24" font-family="serif">z</text><text x="84" y="20" font-size="7" fill="#6b4a24" font-family="serif">z</text>`;
  return `<svg viewBox="0 0 100 100" aria-hidden="true">
<ellipse cx="50" cy="62" rx="40" ry="22" fill="#c9a24f" opacity="0.35"/>
${[[38, 70], [56, 74], [66, 64], [30, 60]].map(([x, y]) => `<ellipse cx="${x}" cy="${y}" rx="3.2" ry="1.6" fill="#d9a63a" stroke="#7a5a1a" stroke-width="0.6"/>`).join("")}
<path d="${d}" fill="none" stroke="#1d2410" stroke-width="11.5" stroke-linecap="round"/>
<path d="${d}" fill="none" stroke="#5b6b34" stroke-width="9" stroke-linecap="round"/>
<path d="${d}" fill="none" stroke="#8a9a52" stroke-width="2" stroke-linecap="round" stroke-dasharray="1.5 3"/>
${head}
${moves}
</svg>`;
}

// miniCard är ett litet ormkort bredvid en rad i rundans kort och krönikan:
// sover (en sovande slinga), rör sig (en våg) eller vaknar (gula ögon).
function miniCard(id) {
  const s = document.createElement("span");
  s.className = "minikort " + (id === "vaknar" ? "vaknar" : id === "ror_sig" ? "ror_sig" : "sover");
  s.setAttribute("aria-hidden", "true");
  return s;
}

// cardRow är en rad små ormkort, med namnen när musen är över dem.
function cardRow(list) {
  const r = document.createElement("span");
  r.className = "kortrad";
  list.forEach((c) => {
    const m = miniCard(c.id);
    m.title = c.namn;
    r.appendChild(m);
  });
  r.setAttribute("aria-label", list.map((c) => c.namn).join(", "));
  return r;
}

// buildSerpentCard bygger ett kort med baksidan och framsidan, så att det
// kan vändas.
function buildSerpentCard(el, msg, extra) {
  el.innerHTML = "";
  const inner = document.createElement("div");
  inner.className = "ok-inner";
  const back = document.createElement("div");
  back.className = "ok-back";
  const front = document.createElement("div");
  front.className = "ok-front" + (msg.vaknar || msg.kort === "vaknar" ? " vaknar" : "");
  const name = document.createElement("div");
  name.className = "ok-namn";
  name.textContent = msg.namn || "";
  const art = document.createElement("div");
  art.className = "ok-bild";
  art.innerHTML = serpentArt(msg.kort, msg.vaknar);
  const text = document.createElement("div");
  text.className = "ok-text";
  text.textContent = msg.text || "";
  front.append(name, art, text);
  if (extra) front.appendChild(extra);
  inner.append(back, front);
  el.appendChild(inner);
  el.hidden = false;
  el.classList.remove("leaving", "vand");
  void el.offsetWidth; // börja med baksidan upp
  el.classList.add("vand");
}

function showSerpentCard(msg) {
  const el = $("ormkort");
  const gold = document.createElement("div");
  gold.className = "ok-guld";
  gold.textContent = t("ormkort_guld", { antal: msg.guld });
  buildSerpentCard(el, msg, gold);
  clearTimeout(serpentTimer);
  serpentTimer = setTimeout(hideSerpentCard, SERPENT_MS);
}
function hideSerpentCard() {
  const el = $("ormkort");
  if (el.hidden || el.classList.contains("leaving")) return;
  clearTimeout(serpentTimer);
  el.classList.add("leaving");
  setTimeout(() => {
    if (el.classList.contains("leaving")) el.hidden = true;
  }, 400);
}
$("ormkort").addEventListener("click", hideSerpentCard);

// showWakeCard visar för alla att ormen har vaknat: kortet "Ormen vaknar"
// stort mitt på kartan, med det du märkte och, om du stod i kammaren, vad
// du förlorade.
const WAKE_CARD_MS = 5000;
let wakeTimer = null;
function showWakeCard(msg) {
  if (!msg.namn) return;
  hideSerpentCard();
  const el = $("vakenkort");
  let extra = null;
  if (msg.inne) {
    extra = document.createElement("div");
    extra.className = "ok-forlust";
    extra.textContent = t("vakenkort_forlust", { guld: msg.forlust || 0, liv: msg.skada || 0 });
  }
  buildSerpentCard(el, { kort: "vaknar", vaknar: true, namn: msg.namn, text: msg.text }, extra);
  clearTimeout(wakeTimer);
  wakeTimer = setTimeout(hideWakeCard, WAKE_CARD_MS);
}
function hideWakeCard() {
  const el = $("vakenkort");
  if (el.hidden || el.classList.contains("leaving")) return;
  clearTimeout(wakeTimer);
  el.classList.add("leaving");
  setTimeout(() => {
    if (el.classList.contains("leaving")) el.hidden = true;
  }, 400);
}
$("vakenkort").addEventListener("click", hideWakeCard);

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
    // Med en mus visas den redan när musen är över tecknet, så där
    // visar ett klick den bara.
    if (touch && !$("tip").hidden && tipKey === keyOf(el)) hideTip();
    else showTip(el);
  });
}
// Ett tryck någon annanstans döljer förklaringen.
document.addEventListener("click", (e) => {
  if (tipKey && !e.target.closest("#hud, #crew")) hideTip();
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
const TOUCH_NAMES = { Tab: "knapp_byt", Enter: "knapp_redo" };
function keyName(k) {
  return touch && TOUCH_NAMES[k] ? t(TOUCH_NAMES[k]) : k;
}

// pressWord säger "Tryck" och knappen på en pekskärm, och annars tangenten.
function pressWord(button, key) {
  return touch ? t("tryck_knapp", { knapp: t(button) }) : t("tryck_tangent", { tangent: key });
}

// newTowerWord säger hur man går vidare till ett nytt torn i slutet.
function newTowerWord() {
  return touch ? t("tryck_knapp", { knapp: t("knapp_nytt_torn") }) : t("slut_nytt_torn");
}

// playerName ger vandrarens namn och om det är du, en AI eller datorn.
function playerName(p, aiKey) {
  if (p.id === state.you) return t("vandrare_namn_du", { namn: p.name });
  if (p.ai) return t(aiKey, { namn: p.name });
  if (p.bot) return t("vandrare_namn_dator", { namn: p.name });
  return p.name;
}

// Pekknappen för Enter heter olika saker beroende på vad den gör.
function setEnterLabel(text) {
  const b = $("padEnter");
  if (b.textContent !== text) b.textContent = text;
}

// Knappen Redo i översta raden, i lobbyn och slutet (Nytt torn). Enter gör
// samma sak. På en pekskärm är det knappen längst ner som gäller.
function renderReady() {
  const me = myPlayer();
  const waiting = state.phase === "lobby" || state.phase === "slut";
  const done = !!(me && me.ready);
  const label = done ? t("knapp_redo_klar") : t(state.phase === "slut" ? "knapp_nytt_torn" : "knapp_redo");
  setEnterLabel(waiting ? label : t("knapp_redo"));
  const b = $("readyBtn");
  b.hidden = !waiting || !me;
  b.disabled = done;
  b.classList.toggle("done", done);
  b.textContent = label;
  if (!done) {
    const k = document.createElement("span");
    k.className = "key";
    k.textContent = "Enter";
    b.appendChild(k);
  }
}
$("readyBtn").addEventListener("click", () => send({ t: "ready" }));

function hpMeter(hp, maxHp) {
  return meter("hp", t("matare_liv"), hp, maxHp);
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
  const m = meter("sun", t("topp_sol"), filled, SUN_WIDTH, String(state.sun));
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
    // I lobbyn och slutet står namnet bredvid pricken.
    if (state.phase !== "spel") {
      d.classList.add("named");
      const n = document.createElement("span");
      n.className = "dotname";
      n.textContent = p.name;
      d.appendChild(n);
    }
    const who = playerName(p, "vandrare_namn_ai_spelar");
    const tag = playerTag(p);
    d.dataset.tip = (p.klass ? t("vandrare_med_klass", { vem: who, klass: p.klass }) : who) + (tag ? "\n" + tag : "");
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
    const lage = t(n < state.minPlayers ? "topp_vantar_vandrare" : "topp_vantar_redo");
    el.textContent = t("topp_lobby", { lage, storlek: state.size });
    el.classList.remove("urgent");
    return;
  }
  if (state.phase === "klassval") {
    el.textContent = t("topp_klassval");
    updateChooseStatus();
    return;
  }
  if (state.phase === "slut") {
    el.textContent = t(state.sun > 0 ? "topp_slut" : "topp_solen_nere");
    el.classList.remove("urgent");
    return;
  }
  if (state.paused) {
    el.textContent = t("topp_pausad", { runda: state.round });
    el.classList.remove("urgent");
    return;
  }
  const secs = Math.ceil(remainingMs() / 1000);
  el.textContent = t("topp_runda", { runda: state.round, sekunder: String(secs).padStart(2, " ") });
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
  $("conn").textContent = "";
  $("overlay").hidden = true;
  $("log").textContent = "";
  $("nameerr").textContent = "";
  Intro.reset();
  show("start");
  $("enter").focus();
}

// resultWord säger hur det gick för en vandrare: guldet, tomhänt, död eller förlorad.
function resultWord(r) {
  if (r.status === "ute") return r.treasure > 0 ? t("slut_guld", { n: r.treasure }) : t("slut_tomhant");
  if (r.status === "dod") return t("slut_dod");
  return t("slut_forlorad");
}

function endTitle() {
  return t(state.sun > 0 ? "slut_rubrik" : "slut_rubrik_solen");
}

function renderOverlay() {
  const el = $("overlay");
  if (state.phase !== "slut" || !state.results) {
    el.hidden = true;
    return;
  }
  const lines = [{ text: endTitle() }, { text: "" }];
  state.results.forEach((r, i) => {
    const what = resultWord(r);
    const name = (r.winner ? "*" : " ") + (r.ai ? "~" : " ") + r.name;
    const row = name.padEnd(20, " ") + what.padStart(10, " ");
    lines.push({ text: row, cls: r.winner ? "win" : r.status !== "ute" ? "lost" : "" });
  });
  if (state.results.some((r) => r.ai)) lines.push({ text: "" }, { text: t("slut_ai") });
  lines.push({ text: "" }, { text: newTowerWord() });
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
  // På tornets egen dator leder introt till värdsidan, utom när man har
  // valt att spela själv (adressen slutar med ?spela).
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
      if (info.sprak && info.sprak !== UI.sprak) changeLanguage(info.sprak, false);
      reveal = !!info.avslojande;
      if (!reveal) $("intro").classList.add("invigd");
      toHostPage = info.vard && !wantsToPlay;
      $("listen").textContent = (touch && info.ljudTryck) || info.ljud || "";
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
      location.href = "spelledare.html";
      return;
    }
    // Tornet ger namnet. Sidan skickar det namn tornet gav förra gången,
    // inom 12 timmar, så att den som kommer tillbaka får sin vandrare.
    // Fliken kommer ihåg sitt eget namn, så att två flikar i samma
    // webbläsare inte byter vandrare med varandra.
    myName = loadName();
    connect();
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

// padArrows visar pilarna på pekskärmens knappar åt det håll figuren går.
// Snett ovanifrån pekar norr upp åt höger, så pilarna blir sneda.
function padArrows() {
  const glyphs = { n: "↗", e: "↘", s: "↙", w: "↖" };
  Object.entries(glyphs).forEach(([c, ch]) => {
    const b = document.querySelector("#pad .dirs ." + c);
    if (b && b.textContent !== ch) b.textContent = ch;
  });
}

// Zoomen: knapparna i kartans hörn, scrollhjulet och två fingrar väljer hur
// nära kartan går. + och − på tangentbordet, och M för allt man känner till,
// går genom pressKey.
let zoomKey = ""; // nyckeln till zoomnivåns namn
MapView.setZoomListener((name, n) => {
  zoomKey = name;
  $("zoomName").textContent = t(name);
  $("zoomIn").disabled = n === 0;
  $("zoomOut").disabled = n === 3;
});
$("zoomIn").addEventListener("click", () => MapView.zoom(-1));
$("zoomOut").addEventListener("click", () => MapView.zoom(1));
$("zoomName").addEventListener("click", () => MapView.toggleWhole());
{
  // Ett hjul ger många små steg; ett steg i taget räcker.
  let wheelAt = 0;
  $("mapbox").addEventListener("wheel", (e) => {
    if (!state || state.phase !== "spel" || Math.abs(e.deltaY) < 1) return;
    e.preventDefault();
    const now = performance.now();
    if (now - wheelAt < 250) return;
    wheelAt = now;
    MapView.zoom(e.deltaY > 0 ? 1 : -1);
  }, { passive: false });
  // Två fingrar: isär går närmare, ihop går längre ut.
  let pinch = 0;
  const dist = (t) => Math.hypot(t[0].clientX - t[1].clientX, t[0].clientY - t[1].clientY);
  $("mapbox").addEventListener("touchstart", (e) => {
    if (e.touches.length === 2) pinch = dist(e.touches);
  }, { passive: true });
  $("mapbox").addEventListener("touchmove", (e) => {
    if (e.touches.length !== 2 || !pinch) return;
    e.preventDefault();
    const d = dist(e.touches);
    if (d > pinch * 1.35) {
      MapView.zoom(-1);
      pinch = d;
    } else if (d < pinch / 1.35) {
      MapView.zoom(1);
      pinch = d;
    }
  }, { passive: false });
  $("mapbox").addEventListener("touchend", () => (pinch = 0));
}

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

const KEYS = {
  ArrowUp: "N", ArrowRight: "E", ArrowDown: "S", ArrowLeft: "W",
  " ": "vanta", t: "ta", T: "ta", s: "sok", S: "sok", d: "drick", D: "drick",
  a: "anfall", A: "anfall", p: "parera", P: "parera", u: "undan", U: "undan",
  // E (evade) är undanmanövern på engelska; båda bokstäverna fungerar på båda språken.
  e: "undan", E: "undan",
};
document.addEventListener("keydown", (e) => {
  if (e.metaKey || e.ctrlKey || e.altKey) return;
  if (e.target && e.target.id === "whisperText") return; // viskningen skrivs i sin egen ruta
  if (pressKey(e.key)) e.preventDefault();
});

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
  // Esc på skärmen med felet går tillbaka till startsidan.
  if (key === "Escape" && $("namescreen").classList.contains("active")) {
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
  if (state && state.phase === "spel") {
    if (key === "+" || key === "=") {
      MapView.zoom(-1);
      return true;
    }
    if (key === "-") {
      MapView.zoom(1);
      return true;
    }
    if (key === "m" || key === "M") {
      MapView.toggleWhole();
      return true;
    }
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

// Tecknen vid dörrarna på kartan snett ovanifrån går att trycka på, på en
// pekskärm (se .isotangenter i style.css).
document.addEventListener("click", (e) => {
  const s = e.target.closest && e.target.closest(".isotangenter span[data-key]");
  if (!s || !touch) return;
  e.stopPropagation();
  pressKey(s.dataset.key);
}, true);

// Pekskärm: knapparna visas från början på en telefon, och annars så fort
// någon rör vid skärmen.
let touch = !!(window.matchMedia && matchMedia("(pointer: coarse)").matches);
function setTouch() {
  document.body.classList.toggle("touch", touch);
  if (state) {
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
