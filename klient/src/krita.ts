// Krita på canvas: texturerna för golvet, väggarna och tecknen på golvet.
// De ritas vita och färgas av materialet, så att en textur räcker för alla
// färger.
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

export function skiffer(): THREE.Texture {
  return spara("skiffer", () => {
    const [c, x] = duk(256, 256);
    x.fillStyle = "#1c1f23";
    x.fillRect(0, 0, 256, 256);
    for (let i = 0; i < 6; i++) {
      const g = x.createRadialGradient(slump() * 256, slump() * 256, 0, slump() * 256, slump() * 256, 120);
      g.addColorStop(0, "rgba(255,255,255,0.035)");
      g.addColorStop(1, "rgba(255,255,255,0)");
      x.fillStyle = g;
      x.fillRect(0, 0, 256, 256);
    }
    for (let i = 0; i < 2600; i++) {
      const v = 40 + slump() * 60;
      x.fillStyle = `rgba(${v},${v + 3},${v + 6},${0.25 + slump() * 0.3})`;
      x.fillRect(slump() * 256, slump() * 256, 1, 1);
    }
    x.strokeStyle = "rgba(255,255,255,0.05)";
    x.lineWidth = 2;
    x.strokeRect(2, 2, 252, 252);
    return textur(c);
  });
}

export function vagg(): THREE.Texture {
  return spara("vagg", () => {
    const [c, x] = duk(128, 128);
    x.fillStyle = "#fff";
    x.fillRect(0, 0, 128, 128);
    for (let i = 0; i < 1800; i++) {
      x.fillStyle = `rgba(60,60,60,${slump() * 0.28})`;
      x.fillRect(slump() * 128, slump() * 128, 1 + slump() * 2, 1);
    }
    const t = textur(c);
    t.wrapS = t.wrapT = THREE.RepeatWrapping;
    return t;
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
