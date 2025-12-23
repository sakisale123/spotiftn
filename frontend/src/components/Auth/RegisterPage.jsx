import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import './Auth.css';

const RegisterPage = () => {
    const navigate = useNavigate();

    const [formData, setFormData] = useState({
        name: '',
        email: '',
        password: '',
        confirmPassword: ''
    });

    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');

    const handleChange = (e) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value
        });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setSuccess('');

        if (formData.password !== formData.confirmPassword) {
            setError("Lozinke se ne poklapaju!");
            return;
        }

        const strongPasswordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]).{8,}$/;

        if (!strongPasswordRegex.test(formData.password)) {
            setError("Lozinka mora imati bar 8 karaktera, jedno veliko slovo, jedno malo slovo, jedan broj i jedan specijalni karakter.");
            return;
        }

        try {
            const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';

            const response = await axios.post(`${apiUrl}/api/users/auth/register`, {
                name: formData.name,
                email: formData.email,
                password: formData.password,
                confirmPassword: formData.confirmPassword
            });

            setSuccess("Registracija uspešna! Proverite email za aktivacioni link. (Za development: pogledajte konzolu backend servisa)");

        } catch (err) {
            console.error("Registration Error:", err);
            let msg = "Došlo je do greške prilikom registracije.";

            if (err.response) {
                msg = err.response.data?.message || err.response.data || msg;
                if (typeof msg === 'object') msg = JSON.stringify(msg);
            } else if (err.request) {
                msg = `Nema odgovora od servera. Proverite da li je Gateway pokrenut na portu 8080. Detalji: ${err.message}`;
            } else {
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