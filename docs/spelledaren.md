# Förslag: Spelledaren

Status: förslaget skrevs 2026-11-18. Det mesta av del 1–3 byggdes, se
läget nedan. Spelledarläget med översikten och ingreppen (del 2–3) togs
bort 2026-11-22. Resten av dokumentet är det ursprungliga förslaget.

## Läget

Byggt (README beskriver hur det fungerar):

- Värdsidan på `/spelledare`, bara från tornets egen dator.
- Inställningar som ändras utan omstart: svårighet, AI-spelare, språket, avslöjandet och frånvarande vandrare. De sparas
  till nästa gång. ("Vandrarna syns" togs bort 2026-11-21: egen karta är
  det enda läget. Minsta antal, sol, turtid, strid mellan vandrare, tornets
  storlek, klasser, invigningen, tornet förändras och kartan togs bort
  2026-11-22: de bestäms automatiskt eller gäller alltid.)
- Starta nu, Stäng tornet, och inbjudan att kopiera till en gruppchatt.
- AI tar över en frånvarande vandrare, och sedan 2026-11-22 också en som
  står still utan att välja drag.
- Spelledaren spelar i samma flik (2026-11-22): Gå in själv byter till
  spelarsidan, som visar märket spelledare och knappen Värdsidan i lobbyn
  och slutet. Adressen under Bjud in visas bara med Visa adressen.
- `Ormtornet.app` via `./skapa-app.sh`. Den bygger från repot och kräver Go.
- Värdsidan på GitHub Pages (2026-11-21, [github-pages.md](github-pages.md)):
  den visas i tornets flik, skalet, och har under Bjud in en länk i
  stället för adressen på samma wifi (QR-koden togs bort 2026-11-22), en varning om att hålla
  fliken öppen och egna råd under Kommer någon inte in?.
- Pekkontroller, så att telefoner och surfplattor kan spela.
- Krönikan: en sammanfattning efter varje omgång.
- Förval för svårighet: Lätt, Normal, Svår och Mardröm.
- Tornet pausar rundan själv när ingen vandrare är ansluten.
- Läs om innehållet i `data/` utan omstart, med felen på värdsidan om
  filerna inte går ihop. Vandringens frö visas under Vid porten.
- Värdsidans avsnitt står alltid i samma ordning. Under vandringen visar
  Vid porten bara raden om vandringen och vandrarna (2026-11-21).

Borttaget 2026-11-22: spelledarläget, spelledarbordet på
`/spelledare/bord`, översikten (hela tornet, vandrarna, ormens sömn och
risk, viskningarna och loggen) och ingreppen (tid, tornets röst, fällor
och varelser, ormen, avslöjandets nivå, liv, guld, drycker, föremål och
invigningen). Spelledaren spelar nu alltid med eller låter bli.

På gång: inget just nu.

Inte byggt, och inte planerat just nu:

- Fler enskilda inställningar (chans för fällor, varelser och dörrar, frö,
  liv, ormkort, nivågränser, personlighet för varje AI-spelare), namngivna
  förval och paus mellan rundor.
- Statistik över flera kvällar och redigering av datafilerna i webbläsaren.
- Kod för att nå sidan från andra datorer.
- En färdigbyggd app som går att dela utan Go.

Spelledaren är en egen sida för den som leder spelet. Den ligger i samma
Go-server på `/spelledare`, inte i ett separat program. Då behövs ingen extra
installation, inställningar kan ändras medan servern kör, och spelledaren
delar tillstånd med spelet direkt. Sidan har samma gamla skärmkänsla som
spelet.

## Del 1: inställningar och lobby

Dagens flaggor blir ett formulär som kan ändras i lobbyn utan omstart.

