# Bara anfall

Striden har ett enda stridsval, anfall. Parera och undanmanöver är borta.
Står en varelse i rummet kan man ändå söka, dricka, vänta eller gå, men då
riskerar man att varelsen anfaller. Vem som träffar avgörs av värden.

## Valen i strid

| Val | Vad som händer |
|---|---|
| Anfall (A) | du och varelsen slår båda; var och en träffar med sin chans |
| Vänta (mellanslag) | du håller garden: varelsen kanske anfaller, men träffar sämre |
| Sök, Drick | varelsen kanske anfaller, och träffar då lättare |
| Gå | du lämnar rummet; varelsen kanske slår mot din rygg först |

Den som inte väljer något står still utan att hålla garden.

## Värdena

- **Kraft**: hur lätt man träffar.
- **Försvar**: hur svårt man är att träffa.
- **Lynne** (bara varelserna): chansen i procent att varelsen anfaller en
  runda då vandraren inte anfaller.

Träffchansen är `bas + per_steg × (Kraft − motståndarens Försvar)`, inom
`min` och `max` (`traff` i `data/monster.json`). Den som håller garden är
svårare att träffa (`garden`), den som söker eller dricker lättare
(`oberedd`). Spelarna ser aldrig talen, bara det som händer.

Invigda och den som bär Myntet med hål gör varelserna mer lynniga
(`extra_lynne`).

## Mellan vandrare

Den som anfaller slår med sin chans; anfaller båda slår båda. Den som väntar
håller garden, den som söker eller dricker är lättare att träffa. Bytet,
portrummens fred och att tornet väljer motståndaren är som förut.

## Klasserna

Kraft och Försvar står i `data/klasser.json`. Legoknektens "ingen skada när
båda anfaller" är borta; hen har i stället högst Kraft. Kartritarens
svaghet är nu att den som slår mot ryggen när hen flyr alltid träffar.

## Stegen

1. **Värden och anfall** (byggt 2026-11-24): värdena i datan, anfall som
   enda stridsval, garden när man väntar, strid mellan vandrare, AI:n,
   spelarsidan och texterna på båda språken.
2. **Föremålen** (inte byggt): Rostigt svärd, Fjällbrynja, skölden som
   Försvar, Svärdet som törstar; värdena bland tecknen och på klasskorten.
3. **Balansen** (inte byggt): simuleringen av strategierna igen, och
   panelen Nu som säger hur striden ser ut.
