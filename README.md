# Ormtornet

En dungeon crawler för 2–8 spelare i webbläsaren, på ett lokalt nätverk.
Mitt i ett torn på 9x9 rum (11x11 för fem eller sex, 13x13 för sju eller
åtta) sägs en jättelik orm sova kring en skatt.
Ta så mycket du vågar och kom ut innan solen går ner.

Vad som är nytt står i [docs/andringar.md](docs/andringar.md).

## Starta

Kräver Go (`brew install go`). Det är bara den som kör tornet som behöver
spelet. De andra öppnar en adress i webbläsaren.

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

skapar `~/Projekt/Ormtornet.app` med tornets ikon (ange en annan mapp som
argument om du vill). Dubbelklick på appen bygger spelet från den här mappen,
startar tornet i bakgrunden och öppnar introt, som leder till spelledarens
sida. Inget terminalfönster behövs. Loggen hamnar i `/tmp/ormtornet.log`.
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

Inställningarna ändras enklast på värdsidan (se nedan) och sparas
till nästa gång. Flaggorna nedan gäller före det sparade, men bara den gången.

Flaggor:

- `-addr 0.0.0.0:8080` adress och port
- `-data data` mapp med datafilerna
- `-sol 40` antal rundor innan solen går ner i ett torn på 9 × 9; större torn får fler (50 på 11 × 11, 60 på 13 × 13)
- `-turtid 30s` längsta tid att välja drag per runda
- `-min 2` minsta antal spelare för att börja, datorspelare inräknade (sätt `-min 1` för att prova ensam)
- `-ai 2` antal datorstyrda vandrare (0–7, eller 8 i spelledarläget)
- `-svarighet normal` svårighet: `latt`, `normal`, `svar` eller `mardrom`
- `-klasser=false` vandrarna väljer ingen klass
- `-simulera 500` spelar 500 vandringar med bara AI-spelare och skriver hur det gick, utan att öppna tornet (se Simulera nedan)
- `-kontrollera` kontrollerar datafilerna och avslutar, utan att öppna tornet
- `-fro 0` frö för den första vandringen. Varje vandring skriver sitt frö i terminalen och på värdsidan. Startar du med samma frö och vandrarna gör samma drag blir tornet detsamma, rum för rum, så att en konstig vandring kan spelas upp igen. `0` ger ett slumpat frö
- `-storlek 0` tornets storlek i rum: `9`, `11` eller `13`, eller `0` för automatisk (efter antalet vandrare). Ett torn blir aldrig för litet för alla
- `-synlighet egen` vilka vandrare och rum man ser: `egen` (andra vandrare bara i samma rum, och bara egna rum på kartan), `moten` (andra vandrare bara i samma rum) eller `alla`
- `-avslojande=false` stänger av avslöjandet, för spelare som redan känner till tornets hemlighet
- `-oppna` öppnar introt i webbläsaren när servern har startat; på den här datorn leder det till värdsidan
- `-vaken` håller datorn vaken så länge servern är igång
- `-autostopp 10m` stänger tornet när ingen spelare är ansluten och ingen sida har synts på tornets dator så länge (`0` = aldrig)
- `-autostopp-max 12h` stänger tornet när ingen har gjort något så länge, även om flikar står öppna (`0` = aldrig)
- `-installningar fil` var spelledarens inställningar sparas (standard `~/Library/Application Support/Ormtornet/installningar.json`, tom för att inte spara)

### Värdsidan

Öppna `http://localhost:8080/spelledare` på datorn som kör servern. Dit
kommer man också efter introt på den datorn. Där står
adressen spelarna ska öppna, med knappar som kopierar en färdig inbjudan att
skicka i en gruppchatt. Sidan visar också vilka som har anslutit och är redo,
och vad man kan göra om någon inte kommer in.

Överst står en rad som följer med när sidan rullar: lobbyn, klassvalet,
rundan och solen, eller att vandringen är slut, hur många människor och
AI-spelare som är här och vilka som är borta. Där ligger också knappen som
behövs just nu: Starta nu före vandringen och Pausa eller Fortsätt under
vandringen i spelledarläget.

Avsnitten byter plats efter fasen, så att det som behövs just nu står
först. I lobbyn och klassvalet kommer Bjud in och Vid porten först. Under
vandringen kommer översikten först, och Bjud in krymper till adressen och
en knapp som kopierar den. När vandringen är slut kommer Inför nästa
vandring först. Längst ner ligger Kommer någon inte in? och Läs om
innehållet, hopfällda, och en liten knapp som stänger tornet. Har ingen
människa kommit in en och en halv minut efter att sidan öppnades fälls
Kommer någon inte in? ut av sig själv.

