// Kartan snett ovanifrån. Tornet byggs av golv, låga väggar och portar, och
// vandrarna är små figurer i sin färg. Kameran är ortografisk och vrids
// aldrig: norr pekar upp åt höger. Utsnittet, zoomen och teckenförklaringen
// räknas ut i web/map.js, som den här kartan läser utsnittet från.
import * as THREE from "three";
import { mergeGeometries } from "three/examples/jsm/utils/BufferGeometryUtils.js";
import * as I from "./inredning";
import * as K from "./krita";
import * as M from "./modeller";
import type { Figur } from "./modeller";
import type { Rum, Tillstand, Utsnitt, Val } from "./typer";

const S = 6; // ett rum i världen
const VH = 0.8; // murarnas höjd
const VT = 0.4; // murarnas tjocklek
const LUCKA = 1.9; // dörrarnas bredd

// Ett rum delas i celler, 7 på bredden och 4 på höjden, för vandrarnas
// platser i rummet. Cellerna kommer från den platta kartan som fanns förut.
const RW = 7;
const RH = 4;
const SLOTS = [[2, 1], [5, 1], [2, 3], [5, 3], [1, 2], [6, 2], [3, 1], [4, 3]];
// Varelsen står mellan cellerna, snett framför mitten. Kameran tittar längs
// diagonalen, så i (4, 2) stod den precis framför vandraren i (2, 1) och
// skymde henom helt. Härifrån står ingen vandrare rakt bakom den.
const VARELSE = [3.5, 2.5];

const FARG = {
  okand: "#6a6e76",
  port: "#d8cfbf",
  stangd: "#d0503a",
  guld: "#ffd36b",
  forlorad: "#a04a3c",
  mark: "#ff8a6a",
  varelse: "#ff7a5a",
  vald: "#ffe14a",
  krita: "#f2f0e8",
  borta: "#5a5d63",
  vagg: "#c9c6bc",
  jarn: "#2e2a27",
};
const KANT_VARELSE = "#7a2418"; // varelsernas socklar har en mörkröd kant
const RING = ["#c9a24f", "#c9a24f", "#c9a24f", "#c98a5a", "#c4483a"]; // per nivå

// Avslöjandet: tornets utseende följer nivån. På
// nivå 2 andas kritan, på nivå 3 syns revbenen och hjärtslagen, på nivå 4
// är det kött och tänder. vagg och hi är väggens färg i vila och i slaget,
// golv golvets ton, period hur länge ett andetag eller hjärtslag tar och
// revben hur tydliga revbenen på golvet är.
const LOOK: Record<number, { vagg: number[]; hi: number[]; golv: number[]; period: number; revben: number }> = {
  1: { vagg: [236, 228, 214], hi: [236, 228, 214], golv: [255, 255, 255], period: 0, revben: 0 },
  2: { vagg: [228, 220, 206], hi: [255, 246, 230], golv: [255, 248, 236], period: 4200, revben: 0 },
  3: { vagg: [220, 172, 152], hi: [255, 182, 152], golv: [230, 170, 160], period: 1250, revben: 0.35 },
  4: { vagg: [196, 102, 84], hi: [255, 118, 88], golv: [210, 130, 120], period: 950, revben: 0.55 },
};
const PORT_KOTT = "#a4483a";
const MS_PER_ENHET = 16; // hjärtslaget sprider sig från tornets mitt så här fort

// puls ger slaget just nu, 0–1, för en punkt dist enheter från mitten.
// Samma form som pulse i web/map.js. Väggarna räknar ut samma sak i sin
// shader (se vaggMaterial).
function puls(niva: number, t: number, dist: number): number {
  const look = LOOK[niva];
  if (!look.period) return 0;
  const tt = t - dist * MS_PER_ENHET;
  if (niva === 2) return 0.5 + 0.5 * Math.sin((2 * Math.PI * tt) / look.period);
  const ph = (((tt % look.period) + look.period) % look.period) / look.period;
  const g = (c: number) => Math.exp(-(((ph - c) / 0.07) ** 2));
  return Math.max(g(0), 0.7 * g(0.22), g(1));
}

function srgb(c: number[]): THREE.Color {
  return new THREE.Color().setRGB(c[0] / 255, c[1] / 255, c[2] / 255, THREE.SRGBColorSpace);
}
const EKO_MS = 4500;
const LANDA_MS = 650; // så länge en ny bricka faller
const BRICKA = 0.3; // brickans tjocklek
// Den som har bett om mindre rörelse får rummen utan fallande brickor.
const LUGN = typeof matchMedia === "function" && matchMedia("(prefers-reduced-motion: reduce)").matches;
const FART = 3.4; // så fort figurerna går, i världens enheter per sekund
const VAKEN_MS = 3500;

// Telefoner och svaga enheter: på en pekskärm ritas högst 30 bilder i
// sekunden, för batteriets skull. Går det för trögt ändå slår kartan av
// skuggorna och ritar med lägre upplösning (sankKvalitet). Adressen kan
// tvinga fram det med ?kvalitet=lag, för att prova.
const PEK = typeof matchMedia === "function" && matchMedia("(pointer: coarse)").matches;
const MAX_BILDER = PEK ? 30 : 60;
const TROGT_MS = PEK ? 45 : 30; // medianen för en bildruta som räknas som trög
const MATT_RUTOR = 90; // så många bildrutor mäts innan kartan bestämmer sig

// Kameran tittar från sydost, snett uppifrån.
const FRAN = new THREE.Vector3(30, 26, 30).normalize();
const LUTNING = FRAN.y; // hur mycket golvet trycks ihop på höjden
const SOL = new THREE.Vector3(-10, 22, -4);

// plats ger punkten i världen för cellen (cx, cy) i rummet (rx, ry), med
// rummets celler: 1–6 på bredden och 1–3 på höjden.
function plats(rx: number, ry: number, cx: number, cy: number): [number, number] {
  return [rx * S + S * (0.1 + (0.8 * (cx - 0.5)) / (RW - 1)), ry * S + S * (0.1 + (0.8 * (cy - 0.5)) / (RH - 1))];
}

// vidDorr ger en punkt innanför dörren åt hållet d (0–3) och vinkeln dit.
function vidDorr(rx: number, ry: number, d: number): [number, number, number] {
  const m = S / 2, in_ = 0.75;
  const p = [[m, in_], [S - in_, m], [m, S - in_], [in_, m]][d];
  return [rx * S + p[0], ry * S + p[1], [Math.PI / 2, 0, -Math.PI / 2, Math.PI][d]];
}

// En figur på kartan: en vandrare eller en varelse. Den står kvar mellan
// tillstånden och går dit den ska i stället för att hoppa.
interface Aktor {
  hall: THREE.Group; // platsen och riktningen
  fig: Figur;
  nyckel: string; // modellen: byts den, byggs figuren om
  rum: [number, number];
  vag: THREE.Vector3[]; // punkterna den går mot
  vinkel: number;
  dod: boolean;
  ring?: THREE.Mesh; // den egna vandrarens ring
}

// kompassSvg ritar ett N med en pil åt det håll norr ligger på skärmen,
// räknat från kamerans vinkel.
function kompassSvg(): string {
  const cam = new THREE.OrthographicCamera(-1, 1, 1, -1, 0.1, 10);
  cam.position.copy(FRAN);
  cam.lookAt(0, 0, 0);
  cam.updateMatrixWorld();
  const a = new THREE.Vector3(0, 0, 0).project(cam);
  const b = new THREE.Vector3(0, 0, -1).project(cam);
  const vx = b.x - a.x, vy = b.y - a.y;
  const l = Math.hypot(vx, vy);
  const ux = vx / l, uy = -vy / l; // på skärmen, y nedåt
  const c = 22;
  const tx = c + ux * 6, ty = c + uy * 6, fx = c - ux * 13, fy = c - uy * 13;
  const px = -uy, py = ux; // vinkelrätt, för pilens hullingar
  const h1x = tx - ux * 5 + px * 4, h1y = ty - uy * 5 + py * 4, h2x = tx - ux * 5 - px * 4, h2y = ty - uy * 5 - py * 4;
  const f = (n: number) => n.toFixed(1);
  return `<svg viewBox="0 0 44 44" width="44" height="44">
<circle cx="${c}" cy="${c}" r="19" fill="none" stroke="currentColor" stroke-opacity="0.25"/>
<path d="M${f(fx)} ${f(fy)}L${f(tx)} ${f(ty)}M${f(h1x)} ${f(h1y)}L${f(tx)} ${f(ty)}L${f(h2x)} ${f(h2y)}" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
<text x="${f(c + ux * 13)}" y="${f(c + uy * 13 + 4)}" font-size="11" text-anchor="middle" fill="currentColor">N</text>
</svg>`;
}

// skala gör figuren storlek gånger större (modell.storlek i datafilerna).
function skala(f: Figur, storlek?: number): Figur {
  if (storlek && storlek > 0) f.grupp.scale.setScalar(storlek);
  return f;
}

// guldhog är kammarens hög av mynt, med mitten i origo.
function guldhog(): THREE.Mesh {
  const mynt: THREE.BufferGeometry[] = [];
  let s = 3;
  const rnd = () => ((s = (s * 16807) % 2147483647) / 2147483647);
  for (let lager = 0; lager < 5; lager++) {
    const rad = 0.9 - lager * 0.17;
    for (let i = 0; i < 12 - lager * 2; i++) {
      const a = rnd() * Math.PI * 2, d = rnd() * rad;
      const g = new THREE.CylinderGeometry(0.16, 0.16, 0.05, 10);
      g.rotateX((rnd() - 0.5) * 0.6);
      g.rotateZ((rnd() - 0.5) * 0.6);
      g.translate(Math.cos(a) * d, 0.03 + lager * 0.08 + rnd() * 0.03, Math.sin(a) * d);
      mynt.push(g);
    }
  }
  const hog = new THREE.Mesh(mergeGeometries(mynt), new THREE.MeshLambertMaterial({ color: 0xe8b84a, emissive: 0x2a1c00 }));
  mynt.forEach((g) => g.dispose());
  hog.castShadow = true;
  return hog;
}

