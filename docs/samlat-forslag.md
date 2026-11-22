# Samlat förslag: stycken i krönikan och färre kommandon

Förslag 2026-11-22. Byggt samma dag, steg för steg. Tornet väljer
motståndare efter den som anföll, annars den som kom först in i tornet,
eftersom skatten är hemlig. Det här samlar
[kronikans-stycken.md](kronikans-stycken.md) och
[farre-kommandon.md](farre-kommandon.md), med ett svar på varje öppen
fråga och en ordning att bygga i. Detaljerna står i de två dokumenten.

## Målet

Spelet ska säga mindre och låta spelarna lista ut mer. Krönikan berättar
vad du ser och hör, inte vad saker är. Styrningen visar bara det du kan
göra just nu.

## Krönikan

- **Varje runda är ett stycke**, utan ruta: det som hände dig, sedan
  rummet (namnet i versaler, beskrivningen och det du ser) och sist det du
  hör.
- **Det du ser beskrivs, det förklaras inte.** Porten, guldet, vandrarna,
  föremålen och spåren efter fällor får egna meningar. Ett tomt rum säger
  ingenting.
- **Står du kvar** kommer bara det som har ändrats och det du hör.
- **Samma stycke** står i Din vandring och slutets krönika.

Svar på frågorna:

| Fråga | Svar |
| --- | --- |
| Ska varelsens namn sägas? | Ja. Namnen står i striden, på kortet och i krönikan, och det är svårt att prata med de andra utan dem. Utropet "Fjällvakten!" och "Fjällvakten är här" försvinner; varelsens egen text räcker. |
| Ska fällans namn sägas? | En gång, när den slår till, som i dag. Efteråt syns bara spåret i rummet. |
| Ska kortet efter rundan finnas kvar? | Ja, men kortare: bara det som ändrade dig (liv, guld, föremål, ormkort, gift och mörker). Vart du gick och solen står redan i stycket och i översta raden. |
| Ska ljuden säga riktningen? | Ja. Det är det enda sättet att veta var de andra är, och ljuden är redan vaga om vem det är. |

## Styrningen

- **Sök tar också skatt.** S söker i vanliga rum och tar en näve i
  kammaren. T försvinner.
- **Knapparna följer läget.** Utanför strid: pilarna, ·, Sök och Drick
  (och Viska). I strid: Anfall, Parera, Undan och Drick; att gå är att
  fly. I lobbyn och slutet: Redo eller Nytt torn. Styrningen bakom ? visar
  bara det som gäller just nu.
- **Tornet väljer motståndare.** Tab och Byt försvinner.
- **Z för kartan.** + − och M blir en tangent som går runt mellan Nära,
  Mellan och Allt.

Svar på frågorna:

| Fråga | Svar |
| --- | --- |
| Alla fyra delarna? | Ja. De är oberoende av varandra och små var för sig. |
| Ska drycken verka av sig själv? | Nej. Att dricka kostar en runda, och det är ett av få val som spelar roll mitt i en strid. D och Drick finns kvar. |
| Hur väljer tornet motståndare? | Den som anföll dig förra rundan, annars den som bär mest guld. Panelen Nu säger vem. |

Efteråt: 11 tangenter i stället för 16, och högst 9 knappar på
pekskärmen i stället för 15.

## Ordning

Varje steg är en egen commit, och du provar mellan stegen.

1. **Sök tar också skatt** (styrning). Litet och fristående.
2. **Knapparna följer läget** (styrning). Gör pekskärmen lugnare direkt.
3. **Tornet väljer motståndare och Z för kartan** (styrning).
4. **Det du ser** (krönika): de nya meningarna, spåren efter fällorna,
   bort med utropet och "Rummet är tomt".
5. **Det du hör** (krönika): ljuden samlas och står sist i stycket.
6. **Stycket** (krönika): löpande text utan ruta i loggen, i Din vandring
   och i krönikan, och kortet efter rundan blir kortare.

Styrningen kommer först eftersom den är mindre och syns i varje runda.
Krönikan bygger på rumsblocket som finns i dag och ändrar det i tre steg.
