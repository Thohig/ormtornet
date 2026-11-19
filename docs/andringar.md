# Ändringar

Det som har ändrats i spelet, nyast först. README beskriver hur allt
fungerar. Här står vad som är nytt och vad det betyder för den som spelar
eller leder tornet.

## 2026-11-19

### Automatisk solnedgång och turtid

- Solnedgång och Turtid har flyttat till Fler inställningar på
  värdsidan, efter Tornets storlek. Båda har valet automatisk, som är
  förvalt.
- Automatisk solnedgång följer tornets storlek: 40 rundor på 9 × 9, 50 på
  11 × 11 och 60 på 13 × 13.
- Automatisk turtid följer antalet människor: 30 sekunder för upp till
  fyra, 35 för fem eller sex och 40 för sju eller åtta. Den räknas ut när
  vandringen börjar.
- Ett eget val av solnedgång gäller nu som det står. Förut skalades det
  upp i större torn.
- Bredvid valen står vad automatisk ger just nu. Sammanfattningen av Fler
  inställningar nämner solen och turtiden bara när de inte är automatiska.
- En sparad inställning från förut (till exempel 40 rundor och 30
  sekunder) gäller som eget val. Byt till automatisk en gång på värdsidan.
- Flaggorna `-sol 0` och `-turtid 0` betyder automatiskt.

### Inga droppar

- Droppljudet i bakgrunden är borta. Vinden hörs som förut.

### Teckenförklaringen visar bara tecknen

- Teckenförklaringen uppe till vänster på kartan visar bara tecknen,
  under varandra, så att den skymmer mindre av kartan.
- Vad ett tecken betyder syns när musen är över det, eller efter ett
  klick eller tryck på det, som med tecknen uppe till höger. Ett tryck
  till, eller ett tryck någon annanstans, döljer förklaringen.
- Den är utfälld från början även på en telefon. Den som har fällt ihop
  den har den fortfarande hopfälld.
- Rutan har samma färg som kartan. Den ligger nu under skärmens
  svepande linjer och mörka kanter, som kartan, i stället för över dem,
  så att den inte syns som en ljusare fläck.

### Hela tornet på engelska

- Med Språk: English på värdsidan är allt på engelska: spelarsidan,
  värdsidan, rummen, fällorna, varelserna, klasserna, föremålen, ormen,
  avslöjandet och epilogen. Spelet heter då The Serpent Tower.
- Öppna sidor byter språk utan att laddas om.
- På engelska heter undanmanövern E (evade). U fungerar som förut, på
  båda språken.
- En viskning till tornet räknas på både svenska och engelska.
- Sidornas egna texter ligger nu i `data/spelarsidan.json` och
  `data/vardsidan.json`, och kan ändras där som allt annat innehåll.
- Rundans kort skriver riktningen som texterna vill ("norrut" på svenska,
  "north" på engelska), och beskedet när tornet är fullt säger inte längre
  "Fyra är redan inne".

### Språk: svenska eller engelska

- Värdsidan har en ny inställning, Språk, under Inför nästa vandring, och
  tornet har flaggan `-sprak en`. Bytet gäller genast i lobbyn och annars
  från nästa torn.
- Engelskan läses från `data/en/`, med bara texter. Det som inte är
  översatt visas på svenska. Än så länge är bara svårigheternas namn
  översatta (Easy, Normal, Hard, Nightmare); resten kommer i nästa steg
  enligt [engelska.md](engelska.md).
- `-kontrollera` skriver hur många texter som är översatta.

### Kartritaren och örtkvinnan

- Kartritaren börjar med en dryck.
- Örtkvinnan har inte längre ett liv mindre än de andra. I stället tappar
  hon mycket mer skatt om ormen vaknar medan hon är i kammaren: all
  skatt på Normal och svårare, och 85 % på Lätt. Förut tappade hon en
  fjärdedel av skatten mer än andra.
- I 10 000 simulerade vandringar på Svår kom kartritaren ut med skatt i
  39 % av vandringarna i stället för 30 %, och örtkvinnan i 48 % i stället
  för 34 %. Hon dog i 33 % av dem i stället för 56 %. På Normal ligger alla
  klasser nu mellan 50 och 61 %. På Mardröm dör örtkvinnan något oftare
  än förut, eftersom hon vågar gå längre in.

