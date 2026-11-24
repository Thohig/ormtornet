# Ormtornet

En dungeon crawler för 2–8 spelare i webbläsaren. Spelledaren öppnar
<https://thohig.github.io/ormtornet/> och skickar länken till de andra.
Inget behöver installeras, och alla kan sitta på olika nät (se Starta).
Mitt i ett torn på 9x9 rum (11x11 för fem eller sex, 13x13 för sju eller
åtta) sägs en jättelik orm sova kring en skatt.
Ta så mycket du vågar och kom ut innan solen går ner.

Tornet kan spelas på svenska eller engelska (The Serpent Tower); se Språk.

Vad som är nytt står i [docs/andringar.md](docs/andringar.md).

## Starta

Spelet spelas från en länk på GitHub Pages. Ingen behöver installera
något eller vara på samma wifi.

1. Spelledaren öppnar <https://thohig.github.io/ormtornet/> i Chrome,
   Edge eller Safari på en dator, Mac eller Windows.
2. Introt leder till värdsidan. Där kopierar knapparna under Bjud in en
   inbjudan med länken, att skicka till spelarna.
3. Spelarna öppnar länken på en dator, telefon eller surfplatta, var de
   än är. Spelledaren spelar med genom "Gå in själv".

Tornet körs i spelledarens flik, med samma spel som programmet man kan
ladda ner. Spelarna ansluter direkt till spelledarens webbläsare. Den som
har länken kommer in, så ge den bara till dem som ska spela.

- Tornets flik måste vara öppen hela vandringen. Stängs den eller laddas
  den om stängs tornet. "Gå in själv" byter sida i fliken, så spelledaren
  spelar i tornets egen flik. Fliken får också ligga i bakgrunden.
- Avsluta med "Stäng tornet" på värdsidan innan du stänger fliken. Då får
  spelarna beskedet att tornet stänger. Stängs bara fliken tror deras
  sidor att anslutningen bröts och börjar leta efter tornet igen.
- En dator är bäst. På en telefon kan tornet stanna när skärmen släcks.
- Den som tappar anslutningen öppnar länken igen i samma webbläsare och får
  tillbaka sin vandrare.
- Vissa nät stoppar anslutningen, oftast mobilnät och nät på jobb och
  skolor. Det finns ingen reläserver (TURN) än; prova då ett annat nät.
- Spelarna kan se spelledarens publika IP-adress, som i ett videosamtal.
- Inställningarna sparas i spelledarens webbläsare. Läs om innehållet
  läser bara det innehåll som följde med bygget; ändringar i `data/` kommer
  ut när de pushas till `main`.
- Kommer ni inte in alls går tornet att ladda ner och köra på en Mac (se
  Ladda ner och kör tornet själv).

Hur det är byggt står under Tornet i webbläsaren och i
[docs/github-pages.md](docs/github-pages.md).

### Värdsidan

Värdsidan är det spelledaren ser i tornets flik efter introt. Där står
länken spelarna ska öppna, med knappar som kopierar en färdig inbjudan att
skicka i en gruppchatt. Sidan visar också vilka som har anslutit och är redo,
och vad man kan göra om någon inte kommer in.

Överst står en rad som följer med när sidan rullar: lobbyn, klassvalet,
rundan och solen, eller att vandringen är slut, hur många människor och
AI-spelare som är här och vilka som är borta. Där ligger också knappen som
behövs just nu: Starta nu före vandringen.

Avsnitten står alltid i samma ordning: Bjud in, Vid porten, Inför
nästa vandring och Fröet. På en skärm som är minst 1100 punkter bred står de i två
kolumner under raden: Bjud in och Vid porten till vänster, Inför nästa
vandring och resten till höger, med förklaringen under varje val. På en
smalare skärm är sidan en kolumn. Under Bjud in står inte adressen (eller
länken) utskriven: knapparna kopierar inbjudan och adressen, och Visa
adressen visar den vid behov. Går det inte att kopiera visas den av sig
själv. Under och efter vandringen krymper Bjud in till knapparna för
adressen, så att resten kommer upp. Längst ner ligger Kommer någon inte in? och Läs om
innehållet, hopfällda, och en liten knapp som stänger tornet. Har ingen
människa kommit in en och en halv minut efter att sidan öppnades fälls
Kommer någon inte in? ut av sig själv.

Under Fröet står den här vandringens frö, under vandringen och i slutet,
med Kopiera och, i slutet, Samma igen, som ger nästa vandring samma frö.
Fältet Nästa vandring visar alltid fröet nästa vandring får: tornet
slumpar det i förväg, Slumpa drar ett nytt, och ett eget frö (1–999 999
999) skrivs in och sparas med Använd eller Enter. Fröet gäller bara nästa
vandring; vandringen efter får ett nytt. I lobbyn byggs tornet om med
fröet, så att också portarna blir desamma. Samma frö ger samma torn när
lika många vandrar, med samma svårighet och samma drag.

Under Inför nästa vandring väljer spelledaren det här. Under formuläret förklarar "Vad inställningarna gör", hopfälld, varje
inställning med några meningar (`forkl_*` i `vardsidan.json`).

| Inställning | Val | Förval |
|---|---|---|
| AI-spelare | 0–7 | 0 |
| Svårighet | Lätt, Normal, Svår, Mardröm | Normal |
| Frånvarande vandrare | AI tar över efter 1–3 rundor, eller står still | 2 rundor |
| Språk | svenska eller English | svenska |

Resten går inte att välja. Tornets storlek följer antalet vandrare: 9 × 9
rum för upp till fyra, 11 × 11 för fem eller sex och 13 × 13 för sju eller
åtta. Solen går ner efter 40 rundor i ett torn på 9 × 9, 50 på 11 × 11
och 60 på 13 × 13, lika många fler som vägen till kammaren är längre, och
en fjärdedel fler på Lätt.
Turtiden är 30 sekunder för upp till fyra människor, 35 för fem eller sex
och 40 för sju eller åtta. AI-spelare räknas inte. Turtiden räknas ut när
vandringen börjar och ligger fast under den. Tillägget är litet med flit:
rundan avgörs så fort alla har valt, och ett högre tak låter bara alla
vänta längre på den som dröjer. Vandrarna väljer alltid klass, de kan
alltid slåss med varandra utom i portrummen, och vandringen börjar av sig
själv när minst två är redo, AI-spelare inräknade. Kartan ritas snett
ovanifrån i 3D; en enhet utan WebGL visar en text i stället för kartan.

Under Vid porten står hur stort tornet blir och hur många rundor solen
har. AI-spelare kommer och går direkt i lobbyn, och tornet byggs om i
lobbyn när antalet vandrare ändras. Resten gäller från nästa omgång. Under en pågående vandring
är inställningarna låsta. Om tornet är fullt när en människa ansluter i
lobbyn kliver en AI-spelare åt sidan, och kommer tillbaka om platsen blir
ledig igen. Knappen Starta nu börjar vandringen även om inte alla har
tryckt Enter.

