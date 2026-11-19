// Spelledarbordet: sidan där spelledaren leder vandringen i
// spelledarläget, med hela tornet, vandrarna, ingreppen och loggen.
// Inbjudan och inställningarna ligger på värdsidan (spelledare.html).
// Texterna står i data/vardsidan.json och hämtas med t("nyckel") (texter.js).
"use strict";
const $ = (id) => document.getElementById(id);
const COLORS = MapView.PLAYER_COLORS;
let view = null; // /spelledare/status: fasen, vandrarna och inställningarna
let lastGM = null; // den senaste översikten, för att rita om när språket byts
let closed = false;

function post(path, body) {
  return fetch(path, {
    method: "POST",
    headers: { "X-Ormtornet": "spelledare", "Content-Type": "application/json" },
    body: body === undefined ? undefined : JSON.stringify(body),
  });
}

// Bara en sida som syns på skärmen håller tornet öppet, inte en flik som
// glömts kvar i bakgrunden.
function seenParam() {
  return document.visibilityState === "visible" ? "?synlig" : "";
}

// Raden överst: fasen, vilka som är här och knappen som behövs just nu,
// Starta nu före vandringen och Pausa eller Fortsätt under den.
let paused = false;

function renderBar(v) {
  const phase = { lobby: t("rad_lobby"), klassval: t("rad_klassval"), slut: t("rad_slut") }[v.phase] ||
    t("rad_runda", { runda: v.round, sol: v.sun, max: v.sunMax });
  $("barPhase").textContent = paused && v.phase === "spel" ? t("rad_pausad", { fas: phase }) : phase;
  const who = $("barWho");
  who.textContent = "";
  const here = v.players.filter((p) => p.online);
  const humans = here.filter((p) => !p.bot).length;
  who.append(t(humans === 1 ? "rad_har_manniska" : "rad_har_manniskor", { antal: humans, ai: here.length - humans }));
  const away = v.players.filter((p) => !p.online && !p.bot).map((p) => p.name);
  if (away.length) {
    const s = document.createElement("span");
    s.className = "away";
    s.textContent = " · " + t("rad_borta", { namn: away.join(", ") });
    who.append(s);
  }
  const b = $("barMain");
  const gm = v.settings.spelledarlage;
  b.hidden = !gm || v.phase === "klassval";
  if (v.phase === "spel") {
    b.textContent = paused ? t("fortsatt") : t("pausa");
    b.dataset.act = paused ? "fortsatt" : "paus";
    b.disabled = false;
  } else {
    b.textContent = t("starta_nu");
    b.dataset.act = "starta";
    b.disabled = v.players.length === 0;
  }
}

$("barMain").onclick = async () => {
  const act = $("barMain").dataset.act;
  if (act !== "starta") {
    intervene({ typ: act });
    return;
  }
  try {
    const r = await post("/spelledare/starta");
    if (!r.ok) {
      const data = await r.json().catch(() => ({}));
      alert(data.fel || t("starta_fel"));
    }
  } catch {
    alert(t("svarar_inte"));
  }
};

function render(v) {
  const gm = v.settings.spelledarlage;
  $("off").hidden = gm;
  $("gm").hidden = !gm;
  $("door").hidden = !gm || (v.phase !== "lobby" && v.phase !== "klassval");
  $("end").hidden = !gm || v.phase !== "slut";
  $("turnOn").disabled = v.phase === "spel";
  if (!gm && v.phase === "spel") $("offMsg").textContent = t("bord_sla_pa_spel");
  else if ($("offMsg").dataset.err !== "1") $("offMsg").textContent = "";
  paused = !!v.paused;
  renderBar(v);
  if (gm && !$("door").hidden) renderDoor(v);
}

// Spelledarläget slås på härifrån, med de andra inställningarna som de är.
// Det går bara mellan vandringarna.
$("turnOn").onclick = async () => {
  if (!view) return;
  const msg = $("offMsg");
  msg.dataset.err = "";
  try {
    const r = await post("/spelledare/installningar", { ...view.settings, spelledarlage: true });
    const data = await r.json().catch(() => ({}));
    if (!r.ok || data.varning) {
      msg.dataset.err = "1";
      msg.textContent = data.fel || data.varning || t("inst_fel");
    }
  } catch {
    msg.dataset.err = "1";
    msg.textContent = t("svarar_inte");
  }
};

