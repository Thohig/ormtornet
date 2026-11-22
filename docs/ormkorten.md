# Förslag: ormkorten syns

Status: förslaget skrevs 2026-11-22, och alla fem delarna (A–E) är byggda
samma dag. README beskriver hur det fungerar under Skattkammaren och ormen.

## Varför

Ormleken fanns bara i servern. Den som tog en näve fick en rad i loggen,
och i kammaren stod en text om hur djupt ormen sov. Ingen såg att det var
kort som drogs, eller att risken växte för varje näve.

Det som byggdes visar det man redan visste som kort och bilder. Ingen får
veta mer än förut, och Ormtjusarens fördel, att se exakt hur många
vakna-kort som är kvar, finns kvar.

## Delarna

| Del | Vad | Commit |
|---|---|---|
| A | Kortet du drar vänds upp på kartan | `51eaf5b` |
| B | Sömnmätaren i skattkammaren | `d0587e5` |
| C | Uppvaknandet som ett stort kort för alla | `4563430` |
| E | Korten i rundans kort, Din vandring och krönikan | `ff91fca` |
| D | Leken och de dragna korten bredvid guldhögen | `484667d` |

### A. Kortet du drar

- Servern skickar `ormkort` till den som tar (`cardMsg` i `serpent.go`):
  kortets id, namn, text och guldet. De andra får inget kort.
- Spelarsidan vänder upp kortet till höger på kartan
  (`showSerpentCard` i `web/app.js`). Bilden ritas som SVG
  (`serpentArt`): ormen hoprullad med slutna ögon, med ett öga på glänt,
  eller två gula ögon i mörkret.
- Kortens namn står under `namn` i `data/ormen.json` och
  `data/en/ormen.json`. Bilderna har ingen text och är desamma på båda
  språken.

### B. Sömnmätaren

- I skattkammaren visar panelen Nu en mätare i tre steg som andas
  (`sleepMeter`), med samma tre nivåer som texten: djupt, oroligt och nära
  att vakna.
- Ormtjusaren får `vaknaKort` i tillståndet och ser vakna-korten som små
  kort bredvid mätaren. Andra får inte talet.

### C. Uppvaknandet

- `vaknar` skickas till alla när det är klart vad som hände var och en
  (`wakeSerpent`): kortets namn, det de märkte och, för den som stod i
  kammaren, guldet och livet den förlorade.
- Spelarsidan visar kortet Ormen vaknar stort mitt på kartan
  (`showWakeCard`).

### E. Korten i krönikan

- Servern sparar korten var och en drar (`stats.draws` i `chronicle.go`).
- Rundans sammanfattning får en rad med `ormkort` (`sam_ormkort` i
  `data/texter.json`), och spelarsidan ritar ett litet kort framför den
  (`miniCard`).
- Krönikans rad med siffror och Din vandring får korten som en rad små
  kort (`kort` i `chronLine` och `journeyView`, `cardRow` i `app.js`).

### D. Leken

- Kartan lägger leken med baksidan upp och de dragna korten med
  framsidan upp bredvid guldhögen (`byggLeken` i `klient/src/karta.ts`,
  `kortbaksida` och `kortframsida` i `krita.ts`).
- Högarna har bara tre höjder, samma som sömnmätaren, och ändras bara för
  den som står i kammaren (`sleep` i tillståndet). Andra ser högar som
  inte ändras. Den exakta höjden skulle annars avslöja risken nästan
  precis, och man skulle se ormens sömn på håll.
- När ormen vaknar flyger de dragna korten tillbaka in i leken (`blanda`).

## Inte byggt

- Ett ljud när kortet vänds upp. I dag hörs skattens ljud som förut.
