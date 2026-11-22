# Färre kommandon (förslag)

Förslag 2026-11-22. Byggt samma dag, med svaren i
[samlat-forslag.md](samlat-forslag.md). Tornet väljer motståndare efter
den som anföll, annars den som kom först in i tornet, eftersom skatten
är hemlig.

## Läget i dag

På datorn finns 16 tangenter i styrningen bakom ?, och på en pekskärm 15
knappar längst ner:

| Tangent | Gör | Knapp |
| --- | --- | --- |
| Pilar | gå | ↑ ← → ↓ |
| Mellanslag | vänta | · |
| T | ta skatt (bara i kammaren) | Ta |
| S | sök i rummet (bara i vanliga rum) | Sök |
| D | drick en dryck | Drick |
| A | anfall (i strid) | Anfall |
| P | parera (i strid) | Parera |
| U eller E | undanmanöver (i strid) | Undan |
| Tab | byt mål när flera vandrare står i rummet | Byt |
| V | viska till tornet (med avslöjandet) | Viska |
| + − | zooma kartan | knappar på kartan |
| M | allt du känner till | knapp på kartan |
| Enter | redo | Redo |
| ? | visa och dölj styrningen | ? |

Alla knappar står kvar hela tiden, också de som inte går att använda just
nu: Anfall, Parera och Undan utanför strid, Ta utanför kammaren och Sök i
kammaren.

## Förslag

### 1. Ta och Sök blir en: S

Ta går bara i kammaren och Sök bara i vanliga rum, så de krockar aldrig.
**S** (knappen **Sök**) tar en näve guld i kammaren och söker i andra rum.
T försvinner. Panelen Nu säger som i dag vad S gör där du står ("S tar en
näve ur högen").

### 2. Knapparna följer läget

Bara det som går att göra just nu syns:

- **Utanför strid:** pilarna, ·, **Sök** och **Drick** (och **Viska** när
  du kan viska).
- **I strid:** pilarna (att gå är att fly), **Anfall**, **Parera**,
  **Undan** och **Drick**.
- **I lobbyn och slutet:** **Redo** eller **Nytt torn**.

På en pekskärm blir det högst nio knappar i stället för femton, och de
står på samma ställen varje gång, så att fingret hittar dem. Tangenterna
fungerar som förut; styrningen bakom ? visar bara de som gäller just nu.

### 3. Byt mål försvinner

Tab och **Byt** behövs bara när två eller fler andra vandrare står i
rummet och strid är på. Tornet väljer i stället: den som anföll dig förra
rundan, annars den som bär mest guld. Panelen Nu säger vem du slåss mot,
som i dag.

### 4. Kartan: en tangent

+ − och M blir en tangent, **Z**, som går runt mellan Nära, Mellan och
Allt du känner till. Knapparna på kartan finns kvar för musen och fingret.

### Resultat

| | I dag | Efter |
| --- | --- | --- |
| Tangenter i styrningen | 16 | 11: pilar, mellanslag, S, D, A, P, U, V, Z, Enter, ? |
| Knappar på pekskärmen | 15 | högst 9 åt gången |

## Kanske, men inte i förslaget

- **Drycken verkar av sig själv:** du dricker när ett slag annars skulle
  döda dig. Då försvinner D och Drick, men det tar bort ett val som spelar
  roll: att dricka kostar en runda.
- **Vänta försvinner:** den som inte väljer något väntar redan när tiden
  är slut. Men mellanslaget låter rundan gå fortare när alla är klara.
- **Strid med pilarna:** i strid blir ↑ anfall, ← parera och → undan. Det
  sparar tre tangenter men gör pilarna svårare att förstå.

## Steg

1. **Sök tar också skatt.** Servern tar emot S i kammaren som en näve, och
   T försvinner (texterna på båda språken, styrningen och README).
2. **Knapparna följer läget.** Spelarsidan visar bara de knappar och
   tangenter som gäller just nu, i styrningen och på pekskärmen.
3. **Tornet väljer mål.** Tab och Byt försvinner.
4. **Z för kartan.** + − och M blir en tangent.

## Frågor

- Ska alla fyra delarna byggas, eller några av dem?
- Ska drycken verka av sig själv (se Kanske)?
- Är "den som anföll dig, annars den som bär mest guld" rätt sätt för
  tornet att välja mål?
