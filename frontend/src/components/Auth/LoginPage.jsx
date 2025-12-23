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
    const [otp, setOtp] = useState('');
    const [step, setStep] = useState(1); 

    const [error, setError] = useState('');

    const handleChange = (e) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value
        });
    };

    const handleLoginStep1 = async (e) => {
        e.preventDefault();
        setError('');

        try {
            const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';

            await axios.post(`${apiUrl}/api/users/auth/login`, {
                email: formData.email,
                password: formData.password
            });

            setStep(2);

        } catch (err) {
            console.error(err);
            const errorMessage = err.response?.data?.message || err.response?.data || "Došlo je do greške prilikom prijave.";
            setError(typeof errorMessage === 'object' ? JSON.stringify(errorMessage) : errorMessage);
        }
    };

    const handleLoginStep2 = async (e) => {
        e.preventDefault();
        setError('');

        try {
            const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';

            const response = await axios.post(`${apiUrl}/api/users/auth/verify-otp`, {
                email: formData.email,
                otp: otp
            });

            if (response.data.token) {
                localStorage.setItem('token', response.data.token);
                navigate('/artists', { replace: true });
            }

        } catch (err) {
            console.error(err);
            const errorMessage = err.response?.data?.message || err.response?.data || "Neispravan OTP kod.";
            setError(typeof errorMessage === 'object' ? JSON.stringify(errorMessage) : errorMessage);
        }
    };

    return (
        <div className="auth-container">
            <div className="auth-box">
                <h2>{step === 1 ? 'Prijava' : 'Unesite OTP Kod'}</h2>

                {error && <div className="alert error">{error}</div>}

                {step === 1 && (
                    <form onSubmit={handleLoginStep1}>
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
                        <button type="submit" className="auth-btn">Dalje (Posalji OTP)</button>
                    </form>
                )}

                {step === 2 && (
                    <form onSubmit={handleLoginStep2}>
                        <div className="form-group">
                            <p style={{ marginBottom: '1rem', color: '#b3b3b3', fontSize: '0.9rem' }}>Kod je poslat na vašu email adresu. Molimo proverite vaš Inbox.</p>
                            <input
                                type="text"
                                name="otp"
                                placeholder="Unesite OTP kod"
                                value={otp}
                                onChange={(e) => setOtp(e.target.value)}
                                required
                            />
                        </div>
                        <button type="submit" className="auth-btn">Potvrdi i Prijavi se</button>
                        <button type="button" onClick={() => setStep(1)} className="link" style={{ background: 'none', border: 'none', marginTop: '10px' }}>Nazad</button>
                    </form>
                )}

                {step === 1 && (
                    <>
                        <p>Nemas nalog? <span onClick={() => navigate('/register')} className="link">Registruj se</span></p>
                        <p><span onClick={() => navigate('/forgot-password')} className="link">Zaboravljena lozinka?</span></p>
                    </>
                )}
            </div>
        </div>
    );
};

export default LoginPage;
