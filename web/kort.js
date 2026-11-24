// Spelkortens bilder (docs/spelkort.md). Varje motiv är en liten SVG i
// samma stil som ormkorten: några få former i jordfärger på pergament,
// utan text, så att bilderna är desamma på båda språken. Servern säger
// motivet (bild i data/ eller id:t); ett motiv som saknas här ritas som
// facklan.
//
//   Kort.bild(motiv, { svag, fallen, forbannad })  SVG-text för kortets bild
//   Kort.mini(motiv, { svag, fallen, forbannad })  ett litet kort framför en rad

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

  // ---- Föremålen, dryckerna och guldet ----

  // Föremålen ligger på en fläck ljus, som på ett bord.
  const ljus = `<circle cx="50" cy="52" r="38" fill="#fff6dc" opacity="0.35"/>`;

  // Rep med kroken kvar.
  BILDER.rep = () => `${ljus}
${[0, 1, 2].map((i) => `<ellipse cx="44" cy="${62 - i * 4}" rx="${22 - i * 3}" ry="${10 - i}" fill="none" stroke="#a8823a" stroke-width="4.5"/>`).join("")}
${[0, 1, 2].map((i) => `<ellipse cx="44" cy="${62 - i * 4}" rx="${22 - i * 3}" ry="${10 - i}" fill="none" stroke="${BRUN}" stroke-width="1" stroke-dasharray="2 2"/>`).join("")}
<path d="M64 56q10 -10 12 -26" fill="none" stroke="#a8823a" stroke-width="4"/>
<path d="M76 30v-8q0 -6 -6 -6q-6 0 -6 6" fill="none" stroke="${JARN}" stroke-width="3.5" stroke-linecap="round"/>
<path d="M64 22l-3 4" stroke="${JARN}" stroke-width="3" stroke-linecap="round"/>`;

  // Fjällamuletten i en snodd.
  BILDER.amulett = () => `${ljus}
<path d="M30 14q20 30 20 30q0 0 20 -30" fill="none" stroke="${BRUN}" stroke-width="1.6"/>
<ellipse cx="50" cy="60" rx="17" ry="21" fill="#5b6b34" stroke="${BLAECK}" stroke-width="1.4"/>
${[0, 1, 2, 3].map((r) => [0, 1, 2].map((c) => `<path d="M${38 + c * 8 + (r % 2) * 4} ${48 + r * 7}q4 5 8 0" fill="none" stroke="#a8b47a" stroke-width="1"/>`).join("")).join("")}
<circle cx="50" cy="40" r="3.5" fill="none" stroke="${GULD}" stroke-width="1.6"/>`;

  // En extra veke, hoprullad, och en gnista.
  BILDER.veke = () => {
    const spiral = "M50 60m-16 0a16 16 0 1 0 32 0a12 12 0 1 0 -24 0a8 8 0 1 0 16 0a4 4 0 1 0 -8 0";
    return `${ljus}
<path d="${spiral}" fill="none" stroke="${BRUN}" stroke-width="6"/>
<path d="M34 60q-6 -16 6 -30" fill="none" stroke="${BRUN}" stroke-width="6" stroke-linecap="round"/>
<path d="${spiral}" fill="none" stroke="#e6d8b8" stroke-width="3.6"/>
<path d="M34 60q-6 -16 6 -30" fill="none" stroke="#e6d8b8" stroke-width="3.6" stroke-linecap="round"/>
<path d="M40 30l0 -3" stroke="${BLAECK}" stroke-width="2"/>
${laga(40, 28, 14)}`;
  };

  // En bucklig sköld.
  BILDER.skold = () => `${ljus}
<path d="M50 14l30 10q0 40 -30 62q-30 -22 -30 -62z" fill="#8a7a5a" stroke="${BLAECK}" stroke-width="1.6"/>
<path d="M50 20l24 8q0 32 -24 50q-24 -18 -24 -50z" fill="#a8946a"/>
<path d="M50 22v54M28 40h44" stroke="${BRUN}" stroke-width="2"/>
<ellipse cx="38" cy="54" rx="5" ry="3" fill="#6a5a3a" opacity="0.6"/><ellipse cx="62" cy="34" rx="4" ry="2.5" fill="#6a5a3a" opacity="0.6"/>
<path d="M60 58l6 4" stroke="${BLAECK}" stroke-width="1"/>`;

  // Väskan med dubbel botten.
  BILDER.vaska = () => `${ljus}
<path d="M24 40h52l-4 46H28z" fill="#7a5a34" stroke="${BLAECK}" stroke-width="1.4"/>
<path d="M24 40q26 -16 52 0" fill="#8a6a40" stroke="${BLAECK}" stroke-width="1.4"/>
<path d="M36 34q14 -22 28 0" fill="none" stroke="${BRUN}" stroke-width="3"/>
<path d="M28 74h44" stroke="${GULD}" stroke-width="1.4" stroke-dasharray="3 2"/>
<ellipse cx="44" cy="80" rx="3" ry="1.4" fill="${GULD}"/><ellipse cx="54" cy="81" rx="3" ry="1.4" fill="${GULD}"/>
<rect x="46" y="44" width="8" height="6" rx="1" fill="${GULD}"/>`;

  // En kompass utan nord: nålen pekar snett.
  BILDER.kompass = () => `${ljus}
<circle cx="50" cy="52" r="28" fill="#a8823a" stroke="${BLAECK}" stroke-width="1.6"/>
<circle cx="50" cy="52" r="23" fill="#efe6d0"/>
${[0, 1, 2, 3, 4, 5, 6, 7].map((i) => { const a = (i * Math.PI) / 4; return `<path d="M${f(50 + Math.cos(a) * 19)} ${f(52 + Math.sin(a) * 19)}L${f(50 + Math.cos(a) * 22)} ${f(52 + Math.sin(a) * 22)}" stroke="${BRUN}" stroke-width="1.2"/>`; }).join("")}
<path d="M50 52l14 -12 -10 15z" fill="${ROD}"/><path d="M50 52l-14 12 10 -15z" fill="${JARN}"/>
<circle cx="50" cy="52" r="2" fill="${BLAECK}"/>
<path d="M44 22h12" stroke="#a8823a" stroke-width="4" stroke-linecap="round"/>`;

  // Rökpulver: en liten pung och rök.
  BILDER.rokpulver = () => `${ljus}
<path d="M34 86q-10 -20 6 -32h20q16 12 6 32z" fill="#8a7a64" stroke="${BLAECK}" stroke-width="1.2"/>
<path d="M40 54q10 -6 20 0" fill="none" stroke="${BRUN}" stroke-width="2.4"/>
${[[46, 40, 9], [58, 30, 11], [42, 22, 8], [62, 14, 7]].map(([x, y, r]) => `<circle cx="${x}" cy="${y}" r="${r}" fill="#b8b0a0" opacity="0.6"/>`).join("")}`;

  // En spegelskärva.
  BILDER.spegel = () => `${ljus}
<path d="M30 80l12 -60 30 18 -8 46z" fill="#c8d0d4" stroke="${BLAECK}" stroke-width="1.4"/>
<path d="M38 70l8 -40 20 12" fill="none" stroke="#fff" stroke-width="2" opacity="0.7"/>
<path d="M46 60l12 -6M44 68l10 -4" stroke="#8a9aa4" stroke-width="1.2"/>`;

  // En kastkniv.
  BILDER.kastkniv = () => `${ljus}
<path d="M24 80l40 -44 10 -10 -6 14 -40 42z" fill="#b8b4ac" stroke="${BLAECK}" stroke-width="1.2"/>
<path d="M30 76l40 -44" stroke="#efece6" stroke-width="1" opacity="0.8"/>
<path d="M20 70l14 14" stroke="${BRUN}" stroke-width="3" stroke-linecap="round"/>
<path d="M12 90l12 -12" stroke="${BRUN}" stroke-width="5" stroke-linecap="round"/>
<circle cx="12" cy="90" r="3.5" fill="none" stroke="${JARN}" stroke-width="1.6"/>`;

  // Tornet andas: en vägg som glider isär, med en ny öppning och luft som
  // strömmar ut.
  BILDER.andas = () => `${vagg()}
<path d="M14 20h26v66H14z" fill="${MORK_STEN}" opacity="0.8"/>
<path d="M60 20h26v66H60z" fill="${MORK_STEN}" opacity="0.8"/>
<path d="M40 20h20v66H40z" fill="#1e1610" opacity="0.85"/>
<path d="M40 20q-6 33 0 66M60 20q6 33 0 66" fill="none" stroke="${BLAECK}" stroke-width="1.4"/>
${[34, 50, 66].map((y, i) => `<path d="M44 ${y}q6 -${4 + i} 12 0" fill="none" stroke="#e6d8b8" stroke-width="1.6" opacity="0.75"/>`).join("")}
<path d="M34 50l-8 0M66 50l8 0" stroke="${BRUN}" stroke-width="2" stroke-linecap="round"/>
<path d="M26 50l4 -3v6zM74 50l-4 -3v6z" fill="${BRUN}"/>`;

  // Guld ur sprickan: en vägg med en spricka, och mynt som rasar ut.
  BILDER.spricka = () => `${vagg()}
<path d="M50 10l-6 18 8 10 -10 16 6 12 -4 16" fill="none" stroke="${BLAECK}" stroke-width="2.6" stroke-linejoin="round"/>
<path d="M50 10l-6 18 8 10 -10 16 6 12 -4 16" fill="none" stroke="#1e1610" stroke-width="5" opacity="0.35"/>
${[[46, 58, 0], [54, 66, 1], [40, 72, 0], [58, 78, 1], [48, 84, 0], [36, 86, 1], [62, 88, 0]].map(([x, y, k]) => `<ellipse cx="${x}" cy="${y}" rx="${k ? 4 : 3.4}" ry="1.8" fill="${GULD}" stroke="${BRUN}" stroke-width="0.6"/>`).join("")}
<circle cx="52" cy="46" r="2" fill="${GULD}"/><circle cx="47" cy="52" r="1.6" fill="${GULD}"/>`;

  // Ljus genom gluggarna: smala fönster högt upp och strålar ner i mörkret.
  BILDER.gluggar = () => `<rect x="4" y="4" width="92" height="92" rx="6" fill="#2a2018" opacity="0.55"/>
${[24, 48, 72].map((x) => `<rect x="${x - 3}" y="12" width="6" height="14" rx="3" fill="#fff6dc"/>`).join("")}
${[24, 48, 72].map((x) => `<path d="M${x - 3} 24L${x - 14} 90H${x + 10}L${x + 3} 24z" fill="#fff6dc" opacity="0.22"/>`).join("")}
${golv(86)}`;

  // Ormens ömsade skinn: ett torrt, genomskinligt rör av fjäll i en slinga.
  BILDER.ormskinn = () => `${golv(80)}
<path d="M16 74q8 -30 30 -24t28 -18q6 -8 -2 -16" fill="none" stroke="${BRUN}" stroke-width="13" stroke-linecap="round" opacity="0.55"/>
<path d="M16 74q8 -30 30 -24t28 -18q6 -8 -2 -16" fill="none" stroke="#f3ead6" stroke-width="10" stroke-linecap="round"/>
<path d="M16 74q8 -30 30 -24t28 -18q6 -8 -2 -16" fill="none" stroke="${STEN}" stroke-width="8" stroke-dasharray="2 3" stroke-linecap="round" opacity="0.6"/>
<path d="M68 12q-6 -4 -10 2" fill="none" stroke="${BRUN}" stroke-width="1.2"/>`;

  // Ett ras: tak som rasar ner i en dörröppning.
  BILDER.ras = () => `${vagg()}
<path d="M30 30h40v60H30z" fill="#1e1610" opacity="0.85"/>
<path d="M26 26h48v6H26z" fill="${MORK_STEN}"/>
${[[34, 78, 12, 10], [48, 74, 14, 14], [62, 80, 10, 9], [40, 62, 9, 8], [56, 60, 11, 9], [46, 48, 8, 7]].map(([x, y, w, h]) => `<rect x="${x - w / 2}" y="${y - h / 2}" width="${w}" height="${h}" rx="2" fill="${STEN}" stroke="${BLAECK}" stroke-width="0.8" transform="rotate(${(x * 7) % 30 - 15} ${x} ${y})"/>`).join("")}
${[[22, 40], [78, 44], [30, 20], [70, 18]].map(([x, y]) => `<circle cx="${x}" cy="${y}" r="1.6" fill="${STEN}"/>`).join("")}`;

  // Den döde vandrar: ett kryss på golvet, och ben som reser sig ur det.
  BILDER.doden = () => `${golv(80)}
<path d="M38 74l24 10M62 74l-24 10" stroke="${ROD}" stroke-width="2.2" stroke-linecap="round" opacity="0.7"/>
<circle cx="50" cy="30" r="9" fill="#efe9da" stroke="${BLAECK}" stroke-width="1.2"/>
<circle cx="46.5" cy="29" r="2" fill="${BLAECK}"/><circle cx="53.5" cy="29" r="2" fill="${BLAECK}"/>
<path d="M50 39v28M40 48h20M42 56h16M44 64h12" stroke="#efe9da" stroke-width="3" stroke-linecap="round"/>
<path d="M50 39v28" stroke="${BLAECK}" stroke-width="0.8"/>`;

  // Klockslaget: en klocka i mörkret och ljudvågor, och en röd sol.
  BILDER.klockslag = () => `<rect x="4" y="4" width="92" height="92" rx="6" fill="#2a2018" opacity="0.45"/>
<circle cx="78" cy="22" r="9" fill="${ROD}" opacity="0.8"/>
<path d="M50 18v6" stroke="${JARN}" stroke-width="2"/>
<path d="M36 64q0 -38 14 -40q14 2 14 40z" fill="#a8823a" stroke="${BLAECK}" stroke-width="1.4"/>
<path d="M32 64h36" stroke="${BLAECK}" stroke-width="3" stroke-linecap="round"/>
<circle cx="50" cy="70" r="3.2" fill="${JARN}"/>
<path d="M22 44q-6 12 0 24M16 40q-9 16 0 32M78 44q6 12 0 24M84 40q9 16 0 32" fill="none" stroke="#e6d8b8" stroke-width="1.6" stroke-linecap="round" opacity="0.8"/>`;

  // Klockan stannar: klockan från klockslaget, stilla, med kläppen rakt
  // ner och en spindeltråd över, och en blek sol som står kvar.
  BILDER.klocka_stilla = () => `<rect x="4" y="4" width="92" height="92" rx="6" fill="#2a2018" opacity="0.35"/>
<circle cx="78" cy="22" r="9" fill="${GULD}" opacity="0.6"/>
<path d="M50 18v6" stroke="${JARN}" stroke-width="2"/>
<path d="M36 64q0 -38 14 -40q14 2 14 40z" fill="#a8823a" stroke="${BLAECK}" stroke-width="1.4"/>
<path d="M32 64h36" stroke="${BLAECK}" stroke-width="3" stroke-linecap="round"/>
<path d="M50 60v8" stroke="${JARN}" stroke-width="1.4"/><circle cx="50" cy="70" r="3.2" fill="${JARN}"/>
<path d="M20 14L50 40L84 70M50 40L20 60M50 40L72 12" stroke="#e6d8b8" stroke-width="0.6" opacity="0.7"/>
<path d="M32 26l12 14M60 22l-6 14M30 52l14 -8" stroke="#e6d8b8" stroke-width="0.5" opacity="0.5"/>`;

  // Draget: en fackla vars låga blåses åt sidan och blir till rök, och
  // drag som går genom bilden.
  BILDER.draget = () => `<rect x="4" y="4" width="92" height="92" rx="6" fill="#1a130c" opacity="0.6"/>
<path d="M44 46l6 0 -2 40h-2z" fill="${BRUN}" stroke="${BLAECK}" stroke-width="1"/>
<path d="M42 44h10l-1 5h-8z" fill="${JARN}"/>
<path d="M48 42q8 -2 14 -8q-4 6 -12 9z" fill="#e8a83a" opacity="0.7"/>
<path d="M50 40q10 -10 22 -10t18 -10" fill="none" stroke="${STEN}" stroke-width="3" stroke-linecap="round" opacity="0.6"/>
<path d="M50 34q12 -8 20 -4" fill="none" stroke="${STEN}" stroke-width="2" stroke-linecap="round" opacity="0.4"/>
${[22, 56, 72].map((y, i) => `<path d="M${8 + i * 4} ${y}q14 -5 26 0t26 0" fill="none" stroke="#e6d8b8" stroke-width="1.4" stroke-linecap="round" opacity="0.55"/>`).join("")}`;

  // Det klirrar: en tung pung som har spruckit, och mynt som studsar ut.
  BILDER.klirrar = () => `${golv(82)}
<path d="M30 80q-6 -26 12 -34q-4 -6 4 -10h10q8 4 4 10q18 8 12 34z" fill="${BRUN}" stroke="${BLAECK}" stroke-width="1.4"/>
<path d="M44 36q6 4 12 0" fill="none" stroke="${BLAECK}" stroke-width="2"/>
<path d="M58 60l6 4 -4 5" fill="none" stroke="${BLAECK}" stroke-width="1.4"/>
${[[70, 56, 0], [78, 44, 1], [74, 74, 0], [84, 66, 1], [64, 40, 0]].map(([x, y, k]) => `<ellipse cx="${x}" cy="${y}" rx="${k ? 4 : 3.4}" ry="${k ? 2 : 3.4}" fill="${GULD}" stroke="${BRUN}" stroke-width="0.6"/>`).join("")}
<path d="M80 30l4 -4M88 38l5 -2M86 52l6 0" stroke="${GULD}" stroke-width="1.4" stroke-linecap="round"/>`;

  // Den glömda packningen: en sliten läderväska som ligger öppen, med
  // något som glimmar i, och damm från taket.
  BILDER.packning = () => `${golv(80)}
<path d="M24 78q-2 -22 6 -30h40q8 8 6 30z" fill="${BRUN}" stroke="${BLAECK}" stroke-width="1.4"/>
<path d="M30 48q20 -18 40 0" fill="none" stroke="${BLAECK}" stroke-width="2.4"/>
<path d="M28 50l-6 -10h56l-6 10z" fill="#8a6234" stroke="${BLAECK}" stroke-width="1.2"/>
<path d="M36 62h28M40 70h20" stroke="${BLAECK}" stroke-width="1" opacity="0.5"/>
<circle cx="50" cy="44" r="4" fill="${GULD}"/><path d="M50 34v4M43 37l3 3M57 37l-3 3" stroke="${GULD}" stroke-width="1.2" stroke-linecap="round"/>
${[[30, 14], [46, 20], [62, 12], [70, 24], [38, 26]].map(([x, y]) => `<circle cx="${x}" cy="${y}" r="1.4" fill="${STEN}"/>`).join("")}`;

  // Brevet under dörren: en dörr och ett vikt brev som sticker in under
  // dörrkanten.
  BILDER.brev_dorr = () => `${vagg()}
<path d="M28 14h44v66H28z" fill="#6b4a24" stroke="${BLAECK}" stroke-width="1.6"/>
<path d="M38 14v66M50 14v66M62 14v66" stroke="${BLAECK}" stroke-width="0.8" opacity="0.5"/>
<circle cx="64" cy="48" r="2.4" fill="${JARN}"/>
<path d="M14 80h72" stroke="${BLAECK}" stroke-width="2"/>
<path d="M36 80l4 12h34l-6 -12z" fill="#f3ead6" stroke="${BLAECK}" stroke-width="1"/>
<path d="M44 86h20M46 89h14" stroke="${BRUN}" stroke-width="0.8"/>`;

  // Tornet sluter sig: två stenblock som glider ihop över en dörröppning,
  // med pilar mot mitten.
  BILDER.sluter = () => `${vagg()}
<path d="M40 20h20v66H40z" fill="#1e1610" opacity="0.85"/>
<path d="M18 20h28v66H18z" fill="${MORK_STEN}" stroke="${BLAECK}" stroke-width="1.2"/>
<path d="M54 20h28v66H54z" fill="${MORK_STEN}" stroke="${BLAECK}" stroke-width="1.2"/>
<path d="M18 42h28M18 64h28M54 42h28M54 64h28M32 20v22M68 42v22M32 64v22" stroke="${BLAECK}" stroke-width="0.8" opacity="0.6"/>
<path d="M8 53h8M92 53h-8" stroke="${BRUN}" stroke-width="2" stroke-linecap="round"/>
<path d="M16 53l-4 -3v6zM84 53l4 -3v6z" fill="${BRUN}"/>
${[30, 46, 62, 78].map((y) => `<circle cx="50" cy="${y}" r="1.2" fill="${STEN}"/>`).join("")}`;

  // Yngel ur väggen: en spricka i muren och små slingrande ormar som
  // kryper ut ur den.
  BILDER.yngel_vagg = () => `${vagg()}
<path d="M50 8l-4 14 6 10 -8 14 4 10" fill="none" stroke="${BLAECK}" stroke-width="5" opacity="0.6" stroke-linejoin="round"/>
<path d="M50 8l-4 14 6 10 -8 14 4 10" fill="none" stroke="#1e1610" stroke-width="2.4" stroke-linejoin="round"/>
${[[48, 40, 30, 70, 0], [46, 52, 66, 80, 1], [50, 30, 74, 52, 0]].map(([x0, y0, x1, y1, k]) => `<path d="M${x0} ${y0}q${(x1 - x0) / 2} ${k ? -8 : 8} ${x1 - x0} ${y1 - y0}" fill="none" stroke="#1d2410" stroke-width="5" stroke-linecap="round"/><path d="M${x0} ${y0}q${(x1 - x0) / 2} ${k ? -8 : 8} ${x1 - x0} ${y1 - y0}" fill="none" stroke="${GRON}" stroke-width="3.4" stroke-linecap="round"/><circle cx="${x1}" cy="${y1}" r="2.6" fill="${GRON}" stroke="#1d2410" stroke-width="0.8"/>`).join("")}
${golv(88)}`;

  // Ett rostigt svärd, snett över kortet.
  BILDER.svard = () => `${ljus}
<path d="M30 76l42 -50 6 -6 -2 8 -42 50z" fill="#9a8a7a" stroke="${BLAECK}" stroke-width="1.2"/>
<path d="M36 70l38 -44" stroke="#c8c0b4" stroke-width="1" opacity="0.8"/>
<circle cx="48" cy="56" r="2.4" fill="#8a4a24" opacity="0.8"/><circle cx="60" cy="42" r="1.8" fill="#8a4a24" opacity="0.7"/>
<path d="M22 66l18 16" stroke="${JARN}" stroke-width="4" stroke-linecap="round"/>
<path d="M30 76l-10 10" stroke="${BRUN}" stroke-width="5" stroke-linecap="round"/>
<circle cx="18" cy="88" r="3.4" fill="${JARN}"/>`;

  // En brynja av fjäll.
  BILDER.brynja = () => `${ljus}
<path d="M28 24l14 -6h16l14 6 6 22 -8 4v38H30V50l-8 -4z" fill="#5a6a54" stroke="${BLAECK}" stroke-width="1.4"/>
${[30, 40, 50, 60, 70, 80].map((y) => [36, 44, 52, 60, 68].filter((x) => y > 44 || (x > 40 && x < 64)).map((x) => `<path d="M${x - 4} ${y}q4 6 8 0" fill="none" stroke="#8a9a80" stroke-width="1.2"/>`).join("")).join("")}
<path d="M42 18q8 8 16 0" fill="none" stroke="${BLAECK}" stroke-width="1.4"/>`;

  // Svärdet med rött fäste: tunt och blankt, med en droppe vid spetsen.
  BILDER.tarst_svard = () => `${ljus}
<path d="M30 76l44 -52 4 -4 -1 6 -44 52z" fill="#d8d4cc" stroke="${BLAECK}" stroke-width="1.1"/>
<path d="M34 72l40 -46" stroke="#fff" stroke-width="0.8" opacity="0.8"/>
<path d="M22 66l18 16" stroke="${GULD}" stroke-width="3.6" stroke-linecap="round"/>
<path d="M30 76l-10 10" stroke="${ROD}" stroke-width="5" stroke-linecap="round"/>
<circle cx="18" cy="88" r="3.4" fill="${ROD}"/>
<path d="M80 26q-2 6 0 8q2 -2 0 -8z" fill="${ROD}"/>`;

  // Benluren: ett krökt horn av ben, och ljudet som kommer genom väggen.
  BILDER.benlur = () => `${vagg()}
<path d="M20 74q6 -34 40 -44l6 8q-28 10 -34 40z" fill="#e8e0cc" stroke="${BLAECK}" stroke-width="1.4"/>
<path d="M60 30l6 8 10 -8 -6 -10z" fill="#d8ccb0" stroke="${BLAECK}" stroke-width="1.2"/>
<path d="M28 66q6 -22 26 -30M34 70q4 -14 16 -22" fill="none" stroke="#b8ac90" stroke-width="1"/>
<path d="M80 18q8 6 0 14M86 12q12 12 0 26" fill="none" stroke="#f6f2e8" stroke-width="2" stroke-linecap="round" opacity="0.8"/>`;

  // En sövande flöjt av ben.
  BILDER.flojt = () => `${ljus}
<path d="M18 74l60 -40" stroke="${BRUN}" stroke-width="10" stroke-linecap="round"/>
<path d="M18 74l60 -40" stroke="#e6d8b8" stroke-width="7" stroke-linecap="round"/>
${[0, 1, 2, 3].map((i) => `<circle cx="${32 + i * 10}" cy="${f(65 - i * 6.7)}" r="1.8" fill="${BLAECK}"/>`).join("")}
<text x="64" y="70" font-size="12" fill="${BRUN}" font-family="serif">z</text><text x="74" y="62" font-size="9" fill="${BRUN}" font-family="serif">z</text>`;

  // Ett timglas av ben.
  BILDER.timglas = () => `${ljus}
<path d="M32 18h36M32 86h36" stroke="#efe6d0" stroke-width="5" stroke-linecap="round"/>
<path d="M34 20v4M66 20v4M34 84v-4M66 84v-4" stroke="#efe6d0" stroke-width="3"/>
<path d="M38 22q0 20 12 30q-12 10 -12 30h24q0 -20 -12 -30q12 -10 12 -30z" fill="#e6eef0" opacity="0.6" stroke="${STEN}" stroke-width="1.2"/>
<path d="M42 30q8 10 16 0z" fill="#c9a24f"/><path d="M40 82q10 -12 20 0z" fill="#c9a24f"/>
<path d="M50 52v24" stroke="#c9a24f" stroke-width="1"/>`;

  // Den svarta pärlan i ett skal.
  BILDER.svart_parla = () => `${ljus}
<path d="M18 66q32 22 64 0q-32 -6 -64 0z" fill="#d8cfc0" stroke="${BLAECK}" stroke-width="1.2"/>
<path d="M22 64q28 -40 56 0" fill="#e6ddcc" stroke="${BLAECK}" stroke-width="1.2"/>
${[0, 1, 2, 3, 4].map((i) => `<path d="M50 30L${26 + i * 12} 64" stroke="${STEN}" stroke-width="0.8"/>`).join("")}
<circle cx="50" cy="62" r="10" fill="#1a1820"/><circle cx="46" cy="58" r="3" fill="#6a6a7a" opacity="0.8"/>`;

  // Myntet med hål.
  BILDER.mynt_hal = () => `${ljus}
<circle cx="50" cy="52" r="26" fill="${GULD}" stroke="#7a5a1a" stroke-width="1.6"/>
<circle cx="50" cy="52" r="20" fill="none" stroke="#7a5a1a" stroke-width="1" stroke-dasharray="2 2"/>
<circle cx="50" cy="52" r="7" fill="#efe6d0" stroke="#7a5a1a" stroke-width="1.4"/>
<path d="M50 26q-4 -12 4 -16" fill="none" stroke="${BRUN}" stroke-width="1.4"/>`;

  // Kronan av rost.
  BILDER.rostkrona = () => `${ljus}
<path d="M22 72V36l14 16 14 -22 14 22 14 -16v36z" fill="#8a4a24" stroke="${BLAECK}" stroke-width="1.4"/>
<rect x="22" y="70" width="56" height="10" fill="#7a3a1a" stroke="${BLAECK}" stroke-width="1.4"/>
${[[30, 62], [46, 48], [62, 58], [70, 44], [38, 76]].map(([x, y]) => `<circle cx="${x}" cy="${y}" r="2.6" fill="#b8642a" opacity="0.8"/>`).join("")}
${[36, 50, 64].map((x) => `<circle cx="${x}" cy="75" r="2.6" fill="#5a2a5a"/>`).join("")}`;

  // Ett brev från en tidigare vandrare, med ett sigill.
  BILDER.brev = () => `${ljus}
<path d="M20 32h60v42H20z" fill="#efe6d0" stroke="${BRUN}" stroke-width="1.4"/>
<path d="M20 32l30 24 30 -24" fill="none" stroke="${BRUN}" stroke-width="1.4"/>
<path d="M20 74l22 -20M80 74l-22 -20" stroke="${STEN}" stroke-width="1"/>
<circle cx="50" cy="56" r="6" fill="${ROD}"/><circle cx="50" cy="56" r="3" fill="none" stroke="#c0503a" stroke-width="1"/>`;

  // En liten flaska med en röd dryck.
  BILDER.dryck = () => `${ljus}
<path d="M44 16h12v6h-12z" fill="#8a6a40"/>
<path d="M46 22h8v14q16 8 16 26a20 20 0 0 1 -40 0q0 -18 16 -26z" fill="#e6eef0" opacity="0.7" stroke="${BLAECK}" stroke-width="1.4"/>
<path d="M32 62a18 18 0 0 0 36 0q-18 -6 -36 0z" fill="#b8282a"/>
<path d="M38 52q2 -8 8 -10" fill="none" stroke="#fff" stroke-width="2" opacity="0.7"/>`;

  // Några mynt.
  BILDER.guld = () => `${ljus}
${[[30, 72], [44, 76], [58, 74], [72, 72], [38, 66], [52, 68], [66, 66], [46, 60], [58, 60], [52, 54]].map(([x, y]) => `<ellipse cx="${x}" cy="${y}" rx="9" ry="4" fill="${GULD}" stroke="#7a5a1a" stroke-width="1"/>`).join("")}
<path d="M64 36l3 -6 3 6 -3 6z" fill="#fff0b0"/>`;

  // Ett förbannat föremål: mörkt, med sprickor.
  const FORBANNAD = `<rect x="4" y="4" width="92" height="92" rx="6" fill="#4a0a06" opacity="0.35"/>
<path d="M12 20l14 10 -4 12 12 8M88 80l-12 -8 4 -14 -10 -6M70 12l-6 12 8 6" fill="none" stroke="#2a0604" stroke-width="1.4" opacity="0.8"/>`;

  // ---- Källarvalven och avslöjandet ----

  // Valv som går in i mörkret under tornet.
  BILDER.kallare = () => `<rect x="4" y="4" width="92" height="92" rx="6" fill="#1e1610" opacity="0.8"/>
${[0, 1, 2, 3].map((i) => { const w = 76 - i * 16, x = 50 - w / 2, y = 18 + i * 9, h = 80 - i * 14; return `<path d="M${x} ${y + h}V${y + w / 2}a${w / 2} ${w / 2} 0 0 1 ${w} 0V${y + h}" fill="none" stroke="${STEN}" stroke-width="${4 - i * 0.8}" opacity="${0.9 - i * 0.2}"/>`; }).join("")}
<path d="M14 94l14 -10h44l14 10" fill="none" stroke="${MORK_STEN}" stroke-width="2"/>`;

  // Ett varsel: ett öga i dimman.
  BILDER.varsel = () => `${[[26, 64, 18], [52, 70, 22], [76, 62, 16], [40, 40, 14], [66, 38, 12]].map(([x, y, r]) => `<circle cx="${x}" cy="${y}" r="${r}" fill="#b8b0a0" opacity="0.35"/>`).join("")}
<path d="M24 50q26 -26 52 0q-26 26 -52 0z" fill="#efe6d0" stroke="${BRUN}" stroke-width="1.6"/>
<circle cx="50" cy="50" r="10" fill="#6a7a4a"/><circle cx="50" cy="50" r="4.5" fill="${BLAECK}"/>
<circle cx="47" cy="47" r="1.8" fill="#fff"/>`;

  // En del av namnet: tre tecken ristade i sten, och ett som syns.
  BILDER.namndel = () => `${vagg()}
<rect x="16" y="30" width="68" height="40" rx="4" fill="${STEN}" opacity="0.5" stroke="${MORK_STEN}" stroke-width="1.4"/>
<path d="M26 38v24M26 44l8 -6M26 52l8 6" fill="none" stroke="${GULD}" stroke-width="2.6" stroke-linecap="round"/>
<path d="M46 38l8 24M54 38l-8 24M44 50h12" fill="none" stroke="${MORK_STEN}" stroke-width="2" stroke-linecap="round" opacity="0.6"/>
<path d="M68 38q-6 12 0 24M72 38q6 12 0 24" fill="none" stroke="${MORK_STEN}" stroke-width="2" stroke-linecap="round" opacity="0.6"/>`;

  // Tornet lyssnar: en viskning som går in i en slinga.
  BILDER.invigd = () => {
    const pts = serpentSpiral(62, 52, 2.2, 2, 3.4);
    const d = "M" + pts.map(([x, y]) => f(x) + " " + f(y)).join("L");
    return `${vagg()}
<path d="${d}" fill="none" stroke="#5a3a6a" stroke-width="3" stroke-linecap="round" opacity="0.8"/>
<path d="M14 40q6 -6 12 0q-6 8 -12 0z" fill="${BRUN}"/>
<path d="M28 38q6 2 10 0M28 44q8 2 12 0M28 50q6 2 10 0" fill="none" stroke="${BRUN}" stroke-width="1.4" stroke-linecap="round"/>`;
  };

  // ---- Tornets stora kort och slutet ----

  // Ute: porten med kvällsljuset utanför och någon som går ut.
  BILDER.ut = () => `${vagg()}
<path d="M24 94V44a26 26 0 0 1 52 0v50z" fill="#f0c070"/>
<circle cx="50" cy="66" r="14" fill="#f8e0a0"/>
<path d="M24 78h52v16H24z" fill="#a8823a" opacity="0.6"/>
<path d="M24 94V44a26 26 0 0 1 52 0v50" fill="none" stroke="${MORK_STEN}" stroke-width="3"/>
<circle cx="50" cy="62" r="4" fill="#2a2018"/><path d="M45 88l2 -20h6l2 20z" fill="#2a2018"/>`;

  // Döden: en släckt fackla med rök.
  BILDER.dod = () => `<rect x="4" y="4" width="92" height="92" rx="6" fill="#1a130c" opacity="0.85"/>
<path d="M44 40h12l-3 8h-6z" fill="${JARN}"/>
<rect x="47" y="46" width="6" height="40" rx="2" fill="#4a3420"/>
<path d="M44 40q6 -4 12 0" fill="none" stroke="#2a2018" stroke-width="2"/>
<path d="M50 36q-6 -8 0 -14t0 -14" fill="none" stroke="#8a7f6e" stroke-width="1.6" opacity="0.7"/>
<path d="M54 34q4 -6 0 -10" fill="none" stroke="#8a7f6e" stroke-width="1.2" opacity="0.5"/>`;

  // Solen som går ner bakom horisonten.
  BILDER.sol = () => `<rect x="4" y="4" width="92" height="92" rx="6" fill="#e8a868" opacity="0.45"/>
<rect x="4" y="40" width="92" height="22" fill="#d87848" opacity="0.35"/>
<path d="M22 64a28 28 0 0 1 56 0z" fill="#e8783a"/>
${[-60, -30, 0, 30, 60].map((a) => { const r = (a - 90) * Math.PI / 180; return `<path d="M${f(50 + Math.cos(r) * 34)} ${f(64 + Math.sin(r) * 34)}L${f(50 + Math.cos(r) * 42)} ${f(64 + Math.sin(r) * 42)}" stroke="#c85a2a" stroke-width="2.4" stroke-linecap="round"/>`; }).join("")}
<path d="M4 64h92v28a4 4 0 0 1 -4 4H8a4 4 0 0 1 -4 -4z" fill="#3a2a1c"/>
<path d="M66 64V46h4v-4h3v4h3v-4h3v4h4v18z" fill="#2a1e14"/>`;

  // En port som stängs: två dörrar som går igen.
  BILDER.port = () => `${vagg()}
<path d="M24 94V44a26 26 0 0 1 52 0v50z" fill="#140e09"/>
<path d="M24 94V44a26 26 0 0 1 22 -25.7V94z" fill="#6b4a24" stroke="${BLAECK}" stroke-width="1.2"/>
<path d="M76 94V44a26 26 0 0 0 -22 -25.7V94z" fill="#6b4a24" stroke="${BLAECK}" stroke-width="1.2"/>
<path d="M30 40v50M38 30v60M62 30v60M70 40v50" stroke="${BLAECK}" stroke-width="0.8" opacity="0.6"/>
<path d="M26 56h20M54 56h20M26 78h20M54 78h20" stroke="${JARN}" stroke-width="2.4"/>
<path d="M10 66h8l-3 -3M10 66l5 3M90 66h-8l3 -3M90 66l-5 3" fill="none" stroke="${BRUN}" stroke-width="1.6" stroke-linecap="round"/>`;

  // ---- Klasserna ----

  // Legoknekten: en järnhatt över en sköld och ett svärd.
  BILDER.legoknekt = () => `${ljus}
<path d="M28 82l44 -50" stroke="#b8b4ac" stroke-width="4" stroke-linecap="round"/>
<path d="M30 74l10 10" stroke="${BRUN}" stroke-width="4" stroke-linecap="round"/>
<path d="M50 40l22 7q0 26 -22 40q-22 -14 -22 -40z" fill="#8a7a5a" stroke="${BLAECK}" stroke-width="1.4"/>
<path d="M50 46v34M34 56h32" stroke="${BRUN}" stroke-width="2"/>
<path d="M28 30q22 -22 44 0z" fill="${JARN}" stroke="${BLAECK}" stroke-width="1.2"/>
<path d="M20 31h60" stroke="${JARN}" stroke-width="4" stroke-linecap="round"/>`;

  // Tjuven: en hand som sträcker sig efter en pung.
  BILDER.tjuv = () => `${ljus}
<path d="M58 84q-12 0 -12 -14q0 -10 8 -14h12q8 4 8 14q0 14 -16 14z" fill="#7a5a34" stroke="${BLAECK}" stroke-width="1.2"/>
<path d="M54 56q6 -6 12 0" fill="none" stroke="${BRUN}" stroke-width="2.4"/>
<ellipse cx="60" cy="88" rx="5" ry="2" fill="${GULD}"/><ellipse cx="70" cy="86" rx="4" ry="1.8" fill="${GULD}"/>
<path d="M10 30q14 2 24 10l10 8q4 4 0 7q-4 2 -8 -2l-4 -3 6 8q2 5 -3 5l-8 -8q-12 -6 -17 -14z" fill="#4a3a2c" stroke="${BLAECK}" stroke-width="1"/>
<path d="M36 46l8 -2M38 50l7 0" stroke="${BLAECK}" stroke-width="1" opacity="0.6"/>`;

  // Örtkvinnan: en mortel med örter.
  BILDER.ortkvinna = () => `${ljus}
<path d="M28 58h44q-2 22 -22 26q-20 -4 -22 -26z" fill="#9a8a6e" stroke="${BLAECK}" stroke-width="1.4"/>
<path d="M26 58h48" stroke="${BLAECK}" stroke-width="2"/>
<path d="M56 62l18 -34" stroke="#8a7a64" stroke-width="5" stroke-linecap="round"/>
${[[34, 54, -30], [42, 50, -10], [48, 52, 20], [38, 46, -40]].map(([x, y, a]) => `<path d="M${x} ${y}q-2 -14 4 -22" fill="none" stroke="${GRON}" stroke-width="1.6" transform="rotate(${a} ${x} ${y})"/><ellipse cx="${x}" cy="${y - 16}" rx="3" ry="5" fill="#7a9a4a" transform="rotate(${a} ${x} ${y})"/>`).join("")}`;

  // Kartritaren: en karta och en passare.
  BILDER.kartritare = () => `${ljus}
<path d="M16 30l20 -6 22 6 26 -6v52l-26 6 -22 -6 -20 6z" fill="#efe6d0" stroke="${BRUN}" stroke-width="1.2"/>
<path d="M36 24v52M58 30v52" stroke="${STEN}" stroke-width="0.8"/>
<path d="M22 50h10v-12h14v20h14v-10h16" fill="none" stroke="${BRUN}" stroke-width="1.6"/>
<path d="M24 66l8 -4M68 40l6 6M74 40l-6 6" stroke="${ROD}" stroke-width="1.4"/>
<path d="M60 88l10 -30 10 30" fill="none" stroke="${JARN}" stroke-width="2"/><circle cx="70" cy="58" r="2.4" fill="${JARN}"/>`;

  // Ormtjusaren: en orm som reser sig ur en korg vid en flöjt.
  BILDER.ormtjusare = () => `${ljus}
<path d="M26 64h36l-4 24H30z" fill="#a8823a" stroke="${BLAECK}" stroke-width="1.2"/>
<path d="M28 70h32M29 76h30M30 82h28" stroke="${BRUN}" stroke-width="0.9"/>
<path d="M44 64q-8 -12 2 -20t-2 -20" fill="none" stroke="#1d2410" stroke-width="7" stroke-linecap="round"/>
<path d="M44 64q-8 -12 2 -20t-2 -20" fill="none" stroke="${GRON}" stroke-width="5" stroke-linecap="round"/>
<ellipse cx="45" cy="22" rx="6" ry="4" fill="${GRON}" stroke="#1d2410" stroke-width="1"/><circle cx="47" cy="21" r="1.2" fill="#e8c040"/>
<path d="M66 54l18 -22" stroke="${BRUN}" stroke-width="6" stroke-linecap="round"/><path d="M66 54l18 -22" stroke="#e6d8b8" stroke-width="3.6" stroke-linecap="round"/>
<path d="M76 30q4 -4 8 0M80 24q4 -4 8 0" fill="none" stroke="${BRUN}" stroke-width="1.2" stroke-linecap="round"/>`;

  // Spårfinnaren: fotspår i dammet.
  BILDER.sparfinnare = () => `${golv(62)}
${[[22, 84, -20], [34, 70, 10], [44, 60, -15], [58, 50, 12], [66, 38, -10], [80, 30, 15]].map(([x, y, a]) => `<g transform="rotate(${a} ${x} ${y})"><ellipse cx="${x}" cy="${y}" rx="3.6" ry="6" fill="${BRUN}" opacity="0.8"/><ellipse cx="${x}" cy="${y - 8.5}" rx="2.8" ry="2.4" fill="${BRUN}" opacity="0.8"/></g>`).join("")}`;

  // Lyktbäraren: en lykta som lyser.
  BILDER.lyktbarare = () => `<circle cx="50" cy="54" r="40" fill="#f8d878" opacity="0.3"/><circle cx="50" cy="54" r="26" fill="#f8d878" opacity="0.35"/>
<path d="M50 12v8M42 20q8 -8 16 0" fill="none" stroke="${JARN}" stroke-width="2.4"/>
<path d="M38 24h24l-3 8H41z" fill="${JARN}"/>
<path d="M40 32h20v36H40z" fill="#fff0b0" opacity="0.85" stroke="${JARN}" stroke-width="2.4"/>
<path d="M50 32v36M40 50h20" stroke="${JARN}" stroke-width="1.2"/>
${laga(50, 60, 16)}
<path d="M36 68h28l-4 6H40z" fill="${JARN}"/>`;

  // Pilgrimen: en stav med en mussla och en kalebass.
  BILDER.pilgrim = () => `${ljus}
<path d="M40 92L56 10" stroke="${BRUN}" stroke-width="4" stroke-linecap="round"/>
<path d="M56 10q6 2 4 8" fill="none" stroke="${BRUN}" stroke-width="3" stroke-linecap="round"/>
<path d="M52 32q-6 0 -8 6t8 14q4 -2 6 -8t-6 -12z" fill="#c8a878" stroke="${BLAECK}" stroke-width="1"/>
<circle cx="64" cy="62" r="9" fill="#b8823a" stroke="${BLAECK}" stroke-width="1.2"/><circle cx="64" cy="50" r="5" fill="#b8823a" stroke="${BLAECK}" stroke-width="1.2"/>
<path d="M52 40l-4 6M54 42l-2 7M50 38l-5 4" stroke="${BRUN}" stroke-width="0.9"/>
<path d="M56 46q4 0 6 2" fill="none" stroke="${BRUN}" stroke-width="1.4"/>`;

  const ORM = ["sover", "ror_sig", "vaknar"];

  return {
    bild(motiv, opt = {}) {
      const rita = BILDER[motiv] || BILDER.fackla;
      let inner = rita();
      if (opt.svag) inner = `<g opacity="0.45" style="filter: grayscale(1)">${inner}</g>`;
      if (opt.fallen) inner = `<g opacity="0.6">${inner}</g>${KRYSS}`;
      if (opt.forbannad) inner = `${inner}${FORBANNAD}`;
      return svg(inner);
    },
    // mini ger ett litet kort. Ormkorten har egna små bilder i CSS:en;
    // de andra får sin bild i liten skala.
    // Med en rubrik går kortet att klicka på för att se det stort (app.js
    // lyssnar efter .minikort[data-kort]).
    mini(motiv, opt = {}) {
      const s = document.createElement("span");
      if (opt.rubrik) {
        s.dataset.kort = JSON.stringify({ bild: motiv, rubrik: opt.rubrik, text: opt.text || "", svag: opt.svag, fallen: opt.fallen, forbannad: opt.forbannad });
        s.setAttribute("role", "button");
        s.setAttribute("aria-label", opt.rubrik);
        s.title = opt.rubrik;
        s.tabIndex = 0;
      } else {
        s.setAttribute("aria-hidden", "true");
      }
      if (ORM.includes(motiv)) {
        s.className = "minikort " + motiv;
        if (opt.rubrik) s.innerHTML = this.bild(motiv); // syns när kortet är stort nog (krönikan)
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
