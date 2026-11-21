// Kartan snett ovanifrån. Tornet byggs av golv, låga väggar och portar, och
// vandrarna är små figurer i sin färg. Kameran är ortografisk och vrids
// aldrig: norr pekar upp åt höger. Utsnittet, zoomen och teckenförklaringen
// räknas ut av den platta kartan (web/map.js), som den här kartan läser
// utsnittet från, så att de två alltid visar samma rum.
import * as THREE from "three";
import { mergeGeometries } from "three/examples/jsm/utils/BufferGeometryUtils.js";
import * as K from "./krita";
import * as M from "./modeller";
import type { Figur } from "./modeller";
import type { Rum, Tillstand, Utsnitt, Val } from "./typer";

const S = 6; // ett rum i världen
const VH = 0.7; // väggarnas höjd
const VT = 0.26; // väggarnas tjocklek
const LUCKA = 1.9; // dörrarnas bredd

// Den platta kartans celler, för vandrarnas platser i rummet.
const RW = 7;
const RH = 4;
const SLOTS = [[2, 1], [5, 1], [2, 3], [5, 3], [1, 2], [6, 2], [3, 1], [4, 3]];

const FARG = {
  okand: "#6a6e76",
  port: "#8d8a82",
  stangd: "#d0503a",
  guld: "#ffd36b",
  forlorad: "#a04a3c",
  mark: "#ff8a6a",
  varelse: "#ff7a5a",
  vald: "#ffe14a",
  krita: "#f2f0e8",
  borta: "#5a5d63",
  vagg: "#c9c6bc",
};
const RING = ["#c9a24f", "#c9a24f", "#c9a24f", "#c98a5a", "#c4483a"]; // per nivå

