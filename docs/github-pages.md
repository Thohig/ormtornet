# Ormtornet på GitHub Pages

Status: skrivet som förslag 2026-11-21. Steg 0–4 är byggda och
publicerade på <https://thohig.github.io/ormtornet/>. Kvar är steg 5, ett
provspel över internet. README beskriver under Spela från en länk hur det
används.

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
  `data/` som i dag. Värdsidan visas i samma flik.
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
(`web/torn.html`) som kör `torn.wasm` och visar värdsidan
eller den egna spelarsidan i en ram (`iframe`). Ramen kan byta sida fritt;
skalet och tornet ligger kvar. Ramen och andra flikar i samma webbläsare
pratar med skalet över en `BroadcastChannel`. Själva tornet går i en Web
Worker (`web/torn-arbetare.js`), som skalet för meddelandena vidare till.

Länkar mellan sidorna byter sida i ramen. De får inte öppna en egen flik,
för då hamnar sidan utanför skalet.

Skalets flik heter först "Ormtornet – tornet". När tornet har startat tar
fliken varje sekund titeln från sidan i ramen, till exempel "Ormtornet" på
introt och "Ormtornet – värd" på värdsidan, så att den går att känna igen
bland flikarna. När tornet är stängt står den senaste titeln kvar.

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

- **Ett lager för transporten** (`web/natet.js`). Det byter ut `fetch`
  och `WebSocket` för tornets adresser, så att de går över WebRTC till
  tornet när sidan kommer från länken, eller över en `BroadcastChannel`
  till skalet i spelledarens egen webbläsare. När sidan kommer från
  Go-programmet gör det ingenting, och resten av skripten märker ingen
  skillnad.
- **WebRTC** (`web/rtc/rtc.js`, byggt från `klient/src/rtc.ts` med Vite
  tillsammans med Trystero och en QR-kodsgenerator). Det laddas bara av
  skalet, av värdsidan för QR-koden och av spelarsidan när länken har en
  rumskod.
- **Värdsidan** får under Bjud in länken och QR-koden med rumskoden i
  stället för tornets adress på nätet.
- **Nya texter** för länken, för när tornet inte hittas, varningen om att
  hålla tornets flik öppen och råden under Kommer någon inte in? ligger i
  `data/spelarsidan.json` och `data/vardsidan.json` och i `data/en/`.

### Publiceringen

GitHub Action `Sida` (`.github/workflows/sida.yml`) kör testerna och
`verktyg/bygg-sida.sh`, som skriver sidorna (`go run . -sida`), bygger
`torn.wasm` (`GOOS=js GOARCH=wasm`) och kopierar `wasm_exec.js` från Go.
Resultatet publiceras på Pages vid varje push till `main`. Skripten från
`klient/` är incheckade, så Action behöver inte Node. Sidornas filer får
ett versionsnummer (`?v=` och början av commiten), eftersom Pages låter
webbläsaren spara filerna i tio minuter.

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
  syns. Tornets tid och robotarna behöver gå jämnt. Löst i steg 4 genom att
  köra tornet i en Web Worker, men inte provat med en dold flik i mer än
  fem minuter.
- **Telefonen som torn.** En låst skärm kan stoppa tornet. Värdsidan ska
  säga att en dator är bäst.
- **Nät som blockerar.** Ungefär var tionde anslutning över internet kan
  behöva TURN, oftast på mobilnät eller nät på jobb och skolor.
- **Storleken.** `torn.wasm` är ungefär 9 MB, mindre med gzip. Bara
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
   Byggt: `.github/workflows/sida.yml` kör testerna och
   `verktyg/bygg-sida.sh` och publicerar vid push till `main`. Användaren
   gör repot publikt och väljer GitHub Actions som källa under Settings,
   Pages.
3. **WebRTC.** Trystero, rumskoden och lösenordet, länken och QR-koden på
   värdsidan, och att spelarsidan ansluter från länken.
   Byggt: `klient/src/rtc.ts` byggs till `web/rtc/rtc.js` (Trystero 0.27
   med Nostr, och en QR-kodsgenerator). Skalet går in i rummet med en
   rumskod och en nyckel på 20 tecken var, sparade i `sessionStorage`, och
   `web/natet.js` går in från länken. Förfrågningar som kommer över WebRTC
   räknas inte som spelledarens (`Tornet.fraga` har ett sjätte argument,
   `lokal`). Spelaren ser "Letar efter tornet …" efter tre sekunder och ett
   besked efter 30 om tornet inte hittas. Provat lokalt med två flikar,
   varav den ena över WebRTC: lobbyn, klassvalet, en runda, att
   spelledarens sidor svarar 403 och att fel nyckel inte kommer in. Inte
   provat mellan olika nät.
4. **Det som gör det hållbart.** Spelare som tappar anslutningen kommer
   tillbaka till sin vandrare (`takeover.go`), tydliga texter när något
   inte går, varning om att hålla tornets flik öppen, och åtgärden mot
   bakgrundsfliken om steg 1 visade att den behövs.
   Byggt: tornet körs i en Web Worker (`web/torn-arbetare.js`), eftersom
   "Gå in själv" lägger tornets flik i bakgrunden, där Chrome efter fem
   minuter kör timers en gång i minuten när ingen WebRTC är igång. Pulsen
   till sidorna kommer också från workern. Inställningarna sparas i
   skalets `localStorage` genom workern. Värdsidan säger på Pages att
   tornets flik ska hållas öppen och att en dator är bäst, och Kommer någon
   inte in? har egna råd för länken. Den som laddar om eller öppnar länken
   igen och skriver samma namn får tillbaka sin vandrare; det fanns redan
   och är provat över WebRTC. Bakgrundsfliken är inte provad, eftersom
   webbläsarpanelen håller alla flikar synliga. Efteråt: ett
   versionsnummer (`?v=`) på sidornas filer i bygget (`sida.go`), eftersom
   Pages låter webbläsaren spara filerna i tio minuter.
5. **Provspel över internet.** Spelledaren hemma, spelare på mobil och
   andra nät. Utifrån det: TURN som inställning eller inte. Se också efter
   att tornets flik håller när den har legat i bakgrunden mer än fem
   minuter, och att den som tappar anslutningen kommer tillbaka.

README har avsnittet "Spela från en länk", och `docs/andringar.md` en
rad för varje steg.

## Inte med i förslaget

- En egen server för signaleringen eller för spelet.
- Att spelarna pratar med varandra; alla pratar bara med tornet.
  Trystero kopplar ändå ihop alla i rummet med varandra, men spelarsidan
  lyssnar bara på tornets flik. Med åtta spelare blir det många
  anslutningar; märks det i provspelet kan tornet få en egen roll i
  rummet.
- Att spara en vandring och fortsätta senare.