### Simuleringen jämför klasserna rättvist

- I `-simulera` får AI-spelarna klasserna i tur och ordning, så att varje
  klass, även pilgrimen, spelas lika ofta och lika ofta av den djärva och
  den försiktiga. Förut valde den försiktiga alltid örtkvinnan och
  kartritaren, och tabellen jämförde lika mycket personligheter som
  klasser.
- Utskriften har en ny tabell, Personligheterna, och tabellen Klasserna
  visar andelen som kom ut med skatt för varje personlighet.
- I riktiga vandringar väljer AI-spelarna som förut.

### AI-spelarnas livsgränser följer deras liv

- Hur lite liv en AI-spelare går hem eller dricker vid räknas nu om efter
  dess egna liv. Förut gällde samma antal oavsett klass och svårighet, så
  att en försiktig örtkvinna på Mardröm, med 3 liv, vände redan i porten.
- På Mardröm kommer nu en tredjedel ut med skatt i stället för en
  fjärdedel, men fler dör, omkring 44 %. Mardröm är därmed den farligaste
  svårigheten även för AI-spelarna.

### AI-spelarna hittar ut när de är vilse

- En AI-spelare som har tappat vägen ut, till exempel efter ett fall genom
  golvet, letar sig mot portarna när den bär skatt eller solen börjar
  sjunka. Förut letade den bara mot kammaren och pendlade där tills solen
  gick ner.
- I 10 000 simulerade vandringar på Lätt blev 14 % kvar i mörkret förut,
  nu omkring hälften så många. Fler kommer ut med skatt på alla
  svårigheter.

### Simuleringen skriver "fälla"

- Dödsorsakerna i `-simulera` står som `fälla: pilfällan` i stället för
  `falla: pilfällan`.

### Spelledaren följer en vandrare

- I översikten låter ett klick på en vandrares namn i tabellen kartan
  följa den vandraren, med 5 × 5 rum, ramen runt rummet och pilarna vid
  kanten. Ett klick till visar hela tornet igen.

### Pilar vid kartans kant

- Ligger en port eller kammaren du känner till, eller en vandrare du ser,
  utanför det som syns på kartan pekar en liten pil vid kanten åt det
  hållet, med samma tecken och färg som på kartan. Pilarna visar bara
  det du redan känner till.

### Ditt rum syns tydligt

- Rummet du står i har ljusare golv och en tunn ram i din färg. Inget
  annat på kartan har blivit mörkare.

### Kartan följer snabbare

- När du går till ett nytt rum glider kartan efter på 0,4 sekunder. Den
  långa, mjuka zoomen blir kvar när du byter nivå.

### Zooma kartan

- Kartan har fyra nivåer: Nära (3 × 3 rum), Mellan (5 × 5), Långt (9 × 9)
  och Allt du känner till. Byt med + och −, scrollhjulet, två fingrar
  eller knapparna nere till höger på kartan.
- M, eller ett klick på nivåns namn, växlar mellan din nivå och Allt.
  Ett tryck på själva kartan gör inget längre, så att den inte byts av
  misstag.
- En smal skärm börjar på Nära och en större på Mellan. Valet sparas på
  enheten.

### Kartan följer ditt rum

- Kartan har alltid rummet du står i i mitten, med två rum åt varje håll
  (ett på en smal skärm). Förut visade den allt du kände till och följde
  dig bara när tornet blev för stort för skärmen.
- Vid tornets kanter går rutan bara lite utanför tornet, så att portarna
  syns och ett rum i hörnet inte hamnar under knapparna och tecknen.

### Teckenförklaringen går att läsa

- Teckenförklaringen ligger över skärmens svepande linjer och mörka
  kanter, så att hörnet där den står inte längre gör texten suddig och
  mörk.
- Texten är lite större, står i den vanliga textfärgen i stället för
  den grå kritfärgen, och bakgrunden är mörkare.

### Ingen dämpad text

- Teckenförklaringen på kartan är inte längre dämpad. Den syns med full
  styrka hela tiden, och rubriken "Tecken" står i samma ljusa färg som
  förklaringarna.
