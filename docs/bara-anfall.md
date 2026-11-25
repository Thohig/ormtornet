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

En lång strid blir lättare åt båda hållen (provspelet 2026-11-25,
tredje): efter `trott_efter` rundor i rad mot samma varelse (4) träffar
båda `trott_per_runda` procentenheter lättare (10) för varje runda till,
och den som strider får en rad om det. Med AI-spelare på Normal gick den
längsta striden mot gestalten i fjällrustning från 24 till 12 rundor, och
andelen strider på 8 rundor eller mer från 18 till 12 %, medan andelen
som dog var densamma (13 % på Normal, 23 % på Svår).

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
2. **Föremålen** (byggt 2026-11-24): föremål kan ha `kraft` och `forsvar`
   i `data/foremal.json`. Rostigt svärd (Kraft +1), Fjällbrynja (Försvar
   +2, men slaget mot ryggen träffar alltid när man flyr), skölden som
   Försvar +1 i tre träffar, och det förbannade Svärdet med rött fäste
   (Kraft +2, men ett liv varje runda i strid då man inte anfaller; det
   heter inte "som törstar", eftersom ett förbannat föremål ska se
   värdefullt ut). Kraft och Försvar står bland tecknen (⚔ och ◆, gröna
   när ett föremål höjer dem) och på klasskorten.
3. **Balansen** (byggt 2026-11-24): simuleringen igen och nya tal för
   klasserna och två varelser (se nedan), och panelen Nu som säger i ord
   hur striden ser ut: "Du träffar lätt · den träffar ibland"
   (`fightOutlook` i `monster.go`). Lätt är minst 65 %, ibland minst 35 %,
   annars sällan.

## Balansen

3000 vandringar med sex AI-spelare (`-simulera 3000 -ai 6`), frö 1 och
5001, på Normal och Svår, som i [balans-2026-11-22.md](balans-2026-11-22.md).
Klasserna delas ut i tur och ordning, så att varje klass spelas lika ofta.

Med de första talen blev striden hårdare än förut: den grönaktiga
skepnaden (Kraft 3) stod för 30 % av dödsfallen, och gestalten i
fjällrustning (Försvar 3) gick nästan inte att träffa. Legoknekten och
Kartritaren vann för ofta, Lyktbäraren, Ormtjusaren och Tjuven för sällan.

| Ändring | Före | Efter |
| --- | --- | --- |
| den grönaktiga skepnaden | Kraft 3 | Kraft 2 |
| gestalten i fjällrustning | Försvar 3 | Försvar 2 |
| Kartritaren | 1 / 1 | 0 / 1 |
| Tjuven | 1 / 2 | 2 / 2 |
| Ormtjusaren | 1 / 1 | 2 / 2 |
| Lyktbäraren | 1 / 1 | 2 / 2 |

Klasserna, Kraft / Försvar och andel vinster efter justeringen (två frön):

| Klass | Kraft / Försvar | Normal | Svår |
| --- | --- | --- | --- |
| Kartritaren | 0 / 1 | 19–20 % | 20–21 % |
| Legoknekten | 3 / 2 | 18–20 % | 18–20 % |
| Lyktbäraren | 2 / 2 | 18–19 % | 16–19 % |
| Örtkvinnan | 0 / 2 | 16–17 % | 18–21 % |
| Ormtjusaren | 2 / 2 | 17 % | 14–15 % |
| Pilgrimen | 0 / 3 | 15–16 % | 16–18 % |
| Tjuven | 2 / 2 | 14–15 % | 14–15 % |
| Spårfinnaren | 2 / 1 | 14 % | 15 % |

Före bytet av strid låg klasserna mellan 14 och 22 % på Normal. Nu dör 13 %
på Normal (12 % förut) och 23 % på Svår (24 % förut). Guldet per
vandrare är lägre än förut, på Normal 36–47 mot 39–52: AI-spelarna valde
förut det drag som slog varelsens vanligaste val, och nu anfaller de
alltid. Den grönaktiga skepnaden står för 20 % av dödsfallen.

Provat och inte valt: Legoknekten 2 / 2 (16–18 %, då slår hen inte
hårdast), Lyktbäraren 1 / 2 (13–17 %) och 2 / 1 (12–16 %).

Talen gäller AI-spelare som alltid anfaller. Människor kan söka eller
hålla garden när varelsen är lugn, och kan spela annorlunda.
