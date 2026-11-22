# Värdsidan i två kolumner (förslag)

Förslag 2026-11-22. Byggt samma dag, som det står här. Raden stod redan
fast överst, varningen om tornets flik står kvar över raden, och
förklaringarna till inställningarna står under valen i två kolumner.

## Läget i dag

Värdsidan (`web/spelledare.html`, `web/spelledare.css`) är en kolumn, högst
760 punkter bred, mitt på skärmen. Överst står rubriken och raden, och
sedan avsnitten Bjud in, Vid porten, Inför nästa vandring, Kommer någon
inte in?, Läs om innehållet och Stäng tornet. På en bred skärm blir det
mycket tom yta på sidorna, och spelledaren måste rulla för att se både
vandrarna och inställningarna.

## Förslag

På en bred skärm, från ungefär 1100 punkter, står avsnitten i två
kolumner. Rubriken och raden står över båda.

```
┌ ORMTORNET · VÄRD ─────────────────────────────────────────────┐
│ Lobbyn · 3 människor och 2 AI här                [Starta nu]  │  ← raden
├───────────────────────────────┬───────────────────────────────┤
│ BJUD IN                       │ INFÖR NÄSTA VANDRING          │
│ https://…#torn=…&nyckel=…     │ Svårighet      [Normal ▾]     │
│ [Kopiera inbjudan] [Kopiera…] │ AI-spelare     [2 ▾]          │
│                               │ Språk          [Svenska ▾]    │
│ VID PORTEN                    │ Avslöjandet    [På ▾]         │
│ Väntar på vandrare.           │ Frånvaro       [AI tar över ▾]│
│ ● Dagny Ekfot      redo       │ ▸ Vad inställningarna gör     │
│ ● Sigge            AI         │                               │
│ · ledig plats                 │ ▸ KOMMER NÅGON INTE IN?       │
│ Tornet: 9 × 9 rum, 40 rundor. │ ▸ LÄS OM INNEHÅLLET           │
│ [Starta nu]                   │ STÄNG TORNET                  │
└───────────────────────────────┴───────────────────────────────┘
```

- **Vänster kolumn, det som händer nu:** Bjud in och Vid porten.
- **Höger kolumn, förberedelser och hjälp:** Inför nästa vandring, Kommer
  någon inte in?, Läs om innehållet och Stäng tornet.
- **Ordningen är densamma.** Läser man vänster kolumn och sedan höger är
  avsnitten i samma ordning som i dag, som CLAUDE.md kräver, och en
  skärmläsare läser dem så.
- **Bredden** blir högst ungefär 1240 punkter, med 32 punkter mellan
  kolumnerna. Länken under Bjud in bryts som i dag om den är för lång.
- **Raden står fast överst** när sidan rullar, så att fasen och knappen
  som behövs just nu alltid syns. Det gäller också i en kolumn.
- **På en smal skärm**, som en surfplatta på höjden eller ett litet
  fönster, är sidan en kolumn som i dag.
- **På GitHub Pages** är det likadant. Varningen om att hålla tornets flik
  öppen står över båda kolumnerna, under raden.

## Hur det byggs

Ett steg: två `div` runt avsnitten i `web/spelledare.html` och ett
rutnät i `web/spelledare.css`, med en mediefråga för den smala skärmen.
Inga texter ändras, och servern rörs inte.

## Frågor

- Ska raden stå fast överst när sidan rullar?
- Ska brytpunkten vara ungefär 1100 punkter, eller ska två kolumner komma
  redan på en vanlig bärbar dator med ett mindre fönster?
- Ska Vad inställningarna gör vara utfällt från början när det nu får
  plats i höger kolumn?
