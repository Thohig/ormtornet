// Vägen till tornet när sidorna kommer från GitHub Pages (docs/github-pages.md).
//
// Då finns ingen server att fråga. Tornet körs i stället i en flik i samma
// webbläsare, skalet torn.html, och sidorna når det över en BroadcastChannel.
// Det här skriptet byter ut fetch för tornets adresser (de som börjar med /,
// som /intro och /spelledare/status) och WebSocket för /ws, så att resten av
// sidornas skript inte märker någon skillnad.
//
// När sidorna kommer från tornets program (main.go) gör skriptet ingenting.
// Från Pages sätter texter/*.js window.ORMTORNET_SIDA.
"use strict";

window.Natet = (() => {
  if (!window.ORMTORNET_SIDA || typeof BroadcastChannel === "undefined") return null;

  // Kanalens namn hör till sidans mapp, så att två torn på samma värd inte
  // hör varandra.
  const base = new URL(".", document.currentScript ? document.currentScript.src : location.href).pathname;
  const channel = new BroadcastChannel("ormtornet:" + base);
  const me = Math.random().toString(36).slice(2) + Date.now().toString(36);
  const realFetch = window.fetch.bind(window);
  const RealWebSocket = window.WebSocket;

  let towerSeen = false;
  const pending = new Map(); // förfrågningar som väntar på svar, på id
  const sockets = new Map(); // öppna anslutningar, på id
  let nextId = 1;

  channel.onmessage = (ev) => {
    const m = ev.data;
    if (!m || typeof m !== "object") return;
    if (m.typ === "har" || m.typ === "lever") {
      towerSeen = true;
      lastBeat = Date.now();
      return;
    }
    if (m.typ === "borta") {
      towerSeen = false;
      for (const s of sockets.values()) s._closed();
      return;
    }
    if (m.till !== me) return;
    if (m.typ === "svar") {
      const p = pending.get(m.id);
      if (!p) return;
      pending.delete(m.id);
      clearTimeout(p.timer);
      const empty = m.status === 204 || m.status === 205 || m.status === 304;
      p.resolve(new Response(empty ? null : m.kropp, { status: m.status, headers: m.innehall ? { "Content-Type": m.innehall } : {} }));
      return;
    }
    const s = sockets.get(m.conn);
    if (!s) return;
    if (m.typ === "oppen") s._opened();
    else if (m.typ === "fran") s._message(m.data);
    else if (m.typ === "stangd") s._closed();
  };

  // Tornets flik säger då och då att den lever. Hörs den inte på länge är
  // den stängd, och anslutningarna bryts så att sidan försöker igen.
  let lastBeat = 0;
  setInterval(() => {
    if (sockets.size) channel.postMessage({ typ: "lever", fran: me });
    if (towerSeen && Date.now() - lastBeat > 10000) {
      towerSeen = false;
      for (const s of sockets.values()) s._closed();
    }
  }, 3000);

  function isTowerPath(url) {
    return typeof url === "string" && url.startsWith("/") && !url.startsWith("//");
  }

  window.fetch = function (input, init) {
    if (!isTowerPath(input)) return realFetch(input, init);
    init = init || {};
    const headers = {};
    new Headers(init.headers || {}).forEach((v, k) => (headers[k] = v));
    const id = nextId++;
    return new Promise((resolve, reject) => {
      const timer = setTimeout(() => {
        pending.delete(id);
        reject(new TypeError("tornet svarar inte"));
      }, 10000);
      pending.set(id, { resolve, timer });
      channel.postMessage({
        typ: "fraga", id, fran: me, metod: (init.method || "GET").toUpperCase(), vag: input,
        huvud: headers, kropp: typeof init.body === "string" ? init.body : "",
      });
    });
  };

  // TornSocket ser ut som en WebSocket för web/app.js.
  class TornSocket {
    constructor() {
      this.id = me + "-" + nextId++;
      this.readyState = 0;
      this.onopen = this.onmessage = this.onclose = this.onerror = null;
      sockets.set(this.id, this);
      channel.postMessage({ typ: "oppna", conn: this.id, fran: me });
      this._timer = setTimeout(() => this._closed(), 5000);
    }
    send(data) {
      if (this.readyState !== 1) return;
      channel.postMessage({ typ: "till", conn: this.id, fran: me, data: String(data) });
    }
    close() {
      if (this.readyState >= 2) return;
      channel.postMessage({ typ: "stang", conn: this.id, fran: me });
      this._closed();
    }
    _opened() {
      if (this.readyState !== 0) return;
      clearTimeout(this._timer);
      this.readyState = 1;
      if (this.onopen) this.onopen({});
    }
    _message(data) {
      if (this.readyState === 1 && this.onmessage) this.onmessage({ data });
    }
    _closed() {
      if (this.readyState === 3) return;
      clearTimeout(this._timer);
      this.readyState = 3;
      sockets.delete(this.id);
      if (this.onclose) this.onclose({});
    }
  }
  TornSocket.CONNECTING = 0;
  TornSocket.OPEN = 1;
  TornSocket.CLOSING = 2;
  TornSocket.CLOSED = 3;

  window.WebSocket = function (url, protocols) {
    if (new URL(url, location.href).pathname === "/ws") return new TornSocket();
    return new RealWebSocket(url, protocols);
  };
  Object.assign(window.WebSocket, { CONNECTING: 0, OPEN: 1, CLOSING: 2, CLOSED: 3 });

  // Sidan stängs eller byts: säg till tornet direkt, så att vandraren inte
  // står kvar som ansluten.
  window.addEventListener("pagehide", () => {
    for (const s of sockets.values()) channel.postMessage({ typ: "stang", conn: s.id, fran: me });
  });

  // Introt öppnas utan att tornet är öppet i den här webbläsaren: då är det
  // spelledaren som kommer, och tornet öppnas i skalet.
  channel.postMessage({ typ: "finns", fran: me });
  const page = location.pathname.slice(base.length);
  if (window.top === window && (page === "" || page === "index.html")) {
    setTimeout(() => {
      if (!towerSeen) location.replace(base + "torn.html" + location.search);
    }, 1200);
  }

  return { base, towerSeen: () => towerSeen };
})();
