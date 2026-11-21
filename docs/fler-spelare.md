# Plan: klasser, större torn och upp till 8 spelare

Status: planen skrevs 2026-11-18. Användaren bad om allt, i lämplig
ordning. Varje steg byggs, committas och provas innan nästa.

| Steg | Vad | Läge |
|---|---|---|
| 1 | Klasser | byggt |
| 2 | Tornets storlek blir rörlig i koden, fortfarande högst 4 vandrare | byggt |
| 3 | Upp till 8 vandrare och automatisk skalning | byggt |

Sedan 2026-11-22 går tornets storlek, solen och klasserna inte att välja
på värdsidan: storleken och solen är alltid automatiska, och klasserna är
alltid med.

## Steg 1: klasser (byggt)

README beskriver klasserna under "Klasser". Besluten:

- Klassen väljs på en egen skärm efter lobbyn, med ett spelkort för varje
  klass och en tidsgräns på 30 sekunder. Ingen klass finns två gånger i
  samma torn. Den som inte väljer får ett slumpat ledigt kort (ändrat
  2026-11-19; förut två klasser att välja mellan i lobbyn).
- Klassen gäller en vandring. Nytt torn leder till lobbyn och sedan till
  ett nytt klassval.
- Klassen syns i listan över vandrare, även när man inte ser varandra.
- AI-spelarna väljer efter personlighet (`klasser` i `data/ai.json`),
  men först när alla människor har valt (ändrat 2026-11-19).
- Klassens liv läggs till svårighetens liv, så att Mardröm förblir hårdare.
- "Klasser: på / av" på värdsidan, med `-klasser`. Förval: på.

## Steg 2: tornets storlek i koden (byggt)

Byggt som planerat. Storleken är `g.size`, portarna `g.corners()` och
kammaren `g.center()`. "Automatisk" ger 9 × 9 tills steg 3. Provat med hela
vandringar av AI-spelare på 9, 11 och 13. Kartan fungerar på 13 × 13 på en
dator, men på en telefon blir rummen mycket små: det löses i steg 3.

Det ursprungliga förslaget:

`gridSize` och `maxPlayers` är konstanter, och rummen och de kända rummen
ligger i tabeller med fast storlek. De blir värden som sätts när tornet
byggs. Inget ändras för spelarna ännu, men en storlek kan väljas för hand
på värdsidan ("Tornets storlek: 9 / 11 / 13") för att prova.

- Udda storlekar, så att kammaren står i mitten.
- Portarna räknas från storleken.
- Klienten ritar redan efter den storlek servern skickar. Kartan på en
  telefon behöver provas på 13×13.

## Steg 3: upp till 8 vandrare och automatisk skalning (byggt)

Byggt som planerat, med de här besluten:

- Solen: i stället för ett nytt val "automatisk" räknas värdsidans antal
  rundor för ett torn på 9 × 9 och skalas med storleken. Då gäller sparade
  inställningar som förut.
- Tornet byggs om i lobbyn när antalet vandrare ändras. Portarna delas ut
  i platsordning, och lottas när vandringen börjar om vandrarna inte ser
  varandra.
- Ett för litet val av storlek höjs till det antalet vandrare kräver.
- Kartan visar bara det vandraren känner till och zoomar ut efter hand.
  Blir det kända större än 9 × 9 rum och tecknen mindre än 8 bildpunkter
  följer kartan den egna vandraren. Ett tryck på kartan växlar.

Det ursprungliga förslaget:

| Vandrare | Torn | Väg port–kammare | Portar | Sol |
|---|---|---|---|---|
| 2–4 | 9×9 | 8 steg | 4 i hörnen | 40 |
| 5–6 | 11×11 | 10 steg | 8 | 50 |
| 7–8 | 13×13 | 12 steg | 8 | 60 |

- Storleken bestäms när vandringen börjar, efter antalet vandrare med
  AI-spelarna. Lobbyn visar hur stort tornet blir.
- Fyra nya portar mitt på varje sida när det är fler än 4 vandrare.
- Solen och ormleken skalar: två extra sovkort per vandrare över fyra.
- Värdsidan: "Tornets storlek: automatisk / 9 / 11 / 13" och
  "Solnedgång: automatisk / rundor". AI-spelare 0–7.
- Åtta färger som går att skilja åt, åtta platser per rum och `@×5` när
  fler än fyra står i samma rum.
- På en telefon följer kartan den egna vandraren och visar ett utsnitt av
  ett stort torn, med möjlighet att se hela.
- Minst åtta namn för AI-spelarna i `data/ai.json`.
