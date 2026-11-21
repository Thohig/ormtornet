# Förslag: Synlighet

Status: förslaget skrevs 2026-11-18. Sedan 2026-11-21 styr den döda inte
längre varelserna och ser därför inte mer än andra; det som står om den
döda spelaren nedan gäller inte längre. Alla tre stegen är byggda (README,
"Vem ser vem").

Sedan 2026-11-21 är egen karta (`egen`) spelets enda läge. Inställningen
"Vandrarna syns" och flaggan `-synlighet` är borttagna, och det som står
nedan om `alla` och `moten` gäller inte längre. Portarna lottas alltid, och
spåren och ljuden från grannrummen finns alltid. En sparad inställningsfil
som har kvar `synlighet` läses som förut; valet ignoreras.

Varje vandrare ser bara sig själv. De andra syns bara när man möter dem i
samma rum. Ingen vet säkert var de andra är, hur långt de har kommit eller
vad de gör.

## Så fungerade det före steg 1

- Alla vandrares plats skickas till alla (`stateFor` i `game.go`) och ritas
  på kartan (`web/map.js`). De döljs bara i källarvalven och för den som
  har mörker.
- Kartan är gemensam. Ett rum som någon upptäcker syns för alla, med
  fällornas tecken och varelserna i rummet.
- Loggen berättar när någon upptäcker ett nytt rum och i vilken del av
  tornet (`nytt_rum_andra`), när någon kommer upp ur källaren
  (`kallare_ut_andra`), går ut (`andra_ut`) eller dör.
- Ljudet för nytt rum (`nyttrum`) spelas för alla.
- Ekofällan visar var den drabbade står, för alla.
- Mörkret döljer redan andra vandrare, okända varelser och resten av kartan
  för den drabbade. Källarvalven döljer den drabbade för alla.

AI-spelarna använder samma `stateFor` som människorna. Det de ser följer
därför med automatiskt när `stateFor` ändras.

## Inställningen

En ny inställning under Inför nästa vandring på värdsidan, med flaggan
`-synlighet`:

| Val | Flagga | Vad det betyder |
|---|---|---|
| Alla syns | `alla` | som före steg 1 |
| Bara vid möten | `moten` | steg 1 |
| Egen karta | `egen` | steg 1 och 2 |

Förvalet är `egen`. Det står i koden (`defaultSettings`), och spelledarens
sparade inställningar gäller före det. När vandringen börjar får alla veta vad
som gäller, på samma sätt som svårigheten.

## Steg 1: bara vid möten

Andra vandrare syns bara när de står i samma rum som du.

- `stateFor` döljer andras plats (`Hidden`, `X`, `Y` = -1) utom när de står i
  ditt rum. Samma kontroll som mörkret använder i dag.
- På kartan står ett svagt tecken där du senast såg någon, med namnet och
  rundan. Det sparas i klienten och försvinner när du ser personen igen.
- Listan över vandrare har kvar namnen men inte platsen. Där står bara "i
  tornet", "ute" eller "förlorad".
- Möten sägs i loggen som i dag (`meetIfNew`).
- Loggraderna om att någon går ut eller dör finns kvar. Det är tornet som
  berättar det, och det behövs för att man ska förstå läget.
- Loggraderna om nya rum och källaren tappar namnet och blir vaga: "Sten
  skrapar mot sten i tornets östra del." Texterna ändras i `texter.json`.
- Ekofällan visar fortfarande var den drabbade är. Det blir fällans verkliga
  straff.

Steg 1 är litet: framför allt `stateFor`, några texter och kartans tecken
för senast sedd.

## Steg 2: egen karta

Du ser bara rum du själv har varit i eller tittat in i.

- Varje vandrare får en egen mängd kända rum. Ett rum blir känt när du går
  in i det, kommer upp ur källaren där eller hittar en hemlig gång dit.
- Rummen finns fortfarande en gång i tornet. Fällan i ett rum slår bara till
  för den första som går in, och den som kommer efter ser fällans tecken.
- Varelser syns bara i rum du känner till.
- Ljudet för nytt rum spelas bara för den som går in. Loggraden om nya rum
  tas bort.
- Mörkret döljer fortfarande allt utom rummet du står i.
- Den döda spelaren som styr varelserna ser hela tornet så långt det har
  vuxit fram. Annars kan den inte jaga. AI-spelarens styrning (`botControl`)
  ser redan alla och kan få fortsätta så.
- AI-spelarna utforskar mer själva, eftersom de inte får de andras rum.
  Deras logik i `ai.go` behöver troligen inte ändras, men de bör provas.

Steg 2 passar ihop med avslöjandet, där var och en redan får egna ledtrådar.
Spelarna måste prata med varandra för att dela både karta och ledtrådar.

## Steg 3: spår och ljud (senare)

Mer stämning, men inte nödvändigt.

- Spår: när du går in i ett rum där någon annan har varit nyligen står det
  "färska fotspår leder västerut" i loggen, utan namn.
- Ljud genom väggarna: strid, en fälla eller någon som dricker i ett
  grannrum hörs, med riktning men utan namn.
- Ormen hörs i hela tornet som i dag, både när den rör sig och när den
  vaknar.

Texterna ligger i `texter.json`.

## Det som inte ändras

- Slutet: resultatet och krönikan berättar vem som gjorde vad. Det blir en
  del av upplösningen.
- Strid mellan vandrare: man ser motståndaren när man står i samma rum, och
  får vinken om hur mycket den bär som i dag. Det blir mer spännande när man
  inte ser vem som kommer.
- Andras skatt och val av drag är fortfarande hemliga.

## Beslut för steg 1

- Förval i steg 1: "bara vid möten". Flaggan heter `-synlighet` med
  värdena `moten` och `alla`.
- Tecknet för senast sedd finns. Det står kvar tills man ser vandraren igen,
  eller tills den har gått ut eller dött. Klienten minns det, så det
  försvinner om sidan laddas om.
- Ekofällan visar platsen för alla som förut.
- Den döda spelaren som styr varelserna ser alla vandrare.
- De vaga texterna ligger i datafilerna med `_dolt` och `andra_dolt`.

## Beslut för steg 2

- Steg 2, `egen`, är förval sedan 2026-11-18. Det gäller också den som har
  sparade inställningar utan synlighet.
- Ett rum blir känt när man står i det, även i mörker. En hemlig gång gör
  inte grannrummet känt, bara dörren dit.
- Från början känner var och en bara sitt eget startrum, med alla tre
  valen. De andras portar blir kända när någon går in i dem, och ett rum
  räknas som hittat först när någon har gått in i det (2026-11-19).
- Den döda spelaren som styr varelserna ser hela tornet så långt det har
  vuxit fram, och när vandringen är slut ser alla hela tornet.
- Texterna om att någon faller ner i källaren eller kommer upp, och om
  fällor, finns kvar som ljud från en del av tornet. Bara raden om nya rum
  tas bort.

## Beslut för steg 3

- Spår och ljud från grannrummen finns med både `moten` och `egen`, inte
  med `alla`.
- Spåren visar bara den som senast gick ut ur rummet, och bara de senaste
  3 rundorna (`spar_rundor` i `fallor.json`). Står den som lämnade spåren
  kvar i rummet ser man vandraren i stället. I mörker syns inga spår.
- Ljud hörs genom alla väggar till grannrummen, även utan dörr. Grannarna
  får texten med riktning i stället för den med delen av tornet.
- Ekofällan och klockan hörs i hela tornet som förut.
- Steg hörs inte genom väggarna. Det skulle avslöja för mycket.
