// Ritar tornet som ett rutnät av tecken på en canvas. Varje tecken placeras
// i en egen cell, så att braille och andra symboler alltid hamnar rätt även
// om typsnittet hämtar dem från ett annat teckensnitt.
"use strict";

const MapView = (() => {
  const RW = 7; // tecken per rum på bredden (inklusive en vägg)
  const RH = 4; // rader per rum (inklusive en vägg)

  const COLORS = {
    wall: "#b39663",
    floor: "#4a3f2c",
    unknown: "#2c261b",
    exit: "#8a7a55",
    chamber: "#c9a24f",
    gold: "#ffd36b",
    lost: "#8a3b30",
    mark: "#c0503e",
    monster: "#e0806a",
  };
  const PLAYER_COLORS = ["#ffcc66", "#66e0ff", "#ff7a9a", "#9dff7a"];

  const WALL = ["⣿", "⣷", "⣾", "⡿", "⢿", "⣯", "⣟", "⣻", "⣽"];
  const ARROWS = ["▲", "▶", "▼", "◀"];
  const FLOOR = ["⠂", "⠄", "⠈", "⠐", "⠠"];

  // Kammaren: en ring kring en hög av guld, 6 x 3 tecken.
  const CHAMBER = [
    " ⣠⠶⠶⣄ ",
    "⢸ ¤¤ ⡇",
    " ⠙⠶⠶⠋ ",
  ];
  // Samma ring när ormen har vaknat: slingorna spänns och två ögon lyser.
  const CHAMBER_AWAKE = [
    "⣴⠟⠛⠛⠻⣦",
    "⣿ ◉◉ ⣿",
    "⠻⣦⣤⣤⣴⠟",
  ];
  const WAKE_MS = 3500;
  let wokeAt = -Infinity;

  // Avslöjandet: tornets utseende följer nivån. På nivå 2 andas väggarna,
  // på nivå 3 syns revbenen och hjärtslagen, på nivå 4 är det kött och tänder.
  const LOOK = {
    1: { wall: [179, 150, 99], hi: [179, 150, 99], floor: "#4a3f2c", period: 0 },
    2: { wall: [178, 146, 100], hi: [228, 190, 136], floor: "#4d3f2e", period: 4200 },
    3: { wall: [176, 118, 92], hi: [240, 166, 132], floor: "#55382c", period: 1250 },
    4: { wall: [160, 70, 52], hi: [255, 118, 88], floor: "#5e2a22", period: 950 },
  };
  // Braille i tre tätheter: väggen sväller och krymper när den andas.
  const THIN = ["⡪", "⢕", "⠪", "⡕", "⢪", "⠕"];
  const DENSE = ["⣿"];
  // Revben: tvärstrimmor längs väggarna.
  const RIB_H = ["⣿", "⡇"];
  const RIB_V = ["⣿", "⠶"];

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

  function wallAt(x, y, horizontal) {
    const look = LOOK[lvl];
    const p = pulse(x, y);
    const color = mix(look.wall, look.hi, p);
    let ch;
    if (lvl >= 3 && p < 0.75) {
      ch = horizontal ? RIB_H[x % 2] : RIB_V[y % RIB_V.length];
    } else if (lvl >= 2 && p > 0.7) {
      ch = DENSE[0];
    } else if (lvl >= 2 && p < 0.25) {
      ch = THIN[hash(x, y) % THIN.length];
    } else {
      ch = WALL[hash(x, y) % WALL.length];
    }
    put(x, y, ch, color);
  }

  // Spelarnas platser inne i ett rum, en per hörnplats.
  const SLOTS = [[2, 1], [5, 1], [2, 3], [5, 3]];

  let canvas, ctx, cols, rows, cellW, cellH, fontPx;
  let grid, colr;

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

  function setup(el, size) {
    canvas = el;
    ctx = canvas.getContext("2d");
    cols = size * RW + 1;
    rows = size * RH + 1;
    window.addEventListener("resize", resize);
    resize();
  }

  function resize() {
    if (!canvas) return;
    const box = canvas.parentElement.getBoundingClientRect();
    const byW = box.width / (cols * 0.62);
    const byH = box.height / (rows * 1.12);
    fontPx = Math.max(8, Math.floor(Math.min(byW, byH)));
    cellW = Math.round(fontPx * 0.62);
    cellH = Math.round(fontPx * 1.12);
    const dpr = window.devicePixelRatio || 1;
    canvas.width = cols * cellW * dpr;
    canvas.height = rows * cellH * dpr;
    canvas.style.width = cols * cellW + "px";
    canvas.style.height = rows * cellH + "px";
    ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
    if (lastState) draw(lastState);
  }

  function drawRoom(r) {
    const x0 = r.x * RW;
    const y0 = r.y * RH;
    const floor = LOOK[lvl].floor;

    // Väggar runt om.
    for (let i = 0; i <= RW; i++) {
      wallAt(x0 + i, y0, true);
      wallAt(x0 + i, y0 + RH, true);
    }
    for (let j = 1; j < RH; j++) {
      wallAt(x0, y0 + j, false);
      wallAt(x0 + RW, y0 + j, false);
    }
    // Golv.
    for (let j = 1; j < RH; j++) {
      for (let i = 1; i < RW; i++) {
        const h = hash(x0 + i, y0 + j);
        const ch = h % 7 === 0 ? FLOOR[h % FLOOR.length] : " ";
        put(x0 + i, y0 + j, ch, floor);
      }
    }
    // Dörrar: en lucka mitt på varje öppen sida.
    const mid = Math.floor(RW / 2);
    const gapN = [x0 + mid, x0 + mid + 1].filter((x) => x < x0 + RW);
    if (r.d[0]) gapN.forEach((x) => put(x, y0, " ", COLORS.floor));
    if (r.d[2]) gapN.forEach((x) => put(x, y0 + RH, " ", COLORS.floor));
    const midY = y0 + Math.floor(RH / 2);
    if (r.d[1]) put(x0 + RW, midY, " ", COLORS.floor);
    if (r.d[3]) put(x0, midY, " ", COLORS.floor);

    // Utgång genom ytterväggen. På sista nivån är porten en mun med tänder,
    // och när spelet är slut har den stängts.
    const ex = r.exit === 1 ? x0 + RW : r.exit === 3 ? x0 : -1;
    if (ex >= 0) {
      if (lvl === 4 && lastState && lastState.phase === "slut") put(ex, midY, "⣿", "#d0503a");
      else if (lvl === 4) put(ex, midY, Math.floor(nowMs / 400) % 2 ? "⩚" : "⩛", "#ff7a5a");
      else put(ex, midY, r.exit === 1 ? "›" : "‹", COLORS.exit);
    }

    if (r.mark) put(x0 + Math.floor(RW / 2), y0 + Math.floor(RH / 2), r.mark, COLORS.mark);

    if (r.k === "kammare") {
      const now = performance.now();
      const awake = now - wokeAt < WAKE_MS;
      const flash = awake && Math.floor((now - wokeAt) / 180) % 2 === 0;
      (awake ? CHAMBER_AWAKE : CHAMBER).forEach((line, j) => {
        [...line].forEach((ch, i) => {
          let color = ch === "¤" ? COLORS.gold : lvl === 4 ? "#c4483a" : lvl === 3 ? "#c98a5a" : COLORS.chamber;
          if (awake) color = ch === "◉" ? "#ffe14a" : flash ? "#ff5a3a" : "#a8402c";
          put(x0 + 1 + i, y0 + 1 + j, ch, color);
        });
      });
    }
  }

  let lastState = null;
  const PING_MS = 4500;
  let pings = [];

  function ping(x, y, slot) {
    pings.push({ x, y, slot, at: performance.now() });
  }

  function draw(state, t) {
    lastState = state;
    if (!ctx) return;
    reveal = !!state.reveal;
    lvl = reveal ? Math.min(4, Math.max(1, state.level || 1)) : 1;
    nowMs = t || performance.now();
    grid = Array.from({ length: rows }, () => new Array(cols).fill(" "));
    colr = Array.from({ length: rows }, () => new Array(cols).fill(COLORS.unknown));

    // Okända rum syns bara som en svag prick.
    const size = state.size;
    const mid = Math.floor(size / 2);
    for (let y = 0; y < size; y++) {
      for (let x = 0; x < size; x++) {
        const ch = x === mid && y === mid ? "◌" : "·";
        put(x * RW + Math.floor(RW / 2), y * RH + Math.floor(RH / 2), ch, COLORS.unknown);
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
      const color = Math.floor((now - p.at) / 250) % 2 ? PLAYER_COLORS[p.slot] : COLORS.wall;
      for (let j = 0; j <= RH; j++) {
        for (let i = 0; i <= RW; i++) {
          const edge = i === 0 || j === 0 || i === RW || j === RH;
          if (edge && grid[y0 + j] && grid[y0 + j][x0 + i] !== " ") colr[y0 + j][x0 + i] = color;
        }
      }
      put(x0 + Math.floor(RW / 2), y0 + Math.floor(RH / 2), "◎", PLAYER_COLORS[p.slot]);
    });
    state.players.forEach((p) => {
      if (p.status === "ute" || p.hidden) return;
      const [sx, sy] = SLOTS[p.slot] || SLOTS[0];
      let color = PLAYER_COLORS[p.slot] || "#fff";
      let ch = "@";
      if (!p.online) color = COLORS.floor;
      if (p.status === "forlorad" || p.status === "dod") {
        ch = "x";
        color = COLORS.lost;
      } else if (p.id === state.you && !blink) {
        ch = "◉";
      }
      put(p.x * RW + sx, p.y * RH + sy, ch, color);
    });

    // Varelserna. Den som styrs av en död spelare blinkar när den är vald.
    (state.monsters || []).forEach((m) => {
      const selected = state.controller && m.id === selectedMonster;
      const color = selected && !blink ? "#ffe14a" : COLORS.monster;
      put(m.x * RW + Math.floor(RW / 2) + 1, m.y * RH + Math.floor(RH / 2), m.glyph, color);
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
    }

    ctx.clearRect(0, 0, cols * cellW, rows * cellH);
    ctx.font = fontPx + "px Menlo, Monaco, monospace";
    ctx.textAlign = "center";
    ctx.textBaseline = "middle";
    for (let y = 0; y < rows; y++) {
      for (let x = 0; x < cols; x++) {
        const ch = grid[y][x];
        if (ch === " ") continue;
        ctx.fillStyle = colr[y][x];
        ctx.fillText(ch, x * cellW + cellW / 2, y * cellH + cellH / 2);
      }
    }
  }

  function wake() {
    wokeAt = performance.now();
  }

  return { setup, draw, resize, ping, wake, PLAYER_COLORS };
})();
