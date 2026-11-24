# Ändringar

Det som har ändrats i spelet, nyast först. README beskriver hur allt
fungerar. Här står vad som är nytt och vad det betyder för den som spelar
eller leder tornet.

## 2026-11-24

### Guldboken och en enklare slutskärm

- Slutskärmen börjar med placeringen. Välj en vandrare för att läsa hens
  guldbok: en längre berättelse om hur hen fick sitt guld, det första
  fyndet, nävarna ur kammaren, varelserna, slagen, det ormen tog tillbaka
  och det porten lade till (docs/guldboken.md).
- Knappen Runda för runda i guldboken visar en rad för varje gång guldet
  ändrades, med plus eller minus och vad vandraren bar efteråt.
- Resultaten och raden om var och en är nu samma lista, med titeln och
  den största stunden på vandrarens rad.
- Hela krönikan runda för runda står inte längre på slutskärmen; berättelsen
  och guldböckerna säger samma sak. Knappen heter Visa din vandring och
  visar bara den egna vandringen.
- Tornets dunk, hjärtslagen i ljudet, tystnar på slutskärmen.

### Berättelsen om vandringen

- Krönikan i slutet börjar med en gemensam berättelse i stället för
  vandringens största stunder: vilka som gick in, de viktigaste
  händelserna med var sitt spelkort och vem som kom ut
  (docs/berattelsen.md).
- Höjdpunkterna väljs efter vikt: ett uppvaknande, den första invigde och
  en död väger mest, och en händelse om den som vann väger mer. De
  berättas med när de hände ("Två rundor senare …"), och det som hände
  flera gånger sägs en gång ("Den natten vaknade ormen tre gånger").
- Hela krönikan runda för runda och Din vandring ligger bakom knappen
  Visa alla rundor.
- Under berättelsen står en rad om var och en: en titel till den som
  ensam har mest av något (den girigaste, rövaren, jägaren, upptäckaren
  …), hens största stund och hur det slutade.

### Minikartan

- En liten karta med allt du har upptäckt, vriden som kartan så att norr
  pekar snett uppåt höger: rummen, dörrarna, portarna, kammaren och var du
  står, och en streckad ram för det kartan visar (docs/minikartan.md).
- På datorn står den uppe till vänster på kartan. Ett klick växlar mellan
  din zoom och Allt, och K döljer och visar den.
- På telefonen står den i kompassens ställe nere till vänster, och ett
  tryck gör den stor över kartan.
- Kartritaren ser de andras rum streckade och kryss där någon har dött.

### Ormen kan ligga nästan var som helst

- Kammaren ligger inte längre nära mitten, utan nästan var som helst: inte
  i den yttre ringen, inte närmare någon port än tre rum, och ingen har den
  mer än tre rum närmare än de andra (docs/ormens-plats.md). Solen räcker
  längre när kammaren ligger långt bort.
- Alla hör åt vilket håll någon tar en näve ur ormens hög, utom när Tjuven
  tar.
- Står du bredvid kammaren hör du ormen andas genom väggen. Från nivå 2
  hör alla åt vilket håll andetagen kommer, och golvet i korsvägarna
  sluttar mot ormen. Introt säger att ormen sover någonstans i tornet.
- AI-spelarna följer andetagen och nävarna. Simuleringen visar hur ofta
  den som startade närmast kammaren vann: 21 % mot 17 % i snitt.

### Tornets händelser

- Från runda 4 kan något hända i tornet i slutet av en runda, med var sitt
  spelkort: tornet andas (en ny dörr), guld ur sprickan, ljus genom
  gluggarna, ömsat skinn (ett nytt förbannat föremål), ett ras, den döde
  vandrar och klockslaget (docs/handelser.md).
- Kartritaren har ingen egen dryck längre, och Ormtjusaren har Kraft 2 och
  Försvar 3, efter simuleringen med den slumpade kammaren.
- Webbläsarprovet säger till om tornet avslutas under provet, med de
  sista raderna det skrev.

### Breven pekar mot något

- Ett brev har nu en rad i marginalen som pekar mot något i just det här
  tornet: åt vilket håll ristningen finns om du inte har läst den, annars
  närmaste varelse, och annars kammaren.
- Brevet om att inte gå till mitten "första dagen" säger nu att man ska
  vänta tills man har hört andetagen, eftersom en vandring bara är en dag.

### Ormen flyttar sig när den vaknar

- När ormen har vaknat kryper den till en ny kammare, ett annat av rummen
  i mittens 3 × 3. Den gamla kammaren blir ett vanligt rum med lite guld
  kvar på golvet, som den som kommer först tar. Alla i tornet hör att
  ormen flyttar sig, och den invigdes hjärtslag pekar mot den nya platsen.

### Kammaren ligger nära mitten

- Kammaren ligger inte längre alltid mitt i tornet, utan i ett av rummen i
  mittens 3 × 3, slumpat efter vandringens frö (`spridning` i
  `data/ormen.json`; 0 ger mitten som förut). Den invigdes hjärtslag och
  AI-spelarna letar efter den.

### (du) på slutskärmen

- I resultaten på krönikans skärm och i resultatrutan står "(du)" efter
  ditt eget namn, och din rad är fet.

### Tangenterna vid dörrarna går alltid att trycka på

- På telefonen kunde en tangent vid en dörr hamna under översta raden
  eller under det stora kortet, så att trycket inte gick fram. Nu hålls
  tangenterna inom den del av kartan som syns, och flyttas in till kanten
  åt samma håll när dörren ligger utanför. På telefonen ligger de över det
  stora kortet.
- Webbläsarprovet säger varför ett tryck inte gick, när inget låg över.

### Placering för alla, och vinken jämför

- I slutet får alla en placering: först de som kom ut, efter skatt, sedan
  de som blev kvar inne och sist de döda. Den står framför namnet i
  resultaten, och var och en får raden "Du kom på plats 2 av 5".
- När du möter någon säger vinken också om hen bär mer eller mindre än du,
  till exempel "Något klirrar i packningen. Det är mindre än du bär."

### Var och en för sig

- Kartritaren ger inte längre bort sin karta till den hen möter
  (docs/individuell-tavling.md).
- Timglaset stoppar inte solen för alla. När solen går ner vänds det, och
  porten väntar en runda till på den som bär det, som på de invigda.
- Den invigde som viskar i kammaren tar inget vakna-kort ur leken för
  alla. I stället väcker hens egen nästa näve inte ormen.
- Introt säger att ni kan berätta för varandra vad ni har sett, men att
  ingen måste säga sanningen.
- I simuleringen vinner Kartritaren nu 20–22 % (19–21 % förut), när de
  andra inte längre får hens karta.

### Striden i balans, och panelen Nu säger hur den ser ut

- Simuleringen har körts med den nya striden, och talen är justerade så
  att alla klasser vinner mellan 14 och 21 % (docs/bara-anfall.md,
  Balansen). Den grönaktiga skepnaden slår lite lösare (Kraft 2) och
  gestalten i fjällrustning går lättare att träffa (Försvar 2).
  Kartritaren har Kraft 0, och Tjuven, Ormtjusaren och Lyktbäraren har
  Kraft 2 och Försvar 2.
- I strid säger panelen Nu hur lätt du träffar och hur lätt motståndaren
  träffar dig: "Du träffar lätt · den träffar ibland". Talen syns aldrig.
- Krönikan säger "en handsbredd" i stället för "1 handsbredder" när solen
  står lågt.

### Föremål som gör en starkare i strid

- Föremål kan nu höja Kraft och Försvar (`kraft` och `forsvar` i
  `data/foremal.json`).
- Nya föremål: **Rostigt svärd** (Kraft +1) och **Fjällbrynja** (Försvar
  +2, men den är tung, så slaget mot ryggen träffar alltid när man flyr).
- Nytt förbannat föremål: **Svärdet med rött fäste**, värt 15 guld och
  Kraft +2, men varje runda i strid då man inte anfaller skär det en och
  tar ett liv.
- **Den buckliga skölden** ger Försvar +1 i stället för en skada mindre,
  och spricker som förut efter tre träffar från varelser.
- Kraft och Försvar står bland tecknen (⚔ och ◆), gröna när ett föremål
  höjer dem, och på klasskorten.

### Striden har bara anfall

- Parera och undanmanöver är borta. Det enda stridsvalet är anfall (A):
  då slår både du och varelsen, och var och en träffar med sin chans
  (docs/bara-anfall.md).
- Chansen kommer av Kraft och Försvar. Varje klass och varelse har sina
  (`kraft` och `forsvar` i `data/klasser.json` och `data/monster.json`),
  och reglerna står under `traff` i `data/monster.json`.
- Står en varelse i rummet kan du ändå söka, dricka, vänta eller gå, men
  då kanske den anfaller, så ofta som dess lynne säger (`lynne`). Att
  vänta är att hålla garden, och då träffar varelsen sämre. Den som söker
  eller dricker träffas lättare. Att söka under en varelse gick inte förut.
- Mellan vandrare gäller samma sak: anfaller båda slår båda, och den som
  väntar håller garden.
- Legoknekten tar inte längre ingen skada när båda anfaller, utan slår
  hårdast av alla. Kartritarens svaghet är nu att slaget mot ryggen alltid
  träffar när hen flyr. Den invigde och Myntet med hål gör varelserna mer
  lynniga.
