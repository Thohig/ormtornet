# Att göra och idéer

Status: samlat 2026-11-18, kompletterat efter provspelet 2026-11-21 (se
[provspel-2026-11-21.md](provspel-2026-11-21.md)). Inget här är bestämt. Användaren väljer vad som
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
  nere till höger? Hamnar pilarna vid kanten under rummets namn
  eller tecknen i kartans övre hörn?
- Klassvalet: räcker 30 sekunder för att läsa korten? På en telefon står
  korten sedan 2026-11-21 två och två utan den kursiva beskrivningen, med
  nedräkningen fast överst. Går de att läsa på en riktig telefon?
- Kartritaren som ser in i grannrummen: för stark nu när alla andra bara
  börjar med sitt eget rum?
- Stormen i introt: hur musiken och åskan låter. Den är inte avlyssnad.
- Introts första 5 sekunder, när Enter inte gör något: fastnar någon
  där och tror att sidan har hängt sig?
- Kartan som ett spelbräde (2026-11-22): är rummen längre bort för mörka?
  Blir murarna med krön för plottriga? Skymmer den egna figuren det som
  står i rummets nordvästra hörn, som trappan och ässjan? Går kartan lika
  lätt på en telefon med facklorna, inredningen och brickorna som faller?
  Nivå 3 och 4 är bara sedda genom att tvinga fram dem.
- Ormkorten (2026-11-22): märks det att kortet vänds upp när man tar en
  näve, och läser man det innan det glider undan? Förstår man
  sömnmätaren och högarna bredvid guldet? Uppvaknandets stora kort är
  bara sett med påhittade meddelanden, och krönikans kort i slutet är inte
  sedda alls.
- Färre inställningar (2026-11-22): strid är alltid på från första rundan
  och klasserna alltid med. Blir striden för tidig, och saknar någon att
  kunna spela utan klasser?
- Striderna i krönikan (2026-11-22): blir krönikan för lång när många
  slåss, och är raden med siffror för lång på en telefon?
- Rumsblocket i loggen (2026-11-22): förstår spelarna vad figurerna på
  kartan är utan teckenförklaringen, när bara texten säger vad som finns
  i rummet? Blir rumsbeskrivningarna för långa på en telefon, och hittar
  man "mer ↓"?
- Det nya gränssnittet (2026-11-22, [nytt-granssnitt.md](nytt-granssnitt.md)):
  läser man raden Nu under kartan mitt i en strid? Är loggen i hela
  högerkolumnen bra på dator? På telefon: skymmer den fasta översta raden
  och knapparna längst ner något, och fungerar det att loggen rullar med
  sidan? Bara provat i en emulerad telefon.
- Föremålen som verkar av sig själva (2026-11-22,
  [passiva-foremal.md](passiva-foremal.md)): märker man när rökpulvret,
  spegelskärvan, kastkniven, kritan och timglaset verkar? Blir kritans
  fem pilar för många, och liknar spegelskärvan Kartritaren för mycket?
  Bara provat med tester.
- Stycken och färre kommandon (2026-11-22,
  [samlat-forslag.md](samlat-forslag.md)): går det att förstå rummen när
  de bara beskrivs, och känner man igen spåren efter fällorna? Hörs
  ljuden, eller försvinner de sist i stycket? Saknar någon att välja
  motståndare själv, eller att se vart man gick på kortet? Hittar man Z
  för kartan? Bara provat med tester och i webbläsaren.
- Namnen som tornet ger (2026-11-22, [slumpade-namn.md](slumpade-namn.md)):
  tycker spelarna om namnen, eller saknar de att skriva sitt eget? Får
  den som byter telefon eller webbläsare mitt i kvällen ett nytt namn, och
  gör det något?
- Bara den egna porten: hittar spelarna ut, eller går de vilse när de
  inte ser de andra portarna?
- Tecknen på kartan (2026-11-19): går de att förstå, och skymmer de
  kartan när den har vuxit?
- Loggen (2026-11-19, ändrad 2026-11-22): är takten lagom, ungefär tre
  rader i sekunden? Hittar man knappen "↑ nya rader"? Sedan 2026-11-22
  står bara den senaste rundan i bild och de äldre under strecket
  "tidigare"; saknar någon att se förra rundan utan att rulla?
