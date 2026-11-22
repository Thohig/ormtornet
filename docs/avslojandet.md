# Förslag: avslöjandet och invigningen

Status: förslaget skrevs 2026-11-19 och är byggt i fem steg samma dag.
Sedan 2026-11-22 finns inte spelledarbordet, och inte heller fasta svar
eller inställningarna för invigningen och tornet som förändras; det som
står om dem nedan är hur det var byggt.
README beskriver hur det fungerar (spoilervarning). Med 600 simulerade
vandringar ändrade portarna som sluter sig andelen som blev kvar i mörkret
med högst två procentenheter. Invigningen, bytet och nåden kan inte
simuleras, eftersom datorspelarna inte viskar, och behöver provas i spel.

## Varför

Invigningen var ett pussel som bara gick att lösa en gång: svaret var
detsamma varje vandring, och epilogen sa det rakt ut. Från andra kvällen
viskade alla samma ord så fort det gick. Avslöjandet var ren berättelse:
inget i tornet ändrades av att man förstod.

Målet är att avslöjandet ska vara berättelsen om vad tornet är, och samtidigt
spela roll för vad vandrarna gör, även andra kvällen.

## Inställningarna

Sedan 2026-11-22 finns de två raderna inte längre: med avslöjandet på
gäller alltid invigningen med ormens namn och tornet som förändras, och
fasta svar är borttaget. Så här var det byggt:

Under Fler inställningar, efter Avslöjandet:

| Rad | Val | Förval |
|---|---|---|
| Avslöjandet (finns) | på, av | på |
| Invigningen | ormens namn, fasta svar | ormens namn |
| Tornet förändras | på, av | på |

Med avslöjandet av är de två nya raderna avstängda ("kräver avslöjandet").
Fälten heter `invigning` (`namn` eller `svar`) och `forandras`, flaggorna
`-invigning` och `-forandras`.

## Del 1: invigningen med ormens namn

- Varje vandring lottar ett namn på tre stavelser ur `namn` i
  `avslojande.json`, med tornets slump (samma frö, samma namn).
- Vid nivå 2 får varje människa i tornet en del, så att gruppen har så många
  olika delar som möjligt. Vid nivå 3 och 4 får var och en en del till som
  den saknar. Den som spelar ensam har hela namnet vid nivå 4.
- Texterna (`delar`) följer nivåerna, från "tre tecken som bildar ett ord"
  till "det är tornet självt som heter så".
- Tecknen på spelskärmen visar delarna, till exempel `· RU ·`.
- Viskningen ska innehålla namnet och ett av de vanliga svaren, eller en
  mall ur `svar_med_namn` ("tornet är {namn}"). De vanliga svaren räcker
  inte ensamma. Ett försök per nivå, som förut. Bara namnet förbrukar inte
  försöket: tornet känner igen det och frågar vad man vill säga om det
  (2026-11-22).
- Byte vid dödande: den som dödar en vandrare får de delar offret hade och
  dödaren saknade (bara med strid mellan vandrare på).
- Tornets nåd: när alla människor som är kvar inne är invigda tas hälften av
  vakna-korten som är kvar bort (`nad_andel`, avrundat nedåt, minst ett),
  en gång per vandring. Nåden håller när leken byggs om.
- Fasta svar fungerar som förut, utan namn, byte och nåd.

## Del 2: varslen

- Från nivå 2 är varannan ledtråd ett varsel: tornet förråder ett rum som
  ingen har varit i, två eller tre steg bort, till exempel "Det gnisslar av
  järn i norr, två rum bort".
- Rummets innehåll bestäms när varslet ges och blir så när någon går in.
  Det kan vara en fälla, en varelse eller ingenting.
- Varslet säger aldrig var kammaren ligger; den riktningen är de invigdas.

## Del 3: tornet förändras

- Nivå 3: fällor som hör till ormens kropp blir vanligare i nya rum
  (`forandring.fallvikter`).
- Nivå 4: portarna stängs en i taget (`forandring.port_stangs_var` rundor),
  först den som ligger längst från vandrarna. Den sista stängs aldrig, och
  rundan innan får alla en varning. Datorspelarna söker sig till en öppen.

## Det som följer med

- Epilogen säger namnet, och krönikan berättar om invigningen, bytet,
  nåden och portarna.
- Spelledarbordet visar namnet, delarna, nåden, varslen och portarna.
- Alla texter på svenska och engelska; namn och siffror bara i de svenska
  filerna.

## Byggordning

1. Invigningen: inställningen, namnet, delarna, viskningen och epilogen.
2. Byte vid dödande och tornets nåd, med krönikan.
3. Varslen.
4. Tornet förändras, med datorspelarna och simulering.
5. Spelledarbordet och dokumentationen.
