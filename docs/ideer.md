# Idéer för framtiden

Idéer som spånades 2026-11-25 för hur spelet kan byggas ut. Inget här är
bestämt; ta fram ett förslag och fråga innan något byggs. Alla följer den
individuella tävlingen (docs/individuell-tavling.md), och ingen säger
tornets hemlighet rakt ut. Idé 17, reprisen, är byggd (docs/repris.md),
idé 15, kombinationerna (docs/kombinationer.md), idé 14, föremålen
med två sidor (docs/tva-sidor.md), och idé 13, de nya klasserna
(docs/nya-klasser.md).

## Tornet

1. **Våningar.** En trappa leder ner till en mörkare våning med mer guld
   och farligare varelser. Det ger ett val mellan att gå djupare och att
   hinna ut.
2. **Rum som minns.** Ett rum där någon har slagits eller dött ändrar sin
   beskrivning: blod på golvet, en trasig sköld. Man kan läsa vad som
   hänt utan att någon berättar.
3. **Tornets humör.** Ett frö med ett tema, som "det drunknade tornet"
   med vatten och rost eller "det brända" med aska och värme, ger andra
   rumsbeskrivningar, fällor och bilder. Det mesta skulle ligga i
   `data/`.
4. **Låsta dörrar** som kräver en nyckel någon annan kanske har hittat,
   och som man kan sälja eller ta.

## Vandrarna mot varandra

5. **Fällor man lägger själv.** Man gillrar en fälla i ett rum, och nästa
   som går in, vem det än är, går i den. Spåret säger kanske att någon
   har varit där.
6. **Falska spår.** Ett föremål som lämnar spår åt fel håll eller låter
   som någon annan.
7. **Ett pris på huvudet.** Den som har burit mest längst blir efterlyst,
   och den som tar guld från hen får extra.
8. **Byteshandel.** När två möts kan de byta föremål eller guld. Det är
   ett utbyte, inte samarbete, och kan vara en lögn.
9. **Lögnen.** En vandrare kan skicka ett falskt varsel till den hen
   möter.

## Ormen och hemligheten

10. **Fler ormar.** Ormen är olika varje vandring, med olika lynnen och
    regler för näven: en som vaknar av ljud, en som vaknar av ljus, en som
    räknar guldet.
11. **Ormens gåta.** Den som är invigd får en fråga av ormen. Ett rätt
    svar ger en säker näve, och ett fel svar väcker ormen.
12. **Spår av tidigare vandringar.** Breven skrivs av förra kvällens
    vandrare, med deras namn och vad de gjorde, så att tornet minns
    mellan vandringarna.

## Klasser och föremål

13. ~~**Nya klasser.**~~ Byggda 2026-11-25 (docs/nya-klasser.md).
14. ~~**Föremål med två sidor.**~~ Byggda 2026-11-25 (docs/tva-sidor.md).
15. ~~**Kombinationer.**~~ Byggda 2026-11-25 (docs/kombinationer.md).

## Kvällen och slutet

16. **En kampanj.** Tre vandringar i rad i samma torn, där kartan, den
    döde och guldet som blev kvar följer med. Ställningen för kvällen
    valdes bort förut (docs/individuell-tavling.md, steg 3), så det här
    vore en annan form av det.
17. ~~**Repris.**~~ Byggd 2026-11-25 (docs/repris.md).
18. **Åskådare.** Den som har dött ser hela tornet och kan, en gång,
    viska en ledtråd till en levande.

## Spelledaren

19. **Spelledarens kort.** Spelledaren får tre händelser att spela ut när
    hen vill, utan att se kartan.
20. **Egna torn.** En redigerare på värdsidan för rum, texter och
    händelser, som sparas som en ny `data/`.

Förslaget var att börja med fällor man lägger själv (5), fler ormar (10)
och reprisen (17), som ger mest nytt spel för minst ny kod.