- Tecknen uppe till höger på kartan, frågetecknet och rummets namn syns
  också med full styrka hela tiden.
- Text som stod i den dämpade bruna färgen på spelarsidan står nu i den
  vanliga textfärgen: titeln och tiden i översta raden, rundornas
  rubriker och de dämpade raderna i loggen, de extra raderna i panelerna,
  kortet, tipsen i styrningen, texterna i klassvalet och ryktet och
  uppmaningen i introt. Det utgråade, som frånkopplade vandrare och tagna
  kort, är kvar.
- På värdsidan står rubrikerna, etiketterna och förklaringarna vid
  inställningarna, raden överst, tabellen och vandrarnas logg också i
  den vanliga textfärgen. Knappar och val som inte går att använda just
  nu är fortfarande utgråade.

### Översikten i två kolumner

- På en bred skärm står kartan till vänster och ingreppen till höger i
  översikten, med tabellen över vandrarna under. På en smal skärm står
  de under varandra som förut.

### Fler inställningar

- Under Inför nästa vandring syns Svårighet, AI-spelare, Solnedgång,
  Turtid och Spelledarläge direkt.
- Minst antal, Avslöjandet, Frånvarande vandrare, Strid, Vandrarna syns,
  Tornets storlek och Klasser ligger under Fler inställningar. Raden visar
  de nuvarande valen i korthet.
- Inställningarna är låsta under vandringen som förut, så det behövs
  ingen märkning av vad som gäller från nästa omgång.

### Värdsidan följer fasen

- Avsnitten på värdsidan byter plats efter fasen. I lobbyn och
  klassvalet kommer Bjud in och Vid porten först, under vandringen
  översikten och när vandringen är slut Inför nästa vandring.
- Under vandringen krymper Bjud in till adressen och Kopiera adressen.
- Kommer någon inte in? och Läs om innehållet ligger hopfällda längst ner,
  och Stäng tornet är en liten knapp under dem. Hjälpen fälls ut av sig
  själv om ingen människa har kommit in efter en och en halv minut.

### En rad överst på värdsidan

- Överst på värdsidan står en rad som följer med när sidan rullar: fasen,
  rundan och solen, hur många som är här och vilka som är borta.
- Där ligger Starta nu före vandringen, och Pausa eller Fortsätt under
  vandringen i spelledarläget.

### Vandrarnas logg per runda

- Vandrarnas logg i översikten är hopfälld tills spelledaren klickar på
  den.
- Den är grupperad per runda med den senaste överst, som loggen på
  spelarsidan. De två senaste rundorna står öppna och de äldre går att
  fälla upp.
- Under Visa kan loggen visa bara en vandrare.

### Senaste rundan överst i loggen

- I loggen på spelarsidan står den senaste rundan överst. Inom en runda
  står raderna fortfarande i den ordning de hände.
- Äldre rundor fälls ihop under de två senaste. En ljusare linje skiljer
  den aktuella rundan från de äldre.
- Har man rullat ner för att läsa ligger texten kvar när en ny runda
  kommer, och "↑ nya rader" tar en tillbaka upp.

### Läs om innehållet

- Under Inför nästa vandring på värdsidan finns knappen "Läs om
  innehållet". Den läser och kontrollerar `data/` utan att tornet startas
  om. I lobbyn gäller det nya genast och tornet byggs på nytt. Under och
  efter en vandring gäller det från nästa torn.
- Är något fel i filerna står felen under knappen, och tornet behåller
  det gamla innehållet.

### Simulera många vandringar

- `go run . -simulera 500 -ai 4` spelar 500 vandringar med bara
  AI-spelare på under en sekund och skriver hur det gick: hur många som
  kom ut med skatt, dog eller blev kvar, vad de dog av, hur ofta ormen
  vaknade och hur det gick för varje klass. Till för att justera
  siffrorna i `data/`.
- Samma `-fro` ger samma resultat.

### Automatiska kontroller på GitHub

- Vid varje push kör GitHub formateringen, `go vet`, testerna med
  race-detektorn, kontrollen av datafilerna och en kontroll av
  JavaScript-filerna (`.github/workflows/kontroll.yml`).
