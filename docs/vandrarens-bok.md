# Vandrarens bok

En hjälpsida på spelarsidan med det man behöver veta om tornet: grunderna,
klasserna och, efter hand, föremålen, varelserna, fällorna, händelserna,
kombinationerna och ormen. Föreslagen och vald 2026-11-25 (variant B,
upptäckt).

## Tanken

Boken fylls i när man möter saker i tornet. Det man inte har mött står som
ett spelkort med baksidan upp och "Ännu inte sett". Det skyddar tre saker
i spelet:

- **namnen**: rum, fällor och varelser har inga namn för spelarna, så
  boken visar den korta beskrivningen och aldrig `namn` ur `data/`,
- **hemligheten**: inget om ormens namn, namndelarna, avslöjandets nivåer
  eller att tornet förändras,
- **att upptäcka**: kombinationerna, föremålen med två sidor och
  förbannelserna lär man sig genom att bära dem.

Inga råa siffror som vikt eller chans står i boken.

## Hur den ser ut

- Tangenten B, raden B i styrningen bakom ?, knappen Boken i översta
  raden på en pekskärm (där styrningen inte syns) och en knapp i introts
  hörn. Den öppnas i ett lager över sidan och stängs med Esc, B, ✕ eller
  ett klick utanför; vänster och höger byter flik. Medan den är öppen
  blir inga tangenter drag.
- Flikar överst. Grunder och Klasser är öppna från början.
- Innehållet kommer från tornet (`GET /bok`, `bok.go`) på tornets språk,
  och grundernas texter ur `data/spelarsidan.json` (`bok_*`).
- Skriptet är `web/bok.js`.

## Flikarna

| Flik | Innehåll | Öppet från början? |
|---|---|---|
| Grunder | målet, rundan, solen, rummen, striden, de andra vandrarna, ormen, föremålen, slutet | ja |
| Klasser | klasskorten | ja |
| Föremål | de man kan hitta | nej |
| Kombinationer | de hopfogade föremålen | nej |
| Varelser | den korta beskrivningen och hur de verkar i ord | nej |
| Fällor | fällorna och spåren som varnar | nej |
| Händelser | tornets händelser | nej |
| Ormen | sömnmätaren, vakna-korten, andetagen | delvis |

## Stegen

1. Boken med Grunder och Klasser, knappen B. **Byggt 2026-11-25.**
2. Tornet säger vad man har mött och webbläsaren sparar det
   (`localStorage`), och flikarna Föremål, Varelser och Fällor.
3. Händelser, Kombinationer och Ormen, engelska, prov och README.