Om en spelare lämnar mitt i en vandring tar en AI-spelare över vandraren
efter två rundor (ändras under "Frånvarande vandrare", eller sätts till "står
still"). Spelledaren kan också trycka "AI tar över" direkt. AI:n fortsätter
med vandrarens plats, skatt och liv, och de andra ser namnet märkt "(AI)".
Kommer spelaren tillbaka i samma webbläsare får den genast tillbaka sina
steg. I resultatet märks den som AI har spelat åt med "~". Den som inte kommer
tillbaka innan vandringen är slut är inte med i nästa.

Detsamma gäller den som är kvar men står still: låter tiden gå ut utan
att välja något drag lika många rundor i rad (Vänta är ett drag) tar
AI:n över, och panelen Nu säger att tornet styr stegen. Nästa drag
spelaren själv väljer ger tillbaka vandraren direkt.

Är ingen människa ansluten alls, till exempel när nätet försvinner för
alla, väntar tornet: rundan står pausad tills någon kommer tillbaka, så att
solen inte hinner gå ner medan AI:n spelar ensam. Raden säger att rundan
är pausad.

Värdsidan visar inga skatter, positioner eller drag, så spelledaren kan spela med.
"Gå in själv" byter till spelarsidan i samma flik, utan introt. Där står
"spelledare" i översta raden, och i lobbyn och slutet en knapp Värdsidan;
i slutet finns också Till värdsidan bredvid Nytt torn. Under vandringen
öppnar ett tryck på "spelledare" värdsidan i en ny flik, så att
spelarsidan är kvar och vandraren inte räknas som frånvarande; Till
spelet i den fliken går tillbaka till spelets flik. Har spelledaren en
vandrare under vandringen heter knappen på värdsidan Till spelet. På
GitHub Pages byter ramen i skalet sida, så tornets flik står kvar, och
under vandringen är "spelledare" bara ordet, eftersom den som går till
värdsidan då räknas som frånvarande. I ett nedladdat torn ligger värdsidan
på `http://localhost:8080/spelledare`, och andra datorer i nätverket
kommer inte åt den.

## Spela

1. Startsidan visar tornet i en storm vid solnedgången medan rykten om det
   tonar fram. Moln driver förbi, det regnar och blixtrar, och efter varje
   blixt mullrar åskan. Klicka var som helst, eller tryck Enter, för att
   höra stormen: en låg, hotfull ton, vind, regn och åska. Webbläsaren
   tillåter inget ljud innan man har klickat. Efter 5 sekunder tonar
   "Gå in i tornet" fram efter en blixt, och porten börjar glimma. Klicka
   på den eller tryck Enter. Den som redan har sett introt i fliken, och
   spelledaren på tornets dator, får knappen direkt. Ryktena står i
   `data/texter.json` under `rykten`. Väntetiden står under `intro` → `vanta`
   och uppmaningen att klicka under `intro_ljud`.
2. Tornet ger dig ett namn, till exempel "Dagny Ekfot", och loggen säger
   "Tornet kallar dig Dagny Ekfot." Ingen skriver något namn.
3. Tryck Enter, eller Redo i översta raden, när du är redo. Vandringen
   börjar när alla är redo.

### Vad en spelare ser

En spelare som ansluter från en annan dator eller telefon ser bara
spelarsidan (`/`), med två vyer i tur och ordning, och krönikans skärm för
den som dör och för alla i slutet (se slutet nedan och Strider):

1. **Introt**: tornet i stormen, ryktena och, efter en stund, "Gå in i
   tornet". Hela övergången in genom porten visas bara första gången i
   fliken.
2. **Spelet**, där spelaren stannar resten av kvällen. Det går igenom fyra faser:
   - **lobbyn**: vandrarna som prickar med namn i översta raden, och
     knappen Redo där (Enter, eller Redo längst ner på en pekskärm),
   - **klassvalet**: en egen skärm med ett kort för varje klass, i samma
     stil som spelkorten, med bild, styrka och svaghet. Korten ligger på
     baksidan och vänds upp ett i taget. Välj med 1–8 eller ett klick. Ett
     kort som någon annan har tagit vänds tillbaka och visar baksidan med
     namnet; när du har valt tonas de andra ned. Den som inte har valt när
     tiden är slut får ett av de lediga korten,
   - **vandringen**: kartan, solmätaren, vandrarna som prickar i sin färg
     i översta raden, panelen Nu som en rad under kartan, och loggen. På en
     bred skärm fyller kartan hela bredden och loggen står på kartan, nere
     till vänster (se nedan). En prick är fylld när vandraren har valt drag
     (i lobbyn när den är redo), ✕ är död och ◌ ute; håll musen över den,
     eller tryck, för namn, klass och läge. I lobbyn och slutet står namnet
     bredvid pricken, och knappen Redo (Nytt torn i slutet) står till
     höger i översta raden. Panelen Nu syns bara när den har något att
     säga. Loggen är grupperad under Före vandringen (lobbyn), Vandringen
     börjar och Runda 1, 2 … Det första rumsblocket säger vid vilken port
     du kliver in, hur högt solen står och att portar knarrar långt borta
     när de andra går in. I loggen står vandrarnas namn i deras färg. Överst i
     varje runda står rummet som löpande text, utan ruta och utan namn:
     beskrivningen när
     rummet är nytt, och det man ser där, beskrivet och inte förklarat:
     varelsens egen text, de andra vandrarna, guldet, porten, spåret
     efter en fälla som har slagit till (utan fällans namn) och föremål
     som någon har lämnat. Går man in i ett känt, tomt rum står bara
     "Du känner igen rummet:" och rummets korta beskrivning (`har_igen`),
     så att man ser att draget blev av. Står man kvar
     kommer bara det som har ändrats, som "Gestalten i fjällrustning är borta."
     En varelse som någon annan fällde "ligger fälld på golvet"; den som
     fällde den har redan läst det. Spåret efter en fälla som nyss slog
     till mot en själv står inte i blocket, eftersom loggen säger vad som
     hände. Sist
     står det man hörde under rundan, som stridsljud från ett grannrum
     eller klockan som slår i hela tornet; ljuden står inte längre som
     egna rader. En strid som pågår hörs bara första rundan. I
     striden beskrivs varelsen helt första rundan och sedan kort
     (`kort_igen` i `data/monster.json`), som "det mångbenta". Under
     Vandringen börjar står porten sist, efter klasserna och svårigheten. I
     mörkret och i källarvalven väntar blocket tills man ser igen.
     Texterna står under `har_*` i `data/texter.json` och `spar_efter` i
     `data/fallor.json`. Rummen har inga namn för spelarna; där en text
     behöver nämna ett rum står dess korta beskrivning (`kort` i
     `data/texter.json`, till exempel "en smal gång under låga valv").
     Uppe till vänster på kartan står bara Mörker eller Källarvalven.
     Uppe till höger står statusen
     som tecken: klassen som ett litet klasskort (ett klick visar det
     stort), liv, skatt, drycker,
     tillstånd som fast, gift och mörker, kompassen, dunket, viskningen och
     `⁂` med de tecken i väggarna som vandraren har hittat, och under dem
     väskan med föremålen och de tomma platserna. Håll musen över ett tecken, eller tryck på
     det på en pekskärm, så står det vad det betyder. Styrningen ligger
     bakom **?** (eller tangenten `?` eller `H`) och visas en gång första
     vandringen, tills den tar slut. Panelen Nu visar det valda draget, striden, ormen i
     kammaren (och att S tar skatt), hur länge man sitter fast och vad som
     väntar, och viskningen skrivs i en egen ruta där. Den som har gått ut
     ser där hur mycket guld hen bar ut medan de andra är kvar.
     Ingen text på spelskärmen är dämpad; bara det som inte går att
     använda, som frånkopplade vandrare och tagna kort, är utgråat.
     I loggen säger ett tecken i marginalen vad slags rad det är (`!` fara,
     `›` det du ser, `·` besked, `❝` tornet talar till dig).
     Varje runda har en rubrik. På en bred skärm står loggen som text
     utan ruta nere till vänster på kartan, med en mörk skugga runt
     bokstäverna och en svag toning, så att kartan syns mellan raderna och
     går att klicka genom texten. Den senaste rundan står nederst, den
     förra svagt ovanför, och de äldre syns inte. Kartan ställer sig i den
     del som loggen inte täcker, och på en bred skärm (minst 760 punkter)
     håller den också spelkortets plats till höger fri, så att kortet inte
     ligger över ens egen figur. Tecknen uppe till höger har en mörk
     botten, så att facklorna inte syns igenom. Knappen Loggen under texten, eller L,
     öppnar hela loggen i en låda på kartan; den stängs med L, Esc eller ett
     klick utanför. I lådan, och på en smal skärm, står den senaste rundan
     överst och fyller hela loggens höjd, och de äldre rundorna ligger
     under strecket "tidigare ↓" och syns bara när man rullar ner. Inom en runda står raderna i den ordning de hände. Den första
     raden i en runda skrivs fram tecken för tecken, resten kommer en hel
     rad i taget. Har man rullat ner för att läsa ligger texten kvar där,
     och "↑ nya rader" tar en tillbaka upp. Ligger loggen efter
     skrivs rader från rundor som redan är förbi ikapp fort. Efter varje
     runda vänds ett **spelkort** upp till höger på kartan med det som
     ändrade dig (se Spelkorten): en bild och en rubrik för rundans
     viktigaste händelse, dess text och under den rundans rader, det
     farliga först: liv och guld du vann eller förlorade, föremål du
     hittade, gift, mörker, vem du slogs mot och om ormen vaknade. Resten
     av rundans kort står som små kort längst ner. Vart du gick står i
     loggen, och solen i översta raden, så kortet säger det bara när ett
     steg inte gick eller i mörkret. Har inget ändrat dig kommer inget
     kort. Kortet glider undan när du väljer nästa drag, efter åtta
     sekunder eller vid ett klick. Rundans rader
     står i `data/texter.json` under `sam_`,
   - **slutet**: en egen skärm med resultaten, epilogen, krönikan över hela
     vandringen och den egna vandringen runda för runda. Visa kartan går
     till kartan med resultatrutan och loggen, där epilogen också står, och
     Visa krönikan i panelen Nu går tillbaka.
     Enter eller Nytt torn i översta raden leder tillbaka till lobbyn när alla
     har tryckt, med ett nytt torn. Därefter kommer ett nytt klassval.
     Starta nu på
     värdsidan går från lobbyn till klassvalet, och under klassvalet
     börjar vandringen direkt.

Under spelet kan spelaren också se "tornet är stängt" när tornet stänger och "förbindelsen bröts,
försöker igen" om förbindelsen bryts.

Värdsidan ser spelaren aldrig. Den som öppnar `/spelledare` från en annan
dator får bara veta att sidan visas på datorn som kör tornet. Rutan "DU LEDER
TORNET" visas bara på tornets dator. Där leder introt till värdsidan i
stället för lobbyn. Den som vill spela där öppnar `/?spela` eller trycker
"Gå in själv" på värdsidan.

Spelet går i rundor. Alla väljer sitt drag samtidigt, och när alla har valt
(eller turtiden tar slut) utförs alla drag på en gång. Du kan ändra dig tills
dess. Den som inte hinner välja blir stående. Ett steg mot en vägg avvisas
direkt, så att ett feltryck inte kostar rundan; i mörker syns inte dörrarna
på kartan, och man får pröva sig fram.

| Tangent    | Drag                              |
|------------|-----------------------------------|
| Pilar      | gå åt det hållet                  |
| Mellanslag | vänta (i strid: håll garden)      |
| S          | sök igenom rummet, eller ta skatt i skattkammaren |
| D          | drick en dryck                    |
| A          | anfall (i strid)                  |
| Z          | byt zoom på kartan: (Rummet,) Nära, Mellan, Långt, Allt |
| L          | hela loggen (på en bred skärm)    |
| Enter      | redo / spela igen                 |
| V          | viska till tornet (från nivå 2)   |

### På en telefon eller surfplatta

Kartan snett ovanifrån ritar högst 30 bilder i sekunden på en pekskärm,
och går den trögt slår den av skuggorna och ritar med lägre upplösning.
`/?kvalitet=lag` i adressen gör det från början.

En iPhone, iPad eller Android-telefon i samma nätverk kan vara med. Öppna
samma adress i Safari eller Chrome. På en pekskärm visas knappar i stället
för tangenterna, och bara de som gäller just nu: utanför strid pilarna,
`·` (vänta), Sök och Drick; i strid pilarna (att gå är att fly), Vänta
(håll garden), Anfall, Sök och Drick; Viska när det går att viska; och i lobbyn och
slutet bara Redo (som Enter, och Nytt torn när vandringen är slut).
Styrningen bakom ? visar på samma sätt bara de tangenter som gäller. På kartan snett ovanifrån går det också att
trycka på pilarna vid dörrarna (och klicka på dem med musen). På en smal
skärm står allt i en kolumn: översta raden fast högst upp med tecknen på
en rad under solmätaren (bara tecken och tal; ett tryck förklarar),
kartan över resten av höjden, raden Nu och knapparna fast längst ner.
Där går man genom att trycka på pilarna vid dörrarna, så knapparna är en
rad med Vänta och det som gäller just nu; zoomar man ut till Långt eller
Allt kommer pilknapparna tillbaka. Under kartan, när man rullar ner, står
den senaste rundan, med en lång rumsbeskrivning på två rader och "mer ↓"
för resten, och knappen Loggen öppnar hela loggen. På en
surfplatta, och en telefon på bredden, står loggen på kartan som på en
dator och knapparna nere till vänster under den. En låg skärm (telefonen
på bredden) får en smalare logg med bara den senaste rundan, tre rader
rumsbeskrivning med "mer ↓", mindre knappar och panelen Nu på en rad. Introt fyller skärmen, så ett tryck var som
helst räknas. I klassvalet står korten två och två och nedräkningen fast
överst. Står telefonen på ljudlöst hörs inget ljud.

Solen sjunker en handsbredd per runda och visas som en mätare överst.
Gå ut genom en port för att lämna tornet: en i varje hörn, och i ett torn
på 11 × 11 eller mer även en mitt på varje sida. Vilken port som helst
duger. Den som kommer
ut med skatt räknas, och den med mest skatt vinner. Den som är kvar inne när
solen har gått ner är förlorad. Bär du inget guld frågar porten först:
gå samma håll en gång till, samma eller nästa runda, om du vill ut tomhänt.

Andras skatt och val av drag syns inte förrän spelet är slut.

Rummen bredvid en port har aldrig en fälla (`trygg_vid_porten` i
`data/fallor.json`), så att ingen förlorar liv på första steget in. Varelser
kan finnas där som förut. Två rum bredvid varandra har aldrig samma namn,
och rumstyperna används så jämnt som möjligt, så att det går att säga till
de andra vilket rum man står i. Ingen blir instängd: rummen byggs när någon
ser in i dem, och skulle ett nytt rum stänga den del av tornet det hör till
får det en dörr till mot ett okänt rum (`keepWayOut` i `game.go`).

### Vem ser vem

Du ser bara dig själv och rum du själv har gått igenom på kartan. Det är
spelets enda läge. De andra vandrarna syns bara när ni står i samma rum, redan i
lobbyn. Där du senast såg någon står en svag ring i dess färg, och i
listan över vandrare står "sågs runda 12" eller "inte sedd". Tecknet
försvinner när du ser vandraren igen, eller när den har gått ut eller dött.

Kartan ser ut som ett spelbräde i fackelsken. Golven är huggna
stenplattor, murarna har krön, dörrarna är av trä och portarna är
portvalv med torn. Rummen har inredning efter sin sort, till exempel
pelare, en brunn, bokhyllor eller en ässja. Rummet du står i lyses av en
fackla, grannrummen ligger i dunkel och rummen längre bort nästan i
mörker. Vandrarna och varelserna är miniatyrer på runda socklar, och en
fälla som har slagit till syns som det den är: ett galler, ett hål i
golvet, en klocka och så vidare. När du går in i ett rum som ingen har
varit i faller det ner som en bricka och läggs på plats.

Från början känner du bara rummet du börjar i, och kartan växer därifrån.
Kartan följer rummet du står i: det står i mitten, med ljusare golv och
en tunn ram i din färg. Vid tornets kanter går rutan bara lite utanför
tornet, så att portarna syns och ett rum i hörnet inte hamnar under
knapparna. Hur nära kartan går väljer du själv: på en smal skärm
Rummet (rummet du står i fyller bredden), och överallt Nära
(3 × 3 rum), Mellan (5 × 5), Långt (9 × 9) eller Allt du känner till. Z
går runt mellan dem, och knapparna nere till höger på kartan går en nivå
närmare eller längre ut. Scrollhjulet och två fingrar på en pekskärm zoomar gradvis, så att
du kan stanna mellan nivåerna, och ett klick på
nivåns namn växlar mellan din nivå och Allt. En smal skärm börjar på Rummet och en
större på Mellan, och sidan kommer ihåg ditt val på den enheten. Ligger
en port eller kammaren du känner till, eller en vandrare du ser, utanför
rutan pekar en liten pil vid kanten åt det hållet, med samma tecken som
på kartan. De andras portar är okända tills du själv går in
i dem (kartritaren ser alla rum som någon har gått in i). Vem som börjar vid vilken port lottas när vandringen börjar,
och var och en får bara veta sin egen. Ordningen i listan och färgen säger
inget om porten, och när någon kommer till tornet står det inte vid vilken
port. Står du i ett grannrum till en port när någon går ut hör du porten
gnissla.

Tornet berättar fortfarande när någon faller ner i källaren eller kommer
upp, eller när en fälla slår till, men bara i vilken del av tornet och
inte vem.
De vaga texterna har nycklar som slutar på `_dolt` i datafilerna, och
fällorna har `andra_dolt`. Ekofällan visar som förut var den drabbade är.
Att någon går ut eller dör får alla veta. Skriket när någon dör säger i
vilken del av tornet det hördes, eller från vilket håll för den som står
bredvid, och bar den döda guld hörs mynten rassla. Rummet minns sedan
vem som föll där, och om någon redan har tagit guldet (`har_fall` och
`har_guld_taget`).

När vandringen är slut syns alla igen. Den som har dött ser inte mer än förut.

Varje vandrare har sin egen karta. Den visar bara rum du själv har gått
igenom, med fällornas tecken och varelserna i dem. Rummen finns ändå bara en gång i tornet: en fälla slår bara till för
den som först går in, och den som kommer efter ser tecknet. Ljudet för ett
nytt rum hör bara den som gick in, och tornet berättar inte om andras nya
rum. Vill ni veta hur tornet ser ut får ni berätta för varandra. När vandringen
är slut syns hela tornet för alla.

Även när du inte ser de andra kan du ana var de är:

- **Spår:** går du in i ett rum där någon annan har gått ut de senaste
  3 rundorna ser du färska fotspår och åt vilket håll de leder, men inte
  vems de är. Antalet rundor står som `spar_rundor` i `data/fallor.json`.
  I mörker ser du inga spår.
- **Ljud genom väggarna:** står du i ett grannrum hör du strid, fällor och
  någon som dricker, med riktningen dit men utan namn. Det gäller oavsett
  om det finns en dörr. Längre bort hörs strid och fällor som förut, bara
  med delen av tornet. Ekofällan och klockan hörs lika i hela tornet.

Texterna står i `data/texter.json` under `spar`, `spar_ut` och `ljud_*`,
och texterna utan port under `ansluter_dolt`, `valkommen_dolt` och
`start_port` (i `ai.json` under `ansluter_dolt` och `lamnar_dolt`).

Förut gick det att välja att alla syns hela tiden, eller att kartan är
gemensam. De valen är borttagna; se [docs/synlighet.md](docs/synlighet.md).

När vandringen är slut står krönikan på en egen skärm, under resultaten och
epilogen. Först kommer det
som hände, runda för runda: vem som gick ut och med hur mycket, vem som dog,
av vad och hur mycket guld som blev liggande, vem som tog upp det, vem som fällde vilken
varelse, vem som slog vem och tog hur mycket guld, vem som nådde
skattkammaren först, när klockan fick solen att sjunka fortare och när
ormen vaknade. Sedan kommer en rad för varje vandrare
med hur många nävar den tog ur kammaren och hur mycket guld de gav, hur
många varelser den fällde, fällor den gick i och rum den sökte, och för
den som har träffat andra vandrare träffarna och det stulna guldet. Bara
det som inte är noll står med.
Texterna står i `data/texter.json` under `kronika_*`. Under krönikan står
den egna vandringen, samma som den som dör ser under vandringen.

### Klasser

När alla är redo i lobbyn kommer klassvalet: en skärm med ett spelkort
för varje klass. Var och en väljer ett kort med 1–8, ett klick eller ett
tryck. Ingen klass finns två gånger i samma torn, så kort som andra har
tagit går inte att välja. Vandringen börjar när alla har valt, eller när
tiden är slut (`klassval_sekunder` i `data/klasser.json`, förval 30). Den
som inte har valt får då ett slumpat ledigt kort. Klassen står i listan
över vandrare. AI-spelarna väljer efter sin
personlighet, men först när alla människor har valt, så att människorna
får de klasser de vill ha. Klassen gäller en vandring: efter Nytt torn väljer alla
igen.

| Klass | Kraft / Försvar | Styrka | Svaghet |
|---|---|---|---|
| Legoknekten | 3 / 2 | ett liv extra, och slår hårdast av alla | hittar hälften så mycket guld när hen söker |
| Tjuven | 2 / 2 | en fjärdedel mer guld ur kammaren, hittar oftare föremål, ingen märker när ormen rör sig vid hens näve, river hälften i strid mot vandrare | hittar aldrig drycker när hen söker |
| Örtkvinnan | 0 / 2 | börjar med en dryck, en dryck ger tre liv, tål gift | tappar mycket mer skatt när ormen vaknar |
| Kartritaren | 0 / 1 | ser genom dörrarna in i rummen bredvid, ser alla rum som någon har hittat och var någon har dött (ett x i hörnet), hittar alltid en hemlig gång om det går, har en dryck med sig | flyr hen ur en strid träffar slaget mot ryggen alltid |
| Ormtjusaren | 2 / 2 | ser i kammaren hur många vakna-kort som är kvar, ormen skadar inte hen när den vaknar | hittar hälften så mycket guld när hen söker |
| Spårfinnaren | 2 / 1 | ser spår 6 rundor bakåt och vems de är, märker hälften av fällorna i tid, ser var föremål och guld ligger och vem som tog guldet, känner före andra åt vilket håll ristningen finns | hör inget genom väggarna, utom med benluren |
| Lyktbäraren | 2 / 2 | ett liv extra, tål mörker, ser varelser i grannrummen | syns för andra redan från grannrummet |
| Pilgrimen | 0 / 3 | ledtrådar dubbelt så ofta, en dryck när tornet förändras, känner igen förbannade föremål | slår aldrig andra vandrare, kan bara värja sig |

Kraft är hur lätt klassen träffar i strid och Försvar hur svår den är att
träffa (se Strider). Namn, beskrivningar, texter och siffror står i `data/klasser.json`. Vad
styrkan och svagheten gör styrs av klassens `id` i `classes.go` och där
spelet händer. Klasserna går inte att stänga av.

### Strid mellan vandrare

Står två vandrare i samma rum utan någon varelse kan de slåss med A, som mot
varelserna. Står flera där väljer
tornet motståndaren: den som anföll dig förra rundan, annars den som kom
först in i tornet. Panelen Nu säger vem. Anfaller båda slår båda, och var
och en träffar med sin chans. Den som inte anfaller kan bli träffad även om
den går ut ur rummet; den som väntar håller garden och är svårare att
träffa, den som söker eller dricker lättare. Varje träff kostar 1 liv och river
åt sig en tredjedel av motståndarens guld (minst 1). Den som dödar tar allt.
Hur mycket någon bär syns inte, men man får en vink när man möter dem.
Portrummen är fredade, och varje runda med strid gör ormen lite oroligare.
AI-spelare anfaller bara människor, utom när alla vandrare är AI-spelare.
Striden går inte att stänga av.

### Spela mot datorn

Välj antal AI-spelare på värdsidan, eller med `-ai` i ett nedladdat torn:

```bash
go run . -ai 3
```

Datorspelarna står redo vid portarna. Omgången börjar när alla människor har
tryckt Enter; minst en människa måste vara med. I prickarnas förklaringar
står de som "(dator)".

- De ser bara det en människa skulle se: kända rum, varelser i kända rum och
  sin egen skatt. De ser inte okända rum, ormkorten eller andras skatt.
- De utforskar mot mitten, tar skatt i kammaren så länge ormen sover djupt
  nog, och går mot närmaste port när solen börjar bli låg, väskan är full
  eller livet tryter.
- Har de tappat vägen ut, till exempel efter ett fall genom golvet, letar
  de sig mot portarna i stället för mot mitten när de bär skatt eller
  solen börjar sjunka.
- Hur lite liv de går hem eller dricker vid står i `data/ai.json` för
  liven i datafilerna och räknas om efter vandrarens egna liv, så att en
  försiktig vandrare med få liv på Mardröm inte vänder redan i porten.
- I strid anfaller de, och flyr när nästa slag kan döda dem.
- Hör de mynt rassla när någon dör går de ibland dit, till rummet om de
  stod bredvid och annars till den del av tornet skriket kom från, men
  bara med minst halva livet kvar och gott om tid. Hur gärna står som
  `guldlust` i `data/ai.json`. De ger upp efter tio rundor. På väg ut
  tar de en omväg på högst två steg till ett rum de vet var det ligger.
- Det finns tre personligheter som turas om: en djärv, en försiktig och
  en plundrare.
- Ledtrådarna om tornets hemlighet går bara till människor.

Namn, personligheter och betänketid ligger i `data/ai.json`.

### Spelkorten

Det som händer dig under en runda blir ett kort ur samma lek som
ormkorten (`spelkort.go`, bilderna i `web/kort.js`). Efter rundan vänds
det viktigaste upp till höger på kartan med en bild, en rubrik och
händelsens text, och under den står rundans rader. Händer flera saker
samma runda står resten som små kort längst ner på kortet. Ordningen,
det viktigaste först, står under `spelkort` i `data/texter.json`. En
runda utan eget kort har Runda N som rubrik och en fackla som bild.

Korten visar bara det som hände dig, och säger inte mer än loggen. Fällor
och varelser har sin korta beskrivning som rubrik, aldrig namnet.

- **Fällorna** har var sin bild: gallret, golvet som ger efter, golvet som
  vrider sig, pilarna, dimman, ormgropen, skinnet, klockan, vindpusten
  och ekot. Klev du förbi en fälla, med amuletten, som Spårfinnaren eller
  som invigd, får du samma kort nedtonat.
- **Varelserna** har var sin bild. Kortet kommer när varelsen dyker upp,
  med samma text som i loggen, och i varje stridsrunda med hur rundan
  gick. Fäller du den är bilden överkorsad. I mörkret syns varelsen inte,
  och då kommer inget kort när den dyker upp.
- **Striden mot en vandrare** ger båda ett kort med två skuggor och den
  andras namn. Liv och guld står på kortets rader, inte i loggen.
- **Föremålen** har var sin bild. Kortet har föremålets namn och vad det
  gör; loggen säger bara att du hittade det. Visar ett förbannat föremål
  sin baksida vänds kortet en gång till och blir mörkt.
- **Drycker** och **guld** du hittar när du söker, och drycken du
  dricker, ger var sitt kort. Siffrorna står på kortets rader.
- **Källarvalven** ger ett kort när du faller ner och när du kommer upp,
  och ett **brev** ett kort med brevets text.
- **Varslen**, **namnets delar** och att du blir **invigd** kommer mellan
  rundorna och visas genast som ett eget kort (meddelandet `kort`). Står
  ett kort redan uppe väntar det nya tills det första har gått undan.
- **Stora kort** mitt på kartan: uppvaknandet, en port som stängs och
  solen när tre rundor är kvar och när den går ner visas för alla som är
  kvar inne (meddelandet `storkort`), och att du kommer ut bara för dig.
  En ny runda ersätter spelkortet som står uppe; ett stort kort väntar
  bara på ett annat stort kort.
- **Vandringens sista kort**, ute, död eller kvar när solen gick ner,
  står överst i Din vandring.
- **Vandringens största stunder**: i krönikan i slutet läggs högst fem
  av krönikans händelser fram som stora kort under resultaten och vänds
  upp ett i taget: vinnaren, uppvaknanden, de som dog, den som blev
  invigd, portar som stängdes, klockan, den första näven,
  en varelse som föll, namnets delar och den som blev kvar när solen
  gick ner, i den ordningen, högst två av samma slag. De står i den
  ordning de hände. Slagen och antalet står under `spelkort.stunder` i
  `data/texter.json`.
- **Klasskortet**: klassen står uppe till höger på kartan som ett litet
  klasskort. Ett klick visar det stort, med styrkan och svagheten.
- **Minikorten**: en rad i loggen som hör till ett kort får ett litet
  kort framför sig, och varje runda i Din vandring har rundans kort som
  små spelkort med bilden och rubriken. Ett klick på ett minikort, eller Enter,
  visar kortet stort med sin text; ett klick till eller Escape stänger
  det. På raden med siffror står ormkorten ihopslagna, "Ormen sover ×3".

### Skattkammaren och ormen

I kammaren mitt i tornet kan du ta skatt (S) en gång per runda. Varje gång
dras ett ormkort ur en lek:

- **Sover** eller **rör sig**: du får 15–45 guld och kortet tas bort ur leken.
  Rör sig ormen märker alla i tornet det.
- **Vaknar**: kortet blir kvar i leken. Alla i kammaren förlorar hälften av
  sin skatt och 2 liv och kastas ut genom en slumpad öppning. Väsandet hörs
  i hela tornet.

Eftersom bara sovkorten tas bort ökar risken för varje näve. Leken börjar med
13 sovkort och 2 vakna-kort på Normal. När ormen har vaknat blandas en ny lek med 2
färre sovkort än förra gången. Står du i kammaren ser du hur djupt ormen
sover och hör andetagen bli tätare.

Korten syns så här:

- **Ormkortet.** Kortet du drar är rundans spelkort till höger på
  kartan: namnet (Ormen sover, Ormen rör sig), en bild av ormen, kortets
  text och guldet. Bara du ser det.
- **Sömnmätaren.** I kammaren visar panelen Nu en mätare i tre steg för
  hur djupt ormen sover, som andas långsamt när den sover djupt och fort
  när den är nära att vakna. Ormtjusaren ser också vakna-korten som är
  kvar som små kort bredvid.
- **Uppvaknandet.** När ormen vaknar ser alla kortet Ormen vaknar stort mitt
  på kartan, med det de märkte. Den som stod i kammaren ser vad den
  förlorade.
- **Leken.** Bredvid guldhögen på kartan ligger leken med baksidan upp och
  de dragna korten med framsidan upp. Står du i kammaren krymper leken och
  de dragna korten växer, i samma tre steg som sömnmätaren. Andra ser
  högar som inte ändras. När ormen vaknar flyger de dragna korten tillbaka
  in i leken.
- **Krönikan.** Din vandring visar kortet du drog med ett litet kort
  bredvid, och i krönikan i slutet har var och en sina kort på
  raden med siffror.

Kortens namn står under `namn` i `data/ormen.json`.

### Strider

Ett nytt rum kan i stället för en fälla ha en varelse. Den syns på kartan för
alla. Står du i samma rum som en varelse är du i strid. Det enda stridsvalet
är **anfall** (A): då slår både du och varelsen, och var och en träffar med
sin chans (se [docs/bara-anfall.md](docs/bara-anfall.md)).

Chansen kommer av två värden. **Kraft** är hur lätt man träffar och
**Försvar** hur svår man är att träffa. Varje klass har sina
(`kraft` och `forsvar` i `data/klasser.json`), och varje varelse sina
(`data/monster.json`). Träffchansen är `bas + per_steg × (Kraft −
motståndarens Försvar)` inom `min` och `max` (`traff` i
`data/monster.json`). Föremål som svärdet, skölden och brynjan lägger
till sitt (`kraft` och `forsvar` i `data/foremal.json`). Dina egna värden
står bland tecknen (⚔ Kraft, ◆ Försvar, grönt när ett föremål höjer dem)
och på klasskorten; varelsernas ser du aldrig. Panelen Nu säger i ord hur
striden ser ut: "Du träffar lätt · den träffar ibland".

Du kan också göra något annat medan varelsen står där, men då kanske den
anfaller. Hur ofta beror på dess **lynne** (`lynne`, i procent), så det
lönar sig att lära sig varelserna:

- **vänta** (mellanslag) är att hålla garden: varelsen träffar sämre
  (`garden`), och du slår inte,
- **sök** och **drick** går, men varelsen träffar lättare (`oberedd`),
- **gå** är att fly: varelsen kanske slår mot din rygg innan du är ute.

Står du still utan att välja något hålls ingen garden. Besegrar du en
varelse får du dess guld. Den invigde och den som bär Myntet med hål gör
varelserna mer lynniga.

Varelserna vandrar också av sig själva. I början av varje runda kan en
varelse gå genom en öppen dörr till ett grannrum, ibland ut i delar av
tornet som ingen har hittat än. Där väntar den tills någon går in. En
varelse som står hos en vandrare stannar och strider, och ingen vandrar
på egen hand in till en vandrare. Portrummen och kammaren går varelserna
aldrig in i. Står du bredvid kan du höra den röra sig genom väggen. Hur
ofta de vandrar står i `vandra_chans` i `data/monster.json`, och en sort
kan ha en egen chans i `vandrar` (benväktaren rör sig nästan aldrig,
gallgasten ofta).

Varelserna sköts alltid av tornet. Den som dör väntar rundan inte på.

Den som dör får en egen skärm, **Din vandring**, med vandringens sista
kort överst: rundans kort runda för
runda, med rummet som i loggen: beskrivningen, det hen såg och det hen
hörde, hur hen dog och sina siffror, och hur många som är kvar i tornet och
har gått ut. Där kortet bara säger "−2 liv" säger Din vandring vad det kom
av ("Pilar ur små hål i väggen: −2 liv"), och vad som dödade en och
vilken varelse man fällde står med. Rundor där man bara stod kvar slås
ihop ("Runda 10–12"), och en runda då man satt fast i en grop eller
bakom ett galler säger "Du satt fast", och det står när tornet styrde vandraren. Var de andra är och vad de gör står inte där, så att den döda
inte kan avslöja något för de levande. Visa kartan går tillbaka till
spelet, och Din vandring i panelen Nu hit igen.

### Rum och fällor

När någon går in i ett nytt rum kan en fälla slå till. Ett tecken i rummet
visar sedan för alla att något har hänt där.

| Fälla           | Vad som händer                                      |
|-----------------|-----------------------------------------------------|
| Fallgaller      | du sitter fast i 2 rundor                           |
| Rasande golv    | 1 skada, du faller ner i källarvalven i 2 rundor och kommer upp någon annanstans; ingen ser var du är |
| Roterande rum   | du kastas ut genom en slumpad öppning               |
| Pilar           | 2 skada                                             |
| Giftig dimma    | 1 skada nu och 1 skada per runda i 2 rundor         |
| Ormgrop         | 2 skada och fast i 1 runda                          |
| Ömsat skinn     | fast i 1 runda                                      |
| Klocka          | solen sjunker 2 extra steg för alla                 |
| Mörker          | lyktan slocknar i 4 rundor; du ser bara ditt eget rum |
| Eko             | alla får veta var du är                             |

Du har 7 liv på svårigheten Normal, och klassen kan ge fler (legoknekten
har 8). Den som når noll liv dör.

Ett vanligt rum kan sökas igenom en gång (S). Man kan hitta skatt, en dryck,
ett föremål, en hemlig gång till ett grannrum, eller utlösa en fälla. En
dryck (D) ger tillbaka 2 liv och tar bort gift.

### Föremål

Man bär högst tre föremål (`max_bara` i `data/foremal.json`). De står i
väskan uppe till höger på kartan, under tecknen: en plats för varje
föremål man kan bära, med föremålets tecken och namn, och under vandringen
en svag prick för varje tom plats. Ett nytt föremål lyser en stund. Alla
föremål verkar av sig själva; man använder dem aldrig. Hittar man ett
när väskan är full blir det kvar i rummet, och den som kommer in där
härnäst tar det. Den som dör tappar sina föremål och sitt guld i rummet.
Den som går in där, eller står kvar där när rundan är slut, tar guldet, och
den döda får i slutet veta om guldet kom ut och med vem. Föremålen följer
inte med till nästa torn. Krönikan berättar vem som hittade vad.

| Föremål | Vad det gör |
|---|---|
| Rep med kroken kvar | faller du genom golvet hamnar du inte i källarvalven (en gång) |
| Fjällamuletten | nästa fälla slår inte till, men tecknet syns (en gång) |
| Extra veke | lyktan tål ett mörker (en gång) |
| Bucklig sköld | Försvar +1, tills den har tagit tre träffar från varelser |
| Rostigt svärd | Kraft +1 |
| Fjällbrynja | Försvar +2, men flyr du ur en strid träffar slaget mot ryggen alltid |
| Väskan med dubbel botten | hälften av det du förlorar när ormen vaknar eller någon river guld är gömt |
| Kompass utan nord | visar åt vilket håll närmaste port ligger |
| Sövande flöjt | drar du ett kort som väcker ormen får du dra om (en gång) |
| Rökpulver | flyr du ur en strid hinner varelsen inte slå (en gång) |
| Spegelskärva | i ett nytt rum visar den ett rum bakom en dörr som du inte har sett: fällan, varelsen och vandrarna där (tre gånger) |
| Kastkniv | dyker en varelse upp i ditt rum tar den en skada innan striden börjar (en gång) |
| Benluren | när du kommer in i ett rum hör du om en varelse står i ett grannrum, och åt vilket håll; ljud som andra bara hör från grannrummen hör du också två rum bort, med väderstreck |
| Timglas av ben | när solen går ner väntar porten en runda till på dig, som på de invigda (en gång) |

Förbannade föremål ser värdefulla ut och är värda guld när man kommer ut,
men man märker baksidan först när man bär dem:

| Föremål | Värt | Baksidan |
|---|---|---|
| Den svarta pärlan | 40 | varje näve ur kammaren lägger ett vakna-kort i leken |
| Myntet med hål | 25 | varelserna anfaller oftare |
| Svärdet med rött fäste | 15 | Kraft +2, men varje runda i strid då du inte anfaller tar det ett liv |
| Kronan av rost | 30 | den som står i grannrummet ser dig |

Ibland hittar man ett brev från någon som var i tornet förut. Det läses
direkt och följer avslöjandets nivå. Breven står under `brev` i
`foremal.json`.

Namn, texter, hur ofta de dyker upp och hur länge de håller står i
`data/foremal.json`, och hur ofta man hittar ett när man söker under
`sok` → `vikter` → `foremal` i `fallor.json`.

Om du laddar om sidan får du tillbaka din vandrare: sidan kommer ihåg namnet
tornet gav dig, i fliken och i webbläsaren. Efter 12 timmar utan spel glöms
namnet, så att varje kväll ger ett nytt.
Det gäller också under klassvalet: tornet väntar på dig tills tiden är slut,
och är du inte tillbaka då följer du med som frånvarande och kan komma in
under vandringen.
I lobbyn startar vandringen inte förrän tio sekunder efter att någon har
gått, så att den som laddar om sidan hinner tillbaka.

## Avslöjandet (spoilervarning)

Läs inte det här om du ska spela utan att veta.

Tornet är ormen. Väggarna är dess kropp, valven dess revben, trappan dess
ryggrad, skattkammaren dess mage och porten dess mun. Spelarna upptäcker det
gradvis, i fyra nivåer som följer solen:

| Nivå | När (andel sol kvar) | Text | Bild | Ljud |
|------|----------------------|------|------|------|
| 1 | över 70 % | små konstigheter: varma väggar, söt lukt, fjällformade stenar | vanliga väggar | ett svagt, långsamt dunk |
| 2 | 70–40 % | tornet andas | kritstrecken sväller och krymper | dunket blir ett hjärta |
| 3 | 40–15 % | bevis: revben, kotor, halvt upplösta kvarlevor | väggar som revben, köttfärg, hjärtslag i vågor | snabbare hjärta, kluckande |
| 4 | under 15 % | det sägs nästan rakt ut; porten stänger sig | rött, porten blir en mun med tänder | starkt och snabbt hjärta |

- Varje rum har beskrivningar för alla fyra nivåerna, och spelare som står i
  samma rum kan få olika texter.
- Var fjärde runda, och vid varje nytt steg, får varje spelare en
  ledtråd. Alla drar ur samma texter, var och en ur sin egen lek, så att
  var och en kan förstå hela bilden själv.
- På nivå 4 får varje spelare nya förklaringar till de fällor och varelser
  hen har mött (gallret som rasade ner var tänderna, de roterande rummen tarmarna ...).
- När spelet är slut säger tornet hemligheten rakt ut, och porten stängs.

### Att viska hemligheten

Från nivå 2 kan en människa viska till tornet med V (eller knappen Viska),
en gång per nivå. Viskningen skrivs i en egen ruta och syns bara för
servern och spelledaren. Rätt viskning gör vandraren **invigd**. Tornet
svarar bara "Det vet att du vet", och de andra hör att det blev tyst
någonstans. Rätt viskning innehåller ormens namn.
Varje vandring lottar ett namn i tre delar ur
`namn` i `avslojande.json`, till exempel Karumeth (Ka, ru, meth). Delarna
finns på tre ställen i tornet, och var och en hittar dem själv. Vilken
del som finns var lottas varje vandring:

- **Ristningen.** Tecken i väggen i ett rum en bit in i tornet, minst en
  tredjedel av bredden från varje port. Den hamnar i ett av rummen som
  byggs när någon ser in i dem eller går dit (`ristning_chans` per rum som
  passar), så att den alltid går att nå. Den som står i rummet efter
  rundan läser den. Finns ingen ristning när nivå 3 kommer hamnar den i
  ett känt rum långt från portarna. Från nivå 3 får den som inte har läst
  den veta åt vilket håll den finns, och Spårfinnaren redan från nivå 2.
- **Varelsen.** Den som träffar en varelse ser ett tecken under fjällen.
- **Kammaren.** Den som tar en näve hör ormen väsa delen i sömnen.

Ingen får veta vad någon annan har hittat. På nivå 4 får den som har
hittat något veta att det är tornet självt som heter så.

Texten säger inte vilken i ordningen delen är, bara var på en orm den
sitter: vid huvudet, i buken eller i svansen (`platser`, med några
olika ord för varje). Huvudet är först, buken i mitten och svansen sist,
men det får man lista ut själv. Delarna står i loggen och i tecknen,
i den ordning man hittade dem (`RU · ·`), så att tecknen inte avslöjar
ordningen; förklaringen till tecknet säger det. Viskningen ska innehålla namnet,
ihop eller i delar, och visa att man förstår: ett av svaren under `svar`
("Karumeth, tornet är ormen") eller en mall under `svar_med_namn`
("tornet är Karumeth"). Svaren räcker inte ensamma, så pusslet är nytt
varje gång, också för den som har spelat förut. Viskar man bara namnet
känner tornet igen det och frågar vad man vill säga om det
(`viskning_bara_namn`), och försöket är inte förbrukat. Har viskningen
rätt delar i fel ordning säger tornet det (`viskning_ordning`), men
försöket är förbrukat, så att ingen kan pröva sig igenom ordningarna. En
fel viskning gör ormen oroligare: ett vakna-kort till i leken, som när
vandrare slåss (`viskning_fel`). Saknar
`avslojande.json` namn räcker svaren ensamma.

Till ormens namn hör också två saker:

- **Byte vid dödande.** Den som dödar en annan vandrare hittar de delar av
  namnet som offret hade och dödaren saknade, med var på ormen de sitter.
- **Den förste invigde.** Den som blir invigd först i vandringen får en
  fjärdedel mer guld till när hen kommer ut (`invigd.forst_guld`), och
  krönikan säger att ingen hann före. Tornets nåd, som tog bort vakna-kort
  när alla var invigda, togs bort 2026-11-23: spelet är en tävling.

Den invigde:

- hör dunket och vet åt vilket håll kammaren ligger,
- känner igen hälften av fällorna och kliver förbi,
- kan en gång viska i kammaren och lugna ormen, så att hens nästa näve inte väcker den (leken är densamma för de andra),
- får en fjärdedel mer guld när hen kommer ut,
- och när solen går ner väntar porten en runda till för de invigda.

Men varelserna anfaller den invigde oftare, och en näve ur kammaren väcker
ormen lättare. Siffrorna står under `invigd` i `avslojande.json`.
Krönikan nämner den som förstod, den som var först, bytena och de invigda som kom ut,
och epilogen säger namnet.

### Varslen

Från nivå 2 är varannan ledtråd ett varsel: tornet förråder ett rum som
ingen har varit i, två eller tre rum från vandraren, till exempel "Det
gnisslar av järn två rum bort åt nordöst". Vad som finns i rummet bestäms
när varslet ges, med samma chanser som annars, och blir så när någon går
in: en fälla, en varelse eller ingenting. Varslet gäller aldrig en port
eller kammaren. Varje vandrare får sitt eget, så det lönar sig att berätta,
eller att låta bli. Texterna står under `varsel` i `avslojande.json`.

### Tornet förändras

Tornet förändras:

- Från nivå 3 blir fällorna som hör till ormens kropp vanligare i nya rum:
  det rasande golvet, det roterande rummet, ormgropen och det ömsade
  skinnet (`forandring.fallvikter`).
- På nivå 4 sluter sig portarna en i taget, var tredje runda
  (`port_stangs_var`), först den som ligger längst från vandrarna. Rundan
  innan får alla som är inne en varning. Den sista porten stängs aldrig.
  Kartan visar en stängd port med dörrarna stängda och röda järnband, och
  ingen kommer ut genom den.

Avslöjandet går inte att stänga av (sedan 2026-11-23).

## Ladda ner och kör tornet själv

Det vanliga är att spela från länken (se Starta). Att ladda ner spelet och
köra tornet på en egen dator behövs bara när länken inte räcker:

- länken kommer inte fram, till exempel på ett nät som stoppar
  anslutningen (det finns ingen reläserver, TURN, än),
- ni sitter på samma wifi och vill klara er utan internet,
- du ändrar i `data/` och vill prova direkt, utan att pusha till `main`,
- du vill använda flaggorna nedan, till exempel `-sol`, `-turtid` eller
  `-simulera`.

Det kräver en Mac och Go (`brew install go`). Med en Windows-dator spelar
man från länken. Det är bara den som kör tornet som behöver spelet. De
andra öppnar en adress i webbläsaren, på vilken dator, telefon eller
surfplatta som helst, i samma nätverk.

```bash
./starta.sh
```

Skriptet bygger spelet, öppnar introt i webbläsaren och håller datorn vaken
så länge tornet är öppet. På datorn som kör tornet leder introt till
värdsidan. Flaggor skickas vidare, till exempel `./starta.sh -ai 2`.
Är tornet redan igång öppnas bara introt.

Det går också att starta direkt med `go run .`.

### Ormtornet.app

```bash
./skapa-app.sh
```

skapar `~/Applications/Ormtornet.app` med tornets ikon. Mappen skapas om
den saknas; ange en annan mapp som argument om du vill, till exempel
`./skapa-app.sh ~/Desktop`. Ikonen ritas av `verktyg/appikon`, så
appen går att ta bort och skapa igen när som helst. Dubbelklick på appen bygger spelet från den här mappen,
startar tornet i bakgrunden och öppnar introt, som leder till värdsidan. Inget terminalfönster behövs. Loggen hamnar i `/tmp/ormtornet.log`.
Appen behöver bara skapas om när den här mappen flyttas.

Servern lyssnar på `0.0.0.0:8080` och skriver ut adressen spelarna ska
öppna, till exempel `http://Din-Mac.local:8080`. Alla Macar i samma nätverk
kan öppna den i Safari eller Chrome.

Stäng tornet med Ctrl+C i terminalen eller knappen "Stäng tornet" på
värdsidan. Spelarna får då ett meddelande om att tornet stänger.

Ett bortglömt torn stänger av sig självt, hur det än har startats:

- efter 10 minuter utan någon ansluten spelare och utan att introt eller
  värdsidan har synts på skärmen på tornets dator (en flik som ligger
  gömd bakom andra räknas inte), och
- senast när ingen har gjort något på 12 timmar, även om flikar står öppna.

En minut innan får alla som är anslutna en varning. Tiderna ändras med
`-autostopp` och `-autostopp-max`, och `0` stänger av regeln.

Inställningarna ändras enklast på värdsidan (se Värdsidan) och sparas
till nästa gång. Flaggorna nedan gäller före det sparade, men bara den gången.

### Flaggor

- `-addr 0.0.0.0:8080` adress och port
- `-data data` mapp med datafilerna
- `-sol 0` antal rundor innan solen går ner; `0` är automatiskt efter tornets storlek (40 på 9 × 9, 50 på 11 × 11, 60 på 13 × 13). Inte en inställning på värdsidan; för den som provar
- `-turtid 0` längsta tid att välja drag per runda; `0` är automatiskt efter antalet människor (30–40 sekunder). Inte en inställning på värdsidan; för den som provar
- `-min 2` minsta antal spelare för att börja, datorspelare inräknade (sätt `-min 1` för att prova ensam). Inte en inställning på värdsidan; för den som provar
- `-ai 2` antal datorstyrda vandrare (0–7)
- `-svarighet normal` svårighet: `latt`, `normal`, `svar` eller `mardrom`
- `-sprak en` språket i spelet och på värdsidan: `sv` eller `en`. Det som inte är översatt visas på svenska (se Språk nedan)
- `-simulera 500` spelar 500 vandringar med bara AI-spelare och skriver hur det gick, utan att öppna tornet (se Simulera nedan)
- `-kontrollera` kontrollerar datafilerna och avslutar, utan att öppna tornet
- `-sida mapp` skriver sidorna för GitHub Pages till mappen och avslutar. Använd hellre `verktyg/bygg-sida.sh`, som också bygger tornet (se Tornet i webbläsaren nedan)
- `-fro 0` frö för den första vandringen. Varje vandring skriver sitt frö i terminalen och under Fröet på värdsidan, där nästa vandrings frö också går att skriva in eller slumpa. Startar du med samma frö och vandrarna gör samma drag blir tornet detsamma, rum för rum, så att en konstig vandring kan spelas upp igen. `0` ger ett slumpat frö
- `-oppna` öppnar introt i webbläsaren när servern har startat; på den här datorn leder det till värdsidan
- `-vaken` håller datorn vaken så länge servern är igång
- `-autostopp 10m` stänger tornet när ingen spelare är ansluten och ingen sida har synts på tornets dator så länge (`0` = aldrig)
- `-autostopp-max 12h` stänger tornet när ingen har gjort något så länge, även om flikar står öppna (`0` = aldrig)
- `-installningar fil` var spelledarens inställningar sparas (standard `~/Library/Application Support/Ormtornet/installningar.json`, tom för att inte spara)


## Ändra innehåll

Allt innehåll ligger i mappen `data/`. På länken gäller innehållet i
`main` när sidan senast byggdes, så en ändring syns där först efter en
push. Det som står nedan om att läsa om filerna gäller ett nedladdat torn
(se Ladda ner och kör tornet själv), där filerna läses när servern startar.
Efter en ändring räcker det att fälla ut "Läs om innehållet" längst ner
på värdsidan och trycka på knappen. I lobbyn gäller det nya genast och tornet
byggs på nytt; under och efter en vandring gäller det från nästa torn. Är
något fel i filerna står det under knappen, och tornet behåller det gamla.
Filerna i `web/` byggs in i programmet, så efter en ändring där behöver
servern startas om.

Tornet kontrollerar filerna när det startar och vägrar starta om något är
fel. Felet står i terminalen med filen och fältet, till exempel ett
felstavat fält, en chans som inte ligger mellan 0 och 1, ett id som finns
två gånger eller en klass i `ai.json` som inte finns i `klasser.json`.
Fält som börjar med `_` är kommentarer och får stå var som helst.

Rum, fällor och varelser har inga namn för spelarna. `namn` syns bara på
värdsidan. Spelarna får beskrivningen och det de ser och hör, och där en
text behöver nämna något står dess korta beskrivning: `kort` för rummen
("en smal gång under låga valv") och fällorna ("ett galler av spetsiga
järnstänger"), och `kort`, `kort_bestamd` och `kort_igen` för varelserna
("något blekt och mångbent", "det bleka mångbenta", "det mångbenta").
`kort_igen` används i striden från andra rundan. Ett test kräver att alla
har dem på båda språken.

Fällor, varelser, föremål och klasser har ett fält `bild`: motivet på
deras spelkort eller klasskort, ett av dem som `web/kort.js` kan rita. Tornet säger till om
bilden inte finns. `verktyg/kortbilder.html` visar alla motiv; öppna den
med `python3 -m http.server 8095` i repots rot och
<http://localhost:8095/verktyg/kortbilder.html>.

- `texter.json` rum, beskrivningar, allmänna texter och introts väntetid. Varje rumstyp har sin inredning på kartan under `modell.inredning`, med namnen i `klient/src/inredning.ts`
- `fallor.json` regler (liv, chans för händelse, drycker, hur länge spåren syns), fällor och sökresultat
- `ormen.json` ormkorten, skatten i kammaren och vad som händer när ormen vaknar
- `monster.json` varelserna (liv, skada, guld, Kraft, Försvar och lynne, hur ofta de vandrar, texter, och modellens färg och storlek på kartan snett ovanifrån) och striderna
- `ai.json` datorspelarnas namn, personligheter och betänketid
- `svarighet.json` förvalen för svårighet (se nedan)
- `klasser.json` klasserna: namn, styrka, svaghet, liv, drycker, siffror och texter, och figurens färg och storlek på kartan snett ovanifrån
- `foremal.json` föremålen: tecken i väskan, namn, beskrivning, hur ofta, hur länge och texter
- `avslojande.json` nivågränser, ledtrådar, texter vid nivåbyte, nya
  förklaringar till fällorna, epilogen, svaren, ormens namn och texterna
  för delarna, den förste invigdes guld, varslen och hur tornet förändras. Rumsbeskrivningarna för varje
  nivå (`"1"` till `"4"`) ligger i `texter.json`.
- `spelarsidan.json` spelarsidans egna texter: introt, klassvalet, styrningen, panelerna och knapparna
- `vardsidan.json` värdsidans texter, och tornets besked och loggrader till spelledaren (`server_*`)
- `en/` översättningarna till engelska (se Språk nedan)

### Språk

Tornet kan visas på svenska eller engelska. Spelledaren väljer under Inför
nästa vandring på värdsidan, eller med `-sprak en`. Bytet gäller genast i
lobbyn, där tornet byggs om, och annars från nästa torn.

De svenska filerna är grunden. Översättningen ligger i `data/en/` med samma
filnamn och innehåller bara texter. Siffror och val står bara i de svenska
filerna, så att en ändring av balansen görs på ett ställe. Listor som
klasser, varelser, fällor och svårigheter kopplas ihop med sina id:n:

```json
{ "klasser": [ { "id": "kartritare", "namn": "The Mapmaker" } ] }
```

Det som inte är översatt visas på svenska. Tornet vägrar starta på
engelska om översättningen har ett fält som inte finns på svenska, ett id
som inte finns, ett tal, eller en platshållare som `{namn}` som den
svenska texten inte har. `-kontrollera` skriver hur många texter i varje
fil som är översatta, och ett test kräver att allt är översatt.

Sidorna hämtar sina texter från `/texter/spelarsidan.js` och
`/texter/vardsidan.js` (se `web/texter.js`). HTML:en har bara nycklar
(`data-t="nyckel"`), och skripten skriver texter med `t("nyckel")`. Byts
språket hämtar öppna sidor de nya texterna utan att laddas om.

Undantaget är sidornas titel. Den står som "Ormtornet" redan i HTML:en,
så att en delad länk visar namnet och fliken har ett namn medan sidan
laddar. Sedan byts den mot titeln på tornets språk:

| Sida | Svenska | Engelska |
|---|---|---|
| Spelarsidan | Ormtornet | The Serpent Tower |
| Skalet på Pages, innan tornet har startat | Ormtornet – tornet | The Serpent Tower – the tower |
| Värdsidan | Ormtornet – värd | The Serpent Tower – Keeper |

Titlarna är `sidans_titel` i `data/spelarsidan.json` och `skal_sidtitel`,
och `sidtitel` i `data/vardsidan.json`.

En viskning till tornet räknas på båda språken.

### Svårighet

`svarighet.json` har fyra förval: Lätt, Normal, Svår och Mardröm. Ett förval
kan ändra antal liv (`liv`), chansen för en fälla (`fallchans`) och för en
varelse (`varelsechans`) i ett nytt rum, antalet sovkort (`sovkort`) och
vakna-kort (`vaknakort`) i ormleken, skadan när ormen vaknar (`ormskada`)
och hur stor andel av skatten som går förlorad (`skatt_forlust`). `skada`
läggs till skadan från alla fällor och varelser som skadar (ingen skadar
mindre än 1), och `sol` gångrar solens rundor när solen är automatisk
(1,25 ger 50 rundor i stället för 40). `kammarskatt` gångrar guldet i
varje näve ur skattkammaren. Det som inte
står med i ett förval tas från de andra datafilerna. Normal står därför tom,
och då gäller filerna som de är. `standard` är förvalet som gäller när inget
annat är valt. När vandringen börjar får alla veta svårigheten och hur många
liv de har.

| Förval  | Liv | Fälla | Varelse | Skada | Sol    | Sovkort | Vakna-kort | Ormskada | Förlust | Näven  |
|---------|-----|-------|---------|-------|--------|---------|------------|----------|---------|--------|
| Lätt    | 8   | 35 %  | 12 %    |       | × 1,25 | 17      | 2          | 1        | 35 %    | 15–45  |
| Normal  | 7   | 45 %  | 20 %    |       |        | 13      | 2          | 2        | 50 %    | 15–45  |
| Svår    | 5   | 55 %  | 28 %    |       |        | 11      | 2          | 2        | 60 %    | 19–56  |
| Mardröm | 6   | 65 %  | 35 %    | +1    |        | 9       | 3          | 3        | 75 %    | 23–68  |

För varje fälla kan du ändra namn, vikt (hur ofta den dyker upp), skada,
antal rundor, tecknet på kartan och texterna. Vad fällan gör styrs av dess
`id` i `events.go`. I texterna kan du använda
`{namn}`, `{horn}`, `{riktning}`, `{region}` och `{rum}` där de passar.
Beskrivningarna för rummen ligger under nivåerna `"1"` till `"4"`.

## Struktur

- `main.go` startar webbservern, läser flaggorna och skriver ut adressen
- `game.go` tornet, rummen, vandrarna, rundorna och vad var och en ser
- `seed.go` fröet: nästa vandrings frö, som spelledaren kan skriva in eller slumpa
- `settings.go` spelledarens inställningar och hur de sparas
- `host.go` värdsidan och dess adresser
- `activity.go` autostopp för ett bortglömt torn
- `takeover.go` AI tar över en frånvarande vandrare eller en som står still
- `pvp.go` strid mellan vandrare
- `events.go` fällor, sökande och drycker
- `serpent.go` skattkammaren och ormkorten
- `monster.go` varelser, strider och hur varelserna vandrar
- `reveal.go` avslöjandet: nivåer, ledtrådar och epilog
- `chronicle.go` krönikan över vandringen
- `here.go` rumsblocket i loggen: rummets beskrivning och vad som finns där
- `sounds.go` spår och ljud från grannrummen
- `summary.go` rundans sammanfattning, raderna på spelkortet
- `spelkort.go` spelkorten: rundans kort i ordning, korten mellan
  rundorna och de stora korten (bilderna i `web/kort.js`)
- `stunder.go` vandringens största stunder i krönikan i slutet
- `difficulty.go` förvalen för svårighet
- `classes.go` klasserna och valet av klass
- `items.go` föremålen
- `secret.go` viskningen, den invigde, ormens namn och den förste invigde
- `namndelar.go` var namnets delar finns: ristningen, varelsen och kammaren
- `omens.go` varslen
- `change.go` tornet som förändras: fällvikterna och portarna som sluter sig
- `ai.go` datorstyrda vandrare
- `data.go` läser in datafilerna
- `validate.go` kontrollerar datafilerna när tornet startar
- `reload.go` läser om datafilerna från värdsidan
- `simulate.go` simulerar vandringar med bara AI-spelare (`-simulera`)
- `verktyg/provspelare.mjs` en robot som provspelar som människa över
  websocketen och samlar alla texter (se Provspela med robotar)
- `verktyg/textkontroll.mjs` reglerna som märker trasiga texter, och
  `verktyg/textkontroll.prov.mjs` som provar dem
- `verktyg/webbprov.mjs` webbläsarprovet: en kort vandring på en dator och
  en telefon i Chromium (se Webbläsarprovet)
- `verktyg/kortbilder.html` visar alla spelkortens och klasskortens bilder
- `lang.go` språken: översättningarna i `data/en/` läggs över de svenska texterna
- `web/texter.js` sidornas texter: `t()`, `applyTexts()` och språkbytet
- `ws.go` en liten websocket-implementation utan beroenden
- `wasm.go` tornet i webbläsaren, byggt till WebAssembly (se Tornet i
  webbläsaren)
- `sida.go` skriver sidorna för GitHub Pages (`-sida`)
- `verktyg/bygg-sida.sh` bygger allt för GitHub Pages till `sida/`
- `.github/workflows/kontroll.yml` kontrollerar formateringen, testerna och
  skripten vid varje push, och `sida.yml` publicerar på GitHub Pages vid
  varje push till `main`
- `web/` klienten: HTML, CSS och JavaScript utan ramverk
  - `index.html` och `app.js` spelarsidan: introt, klassvalet
    och spelet
  - `spelledare.html` värdsidan
  - `spelledare.css` värdsidans stil
  - `torn.html` och `torn.js` skalet: tornets flik när spelet kommer från
    GitHub Pages, och `torn-arbetare.js` workern som kör tornet i den
  - `natet.js` vägen till tornet från GitHub Pages, i samma webbläsare
    eller över WebRTC från länken; gör ingenting när sidorna kommer från
    tornets program
  - `rtc/rtc.js` WebRTC med Trystero, byggt från `klient/`
    (rör den inte för hand)
  - `map.js` räknar ut kartans utsnitt och zoom, så att kartan följer
    rummet vandraren står i. Den skapar kartan snett ovanifrån, eller visar en
    text där WebGL saknas
  - `audio.js` allt ljud, skapat med Web Audio API, även stormen i introt
  - `kort.js` spelkortens och klasskortens bilder, ritade som SVG, och
    minikorten
  - `style.css` utseendet
  - `ikon.svg`, `ikon.png`, `ikon-32.png` och `favicon.ico` tornets ikon
    i flikarna och på hemskärmen. De små ritas av `verktyg/appikon`
    (`icon_32x32.png`), eftersom Safari inte använder SVG:n
  - `iso/iso.js` kartan snett ovanifrån, byggd från `klient/` (rör den inte
    för hand). Den läser utsnittet från `map.js`
- `klient/` kartan snett ovanifrån i TypeScript med Three.js (se
  [docs/isometrisk.md](docs/isometrisk.md)) och WebRTC för GitHub Pages
  (`src/rtc.ts`), byggda med Vite. Kartan är `src/karta.ts` (tornet,
  ljuset, brickorna och ormleken), `src/modeller.ts` (vandrarna,
  varelserna, socklarna och fällorna), `src/inredning.ts` (rummens
  inredning) och `src/krita.ts` (texturerna: sten, murverk, trä, krita och
  korten). De byggda skripten `web/iso/iso.js` och
  `web/rtc/rtc.js` checkas in, så att spelet byggs utan Node. Efter en
  ändring i `klient/src/`:

  ```bash
  cd klient && npm install && npm run bygg
  ```

  `starta.sh` bygger om skriptet av sig själv när källkoden är nyare och
  Node finns.

## Testa lokalt

Starta servern (se Ladda ner och kör tornet själv) och öppna
`http://localhost:8080` i två flikar. För ett snabbt test passar `go run . -sol 8 -turtid 10s`.

Spellogiken har tester:

```bash
go test ./...
```

### Tornet i webbläsaren

Spelet ska kunna spelas från en länk på GitHub Pages
([docs/github-pages.md](docs/github-pages.md)). Då körs tornet i
spelledarens webbläsare, i en flik som heter skalet (`web/torn.html`).
Sidorna i samma webbläsare når det över en `BroadcastChannel`
(`web/natet.js`). Spelarna kommer med en länk med en rumskod och en nyckel,
`…/ormtornet/#torn=<kod>&nyckel=<nyckel>`, som värdsidan visar under Bjud
in. Deras webbläsare hittar tornets flik genom
publika Nostr-reläer och ansluter sedan direkt med WebRTC
(`web/rtc/rtc.js`, Trystero). Nyckeln krypterar utbytet, och den som har
länken kommer in. Länken är densamma så länge tornets flik är öppen, också
om den laddas om.

```bash
verktyg/bygg-sida.sh
```

```bash
python3 -m http.server 8000 -d sida
```

Tornet självt går i en Web Worker (`web/torn-arbetare.js`), så att det
inte saktas ner när fliken ligger i bakgrunden; skalet för bara vidare
meddelandena. Bygget sätter ett versionsnummer (`?v=`) på sidornas skript,
stilar och bilder, så att en ny version slår igenom på en gång trots att
Pages låter webbläsaren spara filerna en stund.

Öppna `http://localhost:8000/`. Introt märker att inget torn är öppet och
öppnar skalet, som visar introt och sedan värdsidan. "Gå in själv" på
värdsidan byter till spelarsidan i samma flik. Länken under Bjud in går också
att öppna i en annan webbläsare; den ansluter över WebRTC. Skalets flik
måste vara öppen så länge tornet ska vara det.

Inställningarna sparas i webbläsaren. Läs om innehållet läser bara det
inbyggda innehållet, eftersom det inte finns någon `data/` att läsa.

Vid varje push till `main` bygger GitHub Action `Sida`
(`.github/workflows/sida.yml`) samma sak och publicerar den på GitHub
Pages, på `https://<namn>.github.io/ormtornet/`. Det kräver att repot är
publikt och att Pages är påslaget i repots inställningar (Settings, Pages,
Source: GitHub Actions).

### Simulera

För att se hur siffrorna i `data/` verkar kan tornet spela många
vandringar med bara AI-spelare, utan webbläsare och utan väntetider:

```bash
go run . -simulera 500 -ai 4 -svarighet svar
```

Det tar under en sekund. Utskriften säger hur många som kom ut med och
utan skatt, dog eller blev kvar i mörkret, hur mycket skatt de bar, vad de
dog av, hur ofta ormen vaknade och hur det gick för varje klass.
Inställningarna är desamma som när tornet öppnas (`-sol` och de sparade). Vandringarna får frön i följd
från `-fro`, så samma kommando ger samma resultat, och en enskild vandring
kan köras om med `-simulera 1 -fro` och dess frö.

I simuleringen får AI-spelarna klasserna i tur och ordning i stället för
efter sin personlighet, så att varje klass spelas lika ofta och lika ofta
av varje personlighet. Tabellen Personligheterna visar hur det gick för
den djärva, den försiktiga och plundraren, med andelen som vann (hade
mest skatt av dem som kom ut) och guld per vandrare. Tabellen Klasserna
har längst till höger en kolumn per personlighet med andelen som kom ut
med skatt, och två tabeller till visar guld per vandrare och andelen som
vann för varje klass och personlighet. Med `-ai 6` spelar varje
personlighet lika ofta. I riktiga vandringar väljer AI-spelarna
fortfarande efter sin personlighet.

### Provspela med robotar

Simuleringen spelar inne i Go-koden. För att prova det en människa möter,
med anslutningen, lobbyn, klassvalet, slutet och alla texter, finns en
robot som ansluter som en människa över websocketen:

```bash
go run . -addr 127.0.0.1:8096 -installningar "" -min 4 -ai 2 -turtid 15s
```

och i en annan terminal, en robot per spelare:

```bash
node verktyg/provspelare.mjs 8096 Alva 10 provlogg.txt & node verktyg/provspelare.mjs 8096 Bertil 10 provlogg.txt & wait
```

Varje robot spelar tio vandringar och skriver allt den får i loggen till
`provlogg.txt`, med namnet först på varje rad. Rader som ser trasiga ut
märks med `<<< MISSTÄNKT` och regeln, till exempel en platshållare som
står kvar, "1 rundor", "1 lives", samma ord två gånger eller svenska ord
i en engelsk vandring. Texterna i rundans kort, spelkorten, krönikan och
resten av det tornet skickar kontrolleras också, och en trasig text där
står en gång på en rad som börjar med `MISSTÄNKT`, så att
`grep MISSTÄNKT provlogg.txt` visar allt. Reglerna står i
`verktyg/textkontroll.mjs`; lägg till ett fel i
`verktyg/textkontroll.prov.mjs` när en regel behöver bli skarpare, och kör
`node verktyg/textkontroll.prov.mjs` (det görs också i CI). Resultatet av varje vandring står på en
rad som börjar med `RESULTAT`. `-min` ska räkna med både robotarna och
AI-spelarna, så att alla hinner ansluta innan vandringen börjar. Kräver Node 22 eller senare. Provspelet 2026-11-21
beskrivs i [docs/provspel-2026-11-21.md](docs/provspel-2026-11-21.md).

### Webbläsarprovet

Robotarna ser bara det tornet skickar. Det som går fel på sidan, som ett
minikort som inte går att klicka på eller tangenterna vid dörrarna som
ligger över varandra på telefonen, hittar webbläsarprovet:

```bash
node verktyg/webbprov.mjs webbprov sv
```

Det bygger och startar tornet (solen 8 rundor, en AI-spelare) och spelar
en vandring i två webbläsare samtidigt, en dator (1440 × 900) och en
telefon (iPhone 13, med tryck). Datorn spelar med piltangenterna och
telefonen med tangenterna vid dörrarna. Under vägen letar det efter

- fel i skripten och i konsolen, och filer som inte går att hämta,
- trasiga texter på sidan (samma regler som robotarna),
- en sida som är bredare än skärmen,
- tangenter vid dörrarna som ligger över varandra,
- ett minikort i loggen som inte öppnar kortvisningen, och att mellanslag
  spelar efter att kortvisningen har stängts,
- en omladdning mitt i vandringen som inte berättar rummet igen,
- och vandrare som saknas i resultatet.

Det tar en bild av varje fas i båda webbläsarna (introt, lobbyn,
klassvalet, vandringen, kortvisningen, omladdningen och slutet) och
skriver allt det hittade till `rapport.txt` i mappen (`webbprov/` om
ingen anges; den checkas inte in). Går något inte att trycka på på
telefonen, för att till exempel spelkortet ligger över, står det på en
rad som börjar med `OBS` och säger vad som låg över; det räknas inte som
ett fynd, och provet trycker på tangenten i stället. Med `en` spelas vandringen på
engelska. Det slutar med 1 om något hittades och tar ett par minuter.
Båda webbläsarna står på tornets dator och spelar alltså som spelledaren
med Gå in själv, med märket i översta raden.

Provet kräver Go, Node 22 och Playwright med Chromium. Har du inte
Playwright:

```bash
cd klient && npm install --no-save playwright && npx playwright install chromium
```