- `./ormtornet -kontrollera` kontrollerar datafilerna och avslutar, utan
  att öppna tornet.

### Frö för varje vandring

- Varje vandring har ett frö som bestämmer slumpen: rummen, dörrarna,
  fällorna, varelserna, striderna, ormkorten och AI-spelarnas val. Fröet
  står i terminalen när vandringen börjar och på värdsidan under Vid
  porten.
- Med `-fro 1234` får den första vandringen det fröet. Gör vandrarna samma
  drag blir tornet detsamma, rum för rum, så att en konstig vandring kan
  spelas upp igen. Nästa vandring får ett nytt frö.
- Texterna väljs fortfarande fritt, så samma frö kan ge andra ord.

### Kontroll av datafilerna

- Tornet kontrollerar `data/` när det startar. Ett felstavat fält, en
  chans utanför 0–1, ett id som finns två gånger, en vana som inte är
  anfall, parera eller undan, eller en klass i `ai.json` som inte finns
  stoppar starten med ett meddelande som säger fil och fält. Förut
  ignorerades ett felstavat fält tyst.
- Fält som börjar med `_` är kommentarer, som förut.
- Ett test kontrollerar att alla texter som koden använder finns i
  `data/`, så att ingen `[nyckel]` dyker upp i spelet.

### Varelserna vandrar

- Varelserna står inte längre still i rummet där de dök upp. I början av
  varje runda kan en varelse gå genom en öppen dörr till ett grannrum,
  också till ett rum som ingen har hittat än. Där väntar den tills någon
  går in, och rummet syns inte på kartan förrän dess.
- En varelse som står hos en vandrare stannar och strider, och ingen
  varelse vandrar på egen hand in till en vandrare. Portrummen och
  kammaren går de aldrig in i.
- Står du i ett grannrum kan du höra varelsen röra sig (`ljud_varelse` i
  `data/texter.json`), när vandrarna inte ser varandra.
- Den döda spelaren som styr varelserna kan också föra dem ut i det okända.
- Hur ofta: `vandra_chans` i `data/monster.json` (1 på 4), och en egen
  `vandrar` för varje sort. Benväktaren rör sig nästan aldrig, gallgasten
  och blekkryparen ofta.

### En död människa tar över varelserna från en AI-spelare

- Förut styrde den första som dog varelserna resten av vandringen, även om
  det var en AI-spelare. Nu tar den första människa som dör över
  varelserna från AI-spelaren. Dör en annan människa sedan är det
  fortfarande den första som styr.

### Människorna väljer klass först

- I klassvalet väntar AI-spelarna tills alla människor har valt. Den som
  spelar mot bara AI-spelare kan alltså ta vilken klass som helst, och
  AI-spelarna tar sedan det de helst vill ha av det som är kvar.
- Om tiden tar slut får människorna sina slumpade kort före AI-spelarna.

### Teckenförklaring på kartan

- Uppe till vänster på kartan, under rummets namn, står nu vad tecknen i
  rummen betyder: vandrarna, varelserna och fällorna med sina namn,
  kammaren, portarna, kritpilar, ekon, föremål och okända rum. Bara de
  tecken som står på kartan just nu finns med, så listan växer under
  vandringen.
- Ett klick på "Tecken" fäller ihop och upp den, och valet sparas. På en
  liten skärm är den hopfälld från början.
- Texterna står under `tecken_*` i `data/texter.json`.

### Vandrarna i översta raden och i loggen

- Under vandringen står vandrarna som prickar i sin färg i översta raden,
  bredvid solen. En fylld prick har valt drag (i lobbyn: är redo), en tom
  väljer fortfarande, ✕ är död och ◌ har gått ut. Den egna är
  understruken. Håll musen över en prick, eller tryck på den, så står
  namn, klass och läge, till exempel "sågs runda 12".
- Panelen Vandrare syns bara i lobbyn och i slutet. Under vandringen får
  loggen platsen.
- I loggen står vandrarnas namn i sin färg, så att det syns vem en rad
  handlar om.

### Rundans sammanfattning

- Efter varje runda kommer ett kort nere på kartan med vad som hände dig:
  vart du gick eller vad du gjorde, liv och guld du vann eller förlorade,
  drycker och föremål du hittade, fast, gift, mörker, en varelse i rummet
  och om ormen vaknade, och sist solen. Det farliga står först, med `!`
  och i varningsfärg.
