// Tornet snett ovanifrån, ritat med Three.js. Vite bygger det här till
// web/iso/iso.js, och sidorna når det genom den globala IsoKarta. web/map.js
// väljer mellan den här kartan och den platta. Se docs/isometrisk.md.
import { REVISION } from "three";
import { IsoKarta } from "./karta";
import type { Val } from "./typer";

// Three.js-versionen som är inbyggd, för felsökning i konsolen.
export const three = REVISION;

// Om webbläsaren kan rita med WebGL. Annars används den platta kartan.
export function stods(): boolean {
  try {
    const canvas = document.createElement("canvas");
    return !!(canvas.getContext("webgl2") || canvas.getContext("webgl"));
  } catch {
    return false;
  }
}

// skapa ritar kartan på canvas. Den platta kartan ger utsnittet.
export function skapa(canvas: HTMLCanvasElement, val: Val): IsoKarta {
  return new IsoKarta(canvas, val);
}
