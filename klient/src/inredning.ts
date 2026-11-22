// Inredningen i rummen på kartan snett ovanifrån: pelare, en brunn, bänkar,
// bokhyllor och annat efter rumstypen. Vilken rumstyp som har vad står
// under modell.inredning i data/texter.json, med namnen här nedan (samma
// lista som furnitureNames i data.go).
//
// Allt byggs i ett rum på S × S enheter, med (0, 0) i nordvästra hörnet,
// och läggs i listor efter material, så att kartan kan slå ihop det. Mitten
// av rummet och vandrarnas platser hålls fria: inredningen står i hörnen och
// längs murarna. Det höga står mot norra och västra muren, som vetter bort
// från kameran, så att det inte skymmer rummet.
import * as THREE from "three";

// Materialen inredningen läggs i. "sten", "tra", "jarn" och "lys" är
// kartans egna; andra nycklar är en färg som #rrggbb.
export type Lagg = (material: string, g: THREE.BufferGeometry) => void;

// Ett eget frö per rum, så att samma rum ser likadant ut varje gång.
function slumpFor(fro: number): () => number {
  let s = (fro * 2654435761) >>> 0 || 1;
  return () => {
    s ^= s << 13;
    s ^= s >>> 17;
    s ^= s << 5;
    return (s >>> 0) / 4294967296;
  };
}

interface Verktyg {
  lada(m: string, w: number, h: number, d: number, x: number, y: number, z: number, vrid?: number): void;
  cyl(m: string, r1: number, r2: number, h: number, x: number, y: number, z: number, sidor?: number): void;
  sten(m: string, r: number, x: number, y: number, z: number, sy?: number): void;
  kon(m: string, r: number, h: number, x: number, y: number, z: number, sidor?: number): void;
  geo(m: string, g: THREE.BufferGeometry, x: number, y: number, z: number): void;
  slump: () => number;
}

function verktyg(lagg: Lagg, ox: number, oz: number, fro: number): Verktyg {
  const flytta = (m: string, g: THREE.BufferGeometry, x: number, y: number, z: number) => {
    g.translate(ox + x, y, oz + z);
    lagg(m, g);
  };
  return {
    lada(m, w, h, d, x, y, z, vrid = 0) {
      const g = new THREE.BoxGeometry(w, h, d);
      if (vrid) g.rotateY(vrid);
      flytta(m, g, x, y, z);
    },
    cyl(m, r1, r2, h, x, y, z, sidor = 10) {
      flytta(m, new THREE.CylinderGeometry(r1, r2, h, sidor), x, y, z);
    },
    sten(m, r, x, y, z, sy = 0.75) {
      const g = new THREE.IcosahedronGeometry(r, 0);
      g.scale(1, sy, 1);
      flytta(m, g, x, y, z);
    },
    kon(m, r, h, x, y, z, sidor = 7) {
      flytta(m, new THREE.ConeGeometry(r, h, sidor), x, y, z);
    },
    geo(m, g, x, y, z) {
      flytta(m, g, x, y, z);
    },
    slump: slumpFor(fro),
  };
}

// Murarnas insidor ligger 0,2 från rummets kant (halva murens tjocklek).
const IN = 0.2;

const BOKFARGER = ["#7a2a22", "#2f4a6a", "#3f5a2e", "#6a5a2a", "#4a2a4a", "#8a6a3a"];
const KRUSFARGER = ["#9a6a44", "#7a7a70", "#6a4a34", "#a08a60"];

function pelare(v: Verktyg, x: number, z: number): void {
  v.lada("sten", 0.62, 0.18, 0.62, x, 0.09, z);
  v.cyl("sten", 0.22, 0.25, 1.5, x, 0.93, z, 8);
  v.lada("sten", 0.58, 0.16, 0.58, x, 1.76, z);
}

