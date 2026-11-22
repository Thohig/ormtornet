# Den dödes guld (förslag)

Förslag 2026-11-22. Alla fyra stegen är byggda. Spårfinnaren hör ingenting genom väggarna, så i
stället för att höra hur mycket som rasslar ser hen guldet på kartan och
vem som tog det.

Sedan 2026-11-22 ligger guldet kvar där en vandrare dog. Den som går in
i rummet, eller står kvar där när rundan är slut, tar det (`guld.go`).
Men ingen får veta var det ligger, och inte heller om det ligger något
alls. När någon dör säger tornet bara: "Ett skrik hörs och dör bort. Alva
är död." Det enda sättet att hitta guldet är att råka gå in i rummet.

## Vad som önskas

För att gå efter guldet behöver man veta tre saker:

- **att** det finns guld,
- **var** det ligger, och
- **vad** som väntar där.

Alla vandrare, människor och AI, ska få samma ledtrådar. Ledtrådarna ska
ge en riktning och en chans, inte en karta till guldet. Vad som väntar
får man inte veta; det är risken.

## Steg 1: Skriket säger mer

Skriket när någon dör blir som ljuden från fällorna och striderna. Hur
mycket man hör beror på var man står.

- **Alla** hör vilken del av tornet skriket kom från: "Ett skrik hörs i
  tornets norra del och dör bort. Alva är död."
- **Mynten.** Bar den döda guld hörs också mynten: "… och sedan mynt som
  rasslar mot sten." Man vet att det finns guld, men inte hur mycket. Bar
  hen inget hörs bara skriket.
- **Den som står i ett grannrum** hör det från ett håll: "Genom väggen åt
  väster hörs ett skrik, och sedan mynt som rasslar mot sten. Alva är
  död." Det fungerar som `ljud_strid` och `ljud_falla` i dag (`hear` i
  `sounds.go`).
- **Den som står i samma rum** ser det hända, som i dag.

Nya texter i `data/fallor.json` och `data/en/fallor.json`: `dod_andra`
får `{region}`, och `dod_andra_guld`, `ljud_dod` och `ljud_dod_guld`
tillkommer.

## Steg 2: Rummet minns

När man sedan går in i rummet ska man märka vad som har hänt där, också
om guldet redan är taget.

- Rumsblocket säger var någon föll: "Här föll Alva. Golvet är mörkt av
  något som har torkat."
- Är guldet redan taget står det: "Några mynt har legat här. Någon har
  redan tagit dem." Då vet man att någon annan har varit där före.
- Hela rundor efteråt syns detta, som spåren efter fällorna
  (`spar_efter`), men utan tidsgräns.

## Steg 3: Klasserna hör olika

Två klasser kan få en egen fördel, som passar det de redan kan:

- **Spårfinnaren** hör hur mycket som rasslar: "ett tungt rassel" när
  den döda bar mer än 50 guld, annars "några mynt".
- **Kartritaren**, som redan ser alla rum som någon har hittat, får ett
  litet märke på kartan där någon har dött.

## Steg 4: AI-spelarna går dit

AI:n kommer ihåg var skriket med mynt kom ifrån, med del av tornet eller
riktning, precis som en människa hör det. Den går dit när:

- den har gott om liv (minst hälften),
- den inte redan är på väg ut med skatt, och
- solen räcker för att hinna dit och sedan till porten.

Hur villig den är styrs av personligheten (`data/ai.json`), med ett nytt
tal `guldlust`. Plundraren går gärna, den djärva ibland och den
försiktiga sällan. En AI som går in i ett rum där guldet redan är taget
glömmer det.

## Frågor

1. Ska mynten höras i skriket, eller ska man bara gissa att den döda bar
   något?
2. Räcker del av tornet för alla, eller ska de som står långt bort inte
   höra något alls?
3. Ska rummet minnas med namn ("Här föll Alva") eller utan?
4. Ska steg 3 med klasserna vara med, eller blir det för mycket?
5. Ska AI:n gå dit, och i så fall hur villigt?

## Ordningen

Steg 1 räcker för att guldet ska gå att hitta. Steg 2 gör att man förstår
vad man ser när man kommer fram. Steg 3 och 4 bygger på steg 1 och kan
vänta tills ett provspel har visat om någon går dit alls.