Under Inför nästa vandring väljer spelledaren det här. Svårighet,
AI-spelare, Solnedgång, Turtid och Spelledarläge syns direkt. Resten ligger
under Fler inställningar, som visar de nuvarande valen i korthet tills man
fäller ut den.

| Inställning | Val | Förval |
|---|---|---|
| AI-spelare | 0–7, eller 8 i spelledarläget | 0 |
| Minst antal vandrare | 1–8, AI-spelare inräknade | 2 |
| Svårighet | Lätt, Normal, Svår, Mardröm | Normal |
| Solnedgång | rundor i ett torn på 9 × 9; större torn får fler | 40 |
| Turtid | sekunder per runda | 30 |
| Avslöjandet | på eller av | på |
| Frånvarande vandrare | AI tar över efter 1–3 rundor, eller står still | 2 rundor |
| Strid mellan vandrare | på, efter halva tiden, av | på |
| Vandrarna syns | bara när de möts med egen karta, bara när de möts, alltid | egen karta |
| Tornets storlek | automatisk, 9, 11 eller 13 rum | automatisk |
| Klasser | på eller av | på |
| Spelledarläge | på eller av | av |

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
Kommer spelaren tillbaka och skriver samma namn får den genast tillbaka sina
steg. I resultatet märks den som AI har spelat åt med "~". Den som inte kommer
tillbaka innan vandringen är slut är inte med i nästa.

Är ingen människa ansluten alls, till exempel när nätet försvinner för
alla, väntar tornet: rundan står pausad tills någon kommer tillbaka, så att
solen inte hinner gå ner medan AI:n spelar ensam. I spelledarläget syns
det i översikten, och spelledaren kan släppa pausen med Fortsätt. En
vandring med bara AI-spelare går som vanligt.

Med "Spelledarläge: på" leder spelledaren tornet utan att spela. Ingen kan
då ansluta som vandrare från spelledarens dator, och spelarna får veta att
tornet har en spelledare. Spelledaren startar vandringen med Starta nu, även
när alla vandrare är AI-spelare. Sidan visar då en översikt: hela tornet så
långt det har vuxit fram, med fällor som väntar och alla varelser, varje
vandrares plats, klass, liv, skatt, drycker, föremål, tillstånd och valda
drag, vad var och en har viskat till tornet och vem som är invigd, hur djupt
ormen sover och risken att den vaknar vid nästa näve, avslöjandets nivå och
allt vandrarna har fått läsa. Vandrarnas logg är hopfälld tills man klickar
på den. Där står den senaste rundan överst, de två senaste rundorna är
öppna och de äldre fällda ihop, och loggen kan visa en vandrare i taget.

På en bred skärm står kartan till vänster och ingreppen till höger, med
tabellen över vandrarna under dem. I spelledarläget kan spelledaren också
göra ingrepp under en vandring, och varje ingrepp skrivs i loggen:

- **Tid:** pausa och fortsätta, ge rundan 30 sekunder till, avgöra rundan
  direkt, och flytta solen fram eller tillbaka (den går aldrig ner av ett
  ingrepp, det gör bara rundorna).
- **Tornets röst:** ett meddelande till en vandrare eller alla, i en egen
  färg och med en viskning.
- **I ett rum:** lägga en fälla eller släppa in en varelse i ett rum som
  redan finns.
- **Ormen:** väcka den, eller lugna den genom att ta ett vakna-kort ur leken.
- **Avslöjandet:** flytta fram det en nivå direkt, med ljud, texter och
  ledtrådar som när solen gör det, eller hålla kvar nivån så att den inte
  följer solen. Släpps nivån hinner den ikapp solen efter rundan.
- **Vandrarna:** ge eller ta liv, guld och drycker, och ge ett föremål.
  Vandraren märker att något händer, men inte varifrån.
- **Den invigde:** inviga en vandrare, som om den hade viskat hemligheten,
  eller ta tillbaka det (bara med avslöjandet på).

Utan spelledarläge visar sidan inga skatter, positioner eller drag, så spelledaren kan spela med i
en annan flik ("Gå in själv", som leder till introt och sedan namnrutan). Andra datorer i nätverket kommer inte åt sidan.

## Spela

