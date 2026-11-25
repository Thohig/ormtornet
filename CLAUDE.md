# Ormtornet

En dungeon crawler för 2–8 spelare i webbläsaren, främst från en länk på
GitHub Pages (<https://thohig.github.io/ormtornet/>), eller nedladdad på
en Mac i ett lokalt nätverk när länken inte räcker.
README beskriver spelet, flaggorna och filerna.

## Arbetssätt

- Skriv på svenska, både i svar och i spelets texter. Spelet finns också
  på engelska: en ny eller ändrad text ska in både i `data/` och i
  `data/en/` (se [docs/engelska.md](docs/engelska.md)).
- Bygg i steg. Gör en commit efter varje steg och stanna sedan så att
  användaren kan testa innan nästa steg.
- Allt innehåll (rum, fällor, varelser, texter, ledtrådar) ligger i `data/`
  och ska kunna ändras utan att röra koden.
- Undvik långa rader med upprepade tecken i kod och utdata. Bygg ramar och
  mätare med fast bredd i kod.
- Tornets hemlighet får aldrig sägas rakt ut i spelet förrän i slutet.

## Ord

Så här skiljer vi på sakerna när användaren ger instruktioner:

- **spelledaren**: personen eller rollen som leder tornet.
- **spelledaren** kan också spela själv, från tornets dator eller en
  annan. Spelledarläget, spelledarbordet, översikten och ingreppen togs
  bort 2026-11-22.
- **värdsidan**: sidan på `/spelledare` på tornets dator (`host.go`,
  `web/spelledare.html`). Överst står **raden** med fasen och knappen
  som behövs just nu. Avsnitten heter Bjud in, Vid porten, Inför nästa
  vandring, Fröet, Kommer någon inte in?, Läs om innehållet och Stäng tornet, och
  står alltid i den ordningen. På en bred skärm står de i två kolumner,
  Bjud in och Vid porten till vänster och resten till höger. Under Bjud
  in står adressen (på Pages länken) inte utskriven; knapparna kopierar
  den, och Visa adressen visar den. Ändras AI-spelarna, svårigheten
  eller fröet i slutet stängs slutskärmen och alla kommer till lobbyn med
  ett nytt torn (`towerChanged` i `seed.go`). **Gå in själv** byter till spelarsidan
  i samma flik och heter Till spelet när spelledaren har en vandrare under
  vandringen. På Pages visas den i skalet och har länken och en varning
  om att hålla tornets flik öppen. QR-koden togs bort 2026-11-22.
- **skalet**: tornets flik när spelet kommer från GitHub Pages
  (`web/torn.html`, `web/torn.js`). Det kör tornet i en Web Worker
  (`web/torn-arbetare.js`, `wasm.go`) och visar värdsidan eller
  spelarsidan i en ram.
- **länken**: adressen med rumskod och nyckel
  (`#torn=<kod>&nyckel=<nyckel>`) som spelarna öppnar för att ansluta till
  skalet över WebRTC (`web/natet.js`, `web/rtc/rtc.js` från
  `klient/src/rtc.ts`).
- **spelarsidan**: sidan på `/` som spelarna öppnar, med vyerna introt,
  klassvalet, spelet och krönikan. Namnrutan togs bort 2026-11-22; tornet
  ger varje vandrare ett namn (`names.go`, `data/namn.json`). Spelet har faserna lobbyn,
  klassvalet, vandringen och slutet (`lobby`, `klassval`, `spel`, `slut`). README beskriver under "Vad en spelare ser"
  vad en spelare möter.
- **spelskärmen**: vyn spelet på spelarsidan. Delarna heter:
  - **översta raden**: solmätaren, **prickarna** (vandrarna i sin färg,
    fyllda när de har valt drag eller är redo, med namn i lobbyn och
    slutet), tiden, knappen Redo (Nytt torn i slutet) och, för
    spelledaren, **märket**: ordet spelledare, som i lobbyn och slutet är
    knappen Värdsidan och under vandringen öppnar värdsidan i en ny flik
    (bara ordet på Pages) (`renderHost` i `web/app.js`, `vard` i
    tillståndet för den som är ansluten från tornets dator),
  - **tecknen**: statusen uppe till höger på kartan (klass, liv, skatt,
    drycker, tillstånd) och föremålen under dem, med förklaringar,
  - **frågetecknet**: styrningen, som ligger bakom ? på kartan,
  - **spelkortet**: rundans kort, som vänds upp till höger på kartan
    med bild, rubrik och text för rundans viktigaste händelse, rundans
    rader under (`summary.go`, `sam_` i `data/texter.json`). På datorn
    står rundans två nästa kort lika stora till vänster om det, i
    **raden** (`renderRow` i `web/app.js`), och resten av rundans kort,
    på telefonen alla utom det första, som **minikort** (`spelkort.go`,
    `showRoundCard` i `web/app.js`, bilderna i `web/kort.js`, ordningen
    under `spelkort` i `data/texter.json`). **Ormkortet** är spelkortet
    när du tar en näve,
  - **det stora kortet** mitt på kartan (`#storkort`): **uppvaknandet**,
    kortet Ormen vaknar som alla ser (`showWakeCard` i `web/app.js`), en
    port som stängs, solen och att komma ut,
  - **minikartan**: allt man har upptäckt, platt och vridet som kartan,
    uppe till vänster på kartan på datorn och i kompassens ställe på
    telefonen (`web/minikarta.js`, docs/minikartan.md),
  - **sömnmätaren**: hur djupt ormen sover, i panelen Nu i skattkammaren,
  - **panelen Nu**: draget, striden och vad som väntar, som en rad under
    kartan (`#nowBlock`). Den syns bara när den har något att säga. Panelen Vandrare togs bort 2026-11-22,
  - **loggen**: raderna, på en bred skärm som text utan ruta nere till
    vänster **på kartan** (`placeSide` i `web/app.js`), med den senaste
    rundan nederst och den förra svag, och hela loggen i **lådan** (knappen
    Loggen, L); på telefonen under kartan. Grupperade per runda under
    rubrikerna Före vandringen (lobbyn), Vandringen börjar (klassvalet och
    det första rumsblocket) och Runda 1, 2 … Bara den
    senaste rundan står i bild; de äldre ligger under strecket "tidigare".
    Det första rumsblocket har porten man kliver in genom, solen och
    portarna man hör de andra gå in genom (`atTheGate` i `here.go`), och
    står sist under Vandringen börjar.
    Överst i rundan står **rumsblocket**, som löpande text utan ruta och
    utan namn: beskrivningen när rummet är nytt, det man
    ser (beskrivet, inte förklarat) och sist det man hörde (`here.go`,
    `har_` i `data/texter.json`, `spar_efter` i `data/fallor.json`).
    Ett känt, tomt rum ger bara raden "Du känner igen rummet" med den
    korta beskrivningen (`har_igen`), och en strid som pågår hörs bara
    första rundan (`hearFight` i `sounds.go`).
    Teckenförklaringen togs bort 2026-11-22.
- **namnen**: rum, fällor och varelser har inga namn för spelarna, bara
  beskrivningar. `namn` i `data/` syns bara på värdsidan; där en text
  behöver nämna något står den **korta beskrivningen** (`kort`, och
  `kort_bestamd` för varelserna, och `kort_igen` när striden har pågått
  en runda). Ormen, skattkammaren och föremålen har namn.
- **frånvarande vandrare**: inställningen Frånvarande vandrare säger efter
  hur många rundor AI:n tar över en vandrare vars människa har lämnat
  eller står still, alltså inte väljer något drag (`takeover.go`). Ett eget
  drag, eller att komma tillbaka, ger vandraren tillbaka.
- **tornets dator**: datorn som kör servern.
- **fröet**: talet som avgör tornet; samma frö och samma drag ger samma
  vandring (`seed.go`). Avsnittet Fröet på värdsidan visar den här
  vandringens frö och nästa vandrings, som går att skriva in eller slumpa.
- **språket**: inställningen Språk på värdsidan (`-sprak`), svenska eller
  engelska, för hela tornet (`lang.go`). Översättningarna ligger i
  `data/en/` med bara texter.
- **sidornas texter**: spelarsidans och värdsidans egna texter i
  `data/spelarsidan.json` och `data/vardsidan.json`. HTML:en har nycklar
  (`data-t`), skripten använder `t()` (`web/texter.js`), och Go-koden
  `HostText`, `hostErr` och `gmNote` med nycklar som börjar med `server_`.
- **synligheten**: att vandrarna bara ser varandra när de står i samma rum
  och var och en har sin egen karta med de rum hen själv har gått igenom.
  Det är spelets enda läge; det går inte att välja något annat. Spåren och
  ljuden från grannrummen hör till den (`sounds.go`). Se
  [docs/synlighet.md](docs/synlighet.md).
- **stormen**: introts moln, regn, blixtar, åska och musik.
- **krönikan** (skärmen): vyn `chronscreen` på spelarsidan (`renderChron` i
  `web/app.js`). Den som dör ser **Din vandring**, rundans kort runda för
  runda (`vandring` i tillståndet, `journeyOf` i `chronicle.go`,
  `addJourney` i `summary.go`), med vandringens sista kort överst, vad
  skadan kom av, tysta rundor
  ihopslagna och märket "tornet styrde". I slutet ser alla
  **placeringen**, en rad för var och en med **titeln** (den girigaste,
  jägaren …) och den största stunden, där ett klick öppnar vandrarens
  **guldbok**: hur hen fick sitt guld som berättelse, och raderna **runda
  för runda** bakom en knapp (`guldbok.go`, docs/guldboken.md). Sedan
  **berättelsen** (inträdet, höjdpunkterna med var sitt spelkort och
  slutet, `berattelse.go`, docs/berattelsen.md), epilogen och bakom
  knappen Visa din vandring sin vandring. Knappen Se vandringen igen
  visar **reprisen**: allas vägar runda för runda på tornet som det blev
  (`replay.go`, `web/repris.js`, docs/repris.md). Hela krönikan och
  tornets dunk togs bort från slutskärmen 2026-11-24.
  Spelledaren har där också Till värdsidan bredvid Nytt torn.
- **tangenterna vid dörrarna**: piltangenterna som står vid dörrarna i ens
  rum på kartan snett ovanifrån (`ritaTangenter` i `klient/src/karta.ts`).
- **minikort**: ett litet spelkort framför en rad i loggen, på
  spelkortet och i krönikan (`Kort.mini` i `web/kort.js`). Ett klick
  visar det stort i **kortvisningen**, ett lager över hela sidan
  (`showCardPopup` i `web/app.js`). I Din vandring står rundans kort som
  små spelkort med bild och rubrik.
- **boken**: Vandrarens bok, hjälpsidan i ett lager över spelarsidan
  med flikar (`web/bok.js`, `bok.go`, `GET /bok`, `bok_*` i
  `data/spelarsidan.json`, `bok` i `data/monster.json` och
  `data/fallor.json`). Den fylls i med det man möter (`learn` i `bok.go`,
  sparat i webbläsaren) och får aldrig visa namnen eller tornets hemlighet
  ([docs/vandrarens-bok.md](docs/vandrarens-bok.md)).
- **väskan**: föremålen under tecknen, med en plats för varje föremål man
  kan bära (`tecken` i `data/foremal.json`).
- **klassvalet**: skärmen efter lobbyn med ett **klasskort** för varje
  klass och en tidsgräns (`klassval_sekunder` i `data/klasser.json`).
  Klasskorten är spelkort med bild och baksida, som vänds upp ett i taget
  (`renderChoose` och `buildClassCard` i `web/app.js`). Under vandringen
  står klassen bland tecknen som ett litet klasskort; ett klick visar det
  stort (`showClassCard`).
- **klasser**: det vandrarna väljer i klassvalet (`classes.go`,
  `data/klasser.json`), till exempel Kartritaren. Elva sedan 2026-11-25,
  med Gravplundraren, Lyssnaren och Ormens präst (docs/nya-klasser.md).
- **föremål**: det vandrarna hittar och bär, högst tre (`items.go`,
  `data/foremal.json`). Alla verkar av sig själva sedan 2026-11-22; inget
  används med 1–3. Några är förbannade. Två som hör ihop **fogas ihop**
  till ett nytt, en **kombination** (`av` i `foremal.json`,
  docs/kombinationer.md). Några har **två sidor** och **vänder** efter en
  nivå eller några rundor (`vander_niva` och `vander_rundor`,
  docs/tva-sidor.md).
- **viska** och **invigd**: den som vet tornets hemlighet kan viska den
  med V och blir då invigd (`secret.go`, `svar` och `invigd` i
  `data/avslojande.json`).
- **invigningen**: det som gör någon invigd, **ormens namn**. Varje
  vandring har ett namn i tre **delar** som var och en hittar själv: i
  **ristningen** i ett rum, hos en varelse och i kammaren
  (`namndelar.go`). Texten säger var på ormen delen sitter (huvudet,
  buken, svansen), inte vilken i ordningen den är; till
  det hör **byte vid dödande** och **den förste invigde**, som får mer guld
  (`secret.go`). Tornets nåd togs bort 2026-11-23.
- **varsel**: en ledtråd om ett rum som ingen har varit i, där innehållet
  bestäms när varslet ges (`omens.go`, `varsel` i `avslojande.json`).
- **tornet förändras**: att fällorna ändras från nivå 3 och portarna
  sluter sig på nivå 4 (`change.go`, `forandring` i `avslojande.json`).
- **avslöjandet**: att tornet visar sin hemlighet i fyra nivåer som följer
  solen (`reveal.go`, `data/avslojande.json`). Det är alltid på; det går
  inte att stänga av sedan 2026-11-23.

- **striden**: det enda stridsvalet är anfall (A). **Kraft** är hur lätt
  man träffar och **Försvar** hur svår man är att träffa (klasser,
  varelser och föremål); varelsens **lynne** är chansen att den anfaller
  när man gör något annat. Att vänta i strid är att hålla **garden**
  (`monster.go`, `traff` i `data/monster.json`, docs/bara-anfall.md).
  Parera och undanmanöver togs bort 2026-11-24.
- **kammaren**: ligger nästan var som helst, men ingen har den mycket
  närmare än de andra (`plats` i `data/ormen.json`, `chamberSpots` i
  `game.go`, docs/ormens-plats.md), och flyttar när ormen vaknar; den
  gamla blir ett rum med **löst guld** (`moveChamber` i `serpent.go`,
  `spillGold` i `guld.go`). **Andetagen** hörs genom väggen bredvid den
  och för alla från nivå 2, och alla hör åt vilket håll någon tar en
  näve.
- **händelserna**: tornets sexton slumpade händelser, var och en med ett
  eget spelkort (`tornhandelser.go`, `data/handelser.json`,
  docs/handelser.md).
- **placeringen**: allas plats 1 till N i slutet, och **vinken**: det man
  får veta om hur mycket den man möter bär, och om det är mer eller mindre
  än man själv bär (`results` i `game.go`, `packHint` i `pvp.go`).
- **den individuella tävlingen**: inget i spelet belönar samarbete sedan
  2026-11-24 (docs/individuell-tavling.md). Föreslå inte samarbete utan
  att användaren ber om det.

Spelarna ser ordet "spelledare" i spelets texter, aldrig "värdsidan".

## Ändringar

[docs/andringar.md](docs/andringar.md) säger vad som är nytt, nyast först.
Skriv dit när ett steg är klart.

## Planerat

- [docs/att-gora.md](docs/att-gora.md): det som ska provas och justeras,
  det som återstår av förslagen och nya idéer. Inget där är bestämt; ta
  fram ett förslag och fråga innan du bygger något därifrån.

- [docs/spelledaren.md](docs/spelledaren.md): värdsidan. Del 1 är byggd;
  översikten och ingreppen i del 2–3 byggdes och togs bort 2026-11-22.
  Avsnittet Läget säger vad som är byggt, vad som är på gång och vad som
  inte är planerat. Fråga vilka delar som ska byggas innan du börjar på
  något som inte står under På gång.
- [docs/github-pages.md](docs/github-pages.md): att spela från en länk på
  GitHub Pages, med tornet i spelledarens webbläsare och WebRTC. Steg 0–4
  är byggda och publiceras vid varje push till `main`. Kvar är steg 5, ett
  provspel över internet, och kanske TURN efter det. Fråga innan du börjar
  på nästa steg.
- [docs/fler-spelare.md](docs/fler-spelare.md): klasser, tornets storlek
  och upp till 8 spelare. Allt är byggt.
- [docs/foremal-och-hemligheten.md](docs/foremal-och-hemligheten.md):
  föremålen och viskningen till tornet. Allt är byggt; avsnittet Inte
  byggt säger vad som återstår.
- [docs/engelska.md](docs/engelska.md): engelska som språk, valt på
  värdsidan. Allt är byggt. Nya texter ska in på båda språken (`data/` och
  `data/en/`); ett test kräver att allt är översatt.
- [docs/avslojandet.md](docs/avslojandet.md): invigningen med ormens namn,
  byte vid dödande, tornets nåd (borttagen 2026-11-23), varslen och tornet som förändras. Allt är
  byggt.
- [docs/synlighet.md](docs/synlighet.md): att vandrarna bara ser varandra
  när de möts, och bara sin egen karta. Allt är byggt, och egen karta är
  sedan 2026-11-21 det enda läget.
- [docs/provspel-2026-11-21.md](docs/provspel-2026-11-21.md): provspelet
  2026-11-21, vad som hittades och byggdes, och vad som inte är provat.
  Robotarna som provspelar finns i `verktyg/provspelare.mjs` (README,
  Provspela med robotar), med reglerna för trasiga texter i
  `verktyg/textkontroll.mjs`. Webbläsarprovet `verktyg/webbprov.mjs`
  spelar en kort vandring på dator och telefon i Chromium (README,
  Webbläsarprovet). Kör båda efter en ändring i texterna eller på
  spelarsidan.
- [docs/provspel-2026-11-22.md](docs/provspel-2026-11-22.md): provspelet
  2026-11-22, felen som rättades och åtta saker att fråga om (mörkret som
  visar dörrarna, Kartritaren, AI:n och kammaren med flera).
- [docs/provspel-2026-11-24.md](docs/provspel-2026-11-24.md): provspelet
  2026-11-24 med robotarna, webbläsarprovet och två egna vandringar, och
  de åtta felen och förbättringarna som byggdes efter det.
- [docs/provspel-2026-11-25.md](docs/provspel-2026-11-25.md): provspelet
  2026-11-25 med robotarna, webbläsarprovet och simuleringen, de fyra
  felen som rättades och de sju förslagen, som byggdes samma dag.
- [docs/provspel-2026-11-25-2.md](docs/provspel-2026-11-25-2.md): det
  andra provspelet 2026-11-25, efter kombinationerna, föremålen med två
  sidor och de nya klasserna: tre fel som rättades och sex förslag, av
  vilka fem byggdes samma dag och ett ströks.
- [docs/ormkorten.md](docs/ormkorten.md): att ormkorten syns, med kortet
  som vänds upp, sömnmätaren, uppvaknandet, leken bredvid guldhögen och
  korten i krönikan. Allt är byggt.
- [docs/isometrisk.md](docs/isometrisk.md): kartan snett ovanifrån, den
  enda kartan sedan 2026-11-22, med Three.js och 3D-modeller, i
  TypeScript i `klient/` och byggd med Vite till `web/iso/iso.js` (och
  WebRTC till `web/rtc/rtc.js`). Bygg om och checka in skripten efter en
  ändring i `klient/src/`. Avsnittet Stegen säger vad som är byggt, och
  Spelbrädet hur kartan fick sten, facklor, inredning, miniatyrer och
  brickor 2026-11-22.
- [docs/nytt-granssnitt.md](docs/nytt-granssnitt.md): det nya
  gränssnittet på dator och telefon 2026-11-22: lobbyn utan listan över
  vandrare, rumsblocket i loggen i stället för teckenförklaringen, bara
  den senaste rundan i bild, raden Nu under kartan och en kolumn på
  telefonen. Allt är byggt.
- [docs/passiva-foremal.md](docs/passiva-foremal.md): alla föremål verkar
  av sig själva och 1–3 är borta. Allt är byggt.
- [docs/samlat-forslag.md](docs/samlat-forslag.md): rummet som ett
  stycke i krönikan (det du ser och hör) och färre kommandon (S tar skatt,
  knapparna följer läget, tornet väljer motståndare, Z för kartan). Allt
  är byggt; detaljerna står i kronikans-stycken.md och
  farre-kommandon.md.
- [docs/vardsida-tva-kolumner.md](docs/vardsida-tva-kolumner.md):
  värdsidan i två kolumner på en bred skärm. Allt är byggt.
- [docs/frot-och-bredden.md](docs/frot-och-bredden.md): avsnittet Fröet
  på värdsidan, spelskärmen över hela bredden med loggen på kartan och
  lådan med hela loggen, och märket som öppnar värdsidan i en ny flik.
  Allt är byggt.
- [docs/spelkort.md](docs/spelkort.md): spelkort för alla situationer,
  med spelkortet som ersätter rundans kort och ormkortet, och efteråt
  vandringens största stunder, klasskorten och minikorten, och raden med
  rundans kort på datorn (2026-11-24). Allt är byggt; bilderna ritas i
  `web/kort.js` och syns i `verktyg/kortbilder.html`.
- [docs/balans-2026-11-22.md](docs/balans-2026-11-22.md): simuleringen av
  sex strategier, och hur klasserna, den grönaktiga skepnaden och
  AI-spelarna justerades efter den. Allt är byggt.
- [docs/den-dodes-guld.md](docs/den-dodes-guld.md): guldet som ligger
  kvar där någon har dött, skriket som säger var, rummet som minns,
  Kartritaren och Spårfinnaren, och AI:n som går dit. Allt är byggt.
- [docs/mobilvyn.md](docs/mobilvyn.md): telefonen med kartan över hela
  höjden, en tunnare logg, tecknen på en rad och tryck på dörrarna i
  stället för pilknapparna, och nivån Rummet närmast. Allt är byggt.
- [docs/namnet-som-tavling.md](docs/namnet-som-tavling.md): ormens namn
  som en tävling, där var och en hittar namnet själv i tornet, och
  avslöjandet som alltid är på. Allt är byggt; avsnittet Att prova säger
  vad som behöver provas i spel.
- [docs/bara-anfall.md](docs/bara-anfall.md): striden med anfall som
  enda val, Kraft, Försvar och varelsernas lynne, föremålen som gör en
  starkare och balansen efter simuleringen. Allt är byggt.
- [docs/individuell-tavling.md](docs/individuell-tavling.md): att renodla
  den individuella tävlingen. Steg 1–2 är byggda (Kartritaren ger inte
  bort sin karta, timglaset och viskningen i kammaren gäller bara en
  själv, introt uppmanar inte till samarbete, placering för alla och en
  vink som jämför med det du bär). Steg 3, ställningen för kvällen,
  valdes bort.
- [docs/ormens-plats.md](docs/ormens-plats.md): kammaren på en fri men
  ganska rättvis plats, näven som hörs, andetagen och AI-spelarna som
  följer dem. Allt är byggt.
- [docs/handelser.md](docs/handelser.md): breven som pekar mot något,
  kammaren nära mitten (sedan ersatt av den fria platsen i
  ormens-plats.md), ormen som flyttar sig när den vaknar och tornets
  sexton slumpade händelser med var sitt spelkort (de nio sista
  2026-11-24). Allt är byggt; avsnittet Inte byggt säger vad AI:n inte
  känner till.
- [docs/berattelsen.md](docs/berattelsen.md): slutkrönikan som en
  gemensam berättelse med höjdpunkter i stället för drag för drag, och
  raden om var och en med en titel. Allt är byggt.
- [docs/guldboken.md](docs/guldboken.md): att välja en vandrare på
  slutskärmen och läsa hur hen fick sitt guld, och den enklare
  slutskärmen. Allt är byggt; avsnittet Att prova säger vad som behöver
  provas i spel.
- [docs/minikartan.md](docs/minikartan.md): minikartan med det man har
  upptäckt, en fast ruta på datorn och kompassens ställe på telefonen.
  Allt är byggt; avsnittet Att prova säger vad som behöver provas i spel.
- [docs/repris.md](docs/repris.md): reprisen på slutskärmen, där allas
  vägar växer fram runda för runda. Allt är byggt; avsnittet Att prova
  säger vad som behöver provas i spel.
- [docs/vandrarens-bok.md](docs/vandrarens-bok.md): hjälpsidan
  Vandrarens bok, som fylls i med det man möter. Steg 1–2 är byggda;
  fråga innan steg 3.
- [docs/ideer.md](docs/ideer.md): idéer för framtida utbyggnad, spånade
  2026-11-25. Inget där är bestämt; ta fram ett förslag och fråga innan
  du bygger något därifrån.
- [docs/kombinationer.md](docs/kombinationer.md): två föremål som fogas
  ihop till ett nytt (idé 15). Allt är byggt; avsnittet Att prova säger
  vad som behöver provas i spel.
- [docs/tva-sidor.md](docs/tva-sidor.md): föremål med två sidor som
  vänder (idé 14). Allt är byggt; avsnittet Att prova säger vad som
  behöver provas i spel.
- [docs/nya-klasser.md](docs/nya-klasser.md): Gravplundraren,
  Lyssnaren och Ormens präst (idé 13). Allt är byggt; avsnittet Att prova
  säger vad som behöver provas i spel.
- [docs/slumpade-namn.md](docs/slumpade-namn.md): tornet ger namnen och
  namnrutan är borta. Steg 1–2 är byggda; Nytt namn i lobbyn valdes bort.