- Kortet försvinner när du väljer nästa drag, efter sju sekunder eller
  vid ett klick. En kort rad står kvar i loggen, till exempel "Runda 12:
  österut, −1 liv, +10 guld".
- Servern jämför vandrarens läge före och efter rundan. Texterna står i
  `data/texter.json` under `sam_`, med en kort form under `sam_*_kort`.
- Loggen hamnar inte längre efter när rundorna går fort: rader från en
  runda som redan är förbi skrivs ikapp fort.

### Lugnare logg

- Bara den första raden i en runda skrivs fram tecken för tecken. Resten
  kommer en hel rad i taget, ungefär tre rader i sekunden.
- Varje runda har en rubrik, "Runda 12", i stället för ett streck. Före
  vandringen och slutet har egna rubriker.
- Har man rullat upp för att läsa stannar loggen där, och knappen "↓ nya
  rader" tar en tillbaka ner. Loggen har en smal rullningslist.
- När en ny runda börjar fälls rundorna före den förra ihop till en rad,
  till exempel "▸ Runda 4 · Salpetergrottan · 8 rader". Ett klick fäller
  upp den igen. Loggen sparar upp till 300 rader.

### Statusen på kartan och loggen till höger

- Sidopanelen Du är borta. Klassen, liv, skatt, drycker och tillstånd står
  som halvt genomskinliga tecken uppe till höger på kartan, och föremålen
  under dem. Håll musen över ett tecken, eller tryck på det på en
  pekskärm, så står det vad det betyder.
- Klasserna har fått ett tecken, `tecken` i `data/klasser.json`. Det syns
  på kartan och på korten i klassvalet.
- Rummets namn står uppe till vänster på kartan i stället för i panelen
  Här.
- Styrningen ligger bakom ett **?** på kartan, och tangenterna `?` och `H`.
  Den visas en gång första vandringen.
- Loggen fyller högerkolumnen under panelerna Vandrare och Nu, så att
  många fler rader syns. På en telefon är den högre än förut.
- Panelen Nu visar det som gäller just nu: draget, striden, ormen och vad
  som väntar.
- Introt visar "Gå in i tornet" efter 5 sekunder i stället för 15, precis
  efter en blixt.

### Klassvalet på en egen skärm, och kartan växer fram

- Klassvalet är inte längre i lobbyn. När alla är redo kommer en egen
  skärm med ett spelkort för varje klass: namn, beskrivning, liv och
  drycker, styrka och svaghet. Välj med 1–8 eller ett klick. Kort som
  andra har tagit är gråa med namnet på den som tog dem.
- Klassvalet har en tidsgräns, 30 sekunder (`klassval_sekunder` i
  `data/klasser.json`). Den som inte har valt får ett slumpat ledigt kort,
  och vandringen börjar när alla har en klass.
- Spelledarens Starta nu går från lobbyn till klassvalet, och under
  klassvalet börjar vandringen direkt. Värdsidan visar vem som har valt.
- Kartan visar bara det vandraren känner till. I början är det rummet man
  börjar i, och kartan zoomar ut mjukt efter hand som tornet växer fram.
  Okända rum syns som prickar bara bredvid de kända.
- De andras portar är okända från början, även för Kartritaren och med
  de andra synlighetsvalen, tills någon går in i dem.
- Kartritaren ser genom dörrarna in i rummen bredvid det hen står i. De
  kommer bara med på hens egen karta. En fälla där väntar fortfarande på
  den som går in.
- Ett rum räknas som hittat först när någon går in i det. Ett rum man bara
  har sett in i, med Kartritaren eller spegeln, är nytt för den som går in
  och har kvar sin fälla. Förut gick en fälla inte av i ett rum som någon
  hade sett i spegeln.

### Tydligare paneler

- Sidopanelen är delad i små skiffertavlor med ram, som kartan: Vandrare,
  Här, Du och Styrning.
- Liv, solen och fiendens liv är mätare med rutor och ett tal, till
  exempel 4/6. Rutorna byggs i kod och krymper när de är många.
