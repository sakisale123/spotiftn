import React, { useState } from 'react';
import axios from 'axios';
import { useSearchParams, useNavigate } from 'react-router-dom';
import './Auth.css';

const ResetPasswordPage = () => {
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setSuccess('');

        if (password !== confirmPassword) {
            setError('Lozinke se ne poklapaju!');
            return;
        }

        const strongPasswordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]).{8,}$/;
        if (!strongPasswordRegex.test(password)) {
            setError('Lozinka mora imati bar 8 karaktera, jedno veliko slovo, jedno malo slovo, jedan broj i jedan specijalni karakter.');
            return;
        }

        try {
            const token = searchParams.get('token');
            const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';

            await axios.post(`${apiUrl}/api/users/auth/reset-password`, {
                token,
                newPassword: password
            });

            setSuccess('Lozinka uspešno promenjena! Možete se prijaviti.');
            setTimeout(() => navigate('/login'), 2000);
        } catch (err) {
            setError(err.response?.data || 'Greška pri promeni lozinke.');
        }
    };

    return (
        <div className="auth-container">
            <div className="auth-box">
                <h2>Reset Lozinke</h2>

                {error && <div className="alert error">{error}</div>}
                {success && <div className="alert success">{success}</div>}

                <form onSubmit={handleSubmit}>
                    <div className="form-group">
                        <input
                            type="password"
                            placeholder="Nova lozinka"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            required
                        />
                    </div>
                    <div className="form-group">
                        <input
                            type="password"
                            placeholder="Ponovi lozinku"
                            value={confirmPassword}
                            onChange={(e) => setConfirmPassword(e.target.value)}
                            required
                        />
                    </div>
                    <button type="submit" className="auth-btn">Resetuj Lozinku</button>
                </form>
            </div>
        </div>
    );
};

export default ResetPasswordPage;
