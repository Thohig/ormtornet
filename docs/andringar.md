# Ändringar

Det som har ändrats i spelet, nyast först. README beskriver hur allt
fungerar. Här står vad som är nytt och vad det betyder för den som spelar
eller leder tornet.

## 2026-11-18

### Upp till 8 spelare

Tornet tar emot upp till 8 vandrare, AI-spelarna inräknade, och växer med
dem:

| Vandrare | Torn | Portar | Sol (med 40 på värdsidan) |
|---|---|---|---|
| 2–4 | 9 × 9 | 4 i hörnen | 40 |
| 5–6 | 11 × 11 | 8 | 50 |
| 7–8 | 13 × 13 | 8 | 60 |

- Tornet byggs om i lobbyn när någon kommer eller går, så att det alltid
  har rätt storlek. Den som redan är redo förblir det.
- I ett torn på 11 × 11 eller mer finns fyra nya portar mitt på varje sida,
  med tre dörrar inåt och en utgång rakt ut. De heter norr, söder, väster
  och öster i texterna.
- Solen räknas från värdsidans val för ett torn på 9 × 9 och växer lika
  mycket som vägen till kammaren.
- Ormleken får två extra sovkort per vandrare över fyra
  (`sovkort_per_vandrare_over_fyra` i `ormen.json`), så att risken växer i
  samma takt som med fyra.
- Spelledaren väljer 0–7 AI-spelare, eller 8 i spelledarläget, och minst
  1–8 vandrare. Värdsidan visar tornets storlek och solen under Vid porten.
- Väljer spelledaren en storlek blir tornet aldrig mindre än vad antalet
  vandrare kräver.
- Åtta färger och åtta platser i varje rum. AI-spelarna har tio namn.
- På en liten skärm visar kartan 9 × 9 rum runt den egna vandraren när
  tornet är för stort för att läsas. Ett tryck på kartan visar hela tornet,
  och ett tryck till går tillbaka.

### Tornets storlek

Spelledaren kan välja "Tornets storlek: 9 × 9, 11 × 11 eller 13 × 13" på
värdsidan, eller `-storlek`. Förvalet "automatisk" ger 9 × 9 tills vidare;
när upp till 8 spelare kommer räknas storleken ut efter antalet vandrare.
Ett större torn har längre väg till kammaren (10 och 12 steg i stället för
8), så solen kan behöva fler rundor. Byts storleken i lobbyn byggs ett nytt
torn direkt.

På en telefon blir rummen små i ett torn på 13 × 13.

### Klasser

Vandrarna väljer klass i lobbyn, en av två slumpade, i Talismans anda.
Åtta klasser med var sin styrka och svaghet: Legoknekten, Tjuven,
Örtkvinnan, Kartritaren, Ormtjusaren, Spårfinnaren, Lyktbäraren och
Pilgrimen. README beskriver dem under "Klasser". Spelledaren kan stänga av
klasserna på värdsidan.

Balansen är provad med AI-spelare: med samma klass för alla fyra når
13–35 % kammaren, mot 23 % utan klass. AI-spelarna använder inte alla
styrkor (till exempel tjuvens tysta näve och kartritarens karta), så
människor kan få andra siffror. Siffrorna justeras i `data/klasser.json`.

Rättat: spelledarens översikt kraschade under en vandring med egen karta.

### Ingen vet var de andra börjar

Gäller när vandrarna inte ser varandra (`egen` och `moten`).

- Lobbyn döljer de andra vandrarna på samma sätt som vandringen. Förut såg
  alla i lobbyn var de andra stod, och visste därför var de skulle börja.
- Portarna lottas när vandringen börjar. Var och en får veta sin egen port
  ("Du kliver in genom porten i tornets sydöst"), men inte de andras.
  Ordningen i listan och färgen säger inte längre vilken port någon har.
- Texterna när någon kommer till tornet nämner inte porten.
- Alla fyra portrummen syns på kartan från början, för alla och med alla
  val av synlighet. Man vet var utgångarna är, men inte vem som kom in var.
- Den som står i ett grannrum till en port hör den gnissla när någon går ut.

Med "alltid, för alla" börjar var och en vid porten som hör till dess plats
i listan, som förut.

Koden finns i `game.go` (`drawCorners`, `revealPorts`, `Corner` på
vandraren). Texterna står under `ansluter_dolt`, `valkommen_dolt`,
`start_port` och `ljud_port` i `texter.json`, och under `ansluter_dolt` och
`lamnar_dolt` i `ai.json`.

### Introt är en storm

Startsidan visar tornet i en storm vid solnedgången.

- Moln driver över solen, det regnar och blixtrar med 8–20 sekunders
  mellanrum. I blixten blir tornet en svart siluett mot himlen.
- Ett klick var som helst, eller Enter, startar ljudet: en låg drönarton i
  moll med en tritonus som sväller in, vind, regn, åska efter blixten och
  ibland ett dunk långt inne i tornet. Webbläsaren tillåter inget ljud
  innan man har klickat, så utan klick är introt tyst.
- Ryktena fladdrar lite, och ögat i tornet glöder ibland efter en blixt.
- "Gå in i tornet" visas först efter 15 sekunder, precis efter en blixt,
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