- Skatten står i guld, och det valda draget står sist under Du, avskilt
  med ett streck.
- Loggen har ett tecken i marginalen för vad slags rad det är: `!` fara,
  `›` det du ser, `▸` rummets namn, `·` besked, `❝` tornets röst och `✎`
  krita. Ett streck skiljer rundorna åt, och rader från äldre rundor än
  den förra bleknar.

## 2026-11-18

### Kartan i krita

- Kartan är krita på en skiffertavla. Väggarna är handritade streck med
  luckor för dörrarna, och kända rum har en svag kritdimma som golv.
- Vandrarna är prickar i sin färg, och den egna har en ring som blinkar.
  Varelser, fällor och märken är tecken som förut, men i fet stil.
- Kammaren är en kritring kring guldet. När ormen vaknar spänns ringen
  och ögonen lyser.
- Avslöjandet syns i kritan: på nivå 2 sväller strecken, på nivå 3 får de
  revben och blir köttfärgade, på nivå 4 är de röda och porten en mun.
- Med egen karta (förvalet) syns bara den egna porten, redan i lobbyn.
  De andra portarna hittar man när man går dit. Förut syntes alla fyra
  från början.

### Efter provspelet

- Nytt torn leder tillbaka till lobbyn, och alla väljer klass på nytt
  bland två nya. Förut följde klassen med till nästa vandring.
- Tornet väntar när ingen människa är ansluten. Förut gick rundorna så fort
  AI:n hann välja, så ett nätavbrott kunde kosta halva vandringen.
- Tab heter "byt mål" i styrningen, så att raden får plats.
- README, `docs/fler-spelare.md` och `docs/spelledaren.md` beskriver det
  nya. Introts första 15 sekunder ligger under "Prova med människor" i
  `docs/att-gora.md`.
- Tornet kraschade om en vandrare kom tillbaka just när AI:n hade tagit
  över och skulle välja drag. Nu tar människan över utan att något händer.
- Hjälptexten för `-ai` säger 0–7, eller 8 i spelledarläget, som README.

### Viska hemligheten

Den som vet tornets hemlighet kan viska den till tornet med V, från nivå 2
och en gång per nivå. Rätt viskning gör vandraren invigd: hen hör åt
vilket håll kammaren ligger, kliver förbi hälften av fällorna, kan lugna
ormen en gång, får mer guld ute och porten väntar en runda för hen. Men
varelserna anfaller oftare och ormen vaknar lättare. Tornet säger aldrig
vad hemligheten är. Spelledaren ser viskningarna och kan inviga. Med
avslöjandet avstängt går det inte att viska. README beskriver det under
"Att viska hemligheten", i avsnittet med spoilervarning.

### Föremål, steg 2 och 3

- Fem föremål används med 1–3: rökpulver, spegelskärva, kastkniv, krita
  och timglas. Den sövande flöjten verkar av sig själv. AI-spelarna
  använder rökpulver, kastkniv och timglas.
- Tre förbannade föremål: den svarta pärlan, myntet med hål och kronan av
  rost. "Kartan som ljuger" från förslaget blev kronan i stället, eftersom
  en karta med falska rum skulle kräva rum som inte finns.
- Brev från tidigare vandrare, efter avslöjandets nivå.
- Den som dör tappar sina föremål i rummet.
- Klasserna: tjuven hittar oftare föremål, kartritaren får fler streck med
  kritan, spårfinnaren ser var föremål ligger och pilgrimen låter
  förbannade föremål ligga.
- Spelledaren kan ge föremål. Valet av föremål på värdsidan står kvar när
  översikten ritas om.

### Föremål, steg 1

När man söker igenom ett rum kan man hitta ett föremål. Man bär högst tre,
och ett som inte får plats blir kvar i rummet åt den som kommer efter.
De sex första föremålen verkar av sig själva: rep, fjällamulett, extra
veke, bucklig sköld, väska med dubbel botten och kompass. README beskriver
dem under "Föremål". De står i panelen och i spelledarens översikt.

AI-spelarna söker sällan och hittar i snitt 0,13 föremål per vandring, så
balansen för dem är densamma. Människor som söker mer får mer nytta av dem.
Steg 2 (föremål man använder själv) och steg 3 (förbannade föremål, brev,
klasserna) är inte byggda.

