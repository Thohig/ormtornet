// Minikartan: allt vandraren har upptäckt, ritat platt och vridet som
// kartan snett ovanifrån, så att norr pekar snett uppåt höger (↗) och
// tangenterna vid dörrarna pekar åt samma håll som på minikartan.
//
// På datorn står den uppe till vänster på kartan; ett klick växlar till
// zoomen Allt och tillbaka, och K döljer och visar den. På telefonen står
// den nere till vänster i kompassens ställe, och ett tryck gör den stor
// över kartan tills man trycker igen.
//
// Den visar bara det vandraren själv vet: rummen på den egna kartan
// (state.rooms), dörrarna, portarna, kammaren och var man själv står.
// Kartritaren ser också de rum andra har hittat (a i rummet), blekare, och
// var någon har dött (dod).
"use strict";

const Minikarta = (() => {
  const NS = "http://www.w3.org/2000/svg";
  const KEY = "ormtornet.minikarta";
  const SQUASH = 0.55; // som kartan snett ovanifrån, där golvet ses från omkring 33°
  const MIN_SPAN = 3; // så nära går minikartan som mest, i rum
  const DX = [0, 1, 0, -1]; // N, Ö, S, V
  const DY = [-1, 0, 1, 0];
  const narrow = window.matchMedia("(max-width: 700px)");

  let box = null; // #minikarta
  let svg = null;
  let on = load();
  let big = false; // stor över kartan, på telefonen
  let lastState = null, lastView = null, lastColor = "";

  function load() {
    try {
      return localStorage.getItem(KEY) !== "av";
    } catch {
      return true;
    }
  }

  function el(name, attrs, parent) {
    const e = document.createElementNS(NS, name);
    for (const k in attrs) e.setAttribute(k, attrs[k]);
    if (parent) parent.appendChild(e);
    return e;
  }

  function setup() {
    const mapbox = document.getElementById("mapbox");
    if (!mapbox || box) return;
    box = document.createElement("div");
    box.id = "minikarta";
    box.hidden = true;
    box.setAttribute("role", "button");
    box.tabIndex = -1;
    svg = el("svg", { "aria-hidden": "true" }, box);
    mapbox.appendChild(box);
    texts();
    box.addEventListener("click", (e) => {
      e.stopPropagation();
      if (narrow.matches) {
        big = !big;
        redraw();
      } else MapView.toggleWhole();
    });
    narrow.addEventListener("change", () => {
      big = false;
      redraw();
    });
    window.addEventListener("resize", () => redraw());
  }

  function texts() {
    if (!box || typeof t !== "function") return;
    box.setAttribute("aria-label", t("minikarta_namn"));
    box.title = narrow.matches ? "" : t("minikarta_titel");
  }

  // toggle döljer och visar minikartan på datorn (tangenten K) och
  // kommer ihåg valet på den här enheten. På telefonen stänger den bara
  // den stora minikartan.
  function toggle() {
    if (narrow.matches) {
      big = false;
    } else {
      on = !on;
      try {
        localStorage.setItem(KEY, on ? "pa" : "av");
      } catch {}
    }
    redraw();
  }

  function close() {
    if (!big) return false;
    big = false;
    redraw();
    return true;
  }

  // draw ritar minikartan för tillståndet state. view är rutan kartan
  // visar, i rum (x, y, bredd, höjd), eller null. color är vandrarens färg.
  function draw(state, view, color) {
    lastState = state;
    lastView = view;
    lastColor = color || lastColor;
    redraw();
  }

  function redraw() {
    if (!box) setup();
    if (!box) return;
    const state = lastState;
    const me = state && state.you && state.players.find((p) => p.id === state.you);
    const show = !!(state && state.phase === "spel" && me && me.status === "inne" && (on || narrow.matches));
    box.hidden = !show;
    if (!show) big = false;
    box.classList.toggle("stor", big);
    document.body.classList.toggle("minikarta", show);
    if (!show) return;
    texts();
    render(state, me);
  }

  // render ritar rummen i en romb. Ett rum (x, y) har hörnen (x, y) till
  // (x + 1, y + 1), och en punkt (x, y) står på (x − y) · a åt höger och
  // (x + y) · a · SQUASH nedåt, så att öster pekar snett nedåt höger och
  // norr snett uppåt höger, som på kartan.
  function render(state, me) {
    const W = box.clientWidth, H = box.clientHeight;
    if (!W || !H) return;
    svg.setAttribute("viewBox", `0 0 ${W} ${H}`);
    svg.setAttribute("width", W);
    svg.setAttribute("height", H);
    while (svg.firstChild) svg.firstChild.remove();

    const rooms = state.rooms.filter((r) => r.x >= 0 && r.y >= 0);
    let x0 = Infinity, y0 = Infinity, x1 = -Infinity, y1 = -Infinity;
    const take = (x, y) => {
      x0 = Math.min(x0, x);
      y0 = Math.min(y0, y);
      x1 = Math.max(x1, x + 1);
      y1 = Math.max(y1, y + 1);
    };
    rooms.forEach((r) => take(r.x, r.y));
    if (me.x >= 0) take(me.x, me.y);
    if (x0 === Infinity) return;
    // Lite luft runt omkring, och aldrig närmare än MIN_SPAN rum.
    const grow = (a, b) => {
      const pad = Math.max(0, MIN_SPAN - (b - a)) / 2 + 0.6;
      return [a - pad, b + pad];
    };
    [x0, x1] = grow(x0, x1);
    [y0, y1] = grow(y0, y1);

    // Rombens bredd är (bredd + höjd) · a och dess höjd samma gånger SQUASH.
    const m = big ? 18 : 8;
    const span = x1 - x0 + (y1 - y0);
    const a = Math.min((W - 2 * m) / span, (H - 2 * m) / (span * SQUASH), big ? 64 : 20);
    const cx = (x0 + x1) / 2, cy = (y0 + y1) / 2;
    const P = (x, y) => [W / 2 + (x - cx - (y - cy)) * a, H / 2 + (x - cx + (y - cy)) * a * SQUASH];
    const pts = (list) => list.map(([x, y]) => P(x, y).map((n) => n.toFixed(1)).join(",")).join(" ");
    const line = (ax, ay, bx, by, cls) => {
      const [p, q] = [P(ax, ay), P(bx, by)];
      el("line", { x1: p[0].toFixed(1), y1: p[1].toFixed(1), x2: q[0].toFixed(1), y2: q[1].toFixed(1), class: cls }, svg);
    };

    // Rutan kartan visar, som en streckad ram under rummen.
    if (lastView && !big) {
      const v = lastView;
      const vx0 = Math.max(v.x, x0), vy0 = Math.max(v.y, y0);
      const vx1 = Math.min(v.x + v.w, x1), vy1 = Math.min(v.y + v.h, y1);
      if (vx1 > vx0 && vy1 > vy0) el("polygon", { points: pts([[vx0, vy0], [vx1, vy0], [vx1, vy1], [vx0, vy1]]), class: "mk-vy" }, svg);
    }

    const g = 0.12; // mellanrummet mellan rummen
    rooms.forEach((r) => {
      const cls = ["mk-rum"];
      if (r.a) cls.push("mk-andras");
      if (r.k === "kammare") cls.push("mk-kammare");
      el("polygon", { points: pts([[r.x + g, r.y + g], [r.x + 1 - g, r.y + g], [r.x + 1 - g, r.y + 1 - g], [r.x + g, r.y + 1 - g]]), class: cls.join(" ") }, svg);
    });
    // Dörrarna som streck från mitten ut över kanten, så att en dörr till
    // ett okänt rum syns som en väg man inte har gått.
    rooms.forEach((r) => {
      const mx = r.x + 0.5, my = r.y + 0.5;
      (r.d || []).forEach((open, d) => {
        if (open) line(mx + DX[d] * 0.3, my + DY[d] * 0.3, mx + DX[d] * 0.62, my + DY[d] * 0.62, "mk-dorr" + (r.a ? " mk-andras" : ""));
      });
      if (r.exit >= 0) {
        // Porten ut: ett streck ut ur tornet med en tvärslå.
        const d = r.exit, px = -DY[d], py = DX[d];
        const ex = mx + DX[d] * 0.62, ey = my + DY[d] * 0.62;
        line(mx + DX[d] * 0.3, my + DY[d] * 0.3, ex, ey, "mk-port" + (r.stangd ? " mk-stangd" : ""));
        line(ex - px * 0.22, ey - py * 0.22, ex + px * 0.22, ey + py * 0.22, "mk-port" + (r.stangd ? " mk-stangd" : ""));
      }
      if (r.dod) {
        const s = 0.16;
        line(mx - s, my - s, mx + s, my + s, "mk-dod");
        line(mx + s, my - s, mx - s, my + s, "mk-dod");
      }
    });

    // Vandraren själv, i sin färg.
    if (me.x >= 0) {
      const [px, py] = P(me.x + 0.5, me.y + 0.5);
      const c = el("circle", { cx: px.toFixed(1), cy: py.toFixed(1), r: Math.max(2.6, a * 0.24).toFixed(1), class: "mk-du" }, svg);
      if (lastColor) c.style.fill = lastColor;
    }

    // N i hörnet åt det håll norr ligger.
    const nx = W - (big ? 18 : 9), ny = big ? 20 : 11;
    const n = el("text", { x: nx, y: ny, class: "mk-n", "text-anchor": "middle" }, svg);
    n.textContent = "N";
    const ul = Math.hypot(1, SQUASH), ux = 1 / ul, uy = -SQUASH / ul;
    const tx = nx - 8, ty = ny - 3; // pilens spets, strax till vänster om N
    const bx = tx - ux * 9, by = ty - uy * 9;
    const h1 = [tx - ux * 3.5 - uy * 2.5, ty - uy * 3.5 + ux * 2.5], h2 = [tx - ux * 3.5 + uy * 2.5, ty - uy * 3.5 - ux * 2.5];
    const f = (n) => n.toFixed(1);
    el("path", { d: `M${f(bx)} ${f(by)}L${f(tx)} ${f(ty)}M${f(h1[0])} ${f(h1[1])}L${f(tx)} ${f(ty)}L${f(h2[0])} ${f(h2[1])}`, class: "mk-npil" }, svg);
  }

  return { setup, draw, toggle, close };
})();