// En port i ytterväggen är ett portvalv: två torn med spetsiga tak, en
// överliggare med krön och två dörrar av trä. Öppen står dörrarna utåt;
// stängd fyller de valvet, med röda järnband. portGeometri bygger den med
// mitten i origo, längs x, och utsidan åt -z.
interface PortDelar {
  sten: THREE.BufferGeometry[];
  tra: THREE.BufferGeometry[];
  jarn: THREE.BufferGeometry[];
  rott: THREE.BufferGeometry[];
}
function portGeometri(stangd: boolean): PortDelar {
  const d: PortDelar = { sten: [], tra: [], jarn: [], rott: [] };
  const lada = (lista: THREE.BufferGeometry[], w: number, h: number, dj: number, x: number, y: number, z: number) => {
    const g = new THREE.BoxGeometry(w, h, dj);
    g.translate(x, y, z);
    lista.push(g);
  };
  const torn = LUCKA / 2 + 0.35;
  [-1, 1].forEach((s) => {
    lada(d.sten, 0.7, 1.9, 0.75, s * torn, 0.95, 0);
    const tak = new THREE.ConeGeometry(0.52, 0.55, 4);
    tak.rotateY(Math.PI / 4);
    tak.translate(s * torn, 2.17, 0);
    d.sten.push(tak);
  });
  lada(d.sten, LUCKA + 0.2, 0.45, 0.6, 0, 1.55, 0);
  [-0.6, 0, 0.6].forEach((x) => lada(d.sten, 0.3, 0.22, 0.6, x, 1.885, 0));
  if (stangd) {
    lada(d.tra, LUCKA, 1.3, 0.12, 0, 0.65, 0);
    [0.3, 1.0].forEach((y) => lada(d.rott, LUCKA + 0.02, 0.1, 0.18, 0, y, 0));
  } else {
    [-1, 1].forEach((s) => {
      const x = s * (LUCKA / 2 - 0.04), z = -(0.38 + LUCKA / 4);
      lada(d.tra, 0.08, 1.3, LUCKA / 2, x, 0.65, z);
      [0.3, 1.0].forEach((y) => lada(d.jarn, 0.11, 0.08, LUCKA / 2, x, y, z));
    });
  }
  return d;
}

// portbage är en port för teckenförklaringen, stängd eller öppen, och av
// kött på nivå 4.
function portbage(stangd: boolean, kott: boolean): THREE.Group {
  const g = new THREE.Group();
  const d = portGeometri(stangd);
  const lagg = (lista: THREE.BufferGeometry[], farg: string) =>
    lista.forEach((geo) => g.add(new THREE.Mesh(geo, new THREE.MeshLambertMaterial({ color: farg }))));
  lagg(d.sten, kott ? PORT_KOTT : FARG.port);
  lagg(d.tra, "#7a4e2c");
  lagg(d.jarn, FARG.jarn);
  lagg(d.rott, FARG.stangd);
  g.rotation.y = Math.PI / 4; // vänd mot kameran
  return g;
}

// dunklaShader gör en shader mörkare ju längre bort från fokus den ritar:
// rummet kartan följer ligger i facklans sken, grannrummen i dunkel och
// rummen längre bort nästan i mörker.
function dunklaShader(sh: THREE.WebGLProgramParametersWithUniforms, u: { uFokus: { value: THREE.Vector2 }; uRadie: { value: number } }): void {
  Object.assign(sh.uniforms, u);
  sh.vertexShader = sh.vertexShader
    .replace("#include <common>", "#include <common>\nvarying vec3 vDunk;")
    .replace("#include <begin_vertex>", "#include <begin_vertex>\nvDunk = (modelMatrix * vec4(transformed, 1.0)).xyz;");
  sh.fragmentShader = sh.fragmentShader
    .replace("#include <common>", "#include <common>\nvarying vec3 vDunk;\nuniform vec2 uFokus;\nuniform float uRadie;")
    .replace(
      "#include <dithering_fragment>",
      "gl_FragColor.rgb *= mix(1.0, 0.22, smoothstep(uRadie * 0.7, uRadie * 2.4, distance(vDunk.xz, uFokus)));\n#include <dithering_fragment>",
    );
}

// varldsUv låter en shader läsa sin textur efter platsen i världen i stället
// för efter ytans egna koordinater, så att murverket har lika stora block
// på alla murar, oavsett hur långa de är. skala är hur många enheter ett
// helt mönster täcker.
function varldsUv(sh: THREE.WebGLProgramParametersWithUniforms, skala: number): void {
  sh.vertexShader = sh.vertexShader
    .replace("#include <common>", "#include <common>\nvarying vec3 vVarld;\nvarying vec3 vNormV;")
    .replace("#include <begin_vertex>", "#include <begin_vertex>\nvVarld = (modelMatrix * vec4(transformed, 1.0)).xyz;\nvNormV = normalize(mat3(modelMatrix) * objectNormal);");
  sh.fragmentShader = sh.fragmentShader
    .replace("#include <common>", "#include <common>\nvarying vec3 vVarld;\nvarying vec3 vNormV;")
    .replace(
      "#include <map_fragment>",
      `vec2 wuv = abs(vNormV.y) > 0.5 ? vVarld.xz : vec2(abs(vNormV.x) > 0.5 ? vVarld.z : vVarld.x, vVarld.y);
diffuseColor *= texture2D(map, wuv / ${skala.toFixed(2)});`,
    );
}

export class IsoKarta {
  private renderer: THREE.WebGLRenderer;
  private scene = new THREE.Scene();
  private cam = new THREE.OrthographicCamera(-1, 1, 1, -1, 1, 400);
  private sol = new THREE.DirectionalLight(0xc8d2ea, 1.25); // svagt månljus genom gluggarna
  // Facklan i rummet kartan följer: ett varmt, fladdrande ljus.
  private fackla = new THREE.PointLight(0xffa64a, 0, 12, 1.5);
  private himmel = new THREE.HemisphereLight(0x9a9284, 0x0c0a08, 0.8);
  // Rum längre bort från rummet kartan följer ligger i dunkel (dunkla).
  private dunkel = { uFokus: { value: new THREE.Vector2() }, uRadie: { value: 1e5 } };
  private ram: HTMLDivElement; // pergamentkanten runt kartan
  private torn = new THREE.Group(); // golv, väggar och portar
  private rorligt = new THREE.Group(); // allt som ändras med varje tillstånd
  private ekon = new THREE.Group();
  private figurer = new THREE.Group();
  private aktorer = new Map<string, Aktor>();
  private tornNyckel = "";
  // Brickorna: rummen som har upptäckts, och de som faller just nu.
  private kanda = new Set<string>();
  private kandaStorlek = 0;
  private landar = new Map<string, { at: number; bricka: THREE.Group }>();
  private brickor = new THREE.Group();
  private brickGeo: THREE.BoxGeometry | null = null;
  // Ormleken bredvid guldhögen: leken med baksidan upp och de dragna korten
  // med framsidan upp (byggLeken).
  private leken = new THREE.Group();
  private lekNyckel = "";
  private lekNiva = -1;
  private lekDragna: THREE.Mesh[] = [];
  private lekTopp = new THREE.Vector3();
  private blandar = -Infinity;
  private state: Tillstand | null = null;
  private storleken = 0;
  private aktiv = false;
  private bredd = 0;
  private hojd = 0;
  private tackt = 0; // den del av bredden som loggen täcker (Val.tackt)
  private ringar: THREE.MeshBasicMaterial[] = [];
  private munnar: { upper: THREE.Group; lower: THREE.Group }[] = [];
  private ogon: THREE.Group[] = [];
  private guldljus: THREE.PointLight | null = null;
  private revbenMat = new THREE.MeshBasicMaterial({ map: K.revben(), color: 0xa03a2c, transparent: true, opacity: 0, depthWrite: false });
  // Väggarnas shader läser de här värdena varje bildruta.
  private vaggar = {
    uTid: { value: 0 },
    uPeriod: { value: 0 },
    uNiva: { value: 1 },
    uVagg: { value: srgb(LOOK[1].vagg) },
    uHi: { value: srgb(LOOK[1].hi) },
    uMitt: { value: new THREE.Vector2() },
  };
  private vaknade = -Infinity;
  private lag = false; // lägre kvalitet: inga skuggor och lägre upplösning
  private rutor: number[] = [];
  private senasteRuta = 0;
  private mitten: THREE.Vector3 | null = null; // kamerans mitt, som glider efter
  private senast = 0;
  private material = new Map<string, THREE.Material>();
  private plan = new Map<number, THREE.PlaneGeometry>();
  private canvas: HTMLCanvasElement;
  private val: Val;
  private kompass: HTMLDivElement; // N och en pil i hörnet nere till vänster
  private kant: HTMLDivElement; // pilarna vid kanten för det som ligger utanför
  private kantPilar: HTMLDivElement[] = [];
  // Tangenterna vid dörrarna: pilen på tangentbordet som går åt varje håll
  // ur rummet vandraren står i, och vilken av dem som är vald.
  private tangenter: HTMLDivElement;
  private tangent: HTMLSpanElement[] = [];
  private tangentRum: { x: number; y: number; oppen: boolean[]; vald: number; farg: string } | null = null;

