# Samlat förslag: fröet, hela bredden och värdsidan under vandringen

Förslag 2026-11-22, byggt samma dag i steg 1–6. En avvikelse: loggen
står **nere** till vänster på kartan, med den senaste rundan nederst, och
inte uppe. Då syns alltid de senaste raderna även när en runda är lång,
och knapparna på en pekskärm står under loggen. Förslaget samlar tre saker: fröet
som ett eget avsnitt på värdsidan, spelskärmen över hela bredden med
loggen på kartan, och hur spelledaren når värdsidan under vandringen
utan att räknas som frånvarande. Varje öppen fråga har ett svar nedan,
och sist står en ordning att bygga i.

## Målet

Kartan ska få hela skärmen, och loggen ska läsas utan att skymma den.
Spelledaren ska kunna se, skriva in och slumpa fröet, och titta på
värdsidan mitt i en vandring utan att vandraren står still.

## 1. Fröet på värdsidan

**I dag** väljs fröet först när vandringen börjar (`seedWandering` i
`game.go`). Det kommer från `-fro` för den första vandringen och slumpas
sedan. Värdsidan visar det som en bit av raden, bara under vandringen och
i slutet.

**Förslaget** är ett eget avsnitt, **Fröet**, i högra kolumnen efter
Inför nästa vandring. Avsnittens ordning blir Bjud in, Vid porten, Inför
nästa vandring, Fröet, Kommer någon inte in?, Läs om innehållet och Stäng
tornet.

- **Den här vandringen** (under vandringen och i slutet): `Frö: 123456789`
  på en egen rad, med en knapp som kopierar det. I slutet finns också
  **Samma igen**, som ger nästa vandring samma frö.
- **Nästa vandring** (alltid): ett fält med fröet som nästa vandring får.
  Tornet slumpar det i förväg, så det står alltid ett tal där.
  - **Slumpa** drar ett nytt.
  - Ett eget frö skrivs i fältet: bara siffror, 1–999 999 999. Annat ger
    ett fel under fältet.
  - Fröet gäller nästa vandring. Vandringen efter den får ett nytt slumpat.
  - `-fro` fyller fältet för den första vandringen.
  - En rad under fältet säger att samma frö ger samma torn bara med lika
    många vandrare, samma svårighet och samma drag.
- **Raden** säger inte längre fröet, bara fasen och knappen som behövs.

Servern drar nästa frö när lobbyn öppnas och när en vandring tar slut,
och `seedWandering` tar det. `hostView` får `nextSeed`, och ett nytt anrop
`POST /spelledare/fro` tar ett tal eller `slumpa`. Det fungerar likadant
på Pages genom arbetaren.

| Fråga | Svar |
| --- | --- |
| Eget avsnitt eller en rad i Inför nästa vandring? | Eget avsnitt, så att det syns och kan visa den här vandringens frö också. |
| Färdigt tal eller tomt fält? | Alltid ett färdigt tal, så att man kan skriva upp det innan man startar. |
| Står fröet kvar i raden? | Nej, raden ska bara säga vad som behövs nu. |
| Samma igen i slutet? | Ja, det är det vanligaste skälet att röra fröet. |

## 2. Spelskärmen över hela bredden

**I dag** står kartan till vänster och en kolumn på 300–440 punkter till
höger med loggen (och knapparna på en pekskärm).

**Förslaget** gäller en skärm bredare än 700 punkter. Telefonen behåller
sin enda kolumn med loggen under kartan.

- **Kartan fyller hela bredden** under översta raden. Panelen Nu står
  kvar som en rad under kartan, över hela bredden.
- **Loggen står på kartan, nere till vänster** (byggt så; först tänkt
  uppe), cirka 40 tecken bred.
  Höger sida har redan tecknen, ormkortet och zoomen, och kortet står nere
  i mitten.
- **Ingen ruta.** Texten är ljus med en mjuk, mörk skugga runt varje
  bokstav, och kartan syns mellan raderna. Räcker inte skuggan på ljusa
  golv läggs en svag mörk toning längst till vänster på kartan, aldrig en
  ruta.
- **Den senaste rundan står i full styrka.** När en ny runda börjar tonas
  den förra ner till svag text under den nya, och rundan före den
  försvinner. Rumsblocket står överst i rundan som i dag.
