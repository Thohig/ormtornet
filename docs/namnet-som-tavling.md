# Ormens namn som en tävling

Status: förslaget skrevs 2026-11-23 och är byggt i fyra steg samma dag.
Avsnittet Stegen säger hur, och Att prova vad som behöver provas i spel.

## Varför

Pusslet med ormens namn byggde på samarbete: tornet delade ut namnets
delar så att gruppen tillsammans hade så många olika som möjligt, och
tornets nåd belönade gruppen när alla var invigda. Men samarbetet var
svagt. Alla fick hela namnet på nivå 4 utan att göra något, och delarna
kunde sägas högt vid bordet utan risk.

Spelet är en tävlan mellan vandrarna. Pusslet ska därför vara svårare,
men var och en ska kunna lösa det själv, och ingen ska tjäna på att
berätta vad hen vet. Det svåra ska vara namnet, inte vad man ska säga om
det, så att den som har spelat förut inte får något stort försprång.
Sedan 2026-11-25 räcker namnet ensamt i viskningen, och Vandrarens bok
säger att namnet finns och var delarna hittas, men inte att det ska
viskas (docs/avslojandet.md, docs/vandrarens-bok.md).

## Stegen

1. **Avslöjandet går inte att stänga av.** Byggt 2026-11-23.
2. **Samarbetet tas bort.** Delarna delas inte ut, tornets nåd ersätts av
   den förste invigde, och ledtrådarna är desamma för alla. Byggt
   2026-11-23: delarna delades först ut till var och en för sig, lottat,
   tills steg 3 lade dem i tornet. Den förste invigde får en fjärdedel mer
   guld till när hen kommer ut (`invigd.forst_guld`).
3. **Delarna finns i tornet.** En del är ristad i ett rum, en finns hos en
   varelse och en hörs i kammaren. Byggt 2026-11-23 (`namndelar.go`).
   Ristningen bestäms när tornet bygger ett rum som passar, så att den
   alltid går att nå; finns den inte på nivå 3 läggs den i ett känt rum.
   Riktningen till den ges från nivå 3, och till Spårfinnaren från nivå 2.
   Varelsens del kommer när man träffar den, inte först när den faller,
   så att den som slåss en runda får den.
4. **Ordningen och straffet.** Delarna kommer utan plats, men med en
   ledtråd om var de står. Fel ordning förbrukar inte försöket. En fel
   viskning gör ormen oroligare. Byggt 2026-11-23, med två ändringar mot
   förslaget:
   - Ledtråden om platsen är var på en orm delen sitter: vid huvudet, i
     buken eller i svansen (`platser`). Huvudet först och svansen sist
     får man lista ut själv. Tecknen på spelskärmen står i den ordning
     delarna hittades, så att de inte avslöjar ordningen.
   - Fel ordning **förbrukar** försöket. Annars kunde man pröva alla sex
     ordningarna samma runda utan att förlora något. Tornet säger ändå att
     ordningen är fel, och ormen blir inte oroligare av det.

## Att prova

- Hittar vandrarna ristningen i tid? Chansen per rum som passar är
  `ristning_chans` (0,3); riktningen ges från nivå 3.
- Får de som undviker strid varelsens del? Den kommer vid första träffen.
- Är huvudet, buken och svansen för lätt eller för svårt att förstå?
- Blir det för dyrt att gissa fel nu när ormen blir oroligare?

Att skrapa bort en ristning är inte med.
