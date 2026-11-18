# Ormtornet

En dungeon crawler för 2–4 spelare i webbläsaren på ett lokalt nätverk.
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
- **spelarsidan**: sidan på `/` som spelarna öppnar: introt, namnrutan och
  spelet.
- **tornets dator**: datorn som kör servern.

Spelarna ser ordet "spelledare" i spelets texter, aldrig "värdsidan".

## Planerat

- [docs/spelledaren.md](docs/spelledaren.md): värdsidan. Det mesta
  av del 1–3 är byggt. Avsnittet Läget säger vad som är byggt, vad som är på
  gång och vad som inte är planerat. Fråga vilka delar som ska byggas innan
  du börjar på något som inte står under På gång.
