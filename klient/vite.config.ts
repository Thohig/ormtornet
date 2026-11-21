// Bygger klient/src till ett enda skript, web/iso/iso.js, som laddas som de
// andra skripten i web/ och nås genom den globala IsoKarta. Go bäddar in det
// med resten av web/, så det byggda skriptet checkas in.
import { defineConfig } from "vite";

export default defineConfig({
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
