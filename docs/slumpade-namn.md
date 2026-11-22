# Slumpade namn (förslag)

Förslag 2026-11-22. Steg 1 och 2 är byggda, med förnamn och tillnamn och
samma namn hela kvällen. Steg 3, Nytt namn i lobbyn, är inte byggt.
Robotarna skickar fortfarande sina egna namn, så att deras loggar går att
läsa.

## Vad som önskas

Spelaren skriver inte längre sitt namn. Tornet ger var och en ett
slumpmässigt namn som passar spelet.

## Så här ser det ut för spelaren

- **Namnrutan försvinner.** "Gå in i tornet" i introt leder direkt till
  lobbyn.
- **Tornet ger namnet.** Den första raden i loggen säger det, till exempel:
  "Tornet kallar dig Märta Vindögd." Namnet står också vid din prick i
  översta raden, som i dag.
- **Samma namn hela kvällen.** Sidan kommer ihåg namnet. Laddar du om
  sidan, tappar kontakten eller kommer tillbaka till nästa torn, får du
  samma namn igen och tar tillbaka din vandrare, precis som i dag.
- **Byta namn** (valfritt, se frågorna): en liten knapp **Nytt namn** i
  lobbyn slumpar ett nytt. Under vandringen går det inte.

## Namnen

- Ett förnamn och ett tillnamn, i stil med tornet: "Märta Vindögd", "Tore
  Askhand", "Ylva från Kärret", "Folke Tredje Sonen". Förnamnen är gamla
  nordiska namn, som datorspelarnas. Tillnamnen ger lite karaktär utan att
  säga något om tornets hemlighet.
- Listorna ligger i en ny fil, `data/namn.json`, med `fornamn` och
  `tillnamn`. På engelska (`data/en/namn.json`) är förnamnen desamma och
  tillnamnen översatta ("Märta Windeye", "Tore Ashhand"). Det går att lägga
  till och ändra namn utan att röra koden.
- **Två vandrare får aldrig samma namn**, inte heller samma förnamn, så
  att det går att skilja dem åt i loggen. Datorspelarna tar sina namn ur
  samma förnamn och undviks på samma sätt.
- Namnen får vara högst 16 tecken långa, som i dag, så att de ryms i
  översta raden och loggen. Ett test kontrollerar det.

## Hur det fungerar

- Spelarsidan skickar det namn den kommer ihåg, eller inget. Servern
  (`Join` i `game.go`) slumpar ett namn när inget kommer, och skickar
  tillbaka vilket det blev. Sidan sparar det som i dag (`ormtornet-namn`).
- Är det sparade namnet redan upptaget av någon som är ansluten, till
  exempel när två flikar i samma webbläsare öppnar tornet, får den nya
  fliken ett eget namn i stället för felet "namnet är upptaget".
- Spelledaren som går in själv från värdsidan får också ett slumpat namn.
- Robotarna i `verktyg/provspelare.mjs` skickar inget namn och får ett.
- Texterna `namn_fraga`, `namn_tips`, `namn_tillbaka`, `namn_ga_in` och
  `namn_saknas` tas bort, och `tomt_namn` på servern. Nya texter kommer på
  båda språken.

## Steg

1. **Namnen och servern.** `data/namn.json` på båda språken, slumpningen i
   `Join` och test för att namnen är unika och korta.
2. **Spelarsidan utan namnruta.** Introt leder direkt till lobbyn, och
   loggen säger vad tornet kallar dig. Sedan uppdateras värdsidan,
   robotarna och dokumentationen.
3. **Nytt namn** i lobbyn, om du vill ha det.

## Frågor

- Ska namnen vara förnamn och tillnamn ("Märta Vindögd"), eller bara
  förnamn som datorspelarnas ("Märta")?
- Ska man kunna slumpa ett nytt namn i lobbyn, och i så fall hur många
  gånger?
- Ska man få behålla namnet till nästa torn samma kväll, eller få ett nytt
  namn för varje vandring?