// Avslöjandet: tornets utseende följer nivån, som på den platta kartan. På
// nivå 2 andas kritan, på nivå 3 syns revbenen och hjärtslagen, på nivå 4
// är det kött och tänder. vagg och hi är väggens färg i vila och i slaget,
// golv golvets ton, period hur länge ett andetag eller hjärtslag tar och
// revben hur tydliga revbenen på golvet är.
const LOOK: Record<number, { vagg: number[]; hi: number[]; golv: number[]; period: number; revben: number }> = {
  1: { vagg: [201, 198, 188], hi: [201, 198, 188], golv: [255, 255, 255], period: 0, revben: 0 },
  2: { vagg: [196, 192, 182], hi: [240, 234, 220], golv: [255, 248, 236], period: 4200, revben: 0 },
  3: { vagg: [196, 160, 140], hi: [240, 170, 140], golv: [230, 170, 160], period: 1250, revben: 0.35 },
  4: { vagg: [176, 92, 74], hi: [255, 118, 88], golv: [210, 130, 120], period: 950, revben: 0.55 },
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
const FART = 3.4; // så fort figurerna går, i världens enheter per sekund
const VAKEN_MS = 3500;

// Kameran tittar från sydost, snett uppifrån.
const FRAN = new THREE.Vector3(30, 26, 30).normalize();
const LUTNING = FRAN.y; // hur mycket golvet trycks ihop på höjden
const SOL = new THREE.Vector3(-10, 22, -4);

// plats ger punkten i världen för cellen (cx, cy) i rummet (rx, ry), med
// den platta kartans celler: 1–6 på bredden och 1–3 på höjden.
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

export class IsoKarta {
  private renderer: THREE.WebGLRenderer;
  private scene = new THREE.Scene();
  private cam = new THREE.OrthographicCamera(-1, 1, 1, -1, 1, 400);
  private sol = new THREE.DirectionalLight(0xfff0dc, 2.2);
  private torn = new THREE.Group(); // golv, väggar och portar
  private rorligt = new THREE.Group(); // allt som ändras med varje tillstånd
  private ekon = new THREE.Group();
  private figurer = new THREE.Group();
  private aktorer = new Map<string, Aktor>();
  private tornNyckel = "";
  private state: Tillstand | null = null;
  private storleken = 0;
  private aktiv = false;
  private bredd = 0;
  private hojd = 0;
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
  private blinkar: THREE.Object3D[] = [];
  private vaknade = -Infinity;
  private mitten: THREE.Vector3 | null = null; // kamerans mitt, som glider efter
  private senast = 0;
  private material = new Map<string, THREE.Material>();
  private plan = new Map<number, THREE.PlaneGeometry>();
  private canvas: HTMLCanvasElement;
  private val: Val;
  private kompass: HTMLDivElement; // N och en pil i hörnet nere till vänster
  private kant: HTMLDivElement; // pilarna vid kanten för det som ligger utanför
  private kantPilar: HTMLDivElement[] = [];

  constructor(canvas: HTMLCanvasElement, val: Val) {
    this.canvas = canvas;
    this.val = val;
    this.renderer = new THREE.WebGLRenderer({ canvas, antialias: true });
    this.renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 2));
    this.renderer.setClearColor(0x15171a);
    this.renderer.shadowMap.enabled = true;
    this.renderer.shadowMap.type = THREE.PCFShadowMap;

    this.scene.add(new THREE.HemisphereLight(0xc9cfdc, 0x1a1c20, 1.1));
    this.sol.castShadow = true;
    this.sol.shadow.mapSize.set(2048, 2048);
    this.sol.shadow.bias = -0.0006;
    this.scene.add(this.sol, this.sol.target, this.torn, this.rorligt, this.ekon, this.figurer);

    this.kompass = document.createElement("div");
    this.kompass.className = "isokompass";
    this.kompass.setAttribute("aria-hidden", "true");
    this.kompass.hidden = true;
    this.kompass.innerHTML = kompassSvg();
    this.kant = document.createElement("div");
    this.kant.className = "isokant";
    this.kant.setAttribute("aria-hidden", "true");
    this.kant.hidden = true;
    canvas.after(this.kompass, this.kant);

    new ResizeObserver(() => this.anpassa()).observe(canvas.parentElement!);
    const loop = (t: number) => {
      this.rita(t);
      requestAnimationFrame(loop);
    };
    requestAnimationFrame(loop);
  }

  // visa slår på eller av kartan. Den platta kartan visas när den här är av.
  visa(pa: boolean): void {
    this.aktiv = pa;
    this.canvas.hidden = !pa;
    this.kompass.hidden = !pa;
    this.kant.hidden = !pa;
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
    const niva = this.niva();
    const slutet = niva === 4 && state.phase === "slut";
    const nyckel = this.storleken + "|" + (slutet ? 1 : 0) + "|" +
      state.rooms.map((r) => [r.x, r.y, r.d.map(Number).join(""), r.k, r.exit, r.stangd ? 1 : 0].join(",")).join(";");
    if (nyckel !== this.tornNyckel) {
      this.tornNyckel = nyckel;
      this.byggTorn(state, slutet);
    }
    this.byggRorligt(state);
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
    });

    const lada = (lista: THREE.BufferGeometry[], w: number, h: number, d: number, x: number, y: number, z: number) => {
      const g = new THREE.BoxGeometry(w, h, d);
      g.translate(x, y, z);
      lista.push(g);
    };
    // Ett väggstycke längs kanten från a till b (längs x om h, annars z).
    const stycke = (k: { h: boolean; x: number; z: number }, a: number, b: number) => {
      const len = b - a + VT;
      if (k.h) lada(vaggar, len, VH, VT, k.x + (a + b) / 2, VH / 2, k.z);
      else lada(vaggar, VT, VH, len, k.x, VH / 2, k.z + (a + b) / 2);
    };
    kanter.forEach((k) => {
      if (!k.oppen && !k.port) {
        stycke(k, 0, S);
        return;
      }
      stycke(k, 0, S / 2 - LUCKA / 2);
      stycke(k, S / 2 + LUCKA / 2, S);
      if (!k.port) return;
      // Porten: två stolpar och en överliggare, och en mur när den har slutit sig.
      const at = (u: number, y: number, w: number, h: number, d: number, lista = portar) =>
        k.h ? lada(lista, w, h, d, k.x + u, y, k.z) : lada(lista, d, h, w, k.x, y, k.z + u);
      at(S / 2 - LUCKA / 2 + VT * 0.2, 0.7, VT * 1.5, 1.4, VT * 1.5);
      at(S / 2 + LUCKA / 2 - VT * 0.2, 0.7, VT * 1.5, 1.4, VT * 1.5);
      at(S / 2, 1.5, LUCKA + VT * 1.2, 0.3, VT * 1.6);
      if (k.stangd) at(S / 2, 0.68, LUCKA, 1.36, VT, stangda);
      else this.mun(k);
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
    lagg(portar, this.enkelt("port", FARG.port), true);
    lagg(stangda, this.enkelt("stangd", FARG.stangd), true);
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

  private golvMaterial(): THREE.Material {
    let m = this.material.get("golv");
    if (!m) {
      m = new THREE.MeshLambertMaterial({ map: K.skiffer() });
      this.material.set("golv", m);
    }
    return m;
  }

  private vaggMaterial(): THREE.Material {
    let m = this.material.get("vagg");
    if (!m) {
      // Väggarnas färg räknas ut i shadern, så att hjärtslaget kan gå som en
      // våg från tornets mitt: samma puls som i puls(), för varje punkt.
      const mat = new THREE.MeshLambertMaterial({ map: K.vagg(), color: FARG.vagg });
      const u = this.vaggar;
      mat.onBeforeCompile = (sh) => {
        Object.assign(sh.uniforms, u);
        sh.vertexShader = sh.vertexShader
          .replace("#include <common>", "#include <common>\nvarying vec3 vVarld;")
          .replace("#include <begin_vertex>", "#include <begin_vertex>\nvVarld = (modelMatrix * vec4(transformed, 1.0)).xyz;");
        sh.fragmentShader = sh.fragmentShader
          .replace(
            "#include <common>",
            `#include <common>
varying vec3 vVarld;
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
      };
      m = mat;
      this.material.set("vagg", m);
    }
    return m;
  }

  private enkelt(nyckel: string, farg: string): THREE.Material {
    let m = this.material.get(nyckel);
    if (!m) {
      m = new THREE.MeshLambertMaterial({ color: farg });
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
        g.translate(cx + Math.cos(a) * d, 0.03 + lager * 0.08 + rnd() * 0.03, cz + Math.sin(a) * d);
        mynt.push(g);
      }
    }
    const hog = new THREE.Mesh(mergeGeometries(mynt), new THREE.MeshLambertMaterial({ color: 0xe8b84a, emissive: 0x2a1c00 }));
    mynt.forEach((g) => g.dispose());
    hog.userData.egen = true;
    hog.castShadow = true;
    this.torn.add(hog);
  }

  // ---- Det som rör sig ----

  private byggRorligt(state: Tillstand): void {
    this.rorligt.clear();
    this.blinkar = [];
    const g = this.rorligt;
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
      if (r.mark) {
        const [x, z] = plats(r.x, r.y, 3, 2);
        this.golvkrita(g, K.tecken(r.mark), "t" + r.mark, FARG.mark, 1.5, x, z);
      }
      if (r.f) {
        const [x, z] = plats(r.x, r.y, RW - 1, 1);
        this.golvkrita(g, K.tecken("+"), "t+", FARG.guld, 0.9, x, z);
      }
      if (r.c >= 0 && r.c !== undefined) {
        const [x, z] = plats(r.x, r.y, 1, 1);
        this.golvkrita(g, K.pil(), "pil", FARG.krita, 1.2, x, z, [Math.PI / 2, 0, -Math.PI / 2, Math.PI][r.c]);
      }
    });

    // Rummet vandraren står i, eller den spelledaren följer: ljusare golv
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
      const a = this.flytta("v" + p.id, klass + farg, () => M.vandrare(klass, farg), p.x, p.y, x, z);
      a.dod = dod;
      this.egenRing(a, p.id === state.you && !dod ? farg : "");
      sedda.add("v" + p.id);
    });
    (state.monsters || []).forEach((m) => {
      const [x, z] = plats(m.x, m.y, Math.floor(RW / 2) + 1, Math.floor(RH / 2));
      const typ = m.typ || "";
      this.flytta("m" + m.id, "m" + typ + m.glyph, () => M.varelse(typ) || this.tecknad(m.glyph), m.x, m.y, x, z);
      sedda.add("m" + m.id);
    });
    this.aktorer.forEach((a, id) => {
      if (sedda.has(id)) return;
      this.figurer.remove(a.hall);
      this.aktorer.delete(id);
    });

    // Den döda spelarens drag för en varelse, och det egna draget, är en pil
    // vid dörren som blinkar.
    const riktning: Record<string, number> = { N: 0, E: 1, S: 2, W: 3 };
    if (state.controller && state.action && state.action[0] === "m") {
      const id = parseInt(state.action.slice(1), 10);
      const d = riktning[state.action.slice(-1)];
      const m = (state.monsters || []).find((q) => q.id === id);
      if (m && d !== undefined) {
        const [x, z, v] = vidDorr(m.x, m.y, d);
        this.blinkar.push(this.golvkrita(g, K.pil(), "pil", FARG.varelse, 1.3, x, z, v, 0.02));
      }
    }
    const d = state.action ? riktning[state.action] : undefined;
    if (jag && jag.status === "inne" && d !== undefined) {
      const [x, z, v] = vidDorr(jag.x, jag.y, d);
      this.blinkar.push(this.golvkrita(g, K.pil(), "pil", this.val.farger[jag.slot] || "#ffffff", 1.3, x, z, v, 0.02));
    }
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
      a.fig.grupp.rotation.z = a.dod ? Math.PI / 2 : 0;
      a.fig.grupp.position.y = a.dod ? 0.2 : 0;
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

    (this.enkelt("port", FARG.port) as THREE.MeshLambertMaterial).color.set(niva === 4 ? PORT_KOTT : FARG.port);
    const tugga = (Math.sin(t / 260) + 1) / 2;
    this.munnar.forEach(({ upper, lower }) => {
      upper.parent!.visible = niva === 4;
      upper.position.y = -tugga * 0.22;
      lower.position.y = tugga * 0.12;
    });
  }

  // ---- Varje bildruta ----

  private rita(t: number): void {
    if (!this.aktiv || !this.bredd || !this.hojd || this.canvas.offsetParent === null) return;
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
    const skarp = f ? 0.72 : 1.03;
    let hw = (diag / 2) * skarp;
    let hh = (diag / 2) * LUTNING * skarp + 1.2;
    const asp = this.bredd / this.hojd;
    if (hw / hh > asp) hh = hw / asp;
    else hw = hh * asp;
    Object.assign(this.cam, { left: -hw, right: hw, top: hh, bottom: -hh });
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

    const blink = Math.floor(t / 550) % 2 === 0;
    this.blinkar.forEach((o) => (o.visible = blink));
    this.rorFigurer(t, dt);

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
  }

  // ritaKanten sätter en liten pil vid kanten för det som ligger utanför
  // bilden, åt det håll det ligger, med samma tecken som på den platta
  // kartan: portarna och kammaren man känner till, och vandrarna man ser.
  // Som på den platta kartan bara när kartan följer någon.
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
    const W = this.bredd, H = this.hojd, m = 16, mx = W / 2, my = H / 2;
    const sedda: { nyckel: string; x: number; y: number }[] = [];
    const v = new THREE.Vector3();
    let n = 0;
    mal.forEach((q) => {
      v.set(q.x, 0.3, q.z).project(this.cam);
      const tx = ((v.x + 1) / 2) * W, ty = ((1 - v.y) / 2) * H;
      if (tx >= 0 && tx <= W && ty >= 0 && ty <= H) return;
      const dx = tx - mx, dy = ty - my;
      const k = Math.min((mx - m) / Math.abs(dx || 1e-6), (my - m) / Math.abs(dy || 1e-6));
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
