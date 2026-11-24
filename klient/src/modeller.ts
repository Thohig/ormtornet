// Modellerna på kartan snett ovanifrån: vandrarna med det deras klass bär,
// och varelserna. Allt är enkla low-poly-former byggda i kod. Varje modell
// tittar mot +z och står på golvet (y = 0); kartan vrider och flyttar den.
//
// Klasserna och varelserna hittas med sitt id i data/klasser.json och
// data/monster.json. En klass utan modell får en vanlig huva, och en
// varelse utan modell visas med sitt tecken (se karta.ts).
import * as THREE from "three";

export interface Figur {
  grupp: THREE.Group;
  // steg rör modellen varje bildruta: t är tiden i millisekunder och gar
  // sant medan figuren går.
  steg(t: number, gar: boolean): void;
}

// ---- Delade material och geometrier ----

const material = new Map<string, THREE.Material>();

function lam(farg: THREE.ColorRepresentation, extra: { emissive?: number; opacity?: number } = {}): THREE.Material {
  const c = new THREE.Color(farg);
  const k = c.getHexString() + JSON.stringify(extra);
  let m = material.get(k);
  if (!m) {
    m = new THREE.MeshLambertMaterial({
      color: c,
      flatShading: true,
      emissive: extra.emissive ?? 0,
      transparent: extra.opacity !== undefined,
      opacity: extra.opacity ?? 1,
      side: extra.opacity !== undefined ? THREE.DoubleSide : THREE.FrontSide,
    });
    material.set(k, m);
  }
  return m;
}

function lysande(farg: THREE.ColorRepresentation): THREE.Material {
  const k = "lys" + new THREE.Color(farg).getHexString();
  let m = material.get(k);
  if (!m) {
    m = new THREE.MeshBasicMaterial({ color: farg });
    material.set(k, m);
  }
  return m;
}

function morkare(farg: string, f: number): THREE.Color {
  return new THREE.Color(farg).multiplyScalar(f);
}

const geometrier = new Map<string, THREE.BufferGeometry>();
function geo(k: string, gor: () => THREE.BufferGeometry): THREE.BufferGeometry {
  let g = geometrier.get(k);
  if (!g) {
    g = gor();
    geometrier.set(k, g);
  }
  return g;
}

// hangande ger en cylinder som hänger nedåt från sin övre ände, för ben och
// armar som svänger kring höften eller axeln.
function hangande(k: string, r1: number, r2: number, h: number, sidor = 6): THREE.BufferGeometry {
  return geo(k, () => {
    const g = new THREE.CylinderGeometry(r1, r2, h, sidor);
    g.translate(0, -h / 2, 0);
    return g;
  });
}

function del(g: THREE.BufferGeometry, m: THREE.Material, x = 0, y = 0, z = 0): THREE.Mesh {
  const o = new THREE.Mesh(g, m);
  o.position.set(x, y, z);
  o.castShadow = true;
  return o;
}

// ---- Socklarna ----

// miniatyr ställer figuren på en rund sockel av tenn, som en målad
// spelfigur: en mörk fot, en kant målad i färgen kant och en jordfärgad
// ovansida. Själva modellen ligger i grupp.userData.modell, så att kartan
// kan lägga ner den när vandraren dör medan sockeln står kvar. storlek är
// figurens storlek (modell.storlek), och en stor figur får en större sockel.
export function miniatyr(f: Figur, kant: THREE.ColorRepresentation, storlek = 1): Figur {
  const g = new THREE.Group();
  const fot = new THREE.Group();
  fot.add(del(geo("sockel", () => new THREE.CylinderGeometry(0.45, 0.5, 0.09, 24)), lam(0x3b3733), 0, 0.045, 0));
  fot.add(del(geo("sockelkant", () => new THREE.CylinderGeometry(0.505, 0.505, 0.035, 24)), lam(kant), 0, 0.03, 0));
  fot.add(del(geo("sockeltopp", () => new THREE.CylinderGeometry(0.44, 0.44, 0.012, 24)), lam(0x4a4034), 0, 0.093, 0));
  fot.scale.setScalar(Math.max(0.85, Math.min(1.6, storlek)));
  fot.children.forEach((o) => (o.receiveShadow = true));
  const modell = new THREE.Group();
  modell.add(f.grupp);
  modell.position.y = SOCKEL;
  g.add(fot, modell);
  g.userData.modell = modell;
  return { grupp: g, steg: f.steg };
}

// SOCKEL är sockelns höjd, där figuren står.
export const SOCKEL = 0.1;

// ---- Vandrarna ----

const HUD = 0xe4cfb2;
const STOVEL = 0x2e2a26;
const TRA = 0x6b4f2e;
const METALL = 0x9a9a92;
const PERGAMENT = 0xe9dcc0;

// Formen ger varje klass en egen kropp, så att klasserna går att skilja åt
// uppifrån: bålens bredd upptill och nedtill, var bålen slutar nedtill (0
// är en kåpa ända ner till golvet), hur långt fram vandraren lutar, hur
// lång kappan på ryggen är (0 är ingen) och hur långt ut armarna sitter.
interface Form {
  topp: number;
  botten: number;
  nere: number;
  lut: number;
  kappa: number;
  axel: number;
}

