# Admin Testing Guide

## Kreiranje Admin Korisnika

Da bi testirao admin funkcionalnost, potreban ti je korisnik sa `role: "admin"` u bazi.

### Opcija 1: Ručno u MongoDB

Ako već imaš registrovanog korisnika, možeš ga ažurirati direktno u MongoDB:

```javascript
// Konektuj se na MongoDB
use spotiftn_users

// Pronađi svog korisnika i postavi ga kao admina
db.users.updateOne(
  { email: "tvoj-email@example.com" },
  { $set: { role: "admin" } }
)

// Verifikuj
db.users.findOne({ email: "tvoj-email@example.com" })
```

### Opcija 2: Modifikuj Backend za Registraciju

Privremeno možeš modifikovati backend da automatski postavlja role na "admin" pri registraciji.

U `services/users/handlers/auth.go` ili sličnom fajlu, pronađi gde se kreira novi korisnik i dodaj:

```go
newUser := models.User{
    // ... ostala polja
    Role: "admin",  // Dodaj ovu liniju
}
```

**NAPOMENA**: Ovo je samo za testiranje! U produkciji, role se dodeljuje drugačije.

### Opcija 3: Seed Script

Možeš kreirati seed script koji dodaje test admin korisnika.

## Testiranje Admin Funkcionalnosti

1. **Pokreni aplikaciju**:
   ```bash
   # U root direktorijumu projekta
   docker-compose up
   
   # U drugom terminalu, pokreni frontend
   cd frontend
   npm run dev
   ```

2. **Uloguj se kao admin**:
   - Otvori http://localhost:5173
   - Uloguj se sa admin kredencijalima

3. **Testiraj Artist Management**:
   - ✅ Vidi "Create Artist" dugme na `/artists` stranici
   - ✅ Klikni "Create Artist" i kreiraj novog umetnika
   - ✅ Vidi "Edit" dugme na svakoj artist kartici
   - ✅ Klikni "Edit" i izmeni umetnika

4. **Testiraj Album Creation**:
   - ✅ Klikni na umetnika
   - ✅ Vidi "Create Album" dugme
   - ✅ Kreiraj novi album

5. **Testiraj Song Creation**:
   - ✅ Klikni na album
   - ✅ Vidi "Create Song" dugme
   - ✅ Kreiraj novu pesmu

## Verifikacija Non-Admin Korisnika

1. Registruj novog korisnika (bez admin role)
2. Uloguj se
3. **Verifikuj da admin dugmad NISU vidljiva**

## Troubleshooting

### Admin dugmad se ne prikazuju

Proveri JWT token:
1. Otvori Developer Tools (F12)
2. Idi na Console
3. Unesi:
   ```javascript
   const token = localStorage.getItem('token');
   const payload = JSON.parse(atob(token.split('.')[1]));
   console.log(payload);
   ```
4. Proveri da li `role` polje postoji i da li je `"admin"`

### Backend greške

Proveri da li backend servisi rade:
```bash
docker-compose ps
```

Proveri logove:
```bash
docker-compose logs -f content
docker-compose logs -f users
```
