# Föremål och hemligheten

Status: byggt 2026-11-18, i den här ordningen. README beskriver hur allt
fungerar under "Föremål" och "Att viska hemligheten" (spoilervarning).

| Steg | Vad | Commit |
|---|---|---|
| Föremål 1 | hitta och bära, sex föremål som verkar av sig själva | dc29fdf |
| Hemligheten 1 | viska, bli invigd, dunket och porten som väntar | 3dbd1a6 |
| Hemligheten 2 | den invigdes fördelar och svaghet | 9b2067b |
| Föremål 2 | föremål man använder själv, AI-spelarnas regler | f405825 |
| Föremål 3 | förbannade föremål, brev, döda tappar, klasserna | e0a53b2 |
| Hemligheten 3 | spelledarens ingrepp och epilogen | a28fcd1 |

## Föremål

- Man bär högst tre (`max_bara` i `foremal.json`). Ett föremål som inte får
  plats blir kvar i rummet, och nästa som kommer in tar det.
- Man hittar föremål när man söker (`foremal` under `sok` → `vikter` i
  `fallor.json`), och den som dör tappar sina.
- Föremålen följer inte med till nästa torn.
- Alla föremål verkar av sig själva sedan 2026-11-22; inget används med
  1–3 längre (se [passiva-foremal.md](passiva-foremal.md)).
- Förbannade föremål visar sin baksida första gången den verkar.
- "Kartan som ljuger" i förslaget blev Kronan av rost, eftersom falska rum
  inte går att göra på ett enkelt sätt.
- AI-spelarna söker sällan och hittar i snitt 0,13 föremål per vandring.

## Hemligheten

- Viskningen jämförs med svaren under `svar` i `avslojande.json`, som hela
  ord. "Ormen sover i tornet" (det alla vet) räknas inte.
- En viskning per nivå, från `viska_fran_niva` (2).
- Tornet bekräftar aldrig hemligheten: den invigde får bara veta att tornet
  har hört. Krönikan och epilogen nämner den invigde i slutet.
- Fördelarna och svagheten står under `invigd` i `avslojande.json`.
- Porten som väntar: när solen går ner är de som inte är invigda förlorade
  direkt, och de invigda får en runda till.
- Med avslöjandet avstängt går det inte att viska, och ingen blir invigd.
- AI-spelarna viskar aldrig.

Hemligheten har sedan byggts ut med ormens namn, byte vid dödande, tornets
nåd, varslen och tornet som förändras; se [avslojandet.md](avslojandet.md).

## Inte byggt

- Att släppa ett föremål med flit, eller byta föremål med en annan vandrare.
- Föremål i spelledarens rum-ingrepp (lägga ett föremål i ett rum).