const FORM: Record<string, Partial<Form>> = {
  legoknekt: { topp: 0.24, botten: 0.33, kappa: 0.55, axel: 0.28 },
  tjuv: { topp: 0.14, botten: 0.22, nere: 0.5, lut: 0.14, kappa: 0.78, axel: 0.18 },
  ortkvinna: { topp: 0.16, botten: 0.37, nere: 0.1, kappa: 0.42 },
  kartritare: { kappa: 0 },
  ormtjusare: { kappa: 0.5 },
  sparfinnare: { lut: 0.08, kappa: 0.6 },
  lyktbarare: { kappa: 0.45 },
  pilgrim: { topp: 0.17, botten: 0.4, nere: 0.02, kappa: 0 },
};
const VANLIG: Form = { topp: 0.17, botten: 0.3, nere: 0.44, lut: 0, kappa: 0.5, axel: 0.21 };

const AXLAR = 1.08; // där bålen slutar upptill
const HUVUD = 1.29; // huvudets mitt

// vandrare bygger en vandrare i färgen farg. klassFarg är färgen på det
// klassen bär (modell.farg i klasser.json), annars klassens egen. Manteln
// över axlarna och kappan på ryggen har vandrarens färg, så att den syns
// uppifrån även när hatten har klassens.
export function vandrare(klass: string, farg: string, foremal: string[] = [], klassFarg?: string): Figur {
  const kf = (egen: THREE.ColorRepresentation): THREE.ColorRepresentation => klassFarg || egen;
  const form: Form = { ...VANLIG, ...FORM[klass] };
  const g = new THREE.Group();
  const tyg = lam(farg);
  const huva = lam(morkare(farg, 0.5));

  const ben = [-0.11, 0.11].map((x) => {
    const b = del(hangande("ben", 0.075, 0.07, 0.46), lam(STOVEL), x, 0.46, 0);
    g.add(b);
    return b;
  });
  const kropp = new THREE.Group();
  g.add(kropp);
  // Lutar vandraren framåt vrids kroppen kring höften, inte kring fötterna.
  kropp.rotation.x = form.lut;
  kropp.position.z = -Math.sin(form.lut) * 0.46;

  const hojd = AXLAR - form.nere;
  const k = `torso${form.topp}|${form.botten}|${form.nere}`;
  kropp.add(del(geo(k, () => new THREE.CylinderGeometry(form.topp, form.botten, hojd, 8)), tyg, 0, form.nere + hojd / 2, 0));
  kropp.add(del(geo("hals", () => new THREE.CylinderGeometry(0.07, 0.08, 0.1, 6)), lam(HUD), 0, AXLAR + 0.02, 0));
  kropp.add(del(geo("stortHuvud", () => new THREE.IcosahedronGeometry(0.2, 1)), lam(HUD), 0, HUVUD, 0));
  // Ögonen, så att man ser vart vandraren tittar.
  [-1, 1].forEach((s) => kropp.add(del(geo("vandraröga", () => new THREE.BoxGeometry(0.045, 0.05, 0.04)), lam(0x1a1410), s * 0.075, HUVUD - 0.02, 0.18)));

  // Manteln: en krage över axlarna och en kappa som hänger på ryggen.
  kropp.add(del(geo("mantel" + form.topp, () => new THREE.CylinderGeometry(0.15, form.topp + 0.16, 0.17, 9)), tyg, 0, AXLAR - 0.04, 0));
  let kappa: THREE.Mesh | null = null;
  if (form.kappa > 0) {
    const langd = form.kappa;
    const b = geo("kappa" + langd, () => {
      const l = new THREE.BoxGeometry(0.44, langd, 0.035);
      l.translate(0, -langd / 2, 0);
      return l;
    });
    kappa = del(b, tyg, 0, AXLAR - 0.04, -form.topp - 0.04);
    kappa.rotation.x = 0.3;
    kropp.add(kappa);
  }

  // Armarna svänger kring axlarna. Den högra (vid +x) bär klassens föremål.
  const armar = [1, -1].map((s) => {
    const a = new THREE.Group();
    a.position.set(s * form.axel, 1.02, 0);
    a.rotation.z = s * 0.22;
    a.add(del(hangande("arm", 0.05, 0.045, 0.42, 5), tyg));
    a.add(del(geo("hand", () => new THREE.IcosahedronGeometry(0.055, 0)), lam(HUD), 0, -0.44, 0));
    kropp.add(a);
    return a;
  });
  const [bar, annan] = armar;
  let fram = 0; // hur långt fram armen med föremålet hålls

  // Huvudbonaden.
  const hatt = (brattet: number, farg: THREE.ColorRepresentation) => {
    kropp.add(del(geo("brätte" + brattet, () => new THREE.CylinderGeometry(brattet, brattet, 0.03, 10)), lam(farg), 0, HUVUD + 0.15, 0));
    kropp.add(del(geo("kulle", () => new THREE.CylinderGeometry(0.15, 0.19, 0.2, 8)), lam(farg), 0, HUVUD + 0.26, 0));
  };
  const luva = (m: THREE.Material) => {
    const h = del(geo("storHuva", () => new THREE.ConeGeometry(0.25, 0.44, 8)), m, 0, HUVUD + 0.2, -0.05);
    h.rotation.x = -0.28;
    kropp.add(h);
  };
  // stav är en stav i handen, med något överst.
  const stav = (langd: number, topp?: THREE.Object3D) => {
    const s = del(geo("stav" + langd, () => new THREE.CylinderGeometry(0.025, 0.03, langd, 5)), lam(TRA), 0, -0.44 + langd * 0.4, 0.05);
    bar.add(s);
    if (topp) {
      topp.position.set(0, -0.44 + langd * 0.9, 0.05);
      bar.add(topp);
    }
    fram = -0.2;
  };

  switch (klass) {
    case "legoknekt": {
      // Bred, med hjälm och axelskydd av järn.
      const metall = lam(METALL);
      kropp.add(del(geo("storHjälm", () => new THREE.SphereGeometry(0.23, 9, 5, 0, Math.PI * 2, 0, Math.PI / 2)), metall, 0, HUVUD + 0.01, 0));
      kropp.add(del(geo("hjälmkant", () => new THREE.CylinderGeometry(0.245, 0.245, 0.04, 9)), metall, 0, HUVUD + 0.02, 0));
      kropp.add(del(geo("näsa", () => new THREE.BoxGeometry(0.035, 0.13, 0.03)), metall, 0, HUVUD - 0.03, 0.22));
      [-1, 1].forEach((s) => {
        const p = del(geo("axelskydd", () => new THREE.SphereGeometry(0.14, 8, 4, 0, Math.PI * 2, 0, Math.PI / 2)), metall, s * 0.27, 1.02, 0);
        p.rotation.z = -s * 0.35;
        kropp.add(p);
      });
      const skold = del(geo("sköld", () => new THREE.CylinderGeometry(0.28, 0.28, 0.05, 8)), lam(kf(0x7a5a32)), -0.42, 0.78, 0.06);
      skold.rotation.z = Math.PI / 2;
      kropp.add(skold);
      kropp.add(del(geo("buckla", () => new THREE.IcosahedronGeometry(0.06, 0)), metall, -0.46, 0.78, 0.06));
      bar.add(del(geo("svärd", () => new THREE.BoxGeometry(0.04, 0.55, 0.06)), metall, 0, -0.64, 0.08));
      break;
    }
    case "tjuv": {
      // Smal och framåtböjd, med en lång kappa och säcken på ryggen.
      luva(lam(kf(0x2a2a30)));
      kropp.add(del(geo("säck", () => new THREE.IcosahedronGeometry(0.17, 0)), lam(0x6a5434), 0.1, 0.78, -0.3));
      bar.add(del(geo("dolk", () => new THREE.BoxGeometry(0.03, 0.22, 0.05)), lam(METALL), 0, -0.52, 0.06));
      break;
    }
    case "ortkvinna":
      // Kjol ända ner till fötterna, staven med en blomma och korgen.
      luva(huva);
      stav(1.35, del(geo("blomma", () => new THREE.IcosahedronGeometry(0.08, 0)), lam(kf(0xe070b0))));
      kropp.add(del(geo("korg", () => new THREE.CylinderGeometry(0.14, 0.1, 0.15, 8)), lam(0x8a6a3a), -0.36, 0.62, 0.1));
      break;
    case "kartritare": {
      // Ett koger fullt av kartrullar på ryggen och en rulle i handen.
      luva(huva);
      const pergament = lam(kf(PERGAMENT));
      kropp.add(del(geo("koger", () => new THREE.CylinderGeometry(0.13, 0.11, 0.42, 8)), lam(0x7a5a34), 0, 0.86, -0.27));
      [[-0.06, 0.02], [0.05, 0.05], [0, -0.05]].forEach(([x, z], i) => {
        const r = del(geo("ryggrulle", () => new THREE.CylinderGeometry(0.045, 0.045, 0.5, 7)), pergament, x, 1.05 + i * 0.03, -0.27 + z);
        r.rotation.z = x * 2;
        kropp.add(r);
      });
      const rulle = del(geo("rulle", () => new THREE.CylinderGeometry(0.055, 0.055, 0.44, 8)), pergament, 0, -0.46, 0.08);
      rulle.rotation.z = Math.PI / 2;
      bar.add(rulle);
      fram = -0.6;
      break;
    }
    case "ormtjusare": {
      // En turban och ormen kring axlarna.
      const vit = lam(0xd9cdae);
      const turban = del(geo("turban", () => new THREE.TorusGeometry(0.17, 0.08, 6, 12)), vit, 0, HUVUD + 0.13, 0);
      turban.rotation.x = Math.PI / 2;
      kropp.add(turban);
      kropp.add(del(geo("turbantopp", () => new THREE.SphereGeometry(0.17, 8, 5, 0, Math.PI * 2, 0, Math.PI / 2)), vit, 0, HUVUD + 0.15, 0));
      const orm = del(geo("axelorm", () => new THREE.TorusGeometry(0.23, 0.05, 5, 14)), lam(kf(0x6d8040)), 0, AXLAR + 0.05, 0);
      orm.rotation.x = Math.PI / 2;
      kropp.add(orm);
      const flojt = del(geo("flöjt", () => new THREE.CylinderGeometry(0.025, 0.025, 0.36, 6)), lam(TRA), 0, -0.48, 0.1);
      flojt.rotation.x = Math.PI / 2;
      bar.add(flojt);
      fram = -0.9;
      break;
    }
    case "sparfinnare":
      // Bred hatt, framåtlutad med blicken i marken.
      hatt(0.32, kf(0x6a4a2a));
      stav(1.0);
      break;
    case "lyktbarare": {
      // Lyktan hänger från en hög stång, så att den syns uppifrån.
      luva(huva);
      stav(1.55);
      const lykta = new THREE.Group();
      lykta.position.set(0, 0.88, 0.2);
      lykta.add(del(geo("lyktkrok", () => new THREE.BoxGeometry(0.03, 0.03, 0.18)), lam(0x3a3530), 0, 0.14, -0.08));
      lykta.add(del(geo("lyktram", () => new THREE.BoxGeometry(0.17, 0.22, 0.17)), lam(0x3a3530)));
      lykta.add(new THREE.Mesh(geo("lågan", () => new THREE.BoxGeometry(0.13, 0.16, 0.18)), lysande(kf(0xffd36b))));
      bar.add(lykta);
      break;
    }
    case "pilgrim":
      // En lång kåpa ända ner till golvet, bred hatt och vandringsstav.
      hatt(0.4, kf(0x3a3530));
      stav(1.7, del(geo("kalebass", () => new THREE.IcosahedronGeometry(0.09, 0)), lam(0x9a6a3a)));
      break;
    default:
      luva(huva);
  }
  foremal.forEach((id) => fastPa(kropp, id, form.topp - VANLIG.topp));

  return {
    grupp: g,
    steg(t, gar) {
      const s = gar ? Math.sin(t / 115) : 0;
      ben[0].rotation.x = s * 0.6;
      ben[1].rotation.x = -s * 0.6;
      annan.rotation.x = s * 0.5;
      bar.rotation.x = fram + (fram ? s * 0.15 : -s * 0.5);
      kropp.position.y = gar ? Math.abs(s) * 0.04 : Math.sin(t / 700) * 0.01;
      if (kappa) kappa.rotation.x = 0.3 + (gar ? 0.15 + Math.abs(s) * 0.12 : Math.sin(t / 900) * 0.02);
    },
  };
}