// En hylla längs muren: längs x mot norra muren (langsX) eller längs z mot
// västra. bredd är längs muren.
function hylla(v: Verktyg, x: number, z: number, bredd: number, langsX: boolean, bocker: boolean): void {
  const djup = 0.34, hojd = 1.35;
  const L = (w: number, h: number, d: number, px: number, py: number, pz: number, m = "tra") =>
    langsX ? v.lada(m, w, h, d, x + px, py, z + pz) : v.lada(m, d, h, w, x + pz, py, z + px);
  L(0.06, hojd, djup, -bredd / 2, hojd / 2, 0);
  L(0.06, hojd, djup, bredd / 2, hojd / 2, 0);
  L(bredd, 0.04, djup * 0.2, 0, hojd / 2, -djup / 2 + 0.02);
  const hyllor = [0.05, 0.47, 0.89, 1.31];
  hyllor.forEach((y) => L(bredd, 0.05, djup, 0, y, 0));
  hyllor.slice(0, 3).forEach((y) => {
    let px = -bredd / 2 + 0.07;
    while (px < bredd / 2 - 0.1) {
      if (bocker) {
        const w = 0.05 + v.slump() * 0.05, h = 0.24 + v.slump() * 0.1;
        if (v.slump() < 0.88) L(w, h, djup * 0.75, px + w / 2, y + 0.025 + h / 2, 0.02, BOKFARGER[Math.floor(v.slump() * BOKFARGER.length)]);
        px += w + 0.008;
      } else {
        const r = 0.06 + v.slump() * 0.05, h = 0.14 + v.slump() * 0.12;
        const f = KRUSFARGER[Math.floor(v.slump() * KRUSFARGER.length)];
        if (langsX) v.cyl(f, r * 0.8, r, h, x + px + r, y + 0.025 + h / 2, z, 7);
        else v.cyl(f, r * 0.8, r, h, x, y + 0.025 + h / 2, z + px + r, 7);
        px += r * 2 + 0.06 + v.slump() * 0.1;
      }
    }
  });
}

function bank(v: Verktyg, x: number, z: number, langd: number): void {
  v.lada("tra", 0.42, 0.08, langd, x, 0.44, z);
  [-1, 1].forEach((s) => v.lada("tra", 0.34, 0.4, 0.08, x, 0.2, z + s * (langd / 2 - 0.15)));
}

function lykta(v: Verktyg, x: number, z: number): void {
  v.cyl("jarn", 0.04, 0.05, 1.6, x, 0.8, z, 6);
  v.lada("jarn", 0.34, 0.12, 0.34, x, 0.06, z);
  v.lada("jarn", 0.4, 0.04, 0.04, x + 0.2, 1.58, z);
  v.lada("jarn", 0.22, 0.04, 0.22, x + 0.36, 1.42, z);
  v.lada("lys", 0.15, 0.2, 0.15, x + 0.36, 1.29, z);
  v.lada("jarn", 0.22, 0.04, 0.22, x + 0.36, 1.17, z);
}

