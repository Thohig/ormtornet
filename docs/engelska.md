# Engelska

Plan för att tornet ska kunna spelas på engelska. Spelledaren väljer
språket på värdsidan, och hela tornet följer valet: spelarsidan,
värdsidan och spelets texter.

## Läget

Inget är byggt. Steg 1 är nästa steg.

## Beslut

- **Ett språk för hela tornet.** Inställningen Språk ligger på värdsidan
  under Inför nästa vandring, med valen Svenska och English, och som
  flaggan `-sprak en`. Den sparas som de andra inställningarna.
- **Bytet gäller som Läs om innehållet:** genast i lobbyn, där tornet
  byggs om, och annars från nästa torn. Ingen vandring byter språk mitt i.
  Öppna sidor byter text utan att laddas om.
- **Värdsidan följer språket**, liksom spelarsidan.
- **Terminalen förblir svensk:** utskrifterna från `main.go`,
  `-simulera` och `-kontrollera`.
- **Bara texter översätts.** Engelskan ligger i `data/en/` med samma
  filnamn som i `data/`. Filerna innehåller bara texter, och listor som
  klasser, varelser, fällor och föremål kopplas ihop med sina id:n. Alla
  siffror finns bara i de svenska filerna, så en balansändring görs på
  ett ställe. Text som saknas på engelska visas på svenska.
- **`-kontrollera` och ett test** listar engelska texter som saknas och
  platshållare som inte stämmer med de svenska, till exempel `{namn}`.
- **Tornets hemlighet** översätts lika beslöjat som på svenska och sägs
  inte rakt ut förrän i slutet.
- **AI-namnen** (Ingrid, Sigge och de andra) behålls.
- **Tangenterna:** på engelska heter undanmanövern E (evade), eftersom D
  redan är drink. De svenska tangenterna fungerar i båda språken.

## Ordlista

| Svenska | Engelska |
|---|---|
| Ormtornet | The Serpent Tower |
| vandrare | wanderer |
| spelledaren | the Keeper |
| kammaren | the chamber |
| ormen | the serpent |
| solen | the sun |
| viska, invigd | whisper, initiated |
| hen, henom | they, them |
| Legoknekten | the Sellsword |
| Tjuven | the Thief |
| Örtkvinnan | the Herbwife |
| Kartritaren | the Mapmaker |
| Ormtjusaren | the Snake Charmer |
| Spårfinnaren | the Tracker |
| Lyktbäraren | the Lantern Bearer |
| Pilgrimen | the Pilgrim |
| Blekkryparen | the Pale Crawler |
| Ormynglet | the Serpent Spawn |
| Fjällvakten | the Scale Warden |
| Gallgasten | the Bile Ghast |
| Benväktaren | the Bone Keeper |
| Lätt, Normal, Svår, Mardröm | Easy, Normal, Hard, Nightmare |

## Vad som ska översättas

| Var | Vad | Mängd |
|---|---|---|
| `data/*.json` | rum, fällor, varelser, klasser, föremål, ormen, avslöjandet och loggtexterna | cirka 60 kB |
| `web/index.html`, `app.js`, `map.js` | spelarsidan: introt, namnrutan, klassvalet, frågetecknet, tecknen och kortet | cirka 110 texter |
| `web/spelledare.html` | värdsidan: raden, avsnitten och inställningarna | cirka 140 texter |
| Go-koden | besked och fel till värdsidan (`host.go`, `ingrepp.go`) och raderna i översikten (`gm.go`) | cirka 50 texter |

## Steg

Varje steg blir en commit, och arbetet stannar efter varje steg så att
det kan provas.

1. **Inställningen och strukturen:** Språk på värdsidan, `-sprak`,
   inläsningen av `data/en/` med svenska som reserv, och kontrollen.
   Engelska går att välja men visar svenska tills texterna finns.
2. **Gränssnittets texter flyttas till `data/`:** spelarsidans och
   värdsidans texter, och Go-kodens texter till `texter.json`. Ingenting
   ändras för spelarna.
3. **Spelets innehåll på engelska:** klasser, varelser, fällor, föremål,
   ormen, AI-texterna och loggtexterna.
4. **Gränssnittet på engelska:** spelarsidan och värdsidan, med
   tangenterna.
5. **Avslöjandet och slutet på engelska**, som ett eget steg, så att
   tonen och hemligheten kan granskas noga.

Engelskan bör läsas av någon som spelar, särskilt introt och avslöjandet.