// fastPa sätter ett föremål vandraren bär på kroppen: vid bältet, på bröstet,
// på ryggen eller på huvudet. Kroppen tittar mot +z; midjan är vid y 0,6.
// Föremål utan egen form syns inte.
//
// bredd är hur mycket bredare klassens bål är än den vanliga, så att det som
// sitter på bröstet eller ryggen hamnar utanpå.
function fastPa(kropp: THREE.Group, id: string, bredd = 0): void {
  const satt = (o: THREE.Object3D, x: number, y: number, z: number, rx = 0, ry = 0, rz = 0) => {
    o.position.set(x, y, z + Math.sign(z) * bredd);
    o.rotation.set(rx, ry, rz);
    kropp.add(o);
  };
  switch (id) {
    case "rep":
      satt(del(geo("repring", () => new THREE.TorusGeometry(0.1, 0.03, 4, 10)), lam(0x9a7a4a)), -0.27, 0.6, 0.06, 0, Math.PI / 2, 0);
      break;
    case "amulett":
      satt(del(geo("amulett", () => new THREE.OctahedronGeometry(0.06, 0)), lam(0x5fb0a0, { emissive: 0x0a2a24 })), 0, 0.9, 0.23);
      break;
    case "veke":
      satt(del(geo("veke", () => new THREE.CylinderGeometry(0.02, 0.02, 0.16, 5)), lam(0xf0ead8)), 0.16, 0.6, 0.25, 0, 0, 0.4);
      break;
    case "skold":
      satt(del(geo("ryggsköld", () => new THREE.CylinderGeometry(0.24, 0.24, 0.05, 8)), lam(0x6a5030)), 0, 0.86, -0.3, Math.PI / 2, 0, 0);
      break;
    case "vaska":
      satt(del(geo("väska", () => new THREE.BoxGeometry(0.12, 0.18, 0.2)), lam(0x7a5a34)), 0.3, 0.56, 0);
      break;
    case "kompass":
      satt(del(geo("kompass", () => new THREE.CylinderGeometry(0.06, 0.06, 0.025, 10)), lam(0xc9a24f)), 0, 0.62, 0.27, Math.PI / 2, 0, 0);
      break;
    case "rokpulver":
      satt(del(geo("pung", () => new THREE.IcosahedronGeometry(0.07, 0)), lam(0x8a8a8a)), -0.24, 0.55, 0.14);
      break;
    case "spegel":
      satt(del(geo("skärva", () => new THREE.TetrahedronGeometry(0.07, 0)), lam(0xd8e6f2, { emissive: 0x1a2630 })), 0.2, 0.66, 0.22);
      break;
    case "kastkniv":
      satt(del(geo("kniv", () => new THREE.BoxGeometry(0.03, 0.2, 0.04)), lam(METALL)), -0.12, 0.6, 0.27, 0, 0, 0.5);
      break;
    case "benlur":
      // Ett krökt horn av ben vid höften.
      satt(del(geo("benlur", () => new THREE.TorusGeometry(0.09, 0.025, 5, 10, Math.PI)), lam(0xe8e0cc)), -0.22, 0.62, 0.18, 0, 0, 0.6);
      break;
    case "flojt":
      satt(del(geo("flöjt", () => new THREE.CylinderGeometry(0.025, 0.025, 0.36, 6)), lam(TRA)), 0.08, 0.95, -0.27, 0, 0, 0.8);
      break;
    case "timglas": {
      const g = new THREE.Group();
      const ben = lam(0xe8e0cc);
      const kon = geo("timglaskon", () => new THREE.ConeGeometry(0.05, 0.07, 6));
      const ovre = del(kon, ben, 0, 0.035, 0);
      ovre.rotation.x = Math.PI;
      g.add(ovre, del(kon, ben, 0, -0.035, 0));
      satt(g, -0.28, 0.62, 0.06);
      break;
    }
    case "svart_parla":
      satt(del(geo("pärla", () => new THREE.SphereGeometry(0.055, 10, 8)), lam(0x15151a)), 0.1, 0.88, 0.23);
      break;
    case "mynt_hal":
      satt(del(geo("myntring", () => new THREE.TorusGeometry(0.05, 0.02, 4, 10)), lam(0xe8b84a)), -0.1, 0.88, 0.23);
      break;
    case "rostkrona": {
      const g = new THREE.Group();
      const rost = lam(0xa0522d);
      g.add(del(geo("kronring", () => new THREE.CylinderGeometry(0.15, 0.16, 0.08, 8, 1, true)), rost));
      for (let i = 0; i < 5; i++) {
        const a = (i / 5) * Math.PI * 2;
        g.add(del(geo("kronudd", () => new THREE.ConeGeometry(0.03, 0.08, 4)), rost, Math.sin(a) * 0.15, 0.08, Math.cos(a) * 0.15));
      }
      satt(g, 0, 1.46, 0);
      break;
    }
  }
}

