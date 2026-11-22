# Ormtornet

En dungeon crawler för 2–8 spelare i webbläsaren, på ett lokalt nätverk
eller från en länk på GitHub Pages (<https://thohig.github.io/ormtornet/>).
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
  den, och Visa adressen visar den. **Gå in själv** byter till spelarsidan
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
    rader under (`summary.go`, `sam_` i `data/texter.json`) och resten
    av rundans kort som **minikort** (`spelkort.go`, `showRoundCard` i
    `web/app.js`, bilderna i `web/kort.js`, ordningen under `spelkort` i
    `data/texter.json`). **Ormkortet** är spelkortet när du tar en näve,
  - **det stora kortet** mitt på kartan (`#storkort`): **uppvaknandet**,
    kortet Ormen vaknar som alla ser (`showWakeCard` i `web/app.js`), en
    port som stängs, solen och att komma ut,
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
    Ett känt, tomt rum ger inget block, och en strid som pågår hörs bara
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
  ihopslagna och märket "tornet styrde". I slutet ser alla resultaten,
  **vandringens största stunder** (stora spelkort, `stunder.go`),
  epilogen, hela krönikan och sin vandring.
  Spelledaren har där också Till värdsidan bredvid Nytt torn.
- **tangenterna vid dörrarna**: piltangenterna som står vid dörrarna i ens
  rum på kartan snett ovanifrån (`ritaTangenter` i `klient/src/karta.ts`).
- **väskan**: föremålen under tecknen, med en plats för varje föremål man
  kan bära (`tecken` i `data/foremal.json`).
- **klassvalet**: skärmen efter lobbyn med ett spelkort för varje klass
  och en tidsgräns (`klassval_sekunder` i `data/klasser.json`).
- **klasser**: det vandrarna väljer i klassvalet (`classes.go`,
  `data/klasser.json`), till exempel Kartritaren.
- **föremål**: det vandrarna hittar och bär, högst tre (`items.go`,
  `data/foremal.json`). Alla verkar av sig själva sedan 2026-11-22; inget
  används med 1–3. Några är förbannade.
- **viska** och **invigd**: den som vet tornets hemlighet kan viska den
  med V och blir då invigd (`secret.go`, `svar` och `invigd` i
  `data/avslojande.json`). Bara med avslöjandet på.
- **invigningen**: det som gör någon invigd, **ormens namn**. Varje
  vandring har ett namn i tre **delar** som delas ut till vandrarna; till
  det hör **byte vid dödande** och **tornets nåd** (`secret.go`). Gäller
  alltid när avslöjandet är på; det finns ingen inställning.
- **varsel**: en ledtråd om ett rum som ingen har varit i, där innehållet
  bestäms när varslet ges (`omens.go`, `varsel` i `avslojande.json`).
- **tornet förändras**: att fällorna ändras från nivå 3 och portarna
  sluter sig på nivå 4 (`change.go`, `forandring` i `avslojande.json`).
  Gäller alltid när avslöjandet är på; det finns ingen inställning.

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
  byte vid dödande, tornets nåd, varslen och tornet som förändras. Allt är
  byggt.
- [docs/synlighet.md](docs/synlighet.md): att vandrarna bara ser varandra
  när de möts, och bara sin egen karta. Allt är byggt, och egen karta är
  sedan 2026-11-21 det enda läget.
- [docs/provspel-2026-11-21.md](docs/provspel-2026-11-21.md): provspelet
  2026-11-21, vad som hittades och byggdes, och vad som inte är provat.
  Robotarna som provspelar finns i `verktyg/provspelare.mjs` (README,
  Provspela med robotar).
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
  med spelkortet som ersätter rundans kort och ormkortet. Allt är byggt;
  bilderna ritas i `web/kort.js` och syns i `verktyg/kortbilder.html`.
- [docs/slumpade-namn.md](docs/slumpade-namn.md): tornet ger namnen och
  namnrutan är borta. Steg 1–2 är byggda; Nytt namn i lobbyn valdes bort.