1. Startsidan visar tornet i en storm vid solnedgången medan rykten om det
   tonar fram. Moln driver förbi, det regnar och blixtrar, och efter varje
   blixt mullrar åskan. Klicka var som helst, eller tryck Enter, för att
   höra stormen: en låg, hotfull ton, vind, regn och åska. Webbläsaren
   tillåter inget ljud innan man har klickat. Efter 5 sekunder tonar
   "Gå in i tornet" fram efter en blixt, och porten börjar glimma. Klicka
   på den eller tryck Enter. Den som redan har sett introt i fliken, och
   spelledaren på tornets dator, får knappen direkt. Ryktena står i
   `data/texter.json` under `rykten`, och under `rykten_invigda` för den som
   spelar med `-avslojande=false`. Väntetiden står under `intro` → `vanta`
   och uppmaningen att klicka under `intro_ljud`.
2. Skriv ditt namn och tryck Enter.
3. Tryck Enter igen när du är redo. Vandringen börjar när alla är redo.

### Vad en spelare ser

En spelare som ansluter från en annan dator eller telefon ser bara
spelarsidan (`/`), med tre vyer i tur och ordning:

1. **Introt**: tornet i stormen, ryktena och, efter en stund, "Gå in i
   tornet". Hela övergången in genom porten visas bara första gången i
   fliken.
2. **Namnrutan**: "Vad heter du, vandrare?" Är tornet fullt eller namnet
   upptaget står felet här.
3. **Spelet**, där spelaren stannar resten av kvällen. Det går igenom fyra faser:
   - **lobbyn**: vandrarna och "Tryck Enter när du är redo" (Redo på en
     pekskärm),
   - **klassvalet**: en egen skärm med ett spelkort för varje klass, med
     styrka och svaghet. Välj med 1–8 eller ett klick. Kort som andra har
     tagit är gråa. Den som inte har valt när tiden är slut får ett av de
     lediga korten,
   - **vandringen**: kartan, solmätaren, vandrarna som prickar i sin färg
     i översta raden, och högerkolumnen med panelen Nu och loggen, som
     fyller resten av höjden. En prick är fylld när vandraren har valt drag
     (i lobbyn när den är redo), ✕ är död och ◌ ute; håll musen över den,
     eller tryck, för namn, klass och läge. Panelen Vandrare syns i lobbyn
     och i slutet. I loggen står vandrarnas namn i deras färg. Uppe till
     vänster på kartan står rummets namn, och under det
     teckenförklaringen: de tecken som står i rummen just nu, som
     vandrarna, varelserna och fällorna med sina namn, kammaren, portarna,
     krita och okända rum. Den fälls ihop med ett klick på rubriken (på en
     liten skärm är den hopfälld från början), och texterna står under
     `tecken_*` i `data/texter.json`. Uppe till höger står statusen
     som tecken: klassens tecken, liv, skatt, drycker,
     tillstånd som fast, gift och mörker, kompassen, dunket och viskningen,
     och under dem föremålen. Håll musen över ett tecken, eller tryck på
     det på en pekskärm, så står det vad det betyder. Styrningen ligger
     bakom **?** (eller tangenten `?` eller `H`) och visas en gång första
     vandringen. Panelen Nu visar det valda draget, striden, ormen i
     kammaren och vad som väntar, och viskningen skrivs i en egen ruta där.
     Ingen text på spelskärmen är dämpad; bara det som inte går att
     använda, som frånkopplade vandrare och tagna kort, är utgråat.
     I loggen säger ett tecken i marginalen vad slags rad det är (`!` fara,
     `›` det du ser, `▸` rummets namn, `·` besked, `❝` tornets röst,
     `✎` krita, `≡` rundans sammanfattning). Varje runda har en rubrik, och den senaste rundan står
     överst. Inom en runda står raderna i den ordning de hände. Den första
     raden i en runda skrivs fram tecken för tecken, resten kommer en hel
     rad i taget. Har man rullat ner för att läsa ligger texten kvar där,
     och "↑ nya rader" tar en tillbaka upp. Rundorna före den förra fälls
     ihop till en rad med rummen man kom till, och ett klick fäller upp dem
     igen. Ligger loggen efter
     skrivs rader från rundor som redan är förbi ikapp fort. Efter varje
     runda kommer ett kort nere på kartan med vad som hände dig, det
     farliga först: vart du gick, liv och guld du vann eller förlorade,
     föremål du hittade, gift, mörker, en varelse i rummet och om ormen
     vaknade. Kortet försvinner när du väljer nästa drag, efter sju
     sekunder eller vid ett klick, och en kort rad med `≡` står kvar i
     loggen, till exempel "Runda 12: österut, −1 liv, +10 guld". Texterna
     står i `data/texter.json` under `sam_`,
   - **slutet**: resultatrutan över kartan, och i loggen vinnaren, krönikan
     och epilogen. Enter (Nytt torn) leder tillbaka till lobbyn när alla
     har tryckt, med ett nytt torn. Därefter kommer ett nytt klassval.
     Utan klasser börjar nästa vandring direkt efter lobbyn. Starta nu på
     värdsidan går från lobbyn till klassvalet, och under klassvalet
     börjar vandringen direkt.