// ---- Varelserna ----

// varelse bygger varelsen typ, i färgen farg (modell.farg i monster.json)
// eller sin egen.
export function varelse(typ: string, farg?: string): Figur | null {
  switch (typ) {
    case "krypare":
      return krypare(farg);
    case "yngel":
      return yngel(farg);
    case "fjallvakt":
      return fjallvakt(farg);
    case "gallgast":
      return gallgast(farg);
    case "benvaktare":
      return benvaktare(farg);
  }
  return null;
}

// Blekkryparen: något blekt och mångbent, utan ögon.
function krypare(farg?: string): Figur {
  const g = new THREE.Group();
  const blek = lam(farg || 0xd8d0bf);
  const kropp = del(geo("kryparkropp", () => new THREE.SphereGeometry(0.5, 10, 6)), blek, 0, 0.26, 0);
  kropp.scale.set(0.75, 0.38, 1);
  g.add(kropp);
  const huvud = del(geo("kryparhuvud", () => new THREE.SphereGeometry(0.2, 8, 6)), blek, 0, 0.24, 0.5);
  huvud.scale.set(1, 0.65, 1);
  g.add(huvud);
  const ben: { b: THREE.Mesh; i: number; s: number }[] = [];
  [-1, 1].forEach((s) => {
    for (let i = 0; i < 3; i++) {
      const b = del(hangande("kryparben", 0.035, 0.03, 0.56, 4), blek, s * 0.3, 0.32, -0.28 + i * 0.27);
      b.rotation.z = s * 0.95;
      g.add(b);
      ben.push({ b, i, s });
    }
  });
  return {
    grupp: g,
    steg(t, gar) {
      const fart = gar ? 50 : 90;
      ben.forEach(({ b, i, s }) => (b.rotation.x = Math.sin(t / fart + i * 1.6 + (s > 0 ? Math.PI : 0)) * 0.25));
    },
  };
}

