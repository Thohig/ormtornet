// Ritar tornet som krita på en skiffertavla. Tornet läggs ut som ett rutnät
// av celler: väggarna ritas som kritstreck mellan cellernas mitt, och tecken
// (varelser, märken, pilar) placeras i en egen cell, så att de alltid hamnar
// rätt även om typsnittet hämtar dem från ett annat teckensnitt.
"use strict";

const MapView = (() => {
  const RW = 7; // celler per rum på bredden (inklusive en vägg)
  const RH = 4; // celler per rum på höjden (inklusive en vägg)

  const SLATE = "#15171a";
  const COLORS = {
    unknown: "#3a3d43",
    exit: "#8d8a82",
    gold: "#ffd36b",
    lost: "#a04a3c",
    mark: "#ff8a6a",
    monster: "#ff7a5a",
    chalk: "#f2f0e8",
    offline: "#5a5d63",
  };
  // Åtta färger som skiljer sig i ljushet och ton, så att de går att skilja
  // åt även för den som har svårt att se skillnad på rött och grönt.
  const PLAYER_COLORS = ["#ffcc66", "#66e0ff", "#ff7a9a", "#9dff7a", "#c9a0ff", "#ff9a3c", "#3fd0b0", "#f2f2f2"];

  const ARROWS = ["▲", "▶", "▼", "◀"];

  const WAKE_MS = 3500;
  let wokeAt = -Infinity;

  // Avslöjandet: tornets utseende följer nivån. På nivå 2 andas kritan, på
  // nivå 3 syns revbenen och hjärtslagen, på nivå 4 är det kött och tänder.
  // ring är kammarens färg, floor golvets svaga ton.
  const LOOK = {
    1: { wall: [201, 198, 188], hi: [201, 198, 188], floor: "rgba(255,255,255,0.035)", ring: "#c9a24f", period: 0 },
    2: { wall: [196, 192, 182], hi: [240, 234, 220], floor: "rgba(255,250,240,0.045)", ring: "#c9a24f", period: 4200 },
    3: { wall: [196, 160, 140], hi: [240, 170, 140], floor: "rgba(140,50,36,0.12)", ring: "#c98a5a", period: 1250 },
    4: { wall: [176, 92, 74], hi: [255, 118, 88], floor: "rgba(150,34,24,0.2)", ring: "#c4483a", period: 950 },
  };
  const STEPS = 6; // så många nyanser kritan andas i

  let lvl = 1;
  let reveal = false;
  let nowMs = 0;

  function pulse(x, y) {
    const look = LOOK[lvl];
    if (!look.period) return 0;
    const dist = Math.hypot(x - cols / 2, (y - rows / 2) * 1.6);
    const t = nowMs - dist * 14;
    if (lvl === 2) return 0.5 + 0.5 * Math.sin((2 * Math.PI * t) / look.period);
    // Hjärtslag: två slag tätt efter varandra, sedan vila.
    const ph = (((t % look.period) + look.period) % look.period) / look.period;
    const g = (c) => Math.exp(-(((ph - c) / 0.07) ** 2));
    return Math.max(g(0), 0.7 * g(0.22), g(1));
  }

  function mix(a, b, f) {
    const c = a.map((v, i) => Math.round(v + (b[i] - v) * f));
    return "rgb(" + c.join(",") + ")";
  }

  // Spelarnas platser inne i ett rum, en per plats i listan. Mitten hålls
  // fri åt tecken för fällor och varelser.
  const SLOTS = [[2, 1], [5, 1], [2, 3], [5, 3], [1, 2], [6, 2], [3, 1], [4, 3]];

  let canvas, ctx, cols, rows, cellW, cellH, fontPx;
  // Utsnittet: kartan följer det rum vandraren står i, som står i mitten
  // med lika många rum åt varje håll. Vid tornets kanter stannar rutan, så
  // att det inte blir tomt utanför tornet. Ett tryck på kartan visar allt
  // man känner till. Spelledaren, som inte är någon vandrare, ser hela
  // tornet.
  const MIN_SPAN = 3; // så nära går kartan som mest, i rum
  const ZOOM_MS = 1100;
  let towerSize = 0;
  let wantWhole = false; // den som spelar vill se allt den känner till
  let boxW = 0, boxH = 0, ox = 0, oy = 0;
  // view är rutan som syns, i rum: x, y, bredd och höjd. Den glider mot
  // target när kartan växer.
  let view = null, from = null, target = null, zoomAt = 0, zooming = false;
  // grid och colr: tecknen. wall: väggceller (true), wallc: en färg som
  // tar över kritans, för ekot. Golv, kammare och vandrare samlas i listor.
  let grid, colr, wall, wallc, floors, rings, figures;

  function hash(x, y) {
    let h = (x * 374761393 + y * 668265263) | 0;
    h = (h ^ (h >>> 13)) * 1274126177;
    return (h ^ (h >>> 16)) >>> 0;
  }

  function put(x, y, ch, color) {
    if (x < 0 || y < 0 || x >= cols || y >= rows) return;
    grid[y][x] = ch;
    colr[y][x] = color;
  }

  // Teckenförklaringen: de tecken som står på kartan just nu. Nyckeln säger
  // vilken förklaring som gäller, name är fällans eller varelsens namn.
  let legend = new Map();
  let legendSig = "";
  let onLegend = null;
  function note(key, glyph, color, name) {
    if (!legend.has(key)) legend.set(key, { key, glyph, color, name: name || "" });
  }

  function setWall(x, y, on) {
    if (x < 0 || y < 0 || x >= cols || y >= rows) return;
    wall[y][x] = on;
  }

  function isWall(x, y) {
    return x >= 0 && y >= 0 && x < cols && y < rows && wall[y][x];
  }

  // setup ritar upp kartan för ett torn av storleken size. Den anropas igen
  // när ett nytt torn har en annan storlek.
  let listening = false;
  function setup(el, size) {
    canvas = el;
    ctx = canvas.getContext("2d");
    cols = size * RW + 1;
    rows = size * RH + 1;
    towerSize = size;
    if (!listening) {
      window.addEventListener("resize", resize);
      listening = true;
    }
    resize();
  }

  function resize() {
    if (!canvas) return;
    const box = canvas.parentElement.getBoundingClientRect();
    boxW = Math.floor(box.width);
    boxH = Math.floor(box.height);
    if (boxW <= 0 || boxH <= 0) return; // dold just nu
    const dpr = window.devicePixelRatio || 1;
    canvas.width = boxW * dpr;
    canvas.height = boxH * dpr;
    canvas.style.width = boxW + "px";
    canvas.style.height = boxH + "px";
    ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
    view = null; // hoppa direkt till rätt utsnitt
    if (lastState) draw(lastState);
  }

  // scaleFor ger teckenstorleken när rutan v (i rum) ska få plats.
  function scaleFor(v) {
    return Math.min(boxW / (v.w * RW * 0.62), boxH / (v.h * RH * 1.12));
  }

  // wanted räknar ut rutan kartan ska visa: allt man känner till, med lite
  // luft omkring, och aldrig närmare än MIN_SPAN rum.
  function wanted(state) {
    if (!state.you) return { x: 0, y: 0, w: towerSize, h: towerSize };
    let x0 = Infinity, y0 = Infinity, x1 = -Infinity, y1 = -Infinity;
    const take = (x, y) => {
      if (x < 0 || y < 0) return;
      x0 = Math.min(x0, x);
      y0 = Math.min(y0, y);
      x1 = Math.max(x1, x + 1);
      y1 = Math.max(y1, y + 1);
    };
    state.rooms.forEach((r) => take(r.x, r.y));
    // De andra vandrarna räknas bara med när alla syns hela tiden.
    state.players.forEach((p) => {
      if (p.id === state.you || (state.sight === "alla" && !p.hidden && p.status !== "ute")) take(p.x, p.y);
    });
    if (x0 === Infinity) return { x: 0, y: 0, w: towerSize, h: towerSize };
    const all = grow({ x: x0, y: y0, w: x1 - x0, h: y1 - y0 }, MIN_SPAN);
    const me = state.players.find((p) => p.id === state.you);
    if (wantWhole || !me || me.status !== "inne" || me.x < 0) return all;
    return around(me, span());
  }

  // span är hur många rum utsnittet visar åt varje håll: färre på en smal
  // skärm, så att tecknen går att läsa.
  function span() {
    return window.matchMedia("(max-width: 700px)").matches ? 3 : 5;
  }

  // around ger en ruta på n × n rum med rummet där p står i mitten. Vid
  // tornets kanter skjuts rutan in, så att den stannar inom tornet.
  function around(p, n) {
    const fit = (c) => (n >= towerSize ? (towerSize - n) / 2 : Math.max(0, Math.min(towerSize - n, c + 0.5 - n / 2)));
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
  // den som syns just nu; första gången hoppar kartan direkt dit.
  function aim(t, now) {
    if (!view) {
      view = target = t;
      zoomAt = -Infinity;
      return;
    }
    if (same(t, target)) return;
    from = current(now);
    target = t;
    zoomAt = now;
    if (!zooming) {
      zooming = true;
      requestAnimationFrame(step);
    }
  }

  function step(now) {
    if (lastState) draw(lastState, now);
    if (zoomAt > -Infinity) requestAnimationFrame(step);
    else zooming = false;
  }

  function current(now) {
    if (zoomAt === -Infinity) return target;
    const f = Math.min(1, (now - zoomAt) / ZOOM_MS);
    if (f >= 1) {
      zoomAt = -Infinity;
      return target;
    }
    const e = f < 0.5 ? 2 * f * f : 1 - (-2 * f + 2) ** 2 / 2;
    const v = {};
    ["x", "y", "w", "h"].forEach((k) => (v[k] = from[k] + (target[k] - from[k]) * e));
    return v;
  }

  // place räknar ut teckenstorlek och förskjutning för utsnittet just nu.
  function place(now) {
    view = current(now);
    fontPx = scaleFor(view);
    cellW = fontPx * 0.62;
    cellH = fontPx * 1.12;
    // Rutan läggs mitt på tavlan.
    ox = view.x * RW - (boxW / cellW - view.w * RW) / 2;
    oy = view.y * RH - (boxH / cellH - view.h * RH) / 2;
  }

  function drawRoom(r) {
    const x0 = r.x * RW;
    const y0 = r.y * RH;
    floors.push([x0, y0]);

    // Väggar runt om.
    for (let i = 0; i <= RW; i++) {
      setWall(x0 + i, y0, true);
      setWall(x0 + i, y0 + RH, true);
    }
    for (let j = 1; j < RH; j++) {
      setWall(x0, y0 + j, true);
      setWall(x0 + RW, y0 + j, true);
    }
    // Dörrar: en lucka mitt på varje öppen sida.
    const mid = Math.floor(RW / 2);
    const gapN = [x0 + mid, x0 + mid + 1].filter((x) => x < x0 + RW);
    if (r.d[0]) gapN.forEach((x) => setWall(x, y0, false));
    if (r.d[2]) gapN.forEach((x) => setWall(x, y0 + RH, false));
    const midY = y0 + Math.floor(RH / 2);
    if (r.d[1]) setWall(x0 + RW, midY, false);
    if (r.d[3]) setWall(x0, midY, false);

    // Utgång genom ytterväggen. På sista nivån är porten en mun med tänder,
    // och när spelet är slut har den stängts.
    // Portarna i hörnen går ut åt sidan, portarna mitt på sidorna rakt ut.
    const ex = r.exit === 1 ? x0 + RW : r.exit === 3 ? x0 : r.exit >= 0 ? x0 + mid : -1;
    const ey = r.exit === 0 ? y0 : r.exit === 2 ? y0 + RH : midY;
    if (ex >= 0) {
      setWall(ex, ey, false);
      if (lvl === 4 && lastState && lastState.phase === "slut") {
        put(ex, ey, "█", "#d0503a");
        note("port", "█", "#d0503a");
      } else if (lvl === 4) {
        put(ex, ey, Math.floor(nowMs / 400) % 2 ? "⩚" : "⩛", "#ff7a5a");
        note("port", "⩛", "#ff7a5a");
      } else {
        put(ex, ey, ["˄", "›", "˅", "‹"][r.exit], COLORS.exit);
        note("port", "˄", COLORS.exit);
      }
    }

    if (r.mark) {
      put(x0 + Math.floor(RW / 2), y0 + Math.floor(RH / 2), r.mark, COLORS.mark);
      note("falla:" + r.mark, r.mark, COLORS.mark, r.markn);
    }
    // Ett föremål som någon har lämnat, för spårfinnaren.
    if (r.f) {
      put(x0 + RW - 1, y0 + 1, "+", COLORS.gold);
      note("foremal", "+", COLORS.gold);
    }
    // En pil i krita, ritad av någon vandrare.
    if (r.c >= 0 && r.c !== undefined) {
      put(x0 + 1, y0 + 1, ARROWS[r.c], COLORS.chalk);
      note("krita", ARROWS[r.c], COLORS.chalk);
    }

    // Kammaren: en kritring kring en hög av guld. När ormen har vaknat
    // spänns ringen och två ögon lyser.
    if (r.k === "kammare") {
      const now = performance.now();
      const awake = now - wokeAt < WAKE_MS;
      const flash = awake && Math.floor((now - wokeAt) / 180) % 2 === 0;
      rings.push({ x0, y0, awake, flash });
      const ch = awake ? "◉" : "¤";
      const color = awake ? "#ffe14a" : COLORS.gold;
      put(x0 + 3, y0 + 2, ch, color);
      put(x0 + 4, y0 + 2, ch, color);
      note("kammare", "¤", COLORS.gold);
      if (awake) note("vaken", "◉", "#ffe14a");
    }
  }

  let lastState = null;
  const PING_MS = 4500;
  let pings = [];

  function ping(x, y, slot) {
    pings.push({ x, y, slot, at: performance.now() });
  }

  function draw(state, t) {
    const first = !lastState;
    lastState = state;
    if (first) {
      resize(); // nu vet kartan vem som tittar, och om utsnittet behövs
      return;
    }
    if (!ctx || !boxW) return;
    reveal = !!state.reveal;
    lvl = reveal ? Math.min(4, Math.max(1, state.level || 1)) : 1;
    nowMs = t || performance.now();
    grid = Array.from({ length: rows }, () => new Array(cols).fill(" "));
    colr = Array.from({ length: rows }, () => new Array(cols).fill(COLORS.unknown));
    wall = Array.from({ length: rows }, () => new Array(cols).fill(false));
    wallc = Array.from({ length: rows }, () => new Array(cols).fill(null));
    floors = [];
    rings = [];
    figures = [];
    legend = new Map();

    // Okända rum syns bara som en svag prick. För en vandrare bara rummen
    // bredvid dem den känner till, så att kartan inte avslöjar mer.
    const size = state.size;
    const mid = Math.floor(size / 2);
    const known = new Set(state.rooms.map((r) => r.x + "," + r.y));
    const near = (x, y) => [[0, 0], [1, 0], [-1, 0], [0, 1], [0, -1]].some(([dx, dy]) => known.has(x + dx + "," + (y + dy)));
    for (let y = 0; y < size; y++) {
      for (let x = 0; x < size; x++) {
        if (state.you && !near(x, y)) continue;
        const ch = x === mid && y === mid ? "◌" : "·";
        put(x * RW + Math.floor(RW / 2), y * RH + Math.floor(RH / 2), ch, COLORS.unknown);
        if (!known.has(x + "," + y)) note(ch === "◌" ? "mitten" : "okand", ch, COLORS.unknown);
      }
    }
    state.rooms.forEach(drawRoom);

    const now = t || performance.now();
    const blink = Math.floor(now / 550) % 2 === 0;

    // Eko: rummets väggar lyser i spelarens färg en stund.
    pings = pings.filter((p) => now - p.at < PING_MS);
    pings.forEach((p) => {
      const x0 = p.x * RW;
      const y0 = p.y * RH;
      const on = Math.floor((now - p.at) / 250) % 2;
      for (let j = 0; j <= RH; j++) {
        for (let i = 0; i <= RW; i++) {
          const edge = i === 0 || j === 0 || i === RW || j === RH;
          if (edge && on && wallc[y0 + j]) wallc[y0 + j][x0 + i] = PLAYER_COLORS[p.slot];
        }
      }
      put(x0 + Math.floor(RW / 2), y0 + Math.floor(RH / 2), "◎", PLAYER_COLORS[p.slot]);
      note("eko", "◎", PLAYER_COLORS[p.slot]);
    });
    // Där en dold vandrare sågs senast står en svag ring i dess färg.
    state.players.forEach((p) => {
      const seen = state.lastSeen && state.lastSeen[p.id];
      if (!seen || !p.hidden) return;
      const [sx, sy] = SLOTS[seen.slot] || SLOTS[0];
      figures.push({ x: seen.x * RW + sx, y: seen.y * RH + sy, kind: "ghost", color: PLAYER_COLORS[seen.slot] || "#ffffff" });
      note("skymd", "○", PLAYER_COLORS[seen.slot] || "#ffffff");
    });
    // Vandrarna är kritprickar i sin färg. Den egna har en ring som blinkar.
    state.players.forEach((p) => {
      if (p.status === "ute" || p.hidden) return;
      const [sx, sy] = SLOTS[p.slot] || SLOTS[0];
      const x = p.x * RW + sx;
      const y = p.y * RH + sy;
      if (p.status === "forlorad" || p.status === "dod") {
        put(x, y, "x", COLORS.lost);
        note("forlorad", "x", COLORS.lost);
        return;
      }
      const color = p.online ? PLAYER_COLORS[p.slot] || "#fff" : COLORS.offline;
      figures.push({ x, y, kind: p.id === state.you ? (blink ? "me" : "dot") : "dot", color });
      if (p.id === state.you) note("du", "◉", color);
      else if (!p.online) note("borta", "●", color);
      else note("vandrare", "●", color);
    });

    // Varelserna. Den som styrs av en död spelare blinkar när den är vald.
    (state.monsters || []).forEach((m) => {
      const selected = state.controller && m.id === selectedMonster;
      const color = selected && !blink ? "#ffe14a" : COLORS.monster;
      put(m.x * RW + Math.floor(RW / 2) + 1, m.y * RH + Math.floor(RH / 2), m.glyph, color);
      note("varelse:" + m.glyph, m.glyph, COLORS.monster, m.name);
    });

    // Den döda spelarens drag för en varelse visas som en pil.
    if (state.controller && state.action && state.action[0] === "m") {
      const id = parseInt(state.action.slice(1), 10);
      const d = { N: 0, E: 1, S: 2, W: 3 }[state.action.slice(-1)];
      const m = (state.monsters || []).find((x) => x.id === id);
      if (m && d !== undefined && blink) {
        const x0 = m.x * RW;
        const y0 = m.y * RH;
        const mx = x0 + Math.floor(RW / 2);
        const my = y0 + Math.floor(RH / 2);
        const pos = [[mx, y0], [x0 + RW, my], [mx, y0 + RH], [x0, my]][d];
        put(pos[0], pos[1], ARROWS[d], COLORS.monster);
      }
    }

    // Det egna, valda draget visas som en pil i dörröppningen.
    const me = state.players.find((p) => p.id === state.you);
    const dir = { N: 0, E: 1, S: 2, W: 3 }[state.action];
    if (me && me.status === "inne" && dir !== undefined) {
      const x0 = me.x * RW;
      const y0 = me.y * RH;
      const mx = x0 + Math.floor(RW / 2);
      const my = y0 + Math.floor(RH / 2);
      const pos = [[mx, y0], [x0 + RW, my], [mx, y0 + RH], [x0, my]][dir];
      const color = PLAYER_COLORS[me.slot];
      if (blink) put(pos[0], pos[1], ARROWS[dir], color);
      note("drag", ARROWS[dir], color);
    }

    aim(wanted(state), nowMs);
    place(nowMs);
    paint();
    tellLegend();
  }

  // Förklaringen byggs bara om när tecknen på kartan har ändrats.
  const LEGEND_ORDER = ["du", "vandrare", "borta", "skymd", "forlorad", "drag", "varelse", "falla",
    "kammare", "vaken", "port", "foremal", "krita", "eko", "okand", "mitten"];
  function tellLegend() {
    if (!onLegend) return;
    const rank = (e) => LEGEND_ORDER.indexOf(e.key.split(":")[0]);
    const list = [...legend.values()].sort((a, b) => rank(a) - rank(b));
    const sig = list.map((e) => e.key + e.glyph + e.color + e.name).join("|");
    if (sig === legendSig) return;
    legendSig = sig;
    onLegend(list);
  }
  function setLegend(fn) {
    onLegend = fn;
    legendSig = "";
  }

  // Cellens mitt på tavlan, med ett litet darr så att strecket ser handritat
  // ut. Darret beror bara på cellen, så att strecken möts i hörnen.
  function px(x, y, shake) {
    const h = hash(x, y);
    const jx = shake ? ((h % 997) / 997 - 0.5) * cellW * 0.2 : 0;
    const jy = shake ? (((h >>> 10) % 997) / 997 - 0.5) * cellH * 0.12 : 0;
    return [(x - ox + 0.5) * cellW + jx, (y - oy + 0.5) * cellH + jy];
  }

  function paint() {
    const W = boxW;
    const H = boxH;
    ctx.fillStyle = SLATE;
    ctx.fillRect(0, 0, W, H);
    const look = LOOK[lvl];

    // Golvet i kända rum är en svag kritdimma.
    ctx.fillStyle = look.floor;
    floors.forEach(([x0, y0]) => {
      const [ax, ay] = px(x0, y0, false);
      const [bx, by] = px(x0 + RW, y0 + RH, false);
      ctx.fillRect(ax, ay, bx - ax, by - ay);
    });

    // Väggarna samlas i en bana per nyans, så att det går fort att rita.
    const base = Math.max(1.4, cellW * 0.3);
    const groups = new Map();
    function group(x, y) {
      const own = wallc[y][x];
      if (own) {
        if (!groups.has(own)) groups.set(own, { color: own, lw: base * 1.3, main: new Path2D(), dust: new Path2D() });
        return groups.get(own);
      }
      const q = Math.round(pulse(x, y) * STEPS);
      if (!groups.has(q)) {
        const lw = lvl >= 2 ? base * (1 + (0.5 * q) / STEPS) : base;
        groups.set(q, { color: mix(look.wall, look.hi, q / STEPS), lw, q, main: new Path2D(), dust: new Path2D() });
      }
      return groups.get(q);
    }
    function stroke(g, x1, y1, x2, y2) {
      const [ax, ay] = px(x1, y1, true);
      const [bx, by] = px(x2, y2, true);
      // Ibland tar kritan inte, och strecket får ett litet glapp.
      if (hash(x1 * 3 + x2, y1 * 7 + y2) % 9 === 0) {
        g.main.moveTo(ax, ay);
        g.main.lineTo(ax + (bx - ax) * 0.4, ay + (by - ay) * 0.4);
        g.main.moveTo(ax + (bx - ax) * 0.62, ay + (by - ay) * 0.62);
        g.main.lineTo(bx, by);
      } else {
        g.main.moveTo(ax, ay);
        g.main.lineTo(bx, by);
      }
      // Kritdamm: ett tunnare streck lite vid sidan om.
      const h = hash(x1 + 11, y1 + 5);
      const d = ((h % 5) - 2) * 0.35;
      g.dust.moveTo(ax + d, ay - d);
      g.dust.lineTo(bx - d, by + d);
    }
    for (let y = 0; y < rows; y++) {
      for (let x = 0; x < cols; x++) {
        if (!wall[y][x]) continue;
        const g = group(x, y);
        const right = isWall(x + 1, y);
        const down = isWall(x, y + 1);
        if (right) stroke(g, x, y, x + 1, y);
        if (down) stroke(g, x, y, x, y + 1);
        if (!right && !down && !isWall(x - 1, y) && !isWall(x, y - 1)) stroke(g, x, y, x, y);
        // Revben: tvärstreck längs väggarna mellan hjärtslagen.
        if (lvl >= 3 && g.q !== undefined && g.q < STEPS * 0.75) {
          const along = right || isWall(x - 1, y);
          if ((along ? x : y) % 2 === 0) {
            const [cx, cy] = px(x, y, true);
            const len = cellH * 0.32;
            if (along) {
              g.main.moveTo(cx, cy - len);
              g.main.lineTo(cx, cy + len);
            } else {
              g.main.moveTo(cx - len, cy);
              g.main.lineTo(cx + len, cy);
            }
          }
        }
      }
    }
    ctx.lineCap = "round";
    ctx.lineJoin = "round";
    groups.forEach((g) => {
      ctx.strokeStyle = g.color;
      ctx.globalAlpha = 0.9;
      ctx.lineWidth = g.lw;
      ctx.stroke(g.main);
      ctx.globalAlpha = 0.35;
      ctx.lineWidth = Math.max(0.8, g.lw * 0.4);
      ctx.stroke(g.dust);
    });
    ctx.globalAlpha = 1;

    // Kammarens ring.
    rings.forEach((r) => {
      const [cx, cy] = px(r.x0 + RW / 2, r.y0 + RH / 2, false);
      ctx.beginPath();
      ctx.ellipse(cx, cy, cellW * 2.6, cellH * 1.15, 0, 0, Math.PI * 2);
      ctx.strokeStyle = r.awake ? (r.flash ? "#ff5a3a" : "#a8402c") : look.ring;
      ctx.lineWidth = r.awake ? base * 1.5 : base * 0.8;
      ctx.setLineDash(r.awake ? [] : [cellW * 0.9, cellW * 0.35]);
      ctx.stroke();
      ctx.setLineDash([]);
    });

    // Vandrarna.
    const rad = Math.max(2.5, Math.min(cellW, cellH) * 0.55);
    figures.forEach((f) => {
      const [cx, cy] = px(f.x, f.y, false);
      ctx.beginPath();
      ctx.arc(cx, cy, rad, 0, Math.PI * 2);
      if (f.kind === "ghost") {
        ctx.strokeStyle = f.color;
        ctx.globalAlpha = 0.4;
        ctx.lineWidth = 1.2;
        ctx.stroke();
        ctx.globalAlpha = 1;
        return;
      }
      ctx.fillStyle = f.color;
      ctx.fill();
      if (f.kind === "me") {
        ctx.beginPath();
        ctx.arc(cx, cy, rad + Math.max(2, rad * 0.6), 0, Math.PI * 2);
        ctx.strokeStyle = COLORS.chalk;
        ctx.lineWidth = Math.max(1.2, rad * 0.35);
        ctx.stroke();
      }
    });

    // Tecknen ovanpå.
    ctx.font = "bold " + fontPx + "px Menlo, Monaco, monospace";
    ctx.textAlign = "center";
    ctx.textBaseline = "middle";
    const xa = Math.max(0, Math.floor(ox)), xb = Math.min(cols, Math.ceil(ox + W / cellW));
    const ya = Math.max(0, Math.floor(oy)), yb = Math.min(rows, Math.ceil(oy + H / cellH));
    for (let y = ya; y < yb; y++) {
      for (let x = xa; x < xb; x++) {
        const ch = grid[y][x];
        if (!ch || ch === " ") continue;
        ctx.fillStyle = colr[y][x];
        const [cx, cy] = px(x, y, false);
        ctx.fillText(ch, cx, cy);
      }
    }
  }

  // toggleWhole växlar mellan utsnittet runt vandraren och allt den
  // känner till, när det är för mycket för att läsa på en gång.
  function toggleWhole() {
    wantWhole = !wantWhole;
    if (lastState) draw(lastState);
  }

  function wake() {
    wokeAt = performance.now();
  }

  return { setup, draw, resize, ping, wake, toggleWhole, setLegend, PLAYER_COLORS };
})();
