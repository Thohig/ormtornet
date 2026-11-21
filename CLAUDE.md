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
- **spelledarläget**: inställningen där spelledaren inte spelar själv, ser
  allt och kan göra ingrepp (`gm.go`, `ingrepp.go`).
- **värdsidan**: sidan på `/spelledare` på tornets dator (`host.go`,
  `web/spelledare.html`). Överst står **raden** med fasen och knappen
  som behövs just nu. Avsnitten heter Bjud in, Vid porten, Inför nästa
  vandring, Kommer någon inte in?, Läs om innehållet och Stäng tornet, och
  står alltid i den ordningen. På Pages visas den i skalet och har länken,
  QR-koden och en varning om att hålla tornets flik öppen.
- **skalet**: tornets flik när spelet kommer från GitHub Pages
  (`web/torn.html`, `web/torn.js`). Det kör tornet i en Web Worker
  (`web/torn-arbetare.js`, `wasm.go`) och visar värdsidan, spelledarbordet
  eller spelarsidan i en ram.
- **länken**: adressen med rumskod och nyckel
  (`#torn=<kod>&nyckel=<nyckel>`) som spelarna öppnar för att ansluta till
  skalet över WebRTC (`web/natet.js`, `web/rtc/rtc.js` från
  `klient/src/rtc.ts`).
- **spelledarbordet**: sidan på `/spelledare/bord` på tornets dator
  (`web/bord.html`, `web/bord.js`), där spelledaren leder vandringen i
  spelledarläget. Värdsidans rad har en knapp dit, och introt på tornets
  dator leder dit när spelledarläget är på. Överst står samma slags rad
  som på värdsidan; avsnitten är Vid porten (lobbyn och klassvalet),
  Vandringen är slut (slutet) och översikten.
- **översikten**: det spelledarbordet visar, med ingreppen, kartan,
  tabellen över vandrarna och loggen.
- **ingrepp**: det spelledaren gör under en vandring: tid, ormen,
  avslöjandet, tornets röst och rum.
- **spelarsidan**: sidan på `/` som spelarna öppnar, med vyerna introt,
  namnrutan, klassvalet, spelet och krönikan. Spelet har faserna lobbyn,
  klassvalet, vandringen och slutet (`lobby`, `klassval`, `spel`, `slut`). README beskriver under "Vad en spelare ser"
  vad en spelare möter.
- **spelskärmen**: vyn spelet på spelarsidan. Delarna heter:
  - **översta raden**: solmätaren, **prickarna** (vandrarna i sin färg,
    fyllda när de har valt drag) och tiden,
  - **teckenförklaringen**: listan uppe till vänster på kartan, under
    rummets namn, med de tecken som står i rummen just nu (`tecken_` i
    `data/texter.json`),
  - **tecknen**: statusen uppe till höger på kartan (klass, liv, skatt,
    drycker, tillstånd) och föremålen under dem, med förklaringar,
  - **frågetecknet**: styrningen, som ligger bakom ? på kartan,
  - **kortet**: rundans sammanfattning nere på kartan (`summary.go`,
    `sam_` i `data/texter.json`),
  - **panelen Nu**: draget, striden och vad som väntar,
  - **panelen Vandrare**: listan över vandrarna, bara i lobbyn och slutet,
  - **loggen**: raderna i högerkolumnen, grupperade per runda.
- **tornets dator**: datorn som kör servern.
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
  runda (`vandring` i tillståndet, `journeyOf` i `chronicle.go`), och i
  slutet ser alla resultaten, epilogen, hela krönikan och sin vandring.
- **tangenterna vid dörrarna**: piltangenterna som står vid dörrarna i ens
  rum på kartan snett ovanifrån (`ritaTangenter` i `klient/src/karta.ts`).
- **väskan**: föremålen under tecknen, med en plats för varje föremål man
  kan bära (`tecken` i `data/foremal.json`).
- **klassvalet**: skärmen efter lobbyn med ett spelkort för varje klass
  och en tidsgräns (`klassval_sekunder` i `data/klasser.json`).
- **klasser**: det vandrarna väljer i klassvalet (`classes.go`,
  `data/klasser.json`), till exempel Kartritaren.
- **föremål**: det vandrarna hittar och bär, högst tre (`items.go`,
  `data/foremal.json`). Vissa verkar av sig själva, andra används med 1–3,
  och några är förbannade.
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

- [docs/spelledaren.md](docs/spelledaren.md): värdsidan. Det mesta
  av del 1–3 är byggt. Avsnittet Läget säger vad som är byggt, vad som är på
  gång och vad som inte är planerat. Fråga vilka delar som ska byggas innan
  du börjar på något som inte står under På gång.
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
- [docs/isometrisk.md](docs/isometrisk.md): kartan snett ovanifrån med
  Three.js och 3D-modeller, i TypeScript i `klient/` och byggd med Vite
  till `web/iso/iso.js` (och WebRTC till `web/rtc/rtc.js`). Bygg om och
  checka in skripten efter en ändring i `klient/src/`. Avsnittet Stegen säger vad som är byggt.
