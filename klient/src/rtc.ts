// Vägen mellan webbläsare när tornet ligger på GitHub Pages
// (docs/github-pages.md, steg 3). Vite bygger det här till web/rtc/rtc.js,
// och sidorna når det genom den globala TornetRTC.
//
// Tornets flik (web/torn.js) och spelarsidorna (web/natet.js) går in i samma
// rum med rumskoden och nyckeln från länken. Trystero hittar varandra genom
// publika Nostr-reläer och kopplar sedan webbläsarna direkt med WebRTC.
// Nyckeln krypterar utbytet, så reläerna ser aldrig spelet.
import { joinRoom } from "trystero";

// Ett meddelande är samma slags objekt som går över BroadcastChannel i
// web/natet.js och web/torn.js.
type Meddelande = Record<string, unknown>;

export type Rum = {
  // skickar till en vandrare, eller till alla när till saknas
  skicka: (m: Meddelande, till?: string) => void;
  vidMeddelande: (fn: (m: Meddelande, fran: string) => void) => void;
  vidNy: (fn: (vem: string) => void) => void;
  vidBorta: (fn: (vem: string) => void) => void;
  lamna: () => void;
};

const APP = "ormtornet";

export function rum(kod: string, nyckel: string): Rum {
  const room = joinRoom({ appId: APP, password: nyckel }, kod);
  const action = room.makeAction<string>("torn");
  return {
    skicka(m, till) {
      action.send(JSON.stringify(m), till ? { target: till } : undefined).catch(() => {});
    },
    vidMeddelande(fn) {
      action.onMessage((data, ctx) => {
        let m: unknown;
        try {
          m = JSON.parse(String(data));
        } catch {
          return;
        }
        if (m && typeof m === "object" && !Array.isArray(m)) fn(m as Meddelande, ctx.peerId);
      });
    },
    vidNy(fn) {
      room.onPeerJoin((ctx) => fn(ctx.peerId));
    },
    vidBorta(fn) {
      room.onPeerLeave((ctx) => fn(ctx.peerId));
    },
    lamna() {
      room.leave().catch(() => {});
    },
  };
}

// En slumpad kod med bokstäver och siffror som går att skriva i en adress.
export function slump(langd: number): string {
  const tecken = "abcdefghijkmnpqrstuvwxyz23456789";
  const bytes = new Uint8Array(langd);
  crypto.getRandomValues(bytes);
  return Array.from(bytes, (b) => tecken[b % tecken.length]).join("");
}
