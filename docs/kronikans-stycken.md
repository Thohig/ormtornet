# Rummet i krönikan som ett stycke (förslag)

Förslag 2026-11-22. Byggt samma dag, med svaren i
[samlat-forslag.md](samlat-forslag.md).

## Vad som önskas

Rumsblocket ska inte vara en egen ruta. Det ska gå upp i loggen, krönikan
under vandringen, som ett stycke: en beskrivning och vad som syns och hörs.
Resten får spelarna lista ut själva.

## Läget i dag

En runda där vandraren söker, utlöser det roterande rummet och hamnar hos
Fjällvakten ser ut så här i loggen:

```
  Du letar för noga och utlöser något.
  Det roterande rummet!
  Rummet börjar vrida sig med ett långsamt, malande ljud. …
  Fjällvakten!
  Fjällvakten står orörlig mitt i rummet. Dess rustning ser vuxen ut …
  Stridsljud hörs i tornets nordvästra del.
┌ Brunnsrummet ─────────────────────────────┐
│ En brunn utan synlig botten mitt i golvet. │
│ Fjällvakten är här.                        │
└────────────────────────────────────────────┘
```

Det mesta sägs två gånger: varelsen med ett utrop, en beskrivning och
"Fjällvakten är här". Rumsblocket säger rakt ut vad som finns ("Det
roterande rummet har redan slagit till här", "Rummet är tomt"), och ljuden
står som lösa rader utanför.

## Förslag

Varje runda blir ett stycke i krönikan, utan ruta, i tre delar som flyter
ihop:

1. **Det som hände dig**, som i dag: draget, fällan, striden.
2. **Rummet du står i**: namnet, beskrivningen och det du ser, i ord som
   beskriver och inte förklarar.
3. **Det du hör**: ljuden från grannrummen och längre bort, sist i stycket.

```
  Runda 6
  Du letar för noga och utlöser något. Rummet börjar vrida sig med ett
  långsamt, malande ljud, och du knuffas ut genom en öppning.

  BRUNNSRUMMET. En brunn utan synlig botten mitt i golvet. Långt därnere
  droppar vatten. Fjällvakten står orörlig mitt i rummet; dess rustning
  ser vuxen ut snarare än smidd.
  Från nordväst hörs stridsljud.
```

### Det du ser, i stället för det som förklaras

| I dag | Blir |
| --- | --- |
| Fjällvakten är här. | Varelsens egen text, som redan finns (`dyker_upp` i `monster.json`). Utropet "Fjällvakten!" försvinner. |
| Sigge står här med dig. | Sigge står i dörröppningen och ser på dig. |
| Här går en port ut ur tornet. | Kvällsljuset faller in genom porten. |
| Porten här har slutit sig. | Där porten var finns bara mur, varm och fuktig. |
| Här ligger ormens skatt. | Guld ligger i en hög mitt i rummet, och något andas under den. |
| Det roterande rummet har redan slagit till här. | Ett spår som varje fälla lämnar, till exempel "Golvet är repat i långa cirklar." Fällans namn sägs inte. |
| Någon har lämnat rep med kroken kvar här. | I ett hörn ligger ett rep med en krok. |
| Rummet är tomt. | Inget. Står inget i rummet sägs inget. |

Spåren efter fällorna blir nya texter i `data/fallor.json`, ett eller två
för varje fälla (`spar_efter`). Övriga texter står i `data/texter.json`.
Allt kommer på båda språken.

### Står du kvar

Beskrivningen kommer inte igen. Stycket säger bara det som har ändrats i
det du ser ("Fjällvakten har gått."), och det du hör. Har inget ändrats
och inget hörs blir det inget stycke om rummet.

### Ljuden

Ljuden från grannrummen och de avlägsna ljuden (`sounds.go`, `hear`)
samlas under rundan och står sist i stycket i stället för som egna rader.
Riktningen finns kvar ("Från nordväst hörs stridsljud.").

### Var det syns

- I loggen under vandringen, överst i rundan som rumsblocket i dag, men
  som löpande text: rummets namn i versaler och sedan styckena.
- I Din vandring och slutets krönika står samma stycke under rundan.
- "mer ↓" på telefonen finns kvar för långa beskrivningar.

## Steg

1. **Det du ser.** De nya texterna i stället för `har_*`, spåren efter
   fällorna och varelsernas egna texter. Utropet med varelsens namn och
   "Rummet är tomt" försvinner.
2. **Det du hör.** Ljuden samlas under rundan och skickas med rummet.
3. **Stycket.** Spelarsidan ritar rummet och ljuden som löpande text i
   loggen, utan ruta, och Din vandring och krönikan får samma stycke.

## Frågor

- Ska varelsens namn sägas ("Fjällvakten står orörlig …"), eller bara
  beskrivas ("En gestalt i rustning står orörlig …") så att spelarna får
  lära sig vilka som är vilka?
- Ska fällans namn bort helt när den har slagit till, eller sägas en gång
  när den slår till, som i dag?
- Ska kortet efter rundan finnas kvar på kartan? Det säger mycket av samma
  sak, som "Du gick österut till Vaktrummet" och solen.
- Ska ljuden säga riktningen ("från nordväst") eller vara vagare ("någon
  annanstans i tornet")?
