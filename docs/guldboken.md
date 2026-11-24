# Guldboken

På slutskärmen kan man välja en vandrare och läsa hur hen fick allt sitt
guld, som en längre berättelse, med raderna runda för runda bakom en
knapp. Användaren bad om det 2026-11-24, tillsammans med att slutskärmen
skulle bli enklare och att hjärtslagen inte skulle höras där.

## Slutskärmen

Uppifrån och ner:

1. **Placeringen**: en rad för var och en (`renderPeople` i `web/app.js`)
   med platsen, ★ för vinnaren, namnet i sin färg, titeln, (du), ~ för
   AI:n och guldet eller hur det slutade. Under raden står hens största
   stund (`stund` i `vandrarna`). Resultaten och raden om var och en var
   två listor; nu är de en.
2. Ett klick på raden öppnar **guldboken** under den. Ett nytt klick, eller
   ett klick på en annan vandrare, stänger den. Knappen **Runda för runda**
   visar raderna.
3. **Berättelsen om vandringen** (docs/berattelsen.md), utan raden om var
   och en, som nu står i placeringen.
4. Epilogen.
5. **Visa din vandring**: den egna vandringen runda för runda, som förut.
   Hela krönikan står inte längre på skärmen och skickas inte i
   tillståndet; `chronicleLines` finns kvar för testerna.

Tornets dunk (`Sound.level`) är tyst på slutskärmen.

## Hur den byggs

Varje gång en vandrares guld ändras skrivs det ner med `goldNote` i
`guldbok.go`: rundan, slaget, hur mycket (plus eller minus), vad
vandraren bar efteråt och variablerna för texten. Slagen är:

| slag | var | variabler |
| --- | --- | --- |
| `sok` | att söka ett rum (`events.go`) | `rum`, den korta beskrivningen |
| `varelse` | att fälla en varelse (`monster.go`) | `varelse` |
| `kammaren` | en näve ur kammaren (`serpent.go`) | `kort`, ormkortet |
| `ormen` | ormen vaknar och kastar ut (minus) | |
| `slag`, `dodade` | att slå eller döda en vandrare (`pvp.go`) | `namn` |
| `slagen`, `dodad` | att bli slagen eller dödad (minus) | `namn` |
| `dod_guld` | guldet en död bar (`guld.go`) | `namn` |
| `lost`, `lost_orm` | löst guld, från en händelse eller där ormen låg | |
| `tappade` | att dö och lämna guldet (minus) | |
| `foremal` | föremålen blir guld vid porten (`game.go`) | |
| `invigd`, `invigd_forst` | ormens namn ger mer guld vid porten | |

När vandringen är slut skriver `goldBooks` en bok för var och en i
placeringens ordning, och den skickas som `guldbok` i tillståndet:

- **Berättelsen** i tre stycken. Det första säger vem hen var (klassen och
  namnet) och hur det gick (`gb_in_*`). Det andra är vägen: ändringarna i
  tur och ordning, var och en med när den hände räknat från den förra
  (samma `{nar}` som berättelsen om vandringen), och flera av samma slag
  med högst tre rundor emellan som en mening (`gb_flera_*`). Den största
  fångsten pekas ut (`gb_storst`, `gb_storst_av`), och efter en förlust
  står vad som var kvar (`gb_kvar`). Utan guld står `gb_aldrig`. Det
  tredje är slutet: det porten lade till, vem som tog upp guldet efter den
  som dog, att guldet blev kvar i tornet, summan och placeringen för den
  som kom ut med guld.
- **Raderna**: en för varje ändring, `gb_rad_*`, utan beloppet, som
  spelarsidan sätter bredvid med plus eller minus och summan.

Texterna står i `data/texter.json` och `data/en/texter.json` under `gb_`.
Robotarna skriver böckerna i provloggen (`GULDBOK`, `GULDRAD`), och
webbläsarprovet öppnar en bok och rundorna.

## Att prova

- Läses berättelsen bra för den som fick mycket guld från många håll, och
  blir den för lång?
- Säger raderna runda för runda tillräckligt, eller saknas rummet och
  varelsen där?
- Är det tydligt att man kan klicka på en vandrare?
