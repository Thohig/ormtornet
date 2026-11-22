# Nytt gränssnitt på dator och mobil (förslag)

Förslag 2026-11-22. Steg 1–4 är byggda; resten är inte bestämt.

## Vad som önskas

- Ett nytt gränssnitt för spelskärmen, både på dator och på telefon.
- Teckenförklaringen tas bort från kartan och flyttas in i krönikan.
- Krönikan beskriver rummen.
- Krönikan visar bara den aktuella rundan. Tidigare rundor ser man bara om
  man själv rullar ner.
- Panelen Vandrare tas bort.
- Raden "Tryck enter när du är redo" tas bort.

## Ordet krönikan

I dag är krönikan skärmen efter döden och i slutet (`chronscreen`), och
**loggen** är raderna i högerkolumnen under spelet. Förslaget gör om
loggen till **krönikan**: en löpande berättelse under vandringen som har
samma tecken som krönikans skärm (`signs.go`). Skärmen efter döden och i
slutet behåller sina namn, Din vandring och krönikan, och får samma
rumsbeskrivningar som krönikan under spelet (se steg 3).

## Krönikan under vandringen

Krönikan står i högerkolumnen och har bara **den aktuella rundan** i bild:

```
  Runda 7
  ┌ Valvgången ───────────────────────────┐
  │ Valven höjer och sänker sig nästan     │
  │ omärkligt, i samma takt som ditt eget  │
  │ andetag. Fast inte riktigt samma.      │
  │                                        │
  │ I rummet finns en kammare, och Spindeln│
  │ väntar i skuggan. Agnes står vid dörren│
  └────────────────────────────────────────┘
  Spindeln biter dig. −2 liv.
  Du hör steg från rummet i norr.

  ─ tidigare ─────────────────────── ↓
  Runda 6 …
```

- **Rummet överst.** När man går in i ett rum börjar rundans text med
  rummets namn och beskrivningen (`beskrivningar` per nivå, som i dag
  skrivs som en vanlig rad). Efter den kommer en mening i ord om det som
  står i rummet just nu, till exempel "I rummet finns en kammare, och
  Spindeln väntar i skuggan. Agnes står vid dörren." Hur det ser ut syns
  på kartan.
- **Inga tecken eller bilder i krönikan.** Bilderna som kom in på
  krönikans rader 2026-11-22 (`signs.go`) tas bort, både under vandringen
  och i Din vandring och slutets krönika. Allt beskrivs i text.
