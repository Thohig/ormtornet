# Kombinationer

Idé 15 i [ideer.md](ideer.md), byggd 2026-11-25: två föremål som
tillsammans blir något nytt.

## Hur det fungerar

- Bär man två föremål som hör ihop fogas de ihop av sig själva, med en rad
  i loggen och ett spelkort. Det nya föremålet tar den förstas plats, och
  den andras plats i väskan blir ledig. Det kostar inget drag.
- Det går även när väskan är full: ett föremål som passar ihop med något
  man bär får plats ändå (`give` och `combine` i `items.go`).
- Beskrivningarna av myntet, veken, flöjten, amuletten och kronan antyder
  vad de passar ihop med, utan att säga det. Man upptäcker kombinationen
  genom att bära båda.
- Delar som verkar hela vandringen (`behaller`) verkar vidare i det nya
  föremålet, med sina förbannelser. En förbannelse som redan har visat sig
  visas inte igen.
- Ett hopfogat föremål med `kvar` blir det föremålet när användningarna är
  slut, i stället för att försvinna.
- Rundans kort säger "Du fogade ihop …".

## Kombinationerna

| Delarna | Blir | Vad det gör |
|---|---|---|
| Myntet med hål + repet | Pendeln | visar som ett tecken åt vilket håll närmaste guld som ligger på ett golv finns (den dödes guld, den gamla kammarens guld); värd 25, och myntets förbannelse verkar vidare |
| Spegelskärvan + extra veken | Ljusspegeln | lyktan slocknar inte i mörkret eller i draget så länge man bär den, och i ett nytt rum visar den två rum bakom dörrarna, tre gånger |
| Fjällamuletten + ormens skinn | Fjällmanteln | varelsernas lynne halveras; värd 30, och skinnets förbannelse verkar vidare |
| Benluren + den sövande flöjten | Ormens pipa | hör som benluren och drar om ett vakna-kort två gånger; sedan blir den benluren igen |
| Den svarta pärlan + kronan av rost | Kronan med pärlan | värd 100 i stället för 40 + 30, med båda förbannelserna |

Repet och kastkniven skulle bli Snaran, en fälla man lägger själv (idé 5).
Den byggdes inte.

## I data/

De nya föremålen står sist i `lista` i `data/foremal.json`, med vikt 0 så
att de bara går att få genom en kombination:

- `av`: de två föremålen som fogas ihop,
- `behaller`: de av delarna som verkar vidare (bara sådana som verkar hela
  vandringen),
- `kvar`: föremålet det blir när användningarna är slut,
- `fogas`: raden i loggen när de fogas ihop.

Vad de nya föremålen gör styrs av id:t, som för de andra: `pendulum`,
`mirror` och `keepsLight` i `items.go`, `lynne` i `monster.go` och
`flute`. Kontrollen av datafilerna (`validateCombos` i `validate.go`)
säger till om en del inte finns, om två kombinationer har samma par eller
om en del som verkar några gånger står i `behaller`. Bilderna ritas i
`web/kort.js`. På kartan har de hopfogade föremålen ingen egen figur.

## Simuleringen

2000 vandringar med 5 AI-spelare på Normal: Ljusspegeln fogades ihop 13
gånger, Pendeln 6, Ormens pipa 5, Fjällmanteln 1 och Kronan med pärlan
aldrig, alltså ungefär en gång på åttio vandringar. AI-spelarna söker
sällan, och de flesta delar är ovanliga (vikt 1–3, skinnet bara genom en
händelse).

Efter provspelet 2026-11-25 väger den andra delen sex gånger mer för den
som söker och redan bär den ena (`del_vikt` i `foremal.json`,
`randomItemFor` i `items.go`). Det gav 60 kombinationer på 2000
vandringar, ungefär en på 33: Ljusspegeln 27, Pendeln 13, Ormens pipa 12,
Kronan med pärlan 7 och Fjällmanteln 1. Delarna är inte vanligare för
den som inte bär någon.

## Att prova

- Märker man att föremålen fogas ihop, och letar någon efter den andra
  delen när beskrivningen antyder den?
- Händer det tillräckligt ofta med människor, som söker mer? Om inte kan
  delarnas vikter höjas.
- Är Fjällmanteln för stark, och används Pendeln för att hitta den dödes
  guld?
