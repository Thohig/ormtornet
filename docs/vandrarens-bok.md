# Vandrarens bok

En hjälpsida på spelarsidan med det man behöver veta om tornet: grunderna,
klasserna, föremålen, kombinationerna, varelserna, fällorna, händelserna
och ormens kort. Föreslagen och vald 2026-11-25 (variant B, upptäckt) och
byggd samma dag i tre steg. Samma dag öppnades allt från början, med
märket Mött på det man har mött.

## Tanken

Boken visar allt från början. Det man har mött i tornet får märket
**Mött**, så att det finns något att samla. Två saker står ändå inte i
boken:

- **namnen**: rum, fällor och varelser har inga namn för spelarna, så
  boken visar den korta beskrivningen och aldrig `namn` ur `data/`,
- **hemligheten**: att ormens namn ska viskas till tornet, och vad det ger
  att bli invigd. Boken säger att namnet finns, i tre delar, och var
  delarna går att hitta, men inte vad man gör med det.

Tills 2026-11-25 fylldes boken i när man mötte saker, och det man inte
hade mött stod med baksidan upp och "Ännu inte sett". Förbannelserna,
vändningarna och kombinationernas delar syntes först när man hade varit
med om dem.

Inga siffror för chanser, vikter eller varelsernas Kraft och Försvar står
i boken. Det som är siffror i spelet för vandraren själv, som skada,
rundor och klassernas Kraft och Försvar, står kvar.

Varianterna som valdes bort först var A, en uppslagsbok med allt öppet,
och C, boken med allt öppet för spelledaren på värdsidan. Boken är nu
närmast A, med märket Mött från B.

## Att öppna och stänga

- Tangenten **B**, var man än är på spelarsidan, utom när man skriver i
  en ruta (viskningen).
- Raden **B Vandrarens bok** i styrningen bakom ?.
- Knappen **Boken** i översta raden, bara på en pekskärm, där styrningen
  inte syns.
- Knappen **Vandrarens bok** i introts nedre högra hörn. Den tonas bort
  när man går in genom porten.

Boken är ett lager över hela sidan (`#bok`, z-index 48, under
kortvisningen). Den stängs med B, Esc, ✕ eller ett klick på det mörka
utanför sidan. Vänster och höger pil byter flik; upp och ner rullar.
Medan boken är öppen blir inga tangenter drag: `web/bok.js` läses före
`app.js`, så att dess lyssnare kommer först, och stoppar tangenterna med
`stopImmediatePropagation`. När boken öppnas får sidan i den fokus, inte
✕, så att Enter varken stänger den eller gör en redo.

På telefonen fyller boken hela skärmen. Klasserna står i en kolumn med
beskrivningen kvar (klassvalet döljer den på telefonen), och de andra
korten två i bredd.

## Flikarna

| Flik | Innehåll |
|---|---|
| Grunder | målet, rundan, solen och portarna, rummen, striden, de andra vandrarna, skatten och ormen, ormens namn, föremålen, varslen, tornet förändras, slutet |
| Klasser | alla klasskort med framsidan upp |
| Föremål | de föremål som går att hitta, och brevet |
| Kombinationer | de hopfogade föremålen och vad de är fogade av |
| Varelser | den korta beskrivningen, texten när den dyker upp och hur den är i ord |
| Fällor | den korta beskrivningen, vad den gör och spåret efteråt |
| Händelser | tornets händelser och vad de gör |
| Ormen | ormens kort och vad de gör, under en inledning om leken |

Korten står i samma ordning som i `data/`. Flikarna har ingen räknare.

Grundernas avsnitt **Ormens namn** säger att namnet har tre delar, att
var och en hittar dem själv (ristningen, varelsen, kammaren) och att varje
del säger var på ormen den sitter. **Varslen** säger att tornet förråder
rum som ingen har varit i när solen har sjunkit en bit, och att inget
varsel säger var kammaren ligger. **Tornet förändras** säger att fällorna
i nya rum ändras när solen står lågt, och att portarna sluter sig en i
taget mot slutet, med en varning rundan innan och den sista öppen.

