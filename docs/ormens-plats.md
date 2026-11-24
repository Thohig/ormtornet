# Ormens plats

Förslag 2026-11-24, byggt samma dag i steg. Ormen och kammaren kan ligga
nästan var som helst i tornet, inte bara nära mitten, utan att den som
startar nära ormen får det för lätt. En liten fördel får finnas, men den
ska inte avgöra vandringen.

## Steg 1: en fri och ganska rättvis plats

- Kammaren hamnar inte i den yttre ringen av rum, och inte närmare någon
  port än tre rum (`fri_fran_port` i `data/ormen.json`).
- Ingen vandrare får ha kammaren mer än tre rum närmare än snittet av
  allas avstånd från sina portar (`fri_skillnad`). Avståndet räknas i rum
  (rakt, inte genom dörrarna, som inte finns förrän rummen hittas).
- Platsen väljs när portarna har lottats, efter vandringens frö. När ormen
  vaknar kryper den till en ny plats med samma regler.
- Med `plats` `mitten` ligger kammaren som förut inom `spridning` från
  mitten; testerna använder det.
- Solen: med automatisk sol får vandringen två rundor till för varje rum
  som den längsta vägen från en port till kammaren är längre än vägen till
  mitten.

Att kammaren inte kan vara helt fri och samtidigt rättvis är geometri:
används alla fyra hörnen i ett torn på 9 × 9 är summan av avstånden
densamma överallt, så varje plats utanför mitten ligger närmare någon.
Regeln ovan släpper fram platser ända ut mot väggarna men hindrar att någon
har kammaren alldeles bredvid sig.

## Steg 2: näven hörs

När någon tar en näve hör alla andra i tornet åt vilket håll det kom ifrån
(förut bara att ormen rörde sig). Den som är först i kammaren blir känd
för alla. Tjuvens näve märks som förut inte.

## Steg 3: att hitta ormen

- Står du bredvid kammaren hör du andetagen genom väggen, och åt vilket
  håll, i rumsblocket.
- Från nivå 2 hör alla andetagen en gång, åt vilket håll de kommer.
- Texterna som sa att golvet sluttar mot mitten säger åt vilket håll ormen
  ligger, och introt säger att ormen sover någonstans i tornet.
- Den invigdes hjärtslag, breven och varslen pekar som förut mot kammaren.

## Steg 4: AI-spelarna och balansen

- AI-spelarna letar efter kammaren genom att följa andetagen och utforska,
  inte mot mitten.
- Simuleringen visar hur det går efter avståndet från porten till
  kammaren. Målet är att den som startade närmast vinner högst omkring 5
  procentenheter oftare än snittet.
- Räcker det inte växer guldet i kammaren lite med nivån.
