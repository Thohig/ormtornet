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
  c: number; // pilen i krita, -1 om ingen
  f?: boolean; // ett föremål ligger här (för spårfinnaren)
  stangd?: boolean; // porten har slutit sig
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
  lastSeen?: Record<string, { x: number; y: number; slot: number }>;
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

  farger: string[]; // vandrarnas färger, en per plats
}
