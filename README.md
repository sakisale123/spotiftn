# spotiftn

PROCITAJ SVE!
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

RASPORED:

🟩 Nedelja 1: Struktura

Član 1 – Users

inicijalna struktura users servisa
konekcija na MongoDB
model korisnika (id, username, password hash, role)

Član 2 – Content

priprema strukture content servisa (prazan skeleton)

Član 3 – Notifications

priprema strukture notifications servisa (prazan skeleton)

Član 4 – DevOps

inicijalni docker-compose (Mongo + prazni servisi)
Dockerfile za users, content, notifications
README (setup, run instructions)

🟩 Nedelja 2 – Osnovni backend (registracija + CRUD umetnika)
Član 1 – Users:

implementacija: registracija 
bcrypt hashing lozinke
osnovna validacija inputa
kreiranje korisnika u bazi

Član 2 – Content:

modeli:
Artist
Album
Song

endpoint: kreiranje umetnika 

Član 3 – Notifications:

setup baze (ručno punjenje)
endpoint: GET /notifications/{userId} (dummy first version)

Član 4 – DevOps:

Docker Compose doterivanje
integracija sa MongoDB kroz environment varijable
početak frontenda → forma za registraciju

🟩 Nedelja 3 – Login + umetnici/albuma/pesama pregled
Član 1 – Users:

implementacija: login 
JWT generisanje
endpoint za provere tokena
basic OTP skeleton

Član 2 – Content:

endpointi:
GET /artists
GET /artists/{id}/albums
POST /albums (dodavanje albuma)
POST /songs (dodavanje pesama)

Član 3 – Notifications:

definisanje strukture notifikacija
ručno ubacivanje test podataka

Član 4 – Frontend:

UI:
login forma
lista umetnika (poziv Content servisa preko Gateway-a)
Nginx kao API Gateway (routing do users/content/notifications)

🟩 Nedelja 4 – OTP + magični link + prikaz albuma/pesama
Član 1 – Users:

implementacija:
OTP login (kod na email)
Magični link 
validacija inputa 
osnovna kontrola pristupa 

Član 2 – Content:

endpoint:

GET /albums/{id}/songs
validacija ulaza (nazivi, žanrovi, dužine)

Član 3 – Gateway:

završetak Nginx rutiranja
testiranje svih poziva preko gateway-a
logovanje osnovnih stvari

Član 4 – Frontend:

prikaz albuma po umetniku
prikaz pesama po albumu
UI povezivanje sa gateway-em

🟩 Nedelja 5 – Notifications, finalno spajanje backend-a i frontend-a
Član 1 – Users:

poliranje login flow-a
reset lozinke (ako želite)

Član 2 – Content:

testiranje svih CRUD operacija
prečišćavanje modela i endpointa

Član 3 – Notifications:

finalna verzija:
GET /notifications/{userId}
dokumenti u bazi
integracija sa frontendom

Član 4 – Frontend:

stranica „Notifikacije“
stilizovanje minimalno koliko treba za demo
testiranje celog UI flow-a

















