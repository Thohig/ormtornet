// Bygger klient/src till ett enda skript, web/iso/iso.js, som laddas som de
// andra skripten i web/ och nås genom den globala IsoKarta. Go bäddar in det
// med resten av web/, så det byggda skriptet checkas in.
import { readFileSync } from "node:fs";
import { defineConfig } from "vite";

// Three.js licens (MIT) ska följa med i det byggda skriptet. Minifieringen
// tar bort licenskommentarerna, så den läggs först i skriptet.
const threeLicense = readFileSync("node_modules/three/LICENSE", "utf8").trim();
const banner = "/*! Innehåller three.js (https://threejs.org).\n\n" + threeLicense + "\n*/";

export default defineConfig({
  plugins: [
    {
      name: "three-licens",
      generateBundle(_, bundle) {
        for (const chunk of Object.values(bundle)) {
          if (chunk.type === "chunk") chunk.code = banner + "\n" + chunk.code;
        }
      },
    },
  ],
  build: {
    outDir: "../web/iso",
    emptyOutDir: true,
    target: "es2020",
    lib: {
      entry: "src/main.ts",
      name: "IsoKarta",
      formats: ["iife"],
      fileName: () => "iso.js",
    },
  },
});
