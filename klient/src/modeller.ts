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

function lysande(farg: number): THREE.Material {
  const k = "lys" + farg;
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

// ---- Vandrarna ----

const HUD = 0xe4cfb2;
const STOVEL = 0x2e2a26;
const TRA = 0x6b4f2e;
const METALL = 0x9a9a92;
const PERGAMENT = 0xe9dcc0;

export function vandrare(klass: string, farg: string): Figur {
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
  kropp.add(del(geo("torso", () => new THREE.CylinderGeometry(0.17, 0.3, 0.64, 7)), tyg, 0, 0.76, 0));
  kropp.add(del(geo("huvud", () => new THREE.IcosahedronGeometry(0.16, 0)), lam(HUD), 0, 1.22, 0));

  // Armarna svänger kring axlarna. Den högra (vid +x) bär klassens föremål.
  const armar = [1, -1].map((s) => {
    const a = new THREE.Group();
    a.position.set(s * 0.21, 1.02, 0);
    a.rotation.z = s * 0.22;
    a.add(del(hangande("arm", 0.05, 0.045, 0.42, 5), tyg));
    a.add(del(geo("hand", () => new THREE.IcosahedronGeometry(0.055, 0)), lam(HUD), 0, -0.44, 0));
    kropp.add(a);
    return a;
  });
  const [bar, annan] = armar;
  let fram = 0; // hur långt fram armen med föremålet hålls

  // Huvudbonaden.
  const hatt = (brattet: number, farg: number) => {
    kropp.add(del(geo("brätte" + brattet, () => new THREE.CylinderGeometry(brattet, brattet, 0.03, 10)), lam(farg), 0, 1.33, 0));
    kropp.add(del(geo("kulle", () => new THREE.CylinderGeometry(0.13, 0.17, 0.18, 8)), lam(farg), 0, 1.43, 0));
  };
  const luva = (m: THREE.Material) => {
    const h = del(geo("huva", () => new THREE.ConeGeometry(0.2, 0.38, 7)), m, 0, 1.42, -0.04);
    h.rotation.x = -0.25;
    kropp.add(h);
  };
  // stav är en stav i handen, med något överst.
  const stav = (langd: number, topp?: THREE.Mesh) => {
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
      kropp.add(del(geo("hjälm", () => new THREE.SphereGeometry(0.19, 8, 5, 0, Math.PI * 2, 0, Math.PI / 2)), lam(METALL), 0, 1.24, 0));
      kropp.add(del(geo("näsa", () => new THREE.BoxGeometry(0.03, 0.12, 0.03)), lam(METALL), 0, 1.2, 0.17));
      const skold = del(geo("sköld", () => new THREE.CylinderGeometry(0.26, 0.26, 0.05, 8)), lam(0x7a5a32), -0.34, 0.82, 0.05);
      skold.rotation.z = Math.PI / 2;
      kropp.add(skold);
      kropp.add(del(geo("buckla", () => new THREE.IcosahedronGeometry(0.06, 0)), lam(METALL), -0.38, 0.82, 0.05));
      const svard = del(geo("svärd", () => new THREE.BoxGeometry(0.04, 0.5, 0.06)), lam(METALL), 0, -0.62, 0.08);
      bar.add(svard);
      break;
    }
    case "tjuv": {
      luva(lam(0x2a2a30));
      kropp.add(del(geo("säck", () => new THREE.IcosahedronGeometry(0.17, 0)), lam(0x6a5434), 0, 0.92, -0.28));
      const dolk = del(geo("dolk", () => new THREE.BoxGeometry(0.03, 0.22, 0.05)), lam(METALL), 0, -0.52, 0.06);
      bar.add(dolk);
      break;
    }
    case "ortkvinna":
      luva(huva);
      stav(1.35, del(geo("blomma", () => new THREE.IcosahedronGeometry(0.08, 0)), lam(0xe070b0)));
      kropp.add(del(geo("korg", () => new THREE.CylinderGeometry(0.13, 0.1, 0.14, 8)), lam(0x8a6a3a), -0.32, 0.62, 0.08));
      break;
    case "kartritare": {
      luva(huva);
      const rulle = del(geo("rulle", () => new THREE.CylinderGeometry(0.055, 0.055, 0.44, 8)), lam(PERGAMENT), 0, -0.46, 0.08);
      rulle.rotation.z = Math.PI / 2;
      bar.add(rulle);
      fram = -0.6;
      break;
    }
    case "ormtjusare": {
      luva(huva);
      const orm = del(geo("ormen", () => new THREE.TorusGeometry(0.18, 0.045, 5, 12)), lam(0x6d8040), 0, 1.04, 0);
      orm.rotation.x = Math.PI / 2;
      kropp.add(orm);
      const flojt = del(geo("flöjt", () => new THREE.CylinderGeometry(0.025, 0.025, 0.36, 6)), lam(TRA), 0, -0.48, 0.1);
      flojt.rotation.x = Math.PI / 2;
      bar.add(flojt);
      fram = -0.9;
      break;
    }
    case "sparfinnare":
      hatt(0.3, 0x6a4a2a);
      stav(1.0);
      break;
    case "lyktbarare": {
      luva(huva);
      const lykta = new THREE.Group();
      lykta.position.set(0, -0.6, 0.06);
      lykta.add(del(geo("lyktram", () => new THREE.BoxGeometry(0.15, 0.2, 0.15)), lam(0x3a3530)));
      lykta.add(new THREE.Mesh(geo("lågan", () => new THREE.BoxGeometry(0.11, 0.15, 0.16)), lysande(0xffd36b)));
      bar.add(lykta);
      fram = -0.25;
      break;
    }
    case "pilgrim":
      hatt(0.38, 0x3a3530);
      stav(1.6, del(geo("kalebass", () => new THREE.IcosahedronGeometry(0.09, 0)), lam(0x9a6a3a)));
      break;
    default:
      luva(huva);
  }

  return {
    grupp: g,
    steg(t, gar) {
      const s = gar ? Math.sin(t / 115) : 0;
      ben[0].rotation.x = s * 0.6;
      ben[1].rotation.x = -s * 0.6;
      annan.rotation.x = s * 0.5;
      bar.rotation.x = fram + (fram ? s * 0.15 : -s * 0.5);
      kropp.position.y = gar ? Math.abs(s) * 0.04 : Math.sin(t / 700) * 0.01;
    },
  };
}

// ---- Varelserna ----

export function varelse(typ: string): Figur | null {
  switch (typ) {
    case "krypare":
      return krypare();
    case "yngel":
      return yngel();
    case "fjallvakt":
      return fjallvakt();
    case "gallgast":
      return gallgast();
    case "benvaktare":
      return benvaktare();
  }
  return null;
}

// Blekkryparen: något blekt och mångbent, utan ögon.
function krypare(): Figur {
  const g = new THREE.Group();
  const blek = lam(0xd8d0bf);
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
function yngel(): Figur {
  const g = new THREE.Group();
  const fjall = lam(0x6d8040);
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
function fjallvakt(): Figur {
  const g = new THREE.Group();
  const fjall = lam(0x4f6e66), mork = lam(0x2f403c), metall = lam(METALL);
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
function gallgast(): Figur {
  const g = new THREE.Group();
  const kropp = new THREE.Group();
  g.add(kropp);
  const dimma = lam(0x8ac77f, { opacity: 0.72, emissive: 0x1a3a18 });
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
function benvaktare(): Figur {
  const g = new THREE.Group();
  const ben = lam(0xe8e0cc);
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
