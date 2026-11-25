# Slumpade händelser, ormen som flyttar sig och breven

Förslag och bygge 2026-11-24, och nio händelser till samma kväll. Allt är
byggt.

## Breven

Ett brev har en rad i marginalen som pekar mot något i just det här
tornet (`letterTip` i `items.go`, `brev_tips_*` i `foremal.json`):
åt vilket håll ristningen finns om vandraren inte har läst den, annars
närmaste varelse, och annars kammaren. Brevet om att inte gå till mitten
"första dagen" säger nu att man ska vänta tills man har hört andetagen.

## Kammaren och ormen

- **Nära mitten.** Kammaren låg i ett av rummen i mittens 3 × 3,
  slumpat efter vandringens frö (`spridning` i `ormen.json`, `placeChamber`
  i `game.go`). Samma dag ersattes det av en fri men ganska rättvis plats
  ([ormens-plats.md](ormens-plats.md)); mitten finns kvar som
  `"plats": "mitten"`, och testerna använder det med spridning 0.
- **Ormen flyttar sig när den vaknar.** Den kryper till en ny plats med
  samma regler som när den placerades (`moveChamber` i `serpent.go`). Den
  gamla kammaren blir ett vanligt rum med 15–45 guld löst på golvet åt
  den som kommer först (`spillGold` i `guld.go`). Den nya kammaren får de
  dörrar rummet redan hade, och en varelse som stod där flyr.

## Händelserna

Från runda 4 kan en händelse hända i slutet av varje runda, med chansen
0,12 (`data/handelser.json`, `tornhandelser.go`). Alla inne i tornet får
ett spelkort med händelsens bild och rubrik, och åt vilket håll den hände;
den som står där får en egen text. Varje händelse står i krönikan
(`kronika_handelse_<id>`).

| Händelse | Bild | Vad som händer |
| --- | --- | --- |
| Tornet andas | `andas` | en ny dörr öppnas mellan två rum som finns |
| Guld ur sprickan | `spricka` | 20–45 guld rasar ut i ett rum ingen har varit i |
| Ljus genom gluggarna | `gluggar` | mörkret biter inte på någon i tre rundor |
| Ömsat skinn | `ormskinn` | Ormens skinn, ett förbannat föremål värt 30 guld, hamnar i ett rum ingen har varit i; med det i väskan väcker var fjärde näve ormen |
| Ett ras | `ras` | ett rum man har gått i går inte att gå in i på tre rundor; den som stod där tar en skada och faller ner i valven |
| Den döde vandrar | `doden` | ett väsen av vita ben reser sig där någon har dött |
| Klockslaget | `klockslag` | solen sjunker en handsbredd, och den som står i kammaren får en säker näve |
| Klockan stannar | `klocka_stilla` | solen sjunker inte den rundan; aldrig när solen är full |
| Draget | `draget` | alla som inte bär ljus får mörker en runda; inte medan ljuset faller genom gluggarna, och Lyktbäraren och veken klarar sig |
| Det klirrar | `klirrar` | alla hör åt vilket håll den som bär mest står, men inte vem; den rike och de bredvid får egna texter |
| Den glömda packningen | `packning` | ett föremål som inte är förbannat eller ett brev hamnar i ett rum ingen har varit i |
| Brevet under dörren | `brev_dorr` | en vandrare som inte står i mörkret eller valven får ett brev; bara hen märker det |
| Tornet sluter sig | `sluter` | en dörr mellan två rum stängs i tre rundor, men bara där det finns en annan väg mellan rummen |
| Yngel ur väggen | `yngel_vagg` | ett yngel kryper fram hos den som står ensam längst från portarna; bara hen märker det |
| Ormen drömmer | `drom` | två kort läggs i ormens lek, sovkort eller vakna-kort; alla märker att ormen rör sig, och den i kammaren vet åt vilket håll |
| Kryp som flyr | `kryp` | de som står högst tre rum från kammaren ser kryp komma från det hållet; bara de märker det |

Händelserna gynnar bara den som drabbas eller hinner dit först, så att
tävlingen förblir var och ens.

### Händelserna mellan vandrarna

Åtta händelser kom till 2026-11-25, var och en med ett eget spelkort,
där vandrarna påverkar varandra (`vandrarhandelser.go`):

