# Den individuella tävlingen

Förslag 2026-11-24. Ormtornet är en tävling: den som kommer ut med mest
skatt vinner. Det som gynnade andra än en själv tas bort, och tävlingen
skärps så att det märks vem som är före vem. Inget samarbete belönas;
pratet vid bordet blir kvar, men som ett spel där man kan ljuga.

## Steg 1: det som gynnar andra (byggt 2026-11-24)

- **Kartritaren ger inte bort sin karta.** Förut fick den hen mötte
  hennes karta. Förslaget var att hon i stället skulle se den andras
  karta, men hon ser redan alla rum som någon har hittat, så det fanns
  inget att ta. Delningen är borttagen (`shareMap` och texterna
  `klass_kartritare_karta` och `klass_kartritare_ger`).
- **Timglaset gäller bara den som bär det.** Förut stod solen still en
  runda för alla när den hade tre handsbredder kvar. Nu vänds timglaset
  när solen går ner, och porten väntar en runda till på bäraren, som på
  de invigda (`graced` och `portGrace` i `secret.go`).
- **Viskningen i kammaren lugnar ormen bara för den invigde.** Förut tog
  den bort ett vakna-kort ur leken som alla drar ur. Nu blir den invigdes
  nästa näve säker (`safeDraw`), och leken är densamma för de andra.
- **Introt uppmanar inte till samarbete.** "Berätta för varandra vad ni
  har sett" är nu "Ni kan berätta för varandra vad ni har sett, men ingen
  måste säga sanningen" (`synlighet_egen`), och Kartritarens variant "Det
  du vet behöver du inte dela med någon" (`synlighet_kartritare`).

Redan en tävling, och orört: striden mellan vandrare, den dödes guld, att
den som fäller varelsen tar allt guld, ormens namn, den förste invigde
och leken i kammaren, där varje sovkort man tar gör nästa näve farligare
för de andra.

## Steg 2: placeringar och vinken (byggt 2026-11-24)

- **Placering för alla.** I slutet står alla på plats 1 till N: först de
  som kom ut, efter skatt (tomhänt sist av dem), sedan de som blev kvar
  och sist de döda. Lika placering delas. Placeringen står framför namnet
  i resultaten på krönikans skärm och i resultatrutan, och var och en får
  raden "Du kom på plats 2 av 5" (`plats_du`; `Place` i `results` i
  `game.go`). Förut stod de döda före de som blev kvar.
- **Vinken blir skarpare.** När man möter någon säger vinken också om hen
  bär mer eller mindre än du (`pack_mer`, `pack_mindre`, `pack_lika`;
  `packHint` i `pvp.go`). Har ingen av er något säger den inget mer. Hur
  mycket syns fortfarande inte.

## Steg 3: ställningen för kvällen (inte byggt)

Varje Nytt torn ger poäng efter placeringen, till exempel 3, 2 och 1, och
ställningen följer med mellan tornen tills spelledaren nollställer den på
värdsidan. Det rör värdsidan och vad tornet minns mellan vandringarna;
fråga innan det byggs.
