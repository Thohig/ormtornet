// Texturer ritade på canvas: tecknen i krita på golvet, som ritas vita och
// färgas av materialet, och golvets sten, murverket och träet.
import * as THREE from "three";

// Ett eget frö, så att kritan ser likadan ut varje gång.
let seed = 7;
export function slump(): number {
  seed = (seed + 0x6d2b79f5) | 0;
  let t = Math.imul(seed ^ (seed >>> 15), 1 | seed);
  t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t;
  return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
}

function duk(w: number, h: number): [HTMLCanvasElement, CanvasRenderingContext2D] {
  const c = document.createElement("canvas");
  c.width = w;
  c.height = h;
  return [c, c.getContext("2d")!];
}

// grova låter kritan inte ta överallt: en del punkter försvinner och resten
// blir lite olika starka.
function grova(c: HTMLCanvasElement, kvar = 0.9): HTMLCanvasElement {
  const ctx = c.getContext("2d")!;
  const im = ctx.getImageData(0, 0, c.width, c.height);
  const d = im.data;
  for (let i = 3; i < d.length; i += 4) {
    if (!d[i]) continue;
    d[i] = slump() > kvar ? 0 : d[i] * (0.55 + 0.45 * slump());
  }
  ctx.putImageData(im, 0, 0);
  return c;
}

function textur(c: HTMLCanvasElement): THREE.CanvasTexture {
  const t = new THREE.CanvasTexture(c);
  t.colorSpace = THREE.SRGBColorSpace;
  t.anisotropy = 4;
  return t;
}

type Punkt = [number, number];

function streck(ctx: CanvasRenderingContext2D, p: Punkt[], w: number): void {
  ctx.strokeStyle = "#fff";
  ctx.lineCap = "round";
  ctx.lineJoin = "round";
  for (let k = 0; k < 3; k++) {
    ctx.globalAlpha = k ? 0.5 : 0.95;
    ctx.lineWidth = w * (k ? 0.7 : 1);
    ctx.beginPath();
    p.forEach(([x, y], i) => {
      const jx = x + (slump() - 0.5) * w * 0.45;
      const jy = y + (slump() - 0.5) * w * 0.45;
      if (i) ctx.lineTo(jx, jy);
      else ctx.moveTo(jx, jy);
    });
    ctx.stroke();
  }
  ctx.globalAlpha = 1;
}

function ellips(cx: number, cy: number, rx: number, ry: number, n = 40): Punkt[] {
  const p: Punkt[] = [];
  for (let i = 0; i <= n; i++) {
    const a = (i / n) * Math.PI * 2;
    p.push([cx + Math.cos(a) * rx, cy + Math.sin(a) * ry]);
  }
  return p;
}

const sparade = new Map<string, THREE.Texture>();
function spara(nyckel: string, gor: () => THREE.Texture): THREE.Texture {
  let t = sparade.get(nyckel);
  if (!t) {
    t = gor();
    sparade.set(nyckel, t);
  }
  return t;
}

// tecken är ett tecken i krita, till exempel fällans märke eller en varelse.
export function tecken(ch: string): THREE.Texture {
  return spara("t" + ch, () => {
    const [c, x] = duk(128, 128);
    x.fillStyle = "#fff";
    x.font = 'bold 92px Menlo, Monaco, "DejaVu Sans Mono", "Segoe UI Symbol", monospace';
    x.textAlign = "center";
    x.textBaseline = "middle";
    x.fillText(ch, 64, 68);
    return textur(grova(c, 0.88));
  });
}

// pil pekar åt höger, mot öster, och vrids dit den ska.
export function pil(): THREE.Texture {
  return spara("pil", () => {
    const [c, x] = duk(128, 128);
    streck(x, [[16, 64], [108, 64]], 10);
    streck(x, [[78, 38], [110, 64], [78, 90]], 10);
    return textur(grova(c));
  });
}

export function ring(): THREE.Texture {
  return spara("ring", () => {
    const [c, x] = duk(256, 256);
    streck(x, ellips(128, 128, 116, 116, 60), 9);
    streck(x, ellips(128, 128, 102, 102, 60), 3);
    return textur(grova(c, 0.88));
  });
}

