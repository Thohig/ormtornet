# Reprisen

Byggd 2026-11-25, idé 17 i [ideer.md](ideer.md). Allt är byggt.

## Vad den gör

På slutskärmen finns knappen **Se vandringen igen** bredvid Visa kartan.
Den öppnar ett lager över hela sidan med tornet som det blev, platt och
vridet som minikartan, så att norr pekar snett uppåt höger. Reprisen
spelar sig själv, en runda ungefär var 0,65 sekund:

- Allas vägar växer fram runda för runda som linjer i var och ens färg.
  Vägarna går lite vid sidan av rummens mitt, så att två vägar genom
  samma rum syns bredvid varandra.
- Var och en står som en prick där hen är. Den som har kommit ut är en
  ring vid porten, den som har dött ett kryss, och den som blev kvar i
  mörkret eller ligger i valven är blek.
- Under kartan står vandrarna i sin färg med guldet de bar just den
  rundan och hur det har gått (ute, död, kvar i mörkret, i valven).
- Knapparna är Till början, Spela eller Paus och ett reglage för rundan.
  Rundan står till höger: Före vandringen, Runda N och Slutet.
- Ett klick utanför, ✕ eller Escape stänger reprisen.

Rummen är alla som någon har sett, med dörrarna, portarna och kammaren
där den låg när vandringen tog slut. Inget i reprisen säger något om
ormens namn.

## Hur den är byggd

- `replay.go`: tornet skriver ner var alla står när vandringen börjar
  (`replayStart` i `startGame`), efter varje runda (`replayNote` efter
  tornets händelser) och när vandringen tar slut (`replayEnd` i
  `endGame`). Varje steg har rummet, hur det går (`inne`, `ute`, `dod`,
  `forlorad`), guldet och om man ligger i valven.
- Reprisen skickas i tillståndet som `repris` bara när vandringen är
  slut, till alla, och töms med krönikan inför nästa vandring.
- `web/repris.js` ritar den i SVG; stilen står under `#repris` i
  `web/style.css` och texterna under `repris_` i `data/spelarsidan.json`
  och `data/en/spelarsidan.json`.

## Att prova

- Är 0,65 sekund per runda lagom i en lång vandring med 40 rundor?
- Syns vägarna tydligt med åtta vandrare, och går det att följa en?
- Fungerar den på telefonen, där lagret är smalt?
