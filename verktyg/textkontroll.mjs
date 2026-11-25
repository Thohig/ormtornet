// Textkontroll: regler som letar efter texter som ser trasiga ut, till
// exempel en {platshållare} kvar, "1 rundor", "1 lives" eller svenska ord i
// en engelsk vandring. Provspelarna (provspelare.mjs) och webbläsarprovet
// använder den, och textkontroll.prov.mjs provar reglerna mot felen som
// hittades i provspelen.
//
//   import { kontrollera, texterI } from "./textkontroll.mjs";
//   kontrollera("Du sitter fast i 1 rundor.", "sv")  // ["1 med plural"]
//
// Reglerna ska hellre missa ett fel än märka för mycket: en rad som märks
// läses av en människa, och för många märkta rader gör att ingen läser dem.

// Tecken och ord som inte ska stå kvar i en färdig text, på alla språk.
const allmanna = [
  [/[{}]/, "platshållare kvar"],
  [/undefined|<nil>|%!|\bNaN\b|\bnull\b|\[object/, "värde saknas"],
  [/\S {2,}\S/, "dubbla mellanslag"],
  [/(^|[^.])\.\.(?!\.)/, "två punkter"],
  [/\s[.,:;!?](?!\d)/, "mellanslag före skiljetecken"],
  [/(^|[^.])[.!?] [a-zåäö]/u, "liten bokstav efter punkt"],
  // Två tal i rad före ett ord, som "In round 4 24 gold". Översta raden på
  // sidan ("runda 1 8 s") har talen i var sitt fält och slutar på en bokstav.
  [/(^|[^\d.,])\d+ \d+ \p{L}{2,}/u, "två tal i rad"],
];

// Samma ord två gånger i rad, som "det det". Korta ord som kan stå två
// gånger med rätta ("att att" gör det inte, men "så så" kan) är undantagna.
const dubbelOrd = /(^|[^\p{L}])(\p{L}{3,}) \2(?![\p{L}])/iu;

// Svenska: en etta, eller ordet en/ett, följd av ett ord i plural. Talet
// kan ha ett tecken före (+1, −1). "1 av 6" och "1 liv" (samma form) går
// igenom; reglerna gäller ord som bara finns i plural i tornets texter.
const svPlural = new RegExp(
  "(^|[^\\d,.])[+−-]?1 (\\p{L}+ )?(" +
    [
      "rundor", "gånger", "ggr", "drycker", "skador", "sår", "mynt",
      "handsbredder", "steg", "kort", "näver", "nävar", "föremål",
      "vandrare", "varelser", "rum", "dörrar", "portar", "delar",
      "sekunder", "minuter", "spelare", "liv", "guldmynt",
    ].join("|") + ")(?![\\p{L}])",
  "u",
);
// Ord i svPlural som har samma form i singular. Där blir regeln fel om de
// inte har ett pluraladjektiv framför, så de kontrolleras bara med ett.
const svSammaForm = new Set([
  "sår", "mynt", "steg", "kort", "föremål", "vandrare", "rum", "spelare",
  "liv", "guldmynt",
]);

// Engelska: 1 följt av ett ord som slutar på s, med högst ett ord emellan
// ("1 waking cards"). Undantagen är ord som slutar på s i singular.
const enPlural = /(^|[^\d,.])[+−-]?1 (?:(\p{L}+) )?(\p{L}+s)\b/u;
const enSingularS = new Set([
  "is", "was", "has", "does", "goes", "this", "plus", "across", "less",
  "glass", "moss", "boss", "bonus", "status", "its", "his", "chess",
  "always", "us", "yes", "gas", "pass", "class", "grass", "access",
  "darkness", "witness", "fortress", "chorus", "mess", "dress", "press",
  "princess", "abyss", "brass", "cross", "loss", "kiss", "process",
]);
const enAre = /\b(are|were) (\w+ ){0,2}1(?!\d)\b/i;

// Svenska ord som inte ska stå i en engelsk text (en text som inte är
// översatt visas på svenska).
const svenskaOrd = /(^|[^\p{L}])(och|det|du|är|inte|som|på|för|med|till|har|han|hon|hen|rummet|tornet|ormen|vandringen)(?![\p{L}])/iu;

export function kontrollera(text, sprak = "sv") {
  if (typeof text !== "string") return [];
  const fel = [];
  for (const [re, namn] of allmanna) if (re.test(text)) fel.push(namn);
  if (dubbelOrd.test(text)) fel.push("samma ord två gånger");
  if (sprak === "en") {
    const m = text.match(enPlural);
    if (m && !enSingularS.has(m[3].toLowerCase())) fel.push("1 med plural");
    if (enAre.test(text)) fel.push("1 med are");
    if (/[åäöÅÄÖ]/.test(text)) fel.push("svenska bokstäver");
    else if (svenskaOrd.test(text)) fel.push("svenskt ord");
  } else {
    const m = text.match(svPlural);
    if (m && (!svSammaForm.has(m[3]) || (m[2] && /a $/.test(m[2])))) {
      fel.push("1 med plural");
    }
  }
  return fel;
}

// Nycklar i tillståndet som inte är texter för spelarna: id, färger,
// namn på modeller och bilder, och vandrarnas namn (som kan vara vad som
// helst).
const inteText = new Set([
  "t", "id", "bild", "farg", "k", "a", "marki", "in", "klassId", "status",
  "name", "namn", "tecken", "phase", "sprak", "pvp", "action", "mark",
  "modell", "ikon", "c", "s", "slot",
]);

// texterI går igenom ett meddelande från tornet och ger alla texter för
// spelarna i det, utan id:n och namn: bara strängar med ett mellanslag.
export function texterI(msg, ut = []) {
  if (typeof msg === "string") {
    if (msg.includes(" ")) ut.push(msg);
  } else if (Array.isArray(msg)) {
    for (const v of msg) texterI(v, ut);
  } else if (msg && typeof msg === "object") {
    for (const [k, v] of Object.entries(msg)) {
      if (!inteText.has(k)) texterI(v, ut);
    }
  }
  return ut;
}
