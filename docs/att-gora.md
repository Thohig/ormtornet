# Att göra och idéer

Status: samlat 2026-11-18, kompletterat efter provspelen 2026-11-21 och
2026-11-22 (se [provspel-2026-11-21.md](provspel-2026-11-21.md) och
[provspel-2026-11-22.md](provspel-2026-11-22.md)). Inget här är bestämt. Användaren väljer vad som
byggs, och förslag tas fram innan något byggs.

## Att göra

### Prova med människor

Den nya striden (2026-11-24, [bara-anfall.md](bara-anfall.md)) är bara
provad med AI-spelare, som alltid anfaller:

- Lär sig spelarna varelsernas lynne, och vågar de söka eller dricka när
  en varelse står i rummet? Håller de garden?
- Säger panelen Nu ("Du träffar lätt · den träffar ibland") för mycket,
  så att det inte finns något kvar att lära sig?
- Blir striden tråkig med ett enda val, särskilt mellan vandrare, där man
  förut läste den andra?
- Svärdet med rött fäste: märker man förbannelsen, och är Kraft +2 värt
  ett liv per runda man inte anfaller?
- Guldet per vandrare blev lägre än före bytet (36–47 mot 39–52 på
  Normal med AI-spelare). Behöver varelserna ge mer guld?

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
- Ormens namn som en tävling (2026-11-23,
  [namnet-som-tavling.md](namnet-som-tavling.md)): hittar vandrarna
  ristningen i tid (`ristning_chans`), får de som undviker strid
  varelsens del, förstår de att huvudet är först och svansen sist, och
  blir det för dyrt att viska fel nu när ormen blir oroligare? Hinner
  någon bli invigd före nivå 4?
- Inga namn (2026-11-22): går det att prata om rum, fällor och varelser
  utan namn, eller hittar spelarna på egna? Läser de korta beskrivningarna
  (`kort` i `data/`) bra i striden, som "Det bleka mångbenta träffar
  dig"? Alla är skrivna utan att provas.
- AI:n tar över den som står still (2026-11-22): är två rundor lagom, och
  märker spelaren att hen fick tillbaka vandraren med sitt drag?
- Spelledaren i samma flik (2026-11-22): hittar spelledaren tillbaka till
  värdsidan och spelet? Provat lokalt och i skalet, men inte på en riktig
  kväll.
- Kartan (2026-11-19): är Mellan på datorn och Nära på telefonen bra
  förval? Går det att zooma med två fingrar på en riktig telefon (bara
  provat med hjul, tangenter och knappar)? Hittar man M och knapparna
  nere till höger? Hamnar pilarna vid kanten under rutan Mörker eller
  Källarvalven, eller tecknen i kartans övre hörn? Webbläsarprovet
  (2026-11-23) såg på telefonen ibland en tangent vid dörren under
  tecknen (`#status`), solmätaren eller märket.
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
- Spelkorten (2026-11-22, [spelkort.md](spelkort.md)): blir det för
  många kort, och hinner man läsa dem på åtta sekunder? Står kortet i
  vägen för kartan på en telefon? Saknar någon siffrorna i loggen, nu när
  de står på kortet? Förbannelsens andra vändning, porten som stängs och
  döden i Din vandring är bara sedda i tester och med påhittade
  meddelanden; solen, varelserna, namnets delar och vandringens sista
  kort är sedda med AI-spelare.
- Klasskorten (2026-11-22): går det att läsa pergamentet på en telefon,
  och tar vändningen i klassvalet för lång tid av de 30 sekunderna? Ett
  kort som någon annan tar är bara sett med ett påhittat läge.
- Vandringens största stunder (2026-11-22): väljs rätt stunder, och är
  fem lagom? Bara sedda med påhittade stunder; en kort vandring med AI
  gav inga.
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
  läser man raden Nu under kartan mitt i en strid? På telefon: skymmer den
  fasta översta raden och knapparna längst ner något, och fungerar det att
  loggen rullar med sidan? Bara provat i en emulerad telefon.
- Loggen på kartan (2026-11-22, [frot-och-bredden.md](frot-och-bredden.md)):
  går texten att läsa över ljusa golv och facklor, räcker den förra
  rundan svagt ovanför, och hittar man Loggen och L? På en riktig
  surfplatta och en telefon på bredden: får knapparna plats under loggen,
  och är kartan bredvid stor nog? Bara provat i en emulerad skärm, och
  inte på Pages.