| Händelse | Bild | Vad som händer |
| --- | --- | --- |
| Tornet drar tillbaka | `tillbaka` | den inne som bär minst guld får välja att dra tillbaka den som gick ut med mest; valet (T eller knappen Dra tillbaka i panelen Nu) är hens drag och gäller två rundor |
| Ormens blick | `blick` | en vandrare får ormens blick i sex rundor: var tredje näve väcker ormen, blicken går över till den man står i samma rum som, och när den släpper tar den 10 guld |
| Skottpengen | `skottpeng` | den som bär mest får ett pris på sig i tre rundor: den som träffar hen tar dubbelt så mycket |
| Skatan | `skata` | ett föremål som inte är förbannat flyger från en vandrare till en annan i ett annat rum; båda vet åt vilket håll den andre står |
| Väggen faller | `vagg_faller` | väggen mellan två vandrare i grannrum rasar, och var och en ser den andres rum på kartan |
| Kistan med två lås | `kista` | 30–60 guld i en kista i ett rum någon har gått i; den går bara upp när två står där, och den som bänder (S) tar allt |
| Spåren glöder | `spar_glod` | i tre rundor ser alla på kartan rummen där någon går |
| Golvet vrider sig | `vrider` | två vandrare i vanliga rum utan varelser byter plats |

**Tornet drar tillbaka** är aldrig frivilligt för den som dras in, men
den som drar får ormens blick. Den som dras in kommer in genom samma port
med fullt liv och utan föremål (de blev guld vid porten). Hälften av
guldet ligger säkert ute och läggs till när hen går ut igen; dör hen
eller blir kvar räknas bara den säkra hälften i placeringen. Det händer
bara en gång per vandrare, och bara medan porten är öppen och solen har
minst åtta handsbredder kvar (`sol_minst`).

Blicken och skottpengen står bland tecknen (⊙ och ¤) med rundorna kvar.
Kistan står som en rad i panelen Nu i rummet där den finns.

Chansen per runda är fortfarande 0,12, så varje händelse kommer mer
sällan. En simulering av 1500 vandringar med fyra AI-spelare (frö 1)
ändrades lite: guldet hos dem som kom ut med skatt blev 53,2 i stället
för 54,2, och alla klasser vinner 19–31 % som förut.

De nio sista kom till 2026-11-24. En händelse kan sakna `text`, och då
märker bara den det händer något (`har`). `bredvid` är texten för de
andra i samma rum, `har_annars` det den får när det gick åt andra hållet
(ormen som drömmer), `avstand` hur nära kammaren krypen syns och
`ormkort` hur många kort ormen som drömmer lägger i leken. Chansen per
runda är densamma, 0,12, så de gamla händelserna kommer mer sällan än
förut. I en simulering av 1500 vandringar med fyra AI-spelare ändrades
lite: något mindre guld hos dem som kom ut (53 mot 55) och Kartritaren
vann 19 % i stället för 21 %.

## Balansen

Med den slumpade kammaren vann Kartritaren, som ser alla rum som någon har
hittat, för ofta (23–26 %), och Ormtjusaren för sällan (12–14 %).
Kartritaren har därför ingen egen dryck längre, och Ormtjusaren har Kraft
2 och Försvar 3. Efter det vinner alla klasser 13–21 % i simuleringen
(3000 vandringar, frö 1 och 5001, Normal och Svår).

## Inte byggt

- AI:n drar tillbaka den som gick ut ungefär sex gånger av tio när den
  får välja, och bänder upp kistan när någon annan står i rummet, men den
  går inte till kistan och söker inte upp någon för att bli av med
  ormens blick eller jaga den med ett pris på sig.
- Kistan syns inte på kartan, bara i panelen Nu när man står i rummet.
- Den som dras tillbaka och dör står som död på slutskärmen, fast den
  säkra hälften av guldet räknas i placeringen.

- AI-spelarna känner inte till de nya händelserna: de går inte efter den
  glömda packningen och bryr sig inte om klirret eller krypen.
- En sluten dörr syns inte på kartan som sluten; den är borta tills den
  öppnas igen.
- AI-spelarna känner inte till rasade rum och kan försöka gå in i ett
  några rundor; de går inte heller efter löst guld eller skinnet.
- Ett rasat rum syns inte på kartan; man märker det när man försöker gå
  in.
