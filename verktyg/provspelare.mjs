// Provspelare: ansluter till tornet som en människa över websocketen,
// spelar ett antal vandringar och skriver alla rader den får i loggen till
// en fil (eller till terminalen). Rader som ser trasiga ut, till exempel med
// en {platshållare} kvar, "1 rundor" eller svenska ord på engelska, märks
// med "<<< MISSTÄNKT" och regeln (se textkontroll.mjs). Texterna i
// rundans kort, spelkorten, krönikan och resten av tillståndet kontrolleras
// också; en trasig text där skrivs en gång på en rad som börjar med
// MISSTÄNKT. Den används för att hitta fel i texter och lägen som
// sällan dyker upp när man spelar för hand. Se README, "Provspela med
// robotar", och docs/provspel-2026-11-21.md.
//
//   node verktyg/provspelare.mjs <port> <namn> <antal vandringar> [loggfil]
//
// Kräver Node 22 eller senare (inbyggd WebSocket). Den spelar enkelt: slåss
// mot varelser, går mot kammaren, tar skatt högst två gånger, söker ibland,
// använder föremål och viskar, och går ut mot närmaste kända port. Kör
// tornet med -min satt så att alla provspelare hinner ansluta innan
// vandringen börjar.
import fs from "node:fs";
import { kontrollera, texterI } from "./textkontroll.mjs";

const [port, name, gamesArg, logFile] = process.argv.slice(2);
const games = Number(gamesArg || 1);
const out = fs.createWriteStream(logFile || `/dev/stdout`, { flags: "a" });
const w = (s) => out.write(`[${name}] ${s}\n`);

const ws = new WebSocket(`ws://127.0.0.1:${port}/ws`);
let st = null, lastRound = -1, readySent = "", classSent = false, played = 0;
let takes = 0, whispered = false, lastPhase = "";
let sprak = "sv";
const sett = new Set();

// misstankt kontrollerar texterna i ett meddelande från tornet och skriver
// varje trasig text en gång.
function misstankt(msg) {
  for (const text of texterI(msg)) {
    if (sett.has(text)) continue;
    sett.add(text);
    const fel = kontrollera(text, sprak);
    if (fel.length) w(`MISSTÄNKT [${fel.join(", ")}] (${msg.t}) ${text}`);
  }
}

const send = (o) => ws.send(JSON.stringify(o));
const DX = [0, 1, 0, -1], DY = [-1, 0, 1, 0], NESW = "NESW";

function me() { return st.players.find((p) => p.id === st.you); }
function room(x, y) { return st.rooms.find((r) => r.x === x && r.y === y); }

function bfs(goal) {
  const m = me();
  const prev = new Map([[m.x + "," + m.y, null]]);
  const q = [[m.x, m.y]];
  while (q.length) {
    const [x, y] = q.shift();
    const r = room(x, y);
    if (goal(r, x, y)) {
      let k = x + "," + y, path = [];
      while (prev.get(k)) { const [px, py, d] = prev.get(k); path.unshift(d); k = px + "," + py; }
      return path;
    }
    if (!r) continue;
    for (let d = 0; d < 4; d++) {
      if (!r.d[d]) continue;
      const nx = x + DX[d], ny = y + DY[d], k = nx + "," + ny;
      if (nx < 0 || ny < 0 || nx >= st.size || ny >= st.size || prev.has(k)) continue;
      prev.set(k, [x, y, d]);
      q.push([nx, ny]);
    }
  }
  return null;
}