- Zoomen som glider (2026-11-22): är ett hack på hjulet lagom stort
  (talet 300 i `web/app.js`), och känns nypningen på en riktig pekskärm
  naturlig?
- Märket under vandringen (2026-11-22): öppnar spelledare värdsidan i en
  ny flik, och tar Till spelet där tillbaka till spelets flik i Safari och
  Chrome? Inte provat med ett riktigt klick.
- Föremålen som verkar av sig själva (2026-11-22,
  [passiva-foremal.md](passiva-foremal.md)): märker man när rökpulvret,
  spegelskärvan, kastkniven och timglaset verkar? Liknar spegelskärvan
  Kartritaren för mycket? Bara provat med tester.
- Benluren (2026-11-22), som ersatte kritan: hör man varelserna i
  grannrummen i rumsblocket, och är ljuden två rum bort till hjälp eller
  bara brus? Kartritaren har en dryck med sig i stället för kritans extra
  streck; räcker det?
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
- Den dödes guld (2026-11-22): går någon dit en annan har dött för att
  hämta guldet, och blir det för lätt att plocka upp en hel kammares
  skatt? Hör man var skriket kom ifrån, och går någon dit? AI-spelarna går
  dit efter `guldlust` i `data/ai.json`; är det lagom?
- Bara den egna porten: hittar spelarna ut, eller går de vilse när de
  inte ser de andra portarna?
- Tecknen på kartan (2026-11-19): går de att förstå, och skymmer de
  kartan när den har vuxit?
- Loggen och krönikan efter genomläsningen (2026-11-22): märks det att
  en pågående strid bara hörs första rundan, eller tror man att den har
  slutat? Läser man "det mångbenta" som samma varelse? Blir Din vandring
  och krönikan för mycket att läsa, nu när orsakerna och vändpunkterna
  står med? Att visa förra rundan nedtonad när den senaste är kort valdes
  bort, för att inte visa för mycket på en gång.
- Loggen (2026-11-19, ändrad 2026-11-22): är takten lagom, ungefär tre
  rader i sekunden? Hittar man knappen "↑ nya rader"? Sedan 2026-11-22
  står bara den senaste rundan i bild och de äldre under strecket
  "tidigare"; saknar någon att se förra rundan utan att rulla?
- Kortet efter rundan (2026-11-19, sedan 2026-11-22 spelkortet, som står
  åtta sekunder): säger det det viktiga, är sju sekunder
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

- Svårigheterna efter justeringen 2026-11-23 (3000 vandringar, 4
  AI-spelare): ut med skatt 82, 69, 54 och 40 % på Lätt, Normal, Svår och
  Mardröm, döda 9, 14, 27 och 43 %. Normal har 7 liv, Mardröm 6 liv men
  +1 i skada, Lätt en fjärdedel längre sol, och kammaren ger mer guld på
  Svår och Mardröm. Med människor: är Normal lagom nu, känns Mardröms
  hårdare träffar rätt, och går någon till kammaren för guldets skull?
  Klassernas extraliv väger fortfarande tungt på Mardröm (Legoknekten
  49 % ut med skatt, Kartritaren 32 %).
- Kartritaren är fortfarande svagast i simuleringen 2026-11-22 (300
  vandringar, 3 AI-spelare): 49 % överlever och 22 % vinner, mot
  Legoknektens 79 % och 41 %.
- AI-spelarna når sällan kammaren (13–23 % i simuleringarna) och söker
  nästan aldrig, så vandringar med bara AI-spelare blir tama. De kan bli
  modigare, söka oftare och följa vägar som andra redan har gått. I
  provspelet 2026-11-22 kom ingen av två AI-spelare till kammaren på 40
  rundor i en vandring.

### Småsaker från provspelet

Sett 2026-11-21 men inte ändrat:

- "Solen har vandrat halva vägen" kan sägas vid 19 handsbredder kvar när
  klockan har skyndat på solen.
Sett 2026-11-22 men inte ändrat:

- Har förra rundan ingenting i loggen (man stod kvar i ett tomt rum)
  börjar en stöt i muren ändå en ny runda, och står där före rumsblocket
  när rundan sedan händer. Sedan 2026-11-24 ger ett känt, tomt rum som man
  går in i raden "Du känner igen rummet", så det händer bara när man stod
  kvar.
- Lådan med loggen börjar med den senaste rundan över hela höjden, så den
  ser tom ut under den första rundan.
- Spelkortet täcker rummet på telefonen en stund efter varje runda.
  Webbläsarprovet (2026-11-23) bekräftar det: på telefonen låg spelkortet
  (`#spelkort`) eller det stora kortet (`#storkort`) ibland över
  tangenterna vid dörrarna, så att de inte gick att trycka på.
- Värdsidans avsnitt Inför nästa vandring är låst under vandringen, utom
  Fröet.
- Pages lokalt: spelledarens Gå in själv fastnade en gång på introt när
  den lokala filservern tappade filer. Prova på riktiga Pages.

Sett efter provspelet 2026-11-22, med förslag men inte ändrat:

- **Mörkret ser ut som ett rum utan dörrar.** När lyktan är släckt ritar
  kartan hela murar runt rummet, och man går rakt igenom dem. Förslag:
  rita inga murar alls runt rummet i mörkret, så att kanterna ser okända
  ut i stället för stängda. Ordet Mörker och ◐ står kvar.
- **Timglaset märks knappt.** Solen står still en runda, men det är bara
  en rad i loggen. Förslag: ett stort kort för alla inne, som solkortet,
  till exempel "Solen står still".
- **Kompassen kan peka mot en stängd port.** Den räknar inte med att
  portarna sluter sig på nivå 4 (`compass` i `items.go`).

Sett i provspelet 2026-11-24 och prövat:

- Fjällrustningen tog 4 liv på två rundor från Kartritaren på Normal, när
  hen anföll den. Sedan 2026-11-24 finns ingen parering; gestalten i
  fjällrustning har Försvar 2 och går lättare att träffa
  ([bara-anfall.md](bara-anfall.md)).

Från förut:

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

Spelmekaniker (förslag 2026-11-22):

- **Märken på väggen:** var och en kan rita ett av tre märken i rummet:
  fara, skatt eller vägen ut. Den som kommer in efteråt ser det i
  rumsblocket och på kartan. Ett sätt att hjälpa, eller lura, varandra med
  egen karta.
- **Spakar och dörrar:** en spak i ett rum öppnar en stängd dörr i ett
  annat, bara medan någon står vid spaken. Samarbete på avstånd; passar
  med portarna som sluter sig på nivå 4.
- **Smyga eller springa:** smyga ett rum utan att höras eller väcka något,
  gå som i dag, eller springa två rum i ett drag som hörs i hela tornet,
  väcker varelser och gör ormens sömn lättare.
- **Altaret:** offra guld för ett liv, ett varsel eller en handsbredd sol.
  Ju djupare in, desto bättre byte.
- **Gå tillsammans:** vandrare i samma rum kan hålla ihop och flytta som
  en grupp; fällans skada delas, de slåss sida vid sida, men gruppen hörs
  mer.
- **Bonus vid porten:** kommer alla levande ut får var och en en
  guldbonus, men bara om ingen dör. Ett socialt dilemma utan strid.
- (Den dödes väska byggdes 2026-11-22, se
  [den-dodes-guld.md](den-dodes-guld.md).)

Nya fällor (förslag 2026-11-22). Fällornas verkan styrs av id:t i
`events.go`, så varje ny fälla behöver en rad kod, en kortbild och texter
på båda språken:

- **Den sura pölen** ("en pöl som fräser"): fräter bort ungefär en
  femtedel av skatten men inget liv.
- **Svalget** ("ett golv som lutar inåt"): du glider ett rum närmare
  skattkammaren.
- **Hjärtslaget** ("ett dovt dunkande"): ormen sover lättare, sömnmätaren
  sjunker ett steg, och alla hör det.
- **Tänderna** ("två rader spetsiga stenar"): −1 liv och fast 1 runda, en
  mildare ormgrop för de första nivåerna.
- **Den täta dimman** ("en vit dimma som suddar ut"): kartan glömmer de
  två–tre senaste rummen.