// Ormynglet: långt som en man, rest ur dammet.
function yngel(farg?: string): Figur {
  const g = new THREE.Group();
  const fjall = lam(farg || 0x6d8040);
  const kurva = new THREE.CatmullRomCurve3(
    [[0, 0.09, -0.75], [0.38, 0.09, -0.38], [0, 0.1, 0], [-0.32, 0.12, 0.25], [-0.12, 0.45, 0.42], [0.04, 0.85, 0.32], [0.1, 1.06, 0.16]].map(
      ([x, y, z]) => new THREE.Vector3(x, y, z),
    ),
  );
  const kropp = new THREE.Group();
  g.add(kropp);
  kropp.add(del(geo("yngel", () => new THREE.TubeGeometry(kurva, 48, 0.1, 7)), fjall));
  const huvud = new THREE.Group();
  huvud.position.set(0.1, 1.1, 0.2);
  const h = del(geo("yngelhuvud", () => new THREE.SphereGeometry(0.16, 8, 6)), fjall);
  h.scale.set(1, 0.72, 1.45);
  huvud.add(h);
  [-1, 1].forEach((s) => huvud.add(del(geo("öga", () => new THREE.SphereGeometry(0.035, 6, 4)), lysande(0xffcc33), s * 0.09, 0.05, 0.14)));
  kropp.add(huvud);
  return {
    grupp: g,
    steg(t) {
      kropp.rotation.y = Math.sin(t / 650) * 0.25;
      huvud.position.y = 1.1 + Math.sin(t / 400) * 0.04;
    },
  };
}

