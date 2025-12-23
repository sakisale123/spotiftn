# spotiftn

PROCIRAJ SVE!
Primer rasporeda rada projekta.
Postuj verzije koje su dodeljenje u go.mod fajlu (za sada samo go kako budete dodajali i ostali clanovi da znaju da ubace).
Procitajte sve izaberite koje cete funkcionalnosti raditi, ako mislite da nije dobro grupni poziv pa da napravimo novo.

----------------------------------------------------------------------------------------------
TEHNOLOGIJE: 

Servisi: GO + Mux,
Baza: MongoDB,
API Gateway: NGINX,
Docker: Docker Compose,
Auth: JWT + bcrypt,
FrontEnd: React + Vite,


-----------------------------------------------------------------------------------------------


🛠 ZAJEDNIČKI ZADATAK (Svi članovi)

Zahtev 2.19 (Informaciona bezbednost): HTTPS između servisa

Svaki servis mora preći sa http na https.

Svi servisi moraju učitati SSL sertifikate (koje generiše Član 4) i pokrenuti TLS server umesto običnog.

Zahtev 2.7: Ugradnja Timeouts i Circuit Breaker mehanizama u HTTP klijente (kada jedan servis poziva drugi).



👤 Član 1 (Users & Subscriptions)

Glavni fokus: Upravljanje pretplatama i bezbednost korisnika.

Novi Mikroservis: Subscriptions (Zahtev 1.10)

Kreiranje servisa koji čuva ko koga prati (User -> Artist/Genre).

Baza: NoSQL po izboru (može Mongo).

Endpointi: POST /subscribe, POST /unsubscribe, GET /subscriptions/{userId}.

HTTPS Implementacija:

Prebacivanje Users servisa na HTTPS.

Konfiguracija TLS-a u Go-u.



🎵 Član 2 (Content & Streaming)

Glavni fokus: Pretraga i serviranje muzike.

Filtriranje i Pretraga (Zahtev 1.8):

Proširivanje endpointa za GET /artists i GET /songs.

Implementacija pretrage po nazivu (Regex ili Text Search u Mongu).

Implementacija filtriranja po žanru.

Reprodukcija (Zahtev 1.7):

Serviranje audio fajlova. Za sada, audio fajl možeš čuvati u folderu unutar kontejnera ili kao Binary u Mongu (GridFS) dok ne uvedemo HDFS (za ocenu 9).

Endpoint: GET /songs/{id}/stream.

HTTPS Implementacija:

Prebacivanje Content servisa na HTTPS.



🔔 Član 3 (Notifications & Ratings)

Glavni fokus: Interakcija korisnika i sinhrona komunikacija.

Novi Mikroservis: Ratings (Zahtev 1.9)

Servis za ocenjivanje pesama (1-5 zvezdica).

Endpointi: POST /rating, GET /rating/{songId}.

Sinhrona komunikacija (Zahtev 2.5):

Ključno: Kada korisnik oceni pesmu, Ratings servis mora sinhrono (HTTP pozivom) pitati Content servis: "Da li pesma sa ID-jem X postoji?".

Ako Content servis ne radi, ocenjivanje ne sme da prođe (ili treba da se aktivira Circuit Breaker).

HTTPS Implementacija:

Prebacivanje Notifications i novog Ratings servisa na HTTPS.



💻 Član 4 (DevOps & Frontend)

Glavni fokus: Infrastruktura bezbednosti i UI funkcionalnosti.

DevOps - Sertifikati & HTTPS (Zahtev 2.19):

Generisanje Self-Signed SSL sertifikata (.crt i .key).

Distribucija sertifikata svim servisima kroz docker-compose (volumes).

Konfiguracija Nginx Gateway-a da radi na portu 443 (HTTPS) i da priča sa servisima preko HTTPS-a.

Frontend - Pretraga & Player:

Search Page: UI sa poljem za pretragu i dropdown-om za žanrove koji gađa Člana 2.

Music Player: Komponenta u dnu ekrana koja pušta zvuk sa endpointa Člana 2.

Frontend - Interakcije:

Subscribe dugme: Na profilu umetnika (gađa Člana 1).

Rating zvezdice: Pored pesme, mogućnost klika na zvezdicu (gađa Člana 3).



📝 Redosled koraka za tim:

Član 4 pravi sertifikate i šalje ih timu (ili ih stavlja u repo).

Svi članovi implementiraju HTTPS u svojim Go serverima.

Spajanje (Merge) da se proveri da li infrastruktura radi sa HTTPS-om.

Svako radi svoje funkcionalnosti (Search, Ratings, Subscriptions).

Frontend integracija.


----------------------------------------------------------------------------------------------

TRELLO KARTICE:

Clan 1 BORIVOJE:
    Subscribing to Artists and Genres


Clan 2 Cigan:
    Content Search Service
    Audio Content Streaming

Clan 3 Spaki Utoka
    Rating songs
    Synchronous Communication
    Timeout Mechanism
    Circuit Breaker

Clan 4 Gazda:
    Implementacija HTTPS infrastrukture
    UI za Pretragu i Filtriranje
    Interactions UI (Rating & Subscription)