- **Rösterna** ("viskningar i väggarna"): ett falskt ljud från ett
  grannrum, kanske bara från nivå 3.

Nya varelser (förslag 2026-11-22). 1–4 går att lägga till med bara data,
plus en form i `klient/src/modeller.ts` och en kortbild:

1. **Fladdrarna** ("något litet och fjällvingat"): liv 1, skada 1, guld
   5, vandrar ofta (0,6), Kraft 0, Försvar 2, lynne 20.
2. **Fjällhunden** ("något lågt som springer på fyra ben"): liv 3, skada
   1, guld 15, vandrar ofta (0,5), Kraft 2, Försvar 0, lynne 90.
3. **Tungan** ("en kluven tunga ur väggen"): liv 3, skada 2, guld 20,
   vandrar aldrig.
4. **Ruvaren** ("något stort som ligger på ägg"): liv 4, skada 1, guld
   50, vandrar aldrig, Kraft 1, Försvar 3, lynne 10.
5. **Tjuvormen** ("en smal, svart orm med guld i munnen"): stjäl 10–20
   guld per träff i stället för liv och flyr; dödar man den får man
   guldet tillbaka. Kräver en ny regel i `monster.go`.

Chanskort vid sökning (förslag 2026-11-22):

- I stället för att slumpa utfallet efter vikter drar den som söker ett
  kort ur en lek som hela tornet delar, som ormkorten, och kortet visas
  som ett spelkort. Ett draget kort är borta, så den som söker först tar
  det bästa. Mörkare kort blandas in på nivå 2 och 3. Leken ligger i
  `data/` med vikt, nivå, verkan och texter.
- Dagens utfall blir kort: skatt, dryck, hemlig gång, fälla och inget
  (med fler texter).
- Föremålen blir egna kort, lika många som deras vikt i `foremal.json`.
  Varje föremål finns då bara så många gånger per torn; i dag kan samma
  föremål hittas hur många gånger som helst.
- Nya kort som hjälper: kartbiten (visar två–tre rum runt dig), gluggen
  (solen backar en handsbredd, sällsynt), det lösa fjället (ormen sover
  djupare), spåren (åt vilket håll närmaste vandrare gick) och varslet.
  (Ristningen finns sedan 2026-11-23 som en av namnets tre källor, se
  [namnet-som-tavling.md](namnet-som-tavling.md).)
- Nya kort som stjälper: den lösa stenen (ett vakna-kort i leken), boet
  (en varelse kryper fram), falluckan (ner i källarvalven) och det kalla
  guldet (40 guld, men ormen blir oroligare och alla hör det).
- Kort som påverkar andra: röret i väggen (du hör ett rum långt bort, och
  de där hör dig) och gömstället (lämna guld som bara du kan hämta; dör du
  hittar nästa som söker det).
- Klasserna behåller sina fördelar: tjuven drar om ett inget-kort,
  Kartritaren får den hemliga gången om den finns kvar.
- Frågor: gemensam lek eller en per vandrare? Ska man se hur många kort
  som är kvar? Ska ristningen och gluggen vara med, när de rör
  invigningen och solen? Är röret och gömstället värda nya regler?

Andra föremål som föreslogs när kritan ersattes av benluren
(2026-11-22):

- **Slagrutan:** säger när du kommer in i ett osökt rum om det lönar sig
  att söka (skatt eller föremål). Åtta rum. Går bäst ihop med
  chanskorten, där utfallet bestäms i förväg.
- **Ormtanden:** dina träffar mot varelser tar två liv, tre gånger.

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
- **Tester för klienten:** webbläsarprovet (`verktyg/webbprov.mjs`,
  2026-11-23) spelar en kort vandring på dator och telefon och letar efter
  fel på sidan. Kvar som idéer:
  - **Webbläsarprovet i CI:** det kräver Chromium vid varje push och tar
    ett par minuter.
  - **Spara en vandring och spela upp den igen:** fröet och dragen avgör
    allt, så värdsidan kunde spara dem i en fil som `-spela-om fil`
    spelar upp. Då går ett fel från en kväll med människor att återskapa
    exakt.
  - **En vanlig spelare i provet:** båda webbläsarna står på tornets
    dator och spelar som spelledaren; en spelare utan märket provas inte.
