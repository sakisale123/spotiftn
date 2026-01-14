import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useSearchParams, useNavigate } from 'react-router-dom';
import './Auth.css';

const ActivationPage = () => {
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();
    const [status, setStatus] = useState('loading');
    const [message, setMessage] = useState('');

    useEffect(() => {
        const activateAccount = async () => {
            const token = searchParams.get('token');

            if (!token) {
                setStatus('error');
                setMessage('Nedostaje activation token.');
                return;
            }

            try {
                const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';
                await axios.get(`${apiUrl}/api/users/auth/confirm?token=${token}`);

                setStatus('success');
                setMessage('Nalog je uspešno aktiviran! Možete se prijaviti.');
                setTimeout(() => navigate('/login'), 3000);
            } catch (err) {
                setStatus('error');
                setMessage(err.response?.data || 'Greška pri aktivaciji naloga.');
            }
        };

        activateAccount();
    }, [searchParams, navigate]);

    return (
        <div className="auth-container">
            <div className="auth-box">
                <h2>Aktivacija Naloga</h2>

                {status === 'loading' && <p>Aktiviram nalog...</p>}

                {status === 'success' && (
                    <div className="alert success">{message}</div>
                )}

                {status === 'error' && (
                    <div className="alert error">{message}</div>
                )}
            </div>
        </div>
    );
};

export default ActivationPage;