### Upp till 8 spelare

Tornet tar emot upp till 8 vandrare, AI-spelarna inräknade, och växer med
dem:

| Vandrare | Torn | Portar | Sol (med 40 på värdsidan) |
|---|---|---|---|
| 2–4 | 9 × 9 | 4 i hörnen | 40 |
| 5–6 | 11 × 11 | 8 | 50 |
| 7–8 | 13 × 13 | 8 | 60 |

- Tornet byggs om i lobbyn när någon kommer eller går, så att det alltid
  har rätt storlek. Den som redan är redo förblir det.
- I ett torn på 11 × 11 eller mer finns fyra nya portar mitt på varje sida,
  med tre dörrar inåt och en utgång rakt ut. De heter norr, söder, väster
  och öster i texterna.
- Solen räknas från värdsidans val för ett torn på 9 × 9 och växer lika
  mycket som vägen till kammaren.
- Ormleken får två extra sovkort per vandrare över fyra
  (`sovkort_per_vandrare_over_fyra` i `ormen.json`), så att risken växer i
  samma takt som med fyra.
- Spelledaren väljer 0–7 AI-spelare, eller 8 i spelledarläget, och minst
  1–8 vandrare. Värdsidan visar tornets storlek och solen under Vid porten.
- Väljer spelledaren en storlek blir tornet aldrig mindre än vad antalet
  vandrare kräver.
- Åtta färger och åtta platser i varje rum. AI-spelarna har tio namn.
- På en liten skärm visar kartan 9 × 9 rum runt den egna vandraren när
  tornet är för stort för att läsas. Ett tryck på kartan visar hela tornet,
  och ett tryck till går tillbaka.

### Tornets storlek

Spelledaren kan välja "Tornets storlek: 9 × 9, 11 × 11 eller 13 × 13" på
värdsidan, eller `-storlek`. Förvalet "automatisk" ger 9 × 9 tills vidare;
när upp till 8 spelare kommer räknas storleken ut efter antalet vandrare.
Ett större torn har längre väg till kammaren (10 och 12 steg i stället för
8), så solen kan behöva fler rundor. Byts storleken i lobbyn byggs ett nytt
torn direkt.

På en telefon blir rummen små i ett torn på 13 × 13.

### Klasser

Vandrarna väljer klass i lobbyn, en av två slumpade, i Talismans anda.
Åtta klasser med var sin styrka och svaghet: Legoknekten, Tjuven,
Örtkvinnan, Kartritaren, Ormtjusaren, Spårfinnaren, Lyktbäraren och
Pilgrimen. README beskriver dem under "Klasser". Spelledaren kan stänga av
klasserna på värdsidan.

Balansen är provad med AI-spelare: med samma klass för alla fyra når
13–35 % kammaren, mot 23 % utan klass. AI-spelarna använder inte alla
styrkor (till exempel tjuvens tysta näve och kartritarens karta), så
människor kan få andra siffror. Siffrorna justeras i `data/klasser.json`.

Rättat: spelledarens översikt kraschade under en vandring med egen karta.

### Ingen vet var de andra börjar

Gäller när vandrarna inte ser varandra (`egen` och `moten`).

- Lobbyn döljer de andra vandrarna på samma sätt som vandringen. Förut såg
  alla i lobbyn var de andra stod, och visste därför var de skulle börja.
- Portarna lottas när vandringen börjar. Var och en får veta sin egen port
  ("Du kliver in genom porten i tornets sydöst"), men inte de andras.
  Ordningen i listan och färgen säger inte längre vilken port någon har.
- Texterna när någon kommer till tornet nämner inte porten.
- Alla fyra portrummen syns på kartan från början, för alla och med alla
  val av synlighet. Man vet var utgångarna är, men inte vem som kom in var.
- Den som står i ett grannrum till en port hör den gnissla när någon går ut.

Med "alltid, för alla" börjar var och en vid porten som hör till dess plats
i listan, som förut.