export function liten(): THREE.Texture {
  return spara("liten", () => {
    const [c, x] = duk(128, 128);
    streck(x, ellips(64, 64, 52, 52, 40), 9);
    return textur(grova(c));
  });
}

export function prick(): THREE.Texture {
  return spara("prick", () => {
    const [c, x] = duk(64, 64);
    x.fillStyle = "#fff";
    x.beginPath();
    x.arc(32, 32, 9, 0, Math.PI * 2);
    x.fill();
    return textur(grova(c, 0.85));
  });
}

// ram är en tunn ram strax innanför väggarna, för rummet man står i.
export function ram(): THREE.Texture {
  return spara("ram", () => {
    const [c, x] = duk(256, 256);
    streck(x, [[18, 18], [238, 18], [238, 238], [18, 238], [18, 18]], 5);
    return textur(grova(c, 0.92));
  });
}

// revben är en ryggrad med revben, som syns genom golvet från nivå 3.
export function revben(): THREE.Texture {
  return spara("revben", () => {
    const [c, x] = duk(256, 256);
    streck(x, [[128, 18], [128, 238]], 8);
    for (let i = 0; i < 6; i++) {
      const y = 40 + i * 34;
      [-1, 1].forEach((s) => streck(x, [[128, y], [128 + s * 50, y - 14], [128 + s * 92, y + 6], [128 + s * 104, y + 30]], 6));
    }
    return textur(grova(c, 0.85));
  });
}

export function skugga(): THREE.Texture {
  return spara("skugga", () => {
    const [c, x] = duk(64, 64);
    const g = x.createRadialGradient(32, 32, 0, 32, 32, 30);
    g.addColorStop(0, "rgba(0,0,0,0.55)");
    g.addColorStop(1, "rgba(0,0,0,0)");
    x.fillStyle = g;
    x.fillRect(0, 0, 64, 64);
    return textur(c);
  });
}

// ---- Sten och trä ----
//
// Golvet, murarna och dörrarna är målade som spelbrickor: huggna
// stenplattor, murverk och plankor. De ritas en gång, med samma frö varje
// gång, och färgas av materialet där tornet förändras med avslöjandet.

// varie ger en färg nära bas, där varje kanal kan gå upp eller ner med spann.
function varie(bas: number[], spann: number): string {
  const d = (slump() - 0.5) * spann;
  return `rgb(${bas.map((v) => Math.max(0, Math.min(255, Math.round(v + d + (slump() - 0.5) * spann * 0.3)))).join(",")})`;
}

// sten fyller rutan (x, y, w, h) med en huggen sten: grundfärg, fläckar,
// en ljusare övre kant och en mörkare undre, som om ljuset kom snett uppifrån.
function sten(x: CanvasRenderingContext2D, rx: number, ry: number, w: number, h: number, bas: number[], fog: number): void {
  const r = Math.min(w, h) * 0.12;
  x.save();
  x.beginPath();
  x.roundRect(rx + fog, ry + fog, w - fog * 2, h - fog * 2, r);
  x.clip();
  x.fillStyle = varie(bas, 34);
  x.fillRect(rx, ry, w, h);
  for (let i = 0; i < (w * h) / 18; i++) {
    const v = slump() < 0.5 ? 0 : 255;
    x.fillStyle = `rgba(${v},${v},${v},${slump() * 0.07})`;
    const s = 1 + slump() * 3;
    x.fillRect(rx + slump() * w, ry + slump() * h, s, s);
  }
  const g = x.createLinearGradient(rx, ry, rx + w * 0.3, ry + h);
  g.addColorStop(0, "rgba(255,240,215,0.16)");
  g.addColorStop(0.5, "rgba(0,0,0,0)");
  g.addColorStop(1, "rgba(0,0,0,0.28)");
  x.fillStyle = g;
  x.fillRect(rx, ry, w, h);
  // En spricka ibland.
  if (slump() < 0.3) {
    x.strokeStyle = "rgba(20,16,12,0.55)";
    x.lineWidth = 1.2;
    x.beginPath();
    let px = rx + slump() * w, py = ry + slump() * h * 0.3;
    x.moveTo(px, py);
    for (let i = 0; i < 4; i++) {
      px += (slump() - 0.5) * w * 0.3;
      py += h * 0.2;
      x.lineTo(px, py);
    }
    x.stroke();
  }
  x.restore();
}

