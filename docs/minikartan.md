# Minikartan

En liten karta med allt vandraren har upptäckt, byggd 2026-11-24. Förslaget
hade tre varianter: A, en fast ruta på kartan; B, kompassen blir
minikartan; och C, en karta som bara syns när man ber om den. A byggdes för
datorn och B för telefonen.

## Vad den visar

Bara det vandraren själv vet, alltså samma rum som kartan (`state.rooms`).
Den avslöjar aldrig mer än zoomen Allt.

- Rummen som rutor, och kammaren i guld när man har hittat den.
- Dörrarna som streck från mitten ut över kanten. Ett streck som inte möter
  något är en dörr till ett rum man inte har varit i.
- Portarna som ett streck ut med en tvärslå, röda när porten har slutit
  sig.
- Vandraren själv som en prick i sin färg. Andra vandrare syns inte; man
  ser dem bara när man står i samma rum, och då står de på kartan.
- En streckad ram för det kartan visar just nu (inte när minikartan är
  stor).
- Kartritaren ser de rum andra har hittat streckade (`a` i rummet, från
  `stateFor` i `game.go`) och kryss där någon har dött (`dod`).
- Avslöjandet når minikartan: kritan rodnar på nivå 3 och 4.

Den visas bara under vandringen när man är inne i tornet.

## Hur den ritas

`web/minikarta.js` ritar en SVG som `MapView.draw` i `web/map.js` uppdaterar
varje gång kartan får ett nytt tillstånd eller en ny zoom. Rummen ritas
platta och vridna: en punkt (x, y) står på (x − y) · a åt höger och
(x + y) · a · 0,55 nedåt, så att öster pekar snett nedåt höger och norr
snett uppåt höger, som på kartan snett ovanifrån och tangenterna vid
dörrarna. Ett N med en pil står i hörnet. Skalan räcker för allt man känner
till, med lite luft och aldrig närmare än 3 × 3 rum.

## Var den står

- **Datorn** (bredare än 700 punkter): uppe till vänster på kartan, 168 ×
  112 punkter (124 × 80 på en låg skärm). Raden om källarvalven och loggen
  på kartan börjar under den. Ett klick växlar mellan zoomen man har och
  Allt (`MapView.toggleWhole`). K döljer och visar den, och valet sparas i
  webbläsaren (`ormtornet.minikarta`). Kompassen står kvar ovanför zoomen.
- **Telefonen**: nere till vänster i kompassens ställe, som döljs. Ett
  tryck gör minikartan stor över kartan och ett tryck till gör den liten
  igen; Escape stänger den också. Den ligger under tangenterna vid
  dörrarna, så att de alltid går att trycka på.

## Att prova

- Är minikartan till nytta på datorn, eller räcker zoomen Allt? Står den i
  vägen för loggen när en runda har mycket text?
- Går den att läsa på telefonen när den är liten, eller behövs den bara
  stor?
- Hittar man hem till en port med den, och lockar strecken mot okända rum
  till att utforska?
