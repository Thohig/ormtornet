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
}

export interface Varelse {
  id: number;
  name: string;
  glyph: string;
  x: number;
  y: number;
}

export interface Tillstand {
  you: number; // 0 för spelledaren
  size: number;
  level: number;
  reveal: boolean;
  phase: string;
  action?: string;
  controller?: boolean;
  karta?: string;
  rooms: Rum[];
  players: Vandrare[];
  monsters?: Varelse[];
  lastSeen?: Record<string, { x: number; y: number; slot: number }>;
}

// Utsnittet som kartan visar, i rum: x och y är hörnet uppe till vänster
// (nordväst), w och h bredden och höjden. Den platta kartan räknar ut det.
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
