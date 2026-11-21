# Att göra och idéer

Status: samlat 2026-11-18. Inget här är bestämt. Användaren väljer vad som
byggs, och förslag tas fram innan något byggs.

## Att göra

### Prova med människor

Det mesta som byggdes 2026-11-18 är bara provat med AI-spelare och tester.
En riktig kväll behövs för att se:

- Balansen på Normal: når någon kammaren, och vad tar livet?
- Egen karta och lottade portar: pratar spelarna med varandra, eller blir
  det bara rörigt?
- Klasserna: är någon för stark eller för svag? Med AI-spelare
  (2026-11-19, 10 000 vandringar, klasserna rättvist fördelade) kommer
  alla ut med skatt i 50–61 % av vandringarna på Normal. Kartritaren och
  örtkvinnan justerades då. Känns det rätt med människor?
- Örtkvinnan (2026-11-19): hon tappar nu all sin skatt om ormen vaknar
  medan hon är i kammaren (85 % på Lätt). Är det för hårt? Sänk i så fall
  `extra_forlust` i `data/klasser.json` till cirka 0,35.
- Viskningen: gissar någon rätt, och känns fördelen lagom?
- Kartan (2026-11-19): är Mellan på datorn och Nära på telefonen bra
  förval? Går det att zooma med två fingrar på en riktig telefon (bara
  provat med hjul, tangenter och knappar)? Hittar man M och knapparna
  nere till höger? Hamnar pilarna vid kanten under teckenförklaringen
  eller tecknen i kartans övre hörn?
- Klassvalet: räcker 30 sekunder för att läsa korten, och går de att läsa
  på en telefon?
- Kartritaren som ser in i grannrummen: för stark nu när alla andra bara
  börjar med sitt eget rum?
- Stormen i introt: hur musiken och åskan låter. Den är inte avlyssnad.
- Introts första 5 sekunder, när Enter inte gör något: fastnar någon
  där och tror att sidan har hängt sig?
- Kritkartan (2026-11-18): syns vandrarnas prickar och varelserna bra på
  en telefon, och märks avslöjandet lika tydligt i kritan som i braillen?
  Nivå 3 och 4 är bara sedda genom att tvinga fram dem.
- Teckenförklaringen på kartan (2026-11-19): nu bara tecknen, med
  betydelsen vid tryck. Hittar spelarna att man kan trycka på dem, och
  skymmer den fortfarande något på en telefon? I utsnittet tar den med
  tecken från kända rum som ligger utanför bild.
- Bara den egna porten: hittar spelarna ut, eller går de vilse när de
  inte ser de andra portarna?
- Tecknen på kartan (2026-11-19): går de att förstå, och skymmer de
  kartan när den har vuxit?
- Loggen (2026-11-19): är takten lagom, ungefär tre rader i sekunden?
  Hittar man knappen "↑ nya rader", och är det bra att den senaste
  rundan står överst och äldre rundor fälls ihop?
- Kortet efter rundan (2026-11-19): säger det det viktiga, är sju sekunder
  lagom, och skymmer det något på kartan?
- Prickarna och namnen i färg (2026-11-19): syns vem som väljer? Den vita
  färgen (plats 8) syns knappt bland den vanliga texten i loggen. Hjälper tecknen i loggens marginal, eller
  blir de brus?

- Engelskan (2026-11-19): är tonen rätt, särskilt i introt, ledtrådarna
  och avslöjandet? Är det bättre med "the Silent Chapel" än "The Silent
  Chapel" mitt i meningar? Får alla knappar och etiketter plats på en
  telefon? Engelskan är skriven av AI och inte läst av någon som spelar.

### Justera efter provet

- Legoknekten på Mardröm: med AI-spelare kommer hen ut med skatt i 43 %
  av vandringarna, mot 21–32 % för de andra klasserna. Extralivet och
  skölden väger tungt när alla har få liv.
- Svårigheterna med AI-spelare (2026-11-19, 10 000 vandringar, 4
  AI-spelare): ut med skatt 71, 54, 45 och 28 % på Lätt, Normal, Svår och
  Mardröm, döda 10, 24, 26 och 50 %. Fler varelser på Svår
  (`varelsechans` 0,32) gav ingen mätbar skillnad; det som avgör är
  antalet liv.
- Rummen bredvid porterna utan fällor (2026-11-21, 4000 vandringar per
  svårighet, 4 AI-spelare): ut med skatt 74, 59, 50 och 36 % på Lätt,
  Normal, Svår och Mardröm (förut 73, 56, 47 och 33 %), döda 9, 20, 22 och
  44 % (förut 9, 21, 23 och 44 %). Fällor som slog till per vandrare 2,1
  i stället för 2,5. Känns början lagom med människor?

- Göra Normal mildare om den är för hård, i `data/svarighet.json` (färre
  fällor eller ett liv till). Uppdatera då tabellen i README.
- AI-spelarna når sällan kammaren (13–23 % i simuleringarna) och söker
  nästan aldrig, så vandringar med bara AI-spelare blir tama. De kan bli
  modigare, söka oftare och följa vägar som andra redan har gått.

### Städning

- Grenen `claude/vibrant-fermat-berckf` på GitHub har inga committar som
  inte redan finns i `main`, och kan tas bort.

## Kvar från förslagen

Från [spelledaren.md](spelledaren.md):

- Fler enskilda inställningar: chans för fällor, varelser och dörrar, liv,
  ormkort, nivågränser och personlighet för varje AI-spelare.
- Namngivna förval och paus mellan rundor.
- Frö, för att kunna spela samma torn igen.
- Ingrepp: flytta en vandrare, slå klockan, låta ormen röra sig, spela upp
  ett ljud, lägga ett föremål i ett rum.
- Statistik över flera kvällar och redigering av datafilerna i webbläsaren.
- En färdigbyggd app som går att dela utan Go.

Från [foremal-och-hemligheten.md](foremal-och-hemligheten.md):

- Släppa ett föremål med flit, eller byta föremål med en annan vandrare.

## Nya idéer

### Spelet

- **Lag:** två eller tre lag som delar skatten. Med egen karta får
  pratandet ett syfte.
- **Uppdrag:** varje vandrare får ett hemligt mål i lobbyn, till exempel
  "kom ut med kronan" eller "se till att ormen vaknar två gånger", som ger
  extra guld.
- **Tornets humör:** tornet reagerar på vandrarna. Mycket strid gör
  väggarna "irriterade", med fler fällor och varelser som rör sig.
- **Kvällens krönika:** en sammanställning efter flera vandringar samma
  kväll: vem som vann flest gånger, flest invigda och flest döda.

### Tekniskt

- **Ljud:** ett reglage för volymen, och en knapp som stänger av musiken
  men behåller ljuden, på spelarsidan.
- **Färgblindhet:** ett läge där vandrarna skiljs åt med symboler, inte
  bara färger. Med kritkartan kan prickarna få en initial eller en form.
- **Loggens tecken efter innehåll:** i dag följer tecknet radens klass
  från servern. Med en egen klass för strid, ljud och fynd kunde loggen
  visa till exempel ⚔, ♪ och ✦, som i förslaget till panelerna.
- **Andra kartstilar:** av stilförslagen 2026-11-18 valdes krita på
  skiffer. Kvar som idéer: skarp terminal med hela väggar, ritade rum med
  runda märken, kartritarens karta på sotat pergament eller guld på
  läder.
- **Tester för klienten:** i dag testas bara servern, så kartan och
  panelen provas för hand.