- **Texten tar inte emot musen**, så att allt på kartan går att klicka
  genom den. Bara knappen till lådan går att klicka.
- **Hela loggen** öppnas med knappen **Loggen** under texten eller med
  tangenten L. Den är en halvgenomskinlig låda från vänster, med rundorna
  under rubrikerna som i dag, och stängs med L, Esc eller ett klick
  utanför. Knappen Nya rader följer med in i lådan.
- **Kameran flyttas lite åt höger**, så att rummet man står i hamnar
  bredvid loggen och inte under den (`karta.ts`).
- **Surfplattan:** knapparna från högerkolumnen står nere till vänster på
  kartan, under loggen, och loggen blir kortare för att ge plats.

**Telefonen:**

- **På höjden** (smalare än 700 punkter) blir allt som i dag: översta
  raden fast högst upp, kartan, raden Nu, loggen under kartan och
  knapparna fast längst ner. Loggen får inte plats på kartan där; 40
  tecken är hela kartans bredd.
- **På bredden** är en telefon bredare än 700 punkter men bara 350–430
  hög, och får i dag samma två kolumner som en dator. Med förslaget får
  den kartan över hela bredden, men på en låg skärm (lägre än 500 punkter)
  blir loggen på kartan smalare (cirka 30 tecken), rumsbeskrivningen kortas
  till tre rader med "mer" som på telefonen på höjden, och bara den
  senaste rundan står där, utan den förra. Lådan med hela loggen finns
  som på datorn.
- **Knapparna** står på telefonen på bredden nere till vänster på kartan,
  som på surfplattan, och panelen Nu blir en enda rad.
- **Värdsidan** på en telefon (tornet från Pages) har redan en kolumn, och
  Fröet står där i avsnittens ordning. Märket ändras inte, eftersom Pages
  inte kan öppna värdsidan i en ny flik.

| Fråga | Svar |
| --- | --- |
| Uppe till vänster eller nere till vänster? | Nere, med den senaste rundan nederst, så att de senaste raderna alltid syns; knapparna på en surfplatta står under loggen. |
| Tonar den förra rundan bort? | Den står kvar svagt en runda och försvinner sedan. |
| Behövs lådan med hela loggen? | Ja, de äldre rundorna ska gå att läsa under vandringen, inte bara i krönikan. |

## 3. Värdsidan under vandringen

**I dag** finns ingen knapp till värdsidan under vandringen, eftersom den
som lämnar spelarsidan blir frånvarande. Går spelledaren dit ändå står
vandraren still, spelet fortsätter utan hen, och AI:n tar över efter så
många rundor som Frånvarande vandrare säger. Till spelet tar tillbaka
vandraren.

**Förslaget:** på tornets dator blir märket **spelledare** en knapp också
under vandringen. Den öppnar värdsidan i en **ny flik**, så att
spelarsidan är ansluten hela tiden och vandraren inte blir frånvarande.
I lobbyn och slutet gör knappen som i dag och byter sida i samma flik.

På Pages visar skalet värdsidan och spelarsidan i samma ram, och en ny
flik kan inte visa värdsidan. Där står märket kvar som ett ord under
vandringen.

## Ordningen

Varje steg blir en commit, och sedan stannar jag så att det kan provas.

1. **Fröet i servern:** nästa frö i förväg, `nextSeed`, `POST
   /spelledare/fro`, tester.
2. **Avsnittet Fröet** på värdsidan, texterna på båda språken
   (`data/vardsidan.json`, `data/en/vardsidan.json`), raden utan frö.
3. **Kartan över hela bredden** och loggen på kartan med skuggad text och
   den förra rundan svag.
4. **Lådan med hela loggen** (knappen och L) och kameran som flyttas åt
   höger. Bygg om `web/iso/iso.js`.
5. **Surfplattan och telefonen på bredden** med knapparna på kartan och den lägre loggen, och **märket** som öppnar
   värdsidan i en ny flik under vandringen.
6. **Dokumentationen:** README, CLAUDE.md (värdsidans avsnitt och
   spelskärmens delar), `docs/andringar.md` och det här dokumentet.

Steg 1–2 och 3–5 hänger inte ihop och kan byggas i vilken ordning som
helst.