## Vad korten visar

Alla kort har märket Mött under rubriken när man har mött det.

- **Föremål**: namnet, bilden och beskrivningen, förbannelsen (den första
  texten under `forbannelse`) och det som händer när föremålet vänder
  (den första under `vander`).
- **Kombinationer**: namnet, bilden, beskrivningen och "Fogat av …" med
  delarnas namn.
- **Varelser**: `kort` med stor bokstav, den första texten under
  `dyker_upp` och `bok`.
- **Fällor**: `kort` med stor bokstav, `bok` och "Efteråt:" med den första
  texten under `spar_efter`.
- **Händelser**: `rubrik` och `bok`.
- **Ormen**: kortets `namn` och `bok`.

## Det man har mött

Det man har mött ger märket Mött. Tornet skriver in det vandraren möter (`learn` i `bok.go`) och skickar det
i varje tillstånd som `bok`, en sorterad lista med nycklar som
`varelse:krypare`. Listan börjar om för varje vandring (`placeAtCorner`).

| Nyckel | När | Var |
|---|---|---|
| `foremal:<id>` | man får föremålet i väskan, eller läser ett brev | `give` i `items.go` |
| `forbannelse:<id>` | föremålets förbannelse visar sig | `curse` i `items.go` |
| `vander:<id>` | föremålet vänder | `turnItem` i `items.go` |
| `kombination:<id>` | två föremål fogas ihop | `combine` i `items.go` |
| `varelse:<id>` | en varelse syns i ens rum, inte i mörkret | `announceMonster` i `monster.go` |
| `falla:<id>` | en fälla slår till, amuletten tar den, eller Spårfinnaren ser den i tid | `trigger` i `events.go`, och `game.go` |
| `handelse:<id>` | en händelse ger en ett spelkort | `tornhandelser.go` |
| `ormkort:<id>` | man drar ett kort som ger guld | `take` i `serpent.go` |
| `ormkort:<vakna-kortet>` | man ser uppvaknandet, i kammaren eller utanför | `wakeSerpent` i `serpent.go` |

Förbannelserna och vändningarna skrivs fortfarande in men visas inte för
sig; märket gäller föremålet. Spelarsidan lägger ihop listan med det den
redan har (`Bok.merge`, som `app.js` anropar när tillståndet kommer) och
sparar allt i `localStorage` under `ormtornet-bok`. Märkena växer alltså
från vandring till vandring i samma webbläsare, oavsett torn och namn. Den som byter webbläsare eller
tömmer den börjar om. På Pages gäller spelarens egen webbläsare, inte
skalets.

## Texterna

- **Sidans texter** i `data/spelarsidan.json` och
  `data/en/spelarsidan.json`, med nycklar som börjar med `bok_`:
  rubriken, knapparna, flikarna (`bok_flik_*`), märket (`bok_mott`),
  flikarnas inledningar (`bok_*_inledning`) och
  raderna för förbannelse, vändning, spår och delar. Grundernas avsnitt
  står i ordningen under `bok_grunder`; varje id har `bok_<id>_rubrik`
  och `bok_<id>`. Raden i styrningen är `hjalp_bok` och introts knapp
  `intro_bok`.
- **Innehållets texter** under `bok` i `data/monster.json`,
  `data/fallor.json`, `data/handelser.json` och `data/ormen.json`, och
  översatta i `data/en/`. Tornet vägrar starta om någon saknas
  (`validate.go`). Siffror som skada och rundor är skrivna i texten och
  ska ändras där när fälten ändras.
- Varelsernas `bok` beskriver dem i ord utifrån värdena: `lynne` under 30
  är "anfaller sällan", 40–50 "ibland" eller "ofta", 75 och mer "nästan
  alltid"; `vandrar` 0,3 och mer är "vandrar gärna" och under 0,1
  "lämnar nästan aldrig sitt rum"; liv, skada och guld blir "faller
  lätt", "slår hårt" och "bär mest guld". Ändras värdena ska texten följa
  med.

