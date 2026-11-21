// Vägen till tornet när sidorna kommer från GitHub Pages (docs/github-pages.md).
//
// Då finns ingen server att fråga. Tornet körs i stället i en flik i
// spelledarens webbläsare, skalet torn.html. Det här skriptet byter ut fetch
// för tornets adresser (de som börjar med /, som /intro och
// /spelledare/status) och WebSocket för /ws, så att resten av sidornas
// skript inte märker någon skillnad. Vägen dit är en av två:
//
// - I spelledarens egen webbläsare går allt över en BroadcastChannel.
// - Från länken med rumskoden (#torn=...&nyckel=...) går allt över WebRTC
//   till spelledarens webbläsare, genom web/rtc/rtc.js.
//
// Meddelandena är desamma åt båda hållen. När sidorna kommer från tornets
// program (main.go) gör skriptet ingenting. Från Pages sätter texter/*.js
// window.ORMTORNET_SIDA.
"use strict";

window.Natet = (() => {
  if (!window.ORMTORNET_SIDA) return null;

  const script = document.currentScript ? document.currentScript.src : location.href;
  // Mappen sidorna ligger i, så att två torn på samma värd inte hör varandra.
  const base = new URL(".", script).pathname;
  const link = new URLSearchParams(location.hash.slice(1));
  const remote = !!(link.get("torn") && link.get("nyckel"));
  if (!remote && typeof BroadcastChannel === "undefined") return null;

  const me = Math.random().toString(36).slice(2) + Date.now().toString(36);
  const realFetch = window.fetch.bind(window);
  const RealWebSocket = window.WebSocket;

  let towerSeen = false;
  let lastBeat = 0;
  let waiting = []; // det som väntar på att tornet ska höras
  const pending = new Map(); // förfrågningar som väntar på svar, på id
  const sockets = new Map(); // öppna anslutningar, på id
  let nextId = 1;

  // post skickar till tornet.
  let post;
  if (remote) {
    let room = null;
    let tower = null; // tornets flik bland dem i rummet
    const queue = [];
    post = (m) => {
      if (room && tower) room.skicka(m, tower);
      else queue.push(m);
    };
    const start = () => {
      room = window.TornetRTC.rum(link.get("torn"), link.get("nyckel"));
      room.vidMeddelande((m, from) => {
        if (m.typ === "har" || m.typ === "lever") {
          if (tower !== from) {
            tower = from;
            for (const q of queue.splice(0)) room.skicka(q, tower);
          }
        } else if (from !== tower) {
          return; // bara tornet svarar
        }
        receive(m);
      });
      room.vidNy((peer) => room.skicka({ typ: "finns", fran: me }, peer));
      room.vidBorta((peer) => {
        if (peer !== tower) return;
        tower = null;
        lost();
      });
      window.addEventListener("pagehide", () => room.lamna());
    };
    // rtc.js laddas bara när länken behöver det.
    const s = document.createElement("script");
    // Med samma versionsnummer som det här skriptet (sida.go).
    const u = new URL("rtc/rtc.js", script);
    u.search = new URL(script).search;
    s.src = u.href;
    s.onload = start;
    s.onerror = () => status("natet_fel");
    document.head.appendChild(s);
  } else {
    const channel = new BroadcastChannel("ormtornet:" + base);
    channel.onmessage = (ev) => {
      const m = ev.data;
      // Det som andra sidor skickar till tornet har fran; det tornet
      // skickar har det inte.
      if (!m || typeof m !== "object" || m.fran) return;
      if (m.till && m.till !== me) return;
      receive(m);
    };
    post = (m) => channel.postMessage(m);
  }

  function receive(m) {
    if (m.typ === "har" || m.typ === "lever") {
      lastBeat = Date.now();
      if (!towerSeen) {
        towerSeen = true;
        status(null);
        for (const fn of waiting.splice(0)) fn();
      }
      return;
    }
    if (m.typ === "borta") {
      lost();
      return;
    }
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
  }

  // Tornet är borta: anslutningarna bryts, så att sidan försöker igen.
  function lost() {
    if (!towerSeen) return;
    towerSeen = false;
    for (const s of sockets.values()) s._closed();
    if (remote) searching();
  }

  // whenTower kör fn när tornet hörs, direkt om det redan har hörts.
  function whenTower(fn) {
    if (towerSeen) fn();
    else waiting.push(fn);
  }

  // Tornets flik säger då och då att den lever. Hörs den inte på länge är
  // den stängd.
  setInterval(() => {
    if (sockets.size) post({ typ: "lever", fran: me });
    if (towerSeen && Date.now() - lastBeat > 10000) lost();
  }, 3000);

  // Med länken syns en rad längst ned medan sidan letar efter tornet.
  let statusTimers = [];
  function status(key) {
    let el = document.getElementById("natetLage");
    if (!key) {
      if (el) el.remove();
      return;
    }
    if (!el) {
      el = document.createElement("div");
      el.id = "natetLage";
      el.setAttribute("role", "status");
      el.style.cssText = "position:fixed;left:0;right:0;bottom:0;padding:10px 16px;text-align:center;" +
        "font-size:13px;color:var(--warn,#e0b060);background:rgba(0,0,0,.85);z-index:1000";
      document.body.appendChild(el);
    }
    el.textContent = typeof t === "function" ? t(key) : key;
  }
  function searching() {
    for (const id of statusTimers) clearTimeout(id);
    statusTimers = [
      setTimeout(() => !towerSeen && status("natet_letar"), 3000),
      setTimeout(() => !towerSeen && status("natet_hittar_inte"), 30000),
    ];
  }
  if (remote) {
    if (document.body) searching();
    else document.addEventListener("DOMContentLoaded", searching);
  }

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
      const fail = () => {
        pending.delete(id);
        reject(new TypeError("tornet svarar inte"));
      };
      // Med länken får tornet längre tid på sig att hittas första gången.
      const p = { resolve, timer: setTimeout(fail, remote ? 60000 : 10000) };
      pending.set(id, p);
      whenTower(() => {
        if (!pending.has(id)) return;
        clearTimeout(p.timer);
        p.timer = setTimeout(fail, 10000);
        post({
          typ: "fraga", id, fran: me, metod: (init.method || "GET").toUpperCase(), vag: input,
          huvud: headers, kropp: typeof init.body === "string" ? init.body : "",
        });
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
      whenTower(() => {
        if (this.readyState !== 0) return;
        post({ typ: "oppna", conn: this.id, fran: me });
        this._timer = setTimeout(() => this._closed(), remote ? 15000 : 5000);
      });
    }
    send(data) {
      if (this.readyState !== 1) return;
      post({ typ: "till", conn: this.id, fran: me, data: String(data) });
    }
    close() {
      if (this.readyState >= 2) return;
      post({ typ: "stang", conn: this.id, fran: me });
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
    for (const s of sockets.values()) post({ typ: "stang", conn: s.id, fran: me });
  });

  // Introt öppnas utan länk och utan att tornet är öppet i den här
  // webbläsaren: då är det spelledaren som kommer, och tornet öppnas i
  // skalet.
  if (!remote) {
    post({ typ: "finns", fran: me });
    const page = location.pathname.slice(base.length);
    if (window.top === window && (page === "" || page === "index.html")) {
      setTimeout(() => {
        if (!towerSeen) location.replace(base + "torn.html" + location.search);
      }, 1200);
    }
  }

  return { base, remote, towerSeen: () => towerSeen };
})();
