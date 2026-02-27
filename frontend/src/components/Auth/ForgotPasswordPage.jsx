import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate, Link } from 'react-router-dom';
import NavBar from '../NavBar/NavBar';
import '../Pages/Pages.css';

const ForgotPasswordPage = () => {
    const [email, setEmail] = useState('');
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState('');
    const [error, setError] = useState('');
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        setMessage('');
        setError('');

        if (!email.trim()) {
            setError('Please enter your email address');
            return;
        }

        setLoading(true);

        try {
            const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';

            await axios.post(`${apiUrl}/api/users/auth/forgot-password`, { email });

            setMessage('If an account exists with this email, you will receive a password reset link shortly.');

            setEmail('');
        } catch (err) {
            console.error('Forgot password error:', err);
            setError('Failed to process request. Please try again.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="page-container">
            <div className="content-wrap" style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                <div className="form-container" style={{ maxWidth: '400px', marginTop: '50px' }}>
                    <h1>Forgot Password</h1>
                    <p style={{ color: '#b3b3b3', marginBottom: '20px' }}>
                        Enter your email address to receive a magic link for password reset.
                    </p>

                    {error && <div className="error-msg">{error}</div>}
                    {message && <div className="success-msg">{message}</div>}

                    <form onSubmit={handleSubmit}>
                        <div className="form-group">
                            <label htmlFor="email">Email Address</label>
                            <input
                                type="email"
                                id="email"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                placeholder="name@example.com"
                                disabled={loading}
                            />
                        </div>

                        <div className="form-actions" style={{ flexDirection: 'column', gap: '15px' }}>
                            <button
                                type="submit"
                                className="btn-primary"
                                style={{ width: '100%' }}
                                disabled={loading}
                            >
                                {loading ? 'Sending...' : 'Send Magic Link'}
                            </button>

                            <Link to="/login" style={{ color: '#b3b3b3', textDecoration: 'underline', fontSize: '0.9rem' }}>
                                Back to Login
                            </Link>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default ForgotPasswordPage;