// stengolv är ett rums golv: oregelbundna stenplattor i rader, med fogar,
// lite mossa i fogarna och en mörkare kant längs väggarna.
export function stengolv(): THREE.Texture {
  return spara("stengolv", () => {
    const N = 512;
    const [c, x] = duk(N, N);
    x.fillStyle = "#2b2520";
    x.fillRect(0, 0, N, N);
    let y = 0;
    while (y < N - 4) {
      const h = Math.min(N - y, 70 + slump() * 50);
      let xx = -slump() * 60;
      while (xx < N) {
        const w = 80 + slump() * 90;
        sten(x, xx, y, w, h, [112, 102, 88], 3);
        xx += w;
      }
      y += h;
    }
    // Mossa i fogarna.
    for (let i = 0; i < 900; i++) {
      x.fillStyle = `rgba(${70 + slump() * 30},${90 + slump() * 40},${40 + slump() * 20},${slump() * 0.35})`;
      x.fillRect(slump() * N, slump() * N, 2, 2);
    }
    // Mörkare mot väggarna, där ingen går.
    const g = x.createRadialGradient(N / 2, N / 2, N * 0.3, N / 2, N / 2, N * 0.72);
    g.addColorStop(0, "rgba(0,0,0,0)");
    g.addColorStop(1, "rgba(0,0,0,0.45)");
    x.fillStyle = g;
    x.fillRect(0, 0, N, N);
    return textur(c);
  });
}

// murverk är murarnas stenblock i förband. Den upprepas efter murens
// storlek i världen (se varldsUv i karta.ts), så att blocken är lika
// stora överallt.
export function murverk(): THREE.Texture {
  return spara("murverk", () => {
    const N = 256;
    const [c, x] = duk(N, N);
    x.fillStyle = "#3a332c";
    x.fillRect(0, 0, N, N);
    const rad = N / 4;
    // Ett block som går över kanten ritas också på andra sidan, med samma
    // frö, så att muren går ihop när mönstret upprepas.
    const block = (bx: number, by: number, w: number) => {
      const s0 = seed;
      sten(x, bx, by, w, rad, [176, 166, 150], 2.5);
      if (bx + w > N) {
        seed = s0;
        sten(x, bx - N, by, w, rad, [176, 166, 150], 2.5);
      }
    };
    for (let r = 0; r < 4; r++) {
      let bx = r % 2 ? N / 4 : 0;
      const slut = bx + N;
      // Fyra block per rad, olika långa, så att mönstret inte blir rutigt.
      for (let i = 0; i < 4; i++) {
        const w = i === 3 ? slut - bx : (N / 4) * (0.75 + slump() * 0.5);
        block(bx, r * rad, w);
        bx += w;
      }
    }
    const t = textur(c);
    t.wrapS = t.wrapT = THREE.RepeatWrapping;
    return t;
  });
}

// tra är plankor med ådring, för dörrarna och inredningen.
export function tra(): THREE.Texture {
  return spara("tra", () => {
    const N = 128;
    const [c, x] = duk(N, N);
    const plank = N / 4;
    for (let i = 0; i < 4; i++) {
      x.fillStyle = varie([120, 78, 44], 26);
      x.fillRect(i * plank, 0, plank, N);
      for (let j = 0; j < 14; j++) {
        x.strokeStyle = `rgba(40,22,10,${0.15 + slump() * 0.25})`;
        x.lineWidth = 0.6 + slump();
        x.beginPath();
        const px = i * plank + slump() * plank;
        x.moveTo(px, 0);
        x.bezierCurveTo(px + (slump() - 0.5) * 8, N * 0.3, px + (slump() - 0.5) * 8, N * 0.7, px + (slump() - 0.5) * 6, N);
        x.stroke();
      }
      x.fillStyle = "rgba(20,10,4,0.7)";
      x.fillRect(i * plank, 0, 1.5, N);
    }
    const t = textur(c);
    t.wrapS = t.wrapT = THREE.RepeatWrapping;
    return t;
  });
}