// Vid porten: vilka som har kommit, om de är redo och vilken klass de har
// valt. Klassen kommer från översikten.
function doorTag(p, phase) {
  if (!p.online) return [t("borta"), "away"];
  if (p.bot) return [t("tagg_dator"), ""];
  if (phase === "klassval") return p.chosen ? [t("tagg_har_valt"), "ok"] : [t("tagg_valjer"), ""];
  return p.ready ? [t("tagg_redo"), "ok"] : [t("tagg_inte_redo"), ""];
}

function renderDoor(v) {
  $("doorNote").textContent = v.players.length
    ? t("torn_info", { storlek: v.size, max: v.sunMax })
    : t("bord_ingen_vid_porten");
  const ul = $("doorList");
  ul.textContent = "";
  const classes = {};
  (lastGM ? lastGM.players : []).forEach((p) => (classes[p.slot] = p.klass));
  v.players
    .slice()
    .sort((a, b) => a.slot - b.slot)
    .forEach((p) => {
      const li = document.createElement("li");
      const name = document.createElement("span");
      name.className = "name";
      const mark = document.createElement("span");
      mark.textContent = "@ ";
      mark.style.color = COLORS[p.slot];
      name.append(mark, p.name);
      const [text, cls] = doorTag(p, v.phase);
      const tg = document.createElement("span");
      tg.className = "tag " + cls;
      tg.textContent = text;
      li.append(name, tg);
      if (classes[p.slot]) {
        const k = document.createElement("span");
        k.textContent = classes[p.slot];
        li.append(k);
      }
      ul.append(li);
    });
}

// Slutet: hur det gick för var och en, och krönikan.
function renderEnd(o) {
  const ul = $("results");
  ul.textContent = "";
  (o.results || []).forEach((r) => {
    const li = document.createElement("li");
    const name = document.createElement("span");
    name.className = "name" + (r.winner ? " win" : "");
    const mark = document.createElement("span");
    mark.textContent = "@ ";
    mark.style.color = COLORS[r.slot];
    name.append(mark, r.name);
    const st = document.createElement("span");
    const status = { ute: t("ute"), forlorad: t("forlorad"), dod: t("dod") }[r.status] || r.status;
    st.textContent = status + " · " + t("guld", { antal: r.treasure }) + (r.winner ? " · " + t("bord_vann") : "");
    if (r.winner) st.className = "win";
    li.append(name, st);
    ul.append(li);
  });
  const box = $("chronicle");
  box.textContent = "";
  (o.chronicle || []).forEach((l) => {
    const d = document.createElement("div");
    if (l.c) d.className = l.c;
    d.textContent = l.text;
    box.append(d);
  });
}

// Översikten i spelledarläget: hela tornet, alla vandrare, ormen och loggen.
let mapReady = false;
let mapSize = 0;
const chosenItem = {}; // föremålet spelledaren har valt att ge, per vandrare
let following = 0; // vandraren kartan följer, eller 0 för hela tornet

$("gmPlayers").addEventListener("click", (ev) => {
  const b = ev.target.closest("button[data-follow]");
  if (!b) return;
  const id = +b.dataset.follow;
  following = following === id ? 0 : id;
  MapView.follow(following);
  $("gmFollow").textContent = following ? t("folj_text", { namn: b.textContent.replace(/^◉ /, "") }) : "";
});

// whoText: vandrarens namn i tabellen, med AI, dator och klass.
function whoText(p) {
  let s = p.ai ? t("namn_ai", { namn: p.name }) : p.bot ? t("namn_dator", { namn: p.name }) : p.name;
  if (p.klass) s = t("namn_klass", { vem: s, klass: p.klass });
  return s;
}

// serpentText: ormens rad i översikten.
function serpentText(s) {
  const key = !s.wakes ? "top_ormen" : s.wakes === 1 ? "top_ormen_vaknat_en" : "top_ormen_vaknat";
  return t(key, { somn: t("somn_" + s.sleep), sov: s.sleepCards, vakna: s.wakeCards, ganger: s.wakes, risk: s.risk });
}