Under spelet kan spelaren också se att tornet har en spelledare, tornets röst
i loggen, "tornet är stängt" när tornet stänger och "förbindelsen bröts,
försöker igen" om förbindelsen bryts.

Värdsidan ser spelaren aldrig. Den som öppnar `/spelledare` från en annan
dator får bara veta att sidan visas på datorn som kör tornet. Rutan "DU LEDER
TORNET" visas bara på tornets dator. Där leder introt till värdsidan i
stället för namnrutan. Den som vill spela där öppnar `/?spela` eller trycker
"Gå in själv" på värdsidan.

Spelet går i rundor. Alla väljer sitt drag samtidigt, och när alla har valt
(eller turtiden tar slut) utförs alla drag på en gång. Du kan ändra dig tills
dess. Den som inte hinner välja blir stående.

| Tangent    | Drag                              |
|------------|-----------------------------------|
| Pilar      | gå åt det hållet                  |
| Mellanslag | vänta                             |
| T          | ta skatt (bara i skattkammaren)   |
| S          | sök igenom rummet                 |
| D          | drick en dryck                    |
| A          | anfall (i strid)                  |
| P          | parera (i strid)                  |
| U          | undanmanöver (i strid)            |
| Tab        | byt mål: motståndare, eller varelse när du är död |
| Enter      | redo / spela igen                 |
| 1–3        | använd ett föremål i väskan       |
| V          | viska till tornet (från nivå 2)   |

### På en telefon eller surfplatta

En iPhone, iPad eller Android-telefon i samma nätverk kan vara med. Öppna
samma adress i Safari eller Chrome. På en pekskärm visas knappar i stället
för tangenterna: pilarna och `·` (vänta), och Ta, Sök, Drick, Anfall,
Parera, Undan, Byt (som Tab) och Redo (som Enter, och Nytt torn när
vandringen är slut). På en smal skärm ligger knapparna fast längst ner och
sidan rullar. Står telefonen på ljudlöst hörs inget ljud.

Solen sjunker en handsbredd per runda och visas som en mätare överst.
Gå ut genom en port för att lämna tornet: en i varje hörn, och i ett torn
på 11 × 11 eller mer även en mitt på varje sida. Vilken port som helst
duger. Den som kommer
ut med skatt räknas, och den med mest skatt vinner. Den som är kvar inne när
solen har gått ner är förlorad.

Andras skatt och val av drag syns inte förrän spelet är slut.

### Vem ser vem

Som standard ser du bara dig själv och rum du själv har gått igenom på
kartan. De andra vandrarna syns bara när ni står i samma rum, redan i
lobbyn. Där du senast såg någon står en svag ring i dess färg, och i
listan över vandrare står "sågs runda 12" eller "inte sedd". Tecknet
försvinner när du ser vandraren igen, eller när den har gått ut eller dött.

Från början känner du bara rummet du börjar i, och kartan växer därifrån.
Kartan följer rummet du står i: det står i mitten, med två rum åt varje
håll (ett på en smal skärm), och vid tornets kanter stannar rutan inom
tornet. Ett tryck på kartan visar allt du känner till, och ett tryck till
går tillbaka. De andras portar är okända tills någon går in i
dem. Med egen karta ser du bara rum du själv har gått igenom; med de andra
valen syns alla rum som någon har gått in i. Vem som börjar vid vilken port lottas när vandringen börjar,
och var och en får bara veta sin egen. Ordningen i listan och färgen säger
inget om porten, och när någon kommer till tornet står det inte vid vilken
port. Står du i ett grannrum till en port när någon går ut hör du porten
gnissla. Med "alltid, för alla" börjar var och en vid porten som hör till
dess plats i listan, som förut.

Tornet berättar fortfarande när någon faller ner i källaren eller kommer
upp, eller när en fälla slår till, men bara i vilken del av tornet och
inte vem.
De vaga texterna har nycklar som slutar på `_dolt` i datafilerna, och
fällorna har `andra_dolt`. Ekofällan visar som förut var den drabbade är.
Att någon går ut eller dör får alla veta.

Den döda spelaren som styr varelserna ser alla vandrare. Spelledaren ser
allt i översikten, och när vandringen är slut syns alla igen.

