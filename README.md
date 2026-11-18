# Ormtornet

En dungeon crawler för 2–4 spelare i webbläsaren, på ett lokalt nätverk.
Mitt i ett torn på 9x9 rum sägs en jättelik orm sova kring en skatt.
Ta så mycket du vågar och kom ut innan solen går ner.

## Starta

Kräver Go (`brew install go`).

```bash
go run .
```

Servern lyssnar på `0.0.0.0:8080` och skriver ut adressen spelarna ska
öppna, till exempel `http://Din-Mac.local:8080`. Alla Macar i samma nätverk
kan öppna den i Safari eller Chrome.

Flaggor:

- `-addr 0.0.0.0:8080` adress och port
- `-data data` mapp med datafilerna
- `-sol 40` antal rundor innan solen går ner
- `-turtid 30s` längsta tid att välja drag per runda
- `-min 2` minsta antal spelare för att börja (sätt `-min 1` för att prova ensam)

## Spela

1. Klicka på "Gå in i tornet" (det startar ljudet).
2. Skriv ditt namn och tryck Enter.
3. Tryck Enter igen när du är redo. Vandringen börjar när alla är redo.

Spelet går i rundor. Alla väljer sitt drag samtidigt, och när alla har valt
(eller turtiden tar slut) utförs alla drag på en gång. Du kan ändra dig tills
dess. Den som inte hinner välja blir stående.

| Tangent    | Drag                              |
|------------|-----------------------------------|
| Pilar      | gå åt det hållet                  |
| Mellanslag | vänta                             |
| T          | ta skatt (bara i skattkammaren)   |
| S          | sök igenom rummet                 |
| D          | drick en dryck                    |
| Enter      | redo / spela igen                 |

Solen sjunker en handsbredd per runda och visas som en mätare överst.
Gå ut genom en port i något av hörnen för att lämna tornet. Den som kommer
ut med skatt räknas, och den med mest skatt vinner. Den som är kvar inne när
solen har gått ner är förlorad.

Andras skatt och val av drag syns inte förrän spelet är slut.

### Skattkammaren och ormen

I kammaren mitt i tornet kan du ta skatt (T) en gång per runda. Varje gång
dras ett ormkort ur en lek:

- **Sover** eller **rör sig**: du får 15–45 guld och kortet tas bort ur leken.
  Rör sig ormen märker alla i tornet det.
- **Vaknar**: kortet blir kvar i leken. Alla i kammaren förlorar hälften av
  sin skatt och 2 liv och kastas ut genom en slumpad öppning. Väsandet hörs
  i hela tornet.

Eftersom bara sovkorten tas bort ökar risken för varje näve. Leken börjar med
13 sovkort och 2 vakna-kort. När ormen har vaknat blandas en ny lek med 2
färre sovkort än förra gången. Står du i kammaren ser du hur djupt ormen
sover och hör andetagen bli tätare.

### Rum och fällor

När någon går in i ett nytt rum kan en fälla slå till. Ett tecken i rummet
visar sedan för alla att något har hänt där.

| Fälla           | Vad som händer                                      |
|-----------------|-----------------------------------------------------|
| Fallgaller      | du sitter fast i 2 rundor                           |
| Rasande golv    | 1 skada, du faller ner i källarvalven i 2 rundor och kommer upp någon annanstans; ingen ser var du är |
| Roterande rum   | du kastas ut genom en slumpad öppning               |
| Pilar           | 2 skada                                             |
| Giftig dimma    | 1 skada nu och 1 skada per runda i 2 rundor         |
| Ormgrop         | 2 skada och fast i 1 runda                          |
| Ömsat skinn     | fast i 1 runda                                      |
| Klocka          | solen sjunker 2 extra steg för alla                 |
| Mörker          | lyktan slocknar i 4 rundor; du ser bara ditt eget rum |
| Eko             | alla får veta var du är                             |

Du har 6 liv. Den som når noll liv dör.

Ett vanligt rum kan sökas igenom en gång (S). Man kan hitta skatt, en dryck,
en hemlig gång till ett grannrum, eller utlösa en fälla. En dryck (D) ger
tillbaka 2 liv och tar bort gift.

Om du laddar om sidan och skriver samma namn igen får du tillbaka din vandrare.

## Ändra innehåll

Allt innehåll ligger i mappen `data/`. Filerna läses när servern startar,
så starta om servern efter en ändring.

- `texter.json` rum, beskrivningar och allmänna texter
- `fallor.json` regler (liv, chans för händelse, drycker), fällor och sökresultat
- `ormen.json` ormkorten, skatten i kammaren och vad som händer när ormen vaknar

För varje fälla kan du ändra namn, vikt (hur ofta den dyker upp), skada,
antal rundor, tecknet på kartan och texterna. Vad fällan gör styrs av dess
`id` i `events.go`. I texterna kan du använda
`{namn}`, `{horn}`, `{riktning}`, `{region}` och `{rum}` där de passar.
Beskrivningarna för rummen ligger under nivå `"1"`; fler nivåer kommer senare.

## Struktur

- `main.go` startar webbservern och skriver ut adressen
- `game.go` all spellogik (tornet, rummen, spelarna)
- `events.go` fällor, sökande och drycker
- `serpent.go` skattkammaren och ormkorten
- `data.go` läser in datafilerna
- `ws.go` en liten websocket-implementation utan beroenden
- `web/` klienten: HTML, CSS och JavaScript utan ramverk
  - `map.js` ritar kartan med ASCII- och braille-tecken
  - `audio.js` allt ljud, skapat med Web Audio API

## Testa lokalt

Starta servern och öppna `http://localhost:8080` i två flikar med olika
namn. För ett snabbt test passar `go run . -sol 8 -turtid 10s`.

Spellogiken har tester:

```bash
go test ./...
```
