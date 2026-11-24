# Berättelsen om vandringen

Slutkrönikan som en gemensam berättelse med höjdpunkter, i stället för
drag för drag. Förslaget kom 2026-11-24, och användaren valde en gemensam
berättelse för alla, titlar och att vandringens största stunder går upp i
berättelsen.

## Stegen

1. **Berättelsen** (byggt 2026-11-24): inträdet, höjdpunkterna med var
   sitt spelkort och slutet, under resultaten. Hela krönikan och Din
   vandring ligger bakom knappen Visa alla rundor.
2. **Titlarna och raden om var och en** (på gång): under berättelsen en
   rad för varje vandrare med hens bästa stund och hur det slutade, och en
   titel från siffrorna, till exempel den girigaste, den som gick längst
   eller den som gick i flest fällor.

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
- **Flera gånger**: ett slag med `flera` säger en gång hur många gånger
  det hände ("Den natten vaknade ormen tre gånger").
- **Slutet**: vem som bar ut mest (`ber_slut_vinnare`, eller
  `ber_slut_ingen`), vilka andra som kom ut med guld eller tomhänta, vilka
  som var kvar när solen gick ner och vilka som aldrig kom ut. Vinnaren
  får kortet Vinnaren, och när ingen vann får tornet kortet med solen.
  Epilogen står under berättelsen som förut, och tornets hemlighet sägs
  först där.
- **Inget att berätta**: utan höjdpunkter står `ber_tyst`.

Att gå ut och att bli kvar blir inga höjdpunkter; de står i slutet. Ett
validatortest kräver att varje nyckel i `spelkort.berattelse` har sin
berättelsetext, och språktestet att allt finns på engelska.

## Att prova

- Väljer vikterna rätt höjdpunkter? Klockan och portarna kan ta platsen
  från något som kändes viktigare.
- Är sju höjdpunkter lagom, eller blir det för långt?
- Läses "Två rundor senare …" bra flera gånger i rad?