Varje vandrare har sin egen karta. Den visar bara rum du själv har gått
igenom, med fällornas tecken och varelserna i dem. Rummen finns ändå bara en gång i tornet: en fälla slår bara till för
den som först går in, och den som kommer efter ser tecknet. Ljudet för ett
nytt rum hör bara den som gick in, och tornet berättar inte om andras nya
rum. Vill ni veta hur tornet ser ut får ni berätta för varandra. Den döda
spelaren som styr varelserna ser hela tornet så långt det har vuxit fram,
och när vandringen är slut syns hela tornet för alla.

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
Med "alltid, för alla" finns inga spår och inga ljud från grannrummen.

Med "Vandrarna syns: bara när de möts" på värdsidan, eller
`-synlighet moten`, är kartan gemensam igen: rum som någon upptäcker syns
för alla, och tornet berättar i vilken del någon har hittat ett nytt rum.

Med "Vandrarna syns: alltid, för alla" på värdsidan, eller `-synlighet alla`,
syns alla på kartan som förut.

När vandringen är slut berättar tornet krönikan i loggen. Först kommer det
som hände, runda för runda: vem som gick ut och med hur mycket, vem som dog
och av vad, vem som fällde vilken varelse och när ormen vaknade. Sedan kommer
en rad för varje vandrare med hur många nävar den tog ur kammaren och hur
mycket guld de gav, och hur många fällor, fällda varelser och sökta rum den
fick. Texterna står i `data/texter.json` under `kronika_*`.

### Klasser

När alla är redo i lobbyn kommer klassvalet: en skärm med ett spelkort
för varje klass. Var och en väljer ett kort med 1–8, ett klick eller ett
tryck. Ingen klass finns två gånger i samma torn, så kort som andra har
tagit går inte att välja. Vandringen börjar när alla har valt, eller när
tiden är slut (`klassval_sekunder` i `data/klasser.json`, förval 30). Den
som inte har valt får då ett slumpat ledigt kort. Klassen står i listan
över vandrare, och i spelledarens översikt. AI-spelarna väljer efter sin
personlighet, men först när alla människor har valt, så att människorna
får de klasser de vill ha. Klassen gäller en vandring: efter Nytt torn väljer alla
igen.

| Klass | Styrka | Svaghet |
|---|---|---|
| Legoknekten | ett liv extra, tar ingen skada när båda anfaller | hittar hälften så mycket guld när hen söker |
| Tjuven | en fjärdedel mer guld ur kammaren, hittar oftare föremål, ingen märker när ormen rör sig vid hens näve, river hälften i strid mot vandrare | hittar aldrig drycker när hen söker |
| Örtkvinnan | börjar med två drycker, en dryck ger tre liv, tål gift | ett liv mindre, tappar mer skatt när ormen vaknar |
| Kartritaren | ser genom dörrarna in i rummen bredvid, ser alla rum som någon har hittat och ger sin karta till den hen möter, hittar alltid en hemlig gång om det går, tre streck extra med kritan | undanmanöver lyckas aldrig |
| Ormtjusaren | ser i kammaren hur många vakna-kort som är kvar, ormen gör en skada mindre | varelserna anfaller oftare |
| Spårfinnaren | ser spår 6 rundor bakåt och vems de är, märker hälften av fällorna i tid, ser var föremål ligger | hör inget genom väggarna |
| Lyktbäraren | tål mörker, ser varelser i grannrummen | syns för andra redan från grannrummet |
| Pilgrimen | ledtrådar dubbelt så ofta, en dryck när tornet förändras, känner igen förbannade föremål | slår aldrig andra vandrare, kan bara värja sig |

Namn, beskrivningar, texter och siffror står i `data/klasser.json`. Vad
styrkan och svagheten gör styrs av klassens `id` i `classes.go` och där
spelet händer. Spelledaren stänger av klasserna med "Klasser: av" på
värdsidan, eller med `-klasser=false`.

### Strid mellan vandrare

Står två vandrare i samma rum utan någon varelse kan de slåss med samma
sten-sax-påse som mot varelserna (A, P, U, och Tab för att välja motståndare
när flera står där). Den som inte försvarar sig blir träffad om den andra
anfaller, även om den går ut ur rummet. Varje träff kostar 1 liv och river
åt sig en tredjedel av motståndarens guld (minst 1). Den som dödar tar allt.
Hur mycket någon bär syns inte, men man får en vink när man möter dem.
Portrummen är fredade, och varje runda med strid gör ormen lite oroligare.
AI-spelare anfaller bara människor, utom när alla vandrare är AI-spelare.
Spelledaren väljer under "Strid mellan vandrare": på, efter halva tiden
eller av.

