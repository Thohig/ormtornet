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
- Klasserna: är någon för stark eller för svag? Tjuven var svagast med
  AI-spelare.
- Viskningen: gissar någon rätt, och känns fördelen lagom?
- Åtta spelare med telefoner och datorer: zoomen på kartan, utsnittet och
  knappen för allt man känner till.
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
- Bara den egna porten: hittar spelarna ut, eller går de vilse när de
  inte ser de andra portarna?
- Tecknen på kartan (2026-11-19): går de att förstå, och skymmer de
  kartan när den har vuxit? Hjälper tecknen i loggens marginal, eller
  blir de brus?

### Justera efter provet

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
