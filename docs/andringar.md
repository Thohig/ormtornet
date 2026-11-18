# Ändringar

Det som har ändrats i spelet, nyast först. README beskriver hur allt
fungerar. Här står vad som är nytt och vad det betyder för den som spelar
eller leder tornet.

## 2026-11-18

### Introt är en storm

Startsidan visar tornet i en storm vid solnedgången.

- Moln driver över solen, det regnar och blixtrar med 8–20 sekunders
  mellanrum. I blixten blir tornet en svart siluett mot himlen.
- Ett klick var som helst, eller Enter, startar ljudet: en låg drönarton i
  moll med en tritonus som sväller in, vind, regn, åska efter blixten och
  ibland ett dunk långt inne i tornet. Webbläsaren tillåter inget ljud
  innan man har klickat, så utan klick är introt tyst.
- Ryktena fladdrar lite, och ögat i tornet glöder ibland efter en blixt.
- "Gå in i tornet" visas först efter 30 sekunder, precis efter en blixt,
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