- **Står man kvar** i rummet upprepas inte beskrivningen. Rundan börjar då
  med namnet och det som har ändrats ("Spindeln är borta. Agnes har kommit
  in.").
- **Teckenförklaringen tas bort** från kartan (`#legend`,
  `renderLegend`). Rummets namn uppe till vänster på kartan står kvar.
- **Bara aktuell runda i bild.** När en ny runda börjar flyttas den förra
  ner under strecket "tidigare", utanför bild. Krönikan står alltid
  överst, och det äldre når man bara genom att rulla ner. Den nuvarande
  ihopfällningen av äldre rundor tas bort; det räcker att de ligger under
  strecket. Har man rullat ner står man kvar där, och "↑ nya rader" tar en
  upp, som i dag.
- Rundans sammanfattning, kortet på kartan, finns kvar. Dess korta rad i
  loggen (`rundkort_logg`) blir onödig när hela rundan står i krönikan och
  tas bort.

## Lobbyn och slutet utan listan och utan rutan

- **Panelen Vandrare tas bort** helt (`#playersBlock`, `renderSide`).
- **Prickarna i översta raden** tar dess plats även i lobbyn och slutet. Där
  står namnet bredvid varje prick (under vandringen bara pricken, som i
  dag). En fylld prick betyder redo, som ett valt drag gör under
  vandringen. Pekar man på en prick syns namn, klass, och om vandraren är
  frånkopplad eller en robot.
- **Raden "… när du är redo" tas bort** (`nu_tryck_redo`, `nu_redo`). I
  stället står en knapp **Redo** i översta raden på dator (Enter fungerar
  som förut), och den visar "Redo ✓" när man har tryckt. I slutet står där
  **Nytt torn**. På telefon finns knappen redan i knapparna längst ner.
- Hjälpen bakom ? har kvar raden om Enter.

## Dator

```
┌ ORMTORNET ☀☀☀☀○○  ● ● ○ ●  runda 7 · 0:42      [Redo] ┐
├───────────────────────────────────┬─────────────────────┤
│ Valvgången          liv ♥♥♥  ? │ KRÖNIKAN            │
│                     skatt ¤ 12   │ Runda 7             │
│                     [väskan]     │ ┌ Valvgången ────┐  │
│                                  │ │ beskrivning     │  │
│            kartan                │ │ vad som står här│  │
│                                  │ └─────────────────┘  │
│                                  │ rader …             │
│                                  │                     │
│  [rundans kort]                  │ ─ tidigare ──── ↓   │
├───────────────────────────────────┴─────────────────────┤
│ Nu: Du anfaller Spindeln.  Dess liv ▮▮▮▯▯  (A P U)      │
└─────────────────────────────────────────────────────────┘
```

- Kartan får mer plats: ungefär två tredjedelar av bredden.
- Högerkolumnen är **bara krönikan**, i hela höjden.
- **Panelen Nu** blir en smal rad under kartan: draget, striden och
  varelsens liv. Sömnmätaren i skattkammaren står där också. Viskningens
  ruta öppnas i samma rad.

## Telefon

```
┌ ☀☀☀○○  ● ● ○ ●  0:42 ┐
│                       │
│        kartan         │
│  (tecknen, väskan)    │
│                       │
├───────────────────────┤
│ Nu: Du anfaller …     │
├───────────────────────┤
│ Runda 7               │
│ ┌ Valvgången ───────┐ │
│ │ beskrivning        │ │
│ │ vad som står här   │ │
│ └────────────────────┘ │
│ rader …               │
├───────────────────────┤
│  ↑      T S D         │
│ ← · →   A P U         │
│  ↓      ⇥  Redo       │
└───────────────────────┘
```

- Översta raden är fast högst upp och knapparna fast längst ner, som i
  dag.
- Kartan står överst och raden Nu direkt under den, så att striden syns
  utan att man rullar.
- Under den kommer krönikan med bara den aktuella rundan. Rummets
  beskrivning kan vara lång på en smal skärm; den kan visas med tre rader
  och "mer" som fäller ut resten.
- I dag delar Nu och Vandrare två kolumner. De försvinner, och allt står i
  en kolumn.

## Steg

1. **Lobbyn och slutet.** Ta bort panelen Vandrare och raden om Enter, och
   lägg in prickar med namn och knappen Redo/Nytt torn i översta raden.
2. **Rummet i krönikan.** Servern skickar ett rumsblock (namn,
   beskrivning, en mening om vad som står där) i stället för beskrivningen som vanlig
   rad. Spelarsidan ritar blocket, och teckenförklaringen på kartan tas
   bort, liksom bilderna på krönikans rader (`signs.go`). Texterna in på
   båda språken.
3. **Bara aktuell runda.** Krönikan visar rundan i bild och resten under
   strecket. Ihopfällningen och rundkortets rad i loggen tas bort.
   Rumsblocken kommer också med i Din vandring och slutets krönika.
4. **Ny layout på dator.** Kartan större, krönikan i hela högerkolumnen och
   Nu som en rad under kartan.
5. **Ny layout på telefon.** En kolumn med kartan, Nu och krönikan, och
   knapparna fast längst ner.

Varje steg blir en egen commit, och användaren provar mellan stegen.

## Frågor

- Ska loggen under spelet heta krönikan, med rumsblocken där, så som
  förslaget tänker? Eller menas skärmen efter döden och i slutet?
- Räcker prickarna med namn i lobbyn, eller ska man också se klassen där?
- Ska beskrivningen upprepas när man står kvar i ett rum, eller räcker
  namnet och det som har ändrats?
- Ska rundans kort på kartan finnas kvar när hela rundan står överst i
  krönikan?
