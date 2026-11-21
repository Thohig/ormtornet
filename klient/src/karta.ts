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
  private blinkar: THREE.Object3D[] = [];
  private vaknade = -Infinity;
  private mitten: THREE.Vector3 | null = null; // kamerans mitt, som glider efter
  private senast = 0;
  private material = new Map<string, THREE.Material>();
  private plan = new Map<number, THREE.PlaneGeometry>();
  private canvas: HTMLCanvasElement;
  private val: Val;

  constructor(canvas: HTMLCanvasElement, val: Val) {
    this.canvas = canvas;
    this.val = val;
    this.renderer = new THREE.WebGLRenderer({ canvas, antialias: true });
    this.renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 2));
    this.renderer.setClearColor(0x15171a);
    this.renderer.shadowMap.enabled = true;
    this.renderer.shadowMap.type = THREE.PCFSoftShadowMap;

    this.scene.add(new THREE.HemisphereLight(0xc9cfdc, 0x1a1c20, 1.1));
    this.sol.castShadow = true;
    this.sol.shadow.mapSize.set(2048, 2048);
    this.sol.shadow.bias = -0.0006;
    this.scene.add(this.sol, this.sol.target, this.torn, this.rorligt, this.ekon, this.figurer);

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
      m = new THREE.MeshLambertMaterial({ map: K.vagg(), color: FARG.vagg });
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

  // Kammaren: en kritring kring en hög av guld.
  private byggKammare(r: Rum): void {
    const cx = r.x * S + S / 2, cz = r.y * S + S / 2;
    const ringMat = new THREE.MeshBasicMaterial({ map: K.ring(), color: RING[this.niva()], transparent: true, depthWrite: false });
    const ring = new THREE.Mesh(this.yta(S * 0.75), ringMat);
    ring.rotation.x = -Math.PI / 2;
    ring.position.set(cx, 0.012, cz);
    this.torn.add(ring);
    this.ringar.push(ringMat);

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

    [...this.ekon.children].forEach((o) => {
      const gar = t - o.userData.at;
      if (gar > EKO_MS) {
        this.ekon.remove(o);
        ((o as THREE.Mesh).material as THREE.Material).dispose();
      } else o.visible = Math.floor(gar / 250) % 2 === 0;
    });

    this.renderer.render(this.scene, this.cam);
  }
}
