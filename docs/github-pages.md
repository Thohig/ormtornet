# Förslag: Ormtornet på GitHub Pages

Status: förslag, skrivet 2026-11-21. Steg 0 och 1 är byggda.

Målet är att man ska kunna spela Ormtornet från en länk, utan att
installera något och även när vandrarna sitter på olika nät. Tornet som
Go-program på tornets dator finns kvar som i dag; Pages blir ett andra
sätt att öppna tornet, byggt av samma kod.

## Så fungerar det

GitHub Pages lämnar bara ut filer. Därför körs tornet i spelledarens
webbläsare i stället för på en server.

- **Tornets flik.** Spelledaren öppnar sidan på Pages och väljer att öppna
  tornet. Fliken laddar Go-koden kompilerad till WebAssembly
  (`torn.wasm`). Det är samma spelregler, samma synlighet och samma
  `data/` som i dag. Värdsidan och spelledarbordet visas i samma flik.
- **Spelarna** öppnar en länk eller QR-kod från värdsidan. Spelarsidan,
  bilderna och 3D-kartan hämtas från Pages. Sedan ansluter spelarens
  webbläsare direkt till tornets flik med WebRTC, och allt under
  vandringen går den vägen.
- **Signaleringen.** För att två webbläsare ska hitta varandra byter de
  några rader adressuppgifter i början. Det gör biblioteket Trystero via
  publika Nostr-reläer, med en lång slumpad rumskod och ett lösenord som
  krypterar utbytet. Reläerna ser aldrig spelet.
- **Bara tornets flik** har spellogiken och hela kartan. Spelarna får, som
  i dag, bara det tornet skickar till dem.

### Länken

`https://<namn>.github.io/ormtornet/#torn=<kod>&nyckel=<lösenord>`

Koden och lösenordet står efter `#` och skickas därför aldrig till GitHub.
Den som har länken kommer in i lobbyn, precis som den som i dag når
tornets adress på nätet.

### Tornets flik är ett skal

Om fliken byter sida stängs tornet. Därför är tornets flik en skalsida
(`web/torn.html`) som kör `torn.wasm` och visar värdsidan, spelledarbordet
eller den egna spelarsidan i en ram (`iframe`). Ramen kan byta sida fritt;
skalet och tornet ligger kvar. Ramen pratar med skalet med `postMessage`.

Spelledaren kan spela själv i samma flik, eller öppna länken i en annan
flik eller på telefonen som vilken spelare som helst.

## Ändringarna i koden

### Go

- **Byggtaggar.** `main.go`, `ws.go`, öppnandet av webbläsaren och
  liknande får `//go:build !js`. En ny `wasm.go` (`//go:build js`) är
  ingången i webbläsaren. Go-programmet byggs som förut.
- **Anslutningen blir ett gränssnitt.** `Game.serve` tar i dag en
  `*wsConn`. Den får i stället ta något som har `ReadMessage`,
  `WriteText` och `Close`. `wsConn` uppfyller det redan; i wasm finns en
  anslutning för varje WebRTC-kanal och en för ramen i skalet.
- **Förfrågningarna går genom samma `mux`.** `/intro`, `/texter/...` och
  alla `/spelledare/...` hanteras i dag av en `http.ServeMux`. I wasm
  byggs samma `mux`, och skalet skickar förfrågningar till den i minnet
  (`mux.ServeHTTP` med en egen `ResponseWriter`). Handlarna i `host.go`
  behöver därför knappt ändras.
- **Vem är spelledaren.** `isLocal` i `host.go` tittar i dag på
  nätverkskorten. I wasm är det skalets egen ram som är lokal, och allt
  som kommer över WebRTC är det inte. `localOnly` och `hostAction` fungerar
  då som förut.
- **Inställningarna** sparas i dag i en fil (`settings.go`). I wasm sparas
  de i webbläsarens `localStorage`.
- **Omladdningen av `data/`** (`reload.go`) finns inte i wasm. Innehållet
  följer med bygget; man ändrar i repot och Pages bygger om.

### Webben

- **Ett lager för transporten** (`web/natet.js`, byggt från `klient/` med
  Vite tillsammans med Trystero). Det ger spelarsidan något som ser ut som
  en WebSocket och en `fetch`, men som går över WebRTC till tornet, eller
  med `postMessage` till skalet när sidan ligger i tornets flik.
