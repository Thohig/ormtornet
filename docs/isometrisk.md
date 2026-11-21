# Förslag: Tornet snett ovanifrån

Status: förslaget skrevs 2026-11-21. Skisserna visade två alternativ,
platta kritfigurer (A) och 3D-modeller (B). Användaren valde B. Steg 1 är
byggt.

Kartan ritas i dag med Canvas 2D i `web/map.js`, rakt ovanifrån, med tecken
i ett rutnät. Förslaget är en ny karta i isometrisk vy, ritad med Three.js
och WebGL (`THREE.WebGLRenderer`), skriven i TypeScript och byggd med Vite.
Vandrarna och varelserna blir enkla 3D-modeller med ljus och skuggor.

## Så ser det ut

- Kameran tittar snett uppifrån med en ortografisk kamera och vrids aldrig.
  Norr pekar upp åt höger, och en kompass står i hörnet. Piltangenterna
  betyder alltid samma sak.
- Golven är skifferplattor. Ett rum ingen har varit i är en svag platta.
- Väggarna är låga, så att ingen figur göms bakom en vägg. Dörrarna är
  luckor mitt på väggen, och portarna har en båge.
- Fällornas märken, kritpilarna och föremålen står i krita på golvet.
  Kammaren har en kritring och en hög av guld, och ormens ögon lyser när
  den vaknar.
- Vandrarna är low-poly-figurer i sin färg, med det klassen bär: skölden,
  kartan, staven. De vänder sig och går mellan rummen i stället för att
  hoppa. Den egna har en ring på golvet, och draget är en pil vid dörren.
- Varelserna har en modell var. En varelse i `data/` som saknar modell
  visas med sitt tecken.
- Avslöjandet: kritan andas på nivå 2, revben och hjärtslag på nivå 3, och
  på nivå 4 blir väggarna kött och portarna munnar med tänder.

## Det som är bestämt

- Den platta kartan finns kvar. Spelledaren väljer på värdsidan om kartan
  är platt eller snett ovanifrån, och den platta används av sig själv där
  WebGL saknas.
- Spelledarbordet får samma val, eftersom det använder samma `MapView`.
- Den nya kartan har samma gränssnitt som `MapView` (`setup`, `draw`,
  `zoom` och så vidare), så att `app.js` och `bord.js` knappt ändras.
- Teckenförklaringen finns kvar.

## Hur koden byggs

- Källkoden ligger i `klient/`, med TypeScript i `klient/src/`. Vite bygger
  den till ett enda skript, `web/iso/iso.js`, som laddas som de andra
  skripten och nås genom den globala `IsoKarta`.
- Det byggda skriptet checkas in. Då bygger `go build`, `starta.sh` och
  Ormtornet.app spelet utan Node, och Three.js följer med i programmet, så
  tornet fungerar utan internet.
- Node behövs bara för den som ändrar i `klient/`. `starta.sh` bygger om
  skriptet när källkoden är nyare och Node finns.

## Stegen

1. Vite och TypeScript i `klient/`, med bygget till `web/iso/` och
   kopplingen i `starta.sh`. Ingenting ändras på skärmen. **Byggt.**
2. Den isometriska kartan: golv, väggar, dörrar, portar, kammaren och
   kritan på golvet, med enkla figurer för vandrarna och varelserna och
   samma gränssnitt som `MapView`. Valet platt eller snett ovanifrån på
   värdsidan.
3. Modellerna: vandrarna med klassernas föremål, varelserna, guldet, och
   rörelsen mellan rummen. Ekot och de dolda vandrarnas ringar.
4. Avslöjandet: nivåerna 2–4, ormens ögon och munnarna.

Efter varje steg görs en commit, och användaren provar innan nästa steg.
