// Spelkortens bilder (docs/spelkort.md). Varje motiv är en liten SVG i
// samma stil som ormkorten: några få former i jordfärger på pergament,
// utan text, så att bilderna är desamma på båda språken. Servern säger
// motivet (bild i data/ eller id:t); ett motiv som saknas här ritas som
// facklan.
//
//   Kort.bild(motiv, { svag })  SVG-text för kortets bild
//   Kort.mini(motiv)            ett litet kort att ställa framför en rad

const Kort = (() => {
  const f = (n) => n.toFixed(1);
  // Färgerna på pergamentet.
  const BLAECK = "#2a1e12";
  const BRUN = "#6b4a24";
  const STEN = "#8a7f6e";
  const MORK_STEN = "#5f564a";
  const JARN = "#4a4844";
  const GULD = "#d9a63a";
  const GRON = "#5b6b34";
  const ROD = "#8a2a18";
  const svg = (inner) => `<svg viewBox="0 0 100 100" aria-hidden="true">${inner}</svg>`;

  // En vägg av stenar i bakgrunden, svag så att motivet syns.
  const vagg = (y0 = 6, y1 = 94) => {
    let s = "";
    for (let y = y0, rad = 0; y < y1; y += 14, rad++) {
      for (let x = rad % 2 ? -8 : 6; x < 94; x += 28) {
        s += `<rect x="${Math.max(6, x)}" y="${y}" width="${Math.min(26, 94 - Math.max(6, x))}" height="12" rx="2" fill="${STEN}" opacity="0.18"/>`;
      }
    }
    return s;
  };
  // Ett golv i perspektiv, som en ljus fläck.
  const golv = (cy = 78) => `<ellipse cx="50" cy="${cy}" rx="42" ry="14" fill="${STEN}" opacity="0.25"/>`;
  // En låga, med mitten i (x, y) och höjden h.
  const laga = (x, y, h) =>
    `<path d="M${x} ${y - h}q${h * 0.45} ${h * 0.5} ${h * 0.3} ${h * 0.8}q-${h * 0.3} ${h * 0.35} -${h * 0.6} 0q-${h * 0.15} -${h * 0.4} ${h * 0.3} -${h * 0.8}z" fill="#e8a83a"/>` +
    `<path d="M${x} ${y - h * 0.55}q${h * 0.2} ${h * 0.3} ${h * 0.12} ${h * 0.45}q-${h * 0.12} ${h * 0.15} -${h * 0.24} 0q-${h * 0.06} -${h * 0.2} ${h * 0.12} -${h * 0.45}z" fill="#fff0b0"/>`;

  const BILDER = {};

  // ---- Ormen ----

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

  // ormen ritar ormen som sover hoprullad, som rör sig med ett öga på
  // glänt, eller två gula ögon i mörkret.
  function ormen(kind) {
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
      return `<rect x="4" y="4" width="92" height="92" rx="6" fill="#1a0b08"/>
<path d="${d}" fill="none" stroke="#3a1410" stroke-width="11" stroke-linecap="round" opacity="0.8"/>
${[34, 66].map((x) => `<ellipse cx="${x}" cy="46" rx="11" ry="7" fill="#ffd23a"/><ellipse cx="${x}" cy="46" rx="2" ry="6.5" fill="#1a0b08"/>`).join("")}
<circle cx="50" cy="46" r="40" fill="none" stroke="#a8402c" stroke-width="1.5" opacity="0.6"/>`;
    }
    const moves = kind === "ror_sig"
      ? `<path d="M10 34q6 -6 12 0M78 30q6 -6 12 0M84 80q4 -5 8 0" fill="none" stroke="${BRUN}" stroke-width="1.6" stroke-linecap="round"/>`
      : `<text x="78" y="28" font-size="9" fill="${BRUN}" font-family="serif">z</text><text x="84" y="20" font-size="7" fill="${BRUN}" font-family="serif">z</text>`;
    return `<ellipse cx="50" cy="62" rx="40" ry="22" fill="#c9a24f" opacity="0.35"/>
${[[38, 70], [56, 74], [66, 64], [30, 60]].map(([x, y]) => `<ellipse cx="${x}" cy="${y}" rx="3.2" ry="1.6" fill="${GULD}" stroke="#7a5a1a" stroke-width="0.6"/>`).join("")}
<path d="${d}" fill="none" stroke="#1d2410" stroke-width="11.5" stroke-linecap="round"/>
<path d="${d}" fill="none" stroke="#5b6b34" stroke-width="9" stroke-linecap="round"/>
<path d="${d}" fill="none" stroke="#8a9a52" stroke-width="2" stroke-linecap="round" stroke-dasharray="1.5 3"/>
${head}
${moves}`;
  }
  BILDER.sover = () => ormen("sover");
  BILDER.ror_sig = () => ormen("ror_sig");
  BILDER.vaknar = () => ormen("vaknar");

  // ---- Facklan: rundan utan eget motiv ----

  BILDER.fackla = () => `${vagg()}
<path d="M44 40h12l-3 8h-6z" fill="${JARN}"/>
<rect x="47" y="46" width="6" height="40" rx="2" fill="${BRUN}"/>
<path d="M40 52h20" stroke="${JARN}" stroke-width="3" stroke-linecap="round"/>
<circle cx="50" cy="30" r="20" fill="#e8a83a" opacity="0.15"/>
${laga(50, 41, 26)}`;

  const ORM = ["sover", "ror_sig", "vaknar"];

  return {
    bild(motiv, opt = {}) {
      const rita = BILDER[motiv] || BILDER.fackla;
      const inner = rita();
      return svg(opt.svag ? `<g opacity="0.4" filter="grayscale(1)">${inner}</g>` : inner);
    },
    // mini ger ett litet kort. Ormkorten har egna små bilder i CSS:en;
    // de andra får sin bild i liten skala.
    mini(motiv) {
      const s = document.createElement("span");
      s.setAttribute("aria-hidden", "true");
      if (ORM.includes(motiv)) {
        s.className = "minikort " + motiv;
      } else {
        s.className = "minikort bild";
        s.innerHTML = this.bild(motiv);
      }
      return s;
    },
    // Hjälp åt delarna längre ner i filen.
    _: { f, svg, vagg, golv, laga, BILDER, BLAECK, BRUN, STEN, MORK_STEN, JARN, GULD, GRON, ROD },
  };
})();
