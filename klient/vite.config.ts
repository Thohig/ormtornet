// Bygger klient/src till två skript som laddas som de andra skripten i web/:
//
//   vite build              src/main.ts till web/iso/iso.js (IsoKarta)
//   vite build --mode rtc   src/rtc.ts till web/rtc/rtc.js (TornetRTC)
//
// Go bäddar in dem med resten av web/, så de byggda skripten checkas in.
import { readFileSync } from "node:fs";
import { defineConfig } from "vite";

// Licenserna för det som byggs in ska följa med i skriptet. Minifieringen
// tar bort licenskommentarerna, så de läggs först i skriptet.
// licens är en fil, eller rader ur källkoden när paketet saknar en.
function banner(paket: { namn: string; url: string; licens: string; rader?: number }[]): string {
  return paket
    .map((p) => {
      let text = readFileSync(p.licens, "utf8");
      if (p.rader) text = text.split("\n").slice(0, p.rader).map((r) => r.replace(/^\/\/ ?/, "")).join("\n");
      return "/*! Innehåller " + p.namn + " (" + p.url + ").\n\n" + text.trim() + "\n*/";
    })
    .join("\n");
}

const skript = {
  iso: {
    entry: "src/main.ts",
    name: "IsoKarta",
    banner: [{ namn: "three.js", url: "https://threejs.org", licens: "node_modules/three/LICENSE" }],
  },
  rtc: {
    entry: "src/rtc.ts",
    name: "TornetRTC",
    banner: [
      { namn: "Trystero", url: "https://github.com/dmotz/trystero", licens: "node_modules/@trystero-p2p/core/LICENSE" },
      { namn: "noble-secp256k1", url: "https://github.com/paulmillr/noble-secp256k1", licens: "node_modules/@noble/secp256k1/LICENSE" },
      { namn: "QR Code Generator", url: "https://github.com/kazuhikoarase/qrcode-generator", licens: "node_modules/qrcode-generator/dist/qrcode.mjs", rader: 15 },
    ],
  },
};

export default defineConfig(({ mode }) => {
  const s = mode === "rtc" ? skript.rtc : skript.iso;
  const namn = mode === "rtc" ? "rtc" : "iso";
  const text = banner(s.banner);
  return {
    plugins: [
      {
        name: "licenser",
        generateBundle(_, bundle) {
          for (const chunk of Object.values(bundle)) {
            if (chunk.type === "chunk") chunk.code = text + "\n" + chunk.code;
          }
        },
      },
    ],
    build: {
      outDir: "../web/" + namn,
      emptyOutDir: true,
      target: "es2020",
      lib: {
        entry: s.entry,
        name: s.name,
        formats: ["iife"],
        fileName: () => namn + ".js",
      },
    },
  };
});