function renderGM(o) {
  lastGM = o;
  const top = $("gmTop");
  top.textContent = "";
  const m = o.map;
  const parts = [];
  if (m.phase === "spel") parts.push(t("top_runda", { runda: m.round, sol: m.sun, max: m.sunMax }));
  parts.push(serpentText(o.serpent));
  if (m.reveal) {
    const key = !o.held ? "top_niva" : o.sunLevel > o.level ? "top_niva_halls_sol" : "top_niva_halls";
    parts.push(t(key, { niva: o.level, solniva: o.sunLevel, ledtradar: o.clues }));
  }
  parts.push(t(o.traps === 1 ? "top_falla" : "top_fallor", { antal: o.traps }));
  if (o.namn) parts.push(t("top_namn", { namn: o.namn }));
  if (o.nad) parts.push(t("top_nad"));
  if (o.oppnaPortar < o.portar) parts.push(t("top_portar", { oppna: o.oppnaPortar, alla: o.portar }));
  top.textContent = parts.join(" · ");

  if (!mapReady || mapSize !== m.size) {
    MapView.setup($("gmMap"), m.size);
    if (!mapReady && window.ResizeObserver) new ResizeObserver(() => MapView.resize()).observe($("gmMapBox"));
    mapReady = true;
    mapSize = m.size;
  }
  MapView.draw(m);

  const tb = $("gmPlayers").querySelector("tbody");
  tb.textContent = "";
  o.players.forEach((p) => {
    const tr = document.createElement("tr");
    if (p.status !== "inne") tr.className = "dead";
    const cond = [];
    if (p.cellar) cond.push(t("till_kallaren", { antal: p.cellar }));
    if (p.stuck) cond.push(t("till_fast", { antal: p.stuck }));
    if (p.poison) cond.push(t("till_gift", { antal: p.poison }));
    if (p.dark) cond.push(t("till_morker", { antal: p.dark }));
    if (!p.online) cond.push(t("borta"));
    if (p.foremal && p.foremal.length) cond.push(p.foremal.join(", "));
    if (p.namndelar) cond.push(p.namndelar.map((d) => d || "·").join(" "));
    if (p.invigd) cond.push(t("till_invigd"));
    else if (p.viskning) cond.push(t("till_viskade", { text: p.viskning }));
    const who = whoText(p);
    const where = { ute: t("ute"), forlorad: t("forlorad"), dod: t("dod") }[p.status] || p.room || "";
    const action = p.action || (m.phase === "spel" && p.status === "inne" ? t("drag_valjer") : "");
    const cells = [who, where, p.hp + "/" + (p.maxHp || o.maxHp), t("guld", { antal: p.treasure }), p.potions, cond.join(", "), action];
    cells.forEach((c, i) => {
      const td = document.createElement("td");
      td.textContent = c;
      if (i === 0) {
        // Ett klick på namnet låter kartan följa vandraren.
        td.textContent = "";
        const b = document.createElement("button");
        b.type = "button";
        b.className = "follow" + (following === p.id ? " on" : "");
        b.dataset.follow = String(p.id);
        b.title = following === p.id ? t("folj_hela") : t("folj_titel", { namn: p.name });
        b.textContent = (following === p.id ? "◉ " : "") + c;
        b.style.color = COLORS[p.slot];
        td.append(b);
      }
      if (i === 5 && c) td.className = "cond";
      tr.append(td);
    });
    // Följer kartan vandraren skrivs raden under kartan på nytt, så att den
    // får nya texter när språket byts.
    if (following === p.id) $("gmFollow").textContent = t("folj_text", { namn: who });
    const give = document.createElement("td");
    give.className = "give";
    if (m.phase === "spel" && p.status === "inne") {
      [["liv", -1, "ge_liv_minus"], ["liv", 1, "ge_liv_plus"], ["guld", -10, "ge_guld_minus"], ["guld", 10, "ge_guld_plus"], ["dryck", -1, "ge_dryck_minus"], ["dryck", 1, "ge_dryck_plus"]]
        .forEach(([typ, n, key]) => {
          const b = document.createElement("button");
          b.type = "button";
          b.className = "btn small";
          b.dataset.gm = typ;
          b.dataset.antal = String(n);
          b.dataset.spelare = String(p.id);
          b.textContent = t(key);
          give.append(b);
        });
      // Inviga, eller ta tillbaka, och ge ett föremål.
      if (m.reveal) {
        const inv = document.createElement("button");
        inv.type = "button";
        inv.className = "btn small";
        inv.dataset.gm = "invig";
        inv.dataset.spelare = String(p.id);
        inv.textContent = p.invigd ? t("ge_glom_invigd") : t("ge_inviga");
        give.append(inv);
      }
      // Översikten ritas om hela tiden: kom ihåg vilket föremål som är valt.
      const sel = document.createElement("select");
      (o.foremal || []).forEach((it) => sel.append(new Option(it.namn, it.id)));
      if (chosenItem[p.id]) sel.value = chosenItem[p.id];
      sel.addEventListener("change", () => (chosenItem[p.id] = sel.value));
      const gb = document.createElement("button");
      gb.type = "button";
      gb.className = "btn small";
      gb.dataset.gm = "foremal";
      gb.dataset.spelare = String(p.id);
      gb.dataset.fromSelect = "1";
      gb.textContent = t("ge_foremal");
      give.append(sel, gb);
    }
    tr.append(give);
    tb.append(tr);
  });

  renderOmens(o);
  renderLog(o);
  if (m.phase === "slut") renderEnd(o);
}