### Spela mot datorn

Välj antal AI-spelare på värdsidan, eller med `-ai`:

```bash
go run . -ai 3
```

Datorspelarna står redo vid portarna. Omgången börjar när alla människor har
tryckt Enter; minst en människa måste vara med. I listan över vandrare står
de som "(dator)".

- De ser bara det en människa skulle se: kända rum, varelser i kända rum och
  sin egen skatt. De ser inte okända rum, ormkorten eller andras skatt.
- De utforskar mot mitten, tar skatt i kammaren så länge ormen sover djupt
  nog, och går mot närmaste port när solen börjar bli låg, väskan är full
  eller livet tryter.
- I strid gissar de varelsens vanor utifrån dess beskrivning och lär sig av
  det de ser. De flyr när nästa slag kan döda dem.
- Det finns två personligheter som turas om: en djärv och en försiktig.
- Dör en datorspelare först tar den över varelserna och flyttar dem mot
  närmaste levande vandrare, tills en människa dör. Då tar människan över.
- Ledtrådarna om tornets hemlighet går bara till människor.

Namn, personligheter och betänketid ligger i `data/ai.json`.

### Skattkammaren och ormen

I kammaren mitt i tornet kan du ta skatt (T) en gång per runda. Varje gång
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

### Strider

Ett nytt rum kan i stället för en fälla ha en varelse. Den syns på kartan för
alla. Står du i samma rum som en varelse är du i strid. Båda väljer i hemlighet:

- **parera** slår **anfall**
- **anfall** slår **undanmanöver**
- **undanmanöver** slår **parera**

Vinnaren träffar. Väljer båda anfall träffar båda; annars händer inget vid
lika val. Varje sorts varelse har egna vanor, så det lönar sig att lära sig
dem. Besegrar du en varelse får du dess guld.

Gör du något annat än att slåss (söker, dricker, väntar) får varelsen ett
fritt drag mot dig. Du kan fly genom att gå ut ur rummet, men varelsen får
ett fritt drag mot din rygg först.

Varelserna vandrar också av sig själva. I början av varje runda kan en
varelse gå genom en öppen dörr till ett grannrum, ibland ut i delar av
tornet som ingen har hittat än. Där väntar den tills någon går in. En
varelse som står hos en vandrare stannar och strider, och ingen vandrar
på egen hand in till en vandrare. Portrummen och kammaren går varelserna
aldrig in i. Står du bredvid kan du höra den röra sig genom väggen. Hur
ofta de vandrar står i `vandra_chans` i `data/monster.json`, och en sort
kan ha en egen chans i `vandrar` (benväktaren rör sig nästan aldrig,
gallgasten ofta).

Den första spelaren som dör tar över tornets varelser. Var det en
datorspelare tar den första människa som dör över efter den. Med Tab väljer den spelaren en
varelse och med pilarna flyttar den ett rum per runda (till vanliga rum
utan annan varelse, också till rum som ingen har hittat än). Mellanslag låter varelserna vila. Rundan
väntar på det valet som på alla andras.

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

Du har 6 liv. Den som når noll liv dör.

Ett vanligt rum kan sökas igenom en gång (S). Man kan hitta skatt, en dryck,
ett föremål, en hemlig gång till ett grannrum, eller utlösa en fälla. En
dryck (D) ger tillbaka 2 liv och tar bort gift.

### Föremål

Man bär högst tre föremål, och de står under DU i panelen. Hittar man ett
när väskan är full blir det kvar i rummet, och den som kommer in där
härnäst tar det. Den som dör tappar sina föremål i rummet. Föremålen följer
inte med till nästa torn. Krönikan berättar vem som hittade vad.

De här verkar av sig själva:

| Föremål | Vad det gör |
|---|---|
| Rep med kroken kvar | faller du genom golvet hamnar du inte i källarvalven (en gång) |
| Fjällamuletten | nästa fälla slår inte till, men tecknet syns (en gång) |
| Extra veke | lyktan tål ett mörker (en gång) |
| Bucklig sköld | varelser gör en skada mindre, i tre träffar |
| Väskan med dubbel botten | hälften av det du förlorar när ormen vaknar eller någon river guld är gömt |
| Kompass utan nord | visar åt vilket håll närmaste port ligger |
| Sövande flöjt | drar du ett kort som väcker ormen får du dra om (en gång) |

De här används med 1, 2 eller 3 (platsen i väskan), eller ett tryck på
föremålet i panelen. Föremål som pekar åt ett håll väntar på en pil. Att
använda ett föremål är rundans drag, och mitt i en strid får varelsen ett
fritt drag, utom mot rökpulvret.

