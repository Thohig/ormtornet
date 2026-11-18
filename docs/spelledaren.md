# Förslag: Spelledaren

Status: förslaget skrevs 2026-11-18. Det mesta av del 1–3 är byggt, se
läget nedan. Resten av dokumentet är det ursprungliga förslaget.

## Läget

Byggt (README beskriver hur det fungerar):

- Värdsidan på `/spelledare`, bara från tornets egen dator.
- Inställningar som ändras utan omstart: AI-spelare, minsta antal, sol,
  turtid, avslöjandet, strid mellan vandrare, frånvarande vandrare,
  vandrarna syns, tornets storlek, klasser och spelledarläget. De sparas
  till nästa gång.
- Starta nu, Stäng tornet, och inbjudan att kopiera till en gruppchatt.
- Översikten i spelledarläget: hela tornet, vandrarna, varelserna, hur djupt
  ormen sover, avslöjandets nivå och allt vandrarna har läst.
- Ingrepp: tid, tornets röst, fällor och varelser, väcka och lugna ormen,
  liv, guld och drycker, och AI tar över en frånvarande vandrare.
- `Ormtornet.app` via `./skapa-app.sh`. Den bygger från repot och kräver Go.
- Pekkontroller, så att telefoner och surfplattor kan spela.
- Krönikan: en sammanfattning efter varje omgång.
- Förval för svårighet: Lätt, Normal, Svår och Mardröm.
- Risken att ormen vaknar i översikten, och att flytta fram eller hålla kvar
  avslöjandets nivå.
- I översikten: klass, föremål, viskningar och vem som är invigd. Ingrepp:
  ge föremål, och inviga eller ta tillbaka invigningen.

På gång: inget just nu.

Inte byggt, och inte planerat just nu:

- Fler enskilda inställningar (chans för fällor, varelser och dörrar, frö,
  liv, ormkort, nivågränser, personlighet för varje AI-spelare), namngivna
  förval och paus mellan rundor.
- Ingrepp: flytta en vandrare, slå klockan, låta ormen röra sig, spela upp
  ett ljud.
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