  constructor(canvas: HTMLCanvasElement, val: Val) {
    this.canvas = canvas;
    this.val = val;
    this.renderer = new THREE.WebGLRenderer({ canvas, antialias: true });
    this.renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 2));
    this.renderer.setClearColor(0x0f0c09);
    this.renderer.shadowMap.enabled = true;
    this.renderer.shadowMap.type = THREE.PCFShadowMap;

    this.scene.add(this.himmel, this.fackla);
    this.sol.castShadow = true;
    this.sol.shadow.mapSize.set(2048, 2048);
    this.sol.shadow.bias = -0.0006;
    this.scene.add(this.sol, this.sol.target, this.torn, this.rorligt, this.ekon, this.figurer, this.brickor, this.leken);

    this.kompass = document.createElement("div");
    this.kompass.className = "isokompass";
    this.kompass.setAttribute("aria-hidden", "true");
    this.kompass.hidden = true;
    this.kompass.innerHTML = kompassSvg();
    this.kant = document.createElement("div");
    this.kant.className = "isokant";
    this.kant.setAttribute("aria-hidden", "true");
    this.kant.hidden = true;
    this.tangenter = document.createElement("div");
    this.tangenter.className = "isotangenter";
    this.tangenter.setAttribute("aria-hidden", "true"); // pekknapparna har samma val med namn
    this.tangenter.hidden = true;
    // Pilen ritas av stilen: tangentbordets pil (data-pil), eller på en
    // pekskärm samma snedställda pil som knapparna (data-pek). data-key är
    // tangenten, så att ett tryck på tecknet går åt det hållet.
    const pek = ["↗", "↘", "↙", "↖"];
    const tangent = ["ArrowUp", "ArrowRight", "ArrowDown", "ArrowLeft"];
    ["↑", "→", "↓", "←"].forEach((pil, d) => {
      const t = document.createElement("span");
      t.dataset.pil = pil;
      t.dataset.pek = pek[d];
      t.dataset.key = tangent[d];
      t.hidden = true;
      this.tangenter.appendChild(t);
      this.tangent.push(t);
    });
    this.ram = document.createElement("div");
    this.ram.className = "isoram";
    this.ram.setAttribute("aria-hidden", "true");
    this.ram.hidden = true;
    canvas.after(this.ram, this.kompass, this.kant, this.tangenter);

    new ResizeObserver(() => this.anpassa()).observe(canvas.parentElement!);
    const nav = navigator as Navigator & { deviceMemory?: number };
    const tvinga = new URLSearchParams(location.search).get("kvalitet") === "lag";
    if (tvinga || (nav.hardwareConcurrency || 8) <= 2 || (nav.deviceMemory || 8) <= 2) this.sankKvalitet();
    const loop = (t: number) => {
      this.rita(t);
      requestAnimationFrame(loop);
    };
    requestAnimationFrame(loop);
  }

  // visa slår på eller av kartan.
  visa(pa: boolean): void {
    this.aktiv = pa;
    this.canvas.hidden = !pa;
    this.kompass.hidden = !pa;
    this.ram.hidden = !pa;
    this.kant.hidden = !pa;
    this.tangenter.hidden = !pa;
    if (pa) this.anpassa();
  }

  storlek(n: number): void {
    this.storleken = n;
    this.tornNyckel = "";
  }

  anpassa(): void {
    const box = this.canvas.parentElement!.getBoundingClientRect();
    this.bredd = Math.floor(box.width);
    this.hojd = Math.floor(box.height);
    if (this.bredd > 0 && this.hojd > 0) this.renderer.setSize(this.bredd, this.hojd);
  }

  // visaTillstand bygger om det som har ändrats sedan förra tillståndet.
  visaTillstand(state: Tillstand): void {
    this.state = state;
    if (state.size) this.storleken = state.size;

    // Nya rum läggs som brickor (laggBricka). Ett nytt torn, eller många
    // rum på en gång, som när vandringen är slut och allt syns, läggs utan.
    const nu = new Set(state.rooms.map((r) => r.x + "," + r.y));
    const nyttTorn = this.kandaStorlek !== this.storleken || [...this.kanda].some((k) => !nu.has(k));
    if (nyttTorn) {
      this.kanda.clear();
      this.landar.forEach((l) => this.brickor.remove(l.bricka));
      this.landar.clear();
      this.kandaStorlek = this.storleken;
    }
    const nya = state.rooms.filter((r) => !this.kanda.has(r.x + "," + r.y));
    const lagg = !nyttTorn && this.kanda.size > 0 && nya.length <= 3 && !LUGN;
    nya.forEach((r) => {
      this.kanda.add(r.x + "," + r.y);
      if (lagg) this.laggBricka(r);
    });
    this.byggOm();
    this.byggRorligt(state);
  }

  // byggOm bygger tornet på nytt om något har ändrats, utan rummen vars
  // brickor ännu inte har landat.
  private byggOm(): void {
    const state = this.state;
    if (!state) return;
    const niva = this.niva();
    const slutet = niva === 4 && state.phase === "slut";
    const rum = state.rooms.filter((r) => !this.landar.has(r.x + "," + r.y));
    const nyckel = this.storleken + "|" + (slutet ? 1 : 0) + "|" +
      rum.map((r) => [r.x, r.y, r.d.map(Number).join(""), r.k, r.exit, r.stangd ? 1 : 0, (r.in || []).join("+")].join(",")).join(";");
    if (nyckel !== this.tornNyckel) {
      this.tornNyckel = nyckel;
      this.byggTorn({ ...state, rooms: rum }, slutet);
    }
  }

  // laggBricka låter ett nytt rum falla ner som en spelbricka: den kommer
  // uppifrån med baksidan upp, vänds i luften och landar på sin plats. När
  // den har landat byggs rummet, med murar och inredning.
  private laggBricka(r: Rum): void {
    const g = this.brickGeo || (this.brickGeo = new THREE.BoxGeometry(S * 0.98, BRICKA, S * 0.98));
    const kant = this.enkelt("brickkant", "#3a2c20");
    const bak = this.enkelt("brickbak", "#c9b48a");
    // Lådans sidor: +x, -x, +y (ovansidan), -y (baksidan), +z, -z.
    const m = new THREE.Mesh(g, [kant, kant, this.golvMaterial(), bak, kant, kant]);
    m.castShadow = true;
    const bricka = new THREE.Group();
    bricka.add(m);
    bricka.position.set(r.x * S + S / 2, 6, r.y * S + S / 2);
    bricka.rotation.x = Math.PI;
    this.brickor.add(bricka);
    this.landar.set(r.x + "," + r.y, { at: performance.now(), bricka });
  }

  // landa flyttar brickorna som faller, och bygger rummen vars brickor har landat.
  private landa(t: number): void {
    if (!this.landar.size) return;
    let klara = false;
    this.landar.forEach((l, k) => {
      const f = Math.min(1, (t - l.at) / LANDA_MS);
      if (f >= 1) {
        this.brickor.remove(l.bricka);
        this.landar.delete(k);
        klara = true;
        return;
      }
      const fall = f * f; // faller fortare och fortare
      l.bricka.position.y = 6 * (1 - fall) - BRICKA / 2 * fall;
      l.bricka.rotation.x = Math.PI * (1 - Math.min(1, f * 1.4)); // vänds innan den landar
    });
    if (klara && this.state) {
      this.byggOm();
      this.byggRorligt(this.state);
    }
  }

  eko(x: number, y: number, slot: number): void {
    const m = new THREE.Mesh(this.yta(S * 0.96), new THREE.MeshBasicMaterial({
      color: this.val.farger[slot] || "#ffffff", transparent: true, opacity: 0.22, depthWrite: false,
    }));
    m.rotation.x = -Math.PI / 2;
    m.position.set(x * S + S / 2, 0.02, y * S + S / 2);
    m.userData.at = performance.now();
    this.ekon.add(m);
  }

  vakna(): void {
    this.vaknade = performance.now();
    // De dragna korten flyger tillbaka in i leken, som blandas om.
    if (this.lekDragna.length && this.lekNiva >= 0) this.blandar = this.vaknade;
  }

  private niva(): number {
    const s = this.state;
    return s && s.reveal ? Math.min(4, Math.max(1, s.level || 1)) : 1;
  }

  // ---- Delade geometrier och material ----

  private yta(storlek: number): THREE.PlaneGeometry {
    let g = this.plan.get(storlek);
    if (!g) {
      g = new THREE.PlaneGeometry(storlek, storlek);
      this.plan.set(storlek, g);
    }
    return g;
  }

  private kritMaterial(tex: THREE.Texture, farg: string, nyckel: string, opacity: number): THREE.Material {
    const k = nyckel + farg + opacity;
    let m = this.material.get(k);
    if (!m) {
      m = new THREE.MeshBasicMaterial({ map: tex, color: farg, transparent: true, opacity, depthWrite: false });
      this.material.set(k, m);
    }
    return m;
  }

  // golvkrita lägger en textur platt på golvet, vriden vinkel radianer från
  // öster mot norr.
  private golvkrita(grupp: THREE.Group, tex: THREE.Texture, nyckel: string, farg: string, storlek: number, x: number, z: number, vinkel = 0, y = 0.016, opacity = 1): THREE.Mesh {
    const m = new THREE.Mesh(this.yta(storlek), this.kritMaterial(tex, farg, nyckel, opacity));
    m.rotation.set(-Math.PI / 2, 0, vinkel);
    m.position.set(x, y, z);
    m.renderOrder = 1;
    grupp.add(m);
    return m;
  }

  // ---- Tornet ----

  private byggTorn(state: Tillstand, slutet: boolean): void {
    this.torn.traverse((o) => {
      if (o instanceof THREE.Mesh && o.userData.egen) o.geometry.dispose();
    });
    this.torn.clear();
    this.ringar = [];
    this.munnar = [];
    this.ogon = [];

    const golv: THREE.BufferGeometry[] = [];
    const vaggar: THREE.BufferGeometry[] = [];
    const portar: THREE.BufferGeometry[] = [];
    const stangda: THREE.BufferGeometry[] = [];
    const tra: THREE.BufferGeometry[] = [];
    const jarn: THREE.BufferGeometry[] = [];
    const flammor: THREE.BufferGeometry[] = [];
    const inredning = new Map<string, THREE.BufferGeometry[]>();
    const inreda: I.Lagg = (m, g) => {
      const l = inredning.get(m);
      if (l) l.push(g);
      else inredning.set(m, [g]);
    };
    const sken: THREE.BufferGeometry[] = [];
    // En fackla på muren vid (x, z), som pekar in i rummet åt (nx, nz):
    // ett järnfäste, en käpp, en låga och ett sken på golvet.
    const fackla = (x: number, z: number, nx: number, nz: number) => {
      const b = new THREE.BoxGeometry(0.13, 0.13, 0.13);
      b.translate(x + nx * 0.05, 0.95, z + nz * 0.05);
      jarn.push(b);
      const st = new THREE.CylinderGeometry(0.04, 0.028, 0.4, 6);
      if (nz) st.rotateX(0.4 * nz);
      else st.rotateZ(-0.4 * nx);
      st.translate(x + nx * 0.14, 1.12, z + nz * 0.14);
      tra.push(st);
      const fl = new THREE.ConeGeometry(0.1, 0.32, 7);
      fl.translate(x + nx * 0.22, 1.46, z + nz * 0.22);
      flammor.push(fl);
      const p = new THREE.PlaneGeometry(2.8, 2.8);
      p.rotateX(-Math.PI / 2);
      p.translate(x + nx * 1.1, 0.03, z + nz * 1.1);
      sken.push(p);
    };
    const kanter = new Map<string, { oppen: boolean; port: boolean; stangd: boolean; h: boolean; x: number; z: number }>();

    const kant = (nyckel: string, h: boolean, x: number, z: number) => {
      let k = kanter.get(nyckel);
      if (!k) {
        k = { oppen: false, port: false, stangd: false, h, x, z };
        kanter.set(nyckel, k);
      }
      return k;
    };

    state.rooms.forEach((r: Rum) => {
      const g = new THREE.PlaneGeometry(S, S);
      g.rotateX(-Math.PI / 2);
      g.rotateY(((r.x * 7 + r.y * 3) % 4) * (Math.PI / 2)); // plattorna ligger olika i olika rum
      g.translate(r.x * S + S / 2, 0, r.y * S + S / 2);
      golv.push(g);
      // Revbenen syns genom golvet från nivå 3.
      if (r.k !== "kammare") {
        const rev = new THREE.Mesh(this.yta(S * 0.8), this.revbenMat);
        rev.rotation.set(-Math.PI / 2, 0, (r.x + r.y) % 2 ? 0.3 : -0.4);
        rev.position.set(r.x * S + S / 2, 0.008, r.y * S + S / 2);
        this.torn.add(rev);
      }
      // Sidorna: N, Ö, S, V.
      const sidor = [
        kant("h" + r.x + "," + r.y, true, r.x * S, r.y * S),
        kant("v" + (r.x + 1) + "," + r.y, false, (r.x + 1) * S, r.y * S),
        kant("h" + r.x + "," + (r.y + 1), true, r.x * S, (r.y + 1) * S),
        kant("v" + r.x + "," + r.y, false, r.x * S, r.y * S),
      ];
      sidor.forEach((k, i) => {
        if (r.d[i]) k.oppen = true;
        if (r.exit === i) {
          k.port = true;
          k.stangd = !!r.stangd || slutet;
        }
      });
      if (r.k === "kammare") this.byggKammare(r);
      (r.in || []).forEach((namn) => I.bygg(namn, inreda, r.x * S, r.y * S, S, r.x * 31 + r.y * 17 + 1));
      // Facklorna sitter på murarna som vetter mot kameran, norra och västra.
      fackla(r.x * S + S * 0.22, r.y * S + VT / 2, 0, 1);
      fackla(r.x * S + VT / 2, r.y * S + S * 0.78, 1, 0);
    });

    const lada = (lista: THREE.BufferGeometry[], w: number, h: number, d: number, x: number, y: number, z: number) => {
      const g = new THREE.BoxGeometry(w, h, d);
      g.translate(x, y, z);
      lista.push(g);
    };
    // at lägger en låda längs kanten k: u längs muren (0 till S), v ut från
    // den, åt -z eller -x.
    const at = (k: { h: boolean; x: number; z: number }, lista: THREE.BufferGeometry[], u: number, y: number, w: number, h: number, d: number, v = 0) =>
      k.h ? lada(lista, w, h, d, k.x + u, y, k.z + v) : lada(lista, d, h, w, k.x + v, y, k.z + u);
    // Ett murstycke längs kanten från a till b, med krön ovanpå.
    const stycke = (k: { h: boolean; x: number; z: number }, a: number, b: number) => {
      const len = b - a + VT;
      at(k, vaggar, (a + b) / 2, VH / 2, len, VH, VT);
      const antal = Math.floor((len - 0.2) / 0.9);
      for (let i = 0; i < antal; i++) at(k, vaggar, a - VT / 2 + 0.35 + i * 0.9 + (len - 0.2 - antal * 0.9) / 2, VH + 0.11, 0.42, 0.22, VT);
    };
    kanter.forEach((k, nyckel) => {
      if (!k.oppen && !k.port) {
        stycke(k, 0, S);
        return;
      }
      stycke(k, 0, S / 2 - LUCKA / 2);
      stycke(k, S / 2 + LUCKA / 2, S);
      if (!k.port) {
        // En dörr mellan två rum: dörrposter, och en dörr av trä som står
        // uppslagen mot muren i det ena rummet.
        [-1, 1].forEach((sida) => at(k, vaggar, S / 2 + sida * (LUCKA / 2 + 0.06), (VH + 0.3) / 2, 0.3, VH + 0.3, VT + 0.12));
        let h = 0;
        for (const c of nyckel) h = (h * 31 + c.charCodeAt(0)) | 0;
        const sida = h & 1 ? 1 : -1;
        const u = S / 2 + LUCKA / 2 + 0.25 + 0.85;
        at(k, tra, u, 0.39, 1.7, 0.78, 0.07, sida * (VT / 2 + 0.05));
        [0.2, 0.6].forEach((y) => at(k, jarn, u, y, 1.7, 0.07, 0.03, sida * (VT / 2 + 0.1)));
        return;
      }
      // Porten: ett portvalv med dörrar av trä, se portGeometri.
      const d = portGeometri(k.stangd);
      const ut = k.h ? (k.z === 0 ? 0 : Math.PI) : k.x === 0 ? Math.PI / 2 : -Math.PI / 2;
      const flytta = (lista: THREE.BufferGeometry[], till: THREE.BufferGeometry[]) =>
        lista.forEach((g) => {
          g.rotateY(ut);
          g.translate(k.h ? k.x + S / 2 : k.x, 0, k.h ? k.z : k.z + S / 2);
          till.push(g);
        });
      flytta(d.sten, portar);
      flytta(d.tra, tra);
      flytta(d.jarn, jarn);
      flytta(d.rott, stangda);
      if (!k.stangd) this.mun(k);
    });

    const lagg = (lista: THREE.BufferGeometry[], mat: THREE.Material, skugga: boolean) => {
      if (!lista.length) return;
      const m = new THREE.Mesh(mergeGeometries(lista), mat);
      lista.forEach((g) => g.dispose());
      m.userData.egen = true;
      m.castShadow = skugga;
      m.receiveShadow = true;
      this.torn.add(m);
    };
    lagg(golv, this.golvMaterial(), false);
    lagg(vaggar, this.vaggMaterial(), true);
    lagg(portar, this.stenMaterial("port", FARG.port), true);
    lagg(stangda, this.enkelt("stangd", FARG.stangd), true);
    lagg(tra, this.traMaterial(), true);
    lagg(jarn, this.enkelt("jarn", FARG.jarn), true);
    lagg(flammor, this.flammaMaterial, false);
    inredning.forEach((lista, m) => {
      if (m === "lys") lagg(lista, this.flammaMaterial, false);
      else if (m === "sten") lagg(lista, this.stenMaterial("inredning", "#cfc6b4"), true);
      else if (m === "tra") lagg(lista, this.traMaterial(), true);
      else if (m === "jarn") lagg(lista, this.enkelt("jarn", FARG.jarn), true);
      else lagg(lista, this.enkelt("f" + m, m), true);
    });
    lagg(sken, this.skenMaterial, false);
    this.torn.children.forEach((o) => {
      if (o instanceof THREE.Mesh && o.material === this.skenMaterial) o.renderOrder = 1;
    });
  }

  // mun ger en öppen port tänder, som syns och tuggar på nivå 4.
  private mun(k: { h: boolean; x: number; z: number }): void {
    const mun = new THREE.Group();
    mun.position.set(k.h ? k.x + S / 2 : k.x, 0, k.h ? k.z : k.z + S / 2);
    if (!k.h) mun.rotation.y = Math.PI / 2;
    const tand = this.tandGeo || (this.tandGeo = new THREE.ConeGeometry(0.1, 0.36, 5));
    const mat = this.enkelt("tand", "#eadfc8");
    const upper = new THREE.Group(), lower = new THREE.Group();
    for (let i = 0; i < 5; i++) {
      const u = new THREE.Mesh(tand, mat);
      u.rotation.x = Math.PI;
      u.position.set(-0.7 + i * 0.35, 1.17, 0);
      upper.add(u);
      const d = new THREE.Mesh(tand, mat);
      d.position.set(-0.53 + i * 0.35, 0.18, 0);
      lower.add(d);
    }
    mun.add(upper, lower);
    mun.visible = false;
    this.torn.add(mun);
    this.munnar.push({ upper, lower });
  }
  private tandGeo: THREE.ConeGeometry | null = null;
  private flammaMaterial = new THREE.MeshBasicMaterial({ color: 0xffb04a });
  private skenMaterial = new THREE.MeshBasicMaterial({
    map: K.sken(), color: 0xff8a30, transparent: true, opacity: 0.3, depthWrite: false, blending: THREE.AdditiveBlending,
  });

  private golvMaterial(): THREE.Material {
    let m = this.material.get("golv");
    if (!m) {
      const mat = new THREE.MeshLambertMaterial({ map: K.stengolv() });
      mat.onBeforeCompile = (sh) => dunklaShader(sh, this.dunkel);
      m = mat;
      this.material.set("golv", m);
    }
    return m;
  }

  private vaggMaterial(): THREE.Material {
    let m = this.material.get("vagg");
    if (!m) {
      // Väggarnas färg räknas ut i shadern, så att hjärtslaget kan gå som en
      // våg från tornets mitt: samma puls som i puls(), för varje punkt.
      const mat = new THREE.MeshLambertMaterial({ map: K.murverk(), color: FARG.vagg });
      const u = this.vaggar;
      mat.onBeforeCompile = (sh) => {
        Object.assign(sh.uniforms, u);
        sh.fragmentShader = sh.fragmentShader
          .replace(
            "#include <common>",
            `#include <common>
uniform float uTid;
uniform float uPeriod;
uniform float uNiva;
uniform vec3 uVagg;
uniform vec3 uHi;
uniform vec2 uMitt;
float g(float ph, float c) { float d = (ph - c) / 0.07; return exp(-d * d); }
float puls() {
  if (uPeriod <= 0.0) return 0.0;
  float t = uTid - length(vVarld.xz - uMitt) * ${MS_PER_ENHET.toFixed(1)};
  if (uNiva < 2.5) return 0.5 + 0.5 * sin(6.2831853 * t / uPeriod);
  float ph = mod(t, uPeriod) / uPeriod;
  return max(max(g(ph, 0.0), 0.7 * g(ph, 0.22)), g(ph, 1.0));
}`,
          )
          .replace("vec4 diffuseColor = vec4( diffuse, opacity );", "vec4 diffuseColor = vec4( mix(uVagg, uHi, puls()), opacity );");
        varldsUv(sh, 1.6); // efter, så att vVarld står först
        dunklaShader(sh, this.dunkel);
      };
      m = mat;
      this.material.set("vagg", m);
    }
    return m;
  }

  // stenMaterial är murverk i färgen farg, för portvalven.
  private stenMaterial(nyckel: string, farg: string): THREE.Material {
    let m = this.material.get(nyckel);
    if (!m) {
      const mat = new THREE.MeshLambertMaterial({ map: K.murverk(), color: farg });
      mat.onBeforeCompile = (sh) => {
        varldsUv(sh, 1.6);
        dunklaShader(sh, this.dunkel);
      };
      m = mat;
      this.material.set(nyckel, m);
    }
    return m;
  }

  private traMaterial(): THREE.Material {
    let m = this.material.get("tra");
    if (!m) {
      const mat = new THREE.MeshLambertMaterial({ map: K.tra() });
      mat.onBeforeCompile = (sh) => dunklaShader(sh, this.dunkel);
      m = mat;
      this.material.set("tra", m);
    }
    return m;
  }

  private enkelt(nyckel: string, farg: string): THREE.Material {
    let m = this.material.get(nyckel);
    if (!m) {
      const mat = new THREE.MeshLambertMaterial({ color: farg });
      mat.onBeforeCompile = (sh) => dunklaShader(sh, this.dunkel);
      m = mat;
      this.material.set(nyckel, m);
    }
    return m;
  }

  private ogaGeo: THREE.SphereGeometry | null = null;

  // Kammaren: en kritring kring en hög av guld.
  private byggKammare(r: Rum): void {
    const cx = r.x * S + S / 2, cz = r.y * S + S / 2;
    const ringMat = new THREE.MeshBasicMaterial({ map: K.ring(), color: RING[this.niva()], transparent: true, depthWrite: false });
    const ring = new THREE.Mesh(this.yta(S * 0.75), ringMat);
    ring.rotation.x = -Math.PI / 2;
    ring.position.set(cx, 0.012, cz);
    this.torn.add(ring);
    this.ringar.push(ringMat);

    // Ormens ögon över guldet, som lyser när den vaknar. De står bredvid
    // varandra som kameran ser dem.
    const ogon = new THREE.Group();
    const ogonMat = this.material.get("ögon") || new THREE.MeshBasicMaterial({ color: 0xffe14a });
    this.material.set("ögon", ogonMat);
    const oga = this.ogaGeo || (this.ogaGeo = new THREE.SphereGeometry(0.11, 10, 8));
    [-1, 1].forEach((sida) => {
      const e = new THREE.Mesh(oga, ogonMat);
      e.scale.set(1, 0.7, 1);
      e.position.set(sida * 0.25, 0, -sida * 0.25);
      ogon.add(e);
    });
    ogon.position.set(cx, 1.6, cz);
    ogon.visible = false;
    this.torn.add(ogon);
    this.ogon.push(ogon);
    if (!this.guldljus) {
      this.guldljus = new THREE.PointLight(0xffc94a, 0, 10, 2);
      this.scene.add(this.guldljus);
    }
    this.guldljus.position.set(cx, 1.6, cz);

    const hog = guldhog();
    hog.position.set(cx, 0, cz);
    hog.userData.egen = true;
    this.torn.add(hog);
  }

  // ---- Ormleken ----

  // byggLeken lägger ormleken bredvid guldhögen. Den som står i kammaren
  // (niva 0–2, hur djupt ormen sover) ser leken krympa och de dragna korten
  // växa, i tre steg som sömnmätaren, så att högarna inte avslöjar mer än
  // den. Andra ser en hög som inte ändras.
  private byggLeken(r: Rum | undefined, niva: number): void {
    const nyckel = r ? r.x + "," + r.y + "," + niva : "";
    if (nyckel === this.lekNyckel || this.blandar > -Infinity) return;
    this.lekNyckel = nyckel;
    this.lekNiva = niva;
    this.leken.clear();
    this.lekDragna = [];
    if (!r) return;
    const cx = r.x * S + S / 2, cz = r.y * S + S / 2;
    const lek = niva < 0 ? 8 : [14, 8, 3][Math.min(2, niva)];
    const dragna = niva < 0 ? 3 : [2, 7, 12][Math.min(2, niva)];
    const kort = this.kortGeo || (this.kortGeo = new THREE.BoxGeometry(0.66, 0.028, 0.94));
    const kant = this.enkelt("kortkant", "#d8c8a0");
    const bak = this.material.get("kortbak") || new THREE.MeshLambertMaterial({ map: K.kortbaksida() });
    const fram = this.material.get("kortfram") || new THREE.MeshLambertMaterial({ map: K.kortframsida() });
    this.material.set("kortbak", bak);
    this.material.set("kortfram", fram);
    const lagg = (n: number, x: number, z: number, topp: THREE.Material, dragen: boolean) => {
      let s = (x * 13 + z * 7) | 0;
      const rnd = () => ((s = (s * 16807 + 11) % 2147483647) / 2147483647);
      for (let i = 0; i < n; i++) {
        // Lådans sidor: +x, -x, +y (ovansidan), -y, +z, -z.
        const m = new THREE.Mesh(kort, [kant, kant, i === n - 1 || dragen ? topp : bak, kant, kant, kant]);
        m.position.set(x + (rnd() - 0.5) * 0.05, 0.015 + i * 0.03, z + (rnd() - 0.5) * 0.05);
        m.rotation.y = Math.PI / 4 + (rnd() - 0.5) * (dragen ? 0.5 : 0.08);
        m.castShadow = i === n - 1;
        this.leken.add(m);
        if (dragen) this.lekDragna.push(m);
      }
      return new THREE.Vector3(x, 0.015 + n * 0.03, z);
    };
    this.lekTopp = lagg(lek, cx - 0.7, cz + 1.35, bak, false);
    lagg(dragna, cx + 0.75, cz + 1.35, fram, true);
  }
  private kortGeo: THREE.BoxGeometry | null = null;

  // blanda låter de dragna korten flyga i en båge tillbaka in i leken när
  // ormen har vaknat. Sedan byggs högarna om efter den nya leken.
  private blanda(t: number): void {
    if (this.blandar === -Infinity) return;
    const tid = t - this.blandar;
    const n = this.lekDragna.length;
    this.lekDragna.forEach((m, i) => {
      if (!m.userData.fran) m.userData.fran = m.position.clone();
      const f = Math.max(0, Math.min(1, (tid - (n - 1 - i) * 60) / 450));
      const fran = m.userData.fran as THREE.Vector3;
      m.position.lerpVectors(fran, this.lekTopp, f);
      m.position.y += Math.sin(f * Math.PI) * 0.8;
      m.rotation.z = f * Math.PI;
    });
    if (tid > (n - 1) * 60 + 500) {
      this.blandar = -Infinity;
      this.lekNyckel = "";
      if (this.state) this.byggRorligt(this.state);
    }
  }

  // ---- Teckenförklaringen ----

  private typer = new Map<string, { typ: string; farg?: string; storlek?: number }>(); // varelsernas tecken → id och modell
  private fallTyper = new Map<string, string>(); // fällornas tecken → id
  private bilder = new Map<string, string | null>();
  private bildScen: THREE.Scene | null = null;
  private bildMal: THREE.WebGLRenderTarget | null = null;

  // bild ger en liten bild av modellen för ett tecken i teckenförklaringen,
  // som en data-URL, eller null där tecknet själv står på kartan (krita på
  // golvet och liknande). Bilderna ritas en gång och sparas.
  bild(nyckel: string, glyph: string, farg: string): string | null {
    const slag = nyckel.split(":")[0];
    let gor: (() => THREE.Object3D | null) | null = null;
    if (slag === "du" || slag === "vandrare" || slag === "borta") gor = () => M.miniatyr(M.vandrare("", farg), farg).grupp;
    else if (slag === "varelse") {
      const v = this.typer.get(glyph);
      if (!v) return null;
      gor = () => {
        const f = M.varelse(v.typ, v.farg);
        return f ? M.miniatyr(f, KANT_VARELSE).grupp : null;
      };
    } else if (slag === "falla") {
      const id = this.fallTyper.get(glyph);
      if (!id) return null;
      gor = () => M.falla(id);
    } else if (slag === "kammare") gor = () => guldhog();
    else if (slag === "port" || slag === "port_stangd") gor = () => portbage(slag === "port_stangd" || (this.niva() === 4 && this.state?.phase === "slut"), this.niva() === 4);
    if (!gor) return null;
    const v = this.typer.get(glyph);
    const k = slag + "|" + glyph + "|" + farg + "|" + (slag.startsWith("port") ? this.niva() : "") + "|" + (slag === "varelse" && v ? v.farg || "" : "");
    if (this.bilder.has(k)) return this.bilder.get(k)!;
    const o = gor();
    const url = o ? this.ritaBild(o) : null;
    this.bilder.set(k, url);
    return url;
  }

  private ritaBild(o: THREE.Object3D): string {
    const N = 96;
    if (!this.bildScen) {
      this.bildScen = new THREE.Scene();
      this.bildScen.add(new THREE.HemisphereLight(0xdfe4ee, 0x2a2c30, 1.6));
      const sol = new THREE.DirectionalLight(0xfff0dc, 2.0);
      sol.position.set(-4, 8, 6);
      this.bildScen.add(sol);
      this.bildMal = new THREE.WebGLRenderTarget(N, N, { samples: 4 });
      this.bildMal.texture.colorSpace = THREE.SRGBColorSpace;
    }
    const scen = this.bildScen;
    scen.add(o);
    const box = new THREE.Box3().setFromObject(o);
    const mitt = box.getCenter(new THREE.Vector3());
    const storlek = box.getSize(new THREE.Vector3());
    const h = Math.max(storlek.x, storlek.y, storlek.z) * 0.62;
    const cam = new THREE.OrthographicCamera(-h, h, h, -h, 0.1, 100);
    cam.position.copy(mitt).addScaledVector(FRAN, 20);
    cam.lookAt(mitt);
    const r = this.renderer;
    const forra = r.getRenderTarget();
    const farg = r.getClearColor(new THREE.Color());
    const alfa = r.getClearAlpha();
    r.setClearColor(0x000000, 0);
    r.setRenderTarget(this.bildMal);
    r.clear();
    r.render(scen, cam);
    const px = new Uint8Array(N * N * 4);
    r.readRenderTargetPixels(this.bildMal!, 0, 0, N, N, px);
    r.setRenderTarget(forra);
    r.setClearColor(farg, alfa);
    scen.remove(o);
    // Bilden ligger upp och ner i målet.
    const c = document.createElement("canvas");
    c.width = c.height = N;
    const ctx = c.getContext("2d")!;
    const im = ctx.createImageData(N, N);
    for (let y = 0; y < N; y++) im.data.set(px.subarray((N - 1 - y) * N * 4, (N - y) * N * 4), y * N * 4);
    ctx.putImageData(im, 0, 0);
    return c.toDataURL();
  }

  // ---- Det som rör sig ----

  private byggRorligt(state: Tillstand): void {
    this.rorligt.clear();
    const g = this.rorligt;
    const kammare = state.rooms.find((r) => r.k === "kammare" && !this.landar.has(r.x + "," + r.y));
    this.byggLeken(kammare, state.sleep ?? -1);
    const landar = (r: Rum) => this.landar.has(r.x + "," + r.y);
    const kanda = new Set(state.rooms.map((r) => r.x + "," + r.y));

    // Okända rum är en svag prick. En vandrare ser bara prickarna bredvid
    // rummen den känner till, så att kartan inte avslöjar mer.
    const nara = (x: number, y: number) => [[0, 0], [1, 0], [-1, 0], [0, 1], [0, -1]].some(([dx, dy]) => kanda.has(x + dx + "," + (y + dy)));
    for (let y = 0; y < state.size; y++) {
      for (let x = 0; x < state.size; x++) {
        if (kanda.has(x + "," + y) || (state.you && !nara(x, y))) continue;
        const mitt = x === Math.floor(state.size / 2) && y === x;
        this.golvkrita(g, mitt ? K.liten() : K.prick(), mitt ? "liten" : "prick", FARG.okand, mitt ? 1.1 : 0.8, x * S + S / 2, y * S + S / 2);
      }
    }

    state.rooms.forEach((r) => {
      if (landar(r)) return; // brickan har inte landat än
      if (r.mark) {
        // En fälla som har slagit till syns som en modell, eller som sitt
        // tecken i krita om den saknar modell.
        const [x, z] = plats(r.x, r.y, 3, 2);
        const modell = r.marki ? M.falla(r.marki) : null;
        if (modell) {
          this.fallTyper.set(r.mark, r.marki!);
          modell.position.set(x, 0, z);
          modell.rotation.y = Math.PI / 4;
          g.add(modell);
        } else this.golvkrita(g, K.tecken(r.mark), "t" + r.mark, FARG.mark, 1.5, x, z);
      }
      if (r.f) {
        const [x, z] = plats(r.x, r.y, RW - 1, 1);
        this.golvkrita(g, K.tecken("+"), "t+", FARG.guld, 0.9, x, z);
      }
      if (r.dod) {
        // Kartritaren ser var någon har dött: ett litet x i hörnet.
        const [x, z] = plats(r.x, r.y, RW - 1, RH - 1);
        this.golvkrita(g, K.tecken("x"), "tx", FARG.forlorad, 0.9, x, z);
      }
    });

    // Rummet vandraren står i: ljusare golv
    // och en ram i vandrarens färg.
    const jag = state.players.find((p) => p.id === state.you);
    if (jag && jag.status === "inne" && !jag.hidden && jag.x >= 0) {
      const farg = this.val.farger[jag.slot] || "#ffffff";
      this.golvkrita(g, K.ram(), "ram", farg, S * 0.94, jag.x * S + S / 2, jag.y * S + S / 2, 0, 0.011);
    }

    // Där en dold vandrare sågs senast står en svag ring i dess färg.
    state.players.forEach((p) => {
      const sedd = state.lastSeen && state.lastSeen[p.id];
      if (!sedd || !p.hidden) return;
      const [sx, sy] = SLOTS[sedd.slot] || SLOTS[0];
      const [x, z] = plats(sedd.x, sedd.y, sx, sy);
      this.golvkrita(g, K.liten(), "liten", this.val.farger[sedd.slot] || "#ffffff", 0.8, x, z, 0, 0.016, 0.45);
    });

    // Vandrarna och varelserna. De döda ligger kvar, med ett x i krita.
    const sedda = new Set<string>();
    state.players.forEach((p) => {
      if (p.status === "ute" || p.hidden || p.x < 0) return;
      const [sx, sy] = SLOTS[p.slot] || SLOTS[0];
      const [x, z] = plats(p.x, p.y, sx, sy);
      const dod = p.status === "forlorad" || p.status === "dod";
      if (dod) this.golvkrita(g, K.tecken("x"), "tx", FARG.forlorad, 0.9, x, z);
      const farg = p.online && !dod ? this.val.farger[p.slot] || "#ffffff" : FARG.borta;
      const klass = p.klassId || "";
      // Föremålen syns bara på den egna figuren. Andras väskor är hemliga.
      const foremal = p.id === state.you ? (state.foremal || []).map((f) => f.id) : [];
      const km = p.klassModell || {};
      const a = this.flytta(
        "v" + p.id,
        klass + farg + "|" + foremal.join(",") + "|" + (km.farg || "") + "|" + (km.storlek || 1),
        () => M.miniatyr(skala(M.vandrare(klass, farg, foremal, km.farg), km.storlek), farg, km.storlek),
        p.x, p.y, x, z,
      );
      a.dod = dod;
      this.egenRing(a, p.id === state.you && !dod ? farg : "");
      sedda.add("v" + p.id);
    });
    (state.monsters || []).forEach((m) => {
      const [x, z] = plats(m.x, m.y, VARELSE[0], VARELSE[1]);
      const typ = m.typ || "";
      const mm = m.modell || {};
      if (typ) this.typer.set(m.glyph, { typ, farg: mm.farg, storlek: mm.storlek });
      this.flytta(
        "m" + m.id,
        "m" + typ + m.glyph + "|" + (mm.farg || "") + "|" + (mm.storlek || 1),
        () => M.miniatyr(skala(M.varelse(typ, mm.farg) || this.tecknad(m.glyph), mm.storlek), KANT_VARELSE, mm.storlek),
        m.x, m.y, x, z,
      );
      sedda.add("m" + m.id);
    });
    this.aktorer.forEach((a, id) => {
      if (sedda.has(id)) return;
      this.figurer.remove(a.hall);
      this.aktorer.delete(id);
    });

    const riktning: Record<string, number> = { N: 0, E: 1, S: 2, W: 3 };
    // Det egna draget visas på tangenten vid dörren (ritaTangenter).
    const d = state.action ? riktning[state.action] : undefined;
    const rum = jag && jag.status === "inne" && !jag.hidden && jag.x >= 0 && state.phase === "spel"
      ? state.rooms.find((r) => r.x === jag.x && r.y === jag.y)
      : undefined;
    this.tangentRum = rum && jag
      ? {
          x: rum.x,
          y: rum.y,
          oppen: [0, 1, 2, 3].map((i) => rum.d[i] || (rum.exit === i && !rum.stangd)),
          vald: d === undefined ? -1 : d,
          farg: this.val.farger[jag.slot] || "#ffffff",
        }
      : null;
  }

  // flytta ser till att figuren id finns och går mot (x, z) i rummet
  // (rx, ry). Till ett rum bredvid går den genom dörren; längre bort, som
  // ur källarvalven, står den bara där.
  private flytta(id: string, nyckel: string, gor: () => Figur, rx: number, ry: number, x: number, z: number): Aktor {
    const mal = new THREE.Vector3(x, 0, z);
    let a = this.aktorer.get(id);
    if (a && a.nyckel !== nyckel) {
      const fig = gor();
      a.hall.remove(a.fig.grupp);
      a.hall.add(fig.grupp);
      a.fig = fig;
      a.nyckel = nyckel;
    }
    if (!a) {
      const hall = new THREE.Group();
      const fig = gor();
      hall.add(fig.grupp);
      hall.position.copy(mal);
      const vinkel = Math.PI / 4; // mot kameran
      hall.rotation.y = vinkel;
      this.figurer.add(hall);
      a = { hall, fig, nyckel, rum: [rx, ry], vag: [], vinkel, dod: false };
      this.aktorer.set(id, a);
      return a;
    }
    const [ox, oy] = a.rum;
    if (ox === rx && oy === ry) {
      if (a.hall.position.distanceTo(mal) > 0.01) a.vag = [mal];
    } else if (Math.abs(ox - rx) + Math.abs(oy - ry) === 1) {
      const dorr = new THREE.Vector3(((ox + rx + 1) / 2) * S, 0, ((oy + ry + 1) / 2) * S);
      a.vag = [dorr, mal];
    } else {
      a.hall.position.copy(mal);
      a.vag = [];
    }
    a.rum = [rx, ry];
    return a;
  }

  // egenRing ger den egna vandraren en ring i sin färg, eller tar bort den.
  private egenRing(a: Aktor, farg: string): void {
    if (a.ring && (!farg || (a.ring.material as THREE.MeshBasicMaterial).color.getStyle() !== new THREE.Color(farg).getStyle())) {
      a.hall.remove(a.ring);
      a.ring = undefined;
    }
    if (!farg || a.ring) return;
    a.ring = new THREE.Mesh(this.yta(1.05), this.kritMaterial(K.liten(), farg, "liten", 1));
    a.ring.rotation.x = -Math.PI / 2;
    a.ring.position.y = 0.02;
    a.ring.renderOrder = 1;
    a.hall.add(a.ring);
  }

  // tecknad är en varelse utan modell: dess tecken i krita, stående.
  private tecknad(glyph: string): Figur {
    const g = new THREE.Group();
    const sp = new THREE.Sprite(this.spriteMaterial(glyph, FARG.varelse));
    sp.center.set(0.5, 0.08);
    sp.scale.set(1.5, 1.5, 1);
    g.add(sp);
    const sk = new THREE.Mesh(this.yta(1.1), this.kritMaterial(K.skugga(), "#ffffff", "skugga", 1));
    sk.rotation.x = -Math.PI / 2;
    sk.position.y = 0.01;
    g.add(sk);
    return { grupp: g, steg: () => {} };
  }

  // figurerna går, vänder sig och lägger sig när de dör.
  private rorFigurer(t: number, dt: number): void {
    const p = 1 + 0.12 * Math.sin(t / 260);
    this.aktorer.forEach((a) => {
      let gar = false;
      let kvar = (FART * dt) / 1000;
      while (a.vag.length && kvar > 0) {
        const mal = a.vag[0];
        const d = mal.clone().sub(a.hall.position);
        const l = d.length();
        if (l > 0.001) a.vinkel = Math.atan2(d.x, d.z);
        if (l <= kvar) {
          a.hall.position.copy(mal);
          a.vag.shift();
          kvar -= l;
        } else {
          a.hall.position.addScaledVector(d, kvar / l);
          kvar = 0;
        }
        gar = true;
      }
      let dv = a.vinkel - a.hall.rotation.y;
      dv = Math.atan2(Math.sin(dv), Math.cos(dv));
      a.hall.rotation.y += dv * Math.min(1, dt / 90);
      // Den döda läggs ner på sin sockel.
      const modell = (a.fig.grupp.userData.modell as THREE.Object3D | undefined) || a.fig.grupp;
      modell.rotation.z = a.dod ? Math.PI / 2 : 0;
      modell.position.y = a.dod ? M.SOCKEL + 0.2 : M.SOCKEL;
      if (!a.dod) a.fig.steg(t, gar);
      if (a.ring) a.ring.scale.set(p, p, 1);
    });
  }

  private spriteMaterial(ch: string, farg: string): THREE.SpriteMaterial {
    const k = "sprite" + ch + farg;
    let m = this.material.get(k) as THREE.SpriteMaterial | undefined;
    if (!m) {
      m = new THREE.SpriteMaterial({ map: K.tecken(ch), color: farg, transparent: true });
      this.material.set(k, m);
    }
    return m;
  }

  // avslojande ger tornet nivåns utseende: väggarna andas eller slår,
  // golvet mörknar mot kött, revbenen syns och portarna blir munnar.
  private avslojande(t: number): void {
    const niva = this.niva();
    const look = LOOK[niva];
    const n = this.storleken || 9;
    const u = this.vaggar;
    u.uTid.value = t;
    u.uPeriod.value = look.period;
    u.uNiva.value = niva;
    u.uVagg.value.copy(srgb(look.vagg));
    u.uHi.value.copy(srgb(look.hi));
    u.uMitt.value.set((n * S) / 2, (n * S) / 2);

    // Golvet och revbenen slår i takt med väggarna, mätt vid kameran.
    const p = niva >= 3 && this.mitten ? puls(niva, t, Math.hypot(this.mitten.x - (n * S) / 2, this.mitten.z - (n * S) / 2)) : 0;
    const f = niva >= 3 ? 0.75 + 0.25 * p : 1;
    (this.golvMaterial() as THREE.MeshLambertMaterial).color.copy(srgb(look.golv)).multiplyScalar(f);
    this.revbenMat.opacity = look.revben * (niva >= 3 ? 0.7 + 0.3 * p : 1);
    this.revbenMat.visible = look.revben > 0;

    (this.stenMaterial("port", FARG.port) as THREE.MeshLambertMaterial).color.set(niva === 4 ? PORT_KOTT : FARG.port);
    const tugga = (Math.sin(t / 260) + 1) / 2;
    this.munnar.forEach(({ upper, lower }) => {
      upper.parent!.visible = niva === 4;
      upper.position.y = -tugga * 0.22;
      lower.position.y = tugga * 0.12;
    });
  }

  // ---- Varje bildruta ----

  // sankKvalitet slår av skuggorna och ritar med en pixel per punkt.
  private sankKvalitet(): void {
    if (this.lag) return;
    this.lag = true;
    this.sol.castShadow = false;
    this.renderer.setPixelRatio(1);
    this.anpassa();
    console.info("Kartan snett ovanifrån ritas med lägre kvalitet: inga skuggor och lägre upplösning.");
  }

  // matRuta räknar hur lång tid bildrutorna tar och sänker kvaliteten om
  // kartan går trögt.
  private matRuta(ram: number): void {
    if (this.lag || ram > 500) return; // en lång paus är en dold flik, inte trögt
    this.rutor.push(ram);
    if (this.rutor.length < MATT_RUTOR) return;
    const sorterade = this.rutor.sort((a, b) => a - b);
    const median = sorterade[Math.floor(sorterade.length / 2)];
    this.rutor = [];
    if (median > TROGT_MS) this.sankKvalitet();
  }

  private rita(t: number): void {
    if (!this.aktiv || !this.bredd || !this.hojd || this.canvas.offsetParent === null) return;
    if (t - this.senasteRuta < 1000 / MAX_BILDER - 2) return;
    this.matRuta(t - this.senasteRuta);
    this.senasteRuta = t;
    const n = this.storleken || 9;
    const u: Utsnitt = this.val.utsnitt() || { x: 0, y: 0, w: n, h: n };

    // Kameran: rummet kartan följer står i mitten, annars utsnittets mitt.
    // Utsnittet är en romb på skärmen, och när kartan följer någon får
    // hörnen skäras av, så att rummen närmast blir större. Mitten glider
    // efter när vandraren byter rum.
    const f = this.val.fokus();
    const mal = f
      ? new THREE.Vector3((f.x + 0.5) * S, 0.3, (f.y + 0.5) * S)
      : new THREE.Vector3((u.x + u.w / 2) * S, 0.3, (u.y + u.h / 2) * S);
    const dt = Math.min(200, t - this.senast);
    this.senast = t;
    if (!this.mitten) this.mitten = mal.clone();
    else this.mitten.lerp(mal, 1 - Math.exp(-dt / 160));
    const mitt = this.mitten;
    this.cam.position.copy(mitt).addScaledVector(FRAN, 150);
    this.cam.lookAt(mitt);
    const diag = ((u.w + u.h) * S) / Math.SQRT2;
    // Nära (3 rum) går lite närmare än de andra nivåerna, så att rummet man
    // står i fyller mer av en smal skärm.
    const skarp = f ? 0.6 + 0.12 * Math.min(1, Math.max(0, (Math.min(u.w, u.h) - 3) / 2)) : 1.03;
    let hw = (diag / 2) * skarp;
    let hh = (diag / 2) * LUTNING * skarp + 1.2;
    // Täcker loggen kartans vänstra del ställs utsnittet i den del som
    // syns: kartan ritas bredare och flyttas åt höger.
    const tackt = Math.max(0, Math.min(0.5, this.val.tackt ? this.val.tackt() : 0));
    this.tackt = tackt;
    const asp = (this.bredd * (1 - tackt)) / this.hojd;
    // En smal, hög karta (telefonen) går närmare när den följer någon:
    // rummet får mer av bredden, och rombens sidor skärs av.
    if (f && asp < 1) {
      const k = Math.max(0.75, Math.min(1, 0.45 + 0.55 * asp));
      hw *= k;
      hh *= k;
    }
    if (hw / hh > asp) hh = hw / asp;
    else hw = hh * asp;
    const hel = hw / (1 - tackt);
    Object.assign(this.cam, { left: -hel - hel * tackt, right: hel - hel * tackt, top: hh, bottom: -hh });
    this.cam.updateProjectionMatrix();

    // Solen följer kameran, så att skuggorna är skarpa där man tittar.
    this.sol.position.copy(mitt).add(SOL);
    this.sol.target.position.copy(mitt);
    const sk = this.sol.shadow.camera;
    const r = diag / 2 + 6;
    if (sk.right !== r) {
      Object.assign(sk, { left: -r, right: r, top: r, bottom: -r, near: 1, far: 80 });
      sk.updateProjectionMatrix();
    }

    this.rorFigurer(t, dt);
    this.landa(t);
    this.blanda(t);

    // Facklorna fladdrar, och rummet kartan följer ligger i facklans sken.
    const fladder = 0.82 + 0.1 * Math.sin(t / 97) + 0.08 * Math.sin(t / 41 + 1.3);
    this.flammaMaterial.color.setRGB(1, 0.62 * fladder + 0.1, 0.22 * fladder, THREE.SRGBColorSpace);
    this.skenMaterial.opacity = 0.26 * fladder;
    if (f) {
      this.fackla.position.set((f.x + 0.5) * S, 2.6, (f.y + 0.5) * S);
      this.fackla.intensity = 34 * fladder;
      this.dunkel.uFokus.value.set(mitt.x, mitt.z);
      this.dunkel.uRadie.value = S * 0.9;
      this.himmel.intensity = 0.8;
      this.sol.intensity = 1.25;
    } else {
      // Ingen att följa, till exempel i slutet: hela tornet syns, i ett
      // jämnare ljus.
      this.fackla.intensity = 0;
      this.dunkel.uRadie.value = 1e5;
      this.himmel.intensity = 1.6;
      this.sol.intensity = 1.9;
    }

    const vaken = t - this.vaknade < VAKEN_MS;
    const blixt = vaken && Math.floor((t - this.vaknade) / 180) % 2 === 0;
    this.ringar.forEach((m) => m.color.set(vaken ? (blixt ? "#ff5a3a" : "#a8402c") : RING[this.niva()]));
    this.ogon.forEach((o) => (o.visible = vaken && (blixt || t - this.vaknade > 1100)));
    if (this.guldljus) this.guldljus.intensity = vaken ? 30 * (0.7 + 0.3 * Math.sin(t / 60)) : 0;
    this.avslojande(t);

    [...this.ekon.children].forEach((o) => {
      const gar = t - o.userData.at;
      if (gar > EKO_MS) {
        this.ekon.remove(o);
        ((o as THREE.Mesh).material as THREE.Material).dispose();
      } else o.visible = Math.floor(gar / 250) % 2 === 0;
    });

    this.renderer.render(this.scene, this.cam);
    this.ritaKanten(!!f);
    this.ritaTangenter();
  }

  // ritaTangenter sätter en tangent med pilen vid varje öppen dörr och port
  // i rummet vandraren står i, så att man ser vart piltangenterna går när
  // norr inte är uppåt. Den valda tangenten lyser i vandrarens färg.
  private ritaTangenter(): void {
    const r = this.tangentRum;
    const v = new THREE.Vector3();
    const pos: [number, number][] = [];
    this.tangent.forEach((el, d) => {
      if (!r || !r.oppen[d]) {
        el.hidden = true;
        return;
      }
      const [x, z] = vidDorr(r.x, r.y, d);
      // Lite längre in i rummet än pilen, så att tangenten inte täcker dörren.
      const mx = r.x * S + S / 2, mz = r.y * S + S / 2;
      v.set(x + (mx - x) * 0.25, 0.1, z + (mz - z) * 0.25).project(this.cam);
      const tx = ((v.x + 1) / 2) * this.bredd, ty = ((1 - v.y) / 2) * this.hojd;
      pos.push([tx, ty]);
      el.hidden = false;
      el.style.left = tx + "px";
      el.style.top = ty + "px";
      const vald = r.vald === d;
      el.classList.toggle("vald", vald);
      el.style.color = vald ? r.farg : "";
      el.style.borderColor = vald ? r.farg : "";
    });
    // Är rummet litet på skärmen krymper tangenterna, så att de inte hamnar
    // över varandra (på telefonen, och när kartan är utzoomad).
    let narmast = Infinity;
    for (let a = 0; a < pos.length; a++) {
      for (let b = a + 1; b < pos.length; b++) {
        narmast = Math.min(narmast, Math.hypot(pos[a][0] - pos[b][0], pos[a][1] - pos[b][1]));
      }
    }
    const synlig = this.tangent.find((el) => !el.hidden);
    const storlek = synlig ? Math.max(synlig.offsetWidth, synlig.offsetHeight) + 4 : 0;
    const skala = storlek && narmast < storlek ? Math.max(0.5, narmast / storlek) : 1;
    const flytt = skala < 1 ? `translate(-50%, -50%) scale(${skala.toFixed(2)})` : "";
    this.tangent.forEach((el) => {
      if (el.style.transform !== flytt) el.style.transform = flytt;
    });
  }


  // ritaKanten sätter en liten pil vid kanten för det som ligger utanför
  // bilden, åt det håll det ligger, med samma tecken som i
  // teckenförklaringen: portarna och kammaren man känner till, och
  // vandrarna man ser. Bara när kartan följer någon.
  private ritaKanten(foljer: boolean): void {
    const mal: { x: number; z: number; tecken: string; farg: string }[] = [];
    const s = this.state;
    if (foljer && s) {
      s.rooms.forEach((r) => {
        const x = r.x * S + S / 2, z = r.y * S + S / 2;
        if (r.exit >= 0 && !r.stangd) mal.push({ x, z, tecken: "˄", farg: FARG.port });
        if (r.k === "kammare") mal.push({ x, z, tecken: "¤", farg: FARG.guld });
      });
      const f = this.val.fokus();
      s.players.forEach((p) => {
        if (p.status !== "inne" || p.hidden || p.x < 0 || p.id === s.you) return;
        if (f && p.x === f.x && p.y === f.y) return;
        mal.push({ x: p.x * S + S / 2, z: p.y * S + S / 2, tecken: "●", farg: this.val.farger[p.slot] || "#ffffff" });
      });
    }
    // Pilarna står vid kanten av den del av kartan som syns bredvid loggen.
    const W = this.bredd, H = this.hojd, m = 16, x0 = W * this.tackt, mx = (x0 + W) / 2, my = H / 2;
    const sedda: { nyckel: string; x: number; y: number }[] = [];
    const v = new THREE.Vector3();
    let n = 0;
    mal.forEach((q) => {
      v.set(q.x, 0.3, q.z).project(this.cam);
      const tx = ((v.x + 1) / 2) * W, ty = ((1 - v.y) / 2) * H;
      if (tx >= x0 && tx <= W && ty >= 0 && ty <= H) return;
      const dx = tx - mx, dy = ty - my;
      const k = Math.min(((W - x0) / 2 - m) / Math.abs(dx || 1e-6), (my - m) / Math.abs(dy || 1e-6));
      const ex = mx + dx * k, ey = my + dy * k;
      const nyckel = q.tecken + q.farg;
      if (sedda.some((o) => o.nyckel === nyckel && Math.hypot(o.x - ex, o.y - ey) < 22)) return;
      sedda.push({ nyckel, x: ex, y: ey });
      let el = this.kantPilar[n];
      if (!el) {
        el = document.createElement("div");
        el.innerHTML = '<span class="spets"></span><span class="tecken"></span>';
        this.kant.appendChild(el);
        this.kantPilar.push(el);
      }
      el.hidden = false;
      const a = Math.atan2(dy, dx);
      el.style.left = ex + "px";
      el.style.top = ey + "px";
      el.style.color = q.farg;
      const spets = el.firstChild as HTMLElement, tecken = el.lastChild as HTMLElement;
      spets.style.transform = "translate(-50%, -50%) rotate(" + a + "rad)";
      tecken.textContent = q.tecken;
      tecken.style.transform = "translate(calc(-50% + " + -Math.cos(a) * 14 + "px), calc(-50% + " + -Math.sin(a) * 14 + "px))";
      n++;
    });
    for (let i = n; i < this.kantPilar.length; i++) this.kantPilar[i].hidden = true;
  }
}