- Knapparna Parera och Undan och tangenterna P, U och E är borta. I strid
  står Vänta (håll garden), Anfall, Sök och Drick.

### Figurerna ser målade ut

- Vandrarna, varelserna, socklarna och ormen har fått en tunn mörk
  kontur, som målade spelfigurer. Det som vetter uppåt fångar ljuset,
  det som vetter nedåt är mörkare, och kanterna som vänder sig bort
  mörknar (`konturera` och `malad` i `klient/src/modeller.ts`).
- Konturen ritar varje del en gång till. När kartan går trögt slår den
  av konturen tillsammans med skuggorna.

### Ormen ligger på guldet

- I kammaren ligger nu ormen själv hoprullad uppe på guldhögen, med
  huvudet vilande framtill och ögonen slutna. Förut syntes bara två ögon
  över guldet när den vaknade.
- Den andas. Den som står i kammaren ser andningen följa sömnmätaren:
  långsamt när ormen sover djupt, snabbare när den sover lätt, och när
  den nästan vaknar rör sig huvudet oroligt. Andra ser den andas lugnt,
  så att den inte säger mer än sömnmätaren.
- När ormen vaknar reser den huvudet, öppnar lysande gula ögon och
  munnen, smakar luften med tungan och lägger sig sedan igen (`orm` i
  `klient/src/modeller.ts`).

### Figurerna slåss och faller

- På kartan gör figurerna det de gjorde under rundan, när de har gått
  klart. Den som anfaller gör ett utfall mot den den slåss med. Den som
  parerar drar sig bakåt, och den som viker undan kliver åt sidan. Den
  som tar skatt eller söker går ner på knä, och den som dricker lutar
  huvudet bakåt. Den som blir träffad rycker till, och en röd ring syns
  på golvet. Varelserna gör samma sak.
- Den som dör faller framåt och blir liggande på sockeln, i stället för
  att vridas på sidan. Varelserna faller på var sitt sätt: det bleka
  mångbenta viker in benen, det långa sjunker ihop i dammet, gestalten i
  fjällrustning faller baklänges och tappar skölden, den grönaktiga
  skepnaden stiger och bleknar bort, och väsendet av ben rasar ihop till
  en hög.
- Tillståndet säger vad varje figur gjorde och om den blev träffad
  (`gjorde`, `traffad`), bara för de figurer man ser, och vilka varelser
  som föll under rundan (`fallna`). `TestMovesAndFallenForTheMap`.

### Figurerna går att känna igen

- Varje klass har fått en egen kropp på kartan, så att klasserna går att
  skilja åt även uppifrån. Legoknekten är bred med hjälm och axelskydd.
  Tjuven är smal och framåtböjd med en lång kappa. Örtkvinnan har kjol
  och Pilgrimen en kåpa ända ner till golvet. Kartritaren bär ett koger
  fullt av kartor på ryggen och Ormtjusaren en turban. Lyktbäraren har
  lyktan högt på en stång, och Spårfinnaren lutar sig fram mot marken
  (`FORM` och `vandrare` i `klient/src/modeller.ts`).
- Alla vandrare har större huvud, ögon som visar vart de tittar och en
  mantel över axlarna med en kappa på ryggen i vandrarens egen färg. Då
  syns färgen även när hatten eller huvan har klassens färg. Kappan
  fladdrar när vandraren går.

### Varelsen skymmer inte den första vandraren

- På kartan snett ovanifrån stod varelsen precis framför den första
  vandrarens plats, sett från kameran, så att en hög varelse som
  fjällrustningen skymde figuren helt under striden. Varelsen står nu
  snett framför rummets mitt (`VARELSE` i `klient/src/karta.ts`), där
  ingen vandrare står rakt bakom den.

### Det övriga från provspelet

- **Solens varningar i rätt ordning.** Med en kort sol (`-sol` under 8)
  kom "Solen nuddar nästan horisonten" före "halva vägen". Varningarna om
  halva och en fjärdedel hoppas nu över när de skulle komma vid tre rundor
  eller färre (`sunWarning` i `game.go`, `TestShortSunWarningsInOrder`).
- **Striden.** Varelsen "försöker parera" och "försöker vika undan" i
  stället för "parerar" och "gör en undanmanöver", så att det inte låter
  som att den lyckades när man sedan träffar den (`val` i
  `data/monster.json`, på engelska "an attempt to parry/dodge"). Dricker
  man eller flyr medan varelsen inte anfaller står det "Du är inte beredd,
  men … låter dig vara" och "… följer inte efter" (`inte_beredd_lugn` och
  `fly_lugn`, `TestFreeAttackCalmWhenNotAttacking`).
- **Din vandring** säger "Du satt fast" om en runda i en grop eller bakom
  ett galler, i stället för "Du stod kvar" (`sam_tyst_fast`,
  `TestJourneyStuckRound`).
- **Kartritaren** får höra att kartan fylls i med varje rum som någon
  hittar, i stället för att den bara visar de egna rummen
  (`synlighet_kartritare`, `TestCartographerIntro`).
- **Ett känt, tomt rum** ger raden "Du känner igen rummet: …" med rummets
  korta beskrivning, så att loggen säger att draget blev av och var man
  står (`har_igen` i `here.go`).
- **Tecknen på kartan** har en gemensam mörk botten, så att facklor och
  murar inte lyser igenom. Förklaringen till namnets delar säger att de
  står i den ordning man hittade dem, inte i namnets ordning.
- **Spelkortet** täcker inte längre den egna figuren på en bred skärm:
  kartan håller kortets plats till höger fri, som den gör för loggen till
  vänster (`tacktHoger` i `klient/src/karta.ts`, `coverRight` i
  `web/map.js`). På en skärm smalare än 760 punkter är det som förut.
- Fjällrustningen mot Kartritaren är prövad i simuleringen och lämnad som
  den är; se `docs/att-gora.md`.
- Provspelet står i [provspel-2026-11-24.md](provspel-2026-11-24.md).

## 2026-11-23

### "En runda kvar" när solen hoppar

- Webbläsarprovet hittade "1 rounds left" på engelska: när klockfällan
  får solen att hoppa ner till en runda sa varningen "1 rundor kvar".
  Nu säger den "En runda kvar" (`sol_snart_en` i `data/texter.json` och
  `data/en/texter.json`, `TestSunWarningOneRound`).
- Webbläsarprovet är stadigare: det väntar högst tre sekunder på något
  som inte går att trycka på och skriver då på en rad med `OBS` vad som
  låg över, det läser inte en rad som loggen fortfarande skriver ut, och
  det räknar inte med att rummet berättas igen i mörkret.
- På telefonen låg spelkortet eller det stora kortet ibland över
  tangenterna vid dörrarna, och ibland låg en tangent under tecknen eller
  solmätaren. Det står i `docs/att-gora.md` bland det som ska frågas om.

### Webbläsarprovet

- `node verktyg/webbprov.mjs` spelar en kort vandring på en dator och en
  telefon i Chromium och letar efter det robotarna inte ser: fel i
  skripten, trasiga texter på sidan, en sida bredare än skärmen,
  tangenterna vid dörrarna som ligger över varandra, minikort som inte
  öppnas, mellanslag som inte spelar efter kortvisningen, en omladdning
  som inte berättar rummet igen och vandrare som saknas i resultatet.
  Det tar en bild av varje fas. README beskriver det under
  Webbläsarprovet.
- Provet hittade inga fel på svenska eller engelska. Med felet från
  provspelet 2026-11-22 tillbaka (mellanslag efter ett klick på ett
  minikort) märkte det felet.

### Robotarna hittar fler trasiga texter

- Provspelarna (`verktyg/provspelare.mjs`) märker nu också "1 rundor",
  "1 lives", "there are at least 1 of you", samma ord två gånger, liten
  bokstav efter punkt och svenska ord i en engelsk vandring. Ingen av
  dem märktes förut, och flera sådana fel hittades i provspelet
  2026-11-22 för hand.
- De kontrollerar alla texter som tornet skickar, inte bara loggen:
  rundans kort, spelkorten, det stora kortet, krönikan, epilogen och
  stunderna. En trasig text skrivs en gång, på en rad som börjar med
  `MISSTÄNKT`.
- Reglerna står i `verktyg/textkontroll.mjs`, och
  `verktyg/textkontroll.prov.mjs` provar dem mot felen från provspelen i
  CI.
- 36 vandringar med tre robotar och två AI-spelare (18 på svenska, 18 på
  engelska på Mardröm) gav ingen märkt text.

### Mer guld i kammaren på de svårare nivåerna

- En näve ur skattkammaren ger mer guld på Svår (19–56 i stället för
  15–45) och mest på Mardröm (23–68), så att det lönar sig att våga mer
  (`kammarskatt` i `data/svarighet.json`).
- AI-spelarna tar sällan en näve på de svåra nivåerna (0,1–0,3 per
  vandring), så simuleringen visar knappt någon skillnad. Det får ett
  provspel med människor visa.

### Skada och sol efter svårigheten

- En svårighet kan nu göra fällor och varelser farligare (`skada` i
  `data/svarighet.json`) och solen längre eller kortare (`sol`).
- **Mardröm** har 6 liv i stället för 4, men varje fälla och varelse
  som skadar tar ett liv mer. Med 4 AI-spelare blir det ungefär lika
  svårt som förut (39 % ut med skatt, 44 % döda), men varje träff känns
  mer.
