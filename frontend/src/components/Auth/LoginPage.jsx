import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import './Auth.css';

const LoginPage = () => {
    const navigate = useNavigate();

    const [formData, setFormData] = useState({
        email: '',
        password: ''
    });

    const [error, setError] = useState('');

    const handleChange = (e) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value
        });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');

        try {
            const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';

            const response = await axios.post(`${apiUrl}/api/users/login`, {
                email: formData.email,
                password: formData.password
            });

            // Store token (basic implementation for now)
            if (response.data.token) {
                localStorage.setItem('token', response.data.token);
                // Redirect to home or dashboard
                // navigate('/'); 
                alert("Login uspesan! Token sacuvan."); // Temporary feedback
            }

        } catch (err) {
            console.error(err);
            const errorMessage = err.response?.data?.message || err.response?.data || "Došlo je do greške prilikom prijave.";
            setError(typeof errorMessage === 'object' ? JSON.stringify(errorMessage) : errorMessage);
        }
    };

    return (
        <div className="auth-container">
            <div className="auth-box">
                <h2>Prijava</h2>

                {error && <div className="alert error">{error}</div>}

                <form onSubmit={handleSubmit}>
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

                    <button type="submit" className="auth-btn">Prijavi se</button>
                </form>
                <p>Nemas nalog? <span onClick={() => navigate('/register')} className="link">Registruj se</span></p>
            </div>
        </div>
    );
};

export default LoginPage;
