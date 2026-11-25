// Provar textkontroll.mjs mot felen som hittades i provspelen och mot rader
// som är rätt. Körs i CI:
//
//   node verktyg/textkontroll.prov.mjs
import { kontrollera, texterI } from "./textkontroll.mjs";

const fel = [
  ["sv", "Du sitter fast i 1 rundor."],
  ["sv", "Du har varit i kammaren 1 ggr."],
  ["sv", "+1 drycker"],
  ["sv", "Det kommer {n} rundor till."],
  ["sv", "Du hittar undefined."],
  ["sv", "Det är mörkt  här."],
  ["sv", "Du går in i rummet ."],
  ["sv", "Det det luktar rök."],
  ["sv", "Klart.. Nu går du."],
  ["en", "−1 lives"],
  ["en", "You take 1 wounds."],
  ["en", "1 waking cards left"],
  ["en", "there are at least 1 of you"],
  ["en", "You find rep med kroken."],
  ["en", "The door is locked och dark."],
  ["en", "It comes several seconds apart. the snake stirs."],
  ["en", "In round 4 24 gold lay hidden in a room."],
];

const ratt = [
  ["sv", "Du sitter fast i 1 runda."],
  ["sv", "Du har 1 av 6 liv kvar."],
  ["sv", "−1 liv"],
  ["sv", "Solen står 6 handsbredder över horisonten."],
  ["sv", "Du hittar 1 mynt."],
  ["sv", "Det finns 12 rundor kvar."],
  ["sv", "Vänta... något rör sig."],
  ["sv", "Klockan är 21.10 och det är 1,5 meter."],
  ["en", "1 life left"],
  ["en", "You take 1 wound."],
  ["en", "1 waking card left"],
  ["en", "there is at least 1 of you"],
  ["en", "1 gold piece is gone."],
  ["en", "You have 2 lives."],
  ["en", "Wait... something moves."],
  ["en", "The door is locked and dark."],
  ["en", "In round 4, 24 gold lay hidden in a room."],
  ["sv", "runda 1 8 s"],
];

let misslyckade = 0;
for (const [sprak, text] of fel) {
  if (kontrollera(text, sprak).length === 0) {
    console.log(`missade (${sprak}): ${text}`);
    misslyckade++;
  }
}
for (const [sprak, text] of ratt) {
  const f = kontrollera(text, sprak);
  if (f.length) {
    console.log(`märkte fel (${sprak}): ${text}  [${f.join(", ")}]`);
    misslyckade++;
  }
}

const t = texterI({
  t: "state", id: "legoknekt", farg: "#7a5a32", name: "Ett namn med mellanslag",
  vandring: { rundor: [{ rum: "Du kliver in.", rader: ["En rad här."] }] },
});
if (t.join("|") !== "Du kliver in.|En rad här.") {
  console.log(`texterI gav: ${JSON.stringify(t)}`);
  misslyckade++;
}

if (misslyckade) {
  console.log(`${misslyckade} fel i textkontrollen`);
  process.exit(1);
}
console.log(`textkontrollen: ${fel.length} fel hittade, ${ratt.length} rätta rader godkända`);
