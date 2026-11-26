# Berättelsen om vandringen

Slutkrönikan som en gemensam berättelse med höjdpunkter, i stället för
drag för drag. Förslaget kom 2026-11-24, och användaren valde en gemensam
berättelse för alla, titlar och att vandringens största stunder går upp i
berättelsen.

## Stegen

1. **Berättelsen** (byggt 2026-11-24): inträdet, höjdpunkterna med var
   sitt spelkort och slutet, under resultaten. Hela krönikan och Din
   vandring ligger bakom knappen Visa alla rundor.
2. **Titlarna och raden om var och en** (byggt 2026-11-24): under
   berättelsen en rad för varje vandrare, i placeringens ordning, med
   namnet i sin färg, en titel från siffrorna, hens största stund och hur
   det slutade. Sedan guldboken (docs/guldboken.md) står titeln och den
   största stunden på vandrarens rad i placeringen överst på
   slutskärmen, och hela krönikan står inte längre där.
3. **Efter provspelet 2026-11-26** (byggt 2026-11-26,
   [provspel-2026-11-26.md](provspel-2026-11-26.md)): guldet hos den som
   blev kvar när solen gick ner står i slutet och på hens rad, den som
   ensam tog mest ur ormens hög blir en höjdpunkt, och den dödes guld
   säger att hen hade fallit när döden inte är berättad.

## Hur den byggs

`buildStory` i `berattelse.go` skriver berättelsen en gång, när
vandringen är slut, så att texterna inte byter variant varje gång
tillståndet skickas. Den skickas som `berattelse` i tillståndet och ritas
av `renderStory` i `web/app.js`.

- **Inträdet** (`ber_in`): hur många som gick in, var och en med sin
  klass (`ber_in_vem`, klassens `kort`), och hur många rundor vandringen
  varade.
- **Höjdpunkterna**: krönikans händelser (`g.chron`) sparar nu sina
  variabler (vem, offer, guld, varelse …), så att berättelsen kan säga
  samma sak med andra ord. Varje händelse väger sitt slags `vikt`, mer
  (`forst`) om den är den första av sitt slag och mer (`vinnare`) om den
  handlar om den som vann. De tyngsta tas med, högst `per_slag` av varje
  slag och `antal` i allt, och berättas i den ordning de hände. Allt står
  under `spelkort.berattelse` i `data/texter.json`.
- **Texterna**: varje krönikenyckel har en berättelsetext, `ber_` och
  nyckeln utan `kronika_` (till exempel `ber_vaknade`), som börjar med
  `{nar}`. Därför står verbet efter: "I runda 7 vaknade ormen", "Rundan
  efter föll Ylva för …". `{nar}` är `ber_nar_runda` för den första
  höjdpunkten och sedan räknat från den förra: `ber_nar_samma`,
  `ber_nar_efter`, `ber_nar_senare` (med talet i ord, `tal_2` till
  `tal_10`), `ber_nar_langt` efter mer än tio rundor och `ber_nar_sist`
  för den sista rundan. Av en texts varianter tas den som inte redan är
  sagd, så att två uppvaknanden inte berättas med samma ord.
- **Ormens hög** (2026-11-26): den som ensam tog mest ur kammaren under
  vandringen blir en händelse i krönikan, `kronika_hogen`, i rundan då
  hen tog sin sista näve (`noteHoard`). Den väger som ett slag mellan
  vandrare och tas med också när hen aldrig kom ut: "I runda 5 hade Alva
  tagit 157 guld ur ormens hög, mer än någon annan."
- **Den dödes guld** (2026-11-26): tar någon upp en dödes guld innan
  döden är berättad står `ber_guld_kvar_fallen`, "kom Bo till platsen där
  Cecilia hade fallit, och tog upp guldet", så att läsaren vet att hen
  hade dött.
- **Flera gånger**: ett slag med `flera` säger en gång hur många gånger
  det hände ("Den natten vaknade ormen tre gånger").
- **Slutet**: vem som bar ut mest (`ber_slut_vinnare`, eller
  `ber_slut_ingen`), vilka andra som kom ut med guld eller tomhänta, vilka
  som var kvar när solen gick ner, var och en med sitt guld när hen bar
  något (`ber_slut_kvar_guld`), och vilka som aldrig kom ut. Vinnaren
  får kortet Vinnaren, och när ingen vann får tornet kortet med solen.
  Epilogen står under berättelsen som förut, och tornets hemlighet sägs
  först där.
- **Inget att berätta**: utan höjdpunkter står `ber_tyst`.

## Raden om var och en

`storyPersons` skriver en rad för varje vandrare: "Ingrid, jägaren.
Fällde en grönaktig skepnad. Kom ut med 35 guld." Den skickas som
`vandrarna` i tillståndet.

- **Titeln**: titlarna under `spelkort.berattelse.titlar` delas ut i
  listans ordning. Var och en går till den som ensam har mest av sitt mått
  (`matt`) och inte redan har en titel; lika mycket ger ingen titel, och
  ingen får två. Måtten är guld ur kammaren (den girigaste), stulet guld
  (rövaren), fällda varelser (jägaren), slag på vandrare (slagskämpen),
  rum på den egna kartan (upptäckaren), föremål (samlaren), fällor
  (fällornas vän), sökta rum (letaren) och dragna ormkort (`titleMeasures`
  i `berattelse.go`). Texterna heter `ber_titel_*`.
- **Största stunden**: den tyngsta händelse i krönikan som handlar om
  vandraren, med samma vikter som höjdpunkterna. Namnet först i texten tas
  bort, eftersom raden redan börjar med det.
- **Hur det slutade**: `ber_vem_ute`, `ber_vem_tom`, `ber_vem_kvar`
  (`ber_vem_kvar_guld` med guldet hen bar) eller `ber_vem_dod`, utom när
  den största stunden var att dö.

Att gå ut och att bli kvar blir inga höjdpunkter; de står i slutet. Ett
validatortest kräver att varje nyckel i `spelkort.berattelse` har sin
berättelsetext, och språktestet att allt finns på engelska.

## Att prova

- Väljer vikterna rätt höjdpunkter? Klockan och portarna kan ta platsen
  från något som kändes viktigare.
- Är sju höjdpunkter lagom, eller blir det för långt?
- Läses "Två rundor senare …" bra flera gånger i rad?
- Är titlarna roliga, eller blir de tjatiga när samma titel kommer varje
  kväll? Får för få en titel när lika mycket inte ger någon?