| Grupp | Inställningar |
|---|---|
| Tid | solens längd, turtid, paus mellan rundor |
| Spelare | minsta antal, antal datorspelare och personlighet för var och en |
| Svårighet | förval (Lätt, Normal, Svår, Mardröm) som sätter värdena nedan |
| Tornet | chans för fälla, chans för varelse, chans för dörrar, frö för samma torn igen |
| Vandrare | liv, hur mycket en dryck läker, skatt i kammaren |
| Ormen | antal sovkort och vakna-kort, skada, andel skatt man förlorar |
| Avslöjandet | på eller av, nivågränser, hur ofta ledtrådar delas ut |

Inställningarna kan sparas som namngivna förval i `data/installningar.json`.
Sidan har knapparna **Starta** och **Nytt torn**.

## Del 2: översikt under omgången

- Hela kartan, även okända rum, med fällor och varelser utmärkta.
- Varje spelares liv, skatt, drycker, tillstånd (gift, mörker, källare) och
  valda drag.
- Ormkortsleken: kort kvar och risken att ormen vaknar vid nästa dragning.
- Avslöjandet: nivå och vilka ledtrådar varje spelare har fått.
- Hela loggen, med allt som har skickats till varje spelare.
- Tempo: paus, mer tid åt rundan, utför rundan nu, ändra solen.

## Del 3: ingrepp

- Tornets röst: skicka en text till en spelare eller till alla.
- Placera eller utlös fällor och varelser, slå klockan, väck ormen eller låt
  den röra sig.
- Ge eller ta liv, guld och drycker, flytta en vandrare, byt ut en spelare som
  har lämnat mot en datorspelare.
- Styr avslöjandet: gå vidare till nästa nivå eller håll kvar nivån.
- Spela upp ett ljud för alla.

Alla ingrepp skrivs in i en lista som bara spelledaren ser.

## Del 4: historik och redigering

- Sammanfattning per omgång: vem som tog vad, fällor, strider, när ormen
  vaknade.
- Statistik över flera kvällar.
- Redigering av datafilerna i webbläsaren, med kontroll av formatet och
  omladdning till nästa omgång utan omstart.

## Åtkomst

- Fri från Macen själv (`localhost`).
- Från andra datorer krävs en kod som servern skriver ut när den startar.
- Bara en spelledare åt gången. Spelledaren är ingen vandrare och syns inte i
  spelet.

## Hur en omgång startar

I dag: alla anslutna människor trycker Enter, och spelet startar när alla är
redo och de är minst `-min`.

Förslag:

- **Med Värdsidan öppen** startar spelledaren omgången. Lobbyn visar
  "Väntar på spelledaren …". Enter betyder bara "jag är här och klar".
  Spelledaren kan starta även om någon inte är redo, ta bort spelare och fylla
  tomma platser med datorspelare. Efter omgången trycker spelledaren
  **Nytt torn**.
- **Utan spelledare** fungerar det som i dag.
- Den första spelaren ska inte bli värd: rollen blir otydlig, den måste flyttas
  om spelaren tappar förbindelsen, och spelledaren ska se hela tornet.
- Möjligt mellanläge: vem som helst trycker **S** för att starta när minst
  `-min` spelare är redo.

## Hur servern startas

I dag: `go run .` i terminalen på en Mac som har Go och repot.

Förslag: spelledaren kör servern på sin egen Mac genom en färdigbyggd
**Ormtornet.app**.

- Allt är inbakat; ingen Go behövs. Appen kan delas med AirDrop.
- Dubbelklick startar servern, öppnar `/spelledare`, visar adressen till
  spelarna i stor text och håller Macen vaken.
- Inställningarna görs på sidan. Flaggorna finns kvar för terminalen.
- Avsluta med knappen **Stäng tornet** (spelarna får ett meddelande) eller ⌘Q.
- Första gången frågar macOS om appen får ta emot nätverksanslutningar
  (svara Tillåt). Appen är inte signerad, så den måste öppnas med högerklick
  och **Öppna** första gången.

## Föreslagen byggordning

1. Ormtornet.app, inställningar och startknapp (del 1)
2. Översikt (del 2)
3. Ingrepp (del 3)
4. Historik och redigering (del 4)

En commit per del, med paus för test mellan delarna.
