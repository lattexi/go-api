# GO API Harjoitustyö

Tämä on REST API, joka on rakennettu Go-kielellä. API jäljittelee aiemmilla kursseilla käytettyä mediapalvelin API:a, mutta on toteutettu Go:lla. API on toteutettu ilman ulkoisia frameworkkeja, käyttäen vain Go:n standardikirjastoja sekä specifejä tarvittavia kirjastoja kuten mysql-ajuria ja JWT-kirjastoa.

## Ominaisuudet

- Käyttäjien luonti ja kirjautuminen joka palauttaa JWT-tokenin
- Pyyntöjen autentikointi JWT-tokenilla (middleware)
- Mediaitemien luonti ja tiedostojen tallennus palvelimelle
- Mediaitemien haku ja niiden tiedostojen tarjoaminen
- Mediaitemien poistaminen ja päivittäminen (vaiheessa)
