# Tre nya klasser

Idé 13 i [ideer.md](ideer.md), byggd 2026-11-25: Gravplundraren,
Lyssnaren och Ormens präst. Klasserna är elva, och ingen klass finns två
gånger i samma torn.

## Klasserna

| Klass | Kraft / Försvar | Styrka | Svaghet |
|---|---|---|---|
| Gravplundraren ⚰ | 2 / 2 | hör var någon dör inom sju rum (`hor_dod_rum`), hur långt bort och åt vilket håll ("sex rum bort", eller "i rummet bredvid" med `klass_gravplundrare_dod_bredvid`, sedan 2026-11-26); får dubbelt så mycket av den dödes guld, och det extra kommer ur tornet | varelserna anfaller oftare, lynne +10 |
| Lyssnaren ◎ | 2 / 2 | hör två rum bort och om något andas bredvid, som med benluren, och tre rum med benluren; ser sömnmätaren var hen än är | mörkret varar dubbelt så länge |
| Ormens präst ☥ | 1 / 3 | det första vakna-kortet hen drar under vandringen blir ett sovkort; ormen skadar inte hen när den vaknar; varje näve ger hälften mer guld | tar bara en näve per besök i kammaren |

Idén sade att Gravplundraren tar guld från de döda utan straff. Det finns
inget straff för det, så hen får i stället mer.

## I data/ och koden

- `data/klasser.json`, med `styrka` för klasskortet och `du`, samma
  styrka i andra person för raden i loggen när vandringen börjar ("Du är
  gravplundraren. Du hör var någon dör …", sedan 2026-11-26), och
  `varden`:
  - Gravplundraren: `dod_guld` och `extra_lynne`,
  - Lyssnaren: `morker_ganger`,
  - Prästen: `ormskada`, `en_nave` och `kammar_guld`.
- Gravplundraren: `announceDeath` och `pickUpGold` i `guld.go`, och
  `lynne` i `monster.go`.
- Lyssnaren: `hearReach`, `hearsFrom` och `hornListen` i `benlur.go`,
  sömnmätaren i `stateFor` i `game.go`, och `darkRounds` i `classes.go`.
- Prästen: `take` i `serpent.go`, och `tookVisit`, som nollställs i
  `playTurn` när hen inte står i kammaren. AI-prästen går ut genom en
  dörr när hen har tagit sin näve (`anyDoor` i `classes.go`).
- AI-spelarna: Gravplundraren hos den djärva, Lyssnaren hos den
  försiktiga och Prästen hos plundraren (`data/ai.json`).
- Klasskorten ritas i `web/kort.js` och figurerna i
  `klient/src/modeller.ts`, byggd till `web/iso/iso.js`.
- Klassvalet visar "Tryck 1–9" när korten är fler än nio.
- Testerna står i `nya_klasser_test.go`.

## Simuleringen

3000 vandringar med 6 AI-spelare, klasserna rättvist fördelade. De nya
klasserna vann först 11–13 %, mot 13–24 % för de andra, och fick då mer
Kraft och Försvar, Gravplundraren mindre lynne och Prästen mer guld per
näve. Efter det:

| Klass | Vann på Normal | Vann på Svår |
|---|---|---|
| Gravplundraren | 15 % | 14 % |
| Lyssnaren | 14 % | 16 % |
| Ormens präst | 22 % | 21 % |
| de andra åtta | 13–21 % | 15–20 % |

## Att prova

- Går Gravplundraren mot skriken, och blir hen en gam som väntar på att
  andra ska dö?
- Har Lyssnaren nytta av att höra längre, eller blir det bara fler rader?
  Hörs andetagen för mycket när sömnmätaren syns överallt?
- Är en näve per besök för hårt för Prästen, eller för lätt med hälften
  mer guld?
- Går elva klasskort att läsa på 30 sekunder, särskilt på telefonen?
- Figurerna på kartan är bara byggda, inte sedda i en vandring.