| Föremål | Vad det gör |
|---|---|
| Rökpulver | fly genom en dörr utan att varelsen hinner slå (en gång) |
| Spegelskärva | se in i rummet bakom en dörr: fällan, varelsen och vandrarna där (tre gånger) |
| Kastkniv | träffa en varelse i rummet bakom en dörr; missar den blir kniven liggande där |
| Krita | rita en pil på golvet som alla som kommer in ser (fem gånger) |
| Timglas av ben | solen står still en runda för alla (en gång) |

Förbannade föremål ser värdefulla ut och är värda guld när man kommer ut,
men man märker baksidan först när man bär dem:

| Föremål | Värt | Baksidan |
|---|---|---|
| Den svarta pärlan | 40 | varje näve ur kammaren lägger ett vakna-kort i leken |
| Myntet med hål | 25 | varelserna anfaller oftare |
| Kronan av rost | 30 | den som står i grannrummet ser dig |

Ibland hittar man ett brev från någon som var i tornet förut. Det läses
direkt och följer avslöjandets nivå. Breven står under `brev` i
`foremal.json`.

Namn, texter, hur ofta de dyker upp och hur länge de håller står i
`data/foremal.json`, och hur ofta man hittar ett när man söker under
`sok` → `vikter` → `foremal` i `fallor.json`.

Om du laddar om sidan och skriver samma namn igen får du tillbaka din vandrare.

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
- Var fjärde runda, och vid varje nytt steg, får varje spelare en egen
  ledtråd. Ingen får samma som någon annan, så spelarna måste prata med
  varandra för att få hela bilden.
- På nivå 4 får varje spelare nya förklaringar till de fällor och varelser
  hen har mött (fallgallret var tänderna, de roterande rummen tarmarna ...).
- När spelet är slut säger tornet hemligheten rakt ut, och porten stängs.

### Att viska hemligheten

Från nivå 2 kan en människa viska till tornet med V (eller knappen Viska),
en gång per nivå. Viskningen skrivs i en egen ruta och syns bara för
servern och spelledaren. Stämmer den med något av svaren under `svar` i
`avslojande.json`, som hela ord, blir vandraren **invigd**. Tornet svarar
bara "Det vet att du vet", och de andra hör att det blev tyst någonstans.

Den invigde:

- hör dunket och vet åt vilket håll kammaren ligger,
- känner igen hälften av fällorna och kliver förbi,
- kan en gång viska i kammaren och lugna ormen (ett vakna-kort mindre),
- får en fjärdedel mer guld när hen kommer ut,
- och när solen går ner väntar porten en runda till för de invigda.

Men varelserna anfaller den invigde oftare, och en näve ur kammaren väcker
ormen lättare. Siffrorna står under `invigd` i `avslojande.json`.
Spelledaren ser vad var och en har viskat, och kan inviga eller ta
tillbaka det. Krönikan och epilogen nämner den som förstod.

Med `-avslojande=false` stannar tornet på nivå 1: inget dunk, inga
pulserande väggar och ingen epilog. Då går det inte heller att viska, och
ingen blir invigd.

## Ändra innehåll

Allt innehåll ligger i mappen `data/`. Filerna läses när servern startar.
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

- `texter.json` rum, beskrivningar, allmänna texter och introts väntetid
- `fallor.json` regler (liv, chans för händelse, drycker, hur länge spåren syns), fällor och sökresultat
- `ormen.json` ormkorten, skatten i kammaren och vad som händer när ormen vaknar
- `monster.json` varelserna (liv, skada, guld, vanor, hur ofta de vandrar, texter) och striderna
- `ai.json` datorspelarnas namn, personligheter och betänketid
- `svarighet.json` förvalen för svårighet (se nedan)
- `klasser.json` klasserna: namn, styrka, svaghet, liv, drycker, siffror och texter
- `foremal.json` föremålen: namn, beskrivning, hur ofta, hur länge och texter
- `avslojande.json` nivågränser, ledtrådar, texter vid nivåbyte, nya
  förklaringar till fällorna och epilogen. Rumsbeskrivningarna för varje
  nivå (`"1"` till `"4"`) ligger i `texter.json`.

### Svårighet

`svarighet.json` har fyra förval: Lätt, Normal, Svår och Mardröm. Ett förval
kan ändra antal liv (`liv`), chansen för en fälla (`fallchans`) och för en
varelse (`varelsechans`) i ett nytt rum, antalet sovkort (`sovkort`) och
vakna-kort (`vaknakort`) i ormleken, skadan när ormen vaknar (`ormskada`)
och hur stor andel av skatten som går förlorad (`skatt_forlust`). Det som inte
står med i ett förval tas från de andra datafilerna. Normal står därför tom,
och då gäller filerna som de är. `standard` är förvalet som gäller när inget
annat är valt. När vandringen börjar får alla veta svårigheten och hur många
liv de har.

