# Samlat förslag: spelkort för alla situationer

Status: förslag 2026-11-22, byggt samma dag, alla delar A–G. Avsnittet
Läget säger vad som blev annorlunda än förslaget, och
[att-gora.md](att-gora.md) vad som ska provas. Efteråt byggdes
vandringens största stunder, klasskorten och minikorten i krönikan; se
avsnitten sist.

## Målet

Ormkorten gjorde skattkammaren till något man ser: ett kort vänds upp,
med en bild och en text, och samma kort syns sedan litet i krönikan.
Resten av tornet berättas bara i loggen och i rundans kort.

Förslaget gör allt som händer en vandrare till ett kort ur samma lek, och
låter rundans kort gå upp i det. Då säger spelet varje sak på två ställen
i stället för tre:

- **Loggen berättar vad som hände**: rummet, det du ser och hör, fällans
  och varelsens text och allt om de andra.
- **Spelkortet visar vad det gjorde med dig**: en bild, en rubrik och
  siffrorna (liv, guld, föremål, gift, mörker, att du sitter fast).

## Spelkortet

- **Ett kort per runda, när något har ändrat dig.** Det ersätter rundans
  kort (`showRoundCard`) och ormkortet till höger (`showSerpentCard`) och
  vänds upp på samma ställe som ormkortet. En runda där inget ändrade dig
  ger inget kort, som i dag.
- **Framsidan:** rubriken och bilden kommer från rundans viktigaste
  händelse. Under bilden står händelsens text, och längst ner rundans
  rader från i dag (`sam_` i `data/texter.json`): "−2 liv. Du har 4
  kvar.", "+31 guld. Du bär 140." Det farliga står först, som i dag.
- **Fler saker samma runda** blir minikort i en rad längst ner på kortet.
  Sedan 2026-11-24 står de två nästa på datorn i en rad till vänster om
  kortet, lika stora och med sin text, och raden står 4 sekunder längre
  för varje kort till (`renderRow` i `web/app.js`). Ett varsel eller
  viskningen som kommer medan kortet står uppe läggs också där. Fler än
  tre kort, och alla utom det första på telefonen, står som minikort.
  Ordningen är döden, ormen, varelsen, fällan, föremålet och sist resten,
  och den står i `data/`.
- **Kortet glider undan** när du väljer nästa drag, efter en stund eller
  vid ett klick, som rundans kort i dag.
- **De stora korten** mitt på kartan är för det hela tornet märker, som
  uppvaknandet i dag, och för vandringens slut.
- **Minikorten** står framför raderna i Din vandring, i krönikan och i
  loggen där en rad hör till ett kort.

### Reglerna som gäller som förut

- **Namnen.** Fällor och varelser får den korta beskrivningen som rubrik
  (`kort`, med stor bokstav), aldrig namnet. Ormen, föremålen och
  skattkammaren har sina namn.
- **Synligheten.** Du får bara kort för det som hände dig och för det som
  hela tornet märker. Det någon annan råkar ut för i ditt rum står i
  loggen.
- **Hemligheten.** Inget kort säger mer än loggen och rundans kort säger
  i dag.
- **Innehållet i `data/`.** Varje fälla, varelse och föremål får ett fält
  `bild` med motivet. Nya texter kommer in på båda språken.
- **Bilderna** ritas som SVG i koden i ormkortens stil, utan text, och
  är desamma på båda språken.

## Korten

64 motiv, uppdaterat 2026-11-24. De står i `cardImages` i
`spelkort.go` och ritas i `web/kort.js`; alla syns i
`verktyg/kortbilder.html`.

**Ormkorten (3).** Ormen sover, Ormen rör sig och Ormen vaknar (stort,
för alla).

**Fällorna (10).** Klev du förbi får du samma kort nedtonat.
Ett galler av spetsiga järnstänger, ett golv som ger efter, ett golv som
vrider sig, pilar ur små hål i väggen, en söt gulgrön dimma, en grop full
av små ormar, ett väldigt genomskinligt skinn, en väldig klocka i taket,
en varm vindpust som släcker lågan och ett eko som inte vill tystna.

**Varelserna (5).** Kortet kommer när varelsen dyker upp, och ett
överkorsat minikort när den faller eller flyr. Något blekt och
mångbent, något långt och väsande, en gestalt i fjällrustning, en
grönaktig skepnad och ett väsen av vita ben.