function choose() {
  const m = me();
  const r = room(m.x, m.y);
  const foe = (st.monsters || []).find((x) => x.x === m.x && x.y === m.y);
  if (foe && !st.cellar) return { a: ["anfall", "anfall", "anfall", "vanta"][Math.floor(Math.random() * 4)] };
  if (st.kanViska && !whispered && Math.random() < 0.3) {
    whispered = true;
    const parts = st.namndelar || [];
    send({ t: "viska", name: parts.join("") || "orm" });
  }
  if (st.hp <= 2 && st.potions > 0) return { a: "drick" };
  const others = st.players.filter((p) => p.id !== st.you && p.status === "inne" && !p.hidden && p.x === m.x && p.y === m.y);
  if (others.length && st.pvp === "ja" && Math.random() < 0.3) return { a: "anfall", m: -others[0].id };
  if (r && r.k === "kammare" && takes < 2 && st.sleep < 2) { takes++; return { a: "ta" }; }
  if (r && r.k === "rum" && !r.s && Math.random() < 0.3) return { a: "sok" };
  // på väg ut?
  const toExit = bfs((rr) => rr && rr.k === "ingang" && rr.exit >= 0 && !rr.stangd);
  const leaving = m.treasure >= 40 || takes >= 2 || (toExit && st.sun <= toExit.length + 4);
  if (leaving && toExit) {
    if (toExit.length === 0) return { a: NESW[r.exit] };
    return { a: NESW[toExit[0]] };
  }
  const toChamber = bfs((rr) => rr && rr.k === "kammare");
  if (toChamber && toChamber.length) return { a: NESW[toChamber[0]] };
  // utforska: gå mot närmaste okända rum, helst mot mitten
  const c = (st.size - 1) / 2;
  const opts = [];
  for (let d = 0; d < 4; d++) {
    if (!r || !r.d[d]) continue;
    const nx = m.x + DX[d], ny = m.y + DY[d];
    if (nx < 0 || ny < 0 || nx >= st.size || ny >= st.size) continue;
    const score = Math.abs(nx - c) + Math.abs(ny - c) + (room(nx, ny) ? 2 : 0) + Math.random() * 2;
    opts.push([score, d]);
  }
  opts.sort((a, b) => a[0] - b[0]);
  if (opts.length) return { a: NESW[opts[0][1]] };
  return { a: "vanta" };
}

ws.onopen = () => send({ t: "join", name });
ws.onclose = () => { w("STÄNGD"); process.exit(0); };
ws.onmessage = (ev) => {
  const msg = JSON.parse(ev.data);
  if (msg.t === "log" || msg.t === "error") {
    const fel = kontrollera(msg.text, sprak);
    sett.add(msg.text);
    const flag = fel.length ? `  <<< MISSTÄNKT [${fel.join(", ")}]` : "";
    w(`${msg.t === "error" ? "FEL " : ""}r${st ? st.round : "?"} ${msg.text}${flag}`);
    return;
  }
  if (msg.t === "state") sprak = msg.sprak || sprak;
  misstankt(msg);
  if (msg.t !== "state") return;
  st = msg;
  if (st.phase !== lastPhase) { w(`--- fas ${st.phase}`); lastPhase = st.phase; }
  if (st.phase === "lobby" && readySent !== "lobby" + played) {
    readySent = "lobby" + played;
    takes = 0; whispered = false; classSent = false; lastRound = -1;
    setTimeout(() => send({ t: "ready" }), 300);
  }
  if (st.phase === "klassval" && !classSent && st.erbjudande) {
    classSent = true;
    const free = st.erbjudande.filter((c) => !c.tagen);
    const c = free[Math.floor(Math.random() * free.length)] || st.erbjudande[0];
    setTimeout(() => send({ t: "klass", name: c.id }), 500);
  }
  if (st.phase === "spel") {
    const m = me();
    if (m && m.status === "inne" && st.round !== lastRound && !st.action) {
      lastRound = st.round;
      const c = choose();
      setTimeout(() => send({ t: "act", ...c }), 200 + Math.random() * 400);
    }
  }
  if (st.phase === "slut" && !readySent.startsWith("slut")) {
    readySent = "slut" + played;
    const res = (st.results || []).map((r) => `${r.name}:${r.treasure}:${r.status}`).join(" ");
    w(`RESULTAT ${res}`);
    // Berättelsen om vandringen, stycke för stycke, att läsa efteråt.
    for (const b of st.berattelse || []) w(`BERÄTTELSE ${b.text}`);
    for (const v of st.vandrarna || []) w(`BERÄTTELSE ${v.namn}${v.titel ? ", " + v.titel : ""}. ${v.text}`);
    // Guldböckerna: hur var och en fick sitt guld, och raderna runda för runda.
    for (const b of st.guldbok || []) {
      for (const text of b.berattelse) w(`GULDBOK ${b.namn}: ${text}`);
      for (const r of b.rader) w(`GULDRAD ${b.namn} runda ${r.runda} ${r.guld > 0 ? "+" : ""}${r.guld} = ${r.summa}: ${r.text}`);
    }
    played++;
    if (played >= games) { setTimeout(() => process.exit(0), 500); return; }
    setTimeout(() => send({ t: "ready" }), 800);
  }
};