// Fjällvakten: en gestalt i rustning av fjäll, med en sköld framför sig.
function fjallvakt(farg?: string): Figur {
  const g = new THREE.Group();
  const bas = new THREE.Color(farg || 0x62897f);
  const fjall = lam(bas), mork = lam(bas.clone().multiplyScalar(0.64)), metall = lam(METALL);
  const ben = [-0.14, 0.14].map((x) => {
    const b = del(hangande("vaktben", 0.1, 0.09, 0.55), mork, x, 0.55, 0);
    g.add(b);
    return b;
  });
  const kropp = new THREE.Group();
  g.add(kropp);
  kropp.add(del(geo("vakttorso", () => new THREE.CylinderGeometry(0.26, 0.34, 0.8, 8)), fjall, 0, 0.95, 0));
  // Fjällen: rader av små plattor på bröstet.
  for (let r = 0; r < 3; r++) {
    for (let i = 0; i < 4; i++) {
      const p = del(geo("fjäll", () => new THREE.ConeGeometry(0.07, 0.1, 4)), mork, -0.15 + i * 0.1 + (r % 2) * 0.05, 1.18 - r * 0.15, 0.27 - r * 0.012);
      p.rotation.x = Math.PI;
      kropp.add(p);
    }
  }
  kropp.add(del(geo("vakthjälm", () => new THREE.CylinderGeometry(0.17, 0.2, 0.3, 8)), metall, 0, 1.52, 0));
  kropp.add(del(geo("vaktkam", () => new THREE.BoxGeometry(0.04, 0.12, 0.34)), fjall, 0, 1.72, 0));
  [-1, 1].forEach((s) => kropp.add(new THREE.Mesh(geo("vaktöga", () => new THREE.BoxGeometry(0.06, 0.03, 0.02)), lysande(0xffcc44)).translateX(s * 0.07).translateY(1.53).translateZ(0.2)));
  const skold = del(geo("vaktsköld", () => new THREE.BoxGeometry(0.62, 0.86, 0.07)), mork, -0.05, 0.92, 0.42);
  kropp.add(skold);
  const spjut = del(geo("spjut", () => new THREE.CylinderGeometry(0.025, 0.025, 2.0, 5)), lam(TRA), 0.4, 1.0, 0.15);
  kropp.add(spjut);
  const udd = del(geo("udd", () => new THREE.ConeGeometry(0.06, 0.22, 5)), metall, 0.4, 2.1, 0.15);
  kropp.add(udd);
  return {
    grupp: g,
    steg(t, gar) {
      const s = gar ? Math.sin(t / 140) : 0;
      ben[0].rotation.x = s * 0.45;
      ben[1].rotation.x = -s * 0.45;
      kropp.position.y = gar ? Math.abs(s) * 0.03 : 0;
      skold.position.z = 0.42 + Math.sin(t / 900) * 0.02;
    },
  };
}

// Gallgasten: en grönaktig skepnad i trasor, som glider fram.
function gallgast(farg?: string): Figur {
  const g = new THREE.Group();
  const kropp = new THREE.Group();
  g.add(kropp);
  const dimma = lam(farg || 0x8ac77f, { opacity: 0.72, emissive: 0x1a3a18 });
  kropp.add(del(geo("gastkåpa", () => new THREE.ConeGeometry(0.42, 1.25, 9, 1, true)), dimma, 0, 0.72, 0));
  kropp.add(del(geo("gasthuvud", () => new THREE.IcosahedronGeometry(0.18, 1)), dimma, 0, 1.45, 0));
  [-1, 1].forEach((s) => kropp.add(new THREE.Mesh(geo("gastöga", () => new THREE.SphereGeometry(0.035, 6, 4)), lysande(0xd8ffb0)).translateX(s * 0.07).translateY(1.47).translateZ(0.15)));
  const armar = [-1, 1].map((s) => {
    const a = del(hangande("gastarm", 0.05, 0.03, 0.55, 5), dimma, s * 0.24, 1.15, 0.05);
    a.rotation.x = -0.9;
    kropp.add(a);
    return a;
  });
  // Trasor som hänger från kåpans kant.
  for (let i = 0; i < 7; i++) {
    const a = (i / 7) * Math.PI * 2;
    const t = del(geo("trasa", () => new THREE.PlaneGeometry(0.16, 0.3)), dimma, Math.sin(a) * 0.4, 0.02, Math.cos(a) * 0.4);
    t.rotation.y = a;
    kropp.add(t);
  }
  return {
    grupp: g,
    steg(t) {
      kropp.position.y = 0.12 + Math.sin(t / 500) * 0.07;
      armar.forEach((a, i) => (a.rotation.x = -0.9 + Math.sin(t / 420 + i) * 0.12));
    },
  };
}