| Förval  | Liv | Fälla | Varelse | Sovkort | Vakna-kort | Ormskada | Förlust |
|---------|-----|-------|---------|---------|------------|----------|---------|
| Lätt    | 8   | 35 %  | 12 %    | 17      | 2          | 1        | 35 %    |
| Normal  | 6   | 45 %  | 20 %    | 13      | 2          | 2        | 50 %    |
| Svår    | 5   | 55 %  | 28 %    | 11      | 2          | 2        | 60 %    |
| Mardröm | 4   | 65 %  | 35 %    | 9       | 3          | 3        | 75 %    |

För varje fälla kan du ändra namn, vikt (hur ofta den dyker upp), skada,
antal rundor, tecknet på kartan och texterna. Vad fällan gör styrs av dess
`id` i `events.go`. I texterna kan du använda
`{namn}`, `{horn}`, `{riktning}`, `{region}` och `{rum}` där de passar.
Beskrivningarna för rummen ligger under nivåerna `"1"` till `"4"`.

## Struktur

- `main.go` startar webbservern, läser flaggorna och skriver ut adressen
- `game.go` tornet, rummen, vandrarna, rundorna och vad var och en ser
- `settings.go` spelledarens inställningar och hur de sparas
- `host.go` värdsidan och dess adresser
- `gm.go` spelledarens översikt
- `ingrepp.go` spelledarens ingrepp
- `activity.go` autostopp för ett bortglömt torn
- `takeover.go` AI tar över en frånvarande vandrare
- `pvp.go` strid mellan vandrare
- `events.go` fällor, sökande och drycker
- `serpent.go` skattkammaren och ormkorten
- `monster.go` varelser, strider och den döda spelarens styrning
- `reveal.go` avslöjandet: nivåer, ledtrådar och epilog
- `chronicle.go` krönikan över vandringen
- `sounds.go` spår och ljud från grannrummen
- `summary.go` rundans sammanfattning, kortet efter varje runda
- `difficulty.go` förvalen för svårighet
- `classes.go` klasserna och valet av klass
- `items.go` föremålen
- `secret.go` viskningen och den invigde
- `ai.go` datorstyrda vandrare
- `data.go` läser in datafilerna
- `validate.go` kontrollerar datafilerna när tornet startar
- `reload.go` läser om datafilerna från värdsidan
- `simulate.go` simulerar vandringar med bara AI-spelare (`-simulera`)
- `ws.go` en liten websocket-implementation utan beroenden
- `web/` klienten: HTML, CSS och JavaScript utan ramverk
  - `index.html` och `app.js` spelarsidan: introt, namnrutan, klassvalet
    och spelet
  - `spelledare.html` värdsidan
  - `map.js` ritar kartan som krita på en skiffertavla: väggarna är
    handritade streck, vandrarna prickar i sin färg (den egna med en ring)
    och varelser och märken tecken. Kartan visar bara det kända och
    följer rummet vandraren står i. Den antecknar
    också vilka tecken den ritar, och spelarsidan gör teckenförklaringen
    av dem
  - `audio.js` allt ljud, skapat med Web Audio API, även stormen i introt
  - `style.css` utseendet

## Testa lokalt

Starta servern och öppna `http://localhost:8080` i två flikar med olika
namn. För ett snabbt test passar `go run . -sol 8 -turtid 10s`.

Spellogiken har tester:

```bash
go test ./...
```

### Simulera

För att se hur siffrorna i `data/` verkar kan tornet spela många
vandringar med bara AI-spelare, utan webbläsare och utan väntetider:

```bash
go run . -simulera 500 -ai 4 -svarighet svar
```

Det tar under en sekund. Utskriften säger hur många som kom ut med och
utan skatt, dog eller blev kvar i mörkret, hur mycket skatt de bar, vad de
dog av, hur ofta ormen vaknade och hur det gick för varje klass.
Inställningarna är desamma som när tornet öppnas (`-sol`, `-storlek`,
`-synlighet`, `-klasser` och de sparade). Vandringarna får frön i följd
från `-fro`, så samma kommando ger samma resultat, och en enskild vandring
kan köras om med `-simulera 1 -fro` och dess frö. AI-spelarna väljer klass efter sin personlighet, så
klasserna fördelas inte jämnt mellan dem.
