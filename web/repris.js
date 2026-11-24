// Reprisen: vandringen spelad igen på slutskärmen (replay.go,
// docs/repris.md). Tornet som det blev ritas platt och vridet som
// minikartan, så att norr pekar snett uppåt höger, och allas vägar växer
// fram runda för runda som linjer i var och ens färg. Under kartan står
// vandrarna med guldet de bar just då.
//
//   Repris.open(repris)  öppnar reprisen (state.repris) och spelar den
//   Repris.close()       stänger den; svarar sant om den var öppen
"use strict";

const Repris = (() => {
  const NS = "http://www.w3.org/2000/svg";
  const SQUASH = 0.55; // som minikartan
  const DX = [0, 1, 0, -1]; // N, Ö, S, V
  const DY = [-1, 0, 1, 0];
  const STEP_MS = 650; // en runda
  // Var och ens linje går lite vid sidan av rummets mitt, så att vägar
  // som går genom samma rum syns bredvid varandra.
  const SPREAD = 0.13;

  let box = null, svg = null, legend = null, slider = null, label = null, playBtn = null;
  let data = null, at = 0, timer = null;

  function el(name, attrs, parent) {
    const e = document.createElementNS(NS, name);
    for (const k in attrs) e.setAttribute(k, attrs[k]);
    if (parent) parent.appendChild(e);
    return e;
  }

  function button(cls, key, onclick) {
    const b = document.createElement("button");
    b.type = "button";
    b.className = cls;
    b.dataset.key = key;
    b.addEventListener("click", onclick);
    return b;
  }

  function setup() {
    if (box) return;
    box = document.createElement("div");
    box.id = "repris";
    box.hidden = true;
    box.setAttribute("role", "dialog");
    const panel = document.createElement("div");
    panel.className = "rp-panel";
    const head = document.createElement("div");
    head.className = "rp-head";
    const h = document.createElement("h2");
    h.dataset.key = "repris_rubrik";
    head.append(h, button("rp-stang", "repris_stang", close));
    svg = el("svg", { "aria-hidden": "true", class: "rp-karta" });
    const controls = document.createElement("div");
    controls.className = "rp-styr";
    playBtn = button("rp-spela", "repris_spela", () => (timer ? pause() : play()));
    slider = document.createElement("input");
    slider.type = "range";
    slider.min = "0";
    slider.step = "1";
    slider.addEventListener("input", () => {
      pause();
      show(Number(slider.value));
    });
    label = document.createElement("span");
    label.className = "rp-runda";
    controls.append(button("rp-borjan", "repris_borjan", () => {
      pause();
      show(0);
    }), playBtn, slider, label);
    legend = document.createElement("div");
    legend.className = "rp-vandrare";
    panel.append(head, svg, controls, legend);
    box.appendChild(panel);
    box.addEventListener("click", (e) => {
      if (e.target === box) close();
    });
    document.body.appendChild(box);
    window.addEventListener("resize", () => {
      if (!box.hidden) draw();
    });
  }

  // texts sätter rubriken och knapparnas förklaringar på språket som
  // gäller. Knapparna själva är tecken.
  const SYMBOL = { repris_stang: "✕", repris_spela: "▶", repris_paus: "❚❚", repris_borjan: "⏮" };
  function texts() {
    box.querySelectorAll("[data-key]").forEach((e) => {
      const k = e === playBtn && timer ? "repris_paus" : e.dataset.key;
      if (e.tagName === "H2") {
        e.textContent = t(k);
        return;
      }
      e.textContent = SYMBOL[k];
      e.title = t(k);
      e.setAttribute("aria-label", t(k));
    });
    slider.setAttribute("aria-label", t("repris_reglage"));
  }

  function open(repris) {
    if (!repris || !repris.rundor || !repris.rundor.length) return;
    setup();
    data = repris;
    box.hidden = false;
    slider.max = String(data.rundor.length - 1);
    show(0);
    play();
  }

  function close() {
    if (!box || box.hidden) return false;
    pause();
    box.hidden = true;
    return true;
  }

  function play() {
    if (timer) return;
    if (at >= data.rundor.length - 1) show(0);
    timer = setInterval(() => {
      if (at >= data.rundor.length - 1) {
        pause();
        return;
      }
      show(at + 1);
    }, STEP_MS);
    texts();
  }

  function pause() {
    clearInterval(timer);
    timer = null;
    if (box) texts();
  }

  function show(i) {
    at = Math.max(0, Math.min(i, data.rundor.length - 1));
    slider.value = String(at);
    const last = data.rundor.length - 1;
    label.textContent = at === 0 ? t("repris_fore") : at === last ? t("repris_slut") : t("repris_runda", { runda: at });
    texts();
    draw();
    drawLegend();
  }

  // draw ritar tornet och vägarna fram till rundan at. Ett rum (x, y)
  // har hörnen (x, y) till (x + 1, y + 1), och en punkt (x, y) står på
  // (x − y) · a åt höger och (x + y) · a · SQUASH nedåt, som minikartan.
  function draw() {
    const W = svg.clientWidth, H = svg.clientHeight;
    if (!W || !H) return;
    svg.setAttribute("viewBox", `0 0 ${W} ${H}`);
    while (svg.firstChild) svg.firstChild.remove();
    let x0 = Infinity, y0 = Infinity, x1 = -Infinity, y1 = -Infinity;
    data.rum.forEach((r) => {
      x0 = Math.min(x0, r.x);
      y0 = Math.min(y0, r.y);
      x1 = Math.max(x1, r.x + 1);
      y1 = Math.max(y1, r.y + 1);
    });
    if (x0 === Infinity) return;
    x0 -= 0.6; y0 -= 0.6; x1 += 0.6; y1 += 0.6;
    const m = 14;
    const span = x1 - x0 + (y1 - y0);
    const a = Math.min((W - 2 * m) / span, (H - 2 * m) / (span * SQUASH));
    const cx = (x0 + x1) / 2, cy = (y0 + y1) / 2;
    const P = (x, y) => [W / 2 + (x - cx - (y - cy)) * a, H / 2 + (x - cx + (y - cy)) * a * SQUASH];
    const f = (n) => n.toFixed(1);
    const pts = (list) => list.map(([x, y]) => P(x, y).map(f).join(",")).join(" ");
    const line = (ax, ay, bx, by, cls, color) => {
      const [p, q] = [P(ax, ay), P(bx, by)];
      const l = el("line", { x1: f(p[0]), y1: f(p[1]), x2: f(q[0]), y2: f(q[1]), class: cls }, svg);
      if (color) l.style.stroke = color;
      return l;
    };

    const g = 0.12;
    data.rum.forEach((r) => {
      el("polygon", { points: pts([[r.x + g, r.y + g], [r.x + 1 - g, r.y + g], [r.x + 1 - g, r.y + 1 - g], [r.x + g, r.y + 1 - g]]), class: "rp-rum" + (r.k === "kammare" ? " rp-kammare" : "") }, svg);
    });
    data.rum.forEach((r) => {
      const mx = r.x + 0.5, my = r.y + 0.5;
      (r.d || []).forEach((open, d) => {
        if (open) line(mx + DX[d] * 0.3, my + DY[d] * 0.3, mx + DX[d] * 0.62, my + DY[d] * 0.62, "rp-dorr");
      });
      if (r.exit >= 0) {
        const d = r.exit, px = -DY[d], py = DX[d];
        const ex = mx + DX[d] * 0.62, ey = my + DY[d] * 0.62;
        line(mx + DX[d] * 0.3, my + DY[d] * 0.3, ex, ey, "rp-port");
        line(ex - px * 0.22, ey - py * 0.22, ex + px * 0.22, ey + py * 0.22, "rp-port");
      }
    });

    // Vägarna: från rum till rum, runda för runda, i var och ens färg. I
    // valven syns man inte, och vägen fortsätter där man kom upp.
    const n = data.vandrare.length;
    data.vandrare.forEach((v, i) => {
      const color = MapView.PLAYER_COLORS[v.slot] || "#fff";
      const ang = (i / Math.max(1, n)) * Math.PI * 2;
      const ox = n > 1 ? Math.cos(ang) * SPREAD : 0, oy = n > 1 ? Math.sin(ang) * SPREAD : 0;
      const C = (s) => [s.x + 0.5 + ox, s.y + 0.5 + oy];
      let prev = null;
      for (let k = 0; k <= at; k++) {
        const s = data.rundor[k][i];
        if (s.v || s.x < 0) {
          prev = null;
          continue;
        }
        if (prev && (prev.x !== s.x || prev.y !== s.y)) {
          const [ax, ay] = C(prev), [bx, by] = C(s);
          line(ax, ay, bx, by, "rp-vag" + (k === at ? " rp-ny" : ""), color);
        }
        prev = s;
      }
      // Var vandraren är nu, eller hur det slutade.
      const now = data.rundor[at][i];
      let last = now;
      for (let k = at; k >= 0 && (last.v || last.x < 0); k--) last = data.rundor[k][i];
      if (last.x < 0) return;
      const [px, py] = P(...C(last));
      const r = Math.max(3, a * 0.17);
      if (now.s === "dod") {
        const s = r * 0.9;
        el("path", { d: `M${f(px - s)} ${f(py - s)}L${f(px + s)} ${f(py + s)}M${f(px + s)} ${f(py - s)}L${f(px - s)} ${f(py + s)}`, class: "rp-dod" }, svg).style.stroke = color;
      } else {
        const c = el("circle", { cx: f(px), cy: f(py), r: f(r), class: "rp-du" + (now.s === "ute" ? " rp-ute" : "") + (now.s === "forlorad" || now.v ? " rp-svag" : "") }, svg);
        if (now.s === "ute") c.style.stroke = color;
        else c.style.fill = color;
      }
    });
  }

  // drawLegend skriver vandrarna i sin färg med guldet de bar i rundan at,
  // och hur det har gått.
  function drawLegend() {
    legend.innerHTML = "";
    data.vandrare.forEach((v, i) => {
      const s = data.rundor[at][i];
      const d = document.createElement("span");
      d.className = "rp-namn";
      const name = document.createElement("b");
      name.textContent = v.namn;
      name.style.color = MapView.PLAYER_COLORS[v.slot] || "";
      const how = { ute: "repris_ute", dod: "repris_dod", forlorad: "repris_forlorad" }[s.s] || (s.v ? "repris_valven" : "");
      d.append(name, " " + t("repris_guld", { guld: s.g }) + (how ? " · " + t(how) : ""));
      legend.appendChild(d);
    });
  }

  return { open, close };
})();
