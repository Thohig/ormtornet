// Skalet: tornets flik när spelet kommer från GitHub Pages
// (docs/github-pages.md). Det kör tornet, torn.wasm som är byggt av samma
// Go-kod som programmet, och för vidare det som sidorna i webbläsaren
// skickar över kanalen (web/natet.js) till tornet och tillbaka.
"use strict";

(() => {
  const $ = (id) => document.getElementById(id);
  const base = new URL(".", location.href).pathname;
  const channel = new BroadcastChannel("ormtornet:" + base);
  const conns = new Map(); // anslutningar: id -> { fran, seen }
  const heard = new Map(); // sidor: fran -> senast hörd
  let ready = false;
  let closed = false;

  function show(key, vars, again) {
    $("skal").hidden = false;
    $("sida").hidden = true;
    $("skalText").textContent = t(key, vars);
    $("skalIgen").hidden = !again;
  }

  $("skalIgen").onclick = () => location.reload();

  // Tornet stänger med fliken. Fråga först, så att det inte händer av misstag.
  function warnOnLeave(ev) {
    if (!ready || closed) return;
    ev.preventDefault();
    ev.returnValue = "";
  }
  window.addEventListener("beforeunload", warnOnLeave);
  window.addEventListener("pagehide", () => channel.postMessage({ typ: "borta" }));

  // Det som Go anropar (wasm.go).
  window.tornetSkickar = (id, text) => {
    const c = conns.get(id);
    if (c) channel.postMessage({ typ: "fran", till: c.fran, conn: id, data: text });
  };
  window.tornetKopplarNer = (id) => {
    const c = conns.get(id);
    if (!c) return;
    conns.delete(id);
    channel.postMessage({ typ: "stangd", till: c.fran, conn: id });
  };
  window.tornetStangt = () => {
    closed = true;
    channel.postMessage({ typ: "borta" });
    show("skal_stangt", null, true);
  };
  window.tornetFel = (msg) => {
    closed = true;
    show("skal_fel", { fel: msg }, true);
  };
  window.tornetRedo = () => {
    ready = true;
    channel.postMessage({ typ: "har" });
    // Skalets egna texter på tornets språk.
    Tornet.fraga("GET", "/texter/vardsidan.json", {}, "", (status, _type, body) => {
      if (status !== 200) return;
      try {
        window.UI = JSON.parse(body);
        applyTexts();
      } catch {}
    });
    const frame = $("sida");
    frame.src = "./" + location.search;
    frame.hidden = false;
    $("skal").hidden = true;
  };

  // Fliken får sidans titel, så att den känns igen bland flikarna.
  setInterval(() => {
    try {
      const title = $("sida").contentDocument && $("sida").contentDocument.title;
      if (title && !closed) document.title = title;
    } catch {}
  }, 1000);

  channel.onmessage = (ev) => {
    const m = ev.data;
    if (!ready || closed || !m || typeof m !== "object" || !m.fran) return;
    heard.set(m.fran, Date.now());
    switch (m.typ) {
      case "finns":
        channel.postMessage({ typ: "har" });
        break;
      case "fraga":
        Tornet.fraga(String(m.metod), String(m.vag), m.huvud || {}, String(m.kropp || ""), (status, type, body) =>
          channel.postMessage({ typ: "svar", till: m.fran, id: m.id, status, innehall: type, kropp: body }));
        break;
      case "oppna":
        conns.set(m.conn, { fran: m.fran });
        // Allt i den här webbläsaren är spelledarens: lokalt, som på
        // tornets dator.
        if (Tornet.anslut(String(m.conn), true)) channel.postMessage({ typ: "oppen", till: m.fran, conn: m.conn });
        else window.tornetKopplarNer(m.conn);
        break;
      case "till":
        if (conns.has(m.conn)) Tornet.taEmot(String(m.conn), String(m.data));
        break;
      case "stang":
        if (conns.has(m.conn)) Tornet.kopplaNer(String(m.conn));
        break;
    }
  };

  // Säg då och då att tornet lever, och släpp anslutningar från sidor som
  // inte har hörts av på länge (fliken stängdes utan att säga till).
  setInterval(() => {
    if (!ready || closed) return;
    channel.postMessage({ typ: "lever" });
    const now = Date.now();
    for (const [id, c] of conns) {
      if (now - (heard.get(c.fran) || 0) > 120000) Tornet.kopplaNer(id);
    }
  }, 3000);

  async function start() {
    try {
      const go = new Go();
      let result;
      try {
        result = await WebAssembly.instantiateStreaming(fetch("torn.wasm"), go.importObject);
      } catch {
        // Servern anger inte application/wasm: läs filen i stället.
        const bytes = await (await fetch("torn.wasm")).arrayBuffer();
        result = await WebAssembly.instantiate(bytes, go.importObject);
      }
      window.ORMTORNET_ADRESS = new URL(".", location.href).href;
      go.run(result.instance);
    } catch (err) {
      window.tornetFel(String((err && err.message) || err));
    }
  }

  show("skal_vaknar");
  // Bara ett torn i taget i den här webbläsaren.
  if (navigator.locks) {
    navigator.locks.request("ormtornet-torn:" + base, { ifAvailable: true }, (lock) => {
      if (!lock) {
        show("skal_redan", null, true);
        return;
      }
      start();
      return new Promise(() => {}); // håll låset så länge fliken är öppen
    });
  } else {
    start();
  }
})();