## Koden

- `bok.go`: `Book` bygger svaret på `GET /bok` (`bookView`, med
  `classView` för klasserna och `bookCard` för resten), `learn` och
  `learned` håller det vandraren har mött, och `wakeCardID` säger vilket
  ormkort som väcker ormen. Adressen läggs till i `addHostRoutes` i
  `host.go` och fungerar därför både från tornets program och på Pages.
- `web/bok.js`: modulen `Bok` med `open`, `close`, `toggle`, `oppen` och
  `merge`. Innehållet hämtas varje gång boken öppnas, så att det följer
  tornets språk och innehåll.
- `web/index.html`: lagret `#bok`, knapparna `#bokBtn` och `#introBok` och
  raden `#helpBok`. `web/app.js` kopplar knapparna och anropar
  `Bok.merge`.
- `web/style.css`: allt under `#bok`, `#bokSida`, `#introBok` och
  `#bokBtn`. Klasskorten återanvänder `.card`, de andra korten `.okort`
  med framsidan upp, och märket är `.bok-mott`.
- Filerna i `web/` byggs in i programmet, så tornet måste startas om för
  att visa en ändring där.

## Proven

- `bok_test.go`:
  - `TestBookHasAllClasses`: `GET /bok` ger alla klasser.
  - `TestBookBasicsTexts`: Grundernas avsnitt finns på båda språken.
  - `TestBookLearnsWhatYouMeet`: föremål, förbannelse, fälla och varelse
    skrivs in, och inget man inte har mött.
  - `TestBookHasNoNames`: inga namn på varelser eller fällor i boken, och
    inga hopfogade föremål bland föremålen.
  - `TestBookLearnsEventsCombosAndSerpent`: kombinationer och ormkort
    skrivs in, och boken har alla händelser, ormkort och kombinationer.
- Testet för engelska kräver att alla nya texter är översatta.
- Webbläsarprovet (`verktyg/webbprov.mjs`) gick igenom efter varje steg.
- Boken provades för hand i webbläsaren på dator och telefon och på
  engelska. De fyllda flikarna provades med mötta saker inlagda i
  `localStorage`. Två fel hittades och rättades: ett klick på en flik
  stängde boken (fliken ritades om innan klicket nådde lagret), och Enter
  tryckte på ✕.

## Stegen

1. Boken med Grunder och Klasser, och sätten att öppna den. Byggt
   2026-11-25.
2. Tornet säger vad man har mött och webbläsaren sparar det, och flikarna
   Föremål, Varelser och Fällor. Byggt 2026-11-25.
3. Kombinationer, Händelser och Ormen, och allt på engelska. Byggt
   2026-11-25.
4. Allt öppet från början, med märket Mött, och avsnitten Ormens namn,
   Varslen och Tornet förändras i Grunder. Byggt 2026-11-25, samtidigt
   som viskningen blev bara namnet (docs/avslojandet.md).

## Att prova

- Att boken fylls i under en riktig vandring, särskilt händelserna och
  uppvaknandet, som bara testerna och inlagda nycklar har provat.
- Om raden med åtta flikar, som blir två rader på datorn, är för lång.
- Var knappen Boken står i översta raden på telefonen.
- Om det blir för lätt när allt står öppet, och om märket Mött räcker
  för att det ska kännas som att samla.
- Om någon listar ut att namnet ska viskas, när boken inte säger det.

## Inte byggt

- Ett klick på ett kort i boken visar det inte stort.
- Förbannelserna hos de hopfogade föremålen står inte på deras kort.
- Det man har mött följer webbläsaren, inte vandraren; det finns inget
  sätt att flytta det eller börja om, utom att tömma webbläsaren.
- Spelledaren har ingen egen version av boken med allt öppet (variant C).
