// Spelkortens bilder (docs/spelkort.md). Varje motiv är en liten SVG i
// samma stil som ormkorten: några få former i jordfärger på pergament,
// utan text, så att bilderna är desamma på båda språken. Servern säger
// motivet (bild i data/ eller id:t); ett motiv som saknas här ritas som
// facklan.
//
//   Kort.bild(motiv, { svag, fallen })  SVG-text för kortets bild
//   Kort.mini(motiv, { svag, fallen })  ett litet kort framför en rad

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

  // ---- Fällorna ----

  // Ett galler av spetsiga järnstänger som rasar ner i en dörröppning.
  BILDER.fallgaller = () => `${vagg()}
<path d="M24 92V40a26 26 0 0 1 52 0v52z" fill="#1e1610"/>
${[30, 38, 46, 54, 62, 70].map((x) => `<path d="M${x - 1.6} 14h3.2V74h1.6l-3.2 9 -3.2 -9h1.6z" fill="#9a958c" stroke="${BLAECK}" stroke-width="0.7"/>`).join("")}
<rect x="24" y="22" width="52" height="4" fill="#8a857c" stroke="${BLAECK}" stroke-width="0.6"/><rect x="24" y="50" width="52" height="3" fill="#8a857c" stroke="${BLAECK}" stroke-width="0.6"/>
<path d="M20 10l6 4M80 10l-6 4M50 4v6" stroke="${BRUN}" stroke-width="1.6" stroke-linecap="round"/>`;

  // Ett golv som ger efter: plattor som spricker och faller ner i mörker.
  BILDER.rasande_golv = () => `${vagg(6, 40)}
<path d="M6 46h88v48H6z" fill="${STEN}" opacity="0.35"/>
<path d="M22 56l14 -6 24 2 16 10 -6 18 -22 8 -20 -6 -10 -14z" fill="#140e09"/>
<path d="M30 74l6 -4 4 6zM58 60l7 2 -3 6zM44 84l6 -3 2 5z" fill="${MORK_STEN}"/>
<path d="M6 50l16 6M76 62l18 -6M60 52l10 -6M36 50l-8 -4M76 62l6 16" stroke="${BLAECK}" stroke-width="1.2" fill="none"/>`;

  // Ett golv som vrider sig: en skiva med pilar runt.
  BILDER.roterande = () => `${vagg(6, 36)}
<ellipse cx="50" cy="64" rx="40" ry="20" fill="${STEN}" opacity="0.3"/>
<ellipse cx="50" cy="64" rx="30" ry="14" fill="${MORK_STEN}" opacity="0.55" stroke="${BLAECK}" stroke-width="1"/>
<path d="M50 64l30 0M50 64l-24 8M50 64l-6 -14" stroke="${BLAECK}" stroke-width="0.8" opacity="0.6"/>
<path d="M22 50a32 16 0 0 1 50 -6" fill="none" stroke="${BRUN}" stroke-width="2"/><path d="M72 44l-1 -6 6 3z" fill="${BRUN}"/>
<path d="M78 80a32 16 0 0 1 -50 6" fill="none" stroke="${BRUN}" stroke-width="2"/><path d="M28 86l1 6 -6 -3z" fill="${BRUN}"/>`;

  // Pilar ur små hål i väggen.
  BILDER.pilar = () => `${vagg()}
${[24, 46, 68].map((y) => `<circle cx="12" cy="${y}" r="3" fill="#140e09"/>`).join("")}
${[[24, 62], [46, 82], [68, 50]].map(([y, x]) => `<path d="M${x - 34} ${y}H${x}" stroke="${BRUN}" stroke-width="1.8"/><path d="M${x} ${y}l-6 -3v6z" fill="${JARN}"/><path d="M${x - 34} ${y}l-4 -3M${x - 34} ${y}l-4 3M${x - 30} ${y}l-4 -3M${x - 30} ${y}l-4 3" stroke="${ROD}" stroke-width="1"/><path d="M${x - 40} ${y - 4}h-8M${x - 42} ${y + 4}h-6" stroke="${STEN}" stroke-width="1" opacity="0.6"/>`).join("")}`;

  // En söt, gulgrön dimma som väller fram.
  BILDER.dimma = () => `${vagg()}
${[[30, 70, 18], [56, 64, 22], [76, 74, 14], [44, 50, 14], [64, 46, 12], [20, 54, 10]].map(([x, y, r]) => `<circle cx="${x}" cy="${y}" r="${r}" fill="#b8c24a" opacity="0.45"/>`).join("")}
<path d="M14 82q12 -8 24 0t24 0t24 0" fill="none" stroke="#8a9a32" stroke-width="2" opacity="0.7"/>
<path d="M20 40q8 -6 16 0M60 30q8 -6 16 0" fill="none" stroke="#8a9a32" stroke-width="1.4" opacity="0.6"/>`;

  // En grop full av små ormar.
  BILDER.ormgrop = () => `${vagg(6, 40)}
<ellipse cx="50" cy="66" rx="40" ry="22" fill="${STEN}" opacity="0.3"/>
<ellipse cx="50" cy="66" rx="32" ry="16" fill="#140e09"/>
${[[30, 66], [44, 72], [56, 62], [66, 70], [40, 60], [58, 74]].map(([x, y], i) => `<path d="M${x} ${y}q3 -4 6 0t6 0t6 0" fill="none" stroke="${i % 2 ? "#7a8a3a" : GRON}" stroke-width="2.2" stroke-linecap="round"/><circle cx="${x + 18}" cy="${y}" r="1.6" fill="#e8c040"/>`).join("")}`;

  // Ett väldigt, genomskinligt skinn som har fastnat runt en.
  BILDER.skinn = () => `${vagg()}
<path d="M14 30q20 -14 36 0t36 0v50q-18 12 -36 0t-36 0z" fill="#efe6d0" opacity="0.7" stroke="${STEN}" stroke-width="1.2"/>
${[0, 1, 2, 3].map((r) => [0, 1, 2, 3, 4].map((c) => `<path d="M${20 + c * 14 + (r % 2) * 7} ${38 + r * 11}q4 -5 8 0" fill="none" stroke="${STEN}" stroke-width="0.9" opacity="0.8"/>`).join("")).join("")}
<path d="M50 30v52" stroke="${STEN}" stroke-width="0.8" stroke-dasharray="2 3" opacity="0.6"/>`;

  // En väldig klocka i taket som slår.
  BILDER.klocka = () => `<rect x="6" y="6" width="88" height="8" fill="${MORK_STEN}" opacity="0.5"/>
<path d="M50 14v8" stroke="${JARN}" stroke-width="3"/>
<path d="M34 64q0 -40 16 -42q16 2 16 42l6 6H28z" fill="#a8823a" stroke="${BLAECK}" stroke-width="1.2"/>
<path d="M40 60q0 -26 8 -32" fill="none" stroke="#e0c070" stroke-width="2" opacity="0.6"/>
<circle cx="50" cy="76" r="5" fill="${JARN}"/>
<path d="M16 46q-6 10 0 20M10 40q-9 16 0 32M84 46q6 10 0 20M90 40q9 16 0 32" fill="none" stroke="${BRUN}" stroke-width="1.6" stroke-linecap="round"/>`;

  // En varm vindpust som släcker lågan.
  BILDER.morker = () => `<rect x="4" y="4" width="92" height="92" rx="6" fill="#1a130c" opacity="0.85"/>
<path d="M44 34h12M50 28v6" stroke="${JARN}" stroke-width="2"/>
<path d="M40 38h20l-2 34H42z" fill="#2a2018" stroke="${JARN}" stroke-width="1.6"/>
<path d="M40 72h20" stroke="${JARN}" stroke-width="3"/>
<path d="M50 62q-3 -6 2 -12q-1 6 3 8" fill="none" stroke="#8a7f6e" stroke-width="1.4" opacity="0.8"/>
<path d="M50 50q-6 -8 0 -14q4 -6 -2 -12" fill="none" stroke="#8a7f6e" stroke-width="1" opacity="0.5"/>
<path d="M8 50q14 -4 26 0M10 58q12 -3 22 0M12 42q10 -3 18 0" fill="none" stroke="#c89a5a" stroke-width="1.4" stroke-linecap="round" opacity="0.7"/>`;

  // Ett eko som inte vill tystna: ringar som går genom stenen.
  BILDER.eko = () => `${vagg()}
${[8, 16, 25, 35, 46].map((r, i) => `<circle cx="50" cy="54" r="${r}" fill="none" stroke="${BRUN}" stroke-width="${2.2 - i * 0.35}" opacity="${1 - i * 0.17}"/>`).join("")}
<circle cx="50" cy="54" r="3" fill="${BLAECK}"/>`;

  // ---- Varelserna och striden ----

  // Något blekt och mångbent som kryper ur en springa.
  BILDER.krypare = () => {
    const led = [0, 1, 2, 3, 4, 5].map((i) => [26 + i * 10, 56 - Math.sin(i / 1.7) * 8]);
    return `${vagg()}
<path d="M8 22q24 8 24 34q0 20 -6 38" fill="none" stroke="#140e09" stroke-width="4"/>
${led.map(([x, y]) => `<path d="M${f(x)} ${f(y)}l-5 -14 -4 -5M${f(x)} ${f(y)}l-4 16 -5 5" fill="none" stroke="#8a8270" stroke-width="1.6" stroke-linecap="round"/>`).join("")}
${led.map(([x, y]) => `<ellipse cx="${f(x)}" cy="${f(y)}" rx="7.5" ry="6" fill="#e6dfcf" stroke="#8a8270" stroke-width="1"/>`).join("")}
<ellipse cx="88" cy="52" rx="7" ry="6" fill="#efe9dc" stroke="#8a8270" stroke-width="1"/>
<path d="M93 48l4 -4M93 55l4 3" stroke="#8a8270" stroke-width="1.2" stroke-linecap="round"/>`;
  };

  // Något långt och väsande, som reser sig med öppen käft.
  BILDER.yngel = () => `${golv(84)}
<path d="M14 86q20 -6 30 -2t24 -8q10 -10 0 -22t-6 -24" fill="none" stroke="#1d2410" stroke-width="9" stroke-linecap="round"/>
<path d="M14 86q20 -6 30 -2t24 -8q10 -10 0 -22t-6 -24" fill="none" stroke="${GRON}" stroke-width="7" stroke-linecap="round"/>
<path d="M62 30l14 -8 4 6 -6 4 6 4 -4 4 -14 -6z" fill="${GRON}" stroke="#1d2410" stroke-width="1"/>
<circle cx="68" cy="27" r="1.6" fill="#e8c040"/>
<path d="M80 30l8 -2M88 28l3 -3M88 28l3 2" stroke="${ROD}" stroke-width="1.1" stroke-linecap="round"/>
<path d="M84 14q4 -2 6 2M86 8q6 -2 8 4" fill="none" stroke="${BRUN}" stroke-width="1.2" stroke-linecap="round"/>`;

  // En gestalt i fjällrustning, med spjut.
  BILDER.fjallvakt = () => `${vagg()}
<path d="M72 8V92" stroke="${BRUN}" stroke-width="2.4"/><path d="M72 4l-4 10h8z" fill="${JARN}"/>
<path d="M38 90l4 -30h16l4 30z" fill="#3c4a32"/>
<path d="M34 64q-2 -22 16 -26q18 4 16 26z" fill="#4f5e3c" stroke="${BLAECK}" stroke-width="1"/>
${[0, 1, 2, 3].map((r) => [0, 1, 2, 3].map((c) => `<path d="M${38 + c * 7 + (r % 2) * 3.5} ${46 + r * 5}q3.5 4 7 0" fill="none" stroke="#a8b47a" stroke-width="0.9"/>`).join("")).join("")}
<path d="M40 30q10 -16 20 0v6H40z" fill="#4f5e3c" stroke="${BLAECK}" stroke-width="1"/>
<rect x="43" y="31" width="14" height="3" fill="#140e09"/>
<path d="M64 50l8 -6" stroke="#4f5e3c" stroke-width="5" stroke-linecap="round"/>`;

  // En grönaktig skepnad som svävar.
  BILDER.gallgast = () => `<rect x="4" y="4" width="92" height="92" rx="6" fill="#2a2a1c" opacity="0.25"/>
<path d="M30 86q-4 -40 6 -56q14 -20 28 0q10 16 6 56l-7 -6 -6 8 -6 -8 -7 8 -7 -8z" fill="#a8c890" opacity="0.7" stroke="#5a7a48" stroke-width="1.2"/>
<ellipse cx="44" cy="42" rx="4" ry="6" fill="#1e2a16"/><ellipse cx="58" cy="42" rx="4" ry="6" fill="#1e2a16"/>
<ellipse cx="51" cy="58" rx="3.5" ry="5" fill="#1e2a16"/>
<path d="M20 50q-6 6 0 12M82 46q6 6 0 12" fill="none" stroke="#7a9a64" stroke-width="1.4" opacity="0.7"/>`;

  // Ett väsen av vita ben.
  BILDER.benvaktare = () => `${vagg()}
<circle cx="50" cy="24" r="12" fill="#f2ecdd" stroke="#8a8270" stroke-width="1"/>
<ellipse cx="45" cy="23" rx="3" ry="3.6" fill="${BLAECK}"/><ellipse cx="55" cy="23" rx="3" ry="3.6" fill="${BLAECK}"/>
<path d="M50 27l-2 4h4z" fill="${BLAECK}"/><path d="M44 33h12" stroke="${BLAECK}" stroke-width="1" stroke-dasharray="1.5 1"/>
<path d="M50 36v38" stroke="#f2ecdd" stroke-width="3"/>
${[42, 49, 56, 63].map((y, i) => `<path d="M50 ${y}q-${14 - i * 2} 2 -${14 - i * 2} 8M50 ${y}q${14 - i * 2} 2 ${14 - i * 2} 8" fill="none" stroke="#f2ecdd" stroke-width="2.2" stroke-linecap="round"/>`).join("")}
<path d="M50 74l-10 18M50 74l10 18M36 42l-14 18M64 42l14 18" stroke="#f2ecdd" stroke-width="3" stroke-linecap="round"/>
<path d="M22 60l-2 4M78 60l2 4" stroke="#8a8270" stroke-width="2"/>`;

  // Två skuggor som slåss: en annan vandrare.
  const skugga = (x, vand) => `<g transform="translate(${x} 0) scale(${vand} 1)">
<circle cx="0" cy="34" r="7" fill="#2a2018"/>
<path d="M-8 90l2 -30q-4 -12 6 -18q10 6 6 18l2 30z" fill="#2a2018"/>
<path d="M4 50l14 -10" stroke="#2a2018" stroke-width="4" stroke-linecap="round"/></g>`;
  BILDER.vandrare = () => `${golv(86)}
${skugga(30, 1)}${skugga(70, -1)}
<path d="M44 40l18 -22M56 40l-18 -22" stroke="#9a958c" stroke-width="2" stroke-linecap="round"/>
<circle cx="50" cy="29" r="3" fill="#fff0b0" opacity="0.9"/>`;

  // Krysset över en varelse som har fallit.
  const KRYSS = `<path d="M18 18L82 82M82 18L18 82" stroke="${ROD}" stroke-width="5" stroke-linecap="round" opacity="0.75"/>`;

  const ORM = ["sover", "ror_sig", "vaknar"];

  return {
    bild(motiv, opt = {}) {
      const rita = BILDER[motiv] || BILDER.fackla;
      let inner = rita();
      if (opt.svag) inner = `<g opacity="0.45" style="filter: grayscale(1)">${inner}</g>`;
      if (opt.fallen) inner = `<g opacity="0.6">${inner}</g>${KRYSS}`;
      return svg(inner);
    },
    // mini ger ett litet kort. Ormkorten har egna små bilder i CSS:en;
    // de andra får sin bild i liten skala.
    mini(motiv, opt = {}) {
      const s = document.createElement("span");
      s.setAttribute("aria-hidden", "true");
      if (ORM.includes(motiv)) {
        s.className = "minikort " + motiv;
      } else {
        s.className = "minikort bild";
        s.innerHTML = this.bild(motiv, opt);
      }
      return s;
    },
    // Hjälp åt delarna längre ner i filen.
    _: { f, svg, vagg, golv, laga, BILDER, BLAECK, BRUN, STEN, MORK_STEN, JARN, GULD, GRON, ROD },
  };
})();