- **Lätt** har en fjärdedel längre sol, 50 rundor i stället för 40 i ett
  torn på 9 × 9. På Lätt var det fler som blev kvar i mörkret än som dog
  (8 och 6 %); nu blir 4 % kvar.

### Normal har 7 liv

- På Normal har var och en 7 liv i stället för 6 (Legoknekten och
  Lyktbäraren 8). I provspelet tog fällorna 3–6 av 6 liv inom tio rundor.
  I simuleringen gör fällchansen nästan ingen skillnad; det är antalet
  liv som avgör. Med 4 AI-spelare kommer 68 % ut med skatt och 13 % dör
  (förut 66 och 15 %).

### Ormens namn: ordningen och fel viskningar

- Texten för en del av namnet säger inte längre om den är den första,
  andra eller tredje. Den säger var på en orm tecknet sitter: vid huvudet,
  i buken eller i svansen. Ordningen får man lista ut själv.
- Tecknen på spelskärmen står i den ordning du hittade delarna.
- Viskar du rätt delar i fel ordning säger tornet det, men försöket är
  förbrukat.
- En fel viskning gör ormen oroligare: ett vakna-kort till i leken.
- Den som dödar en vandrare och tar hens delar får veta var på ormen de
  sitter.

### Ormens namn finns i tornet

- Tornet delar inte längre ut namnets delar. Var och en hittar dem själv,
  på tre ställen: **ristad i väggen** i ett rum en bit in i tornet, **under
  fjällen** på en varelse man träffar, och **i kammaren**, där ormen väser
  sin del i sömnen när man tar en näve. Vilken del som finns var lottas.
- Från nivå 3 får den som inte har hittat ristningen veta åt vilket håll
  den finns. Spårfinnaren får veta det redan från nivå 2.
- På nivå 4 får den som har hittat något veta vems namnet är.

### Ormens namn: var och en för sig

- Delarna av ormens namn lottas för var och en för sig. Tornet fördelar
  dem inte längre så att gruppen tillsammans vet mest.
- Ledtrådarna dras ur samma texter för alla, var och en ur sin egen lek.
  Två vandrare kan alltså få samma ledtråd.
- **Tornets nåd** är borta. I stället får **den förste invigde** en
  fjärdedel mer guld till när hen kommer ut, och krönikan säger att ingen
  hann före.

### Avslöjandet går inte att stänga av

- Raden Avslöjandet på värdsidan och flaggan `-avslojande` är borta.
  Tornet bär alltid på sin hemlighet, med invigningen, varslen och
  tornet som förändras.
- Introt har alltid tornets antydningar och dunk; ryktena för den som
  redan kände hemligheten (`rykten_invigda`) är borta.
- Sparade inställningar med `avslojande` läses som förut; fältet struntas i.

### Telefonen: mer karta

- Kartan fyller höjden mellan översta raden och knapparna, och går lite
  närmare när den följer dig, så att rummet du står i blir större.
- Tecknen (klass, liv, skatt, drycker, föremål) står på en rad under
  solmätaren, med bara tecken och tal. Ett tryck förklarar, som förut.
- Du går genom att trycka på pilarna vid dörrarna på kartan. Pilknapparna
  är borta, och knapparna är en rad med Vänta och det som gäller just nu.
  Zoomar du ut till Långt eller Allt kommer pilknapparna tillbaka.
- Under kartan står bara den senaste rundan, tätare. Knappen Loggen
  öppnar hela loggen.
- På datorn går pilarna vid dörrarna att klicka på med musen.
- En ny zoomnivå närmast, **Rummet**, bara på telefonen: rummet du står i
  fyller bredden. Telefonen börjar där, om du inte redan har valt en
  annan nivå.

## 2026-11-22

### Benluren i stället för kritan

- Kritan är borttagen; pilarna på golvet var till liten hjälp.
- Det nya föremålet **Benluren** (☊) verkar av sig själv och tar aldrig
  slut. Kommer du in i ett rum hör du om en varelse står i något av
  grannrummen, och åt vilket håll, men inte vad det är. Ljud som andra
  bara hör från grannrummen, som strid, fällor, en port, en dryck och
  skriket när någon dör, hör du också två rum bort, med väderstreck. Med
  benluren hör även Spårfinnaren genom väggarna.
- Kartritaren har en dryck med sig i stället för de tre extra strecken
  med kritan.

### Viskningen

- Viskar man bara ormens namn känner tornet igen det: "Väggarna rycker
  till när du säger namnet … Men vad är det du vill säga om det?" Försöket
  är inte förbrukat, och V finns kvar. Förut svarade tornet som på en fel
  viskning, och då gick det inte att viska igen förrän tornet förändrades.
- Svaret på en viskning står i hemlighetens färg i loggen i stället för
  som en svag rad, så att det syns.

### Den dödes guld

- Den som dör tappar sitt guld i rummet, som föremålen. Förut gick guldet
  förlorat. Den som går in i rummet, eller står kvar där när rundan är
  slut, tar det, med spelkortet Det som blev kvar.
- Krönikan säger vem som tog upp vems guld, och den döda får i slutet
  veta om guldet kom ut och med vem.
- Den som dödar en annan vandrare tar som förut allt guldet direkt, så
  det blir inget kvar efter en strid mellan vandrare.
- Skriket säger var: alla hör i vilken del av tornet någon dog, och den
  som står i ett grannrum hör det åt ett håll. Bar den döda guld hörs
  mynten rassla efter skriket, men inte hur många.
- Rummet minns: den som kommer in där någon har dött läser "Här föll
  Alva" i rumsblocket, och "Någon har redan tagit dem" om guldet är taget.
- Kartritaren ser på kartan var någon har dött, med ett litet x i rummets
  hörn. Spårfinnaren ser guldet som ligger kvar, som föremålen, och läser
  i dammet vem som tog det. Hen hör fortfarande ingenting genom väggarna,
  så skriket säger henom bara del av tornet.
- AI-spelarna går efter guldet när skriket lät som mynt, med samma
  ledtrådar som en människa får, om de har minst halva livet kvar och gott
  om tid. Plundraren går oftast (`guldlust` 0,8), den djärva ibland (0,5)
  och den försiktiga sällan (0,15). På väg ut tar de en omväg på högst
  två steg till guldet, när de vet exakt vilket rum det ligger i.
- Se [den-dodes-guld.md](den-dodes-guld.md).

### Ormtornet.app i ~/Applications

- `./skapa-app.sh` skapar nu appen i `~/Applications`, där macOS letar
  efter appar, i stället för `~/Projekt/Lokalt`. En annan mapp kan
  fortfarande anges som argument.

### Provspelet 2026-11-22

Rättat efter en timmes provspel (se
[provspel-2026-11-22.md](provspel-2026-11-22.md)):

- Spelkortet: ett litet kort från en äldre runda, som Tre tecken eller ett
  varsel, visas inte längre efter en nyare rundas kort.
- Loggen: en stöt i muren, en tom port eller "Här finns ingen skatt att ta"
  börjar inte längre en ny runda i loggen. Raden hamnar under rundan som
  står i bild, så att rummet man står i inte försvinner under "tidigare"
  på telefonen.
- I mörkret syns inte dörrarna och porten i rummet man står i, och inga
  tangenter vid dörrarna: man får pröva sig fram, som reglerna redan sa.
- Knappen Drick syns bara när man har en dryck.
- Tangenterna vid dörrarna krymper när rummet är litet på skärmen och
  ligger inte längre över varandra på telefonen.
- Resultatet: den som gick ut tomhänt står före de döda och förlorade.
- Texterna: "Du sitter fast en runda" (inte "1 rundor"), "kammaren en
  gång", "+2 drycker", och föremålen med en eller ett: "Sigge hittade en extra veke",
  "Här ligger en krita … och får lämna fyndet kvar". På engelska står inte
  längre "−1 lives".
- Slutet: den som laddar om sidan, eller vars telefon somnar, står kvar i
  resultatet och får tillbaka sin vandring när hen kommer tillbaka.
- Den som laddar om sidan mitt i vandringen får beskrivningen av rummet
  hen står i igen, och svaret på en viskning hamnar under rundan som står
  i bild.
- Minikorten i loggen på kartan går att klicka på, och efter ett klick
  väljer mellanslag och Enter draget igen i stället för att öppna kortet.

### Länken först

- README leder nu först till länken på GitHub Pages: spelledaren öppnar
  den och skickar inbjudan till de andra. Att ladda ner spelet och köra
  tornet på en Mac står i ett eget avsnitt längre ner, för när länken inte
  räcker, till exempel på ett nät som stoppar anslutningen.

### Mac och Windows

- README säger nu att det krävs en Mac för att ladda ner spelet och köra
  tornet själv. Med en Windows-dator spelar man från länken på GitHub
  Pages, som spelledare eller spelare.

### Minikorten går att se

- I krönikan i slutet och i Din vandring är rundans kort små spelkort,
  104 pixlar breda, med pergament, bilden och rubriken.
- Ett klick på ett minikort, i loggen, på spelkortet eller i krönikan,
  visar kortet stort med sin text. Ett klick till eller Escape stänger
  det.
- På raden med siffror står ormkorten ihopslagna med namn och antal:
  "Ormen sover ×3 · Ormen rör sig".

### AI-spelarna och simuleringen

- AI-spelarna har en tredje personlighet, plundraren, som söker igenom
  rummen och helst låter kammaren vara. Den var den säkraste strategin i
  simuleringen.
