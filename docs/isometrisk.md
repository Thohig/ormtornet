# Förslag: Tornet snett ovanifrån

Status: förslaget skrevs 2026-11-21. Skisserna visade två alternativ,
platta kritfigurer (A) och 3D-modeller (B). Användaren valde B. Alla fyra
stegen är byggda, och kartan snett ovanifrån är förval sedan 2026-11-21.
Sedan 2026-11-22 är den den enda kartan: den platta kartan i krita och
inställningen Kartan är borttagna. `web/map.js` räknar bara ut utsnittet,
zoomen och teckenförklaringen, och där WebGL saknas visar spelarsidan en
text i stället för kartan. Resten av dokumentet beskriver hur det var när
den platta kartan fanns kvar.

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
   värdsidan. **Byggt.**
3. Modellerna: vandrarna med klassernas föremål, varelserna, guldet, och
   rörelsen mellan rummen. Ekot och de dolda vandrarnas ringar. **Byggt.**
4. Avslöjandet: nivåerna 2–4, ormens ögon och munnarna. **Byggt.**

Efter varje steg görs en commit, och användaren provar innan nästa steg.

## Så är steg 2 byggt

Sedan 2026-11-22 finns inställningen Kartan och den platta kartan inte
längre.

- Inställningen Kartan (`karta` i inställningarna, `-karta`) är `platt`
  eller `snett`. Först var `platt` förval; när alla stegen var byggda blev
  det `snett`. Den följer
  med tillståndet som `karta`, både till spelarna och spelledarbordet.
- `MapView` i `web/map.js` väljer kartan. Den platta (`FlatMap`) räknar
  alltid ut utsnittet, zoomen och teckenförklaringen, men ritar inget när
  kartan snett ovanifrån visas. Den kartan (`klient/src/karta.ts`) ritar
  på en egen canvas bredvid och läser utsnittet och rummet den följer från
  den platta.
- Kameran står över rummet kartan följer och glider efter när vandraren
  byter rum. Utsnittet är en romb på skärmen, och hörnen skärs av så att
  rummen närmast blir större. Utan någon att följa visas hela utsnittet.
- Vandrarna är en kropp och ett huvud i sin färg, den egna med en ring som
  pulserar. Varelserna står som sitt tecken. Fällornas märken, kritpilarna,
  föremålen och de döda (x) står i krita på golvet, och draget är en pil
  vid dörren som blinkar. Ekot färgar rummets golv.
- Kompassen nere till vänster och pilarna vid kanten kom senare (se
  Efter stegen nedan).

## Så är steg 3 byggt

- Modellerna står i `klient/src/modeller.ts`. Tillståndet har klassens id
  (`klassId`) för varje vandrare och varelsens id (`typ`) för varje
  varelse, så att modellen hittas oavsett språk.
- Vandrarna har kropp, huvud, armar och ben, i sin färg. Klasserna bär sitt:
  Legoknekten hjälm, sköld och svärd, Tjuven mörk huva, säck och dolk,
  Örtkvinnan stav med en blomma och korg, Kartritaren en kartrulle,
  Ormtjusaren en orm kring halsen och en flöjt, Spårfinnaren hatt och
  käpp, Lyktbäraren en lysande lykta och Pilgrimen bredbrättad hatt och
  stav med kalebass. Utan klasser bär alla en huva.
- Varelserna: Blekkryparen med sex ben, Ormynglet rest ur dammet,
  Fjällvakten i fjällrustning med sköld och spjut, Gallgasten som en
  genomskinlig grön skepnad i trasor och Benväktaren av ben och revben. En
  ny varelse i `monster.json` utan modell visas med sitt tecken.
- Figurerna står kvar mellan rundorna och går till sin nya plats. Till ett
  rum bredvid går de genom dörren; längre bort, som ur källarvalven, står
  de bara där. De vänder sig åt det håll de går, och benen och armarna
  rör sig. Den som dör lägger sig grå på golvet, med ett x i krita.
- Den egna ringen följer figuren.
- Modellerna står i koden, inte i `data/`. Det går att flytta dem dit
  senare, om de ska kunna ändras utan att bygga om.

## Så är steg 4 byggt

- Utseendet följer avslöjandets nivå med samma färger och takt som den
  platta kartan (`LOOK` i `klient/src/karta.ts` och `web/map.js`).
- Väggarnas färg räknas ut i en shader. På nivå 2 andas kritan; på nivå 3
  och 4 går hjärtslaget som en våg från tornets mitt, två slag tätt efter
  varandra och sedan vila.
- Från nivå 3 mörknar golvet mot kött i takt med slagen, och revben syns
  genom golvet. På nivå 4 är väggarna kött och de öppna portarna munnar med
  tänder som tuggar. När vandringen är slut har munnarna slutit sig.
- När ormen vaknar lyser två ögon över guldet och guldet glöder, medan
  kritringen blinkar som förut.

## Efter stegen

- Kompassen: ett N med en pil nere till vänster på kartan, åt det håll
  norr ligger på skärmen. Riktningen räknas ut från kamerans vinkel.