- Kortet efter rundan (2026-11-19): säger det det viktiga, är sju sekunder
  lagom, och skymmer det något på kartan?
- Prickarna och namnen i färg (2026-11-19): syns vem som väljer? Den vita
  färgen (plats 8) syns knappt bland den vanliga texten i loggen. Hjälper tecknen i loggens marginal, eller
  blir de brus?

- Från provspelet 2026-11-21, bara provat i en emulerad telefon och med
  robotar:
  - Pilarna vid dörrarna går att trycka på på en pekskärm. Hittar man det,
    och träffar man dem med fingret?
  - Ett steg mot en vägg avvisas direkt. Märks det, eller tror man att
    tryckningen inte gick fram?
  - Viska har en fast, osynlig plats bland knapparna tills den behövs.
    Undrar någon över luckan?
  - Tjuven tog 173 guld på fem nävar i en vandring. Är hen för stark i
    kammaren?

- Spela från en länk (2026-11-21, [github-pages.md](github-pages.md),
  steg 5): kommer alla in från mobilnät och andra nät, eller behövs TURN?
  Håller tornets flik när den har legat i bakgrunden mer än fem minuter?
  Kommer den som tappar anslutningen tillbaka? Hur lång tid tar det att
  ladda tornet (9 MB)? Saknar någon QR-koden, som togs bort 2026-11-22?

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

### Småsaker från provspelet

Sett 2026-11-21 men inte ändrat:

- "Solen har vandrat halva vägen" kan sägas vid 19 handsbredder kvar när
  klockan har skyndat på solen.
- I mörker visar kartan snett ovanifrån fortfarande pilarna vid dörrarna.
- De vaga `_dolt`-texterna och originalen finns båda kvar i `data/`. Sedan
  egen karta är det enda läget kan de slås ihop, så att varje text finns
  en gång.

### Städning

- Grenarna `claude/awesome-babbage-ekl3oh`, `claude/ecstatic-dijkstra-w5bt1o`,
  `claude/exciting-mayer-jcrxup`, `claude/relaxed-ritchie-ow2ciq` och
  `claude/vibrant-fermat-berckf` på GitHub har inga committar som inte
  redan finns i `main` (2026-11-21), och kan tas bort.
- Repot är publikt sedan 2026-11-21 men har ingen `LICENSE`. Bestäm om
  spelet ska ha en licens.

## Kvar från förslagen

Från [spelledaren.md](spelledaren.md):

- Fler enskilda inställningar: chans för fällor, varelser och dörrar, liv,
  ormkort, nivågränser och personlighet för varje AI-spelare.
- Namngivna förval och paus mellan rundor.
- Frö, för att kunna spela samma torn igen.
- Statistik över flera kvällar och redigering av datafilerna i webbläsaren.
- En färdigbyggd app som går att dela utan Go.

Från [foremal-och-hemligheten.md](foremal-och-hemligheten.md):

- Släppa ett föremål med flit, eller byta föremål med en annan vandrare.

Från [slumpade-namn.md](slumpade-namn.md):

- Knappen Nytt namn i lobbyn, som lottar ett nytt namn (steg 3). Valdes
  bort 2026-11-22; i stället glöms namnet efter 12 timmar. Kan tas upp
  igen efter ett provspel.

Från [isometrisk.md](isometrisk.md):

- `bild` i `klient/src/karta.ts`, som ritade bilderna i
  teckenförklaringen, används inte längre och kan tas bort nästa gång
  kartan byggs om.

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
  bara färger. Socklarna kan få en initial eller en form.
- **Loggens tecken efter innehåll:** i dag följer tecknet radens klass
  från servern. Med en egen klass för strid, ljud och fynd kunde loggen
  visa till exempel ⚔, ♪ och ✦, som i förslaget till panelerna.
- **Andra kartstilar:** av stilförslagen 2026-11-18 valdes krita på
  skiffer. Kvar som idéer: skarp terminal med hela väggar, ritade rum med
  runda märken, kartritarens karta på sotat pergament eller guld på
  läder.
- **Tester för klienten:** i dag testas bara servern, så kartan och
  panelen provas för hand.
