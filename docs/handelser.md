# Slumpade händelser, ormen som flyttar sig och breven

Förslag och bygge 2026-11-24. Allt är byggt.

## Breven

Ett brev har en rad i marginalen som pekar mot något i just det här
tornet (`letterTip` i `items.go`, `brev_tips_*` i `foremal.json`):
åt vilket håll ristningen finns om vandraren inte har läst den, annars
närmaste varelse, och annars kammaren. Brevet om att inte gå till mitten
"första dagen" säger nu att man ska vänta tills man har hört andetagen.

## Kammaren och ormen

- **Nära mitten.** Kammaren ligger i ett av rummen i mittens 3 × 3,
  slumpat efter vandringens frö (`spridning` i `ormen.json`, `placeChamber`
  i `game.go`). Med spridning 0 ligger den mitt i tornet som förut;
  testerna använder det.
- **Ormen flyttar sig när den vaknar.** Den kryper till ett annat av rummen
  i mittens 3 × 3 (`moveChamber` i `serpent.go`). Den gamla kammaren blir
  ett vanligt rum med 15–45 guld löst på golvet åt den som kommer först
  (`spillGold` i `guld.go`). Den nya kammaren får de dörrar rummet redan
  hade, och en varelse som stod där flyr.

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

Händelserna gynnar bara den som drabbas eller hinner dit först, så att
tävlingen förblir var och ens.

## Balansen

Med den slumpade kammaren vann Kartritaren, som ser alla rum som någon har
hittat, för ofta (23–26 %), och Ormtjusaren för sällan (12–14 %).
Kartritaren har därför ingen egen dryck längre, och Ormtjusaren har Kraft
2 och Försvar 3. Efter det vinner alla klasser 13–21 % i simuleringen
(3000 vandringar, frö 1 och 5001, Normal och Svår).

## Inte byggt

- AI-spelarna känner inte till rasade rum och kan försöka gå in i ett
  några rundor; de går inte heller efter löst guld eller skinnet.
- Ett rasat rum syns inte på kartan; man märker det när man försöker gå
  in.