// Benväktaren: byggd av böjda, vita ben, som revben från samma djur.
function benvaktare(farg?: string): Figur {
  const g = new THREE.Group();
  const ben = lam(farg || 0xe8e0cc);
  const kropp = new THREE.Group();
  g.add(kropp);
  for (let i = 0; i < 8; i++) kropp.add(del(geo("kota", () => new THREE.IcosahedronGeometry(0.065, 0)), ben, 0, 0.55 + i * 0.12, -0.08));
  for (let i = 0; i < 4; i++) {
    const r = 0.26 - i * 0.02;
    const rev = del(geo("revben" + i, () => new THREE.TorusGeometry(r, 0.025, 4, 12, Math.PI * 1.3)), ben, 0, 0.85 + i * 0.13, -0.08);
    rev.rotation.set(Math.PI / 2, 0, Math.PI * 0.5 - Math.PI * 0.65);
    kropp.add(rev);
  }
  const skalle = del(geo("skalle", () => new THREE.IcosahedronGeometry(0.17, 1)), ben, 0, 1.58, 0);
  skalle.scale.set(1, 1.05, 1.2);
  kropp.add(skalle);
  kropp.add(del(geo("käke", () => new THREE.BoxGeometry(0.18, 0.06, 0.16)), ben, 0, 1.44, 0.05));
  [-1, 1].forEach((s) => kropp.add(new THREE.Mesh(geo("benöga", () => new THREE.SphereGeometry(0.035, 6, 4)), lysande(0xff7a5a)).translateX(s * 0.07).translateY(1.6).translateZ(0.17)));
  const armar = [-1, 1].map((s) => {
    const a = del(hangande("benarm", 0.035, 0.025, 0.75, 5), ben, s * 0.3, 1.28, 0);
    a.rotation.z = s * 0.25;
    kropp.add(a);
    return a;
  });
  const lar = [-0.12, 0.12].map((x) => {
    const b = del(hangande("benlår", 0.04, 0.035, 0.56, 5), ben, x, 0.56, 0);
    g.add(b);
    return b;
  });
  return {
    grupp: g,
    steg(t, gar) {
      const s = gar ? Math.sin(t / 130) : 0;
      lar[0].rotation.x = s * 0.5;
      lar[1].rotation.x = -s * 0.5;
      kropp.rotation.z = Math.sin(t / 800) * 0.04;
      armar.forEach((a, i) => (a.rotation.x = -0.3 + Math.sin(t / 600 + i * 2) * 0.15));
    },
  };
}

// ---- Fällorna ----