- Pilarna vid kanten: portarna och kammaren man känner till, och vandrarna
  man ser, får en pil vid kanten med sitt tecken när de ligger utanför
  bilden, bara när kartan följer någon.
- Tangenterna vid dörrarna: i rummet vandraren står i står en tangent med
  pilen vid varje öppen dörr och port, så att man ser vart piltangenterna
  går. Det valda draget lyser på sin tangent i vandrarens färg, i stället
  för den blinkande pilen på golvet. På pekskärmen visar knapparna sneda
  pilar (↗ ↘ ↙ ↖) när kartan är snett ovanifrån, och tangenterna vid
  dörrarna visar samma pilar och går att trycka på för att gå dit.
- Föremålen på figuren: det vandraren bär syns på den egna figuren,
  vid bältet, på bröstet, på ryggen eller på huvudet (`fastPa` i
  `klient/src/modeller.ts`). Andras väskor är hemliga och syns inte, men
  spelledarbordet visade föremålen på alla, tills det togs bort
  2026-11-22. Förbannade föremål ser ut som andra.
- Teckenförklaringen med bilder: med kartan snett ovanifrån står en liten
  bild av modellen i stället för tecknet för vandrarna, varelserna,
  kammaren och portarna (`bild` i `klient/src/karta.ts`, `legendImage` i
  `web/map.js`). Sedan 2026-11-22 har fällorna också en bild. Det som står
  i krita på golvet, som pilarna, har kvar sitt tecken. Bilderna ritas en
  gång och sparas.
- Telefoner och svaga enheter: på en pekskärm ritas högst 30 bilder i
  sekunden. Kartan mäter hur lång tid bildrutorna tar, och går det trögt
  (median över 45 ms på en pekskärm, 30 ms annars) slår den av skuggorna
  och ritar med en pixel per punkt. Enheter med högst två kärnor eller 2 GB
  minne börjar så. `?kvalitet=lag` i adressen tvingar fram det. Zoomen Nära
  går lite närmare än förut.
- Modellerna i `data/`: varje klass i `klasser.json` och varelse i
  `monster.json` har `modell` med `farg` (#rrggbb) och `storlek` (0,5–2).
  För en varelse är färgen kroppens, för en klass färgen på det klassen
  bär. Formerna står i koden. Ändringar gäller efter Läs om innehållet,
  eftersom färgen och storleken följer med tillståndet. Kontrollen säger
  till om en färg eller storlek är fel.

## Spelbrädet (2026-11-22)

Kartan fick en mer fantasyaktig stil, som ett spelbräde i fackelsken, i
fem steg. Allt ritas i koden, utan bilder.

1. **Sten och trä.** Golven är huggna stenplattor (`stengolv` i
   `klient/src/krita.ts`), murarna murverk med krön (`murverk`), tjockare
   och lite högre än förut. Murverket läses efter platsen i världen
   (`varldsUv`), så blocken är lika stora på alla murar. Dörrarna har
   dörrposter och en uppslagen trädörr (`tra`). Portarna är portvalv med
   torn, spetsiga tak och trädörrar som står utåt, eller är stängda med
   röda järnband (`portGeometri`, som teckenförklaringen också använder).
2. **Fackelsken.** Varje rum har två facklor på murarna mot kameran, med
   fladdrande lågor och ett sken på golvet. Rummet kartan följer lyses av
   ett varmt punktljus, och resten av tornet blir mörkare med avståndet
   (`dunklaShader`). Utan någon att följa får tornet ett jämnare ljus.
   Ramen runt kartan (`.isoram` i `web/style.css`) är en tunn
   pergamentkant med mörker mot kanterna.
3. **Inredning.** `klient/src/inredning.ts` bygger pelare, en brunn,
   bänkar, ett spjut, ett altare, en spiraltrappa, hyllor, tunnor, säckar,
   lyktor, stenblock, en ässja, ett städ, bokhyllor, ett bord och en hög
   med stenar och ben. Vilken rumstyp som har vad står under
   `modell.inredning` i `data/texter.json`; tillståndet skickar det som
   `in` för varje rum, och kontrollen säger till om ett namn inte finns.
   Det höga står mot norra och västra muren.
4. **Miniatyrer och fällor.** Vandrarna och varelserna står på socklar
   (`miniatyr` i `klient/src/modeller.ts`) med en kant i vandrarens färg
   eller mörkröd för varelserna. Den döda läggs ner på sockeln. En fälla
   som har slagit till är en modell (`falla`), och tillståndet säger
   vilken fälla märket hör till (`marki`).
5. **Brickorna.** Ett nytt rum faller ner som en bricka, vänds i luften och
   landar, och byggs först då (`laggBricka` och `landa` i `karta.ts`). Ett
   nytt torn eller fler än tre rum på en gång läggs utan brickor, liksom
   för den som har bett om mindre rörelse.

## Inte byggt

- Formerna i `data/`: modellernas delar står i `klient/src/modeller.ts`.
  En ny klass eller varelse får en vanlig figur eller sitt tecken tills
  någon bygger en form åt den där.
