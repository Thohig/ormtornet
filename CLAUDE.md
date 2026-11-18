# Ormtornet

En dungeon crawler för 2–8 spelare i webbläsaren på ett lokalt nätverk.
README beskriver spelet, flaggorna och filerna.

## Arbetssätt

- Skriv på svenska, både i svar och i spelets texter.
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
  `web/spelledare.html`). Avsnitten heter Bjud in, Vid porten, Översikt,
  Vandrarnas logg, Inför nästa vandring, Stäng tornet och Kommer någon
  inte in?
- **översikten**: den del av värdsidan som bara syns i spelledarläget, med
  ingreppen, kartan och tabellen över vandrarna.
- **ingrepp**: det spelledaren gör under en vandring: tid, ormen,
  avslöjandet, tornets röst och rum.
- **spelarsidan**: sidan på `/` som spelarna öppnar, med vyerna introt,
  namnrutan och spelet. Spelet har faserna lobbyn, vandringen och slutet
  (`lobby`, `spel`, `slut`). README beskriver under "Vad en spelare ser"
  vad en spelare möter.
- **tornets dator**: datorn som kör servern.
- **synligheten**: inställningen "Vandrarna syns" (`-synlighet`), med
  valen `egen` (egen karta, förval), `moten` och `alla`. Spåren och ljuden
  från grannrummen hör till den (`sounds.go`). Se
  [docs/synlighet.md](docs/synlighet.md).
- **stormen**: introts moln, regn, blixtar, åska och musik.
- **klasser**: det vandrarna väljer i lobbyn (`classes.go`,
  `data/klasser.json`), till exempel Kartritaren.

Spelarna ser ordet "spelledare" i spelets texter, aldrig "värdsidan".

## Ändringar

[docs/andringar.md](docs/andringar.md) säger vad som är nytt, nyast först.
Skriv dit när ett steg är klart.

## Planerat

- [docs/spelledaren.md](docs/spelledaren.md): värdsidan. Det mesta
  av del 1–3 är byggt. Avsnittet Läget säger vad som är byggt, vad som är på
  gång och vad som inte är planerat. Fråga vilka delar som ska byggas innan
  du börjar på något som inte står under På gång.
- [docs/fler-spelare.md](docs/fler-spelare.md): klasser (byggt), tornets
  storlek och upp till 8 spelare (på gång, i den ordningen).
- [docs/synlighet.md](docs/synlighet.md): att vandrarna bara ser varandra
  när de möts, och på sikt bara sin egen karta. Alla tre stegen är
  byggda.
