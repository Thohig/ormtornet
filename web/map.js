// Kartan på spelskärmen. Tornet ritas snett ovanifrån med Three.js
// (IsoKarta, byggd från klient/ till web/iso/iso.js). Här räknas utsnittet
// ut, alltså vilka rum som syns och hur nära kartan går, och vilka tecken
// som står i teckenförklaringen. Där WebGL saknas visas en text i stället
// för kartan.
"use strict";

const MapView = (() => {
  // Åtta färger som skiljer sig i ljushet och ton, så att de går att skilja
  // åt även för den som har svårt att se skillnad på rött och grönt.
  const PLAYER_COLORS = ["#ffcc66", "#66e0ff", "#ff7a9a", "#9dff7a", "#c9a0ff", "#ff9a3c", "#3fd0b0", "#f2f2f2"];

  // Färgerna på tecknen i teckenförklaringen, där ingen modell står.
  const COLORS = {
    unknown: "#3a3d43",
    exit: "#8d8a82",
    shut: "#d0503a",
    gold: "#ffd36b",
    awake: "#ffe14a",
    lost: "#a04a3c",
    mark: "#ff8a6a",
    monster: "#ff7a5a",
    chalk: "#f2f0e8",
    offline: "#5a5d63",
  };
  const ARROWS = ["▲", "▶", "▼", "◀"];

  const WAKE_MS = 3500; // så länge ormen syns vaken efter att den har vaknat
  const PING_MS = 4500; // så länge ett eko syns
  let wokeAt = -Infinity;
  let pings = [];

  // ---- Utsnittet ----

  // Kartan följer det rum vandraren står i, som står i mitten med lika
  // många rum åt varje håll. Vandraren väljer själv hur nära kartan går,
  // eller att den visar allt man känner till. Nivåernas namn är nycklar
  // till spelarsidans texter.
  const MIN_SPAN = 3; // så nära går kartan som mest, i rum
  const ZOOM_MS = 1100; // en ny nivå: en lång, mjuk zoom
  const PAN_MS = 400; // samma nivå i ett nytt rum: kartan glider snabbt efter
  const LEVELS = [
    { name: "zoom_nara", span: 3 },
    { name: "zoom_mellan", span: 5 },
    { name: "zoom_langt", span: 9 },
    { name: "zoom_allt", span: 0 },
  ];
  const ALL = LEVELS.length - 1;
  const ZOOM_KEY = "ormtornet.zoom";
  let level = loadLevel();
  let before = level === ALL ? 1 : level; // nivån som M går tillbaka till
  let onZoom = null;
  let towerSize = 0;
  // target är rutan kartan ska visa, i rum: x, y, bredd och höjd. Kartan
  // glider dit från from, med start vid zoomAt.
  let from = null, target = null, zoomAt = -Infinity, zoomMs = ZOOM_MS;

  // loadLevel ger nivån som valdes senast på den här enheten, annars Nära
  // på en smal skärm och Mellan på en större, så att tecknen går att läsa.
  function loadLevel() {
    try {
      const n = parseInt(localStorage.getItem(ZOOM_KEY), 10);
      if (n >= 0 && n < LEVELS.length) return n;
    } catch {}
    return window.matchMedia("(max-width: 700px)").matches ? 0 : 1;
  }

  function setLevel(n) {
    n = Math.max(0, Math.min(ALL, n));
    if (n !== ALL) before = n;
    if (n === level) return;
    level = n;
    try {
      localStorage.setItem(ZOOM_KEY, String(level));
    } catch {}
    if (onZoom) onZoom(LEVELS[level].name, level);
    if (lastState) draw(lastState);
  }

  // zoom går närmare (d < 0) eller längre ut (d > 0) ett steg.
  function zoom(d) {
    setLevel(level + d);
  }

  // toggleWhole växlar mellan den valda nivån och allt man känner till.
  function toggleWhole() {
    setLevel(level === ALL ? before : ALL);
  }

  function setZoomListener(fn) {
    onZoom = fn;
    fn(LEVELS[level].name, level);
  }

  // focus ger den egna vandraren, så länge den är inne och syns.
  function focus(state) {
    const p = state.you && state.players.find((q) => q.id === state.you);
    return p && p.status === "inne" && !p.hidden && p.x >= 0 ? p : null;
  }

  // wanted räknar ut rutan kartan ska visa: rutan runt vandraren på den
  // valda nivån, eller allt man känner till, med lite luft omkring och
  // aldrig närmare än MIN_SPAN rum.
  function wanted(state) {
    let x0 = Infinity, y0 = Infinity, x1 = -Infinity, y1 = -Infinity;
    const take = (x, y) => {
      if (x < 0 || y < 0) return;
      x0 = Math.min(x0, x);
      y0 = Math.min(y0, y);
      x1 = Math.max(x1, x + 1);
      y1 = Math.max(y1, y + 1);
    };
    state.rooms.forEach((r) => take(r.x, r.y));
    // Bara den egna vandraren räknas med; de andra syns inte på avstånd.
    state.players.forEach((p) => {
      if (p.id === state.you) take(p.x, p.y);
    });
    if (x0 === Infinity) return { x: 0, y: 0, w: towerSize, h: towerSize };
    const all = grow({ x: x0, y: y0, w: x1 - x0, h: y1 - y0 }, MIN_SPAN);
    const me = focus(state);
    const n = LEVELS[level].span;
    if (!n || !me) return all;
    return around(me, n);
  }

  // around ger en ruta på n × n rum med rummet där p står i mitten. Vid
  // tornets kanter skjuts rutan in, men den får gå lite utanför (EDGE rum),
  // så att ett rum i hörnet inte hamnar under knapparna och tecknen i
  // kartans hörn, och portarna ut syns.
  const EDGE = 0.75;
  function around(p, n) {
    const fit = (c) => (n >= towerSize ? (towerSize - n) / 2 : Math.max(-EDGE, Math.min(towerSize - n + EDGE, c + 0.5 - n / 2)));
    return { x: fit(p.x), y: fit(p.y), w: n, h: n };
  }

  function grow(v, span) {
    const w = Math.max(span, v.w);
    const h = Math.max(span, v.h);
    return { x: v.x - (w - v.w) / 2, y: v.y - (h - v.h) / 2, w, h };
  }

  function same(a, b) {
    return a && b && ["x", "y", "w", "h"].every((k) => Math.abs(a[k] - b[k]) < 0.01);
  }

  // aim sätter rutan kartan ska visa. En ny ruta nås med en mjuk zoom från
  // den som syns just nu, eller glider snabbt dit om bara platsen har bytts;
  // första gången hoppar kartan direkt dit.
  function aim(t, now) {
    if (!target) {
      target = t;
      zoomAt = -Infinity;
      return;
    }
    if (same(t, target)) return;
    from = viewAt(now);
    zoomMs = Math.abs(t.w - from.w) < 0.01 && Math.abs(t.h - from.h) < 0.01 ? PAN_MS : ZOOM_MS;
    target = t;
    zoomAt = now;
  }

  // viewAt ger rutan som syns vid tiden now, i rum. Kartan snett ovanifrån
  // frågar efter den varje gång den ritar.
  function viewAt(now) {
    if (!target) return null;
    if (zoomAt === -Infinity) return target;
    const f = Math.min(1, (now - zoomAt) / zoomMs);
    if (f >= 1) {
      zoomAt = -Infinity;
      return target;
    }
    const e = f < 0.5 ? 2 * f * f : 1 - (-2 * f + 2) ** 2 / 2;
    const v = {};
    ["x", "y", "w", "h"].forEach((k) => (v[k] = from[k] + (target[k] - from[k]) * e));
    return v;
  }

  // focusRoom ger rummet kartan följer, eller null.
  function focusRoom() {
    const p = lastState && focus(lastState);
    return p ? { x: p.x, y: p.y } : null;
  }

  // ---- Teckenförklaringen ----

  // De tecken som står på kartan just nu. Nyckeln säger vilken förklaring
  // som gäller, name är fällans eller varelsens namn. Där kartan har en
  // modell står en bild av den i stället för tecknet (legendImage).
  let legendSig = "";
  let onLegend = null;
  const LEGEND_ORDER = ["du", "vandrare", "borta", "skymd", "forlorad", "drag", "varelse", "falla",
    "kammare", "vaken", "port", "port_stangd", "foremal", "krita", "eko", "okand", "mitten"];

  function legendOf(state, now) {
    const legend = new Map();
    const note = (key, glyph, color, name) => {
      if (!legend.has(key)) legend.set(key, { key, glyph, color, name: name || "" });
    };
    const lvl = state.reveal ? Math.min(4, Math.max(1, state.level || 1)) : 1;

    // Okända rum: för en vandrare bara rummen bredvid dem den känner till.
    const size = state.size;
    const mid = Math.floor(size / 2);
    const known = new Set(state.rooms.map((r) => r.x + "," + r.y));
    const near = (x, y) => [[0, 0], [1, 0], [-1, 0], [0, 1], [0, -1]].some(([dx, dy]) => known.has(x + dx + "," + (y + dy)));
    for (let y = 0; y < size; y++) {
      for (let x = 0; x < size; x++) {
        if (!near(x, y) || known.has(x + "," + y)) continue;
        if (x === mid && y === mid) note("mitten", "◌", COLORS.unknown);
        else note("okand", "·", COLORS.unknown);
      }
    }

    state.rooms.forEach((r) => {
      if (r.exit >= 0) {
        if (r.stangd) note("port_stangd", "█", COLORS.shut);
        else if (lvl === 4 && state.phase === "slut") note("port", "█", COLORS.shut);
        else if (lvl === 4) note("port", "⩛", COLORS.monster);
        else note("port", "˄", COLORS.exit);
      }
      if (r.mark) note("falla:" + r.mark, r.mark, COLORS.mark, r.markn);
      if (r.f) note("foremal", "+", COLORS.gold); // ett föremål som någon har lämnat, för spårfinnaren
      if (r.c >= 0 && r.c !== undefined) note("krita", ARROWS[r.c], COLORS.chalk);
      if (r.k === "kammare") {
        note("kammare", "¤", COLORS.gold);
        if (now - wokeAt < WAKE_MS) note("vaken", "◉", COLORS.awake);
      }
    });

    pings = pings.filter((p) => now - p.at < PING_MS);
    pings.forEach((p) => note("eko", "◎", PLAYER_COLORS[p.slot]));

    state.players.forEach((p) => {
      const seen = state.lastSeen && state.lastSeen[p.id];
      if (seen && p.hidden) note("skymd", "○", PLAYER_COLORS[seen.slot] || "#ffffff");
    });
    state.players.forEach((p) => {
      if (p.status === "ute" || p.hidden) return;
      if (p.status === "forlorad" || p.status === "dod") {
        note("forlorad", "x", COLORS.lost);
        return;
      }
      const color = p.online ? PLAYER_COLORS[p.slot] || "#fff" : COLORS.offline;
      if (p.id === state.you) note("du", "◉", color);
      else if (!p.online) note("borta", "●", color);
      else note("vandrare", "●", color);
    });

    (state.monsters || []).forEach((m) => note("varelse:" + m.glyph, m.glyph, COLORS.monster, m.name));

    // Det egna, valda draget.
    const me = state.players.find((p) => p.id === state.you);
    const dir = { N: 0, E: 1, S: 2, W: 3 }[state.action];
    if (me && me.status === "inne" && dir !== undefined) note("drag", ARROWS[dir], PLAYER_COLORS[me.slot]);

    const rank = (e) => LEGEND_ORDER.indexOf(e.key.split(":")[0]);
    return [...legend.values()].sort((a, b) => rank(a) - rank(b));
  }

  // tellLegend säger till sidan när tecknen på kartan har ändrats.
  function tellLegend(state, now) {
    if (!onLegend) return;
    const list = legendOf(state, now);
    const sig = list.map((e) => e.key + e.glyph + e.color + e.name).join("|");
    if (sig === legendSig) return;
    legendSig = sig;
    onLegend(list);
  }

  function setLegend(fn) {
    onLegend = fn;
    legendSig = "";
    if (lastState) tellLegend(lastState, performance.now());
  }

  // legendImage ger en bild av modellen för ett tecken i teckenförklaringen,
  // eller null där tecknet själv står på kartan.
  function legendImage(e) {
    return iso ? iso.bild(e.key, e.glyph, e.color) : null;
  }

  // ---- Kartan ----

  let canvas = null;
  let iso = null;
  let lastState = null;
  let timer = null; // ritar om teckenförklaringen när ett eko eller ormens uppvaknande är över

  // setup gör i ordning kartan på canvas för ett torn av storleken size.
  // Den anropas igen när ett nytt torn har en annan storlek.
  function setup(el, size) {
    canvas = el;
    towerSize = size;
    target = null; // hoppa direkt till rätt utsnitt
    if (!iso && window.IsoKarta && IsoKarta.stods()) {
      iso = IsoKarta.skapa(canvas, { utsnitt: () => viewAt(performance.now()), fokus: focusRoom, farger: PLAYER_COLORS });
      iso.visa(true);
    }
    if (iso) iso.storlek(size);
    const none = document.getElementById("mapNone");
    if (none) none.hidden = !!iso;
    canvas.hidden = !iso;
  }

  function draw(state) {
    lastState = state;
    if (state.size) towerSize = state.size;
    const now = performance.now();
    aim(wanted(state), now);
    if (iso) iso.visaTillstand(state); // först, så att den vet vilka varelser som finns när förklaringen frågar efter bilder
    tellLegend(state, now);
  }

  // later ritar om teckenförklaringen efter ms, när något tillfälligt är över.
  function later(ms) {
    clearTimeout(timer);
    timer = setTimeout(() => lastState && tellLegend(lastState, performance.now()), ms + 50);
  }

  function resize() {
    if (iso) iso.anpassa();
  }

  function ping(x, y, slot) {
    pings.push({ x, y, slot, at: performance.now() });
    if (iso) iso.eko(x, y, slot);
    if (lastState) tellLegend(lastState, performance.now());
    later(PING_MS);
  }

  function wake() {
    wokeAt = performance.now();
    if (iso) iso.vakna();
    if (lastState) tellLegend(lastState, wokeAt);
    later(WAKE_MS);
  }

  return { setup, draw, resize, ping, wake, zoom, toggleWhole, setZoomListener, setLegend, legendImage, PLAYER_COLORS };
})();
