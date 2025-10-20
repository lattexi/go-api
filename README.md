# GO API Harjoitustyö

Tämä on REST API, joka on rakennettu Go-kielellä. API jäljittelee aiemmilla kursseilla käytettyä mediapalvelin API:a, mutta on toteutettu Go:lla. API on toteutettu ilman ulkoisia frameworkkeja, käyttäen vain Go:n standardikirjastoja sekä specifejä tarvittavia kirjastoja kuten mysql-ajuria ja JWT-kirjastoa.

## Lähteet ja lähtökohdat

Olin aiemmin opiskellut hieman Go:n syntaksia ja toimintaa. Silmäilin tutoriaaleja jotka olivat materiaaleissa. En kuitenkaan aloittanut niiden pohjalta toteutusta, sillä ne olivat joko vanhoja tai käyttivät ulkoisia frameworkkeja. Sen sijaan katsoin [tämän](https://www.youtube.com/watch?v=eqvDSkuBihs) videon ja aloin sen pohjalta työstämään apia.

## API:n toiminnot

- Käyttäjien luonti ja kirjautuminen joka palauttaa JWT-tokenin
- Mediaitemien luonti ja tiedostojen tallennus palvelimelle
- Mediaitemien haku ja järjestäminen (myös otsikolla haku)
- Mediatiedostojen tarjoaminen
- Mediaitemien poistaminen ja päivittäminen

## Ominaisuudet

- RESTful API periaatteiden mukainen toteutus
- MySQL-tietokanta (mediashare) käyttäjien ja mediaitemien tallennukseen
- Tarvittavat pyynnöt autentikoidaan JWT-tokeneilla (toteutettu middlewarella)
- Kaikki pyynnöt validoidaan
- Oikeudet tarkistetaan (esim. vain omistaja voi muokata tai poistaa mediaitemin)
- Tiedostojen upload rajoitettu 5MB kokoisiin tiedostoihin
- Tiedostot validoidaan MIME-tyypin perusteella (vain kuvat ja videot sallittu)
- Virheiden käsittely ja asianmukaiset HTTP-statuskoodit
- Docker ja Docker Compose ympäristön pystyttämiseen helposti

## Mainittavat asiat

- Käytin validointiin funkktioita jotka vähentävät toistoa
- Toteutin yksinkertaisen middleware-kerroksen
- Käytin tools pakettia apufunktioille (en tiedä onko tämä hyvä käytäntö Go:ssa)

## Ongelmia ja jatkokehitys

Ongelmia tuli vastaan median luomisessa, erityisesti form-datan validoinnissa. Siitä tuli ehkä vähän stuntti ratkaisu ilman kirjastoja, mutta se toimii jotenkuten. Upload ja tiedot olisi voinut käsitellä erikseen kuten alkuperäisessä mediapalvelin API:ssa.

Olisin voinut tehdä API:sta täysin yhteensopivan alkuperäisen mediapalvelin API:n kanssa jolloin vanhaa frontendiä olisi voinut käyttää suoraan, mutta en kerennyt tekemään projektia niin syvällisesti.

Validoinnista ja autentikoinnista huolimatta API:sta saattaa löytyä tietoturva-aukkoja.