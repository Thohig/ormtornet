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
- `-data data/texter.json` fil med rum och texter
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
| Enter      | redo / spela igen                 |

Solen sjunker en handsbredd per runda och visas som en mätare överst.
Gå ut genom en port i något av hörnen för att lämna tornet. Den som kommer
ut med skatt räknas, och den med mest skatt vinner. Den som är kvar inne när
solen har gått ner är förlorad.

Andras skatt och val av drag syns inte förrän spelet är slut.

Om du laddar om sidan och skriver samma namn igen får du tillbaka din vandrare.

## Ändra innehåll

Alla rum och texter ligger i `data/texter.json`. Filen läses när servern
startar, så starta om servern efter en ändring. I texterna kan du använda
`{namn}`, `{horn}`, `{riktning}`, `{region}` och `{rum}` där de passar.
Beskrivningarna för rummen ligger under nivå `"1"`; fler nivåer kommer senare.

## Struktur

- `main.go` startar webbservern och skriver ut adressen
- `game.go` all spellogik (tornet, rummen, spelarna)
- `data.go` läser in texterna
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
