// Formen på tillståndet som servern skickar (stateMsg i game.go). Bara det
// kartan använder står här.

export interface Rum {
  x: number;
  y: number;
  d: [boolean, boolean, boolean, boolean]; // dörrar åt N, Ö, S, V
  k: string; // "rum" eller "kammare"
  exit: number; // porten ut: 0–3 (N, Ö, S, V), -1 om ingen
  name: string;
  mark?: string; // fällans tecken
  marki?: string; // fällans id i fallor.json, för modellen
  f?: boolean; // ett föremål eller guld ligger här (för spårfinnaren)
  dod?: boolean; // här har någon dött (för kartritaren)
  stangd?: boolean; // porten har slutit sig
  a?: boolean; // ett rum bara andra har hittat (kartritaren), för minikartan
  in?: string[]; // inredningen, efter rumstypen (inredning.ts)
}

export interface Vandrare {
  id: number;
  slot: number;
  x: number;
  y: number;
  online: boolean;
  status: string; // "inne", "ute", "forlorad" eller "dod"
  hidden: boolean;
  klassId?: string; // klassens id i klasser.json
  klassModell?: Modell; // modell i klasser.json
  gjorde?: string; // det vandraren gjorde under rundan: anfall, parera, undan, ta, sok, drick
  traffad?: boolean; // vandraren blev skadad under rundan
}

// Modell är färgen och storleken på en modell, från modell i klasser.json
// och monster.json.
export interface Modell {
  farg?: string;
  storlek?: number;
}

export interface Varelse {
  id: number;
  name: string;
  glyph: string;
  typ?: string; // varelsens id i monster.json
  modell?: Modell;
  x: number;
  y: number;
  gjorde?: string; // varelsens val under rundan
  traffad?: boolean; // varelsen blev träffad under rundan
}

export interface Tillstand {
  you: number; // den egna vandrarens id
  size: number;
  level: number;
  reveal: boolean;
  phase: string;
  action?: string;
  rooms: Rum[];
  foremal?: { id: string }[]; // den egna väskan
  players: Vandrare[];
  monsters?: Varelse[];
  fallna?: Varelse[]; // varelserna som föll under rundan
  round?: number;
  lastSeen?: Record<string, { x: number; y: number; slot: number }>;
  sleep?: number; // hur djupt ormen sover, 0–2, bara för den som står i kammaren (annars -1)
}

// Utsnittet som kartan visar, i rum: x och y är hörnet uppe till vänster
// (nordväst), w och h bredden och höjden. web/map.js räknar ut det.
export interface Utsnitt {
  x: number;
  y: number;
  w: number;
  h: number;
}

export interface Val {
  utsnitt: () => Utsnitt | null;
  fokus: () => { x: number; y: number } | null; // rummet kartan följer
  // Hur stor del av kartans bredd, från vänster, som loggen täcker (0–0,5).
  // Kameran ställer utsnittet i den del som syns.
  tackt?: () => number;
  // Hur stor del av bredden, från höger, som spelkortet täcker när det
  // vänds upp (0–0,3). Kameran håller den delen fri på en bred skärm.
  tacktHoger?: () => number;
  // Den del av fönstrets höjd som syns, i px från fönstrets överkant: under
  // översta raden och ovanför knapparna på telefonen. Tangenterna vid
  // dörrarna hålls där, så att de går att trycka på.
  fritt?: () => { topp: number; botten: number };

  farger: string[]; // vandrarnas färger, en per plats
}
