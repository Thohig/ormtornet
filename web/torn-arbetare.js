// Tornet i en Web Worker (docs/github-pages.md, steg 4). Skalet
// (web/torn.js) startar det här skriptet, som kör torn.wasm.
//
// Webbläsare saktar ner timers i flikar som inte syns, efter en stund till en
// gång i minuten. Tornets tid, robotarna och pulsen till sidorna går med
// timers, och i en worker saktas de inte ner på det sättet. Skalet för
// bara vidare meddelanden, och det går lika fort i en dold flik.
//
// Meddelanden från skalet: starta, fraga, anslut, taEmot, kopplaNer.
// Till skalet: redo, fel, stangt, svar, ansluten, skickar, kopplarNer,
// spara och puls.
"use strict";

importScripts("wasm_exec.js");

const post = (m) => self.postMessage(m);

// Det som Go anropar (wasm.go).
self.tornetSkickar = (id, text) => post({ typ: "skickar", id, text });
self.tornetKopplarNer = (id) => post({ typ: "kopplarNer", id });
self.tornetRedo = () => post({ typ: "redo" });
self.tornetStangt = () => post({ typ: "stangt" });
self.tornetFel = (fel) => post({ typ: "fel", fel: String(fel) });

// Inställningarna sparas i skalets localStorage, som en worker inte har.
// Skalet skickar med det som finns sparat när tornet startar.
function lagring(saved) {
  const data = Object.assign({}, saved);
  return {
    getItem: (k) => (Object.prototype.hasOwnProperty.call(data, k) ? data[k] : null),
    setItem: (k, v) => {
      data[k] = String(v);
      post({ typ: "spara", nyckel: k, varde: String(v) });
    },
  };
}

async function start(m) {
  try {
    self.ORMTORNET_ADRESS = m.adress;
    self.localStorage = lagring(m.lagring);
    const go = new Go();
    let result;
    try {
      result = await WebAssembly.instantiateStreaming(fetch("torn.wasm"), go.importObject);
    } catch {
      // Servern anger inte application/wasm: läs filen i stället.
      const bytes = await (await fetch("torn.wasm")).arrayBuffer();
      result = await WebAssembly.instantiate(bytes, go.importObject);
    }
    go.run(result.instance);
  } catch (err) {
    self.tornetFel((err && err.message) || err);
  }
}

// Pulsen: skalet säger till sidorna att tornet lever var gång den kommer.
setInterval(() => post({ typ: "puls" }), 3000);

self.onmessage = (ev) => {
  const m = ev.data;
  if (!m || typeof m !== "object") return;
  if (m.typ === "starta") {
    start(m);
    return;
  }
  if (!self.Tornet) return;
  switch (m.typ) {
    case "fraga":
      Tornet.fraga(m.metod, m.vag, m.huvud || {}, m.kropp || "", (status, type, body) =>
        post({ typ: "svar", id: m.id, status, innehall: type, kropp: body }), !!m.lokal);
      break;
    case "anslut":
      post({ typ: "ansluten", id: m.id, ok: !!Tornet.anslut(m.id, !!m.lokal) });
      break;
    case "taEmot":
      Tornet.taEmot(m.id, m.text);
      break;
    case "kopplaNer":
      Tornet.kopplaNer(m.id);
      break;
  }
};
