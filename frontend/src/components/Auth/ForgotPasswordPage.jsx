import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import './Auth.css';

const ForgotPasswordPage = () => {
    const navigate = useNavigate();
    const [email, setEmail] = useState('');
    const [message, setMessage] = useState('');
    const [error, setError] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setMessage('');

        try {
            const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';
            await axios.post(`${apiUrl}/api/users/auth/forgot-password`, { email });

            setMessage('Reset link je poslat na vašu email adresu. Proverite backend konzolu.');
        } catch (err) {
            setError(err.response?.data || 'Greška pri slanju reset linka.');
        }
    };

    return (
        <div className="auth-container">
            <div className="auth-box">
                <h2>Zaboravljena Lozinka</h2>

                {error && <div className="alert error">{error}</div>}
                {message && <div className="alert success">{message}</div>}

                <form onSubmit={handleSubmit}>
                    <div className="form-group">
                        <input
                            type="email"
                            placeholder="Email adresa"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            required
                        />
                    </div>
                    <button type="submit" className="auth-btn">Pošalji Reset Link</button>
                </form>

                <p>
                    <span onClick={() => navigate('/login')} className="link">
                        Nazad na prijavu
                    </span>
                </p>
            </div>
        </div>
    );
};

export default ForgotPasswordPage;