- **`web/app.js`, `web/bord.js` och `web/texter.js`** använder lagret när
  sidan kommer från Pages, och `/ws` och `fetch` som i dag när den kommer
  från Go-programmet. Resten av skripten märker ingen skillnad.
- **Värdsidan** får under Bjud in länken och QR-koden med rumskoden i
  stället för tornets adress på nätet.
- **Nya texter**, till exempel "Kom inte fram till tornet. Prova ett annat
  nät." och "Håll den här fliken öppen, annars stängs tornet.", läggs i
  `data/spelarsidan.json` och `data/vardsidan.json` och i `data/en/`.

### Publiceringen

En GitHub Action bygger `torn.wasm` (`GOOS=js GOARCH=wasm`), kopierar
`wasm_exec.js` från Go, bygger `klient/` och publicerar `web/` på Pages vid
varje push till huvudgrenen.

## Kostnad, konton och integritet

- **Kostnad:** inget, så länge repot är publikt. Pages för privata repon
  kräver en betald plan.
- **Konton:** inga nya. Trystero, Nostr-reläerna och STUN-servrarna är
  gratis och kräver ingen registrering.
- **TURN** (en reläserver när direktanslutning inte går) byggs inte från
  början. Visar provspelen att någon inte kommer in läggs en TURN-nyckel
  till som en inställning; det kräver då ett gratiskonto hos en tjänst.
- **Spelledarens dator:** ingen port öppnas och tornet körs i webbläsarens
  sandlåda. Spelarna kan se spelledarens publika IP-adress, som i ett
  videosamtal.
- **Tornets hemlighet** står i `data/` och går att läsa för den som letar i
  ett publikt repo eller i filerna på Pages. Det går inte att skydda
  tekniskt.

## Risker att prova tidigt

- **Fliken i bakgrunden.** Webbläsare saktar ner timers i flikar som inte
  syns. Tornets tid och robotarna behöver gå jämnt. Prova i steg 1; hjälper
  det inte att fliken har aktiva WebRTC-anslutningar kan tornet köras i en
  Web Worker eller fliken hålla skärmen vaken (Wake Lock).
- **Telefonen som torn.** En låst skärm kan stoppa tornet. Värdsidan ska
  säga att en dator är bäst.
- **Nät som blockerar.** Ungefär var tionde anslutning över internet kan
  behöva TURN, oftast på mobilnät eller nät på jobb och skolor.
- **Storleken.** `torn.wasm` blir några MB, mindre med gzip. Bara
  spelledaren laddar ner den.

## Stegen

Varje steg är en commit, och sedan provar användaren.

0. **Förbered repot för att bli publikt.** Lägg till Three.js licens (Vite
   tar bort den ur `web/iso/iso.js`) och bestäm om spelet ska ha en egen
   licens (`LICENSE`). Användaren gör repot publikt.
1. **Tornet i webbläsaren.** Byggtaggarna, gränssnittet för anslutningen,
   `mux` i minnet, inställningarna i `localStorage`, skalet
   `web/torn.html` och ett skript som bygger `torn.wasm`. Provas lokalt med
   en enkel filserver: spelledaren i skalet och robotar, och en andra flik
   i samma webbläsare via `BroadcastChannel`. Go-testerna ska gå som förut.
2. **Publicering på Pages.** GitHub Action och inställningen för Pages.
   Efter steget går steg 1 att prova på riktiga adressen.
3. **WebRTC.** Trystero, rumskoden och lösenordet, länken och QR-koden på
   värdsidan, och att spelarsidan ansluter från länken.
4. **Det som gör det hållbart.** Spelare som tappar anslutningen kommer
   tillbaka till sin vandrare (`takeover.go`), tydliga texter när något
   inte går, varning om att hålla tornets flik öppen, och åtgärden mot
   bakgrundsfliken om steg 1 visade att den behövs.
5. **Provspel över internet.** Spelledaren hemma, spelare på mobil och
   andra nät. Utifrån det: TURN som inställning eller inte.

README får ett avsnitt "Spela från en länk" och `docs/andringar.md` en
rad efter varje steg.

## Inte med i förslaget

- En egen server för signaleringen eller för spelet.
- Att spelarna ansluter till varandra; alla ansluter bara till tornet.
- Att spara en vandring och fortsätta senare.