// Varslen som har getts, de senaste först: vem som fick dem, vad som väntar
// och om någon redan har gått in i rummet.
function renderOmens(o) {
  const el = $("gmOmens");
  const list = o.varsel || [];
  el.textContent = list.length
    ? t("varsel_rubrik") + " " + list.slice(0, 8).map((v) =>
      t(v.found ? "varsel_rad_funnet" : "varsel_rad", { runda: v.round, vem: v.to, vad: v.what || t("varsel_tomt"), x: v.x, y: v.y })).join(" · ")
    : "";
}

// Vandrarnas logg: en grupp per runda, den senaste överst, och inom
// rundan raderna i den ordning de hände. De två senaste rundorna står
// öppna; spelledaren kan fälla upp och ihop de andra själv.
const logOpen = {}; // runda -> öppen eller stängd, när spelledaren har klickat
let logKey = "";

function renderLog(o) {
  const who = $("gmLogWho");
  fillSelect(who, [{ label: t("logg_alla"), value: "" }].concat(o.players.map((p) => ({ label: p.name, value: p.name }))));
  const entries = o.log.filter((e) => !who.value || e.to.includes(who.value));
  const key = JSON.stringify([entries.length, entries.length && entries[entries.length - 1].text, who.value, logOpen, UI.sprak]);
  if (key === logKey) return;
  logKey = key;
  const rounds = new Map();
  entries.forEach((e) => {
    if (!rounds.has(e.round)) rounds.set(e.round, []);
    rounds.get(e.round).push(e);
  });
  const order = [...rounds.keys()].sort((a, b) => b - a);
  const box = $("gmLog");
  const st = box.scrollTop;
  const before = box.scrollHeight;
  box.textContent = "";
  if (!order.length) {
    const d = document.createElement("div");
    d.className = "dim";
    d.textContent = t("logg_tom");
    box.append(d);
  }
  order.forEach((round, i) => {
    const list = rounds.get(round);
    const open = round in logOpen ? logOpen[round] : i < 2;
    const sec = document.createElement("section");
    sec.className = "round";
    const head = document.createElement("button");
    head.type = "button";
    head.className = "roundhead";
    let label = round ? t("logg_runda", { runda: round }) : t("logg_fore");
    if (!open) label += " · " + t(list.length === 1 ? "logg_rad" : "logg_rader", { antal: list.length });
    head.textContent = (open ? "▾ " : "▸ ") + label;
    head.addEventListener("click", () => {
      logOpen[round] = !open;
      renderLog(o);
    });
    sec.append(head);
    if (open) {
      list.forEach((e) => {
        const d = document.createElement("div");
        const w = document.createElement("span");
        w.className = "who";
        w.textContent = e.to.join(", ") + ": ";
        const line = document.createElement("span");
        if (e.c) line.className = e.c;
        line.textContent = e.text;
        d.append(w, line);
        sec.append(d);
      });
    }
    box.append(sec);
  });
  if (st > 6) box.scrollTop = st + (box.scrollHeight - before);
}

// Ingreppen. Listorna byggs bara om när innehållet ändras, så att ett val
// inte hoppar medan spelledaren håller på.
function fillSelect(sel, options) {
  const key = JSON.stringify(options);
  if (sel.dataset.key === key) return;
  const keep = sel.value;
  sel.textContent = "";
  options.forEach((o) => {
    if (o.group) {
      const g = document.createElement("optgroup");
      g.label = o.group;
      o.items.forEach((it) => g.append(new Option(it.label, it.value)));
      sel.append(g);
    } else {
      sel.append(new Option(o.label, o.value));
    }
  });
  sel.dataset.key = key;
  if ([...sel.options].some((x) => x.value === keep)) sel.value = keep;
}

// roomLabel: rummet i listan I ett rum, med fälla och varelse.
function roomLabel(r) {
  const parts = [t("rum_namn", { namn: r.name, x: r.x, y: r.y })];
  if (r.trap) parts.push(t("rum_falla"));
  if (r.monster) parts.push(t("rum_varelse"));
  return parts.join(" · ");
}