**Föremålen (19).** Kortet har föremålets namn och vad det gör, så att
raden om det kan lämna loggen. Rep med kroken kvar, Fjällamuletten,
Extra veke, Bucklig sköld, Väskan med dubbel botten, Kompass utan nord,
Rökpulver, Spegelskärva, Kastkniv, Benluren, Sövande flöjt, Timglas av
ben, Rostigt svärd och Fjällbrynja. De förbannade vänds en gång till när
de visar sin baksida: Den svarta pärlan, Myntet med hål, Kronan av rost,
Svärdet med rött fäste och Ormens skinn.

**Tornets händelser (16).** Var och en har sitt kort
(`data/handelser.json`, docs/handelser.md): Tornet andas, Guld ur
sprickan, Ljus genom gluggarna, Ömsat skinn, Ett ras, Den döde vandrar,
Klockslaget, Klockan stannar, Draget, Det klirrar, Den glömda
packningen, Brevet under dörren, Tornet sluter sig, Yngel ur väggen,
Ormen drömmer och Kryp som flyr.

**Det som händer dig (9).** En annan vandrare anfaller dig, en dryck
(hittad eller drucken), guld (taget, det som ligger kvar efter en död
och löst guld), källarvalven (när du faller ner eller kommer upp), ett
varsel, en del av ormens namn, viskningen eller att du blir invigd, ett
brev, och döden (stort, bara för dig).

**Tornets stora kort (3).** En port stängs, solen går ner de sista
rundorna, och du kommer ut med skatt (bara för dig).

**Klasskorten (8).** Legoknekten, Tjuven, Örtkvinnan, Kartritaren,
Ormtjusaren, Spårfinnaren, Lyktbäraren och Pilgrimen, i klassvalet och
som ett litet klasskort bland tecknen.

**Facklan (1).** Motivet när en runda inte har något eget kort, och när
ett motiv i `data/` saknas (kontrollen i `validate.go` säger till).

Berättelsen i slutet använder samma motiv för höjdpunkterna.

## Loggen och kortet

Det mesta av siffrorna står redan i rundans kort och inte i loggen. Det
som flyttas:

| Från loggen | Till |
| --- | --- |
| Raden om skatten efter en näve (`skatt_totalt`) | ormkortets botten |
| Siffrorna i striden mot en vandrare (`pvp_byte`, `pvp_bytt`, `pvp_traffad`, `pvp_dodar`): "Du har nu 140", "Du förlorar 2 liv" | kortets botten; raden i loggen behåller vad som hände |
| Raden om vad ett föremål gör när du hittar det | föremålskortet; förklaringen finns kvar i tecknen |
| Rundans kort över kartan | spelkortet |

Det som stannar i loggen: rumsblocket, det du ser och hör, fällornas,
varelsernas och ormens egna texter, allt om de andra och lådan med hela
loggen. Loggen är det som finns kvar under vandringen när kortet har
glidit undan, och det som en skärmläsare läser. Rundans rader finns kvar
i Din vandring som i dag.

Under bygget ska varelsernas och fällornas texter i `data/` gås igenom
efter siffror som då står två gånger.

## Svar på frågorna

| Fråga | Svar |
| --- | --- |
| Hur ofta vänds ett kort? | Lika ofta som rundans kort i dag, eftersom spelkortet ersätter det: en gång per runda, när något ändrade dig. |
| Rubriken på fällor och varelser? | Den korta beskrivningen, med stor bokstav. |
| Ersätter kortet raden i loggen? | Nej för berättelsen, ja för siffrorna. |
| Vilka delar? | Alla, i ordningen nedan. A ger det mesta av vinsten redan med ormkortens och en neutral bild. |

## Delarna

| Del | Vad |
| --- | --- |
| A | Ramen: rundans kort och ormkortet blir ett spelkort, med bild, rubrik, text, rundans rader och minikort. Servern väljer rundans viktigaste händelse (`summary.go`). Rundor utan motiv får en neutral bild, en fackla. Raden `skatt_totalt` lämnar loggen. |
| B | Fällorna: tio bilder och kortet när du klev förbi. |
| C | Varelserna och striden: fem bilder, striden mot en vandrare och siffrorna ur `pvp_`-raderna. |
| D | Föremålen och dryckerna: femton bilder, förbannelsen som vänder kortet och förklaringen som lämnar loggen. |
| E | Avslöjandet, brevet och källarvalven. |
| F | Tornets stora kort och slutet: porten, solen, ut och döden. |
| G | Minikorten i Din vandring, krönikan och loggen. |

Varje del är ett steg med en commit. README, `docs/andringar.md`, CLAUDE.md
(orden kortet och ormkortet) och `docs/att-gora.md` (det som ska provas)
uppdateras i varje steg.

