# Ormtornet

En dungeon crawler för 2–4 spelare i webbläsaren, på ett lokalt nätverk.
Mitt i ett torn på 9x9 rum sägs en jättelik orm sova kring en skatt.
Ta så mycket du vågar och kom ut innan solen går ner.

## Starta

Kräver Go (`brew install go`). Det är bara den som kör tornet som behöver
spelet. De andra öppnar en adress i webbläsaren.

```bash
./starta.sh
```

Skriptet bygger spelet, öppnar introt i webbläsaren och håller datorn vaken
så länge tornet är öppet. På datorn som kör tornet leder introt till
spelledarens sida. Flaggor skickas vidare, till exempel `./starta.sh -ai 2`.
Är tornet redan igång öppnas bara introt.

Det går också att starta direkt med `go run .`.

### Ormtornet.app

```bash
./skapa-app.sh
```

skapar `~/Projekt/Ormtornet.app` med tornets ikon (ange en annan mapp som
argument om du vill). Dubbelklick på appen bygger spelet från den här mappen,
startar tornet i bakgrunden och öppnar introt, som leder till spelledarens
sida. Inget
terminalfönster behövs. Tornet stängs med knappen "Stäng tornet", eller av sig
självt när ingen spelare är ansluten och spelledarens sida har varit stängd i
10 minuter. Loggen hamnar i `/tmp/ormtornet.log`. Appen behöver bara skapas om
när den här mappen flyttas.

Servern lyssnar på `0.0.0.0:8080` och skriver ut adressen spelarna ska
öppna, till exempel `http://Din-Mac.local:8080`. Alla Macar i samma nätverk
kan öppna den i Safari eller Chrome.

Stäng tornet med Ctrl+C i terminalen eller knappen "Stäng tornet" på
spelledarens sida. Spelarna får då ett meddelande om att tornet stänger.

Inställningarna ändras enklast på spelledarens sida (se nedan) och sparas
till nästa gång. Flaggorna nedan gäller före det sparade, men bara den gången.

Flaggor:

- `-addr 0.0.0.0:8080` adress och port
- `-data data` mapp med datafilerna
- `-sol 40` antal rundor innan solen går ner
- `-turtid 30s` längsta tid att välja drag per runda
- `-min 2` minsta antal spelare för att börja, datorspelare inräknade (sätt `-min 1` för att prova ensam)
- `-ai 2` antal datorstyrda vandrare (0–3, eller 4 i spelledarläget)
- `-avslojande=false` stänger av avslöjandet, för spelare som redan känner till tornets hemlighet
- `-oppna` öppnar introt i webbläsaren när servern har startat; på den här datorn leder det till spelledarens sida
- `-vaken` håller datorn vaken så länge servern är igång
- `-autostopp 10m` stänger tornet när ingen spelare är ansluten och spelledarens sida har varit stängd så länge
- `-installningar fil` var spelledarens inställningar sparas (standard `~/Library/Application Support/Ormtornet/installningar.json`, tom för att inte spara)

### Spelledarens sida

Öppna `http://localhost:8080/spelledare` på datorn som kör servern. Dit
kommer man också efter introt på den datorn. Där står
adressen spelarna ska öppna, med knappar som kopierar en färdig inbjudan att
skicka i en gruppchatt. Sidan visar också vilka som har anslutit och är redo,
och vad man kan göra om någon inte kommer in.

Under Inställningar väljer spelledaren antal AI-spelare, minsta antal
vandrare, solnedgång, turtid och avslöjandet. AI-spelare kommer och går
direkt i lobbyn, resten gäller från nästa omgång. Under en pågående vandring
är inställningarna låsta. Om tornet är fullt när en människa ansluter i
lobbyn kliver en AI-spelare åt sidan, och kommer tillbaka om platsen blir
ledig igen. Knappen Starta nu börjar vandringen även om inte alla har
tryckt Enter.