// bygg lägger inredningen namn i rummet med hörnet (ox, oz).
export function bygg(namn: string, lagg: Lagg, ox: number, oz: number, S: number, fro: number): void {
  const v = verktyg(lagg, ox, oz, fro);
  const N = IN, O = S - IN; // norra och västra muren, södra och östra
  switch (namn) {
    case "pelare":
      pelare(v, 1.2, 1.35);
      pelare(v, S - 1.2, 1.35);
      break;
    case "pelare_syd":
      pelare(v, 1.2, S - 1.35);
      pelare(v, S - 1.2, S - 1.35);
      break;
    case "brunn": {
      const x = S - 1.35, z = 1.5;
      v.cyl("sten", 0.72, 0.78, 0.5, x, 0.25, z, 12);
      v.cyl("#0b1216", 0.56, 0.56, 0.02, x, 0.5, z, 12);
      [-1, 1].forEach((s) => v.lada("tra", 0.1, 1.15, 0.1, x + s * 0.7, 0.58, z));
      const stang = new THREE.CylinderGeometry(0.05, 0.05, 1.5, 6);
      stang.rotateZ(Math.PI / 2);
      v.geo("tra", stang, x, 1.12, z);
      v.cyl("jarn", 0.12, 0.1, 0.18, x, 0.82, z, 8);
      break;
    }
    case "bank":
      bank(v, N + 0.35, S - 1.5, 1.7);
      break;
    case "spjut": {
      const skaft = new THREE.CylinderGeometry(0.03, 0.03, 2.1, 6);
      skaft.rotateZ(Math.PI / 2);
      skaft.rotateY(0.5);
      v.geo("tra", skaft, S - 1.6, 0.05, 1.3);
      const spets = new THREE.ConeGeometry(0.06, 0.25, 6);
      spets.rotateZ(-Math.PI / 2);
      spets.rotateY(0.5);
      v.geo("jarn", spets, S - 1.6 + Math.cos(0.5) * 1.15, 0.05, 1.3 - Math.sin(0.5) * 1.15);
      break;
    }
    case "altare": {
      const x = S - 1.3, z = N + 0.55;
      v.lada("sten", 1.2, 0.8, 0.62, x, 0.4, z);
      v.lada("#6b1e1e", 1.24, 0.03, 0.4, x, 0.815, z);
      [-0.38, 0.38].forEach((d) => {
        v.cyl("#e8dcc0", 0.05, 0.05, 0.22, x + d, 0.94, z, 6);
        v.kon("lys", 0.035, 0.1, x + d, 1.1, z, 5);
      });
      break;
    }
    case "trappa": {
      const x = 1.35, z = 1.45;
      v.cyl("sten", 0.15, 0.15, 2.1, x, 1.05, z, 8);
      for (let i = 0; i < 10; i++) {
        const a = -Math.PI / 2 + i * 0.68;
        const g = new THREE.BoxGeometry(0.85, 0.1, 0.36);
        g.translate(0.5, 0, 0);
        g.rotateY(-a);
        v.geo("sten", g, x, 0.1 + i * 0.2, z);
      }
      break;
    }
    case "hylla":
      hylla(v, N + 0.19, 1.25, 1.25, false, false);
      break;
    case "tunnor": {
      const tunna = (x: number, y: number, z: number) => {
        v.cyl("tra", 0.27, 0.27, 0.62, x, y + 0.31, z, 12);
        [0.1, 0.52].forEach((dy) => v.cyl("jarn", 0.29, 0.29, 0.05, x, y + dy, z, 12));
      };
      tunna(S - 1.7, 0, N + 0.5);
      tunna(S - 1.05, 0, N + 0.5);
      tunna(S - 1.38, 0.62, N + 0.5);
      tunna(S - 0.75, 0, N + 1.15);
      break;
    }
    case "sackar":
      [[1.0, S - 1.0], [1.45, S - 0.85], [1.15, S - 1.45]].forEach(([x, z]) => v.sten("#8a7a55", 0.3, x, 0.2, z, 0.7));
      break;
    case "lyktor":
      lykta(v, S - 1.25, 1.2);
      lykta(v, 1.15, S - 1.3);
      break;
    case "stenblock":
      [[S - 1.2, 1.1, 0.42], [S - 0.75, 1.55, 0.3], [S - 1.55, 0.75, 0.26]].forEach(([x, z, r]) => v.sten("#6d665c", r, x, r * 0.6, z));
      [[1.0, S - 1.1, 0.95], [1.4, S - 0.8, 0.6], [0.8, S - 1.6, 0.7]].forEach(([x, z, h]) => v.kon("#7a7266", 0.2, h, x, h / 2, z, 6));
      break;
    case "assja": {
      const x = N + 0.55, z = 1.35;
      v.lada("sten", 0.85, 0.7, 1.1, x, 0.35, z);
      v.lada("#c8461a", 0.6, 0.04, 0.8, x, 0.71, z);
      v.lada("sten", 0.5, 0.9, 0.7, x - 0.1, 1.15, z);
      break;
    }
    case "stad": {
      const x = S - 1.35, z = 1.6;
      v.cyl("tra", 0.27, 0.3, 0.42, x, 0.21, z, 9);
      v.lada("jarn", 0.34, 0.12, 0.2, x, 0.48, z);
      v.lada("jarn", 0.6, 0.12, 0.26, x, 0.6, z);
      const horn = new THREE.ConeGeometry(0.1, 0.32, 6);
      horn.rotateZ(-Math.PI / 2);
      v.geo("jarn", horn, x + 0.45, 0.6, z);
      break;
    }
    case "bokhylla":
      hylla(v, S - 1.25, N + 0.19, 1.4, true, true);
      hylla(v, N + 0.19, 1.25, 1.2, false, true);
      break;
    case "bord": {
      const x = S - 1.45, z = O - 1.2;
      v.lada("tra", 1.1, 0.07, 0.75, x, 0.62, z);
      [[-1, -1], [-1, 1], [1, -1], [1, 1]].forEach(([a, b]) => v.lada("tra", 0.07, 0.6, 0.07, x + a * 0.47, 0.3, z + b * 0.3));
      [-1, 1].forEach((s) => {
        const sida = new THREE.BoxGeometry(0.26, 0.02, 0.34);
        sida.rotateZ(s * 0.12);
        v.geo("#e8dcc0", sida, x + s * 0.13, 0.67, z);
      });
      break;
    }
    case "stenhog": {
      const x = 1.15, z = S - 1.2;
      for (let i = 0; i < 7; i++) {
        const a = v.slump() * Math.PI * 2, d = v.slump() * 0.4;
        v.sten("#6d665c", 0.14 + v.slump() * 0.12, x + Math.cos(a) * d, 0.1 + (i > 4 ? 0.18 : 0), z + Math.sin(a) * d);
      }
      [0.3, 1.4, 2.6].forEach((a) => {
        const ben = new THREE.CylinderGeometry(0.035, 0.035, 0.5, 5);
        ben.rotateZ(Math.PI / 2);
        ben.rotateY(a);
        v.geo("#e8e0cc", ben, x + 0.1, 0.33, z);
      });
      break;
    }
  }
}