## Inte med

- Kort för det man ser hos andra vandrare.
- Ljud när ett kort vänds.
- Kort i lobbyn och klassvalet; klassvalets kort finns redan.

## Läget

| Del | Byggd |
| --- | --- |
| A | ja |
| B | ja |
| C | ja; också varelsernas texter förlorade sina siffror |
| D | ja, med ett kort för guld man hittar när man söker |
| E | ja; varslen, namnets delar och viskningen kommer mellan rundorna som ett eget kort |
| F | ja; solen visas stort när tre rundor är kvar och när den går ner med en invigd kvar, och den som blir kvar får kortet Solen gick ner |
| G | ja; krönikan i slutet hade kvar ormkorten på raden med siffror (borttagen från slutskärmen 2026-11-24, guldboken säger vilket kort varje näve gav), och Din vandring får alla kort |

Annorlunda än förslaget:

- Ett kort som kommer mellan rundorna (varsel, namnets delar, att bli
  invigd) skickas genast som meddelandet `kort` och står i kö efter det
  kort som är uppe. Tornets stora kort (port, sol) kommer som `storkort`.
  En ny runda ersätter genast spelkortet som står uppe.
- Guld man hittar när man söker fick ett eget kort, och sökandets,
  dryckens och varelsernas texter förlorade sina siffror.
- Facklan är bilden för en runda utan eget kort.
- `verktyg/kortbilder.html` visar alla motiv, och kontrollen av
  datafilerna säger till om en bild i `bild` inte finns.

## Efteråt: vandringens största stunder

Byggt 2026-11-22 efter ett förslag om spelkorten i krönikan i slutet (D
av fem: minikort vid krönikans händelser, resultaten som kort, var och
ens lek, största stunderna och att klicka upp ett minikort). Krönikans
händelser vet sin nyckel, vem de gällde och sin bild (`chronEvent`), och
`stunder.go` väljer högst fem: de viktigaste slagen först, högst två av
samma slag, lagda i den ordning de hände. Slagen står under
`spelkort.stunder` i `data/texter.json`. Senare byggdes också att man kan
klicka upp ett minikort (se Efteråt: minikorten).

Stunderna gick upp i berättelsen om vandringen 2026-11-24: varje
höjdpunkt i berättelsen har sitt spelkort bredvid sig, och `stunder.go`
ersattes av `berattelse.go` ([berattelsen.md](berattelsen.md)).

## Efteråt: klasskorten

Byggt 2026-11-22 i fyra steg, efter ett förslag om att klasskorten skulle
likna spelkorten:

1. Klasskorten i klassvalet har spelkortens pergament, ram och mörka
   text, med namnet i versaler, styrkan i grönt och svagheten i rött.
2. Varje klass har en bild i `web/kort.js` och ett fält `bild` i
   `data/klasser.json`; kontrollen säger till om bilden saknas.
3. Klassvalet börjar med korten på baksidan, och de vänds upp ett i taget.
   Korten byggs en gång per klassval (`renderChoose`, `buildClassCard`)
   och uppdateras sedan. Ett kort som någon annan tar vänds tillbaka och
   visar baksidan med namnet.
4. Under vandringen står klassen bland tecknen som ett litet klasskort;
   ett klick visar det stort mitt på kartan (`showClassCard`).

Klasskortet står inte i Din vandring, vilket förslaget sa var möjligt.

## Efteråt: minikorten

Minikorten var 12 × 16 pixlar och gick inte att läsa. Byggt 2026-11-22:

- **Klicka upp ett kort.** Ett minikort med rubrik (`Kort.mini` med
  `rubrik`) bär kortet i `data-kort`. Ett klick, eller Enter, visar det
  stort med sin text i ett lager över hela sidan (`#kortvisning`,
  `showCardPopup` i `web/app.js`); ett klick eller Escape stänger det.
  Det gäller i loggen, på spelkortet och i krönikan. Din vandring sparar
  därför kortens text.
- **Krönikans kort.** I Din vandring står rundans kort som små spelkort,
  104 pixlar breda, med pergament, en bild på 88 × 76 och rubriken, och
  hela kortet går att klicka upp.
- **Ormkorten på raden med siffror** står ihopslagna med ett kort på
  30 × 40, namnet och antalet: "Ormen sover ×3 · Ormen rör sig"
  (`cardRow`).

Inte byggt av förslagen för krönikan: minikort framför krönikans
händelser, resultaten som kort på bordet och var och ens lek som ett
eget avsnitt.
