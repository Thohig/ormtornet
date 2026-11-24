# Passiva föremål (förslag)

Förslag 2026-11-22. Allt är byggt, som det står här: kritan blev passiv,
spegelskärvan visar ett osett grannrum och timglaset verkar vid tre
handsbredder kvar. Kritan togs senare bort och ersattes av benluren (se
[andringar.md](andringar.md)). Sedan 2026-11-24 gäller timglaset bara
den som bär det: porten väntar en runda till på hen när solen går ner
([individuell-tavling.md](individuell-tavling.md)).

## Vad som önskas

Alla föremål verkar av sig själva och behöver aldrig användas. Då kan
tangenterna 1–3 (och tryck på föremålen i väskan) för att använda dem tas
bort.

## Läget i dag

Sju föremål verkar redan av sig själva: repet, fjällamuletten, extra veke,
den bucklade skölden, väskan med dubbel botten, kompassen och den sövande
flöjten. De förbannade föremålen och breven är också passiva.

Fem föremål används med 1–3, och fyra av dem siktar med pilarna
(`aktiv` och `riktning` i `data/foremal.json`):

| Föremål | Gör i dag |
| --- | --- |
| Rökpulver | Fly ur en strid utan att varelsen hinner slå. Välj håll. |
| Spegelskärva | Titta in i rummet bakom en dörr. Tre gånger. |
| Kastkniv | Kasta genom en dörr mot en varelse i nästa rum. |
| Krita | Rita en pil på golvet. Fem gånger. |
| Timglas av ben | Solen står still en runda för alla. |

## Förslag: samma föremål, men de verkar av sig själva

| Föremål | Blir |
| --- | --- |
| Rökpulver | När du lämnar ett rum där en varelse slåss med dig, tänds pulvret och varelsen hinner inte slå. En gång. |
| Spegelskärva | När du kommer in i ett nytt rum visar skärvan rummet bakom en av dörrarna, ett som du inte har sett. Tre gånger. |
| Kastkniv | När en varelse dyker upp i ditt rum kastar du kniven först: den tar en skada innan striden börjar. En gång. |
| Krita | När du lämnar ett rum ritar du en pil åt det håll du går, så att du och andra hittar tillbaka. Fem rum. |
| Timglas av ben | När solen har tre handsbredder kvar vänds glaset av sig självt och solen står still en runda för alla. En gång. |

Beskrivningarna i väskan säger när föremålet verkar, och loggen säger när
det har hänt, som för de andra föremålen i dag.

## Det som försvinner

- Tangenterna 1–3, att sikta med pilarna efter 1–3 och raden om det i
  panelen Nu (`aimSlot` och `nu_sikta` i `web/app.js`).
- Att trycka på ett föremål i väskan för att använda det. Ett tryck visar
  fortfarande vad föremålet gör.
- `aktiv` och `riktning` i `data/foremal.json`, raden om 1–3 i styrningen
  bakom ? (`hjalp_foremal`) och texterna för att sikta.
- Draget "använd föremål" i servern (`items.go`). Datorspelarna
  (`ai.go`) och robotarna (`verktyg/provspelare.mjs`) slutar använda
  föremål.

Väskan står kvar under tecknen, med en plats för varje föremål man kan
bära.

## Steg

1. **Föremålen verkar av sig själva.** Rökpulvret, spegelskärvan,
   kastkniven, kritan och timglaset får sina nya regler i `items.go`, med
   test. Nya beskrivningar och texter kommer på båda språken.
2. **Bort med 1–3.** Spelarsidan, styrningen, datorspelarna, robotarna och
   dokumentationen.

## Frågor

- Låter de fem nya reglerna rätt, eller ska något föremål tas bort i
  stället? Kritan är svårast att göra passiv; den kan också tas bort helt.
- Spegelskärvan liknar Kartritaren, som ser in i grannrummen. Gör det
  något, eller ska skärvan göra något annat?
- Timglaset: ska det verka vid tre handsbredder kvar, eller när du själv
  är längst bort från porten?
