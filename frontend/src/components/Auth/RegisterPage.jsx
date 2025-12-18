import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import './Auth.css'; // Odmah cemo napraviti i CSS

const RegisterPage = () => {
    const navigate = useNavigate();

    // State za podatke forme
    const [formData, setFormData] = useState({
        name: '',
        email: '',
        password: '',
        confirmPassword: ''
    });

    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');

    // Rukovanje promenama u inputima
    const handleChange = (e) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value
        });
    };

    // Slanje forme
    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setSuccess('');

        // 1. Validacija na klijentu (Zahtev 2.18)
        if (formData.password !== formData.confirmPassword) {
            setError("Lozinke se ne poklapaju!");
            return;
        }

        // Provera jake lozinke: min 8, 1 velika, 1 mala, 1 broj, 1 specijalni karakter
        const strongPasswordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/;

        if (!strongPasswordRegex.test(formData.password)) {
            setError("Lozinka mora imati bar 8 karaktera, jedno veliko slovo, jedno malo slovo, jedan broj i jedan specijalni karakter.");
            return;
        }

        try {
            // Čitamo URL iz env varijable ili koristimo Gateway (8080)
            const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';

            // Slanje na Users servis preko Gateway-a
            const response = await axios.post(`${apiUrl}/api/users/register`, {
                name: formData.name,
                email: formData.email,
                password: formData.password
            });

            setSuccess("Registracija uspešna! Možete se prijaviti.");
            // Nakon 2 sekunde prebaci na login (kad ga napravimo)
            setTimeout(() => navigate('/login'), 2000);

        } catch (err) {
            // Prikaz detaljnije greške
            console.error("Registration Error:", err);
            let msg = "Došlo je do greške prilikom registracije.";

            if (err.response) {
                // Server je odgovorio sa statusom van 2xx
                msg = err.response.data?.message || err.response.data || msg;
                if (typeof msg === 'object') msg = JSON.stringify(msg);
            } else if (err.request) {
                // Zahtev je poslat ali nema odgovora (npr. Network Error, CORS)
                msg = `Nema odgovora od servera. Proverite da li je Gateway pokrenut na portu 8080. Detalji: ${err.message}`;
            } else {
                // Greška prilikom podesavanja zahteva
                msg = err.message;
            }

            setError(msg);
        }
    };

    return (
        <div className="auth-container">
            <div className="auth-box">
                <h2>Registracija</h2>

                {error && <div className="alert error">{error}</div>}
                {success && <div className="alert success">{success}</div>}

                <form onSubmit={handleSubmit}>
                    <div className="form-group">
                        <input
                            type="text"
                            name="name"
                            placeholder="Ime i prezime"
                            value={formData.name}
                            onChange={handleChange}
                            required
                        />
                    </div>
                    <div className="form-group">
                        <input
                            type="email"
                            name="email"
                            placeholder="Email adresa"
                            value={formData.email}
                            onChange={handleChange}
                            required
                        />
                    </div>
                    <div className="form-group">
                        <input
                            type="password"
                            name="password"
                            placeholder="Lozinka"
                            value={formData.password}
                            onChange={handleChange}
                            required
                        />
                    </div>
                    <div className="form-group">
                        <input
                            type="password"
                            name="confirmPassword"
                            placeholder="Ponovi lozinku"
                            value={formData.confirmPassword}
                            onChange={handleChange}
                            required
                        />
                    </div>

                    <button type="submit" className="auth-btn">Registruj se</button>
                </form>
                <p>Već imate nalog? <span onClick={() => navigate('/login')} className="link">Prijavi se</span></p>
            </div>
        </div>
    );
};

export default RegisterPage;