- Ormtjusaren som AI tar alltid en näve när inga vakna-kort finns kvar.
- Ormen skadar inte Ormtjusaren när den vaknar, i stället för en skada
  mindre.
- Simuleringen (`-simulera`) visar andelen som vann och guld per
  vandrare, för varje personlighet och för varje klass och personlighet.

### Klasserna och den grönaktiga skepnaden jämnas ut

Efter en simulering av sex strategier i 1000 vandringar (se
[balans-2026-11-22.md](balans-2026-11-22.md)):

- Kartritaren börjar utan dryck, och Örtkvinnan med en dryck i stället
  för två. De vann oftast.
- Lyktbäraren har ett liv extra, och Ormtjusarens svaghet är att hen
  hittar hälften så mycket guld när hen söker, i stället för att
  varelserna anfaller oftare. De vann mest sällan.
- Den grönaktiga skepnaden har två liv i stället för tre och anfaller
  oftast (anfall 3, parera 1, undan 1), så att parera är rätt svar. Den
  stod för 30 % av dödsfallen.

### Klasskorten som spelkort

- Klasskorten i klassvalet har spelkortens pergament, ram och mörka
  text, med namnet i versaler överst, styrkan i grönt och svagheten i
  rött.
- Varje klass har en bild på kortet, i samma stil som spelkorten. Fältet
  `bild` finns nu också i `data/klasser.json`.
- Klassvalet börjar med korten på baksidan, och de vänds upp ett i taget.
  Ett kort som någon annan tar vänds tillbaka och visar baksidan med
  namnet på den som tog det.
- Under vandringen står klassen uppe till höger som ett litet klasskort.
  Ett klick visar klasskortet stort mitt på kartan, med bilden,
  beskrivningen, styrkan och svagheten.

### Vandringens största stunder

- Krönikan i slutet lägger fram vandringens största stunder som stora
  spelkort under resultaten, högst fem, som vänds upp ett i taget:
  vinnaren, uppvaknanden, de som dog, den som blev invigd, portar som
  stängdes och annat, i den ordning de hände. Vilka händelser som räknas
  och hur många står under `spelkort.stunder` i `data/texter.json`.

### Spelkortet

- Rundans kort och ormkortet är ett kort: **spelkortet**, som vänds upp
  till höger på kartan efter varje runda där något ändrade dig. Det har
  en bild och en rubrik för rundans viktigaste händelse, händelsens text
  och rundans rader under. Resten av rundans kort står som små kort
  längst ner. Utan eget kort står Runda N överst med en fackla.
- Tar du en näve är ormkortet rundans spelkort. Raden "Du bär nu på …
  guld" står inte längre i loggen; den står på kortet.
- Kommer uppvaknandet samma runda väntar spelkortet tills det stora
  kortet har gått undan.
- Fällorna har var sin bild på spelkortet, med den korta beskrivningen
  som rubrik. Klev du förbi fällan är kortet nedtonat. Fällorna har ett
  nytt fält `bild` i `data/fallor.json`.
- `verktyg/kortbilder.html` visar alla spelkortens bilder.
- Varelserna har var sin bild på spelkortet: när de dyker upp, i varje
  stridsrunda och överkorsad när de faller. Striden mot en annan vandrare
  ger båda ett kort. Fältet `bild` finns nu också i `data/monster.json`.
- Stridens rader i loggen säger inte längre "Du förlorar 1 liv" eller
  "Du har nu 140"; det står på spelkortet.
- Föremålen har var sin bild. Kortet säger vad föremålet gör, och loggen
  säger bara att du hittade det. Ett förbannat föremål vänds en gång till
  när det visar sin baksida. Drycker och guld får egna kort, och raderna
  om dem i loggen säger inte längre hur mycket du har.
- Källarvalven, brevet, varslen, namnets delar och att bli invigd har
  egna kort. Det som kommer mellan rundorna visas genast, efter kortet
  som står uppe.
- En port som stängs och solen, när tre rundor är kvar och när den går
  ner, vänds upp stort mitt på kartan för alla som är kvar inne, som
  uppvaknandet. Den som kommer ut ser sitt sista kort stort på kartan.
- Din vandring börjar med vandringens sista kort: ute, död eller kvar
  när solen gick ner.
- Raderna i loggen som hör till ett kort får ett litet kort framför sig,
  och Din vandring visar rundans kort som små kort vid varje runda.

### Fröet, hela bredden och loggen på kartan

- Värdsidan har ett nytt avsnitt, **Fröet**, efter Inför nästa vandring.
  Där står den här vandringens frö med Kopiera och, i slutet, Samma igen.
  Fältet Nästa vandring visar alltid fröet nästa vandring får; det går
  att skriva in ett eget eller slumpa ett nytt. Fröet står inte längre i
  Vid porten.
- På en bred skärm fyller kartan hela bredden. Loggen står som text utan
  ruta nere till vänster på kartan, med skugga runt bokstäverna, och
  kartan ställer sig bredvid den. Den senaste rundan står nederst och den
  förra svagt ovanför. Knappen Loggen, eller L, öppnar hela loggen i en
  låda. Kompassen står ovanför zoomen.
- På en surfplatta och en telefon på bredden står knapparna nere till
  vänster på kartan. En låg skärm får en smalare logg med bara den
  senaste rundan. Telefonen på höjden är som förut.
- Under vandringen öppnar märket spelledare värdsidan i en ny flik, så
  att spelledarens vandrare inte blir frånvarande. Inte på Pages.

### Zoomen glider

- Scrollhjulet och två fingrar zoomar kartan gradvis i stället för att
  hoppa mellan Nära, Mellan, Långt och Allt. Ett hack på hjulet är en
  tredjedels nivå, och en styrplatta ger en jämn zoom. Knapparna och Z går
  fortfarande en hel nivå i taget, och namnet visar den närmaste nivån.

### Loggen och krönikan, efter en genomläsning

- Rättat: spåret efter en fälla som nyss slog till står inte längre före
  fällan i rumsblocket. En varelse man fällde "är" inte "borta"; den som
  fällde den får inget, de andra i rummet ser den ligga fälld. Raden om
  svårigheten säger livet efter klassen (legoknekten har 7). ~ i
  resultaten förklaras också på krönikans skärm.