function renderControls(o) {
  const playing = o.map.phase === "spel";
  document.querySelectorAll("#gm .gmrow button[data-gm]").forEach((b) => (b.disabled = !playing));
  $("gmPlace").querySelector("button").disabled = !playing || !o.rooms.length;
  paused = !!o.paused;
  const pause = $("gmPause");
  pause.dataset.gm = o.paused ? "fortsatt" : "paus";
  pause.dataset.t = o.paused ? "fortsatt" : "pausa";
  pause.textContent = t(pause.dataset.t);
  $("gmReveal").hidden = !o.map.reveal;
  $("gmLevelUp").disabled = !playing || o.level >= 4;
  const hold = $("gmHold");
  hold.dataset.gm = o.held ? "niva_slapp" : "niva_hall";
  hold.dataset.t = o.held ? "gm_niva_slapp" : "gm_niva_hall";
  hold.textContent = t(hold.dataset.t);
  if (playing) {
    const left = o.paused ? t("top_pausad", { sek: o.timeLeft }) : t("top_tid", { sek: o.timeLeft });
    $("gmTop").textContent += " · " + left;
  }
  fillSelect($("gmVoiceTo"), [{ label: t("rost_alla"), value: "0" }].concat(o.players.map((p) => ({ label: t("rost_till", { namn: p.name }), value: String(p.id) }))));
  fillSelect($("gmRoom"), o.rooms.length
    ? o.rooms.map((r) => ({ label: roomLabel(r), value: r.x + "," + r.y }))
    : [{ label: t("rum_inga"), value: "" }]);
  fillSelect($("gmWhat"), [
    { group: t("grupp_falla"), items: o.trapTypes.map((x) => ({ label: x.namn, value: "falla:" + x.id })) },
    { group: t("grupp_varelse"), items: o.monsterTypes.map((x) => ({ label: x.namn, value: "varelse:" + x.id })) },
  ]);
}

function gmSay(text, isErr) {
  const el = $("gmMsg");
  el.textContent = text;
  el.classList.toggle("err", !!isErr);
}

async function intervene(req) {
  try {
    const r = await post("/spelledare/ingrepp", req);
    if (!r.ok) {
      const data = await r.json().catch(() => ({}));
      gmSay(data.fel || t("ingrepp_fel"), true);
      return false;
    }
    gmSay("");
    return true;
  } catch {
    gmSay(t("svarar_inte"), true);
    return false;
  }
}

$("gm").addEventListener("click", (ev) => {
  const b = ev.target.closest("button[data-gm]");
  if (!b) return;
  const req = { typ: b.dataset.gm };
  if (b.dataset.antal) req.antal = +b.dataset.antal;
  if (b.dataset.spelare) req.spelare = +b.dataset.spelare;
  if (b.dataset.fromSelect) req.id = b.parentElement.querySelector("select").value;
  if (req.typ === "niva_upp" && !confirm(t("fraga_niva_upp"))) return;
  if (req.typ === "vack" && !confirm(t("fraga_vack"))) return;
  intervene(req);
});

$("gmVoice").addEventListener("submit", async (ev) => {
  ev.preventDefault();
  const text = $("gmVoiceText").value.trim();
  if (!text) {
    gmSay(t("rost_tom"), true);
    return;
  }
  if (await intervene({ typ: "rost", spelare: +$("gmVoiceTo").value, text })) {
    $("gmVoiceText").value = "";
    gmSay(t("rost_klar"));
  }
});

$("gmPlace").addEventListener("submit", async (ev) => {
  ev.preventDefault();
  const [x, y] = $("gmRoom").value.split(",").map(Number);
  const [typ, id] = $("gmWhat").value.split(":");
  if (await intervene({ typ, x, y, id })) gmSay(t("klart"));
});

async function pollGM() {
  if (closed) return;
  if (view && view.settings.spelledarlage) {
    try {
      const r = await fetch("/spelledare/oversikt" + seenParam(), { cache: "no-store" });
      if (r.ok) {
        const o = await r.json();
        if (await refreshTexts("vardsidan", o.sprak)) redraw();
        renderGM(o);
        renderControls(o);
      }
    } catch {}
  }
  setTimeout(pollGM, 1000);
}
pollGM();

// redraw ritar om det skriptet själv har skrivit, när språket har bytts.
// applyTexts har redan fyllt i det som står i HTML:en.
function redraw() {
  logKey = "";
  if (view) render(view);
  if (lastGM) {
    renderGM(lastGM);
    renderControls(lastGM);
  }
}

async function poll() {
  if (closed) return;
  try {
    const r = await fetch("/spelledare/status" + seenParam(), { cache: "no-store" });
    if (!r.ok) throw new Error(r.status);
    view = await r.json();
    if (await refreshTexts("vardsidan", view.sprak)) redraw();
    else render(view);
  } catch {
    $("barPhase").textContent = t("status_svarar_inte");
  }
  setTimeout(poll, 1000);
}
poll();
