# Ormtornet

En dungeon crawler för 2–8 spelare i webbläsaren på ett lokalt nätverk.
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
  står alltid i den ordningen.
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
  namnrutan, klassvalet och spelet. Spelet har faserna lobbyn,
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
- **synligheten**: inställningen "Vandrarna syns" (`-synlighet`), med
  valen `egen` (egen karta, förval), `moten` och `alla`. Spåren och ljuden
  från grannrummen hör till den (`sounds.go`). Se
  [docs/synlighet.md](docs/synlighet.md).
- **stormen**: introts moln, regn, blixtar, åska och musik.
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
- [docs/fler-spelare.md](docs/fler-spelare.md): klasser, tornets storlek
  och upp till 8 spelare. Allt är byggt.
- [docs/foremal-och-hemligheten.md](docs/foremal-och-hemligheten.md):
  föremålen och viskningen till tornet. Allt är byggt; avsnittet Inte
  byggt säger vad som återstår.
- [docs/engelska.md](docs/engelska.md): engelska som språk, valt på
  värdsidan. Allt är byggt. Nya texter ska in på båda språken (`data/` och
  `data/en/`); ett test kräver att allt är översatt.
- [docs/synlighet.md](docs/synlighet.md): att vandrarna bara ser varandra
  när de möts, och på sikt bara sin egen karta. Alla tre stegen är
  byggda.
