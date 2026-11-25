# Samlat förslag: mobilvyn

Förslag 2026-11-23. Steg 1–5 byggdes samma dag.
Avsnittet Byggt säger var bygget skiljer sig från förslaget.

Sedan 2026-11-24 finns inte Parera och Undan, och i strid står Vänta,
Anfall, Sök och Drick ([bara-anfall.md](bara-anfall.md)).

## Målet

På telefonen ska kartan vara det man ser. Rummet man står i ska fylla
skärmen, och loggen, tecknen och knapparna ska ta så lite plats som
möjligt utan att något försvinner.

## Läget i dag

På en telefon (smalare än 700 px) står, uppifrån och ned:

- översta raden, fast överst,
- kartan, nästan kvadratisk (`aspect-ratio: 0.95`), med tecknen uppe till
  höger (upp till 58 % av bredden),
- panelen Nu,
- loggen, där rundan i bild tar minst 60 % av skärmens höjd,
- knapparna, fasta längst ned, omkring 170 px höga.

Nära visar 3 rum och kameran kortar av till omkring 0,6 av det
(`skarp` i `rita`, `klient/src/karta.ts`). Det mesta av platsen går alltså
till det som ligger runt kartan, inte till zoomen.

## Det nya

Uppifrån och ned på telefonen:

1. **Översta raden** som i dag, plus **tecknen** på en egen tunn rad
   under den (steg 3). Båda står fast överst.
2. **Kartan** fyller höjden mellan tecknen och knapparna.
3. **Nu** som en tunn rad längst ned på kartan, ovanpå den, bara när den
   har något att säga.
4. **Knapparna** som en rad längst ned, omkring 56 px.
5. **Loggen** under kartan, när man rullar ned: rumsblocket och rundans
   minikort, och hela loggen i lådan.

## Stegen

Varje steg är en commit och går att prova för sig.

### Steg 1: kartan fyller skärmen

Bara CSS. `#mapbox` får höjden
`calc(100svh - översta raden - knapparna)` i stället för
`aspect-ratio`, och loggen hamnar under den. Kartan blir 40–60 % högre
på en vanlig telefon. Kameran räknar redan ut utsnittet efter kartans
mått, så rummet blir större av sig självt.

### Steg 2: en tunnare logg

Bara CSS och lite i `web/app.js`.

- Den första rundan tar inte längre minst 60 % av höjden.
- Rumsblocket visar 2 rader i stället för 3; "mer" fäller ut resten som i
  dag.
- Rundans rader står som minikort på en rad med en kort text. Ett tryck
  visar kortet som i dag.
- Knappen Loggen öppnar lådan med hela loggen, som på datorn.

### Steg 3: tecknen på en rad

Tecknen flyttas från kartans hörn till en tunn rad under översta raden:
klasskortet, liv, skatt, drycker, tillstånd och väskans platser, som
symboler med tal ("♥ 7") men utan ord. Ett tryck visar förklaringen som i
dag. Raden rullar i sidled om den inte får plats, och bryts aldrig.
Kartan blir då helt fri.

### Steg 4: tryck på dörrarna

Sedan 2026-11-24 hålls tangenterna inom den del av kartan som syns:
en tangent som skulle hamna under översta raden, bakom knapparna eller
utanför kanten flyttas in till kanten åt samma håll, och två som hamnar
vid samma kant skjuts isär. På telefonen ligger de över det stora kortet,
så att de alltid går att trycka på (`ritaTangenter` i
`klient/src/karta.ts`, `freeBand` i `web/map.js`).

Tangenterna vid dörrarna blir knappar: ett tryck på en tangent väljer
det hållet, precis som piltangenten. Det gäller på alla skärmar, också
med musen på datorn, eftersom det inte ändrar något som finns där i dag.

I mörkret syns inga dörrar. Sedan 2026-11-25 står då en streckad tangent
vid varje vägg (`trevar` i `ritaTangenter`), så att det alltid finns
något att trycka på. En vägg utan dörr ger en stöt i muren. Pilknapparna
som visades i mörkret samma morgon togs bort igen.

På telefonen försvinner pilknapparna. Kvar blir en rad med de knappar
som gäller just nu: **Vänta** (som pricken i mitten i dag), Sök, Drick,
Anfall, Parera, Undan, Viska och Redo. De följer läget som i dag
(`renderPad`). Tangenterna får minst 40 px att trycka på när rummet är
stort på skärmen; när de krymper (utzoomat) gäller raden ovan och
pilarna kommer tillbaka.

Kräver en ändring i `klient/src/karta.ts` och att `web/iso/iso.js` byggs
om.

### Steg 5: Rummet (om det behövs)

En nivå närmare än Nära, bara på telefon, där ens rum fyller bredden och
grannrummens kanter syns. Den blir telefonens nya standard; den som
redan har valt en zoom behåller den. Två fingrar når den som i dag.
Kameran ställer ens rum lite under mitten så att översta raden inte
skuggar det. Bygg det bara om steg 1–4 inte räcker.

## Vad som påverkas

| Steg | Var | Datorn |
| --- | --- | --- |
| 1, 2, 3 | CSS i `@media (max-width: 700px)` och `web/app.js` | oförändrad |
| 4 | `klient/src/karta.ts`, CSS | kan klicka på tangenterna; inget annat ändras |
| 5 | `web/map.js`, `klient/src/karta.ts`, bara på smal skärm | oförändrad |

Gränsen 700 px gäller också ett smalt datorfönster. En iPad på höjden
(768 px) får datorns vy. Värdsidan, `data/` och reglerna berörs inte;
de nya texterna (Vänta, namnet på nivån Rummet) går in i både
`data/spelarsidan.json` och `data/en/`.

## Frågor

| Fråga | Förslag |
| --- | --- |
| Tecknen i översta raden eller på en egen rad? | En egen rad: översta raden är redan full. |
| Ska tangenterna gå att klicka på datorn? | Ja, det kostar inget. |
| Var står Vänta när pilarna är borta? | Först i knapparnas rad. |
| Steg 5 nu eller senare? | Byggt samma dag som 1–4. |

## Byggt

- **Steg 1**: en högre karta visade inte rummet större av sig själv,
  eftersom kameran passar utsnittet efter bredden. En hög, smal karta går
  därför närmare när den följer någon (`rita` i `klient/src/karta.ts`).
  Raden Nu står kvar under kartan i stället för på den, och kartan mäter
  den (`measureBars` i `web/app.js`), så att den syns hel ovanför
  knapparna också i strid.
- **Steg 2**: rundans rader står som text, tätare, inte som en rad
  minikort med kortare text. Rundans rubrik syns inte under kartan.
- **Steg 3**: tecknen står på en egen rad i översta raden (`placeSide`
  flyttar `#hud`).
- **Steg 4**: tangenterna gick redan att trycka på en pekskärm, så
  ändringen i `karta.ts` behövdes inte. Pilknapparna kommer tillbaka på
  nivåerna Långt och Allt (`langtut`), inte efter tangenternas storlek.
- **Steg 5**: Rummet har nummer -1 bland zoomnivåerna (`ROOM` i
  `web/map.js`), så att en sparad nivå betyder samma sak som förut. Den
  visar 2 × 2 rum. Kameran flyttar inte ned rummet under mitten, eftersom
  tecknen inte längre står på kartan. Blir skärmen bred går kartan till
  Nära.
