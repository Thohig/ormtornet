# Föremål med två sidor

Idé 14 i [ideer.md](ideer.md), byggd 2026-11-25: föremål som är bra i
början och sämre mot slutet, eller som blir starkare ju längre man bär dem.

## Hur det fungerar

- Föremålen tas upp av sig själva och går inte att släppa. Ett föremål
  med två sidor blir därför inget val när man hittar det; valet är när man
  ska gå ut, innan det vänder.
- Ett föremål **vänder** när avslöjandet når dess nivå (`vander_niva`),
  alltså när solen sjunker, eller när det har burits ett visst antal
  rundor (`vander_rundor`). Rundorna räknas i slutet av varje runda man är
  i tornet och inte i källarvalven (`ageItems` i `items.go`).
- När det vänder kommer texten under `vander` och ett spelkort, som när en
  förbannelse visar sig (`turnItem`).
- Beskrivningen antyder vad som väntar, utan att säga det rakt ut.
- De är inte förbannade, så Pilgrimen tar dem. De fogas inte ihop med
  något.

## Föremålen

| Föremål | Den goda sidan | Den andra sidan |
|---|---|---|
| Ormens öga | ett tecken (◉) visar åt vilket håll kammaren ligger | från nivå 3 lägger varje näve ett vakna-kort i leken, som den svarta pärlan; stenen pekar fortfarande |
| Ägget | värt 50 guld när du kommer ut | efter 12 rundor i väskan kläcks det: ett yngel kryper fram i ditt rum, och ägget är borta |
| Kniven som minns | Kraft +1 för varje varelse eller vandrare du dödar medan du bär den, högst +3 | varelsernas lynne +10 för varje dödande; kortet vänds vid det första |
| Glödmasken | mörkret biter inte, och du ser varelser i grannrummen, som Lyktbäraren | efter 10 rundor lyser den genom väskan, och de andra ser dig från grannrummet |

Alla har vikt 1 i `data/foremal.json`.

## I koden

- `ageItems`, `turnItem`, `eye`, `knifeRemembers`, `knifeScent` och
  `glowing` i `items.go`.
- Ögat i näven: `take` i `serpent.go`.
- Kniven: `defeat` i `monster.go`, och träffen mellan vandrare i `pvp.go`.
- Glödmasken: `keepsLight` i `items.go`, `seesMonsterAt` i `monster.go`
  och `nearLantern` i `classes.go`.
- Ögats tecken: `oga` i tillståndet, och i `web/app.js`.
- Bilderna ritas i `web/kort.js`, och testerna står i `tva_sidor_test.go`.

## Att prova

- Märker man när ett föremål vänder, och går man ut tidigare för att
  hinna före det?
- Är ägget värt att bära, eller blir ynglet för hårt?
- Används ögat för att hitta kammaren, nu när andetagen också gör det?
- AI-spelarna känner inte till föremålen och tar dem som andra föremål.
