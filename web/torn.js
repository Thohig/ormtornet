// Skalet: tornets flik när spelet kommer från GitHub Pages
// (docs/github-pages.md). Det kör tornet, torn.wasm som är byggt av samma
// Go-kod som programmet, och för vidare det som sidorna skickar
// (web/natet.js) till tornet och tillbaka: från sidorna i den här
// webbläsaren över en BroadcastChannel, och från spelarna som kommer med
// länken över WebRTC (web/rtc/rtc.js) i ett rum med en egen rumskod.
"use strict";

(() => {
  const $ = (id) => document.getElementById(id);
  const base = new URL(".", location.href).pathname;
  const channel = new BroadcastChannel("ormtornet:" + base);
  const conns = new Map(); // anslutningar: id -> { fran, svara }
  const heard = new Map(); // sidor: fran -> senast hörd
  let room = null; // rummet för spelarna som kommer med länken

  // Rumskoden och nyckeln i länken. De gäller så länge fliken är öppen,
  // också om den laddas om, men ett nytt torn får en ny länk.
  const link = (() => {
    try {
      const saved = JSON.parse(sessionStorage.getItem("ormtornet-rum") || "null");
      if (saved && saved.kod && saved.nyckel) return saved;
    } catch {}
    const made = { kod: TornetRTC.slump(20), nyckel: TornetRTC.slump(20) };
    try {
      sessionStorage.setItem("ormtornet-rum", JSON.stringify(made));
    } catch {}
    return made;
  })();
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
  window.addEventListener("pagehide", () => {
    channel.postMessage({ typ: "borta" });
    if (room) {
      room.skicka({ typ: "borta" });
      room.lamna();
    }
  });

  // Det som Go anropar (wasm.go).
  window.tornetSkickar = (id, text) => {
    const c = conns.get(id);
    if (c) c.svara({ typ: "fran", till: c.fran, conn: id, data: text });
  };
  window.tornetKopplarNer = (id) => {
    const c = conns.get(id);
    if (!c) return;
    conns.delete(id);
    c.svara({ typ: "stangd", till: c.fran, conn: id });
  };
  window.tornetStangt = () => {
    closed = true;
    channel.postMessage({ typ: "borta" });
    if (room) room.skicka({ typ: "borta" });
    show("skal_stangt", null, true);
  };
  window.tornetFel = (msg) => {
    closed = true;
    show("skal_fel", { fel: msg }, true);
  };
  window.tornetRedo = () => {
    ready = true;
    channel.postMessage({ typ: "har" });
    openRoom();
    // Skalets egna texter på tornets språk.
    Tornet.fraga("GET", "/texter/vardsidan.json", {}, "", (status, _type, body) => {
      if (status !== 200) return;
      try {
        window.UI = JSON.parse(body);
        applyTexts();
      } catch {}
    }, true);
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

  // handle tar emot ett meddelande från en sida. svara skickar tillbaka till
  // den. lokal är sant för sidorna i den här webbläsaren, som är
  // spelledarens, som på tornets dator.
  function handle(m, svara, lokal) {
    if (!ready || closed || !m || typeof m !== "object" || !m.fran) return;
    heard.set(m.fran, Date.now());
    const mine = (id) => {
      const c = conns.get(id);
      return c && c.fran === m.fran;
    };
    switch (m.typ) {
      case "finns":
        svara({ typ: "har" });
        break;
      case "fraga":
        Tornet.fraga(String(m.metod), String(m.vag), m.huvud || {}, String(m.kropp || ""), (status, type, body) =>
          svara({ typ: "svar", till: m.fran, id: m.id, status, innehall: type, kropp: body }), lokal);
        break;
      case "oppna":
        if (conns.has(m.conn) && !mine(m.conn)) return;
        conns.set(m.conn, { fran: m.fran, svara });
        if (Tornet.anslut(String(m.conn), lokal)) svara({ typ: "oppen", till: m.fran, conn: m.conn });
        else window.tornetKopplarNer(m.conn);
        break;
      case "till":
        if (mine(m.conn)) Tornet.taEmot(String(m.conn), String(m.data));
        break;
      case "stang":
        if (mine(m.conn)) Tornet.kopplaNer(String(m.conn));
        break;
    }
  }

  channel.onmessage = (ev) => handle(ev.data, (svar) => channel.postMessage(svar), true);

  // Spelarna som kommer med länken. Varje vandrare i rummet räknas som en
  // egen sida, vad sidan än säger om sig själv, och är aldrig lokal.
  function openRoom() {
    room = TornetRTC.rum(link.kod, link.nyckel);
    room.vidMeddelande((m, peer) => {
      if (!m || typeof m !== "object") return;
      m.fran = peer;
      handle(m, (svar) => room.skicka(svar, peer), false);
    });
    room.vidNy((peer) => ready && !closed && room.skicka({ typ: "har" }, peer));
    room.vidBorta((peer) => {
      for (const [id, c] of conns) if (c.fran === peer) Tornet.kopplaNer(id);
    });
  }

  // Säg då och då att tornet lever, och släpp anslutningar från sidor som
  // inte har hörts av på länge (fliken stängdes utan att säga till).
  setInterval(() => {
    if (!ready || closed) return;
    channel.postMessage({ typ: "lever" });
    if (room) room.skicka({ typ: "lever" });
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
      // Adressen på värdsidan är länken med rumskoden.
      window.ORMTORNET_ADRESS = new URL(".", location.href).href + "#torn=" + link.kod + "&nyckel=" + link.nyckel;
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