// falla bygger fällan id (id i data/fallor.json) som den ser ut när den
// har slagit till, med mitten i origo. En fälla utan modell visas med sitt
// tecken i krita (se karta.ts).
export function falla(id: string): THREE.Group | null {
  const g = new THREE.Group();
  const jarn = lam(0x4a4744);
  const hal = lam(0x060504);
  const sten = lam(0x6e665a);
  switch (id) {
    case "fallgaller": {
      // Ett galler av spetsiga järnstänger som har rasat ner.
      [-1, 1].forEach((s) => g.add(del(geo("gallerstolpe", () => new THREE.BoxGeometry(0.09, 1.15, 0.09)), jarn, s * 0.62, 0.575, 0)));
      g.add(del(geo("galleroverstang", () => new THREE.BoxGeometry(1.33, 0.09, 0.09)), jarn, 0, 1.12, 0));
      g.add(del(geo("gallertvarstang", () => new THREE.BoxGeometry(1.25, 0.05, 0.05)), jarn, 0, 0.55, 0));
      for (let i = 0; i < 6; i++) {
        const x = -0.45 + i * 0.18;
        g.add(del(geo("gallerstang", () => new THREE.CylinderGeometry(0.025, 0.025, 1.0, 5)), jarn, x, 0.6, 0));
        const spets = del(geo("gallerspets", () => new THREE.ConeGeometry(0.04, 0.12, 5)), jarn, x, 0.06, 0);
        spets.rotation.x = Math.PI;
        g.add(spets);
      }
      break;
    }
    case "rasande_golv": {
      // Ett hål i golvet, med lösa plattor runt kanten.
      g.add(del(geo("hal", () => new THREE.CylinderGeometry(0.62, 0.62, 0.02, 18)), hal, 0, 0.012, 0));
      for (let i = 0; i < 6; i++) {
        const a = i * 1.05 + 0.3;
        const p = del(geo("platta", () => new THREE.BoxGeometry(0.34, 0.07, 0.26)), sten, Math.cos(a) * 0.72, 0.05, Math.sin(a) * 0.72);
        p.rotation.set(0.25 * Math.sin(i * 2.1), a, 0.3 * Math.cos(i * 1.7));
        g.add(p);
      }
      break;
    }
    case "roterande": {
      // En stor kvarnsten i golvet som har vridit sig.
      g.add(del(geo("kvarnsten", () => new THREE.CylinderGeometry(0.8, 0.82, 0.08, 22)), lam(0x7d756a), 0, 0.04, 0));
      g.add(del(geo("kvarnnav", () => new THREE.CylinderGeometry(0.13, 0.13, 0.14, 10)), jarn, 0, 0.08, 0));
      for (let i = 0; i < 6; i++) {
        const r = del(geo("kvarnranna", () => new THREE.BoxGeometry(0.62, 0.02, 0.05)), hal, 0, 0.085, 0);
        r.rotation.y = (i * Math.PI) / 3;
        r.translateX(0.42);
        g.add(r);
      }
      break;
    }
    case "pilar": {
      // Benvita pilar som sitter i golvet.
      const ben = lam(0xe8e0cc);
      [[-0.35, -0.2, 0.3, 0.2], [0.25, -0.35, -0.4, 0.1], [0.1, 0.3, 0.2, -0.3], [-0.4, 0.35, -0.2, -0.2], [0.45, 0.15, 0.35, 0.25]].forEach(([x, z, rx, rz]) => {
        const pil = new THREE.Group();
        pil.add(del(geo("pilskaft", () => new THREE.CylinderGeometry(0.018, 0.018, 0.55, 5)), ben, 0, 0.27, 0));
        pil.add(del(geo("pilfjader", () => new THREE.BoxGeometry(0.1, 0.12, 0.01)), lam(0x8a2a22), 0, 0.5, 0));
        pil.position.set(x, 0, z);
        pil.rotation.set(rx, 0, rz);
        g.add(pil);
      });
      break;
    }
    case "dimma": {
      // Gulgrön dimma som ligger kvar över golvet.
      const dimma = lam(0xa8cf5a, { emissive: 0x1e2a08, opacity: 0.38 });
      [[0, 0.22, 0, 0.55], [0.4, 0.18, 0.25, 0.4], [-0.35, 0.16, -0.2, 0.42], [0.1, 0.32, -0.35, 0.33]].forEach(([x, y, z, r]) => {
        const m = new THREE.Mesh(geo("dimboll", () => new THREE.IcosahedronGeometry(1, 1)), dimma);
        m.scale.set(r, r * 0.45, r);
        m.position.set(x, y, z);
        g.add(m);
      });
      break;
    }
    case "ormgrop": {
      // En grop full av små ormar.
      g.add(del(geo("hal", () => new THREE.CylinderGeometry(0.62, 0.62, 0.02, 18)), hal, 0, 0.012, 0));
      const orm = lam(0x5a7a34);
      for (let i = 0; i < 5; i++) {
        const o = del(geo("smaorm", () => new THREE.TorusGeometry(0.13, 0.028, 4, 10, Math.PI * 1.5)), orm, Math.cos(i * 1.3) * 0.35, 0.04, Math.sin(i * 1.3) * 0.32);
        o.rotation.set(Math.PI / 2, 0, i * 1.1);
        g.add(o);
      }
      break;
    }
    case "skinn": {
      // Ett ömsat ormskinn, blekt och genomskinligt.
      const skinn = lam(0xd8d2b8, { opacity: 0.7 });
      const s = del(geo("ormskinn", () => new THREE.TorusGeometry(0.45, 0.1, 6, 18, Math.PI * 1.6)), skinn, 0, 0.06, 0);
      s.rotation.x = Math.PI / 2;
      s.scale.set(1, 1, 0.5);
      g.add(s);
      break;
    }
    case "klocka": {
      // En klocka i en ställning av trä.
      const tra = lam(TRA);
      [-1, 1].forEach((s) => g.add(del(geo("klockstolpe", () => new THREE.BoxGeometry(0.09, 1.45, 0.09)), tra, s * 0.5, 0.72, 0)));
      g.add(del(geo("klockbjalke", () => new THREE.BoxGeometry(1.15, 0.1, 0.12)), tra, 0, 1.45, 0));
      g.add(del(geo("klocka", () => new THREE.CylinderGeometry(0.12, 0.34, 0.45, 14)), lam(0xb0803a), 0, 1.15, 0));
      g.add(del(geo("kolv", () => new THREE.SphereGeometry(0.07, 8, 6)), jarn, 0, 0.9, 0));
      break;
    }
    case "morker": {
      // Svart rök som har släckt ljuset.
      const rok = lam(0x0a0908, { opacity: 0.72 });
      [[0, 0.3, 0, 0.45], [0.3, 0.55, 0.15, 0.32], [-0.25, 0.5, -0.1, 0.35], [0.05, 0.8, -0.05, 0.25]].forEach(([x, y, z, r]) => {
        const m = new THREE.Mesh(geo("rokboll", () => new THREE.IcosahedronGeometry(1, 1)), rok);
        m.scale.setScalar(r);
        m.position.set(x, y, z);
        g.add(m);
      });
      break;
    }
    case "eko": {
      // Ringar i golvet, som ljudet som sprider sig.
      const ring = lam(0xc8c0b0, { opacity: 0.6 });
      [0.3, 0.55, 0.8].forEach((r, i) => {
        const m = new THREE.Mesh(geo("ekoring" + i, () => new THREE.TorusGeometry(r, 0.025, 4, 28)), ring);
        m.rotation.x = Math.PI / 2;
        m.position.y = 0.03;
        g.add(m);
      });
      break;
    }
    default:
      return null;
  }
  return g;
}
