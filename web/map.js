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
    const wall = (x, y) => put(x, y, WALL[hash(x, y) % WALL.length], COLORS.wall);

    // Väggar runt om.
    for (let i = 0; i <= RW; i++) {
      wall(x0 + i, y0);
      wall(x0 + i, y0 + RH);
    }
    for (let j = 1; j < RH; j++) {
      wall(x0, y0 + j);
      wall(x0 + RW, y0 + j);
    }
    // Golv.
    for (let j = 1; j < RH; j++) {
      for (let i = 1; i < RW; i++) {
        const h = hash(x0 + i, y0 + j);
        const ch = h % 7 === 0 ? FLOOR[h % FLOOR.length] : " ";
        put(x0 + i, y0 + j, ch, COLORS.floor);
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

    // Utgång genom ytterväggen.
    if (r.exit === 1) put(x0 + RW, midY, "›", COLORS.exit);
    if (r.exit === 3) put(x0, midY, "‹", COLORS.exit);

    if (r.mark) put(x0 + Math.floor(RW / 2), y0 + Math.floor(RH / 2), r.mark, COLORS.mark);

    if (r.k === "kammare") {
      CHAMBER.forEach((line, j) => {
        [...line].forEach((ch, i) => {
          const color = ch === "¤" ? COLORS.gold : COLORS.chamber;
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

  return { setup, draw, resize, ping, PLAYER_COLORS };
})();