Om en spelare lämnar mitt i en vandring tar en AI-spelare över vandraren
efter två rundor (ändras under "Frånvarande vandrare", eller sätts till "står
still"). Spelledaren kan också trycka "AI tar över" direkt. AI:n fortsätter
med vandrarens plats, skatt och liv, och de andra ser namnet märkt "(AI)".
Kommer spelaren tillbaka och skriver samma namn får den genast tillbaka sina
steg. I resultatet märks den som AI har spelat åt med "~". Den som inte kommer
tillbaka innan vandringen är slut är inte med i nästa.

Med "Spelledarläge: på" leder spelledaren tornet utan att spela. Ingen kan
då ansluta som vandrare från spelledarens dator, och spelarna får veta att
tornet har en spelledare. Spelledaren startar vandringen med Starta nu, även
när alla vandrare är AI-spelare. Sidan visar då en översikt: hela tornet så
långt det har vuxit fram, med fällor som väntar och alla varelser, varje
vandrares plats, liv, skatt, drycker, tillstånd och valda drag, hur djupt
ormen sover, avslöjandets nivå och allt vandrarna har fått läsa.

I spelledarläget kan spelledaren också göra ingrepp under en vandring, och
varje ingrepp skrivs i loggen:

- **Tid:** pausa och fortsätta, ge rundan 30 sekunder till, avgöra rundan
  direkt, och flytta solen fram eller tillbaka (den går aldrig ner av ett
  ingrepp, det gör bara rundorna).
- **Tornets röst:** ett meddelande till en vandrare eller alla, i en egen
  färg och med en viskning.
- **I ett rum:** lägga en fälla eller släppa in en varelse i ett rum som
  redan finns.
- **Ormen:** väcka den, eller lugna den genom att ta ett vakna-kort ur leken.
- **Vandrarna:** ge eller ta liv, guld och drycker. Vandraren märker att
  något händer, men inte varifrån.

Utan spelledarläge visar sidan inga skatter, positioner eller drag, så spelledaren kan spela med i
en annan flik ("Gå in själv", som leder till introt och sedan namnrutan). Andra datorer i nätverket kommer inte åt sidan.

## Spela

1. Startsidan visar tornet i solnedgången medan rykten om det tonar fram.
   Klicka på "Gå in i tornet" (det startar ljudet). Ryktena står i
   `data/texter.json` under `rykten`, och under `rykten_invigda` för den som
   spelar med `-avslojande=false`.
2. Skriv ditt namn och tryck Enter.
3. Tryck Enter igen när du är redo. Vandringen börjar när alla är redo.

Spelet går i rundor. Alla väljer sitt drag samtidigt, och när alla har valt
(eller turtiden tar slut) utförs alla drag på en gång. Du kan ändra dig tills
dess. Den som inte hinner välja blir stående.

| Tangent    | Drag                              |
|------------|-----------------------------------|
| Pilar      | gå åt det hållet                  |
| Mellanslag | vänta                             |
| T          | ta skatt (bara i skattkammaren)   |
| S          | sök igenom rummet                 |
| D          | drick en dryck                    |
| A          | anfall (i strid)                  |
| P          | parera (i strid)                  |
| U          | undanmanöver (i strid)            |
| Tab        | byt motståndare, eller varelse när du är död |
| Enter      | redo / spela igen                 |

Solen sjunker en handsbredd per runda och visas som en mätare överst.
Gå ut genom en port i något av hörnen för att lämna tornet. Den som kommer
ut med skatt räknas, och den med mest skatt vinner. Den som är kvar inne när
solen har gått ner är förlorad.

Andras skatt och val av drag syns inte förrän spelet är slut.

### Strid mellan vandrare

Står två vandrare i samma rum utan någon varelse kan de slåss med samma
sten-sax-påse som mot varelserna (A, P, U, och Tab för att välja motståndare
när flera står där). Den som inte försvarar sig blir träffad om den andra
anfaller, även om den går ut ur rummet. Varje träff kostar 1 liv och river
åt sig en tredjedel av motståndarens guld (minst 1). Den som dödar tar allt.
Hur mycket någon bär syns inte, men man får en vink när man möter dem.
Portrummen är fredade, och varje runda med strid gör ormen lite oroligare.
AI-spelare anfaller bara människor, utom när alla vandrare är AI-spelare.
Spelledaren väljer under "Strid mellan vandrare": på, efter halva tiden
eller av.

### Spela mot datorn

Välj antal AI-spelare på spelledarens sida, eller med `-ai`:

```bash
go run . -ai 3
```

Datorspelarna står redo vid portarna. Omgången börjar när alla människor har
tryckt Enter; minst en människa måste vara med. I listan över vandrare står
de som "(dator)".

- De ser bara det en människa skulle se: kända rum, varelser i kända rum och
  sin egen skatt. De ser inte okända rum, ormkorten eller andras skatt.
- De utforskar mot mitten, tar skatt i kammaren så länge ormen sover djupt
  nog, och går mot närmaste port när solen börjar bli låg, väskan är full
  eller livet tryter.
- I strid gissar de varelsens vanor utifrån dess beskrivning och lär sig av
  det de ser. De flyr när nästa slag kan döda dem.
- Det finns två personligheter som turas om: en djärv och en försiktig.
- Dör en datorspelare först tar den över varelserna och flyttar dem mot
  närmaste levande vandrare.
- Ledtrådarna om tornets hemlighet går bara till människor.

Namn, personligheter och betänketid ligger i `data/ai.json`.

### Skattkammaren och ormen

I kammaren mitt i tornet kan du ta skatt (T) en gång per runda. Varje gång
dras ett ormkort ur en lek:

- **Sover** eller **rör sig**: du får 15–45 guld och kortet tas bort ur leken.
  Rör sig ormen märker alla i tornet det.
- **Vaknar**: kortet blir kvar i leken. Alla i kammaren förlorar hälften av
  sin skatt och 2 liv och kastas ut genom en slumpad öppning. Väsandet hörs
  i hela tornet.

Eftersom bara sovkorten tas bort ökar risken för varje näve. Leken börjar med
13 sovkort och 2 vakna-kort. När ormen har vaknat blandas en ny lek med 2
färre sovkort än förra gången. Står du i kammaren ser du hur djupt ormen
sover och hör andetagen bli tätare.

### Strider

Ett nytt rum kan i stället för en fälla ha en varelse. Den syns på kartan för
alla. Står du i samma rum som en varelse är du i strid. Båda väljer i hemlighet:

- **parera** slår **anfall**
- **anfall** slår **undanmanöver**
- **undanmanöver** slår **parera**

Vinnaren träffar. Väljer båda anfall träffar båda; annars händer inget vid
lika val. Varje sorts varelse har egna vanor, så det lönar sig att lära sig
dem. Besegrar du en varelse får du dess guld.

Gör du något annat än att slåss (söker, dricker, väntar) får varelsen ett
fritt drag mot dig. Du kan fly genom att gå ut ur rummet, men varelsen får
ett fritt drag mot din rygg först.

Den första spelaren som dör tar över tornets varelser. Med Tab väljer den spelaren en
varelse och med pilarna flyttar den ett rum per runda (bara till kända,
vanliga rum utan annan varelse). Mellanslag låter varelserna vila. Rundan
väntar på det valet som på alla andras.

### Rum och fällor

När någon går in i ett nytt rum kan en fälla slå till. Ett tecken i rummet
visar sedan för alla att något har hänt där.

| Fälla           | Vad som händer                                      |
|-----------------|-----------------------------------------------------|
| Fallgaller      | du sitter fast i 2 rundor                           |
| Rasande golv    | 1 skada, du faller ner i källarvalven i 2 rundor och kommer upp någon annanstans; ingen ser var du är |
| Roterande rum   | du kastas ut genom en slumpad öppning               |
| Pilar           | 2 skada                                             |
| Giftig dimma    | 1 skada nu och 1 skada per runda i 2 rundor         |
| Ormgrop         | 2 skada och fast i 1 runda                          |
| Ömsat skinn     | fast i 1 runda                                      |
| Klocka          | solen sjunker 2 extra steg för alla                 |
| Mörker          | lyktan slocknar i 4 rundor; du ser bara ditt eget rum |
| Eko             | alla får veta var du är                             |

Du har 6 liv. Den som når noll liv dör.

Ett vanligt rum kan sökas igenom en gång (S). Man kan hitta skatt, en dryck,
en hemlig gång till ett grannrum, eller utlösa en fälla. En dryck (D) ger
tillbaka 2 liv och tar bort gift.

Om du laddar om sidan och skriver samma namn igen får du tillbaka din vandrare.

## Avslöjandet (spoilervarning)

Läs inte det här om du ska spela utan att veta.

Tornet är ormen. Väggarna är dess kropp, valven dess revben, trappan dess
ryggrad, skattkammaren dess mage och porten dess mun. Spelarna upptäcker det
gradvis, i fyra nivåer som följer solen:

| Nivå | När (andel sol kvar) | Text | Bild | Ljud |
|------|----------------------|------|------|------|
| 1 | över 70 % | små konstigheter: varma väggar, söt lukt, fjällformade stenar | vanliga väggar | ett svagt, långsamt dunk |
| 2 | 70–40 % | tornet andas | väggarna sväller och krymper i braille | dunket blir ett hjärta |
| 3 | 40–15 % | bevis: revben, kotor, halvt upplösta kvarlevor | väggar som revben, köttfärg, hjärtslag i vågor | snabbare hjärta, kluckande |
| 4 | under 15 % | det sägs nästan rakt ut; porten stänger sig | rött, porten blir en mun med tänder | starkt och snabbt hjärta |

- Varje rum har beskrivningar för alla fyra nivåerna, och spelare som står i
  samma rum kan få olika texter.
- Var fjärde runda, och vid varje nytt steg, får varje spelare en egen
  ledtråd. Ingen får samma som någon annan, så spelarna måste prata med
  varandra för att få hela bilden.
- På nivå 4 får varje spelare nya förklaringar till de fällor och varelser
  hen har mött (fallgallret var tänderna, de roterande rummen tarmarna ...).
- När spelet är slut säger tornet hemligheten rakt ut, och porten stängs.

Med `-avslojande=false` stannar tornet på nivå 1: inget dunk, inga
pulserande väggar och ingen epilog.

## Ändra innehåll

Allt innehåll ligger i mappen `data/`. Filerna läses när servern startar,
så starta om servern efter en ändring.

- `texter.json` rum, beskrivningar och allmänna texter
- `fallor.json` regler (liv, chans för händelse, drycker), fällor och sökresultat
- `ormen.json` ormkorten, skatten i kammaren och vad som händer när ormen vaknar
- `monster.json` varelserna (liv, skada, guld, vanor, texter) och striderna
- `ai.json` datorspelarnas namn, personligheter och betänketid
- `avslojande.json` nivågränser, ledtrådar, texter vid nivåbyte, nya
  förklaringar till fällorna och epilogen. Rumsbeskrivningarna för varje
  nivå (`"1"` till `"4"`) ligger i `texter.json`.

För varje fälla kan du ändra namn, vikt (hur ofta den dyker upp), skada,
antal rundor, tecknet på kartan och texterna. Vad fällan gör styrs av dess
`id` i `events.go`. I texterna kan du använda
`{namn}`, `{horn}`, `{riktning}`, `{region}` och `{rum}` där de passar.
Beskrivningarna för rummen ligger under nivåerna `"1"` till `"4"`.

## Struktur

- `main.go` startar webbservern och skriver ut adressen
- `game.go` all spellogik (tornet, rummen, spelarna)
- `events.go` fällor, sökande och drycker
- `serpent.go` skattkammaren och ormkorten
- `monster.go` varelser, strider och den döda spelarens styrning
- `reveal.go` avslöjandet: nivåer, ledtrådar och epilog
- `ai.go` datorstyrda vandrare
- `data.go` läser in datafilerna
- `ws.go` en liten websocket-implementation utan beroenden
- `web/` klienten: HTML, CSS och JavaScript utan ramverk
  - `map.js` ritar kartan med ASCII- och braille-tecken
  - `audio.js` allt ljud, skapat med Web Audio API

## Testa lokalt

Starta servern och öppna `http://localhost:8080` i två flikar med olika
namn. För ett snabbt test passar `go run . -sol 8 -turtid 10s`.

Spellogiken har tester:

```bash
go test ./...
```