- Loggen upprepar mindre: en strid som pågår hörs bara första rundan,
  varelsen beskrivs helt första stridsrundan och sedan kort ("det
  mångbenta"), och att man tvekade sägs en gång, i loggen, och inte på
  kortet. Ett känt, tomt rum ger inget block och ingen tom rubrik. Under
  Vandringen börjar står porten sist.
- Din vandring säger vad skadan kom av, vad som dödade en och vilken
  varelse man fällde. Rundor där man bara stod kvar slås ihop till
  "Runda 10–12", och rundor då tornet styrde märks "tornet styrde".
- Krönikan säger hur mycket guld som gick förlorat när någon dog, vem
  som nådde skattkammaren först och när klockan fick solen att sjunka
  fortare. Raden med siffror visar bara det som inte är noll.
- Varelsen som väser heter "något långt och väsande" i de korta formerna,
  så att den inte blandas ihop med Ormen. När den dyker upp är den
  fortfarande en orm, lång som en man.

### Ljudet för spelledaren som går in själv

- Spelledaren som går in med Gå in själv hoppar över introt, och då
  startade ljudet aldrig. Nu startar det direkt, eller vid första tryck
  eller tangent om webbläsaren kräver det.

### Inga namn på rum, fällor och varelser

- Spelet säger inte vad rummen, fällorna och varelserna heter. Rummet
  står i loggen som sin beskrivning, utan namn i versaler, och rutan med
  rummets namn uppe till vänster på kartan visar bara Mörker eller
  Källarvalven.
- Raden med fällans namn när den slår till är borta; texten säger vad som
  händer. Där en text behöver nämna en fälla, ett rum eller en varelse
  står en kort beskrivning, till exempel "Strid: det bleka mångbenta",
  "Du ser tecknen i tid: ett galler av spetsiga järnstänger" och "Bo dog
  av pilar ur små hål i väggen".
- Krönikan och epilogen beskriver också i stället för att använda namn.
  Värdsidans Läs om innehållet har kvar namnen.

### Spelledaren spelar i samma flik

- Gå in själv på värdsidan byter till spelarsidan i samma flik, utan
  introt. På GitHub Pages byter ramen i skalet sida, så tornet går vidare.
- Spelledaren ser "spelledare" i översta raden. I lobbyn och slutet är det
  en knapp, Värdsidan, och i slutet står Till värdsidan bredvid Nytt torn.
  De andra vandrarna ser inte vem som är spelledare.
- Har spelledaren en vandrare under vandringen heter knappen på värdsidan
  Till spelet.
- Rättat: Din vandring och slutet kraschade på spelarsidan när en runda
  hade ett rum men inga rader på kortet.

### Vandringen börjar i det du ser och hör

- Portrummet beskrivs inte längre i lobbyn. När vandringen börjar kommer
  det första rumsblocket under en egen rubrik, Vandringen börjar: vid
  vilken port du kliver in står först i beskrivningen, solen står bland
  det du ser, och portarna där de andra går in är det du hör ("Långt
  borta, åt flera håll, knarrar portar").
- De lösa raderna om porten och solen är borta, och "Välj era drag" sägs
  bara i panelen Nu. Svårigheten och raden om synligheten står kvar under
  blocket.

### Adressen under Bjud in syns bara när man ber om den

- Värdsidan skriver inte ut adressen (på GitHub Pages länken) under Bjud
  in. Kopiera inbjudan och Kopiera adressen räcker oftast, och Visa
  adressen (Visa länken) visar den när den behövs. Går det inte att
  kopiera visas den av sig själv, så att den kan markeras.

### AI:n tar över den som står still

- Den som är kvar i tornet men inte väljer något drag, så att tiden går
  ut, i lika många rundor i rad som inställningen Frånvarande vandrare
  säger (två som standard) får en AI-spelare som styr. Vänta räknas som
  ett drag, och rundor i källaren räknas inte.
- Spelaren ser i loggen och i panelen Nu att tornet styr stegen. Nästa
  drag hen själv väljer ger tillbaka vandraren direkt och gäller redan i
  samma runda.
- Står inställningen på "står still" tar AI:n aldrig över.

### Rummet som ett stycke, och ett kortare kort

- Rummet står som löpande text i loggen, utan ruta: namnet i versaler,
  beskrivningen och det du ser i ett stycke, och det du hör sist. Din
  vandring visar rummet på samma sätt under varje runda.
- Kortet efter rundan säger bara det som ändrade dig: liv, guld,
  ormkortet, föremål, gift, mörker, vem du slogs mot och om ormen vaknade.
  Vart du gick och solen står inte längre där, och har inget ändrat dig
  kommer inget kort.

### Det du hör

- Ljuden under rundan, som stridsljud från ett grannrum, ekot, klockan
  och ormen som rör sig, står sist i rumsblocket i stället för som egna
  rader i loggen. Samma ljud sägs bara en gång, och står du kvar kommer
  ett block också när du bara har hört något. Din vandring visar dem
  också.

### Det du ser beskrivs, det förklaras inte

- Rumsblocket beskriver det du ser: varelsens egen text, "Sigge står i
  rummet och ser på dig.", "Kvällsljuset faller in genom porten." En fälla
  som har slagit till syns som ett spår, till exempel "Golvet är repat i
  långa cirklar.", utan fällans namn. Ett tomt rum säger ingenting.
- Står du kvar kommer bara det som har kommit eller gått, som
  "Fjällvakten är borta."
- Utropet med varelsens namn och den dubbla beskrivningen när en varelse
  dyker upp är borta.

### Tornet väljer motståndare, och Z för kartan

- Står flera vandrare i rummet väljer tornet vem du slåss mot: den som
  anföll dig förra rundan, annars den som kom först in i tornet. Panelen
  Nu säger vem. Tab och knappen Byt är borta. Skatten avgör inte, eftersom
  den är hemlig.
- Z går runt mellan kartans nivåer: Nära, Mellan, Långt och Allt. + − och
  M är borta; knapparna på kartan och scrollhjulet finns kvar.

### Knapparna följer läget

- På en pekskärm syns bara de knappar som går att använda just nu: Sök
  och Drick utanför strid, Anfall, Parera, Undan och Drick i strid, Viska
  när det går att viska och bara Redo i lobbyn och slutet. Styrningen
  bakom ? visar på samma sätt bara de tangenter som gäller.

### Sök tar också skatt

- S (knappen Sök) tar en näve ur högen i kammaren och söker i andra rum.
  T och knappen Ta är borta.

### Värdsidan i två kolumner

- På en skärm som är minst 1100 punkter bred står värdsidan i två
  kolumner under raden: Bjud in och Vid porten till vänster, Inför nästa
  vandring, Kommer någon inte in?, Läs om innehållet och Stäng tornet till
  höger. Förklaringen till varje inställning står under valet. På en
  smalare skärm är sidan en kolumn som förut.

### Ingen blir instängd

- Ibland stängde tornet in den som började vid en port, i ungefär vart
  tjugonde torn och ofta i bara fyra rum. Då gick det bara att komma ut
  genom en hemlig gång. Nu får ett nytt rum en dörr till mot ett okänt rum
  när den del av tornet det hör till annars skulle bli stängd, så att man
  alltid kan nå kammaren.

### Bara länken på värdsidan

- Värdsidan på GitHub Pages visar under Bjud in bara länken, utan
  QR-kod. Spelarna får länken genom Kopiera inbjudan eller Kopiera
  länken.

### Tornet ger namnen

- Ingen skriver längre sitt namn. "Gå in i tornet" leder direkt till
  lobbyn, och tornet ger var och en ett förnamn och ett tillnamn, till
  exempel "Dagny Ekfot". Loggen säger "Tornet kallar dig Dagny Ekfot."
- Två vandrare får aldrig samma förnamn. Namnen står i `data/namn.json`,
  med tillnamnen översatta i `data/en/namn.json`.
- Sidan kommer ihåg namnet, i fliken och i webbläsaren. Den som laddar om
  sidan eller tappar kontakten får samma namn och sin vandrare tillbaka,
  också till nästa torn. Efter 12 timmar utan spel glöms namnet, så att
  varje kväll ger ett nytt namn. Ett namn som skrevs in för hand i en
  tidigare version glöms direkt. En andra flik i samma webbläsare får ett eget
  namn.
- Släpper tornet inte in en, som när det är fullt eller inte svarar, står
  felet på en egen skärm med en knapp tillbaka.

### Föremålen verkar av sig själva

- Inga föremål används längre med 1–3. Tangenterna, att sikta med pilarna
  och att trycka på ett föremål för att använda det är borta. Ett tryck
  eller musen över föremålet visar fortfarande vad det gör.
- Rökpulvret tänds när du flyr ur en strid, och varelsen hinner inte slå.
- Spegelskärvan visar ett osett grannrum när du kommer in i ett nytt rum,
  tre gånger.
- Kastkniven kastas när en varelse dyker upp i ditt rum, och varelsen tar
  en skada innan striden börjar.
- Kritan ritar en pil åt det håll du går när du lämnar ett rum, fem rum.
- Timglaset vänds när solen har tre handsbredder kvar, och solen står
  still en runda för alla.
- Datorspelarna och robotarna använder inte längre föremål; föremålen
  verkar för dem på samma sätt.

### Ny layout på telefon

- Allt står i en kolumn: översta raden fast högst upp, kartan, raden Nu
  och loggen, och knapparna fast längst ner.
- Loggen rullar med sidan i stället för i en egen liten ruta. Den senaste
  rundan tar så mycket plats att de äldre bara syns när du rullar ner.
- En lång rumsbeskrivning visar tre rader, och "mer ↓" visar resten.

### Ny layout på dator

- Panelen Nu är en smal rad under kartan: läget och striden till vänster,
  ditt drag till höger. Viskningens ruta öppnas i samma rad.
- Loggen fyller hela högerkolumnen, som är något bredare på stora skärmar
  (30 procent av bredden, 300–440 punkter).

### Bara den senaste rundan i loggen

- Den senaste rundan fyller hela loggen. De äldre ligger under strecket
  "tidigare ↓" och syns bara när du rullar ner. Rundorna fälls inte
  längre ihop.
- Den korta raden med `≡` efter rundans kort är borta; kortet på kartan
  finns kvar.
- Din vandring visar rummet bredvid rundan ("Runda 4 · Brunnsrummet") och
  vad som fanns där.

### Rummet i loggen, utan teckenförklaring

- Överst i varje runda står ett rumsblock: rummets namn, beskrivningen
  när rummet är nytt, och i ord vad som finns där (varelsen, de andra
  vandrarna, skatten, porten, en fälla som har slagit till, föremål som
  någon har lämnat). Står du kvar kommer ett kort block bara när något
  har ändrats, som när en varelse har kommit in eller en vandrare har
  gått.
- Teckenförklaringen uppe till vänster på kartan är borta. Hur saker ser
  ut syns på kartan, och vad de är står i loggen.
- Bilderna på krönikans rader och på rundans kort i Din vandring är
  borta igen. Krönikan är bara text.

### Lobbyn utan listan över vandrare

- Panelen Vandrare är borta. I lobbyn och slutet står vandrarna som
  prickar med namn i översta raden, och en fylld prick betyder redo.
- Raden "Tryck Enter när du är redo" i panelen Nu är borta. I stället står
  knappen Redo i översta raden (Nytt torn i slutet), och den visar Redo ✓
  när du har tryckt. Enter fungerar som förut, och på en pekskärm är det
  knappen längst ner, som också visar Redo ✓.
- Panelen Nu syns bara när den har något att säga.

### Tecknen från kartan i krönikan

- Raderna i krönikan har samma bilder som teckenförklaringen: varelsen som
  föll eller dödade, fällan, föremålet, skattkammaren när ormen vaknade,
  porten när någon gick ut och prickarna i vandrarnas färger när de slogs.
  Ett tryck på bilden, eller musen över den, visar förklaringen.
- Rundans kort i Din vandring har samma tecken: varelsen eller vandraren
  du slogs mot, föremålet du hittade eller använde, porten när du gick ut
  och det som dödade dig.
- Kartan kan rita varelser och fällor som du aldrig själv såg, så att de
  syns i krönikan även när någon annan mötte dem.

### Striderna mellan vandrare i krönikan

- Krönikan skriver ner varje slag mellan vandrare som träffar men inte
  dödar: "Ada slog Bo och tog 12 guld." Ett bakhåll, när den andra inte
  försvarade sig, står som "Ada föll över Bo bakifrån." Ett slag som dödar
  står som förut bara som en död.
- Vandrarens rad med siffror i slutet visar träffar på andra vandrare och
  stulet guld, men bara hos den som har träffat någon.
- Rundans kort, och därmed Din vandring, säger vem du slogs mot: "Du
  anföll Bo.", "Du parerade mot Bo." eller "Du gjorde en undanmanöver för
  Bo." Mot en varelse står det som förut.

### Ormleken och korten i krönikan

- Bredvid guldhögen ligger ormleken och de dragna korten. Står du i
  skattkammaren ser du leken krympa och de dragna korten växa, i samma tre
  steg som sömnmätaren. När ormen vaknar flyger korten tillbaka in i leken.
- Rundans kort och Din vandring visar kortet du drog, och krönikan i slutet
  visar alla kort var och en drog.

### Ormkorten syns

- När du tar en näve i skattkammaren vänds kortet du drog upp på kartan:
  Ormen sover eller Ormen rör sig, med en bild av ormen, texten och
  guldet. Bara du ser det.
- I skattkammaren visar panelen Nu en sömnmätare i tre steg som andas i
  ormens takt. Ormtjusaren ser vakna-korten som är kvar som små kort.
- När ormen vaknar ser alla kortet Ormen vaknar stort mitt på kartan, och
  den som stod i kammaren ser vad den förlorade.

### Kartan som ett spelbräde

- Kartan ser ut som ett spelbräde i fackelsken. Golven är huggna
  stenplattor, murarna har krön, dörrarna är av trä och portarna är
  portvalv med torn. En stängd port har stängda dörrar med röda järnband.
- Varje rum har facklor. Rummet du står i lyses upp, grannrummen ligger i
  dunkel och rummen längre bort nästan i mörker. Kartan har en ram som ett
  gammalt spelbräde.
- Rummen har inredning efter sin sort: pelare, en brunn, bänkar, ett
  altare, en trappa, hyllor och tunnor, lyktor, stenblock, en ässja och ett
  städ, bokhyllor, eller en hög med stenar och ben. Vad varje rumstyp har
  står i `data/texter.json` under `modell.inredning`.
- Vandrarna och varelserna är miniatyrer på runda socklar, med en kant i
  vandrarens färg eller mörkröd för varelserna. En fälla som har slagit
  till syns som en modell i stället för sitt tecken i krita.
- Ett nytt rum faller ner som en bricka och läggs på plats.

### Bara kartan snett ovanifrån

- Kartan i krita rakt ovanifrån är borttagen. Kartan snett ovanifrån i
  3D är den enda.
- På en enhet utan WebGL står nu en text där kartan skulle vara: att
  kartan inte kan ritas där, och att loggen och panelen Nu fungerar som
  vanligt. Förut visades kritkartan.
- Pilarna på pekskärmens knappar är alltid sneda, som kartan.
- `web/map.js` räknar bara ut utsnittet, zoomen och teckenförklaringen.

### Spelledarläget är borttaget

- Spelledarläget finns inte längre, och inte heller spelledarbordet med
  översikten och ingreppen: pausa, tid, solen, tornets röst, fällor och
  varelser, ormen, avslöjandets nivå, liv, guld, drycker, föremål och
  invigningen.
- Spelledaren leder tornet från värdsidan och kan spela med, från tornets
  dator eller en annan. Introt på tornets dator leder alltid till
  värdsidan.
- Högst sju AI-spelare, eftersom en plats alltid finns kvar åt en
  människa. Starta nu kräver minst en människa.
- Tornet pausar fortfarande rundan när ingen människa är ansluten, och
  fortsätter när någon kommer tillbaka.
- Spelarna får inte längre veta att tornet har en spelledare.

### Färre inställningar på värdsidan

- Inför nästa vandring har nu sex inställningar: Svårighet, AI-spelare,
  Spelledarläge, Språk, Avslöjandet och Frånvarande vandrare. Fler
  inställningar finns inte längre.
- Tornets storlek, solnedgången och turtiden bestäms alltid automatiskt
  efter antalet vandrare.
- Vandrarna väljer alltid klass, och strid mellan vandrare är alltid på
  från första rundan. Portrummen är fortfarande fredade.
- Vandringen börjar av sig själv när minst två är redo, AI-spelare
  inräknade.
- Med avslöjandet på gäller alltid invigningen med ormens namn och tornet
  som förändras. Fasta svar finns inte längre; bara om datafilerna saknar
  namn räcker svaren.
- Kartan ritas alltid snett ovanifrån i 3D. En enhet utan WebGL visar
  fortfarande den platta kartan.
- Flaggorna `-klasser`, `-storlek`, `-invigning`, `-forandras` och
  `-karta` är borttagna. `-sol`, `-turtid` och `-min` finns kvar för den
  som provar.
- Sparade inställningar från förr går att läsa; det som inte finns längre
  hoppas över.

## 2026-11-21

### Ikonen i Safari

- Sidorna har nu också en liten ikon som PNG och en `favicon.ico`, så att
  Safari visar tornets ikon i flikarna. Förut fanns bara SVG:n och en stor
  PNG, och Safari visade ingen ikon.

### Tillbaka från spelledarbordet

- Spelledarbordet har en knapp Värdsidan i raden överst. Värdsidan och
  spelledarbordet byter nu sida i samma flik. Förut öppnade länkarna en
  egen flik. På Pages hamnade den utanför tornets flik, och den som kom
  till bordet från introt fick en ny flik i stället för att komma tillbaka.

### Stäng tornet på Pages

- På GitHub Pages säger texten vid Stäng tornet inte längre att man kan
  trycka Ctrl+C i terminalen, och frågan talar inte om en server. I stället
  står det vad knappen är bra för där: spelarna får beskedet att tornet
  stänger, i stället för att deras sidor börjar leta efter tornet igen.

### Sidorna heter Ormtornet från början

- Sidornas titel står nu som "Ormtornet" redan i HTML:en. Förut var den
  tom tills texterna hade laddats. När man delar länken visar
  förhandsvisningen därför Ormtornet i stället för adressen, och fliken har
  ett namn medan sidan laddar. Sedan byts titeln som förut mot den på
  tornets språk.

### Tornet håller i en dold flik

- Steg 4 i [github-pages.md](github-pages.md). På GitHub Pages körs tornet
  nu i en Web Worker (`web/torn-arbetare.js`). Webbläsare saktar ner flikar
  som inte syns, och tornets flik hamnar i bakgrunden så fort spelledaren
  går in själv. I workern ska tiden, robotarna och pulsen gå som vanligt.
- Värdsidan på Pages säger att tornets flik ska hållas öppen och att en
  dator är bäst. Kommer någon inte in? har råd för länken, för nät som
  stoppar anslutningen och för den som tappar den.
- Den som tappar anslutningen öppnar länken igen och skriver samma namn,
  och får tillbaka sin vandrare.
- Bygget för Pages sätter ett versionsnummer på skripten, stilarna och
  bilderna (`?v=` och början av commiten), och skripten laddar workern,
  `torn.wasm` och `rtc.js` med samma nummer. Pages låter webbläsaren spara
  filerna i tio minuter, och utan numret kunde en ny version blandas med
  en gammal. Sidorna själva (`.html`) kan fortfarande vara upp till tio
  minuter gamla, men då hämtar de sin egen version av resten.

### Spela från en länk

- Steg 3 i [github-pages.md](github-pages.md). När tornet körs på GitHub
  Pages visar värdsidan under Bjud in en länk med rumskod och nyckel, och
  en QR-kod för telefonerna. Den som öppnar länken ansluter till tornets
  flik över WebRTC, var hen än är. Texterna under Bjud in och i inbjudan
  talar om länken i stället för samma wifi.
- Spelaren ser "Letar efter tornet …" medan sidan söker, och ett besked om
  tornet inte går att hitta.
- Det som kommer över WebRTC räknas aldrig som spelledaren: värdsidan och
  spelledarbordet går inte att nå därifrån.
- Kvar är steg 4, det som gör det hållbart, och steg 5, ett provspel över
  internet. Det är inte provat mellan olika nät än.

### Tornet på GitHub Pages

- Steg 2 i [github-pages.md](github-pages.md). GitHub Action `Sida`
  (`.github/workflows/sida.yml`) kör testerna, bygger sidan med
  `verktyg/bygg-sida.sh` och publicerar den på GitHub Pages vid varje push
  till `main`. Det går också att starta den för hand under Actions.
- Sidan fungerar under en undermapp som `/ormtornet/`, som den hamnar på
  Pages. Det är provat lokalt.
- Repot måste vara publikt och Pages påslaget med GitHub Actions som källa
  innan något publiceras. Spelare på andra datorer kan inte ansluta än;
  det är steg 3.

### Tornet i webbläsaren

- Steg 1 i [github-pages.md](github-pages.md). Tornet kan byggas till
  WebAssembly och köras i en flik, skalet (`web/torn.html`), med
  `verktyg/bygg-sida.sh`. Spelledaren och spelarna i samma webbläsare når
  det genom `web/natet.js`. README beskriver hur man provar under Tornet i
  webbläsaren. Spelare på andra datorer kan inte ansluta än.
- Sidorna länkar nu till varandra med relativa adresser, till exempel
  `spelledare.html` och `spelledare/bord.html`, så att samma länkar
  fungerar från tornets program och från en vanlig webbserver. Tornets
  program svarar på de nya adresserna, och de gamla (`/spelledare`,
  `/spelledare/bord`) fungerar som förut.
- Ny flagga `-sida`, som skriver sidorna för GitHub Pages till en mapp.

### Three.js licens följer med kartan

- Kartan snett ovanifrån (`web/iso/iso.js`) börjar nu med three.js
  licens (MIT), som licensen kräver. Bygget i `klient/` lägger dit den.
  Steg 0 i [github-pages.md](github-pages.md).

### Provspelet dokumenterat

- [provspel-2026-11-21.md](provspel-2026-11-21.md) beskriver provspelet:
  hur det gick till, de 16 punkterna, vad som byggdes och vad som inte är
  provat. Småsakerna som inte ändrades står i
  [att-gora.md](att-gora.md).
- Robotarna som provspelar som människor finns nu i
  `verktyg/provspelare.mjs`, beskrivna i README under Provspela med
  robotar.

### Egen karta är det enda läget

- Inställningen "Vandrarna syns" och flaggan `-synlighet` är borttagna.
  Vandrarna ser alltid bara varandra när de står i samma rum, och var och
  en har sin egen karta med de rum hen själv har gått igenom. Det var
  förvalet förut.
- Portarna lottas alltid, och spåren och ljuden från grannrummen finns
  alltid.
- Under "Vad inställningarna gör" på värdsidan står en kort förklaring av
  hur vandrarna ser varandra.
- En sparad inställningsfil med `synlighet` läses som förut; valet
  ignoreras.

### Värdsidan och spelledarbordet

- Under vandringen visar Vid porten bara raden om vandringen och
  vandrarna: de lediga platserna, Starta nu och raden om tornets storlek
  är borta tills vandringen är slut.
- Raden under Vid porten står inte längre indragen eller högerställd.
  Spelarsidans stil för tecknen träffade den av misstag.
- Spelledarbordets valrutor och textrutan för tornets röst har etiketter
  för skärmläsare, på svenska och engelska.

### Tornets rum

- Rummen bredvid en port har aldrig en fälla, så att ingen förlorar liv
  redan på första steget in. Med AI-spelare kommer några procentenheter
  fler ut med skatt på alla svårigheter (se `docs/att-gora.md`).
- Två rum bredvid varandra har aldrig samma namn, och rumstyperna används
  så jämnt som möjligt. Det blir lättare att berätta för de andra var man
  är.
- Kartritaren säger bara att kartan delas när den andra får nya rum, inte
  varje runda de står tillsammans.
- Golvet "gungar till under er" heter nu "gungar till", eftersom det också
  sägs till den som är ensam.

### Bättre på telefonen

- Introt fyller hela skärmen och står i mitten, så att ett tryck var som
  helst startar stormen eller går in. På en pekskärm står det "tryck för
  att höra stormen".
- Pilarna vid dörrarna på kartan snett ovanifrån är desamma som
  knapparnas (↗ ↘ ↙ ↖) på en pekskärm, och ett tryck på en pil går åt det
  hållet.
- Klassvalets nedräkning står ovanför korten och följer med när man
  rullar. På en telefon står den fast överst, korten står två och två och
  den kursiva beskrivningen är borta, så att alla kort ryms på ungefär en
  och en halv skärm.
- Ett långt rumsnamn kortas med … i stället för att krocka med tecknen
  uppe till höger.

### Spelskärmen säger mer

- Ett steg mot en vägg avvisas direkt med väggens text, som ett rum som
  redan är genomsökt, och kostar inte rundan. I mörker, där man inte ser
  dörrarna, får man pröva sig fram som förut.
- Panelen Nu säger hur många rundor till man sitter fast, att T tar skatt i
  kammaren, och, för den som har gått ut, hur mycket guld hen bar ut medan
  de andra är kvar.
- I slutet har panelen Nu knappen Visa krönikan, så att alla kommer
  tillbaka till krönikans skärm efter Visa kartan.
- Kartritarens svaghet heter nu "Inga snabba fötter": det är bara
  undanmanövern som aldrig lyckas. Hen kan slåss som alla andra.

### Rättat efter provspelet

- Den som ansluter medan slutskärmen visas står inte längre i resultatet
  som förlorad. Hen väntar på nästa torn ("väntar på nästa torn" i listan)
  och är med när alla trycker Nytt torn.
- När en varelse träffar dig följer texten det du valde: den parerar ditt
  anfall, glider runt din parering eller hinner före din undanmanöver.
  Förut kunde det stå "Du hinner inte undan" fast du anföll.
- På en pekskärm har Viska en fast plats bland knapparna och är osynlig
  tills den går att använda, så att Ta inte byter plats med Anfall när den
  dyker upp. Redo syns bara i lobbyn och i slutet.
- Värdsidan byter svårighetens förklaring direkt när språket byts, och
  Normal förklaras med liv, fällor och ormen i stället för "som i
  datafilerna".

### Krönikan på egen skärm när vandringen är slut

- När vandringen är slut får alla en egen skärm: resultaten, epilogen,
  krönikan över hela vandringen och sin egen vandring runda för runda.
  Den som dog under vandringen ser sin skärm växa till hela omgången.
- Krönikan står inte längre i loggen. Visa kartan går till kartan med
  resultatrutan och loggen, och Nytt torn (Enter) leder till lobbyn som
  förut.

### Din vandring: en egen skärm för den som dör

- Den som dör ser nu sin egen vandring på en egen skärm: det som hände en
  runda för runda (rundans kort), hur man dog och sina siffror.
- Skärmen säger hur många som är kvar i tornet, men inte var de är eller
  vad de har hittat.
- Visa kartan går tillbaka till spelet, och Din vandring i panelen Nu går
  hit igen.

### Den döda styr inte längre varelserna

- Den första som dör tar inte längre över tornets varelser. Varelserna
  sköts alltid av tornet och vandrar efter sina vanor, som de redan gjorde
  innan någon hade dött.
- Den som har dött ser inte heller alla vandrare eller hela tornet längre.
  Spelledaren ser allt som förut.
- Tab byter bara mål när flera vandrare står i rummet.

### Piltangenterna syns vid dörrarna

- På kartan snett ovanifrån står en liten tangent med pilen vid varje
  öppen dörr och port i rummet du står i: ↑ vid den norra, → vid den
  östra och så vidare. Draget du har valt lyser i din färg.
- På en telefon visar pilknapparna sneda pilar, ↗ ↘ ↙ ↖, åt det håll
  figuren går.

### Modellernas färg och storlek i data/

- Varje klass i `klasser.json` och varelse i `monster.json` har `modell`
  med `farg` och `storlek`, för figuren på kartan snett ovanifrån. Ändra
  till exempel Blekkryparens färg eller gör Fjällvakten större, och tryck
  på Läs om innehållet på värdsidan.
- Kontrollen säger till om en färg inte är skriven som #rrggbb eller om en
  storlek inte ligger mellan 0,5 och 2.

### Kartan snett ovanifrån på telefoner

- På en telefon eller surfplatta ritar kartan högst 30 bilder i sekunden,
  vilket sparar batteri.
- Går kartan trögt slår den av sig själv av skuggorna och ritar med lägre
  upplösning. Adressen `/?kvalitet=lag` gör det direkt.
- Zoomen Nära går lite närmare, så att rummet man står i syns bättre på en
  smal skärm.

### Teckenförklaringen visar modellerna

- Med kartan snett ovanifrån visar teckenförklaringen uppe till vänster en
  liten bild av vandrarna, varelserna, kammaren och portarna, så som de
  ser ut på kartan. Fällornas märken och kritan har kvar sina tecken.

### Föremålen syns på figuren

- På kartan snett ovanifrån bär den egna figuren sina föremål: repet och
  väskan vid höften, amuletten och pärlan på bröstet, skölden på ryggen och
  kronan av rost på huvudet, och så vidare.
- Andras föremål syns inte, eftersom väskorna är hemliga. På
  spelledarbordet syns allas.

### Kartan snett ovanifrån: kompassen och pilarna vid kanten

- Nere till vänster på kartan snett ovanifrån står en kompass, ett N med en
  pil åt det håll norr ligger. Piltangenten upp går dit.
- Portar, kammaren och vandrare som ligger utanför bilden får en pil vid
  kanten med sitt tecken, som på den platta kartan.

### Väskan syns tydligare

- Föremålen står i väskan under tecknen uppe till höger, nu med ett eget
  tecken framför namnet (`tecken` i `data/foremal.json`) och lite större.
- Under vandringen står en svag prick för varje tom plats, så att man ser
  att man kan bära tre.
- Ett föremål man just har hittat lyser en stund.
- Förbannade föremål ser ut som de andra.

### Kartan snett ovanifrån är förval

- Kartan visas nu snett ovanifrån i 3D om inget annat är valt. Den platta
  kritkartan finns kvar: välj Kartan: rakt ovanifrån på värdsidan, eller
  starta med `-karta platt`. En enhet utan WebGL visar den platta kartan.
- Har du sparat inställningar efter att valet Kartan kom, med Kartan på
  rakt ovanifrån, gäller det valet tills du ändrar det.

### Kartan snett ovanifrån: avslöjandet

- Tornet snett ovanifrån förändras nu med avslöjandets nivåer, som den
  platta kartan: kritan andas på nivå 2, hjärtslagen går genom väggarna
  och revben syns genom golvet på nivå 3, och på nivå 4 är väggarna kött
  och portarna munnar som tuggar.
- När ormen vaknar lyser dess ögon över guldet.
- Det var sista steget i [isometrisk.md](isometrisk.md).

### Kartan snett ovanifrån: figurerna

- Vandrarna på kartan snett ovanifrån har nu armar och ben och bär det
  deras klass bär: hjälm och sköld, lykta, kartrulle, stav och så vidare.
- Varelserna har egna modeller: Blekkryparen, Ormynglet, Fjällvakten,
  Gallgasten och Benväktaren.
- Figurerna går genom dörren till nästa rum i stället för att hoppa, och
  den som dör lägger sig på golvet.
- Steg 3 av 4 i [isometrisk.md](isometrisk.md). Kvar är avslöjandets
  utseende.

### Kartan snett ovanifrån, första versionen

- Värdsidan har ett nytt val under Inför nästa vandring: **Kartan**, rakt
  ovanifrån i krita eller snett ovanifrån i 3D (flaggan `-karta`). Förval
  är rakt ovanifrån.
- Snett ovanifrån visar tornet med golv, låga väggar, dörrar och portar,
  och kammaren med sin ring och sitt guld. Vandrarna är små figurer i sin
  färg och varelserna står som sina tecken. Kameran följer rummet man står
  i, och zoomen fungerar som förut. Norr pekar upp åt höger.
- Spelledarbordet visar samma karta.
- En enhet som inte kan rita 3D (WebGL) visar den platta kartan.
- Det här är steg 2 av 4 i [isometrisk.md](isometrisk.md). Klassernas
  figurer, rörelsen mellan rummen och avslöjandets utseende kommer sedan.

### Grunden för kartan snett ovanifrån

- Kartan ska kunna visas snett ovanifrån i 3D, ritad med Three.js. Planen
  står i [isometrisk.md](isometrisk.md). Det här första steget lägger bara
  grunden: TypeScript och Vite i `klient/`, som bygger `web/iso/iso.js`.
  Ingenting ändras på skärmen ännu.
- Spelet byggs som förut utan Node. Den som ändrar i `klient/` behöver
  Node och kör `npm install` där en gång; sedan bygger `starta.sh` om
  kartan av sig själv.

## 2026-11-20

### Styrningen stängs när vandringen är slut

- Styrningen, som visas av sig själv första vandringen, stängs när
  vandringen tar slut, så att den inte ligger över slutkortet. ? visar den
  igen.

### Porten frågar innan du går ut tomhänt

- Den som går mot en öppen port utan att bära något guld får först en
  fråga i loggen. Går hen samma håll en gång till, samma eller nästa
  runda, går hen ut. Förut räckte ett tryck, så att man kunde råka lämna
  tornet redan i första rundan. Med guld går man ut direkt, och AI-spelarna frågas inte.

### Värdsidan: klasserna, AI-spelarna och spelledarläget

- Valet för klasserna heter nu "på, vandrarna väljer klass", och
  förklaringen säger att klassvalet har ett kort för varje klass, inte två.
- När spelledarläget slås på eller av får AI-spelare valet 8, eller
  förlorar det, direkt. Förut syntes det först när man klickade någon
  annanstans. Justerar tornet ett val visas det genast.
- Förklaringen bredvid Spelledarläge följer valet. Med "av" står det att du
  spelar själv.

### Ladda om sidan i lobbyn

- Om alla andra var redo och den enda som inte var redo laddade om sidan,
  började vandringen genast utan hen. Med klasserna av kom hen inte in
  igen förrän solen hade gått ner. Nu väntar lobbyn tio sekunder efter
  att någon har gått innan vandringen kan starta. Starta nu på värdsidan
  startar som förut direkt.

### Ladda om sidan under klassvalet

- Den som laddade om sidan eller tappade anslutningen under klassvalet
  åkte ut, och hade de andra redan valt började vandringen direkt utan
  hen. Nu behåller man sin plats: klassvalet väntar tills hen är tillbaka
  eller tiden är slut.
- Är hen fortfarande borta när tiden är slut får hen ett ledigt kort och
  följer med som frånvarande, som om hen hade gått mitt i vandringen. Med
  samma namn kommer hen in igen, och annars tar AI:n över efter så många
  rundor som "Frånvarande vandrare" säger.

### Ormtornet.app i ~/Projekt/Lokalt

- `./skapa-app.sh` skapar nu appen i `~/Projekt/Lokalt` i stället för
  `~/Projekt`. En annan mapp kan fortfarande anges som argument.
- En app som redan finns kan flyttas fritt, eftersom den bygger spelet från
  repots mapp. Den behöver bara skapas om när repot flyttas.

## 2026-11-19

### Förklaringar till inställningarna

- Under Inför nästa vandring på värdsidan finns "Vad inställningarna gör",
  hopfälld, med några meningar om varje inställning, på svenska och
  engelska. Förklaringarna avslöjar inte tornets hemlighet.
- Efter en viskning syns det direkt om vandraren blev invigd, utan att
  vänta på rundan.

### Spelledarbordet visar namnet, varslen och portarna

- Översikten visar ormens namn, om tornets nåd har kommit och hur många
  portar som är öppna.
- Tabellen visar under Tillstånd vilka delar av namnet var och en har.
- Under tabellen står de senaste varslen: vem som fick dem, vad som väntar
  och om någon har gått in i rummet.
- Förslaget och läget står i `docs/avslojandet.md`.

### Tornet förändras

- Ny inställning under Fler inställningar: Tornet förändras, på som
  förval. Den kräver avslöjandet. Flaggan heter `-forandras`.
- Från nivå 3 blir fällorna som hör till ormens kropp vanligare i nya rum:
  det rasande golvet, det roterande rummet, ormgropen och det ömsade
  skinnet (`forandring.fallvikter` i `avslojande.json`).
- På nivå 4 sluter sig portarna en i taget, var tredje runda
  (`port_stangs_var`). Först den som ligger längst från vandrarna, och
  rundan innan får alla som är inne en varning. Den sista porten stängs
  aldrig. Kartan visar en stängd port som ett rött block, och ingen kommer
  ut genom den. Datorspelarna söker sig till en öppen.
- Med 600 simulerade vandringar för 2, 4 och 7 AI-spelare ändrades andelen
  som blev kvar i mörkret med högst två procentenheter (2–6 %).

### Varslen

- Från avslöjandets nivå 2 är varannan ledtråd ett varsel: tornet förråder
  ett rum som ingen har varit i, två eller tre rum bort, till exempel "Det
  gnisslar av järn två rum bort åt nordöst".
- Vad som finns i rummet bestäms när varslet ges, med samma chanser som
  annars, och blir så när någon går in: en fälla, en varelse eller
  ingenting. Varslet gäller aldrig en port eller kammaren.
- Texterna står under `varsel` i `avslojande.json`, en för varje fälla och
  varelse.

### Byte vid dödande och tornets nåd (spoilervarning)

- Med ormens namn får den som dödar en annan vandrare de delar av namnet
  som offret hade och dödaren saknade.
- Tornets nåd: när alla människor som är kvar inne är invigda tas hälften
  av vakna-korten i leken bort (`nad_andel`, minst ett), en gång per
  vandring. Nåden gäller också när leken byggs om efter att ormen har
  vaknat.
- Krönikan berättar om bytet och nåden, och nämner de invigda som kom ut.

### Invigningen med ormens namn (spoilervarning)

- Ny inställning under Fler inställningar på värdsidan: Invigningen, med
  valen ormens namn (förval) och fasta svar. Den kräver avslöjandet.
- Med ormens namn har varje vandring ett eget namn i tre delar. Vid nivå 2,
  3 och 4 får varje människa i tornet en del, så att gruppen tillsammans
  vet så mycket som möjligt. Delarna syns i tecknen, till exempel `· RU ·`.
- För att bli invigd ska viskningen innehålla namnet och visa att man
  förstår, till exempel "Karumeth, tornet är ormen" eller "tornet är
  Karumeth". De fasta svaren räcker inte ensamma.
- Epilogen säger namnet.
- Fasta svar fungerar som förut. Flaggan heter `-invigning`.

### Spelledarbordet

- Översikten och ingreppen i spelledarläget har flyttat från värdsidan
  till en egen sida, spelledarbordet, på `/spelledare/bord`. Värdsidan
  har kvar inbjudan, porten och inställningarna.
- Värdsidans rad har knappen Spelledarbordet när spelledarläget är på.
  Pausa och Fortsätt finns nu på spelledarbordets rad.
- Spelledarbordet har en egen rad med fasen, vilka som är här, och Starta
  nu före vandringen eller Pausa och Fortsätt under den. Kartan är större
  och loggen står utfälld.
- Värdsidans rad säger när rundan är pausad, även utan spelledarläge.
- Med spelledarläget av säger spelledarbordet det och har knappen Slå på
  spelledarläget, som fungerar mellan vandringarna.
- I lobbyn och klassvalet visar spelledarbordet vilka som står vid porten,
  om de är redo och vilken klass de har valt.
- När vandringen är slut visar spelledarbordet hur det gick för var och en
  och krönikan över vandringen.
- Introt på tornets dator leder till spelledarbordet när spelledarläget är
  på, och annars till värdsidan som förut.
- Värdsidans avsnitt byter inte längre plats efter fasen. De står alltid i
  ordningen Bjud in, Vid porten och Inför nästa vandring. Under och efter
  vandringen krymper Bjud in till adressen, så att resten kommer upp.

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
