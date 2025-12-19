import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import NavBar from '../NavBar/NavBar';
import './Auth.css';

const ChangePasswordPage = () => {
    const navigate = useNavigate();
    const [formData, setFormData] = useState({
        oldPassword: '',
        newPassword: '',
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

        if (formData.newPassword !== formData.confirmPassword) {
            setError('Lozinke se ne poklapaju!');
            return;
        }

        const strongPasswordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]).{8,}$/;
        if (!strongPasswordRegex.test(formData.newPassword)) {
            setError('Lozinka mora imati bar 8 karaktera, jedno veliko slovo, jedno malo slovo, jedan broj i jedan specijalni karakter.');
            return;
        }

        try {
            const token = localStorage.getItem('token');
            const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';

            await axios.post(`${apiUrl}/api/users/auth/change-password`, {
                oldPassword: formData.oldPassword,
                newPassword: formData.newPassword
            }, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            setSuccess('Lozinka uspešno promenjena!');
            setTimeout(() => navigate('/artists'), 2000);
        } catch (err) {
            setError(err.response?.data || 'Greška pri promeni lozinke.');
        }
    };

    return (
        <div className="page-container">
            <NavBar />
            <div className="auth-container">
                <div className="auth-box">
                    <h2>Promena Lozinke</h2>

                    {error && <div className="alert error">{error}</div>}
                    {success && <div className="alert success">{success}</div>}

                    <form onSubmit={handleSubmit}>
                        <div className="form-group">
                            <input
                                type="password"
                                name="oldPassword"
                                placeholder="Stara lozinka"
                                value={formData.oldPassword}
                                onChange={handleChange}
                                required
                            />
                        </div>
                        <div className="form-group">
                            <input
                                type="password"
                                name="newPassword"
                                placeholder="Nova lozinka"
                                value={formData.newPassword}
                                onChange={handleChange}
                                required
                            />
                        </div>
                        <div className="form-group">
                            <input
                                type="password"
                                name="confirmPassword"
                                placeholder="Ponovi novu lozinku"
                                value={formData.confirmPassword}
                                onChange={handleChange}
                                required
                            />
                        </div>
                        <button type="submit" className="auth-btn">Promeni Lozinku</button>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default ChangePasswordPage;