Koden finns i `game.go` (`drawCorners`, `revealPorts`, `Corner` på
vandraren). Texterna står under `ansluter_dolt`, `valkommen_dolt`,
`start_port` och `ljud_port` i `texter.json`, och under `ansluter_dolt` och
`lamnar_dolt` i `ai.json`.

### Introt är en storm

Startsidan visar tornet i en storm vid solnedgången.

- Moln driver över solen, det regnar och blixtrar med 8–20 sekunders
  mellanrum. I blixten blir tornet en svart siluett mot himlen.
- Ett klick var som helst, eller Enter, startar ljudet: en låg drönarton i
  moll med en tritonus som sväller in, vind, regn, åska efter blixten och
  ibland ett dunk långt inne i tornet. Webbläsaren tillåter inget ljud
  innan man har klickat, så utan klick är introt tyst.
- Ryktena fladdrar lite, och ögat i tornet glöder ibland efter en blixt.
- "Gå in i tornet" visas först efter 15 sekunder, precis efter en blixt,
  och då börjar porten glimma. Före det startar klick och Enter bara
  ljudet.
- Den som redan har sett introt i fliken, och spelledaren på tornets dator,
  får knappen direkt.
- Med `-avslojande=false` finns stormen men inte ögat och dunket.
- Den som har valt minskade rörelser i datorn får en mjuk ljusning i
  stället för blinkande blixtar, och regnet står stilla.

Ändra i `data/texter.json`: väntetiden under `intro` → `vanta` och
uppmaningen att klicka under `intro_ljud`. Koden finns i `web/app.js`
(`Intro`), `web/audio.js` (stormen och åskan), `web/style.css` och
`web/index.html`.

### Vandrarna ser inte varandra

En ny inställning på värdsidan, "Vandrarna syns", med flaggan `-synlighet`.

| Val på värdsidan | Flagga | Vad vandrarna ser |
|---|---|---|
| bara när de möts, egen karta | `egen` | andra vandrare bara i samma rum, och bara rum de själva har gått igenom |
| bara när de möts | `moten` | andra vandrare bara i samma rum, men kartan är gemensam |
| alltid, för alla | `alla` | allt, som förut |

Förvalet är `egen`. Den som har sparade inställningar från förut får också
`egen`, eftersom de inte innehåller något val för synlighet.

Med `egen` och `moten`:

- Där du senast såg en annan vandrare står ett svagt `@` i dess färg, och i
  listan över vandrare står "sågs runda 12" eller "inte sedd".
- Tornet berättar fortfarande om fällor och om källaren, men bara i vilken
  del av tornet och inte vem. Att någon går ut eller dör får alla veta.
- Ekofällan visar som förut var den drabbade är, för alla.
- **Spår:** i ett rum där någon annan har gått ut de senaste 3 rundorna ser
  du färska fotspår och åt vilket håll de leder. I mörker syns inga spår.
- **Ljud genom väggarna:** står du i ett grannrum hör du strid, fällor och
  någon som dricker, med riktningen dit men utan namn.
- Den döda spelaren som styr varelserna ser alla vandrare. När vandringen
  är slut syns alla igen.

Med `egen` dessutom:

- Varje vandrare har sin egen karta, med fällornas tecken och varelserna i
  de rum den har varit i. Rummen finns ändå bara en gång i tornet, så en
  fälla slår bara till för den som först går in.
- Ljudet för ett nytt rum hör bara den som gick in, och tornet berättar
  inte om andras nya rum.
- Den döda spelaren som styr varelserna ser hela tornet så långt det har
  vuxit fram, och när vandringen är slut ser alla hela tornet.

AI-spelarna ser spelet genom samma vy som människorna och har därför
samma begränsningar. Spelledarens översikt ser allt som förut.

Ändra i `data/`: de vaga texterna har nycklar som slutar på `_dolt` i
`texter.json` och `fallor.json`, och fällorna har `andra_dolt`. Spåren och
ljuden står under `spar`, `spar_ut` och `ljud_*` i `texter.json`, och hur
länge spåren syns under `spar_rundor` i `fallor.json`. Förslaget och
besluten står i [synlighet.md](synlighet.md).

Koden finns i `game.go` (`sees`, `knowsRoom`, `noteRoom`), `sounds.go`
(spår och ljud), `settings.go` och `web/app.js` (senast sedd).
