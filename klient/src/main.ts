// Tornet snett ovanifrån, ritat med Three.js. Vite bygger det här till
// web/iso/iso.js, och sidorna når det genom den globala IsoKarta.
// Se docs/isometrisk.md.
import { REVISION } from "three";

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
