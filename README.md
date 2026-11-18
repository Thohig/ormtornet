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

## Spela

1. Klicka på "Gå in i tornet" (det startar ljudet).
2. Skriv ditt namn och tryck Enter.
3. Gå med piltangenterna.

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
namn. Båda spelarna ska synas på kartan och röra sig i realtid.
