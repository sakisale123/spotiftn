import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import NavBar from '../NavBar/NavBar';
import './Pages.css';

const ArtistCreate = () => {
    const [formData, setFormData] = useState({
        name: '',
        biography: '',
        genres: ''
    });
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');
    const navigate = useNavigate();

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


        if (!formData.name.trim()) {
            setError('Artist name is required');
            return;
        }
        if (!formData.biography.trim()) {
            setError('Biography is required');
            return;
        }
        if (!formData.genres.trim()) {
            setError('At least one genre is required');
            return;
        }

        setLoading(true);

        try {
            const token = localStorage.getItem('token');
            const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';


            const genresArray = formData.genres
                .split(',')
                .map(g => g.trim())
                .filter(g => g.length > 0);

            const payload = {
                name: formData.name.trim(),
                biography: formData.biography.trim(),
                genres: genresArray
            };

            await axios.post(`${apiUrl}/api/content/artists`, payload, {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            });

            setSuccess('Artist created successfully!');
            setTimeout(() => {
                navigate('/artists');
            }, 1500);
        } catch (err) {
            console.error('Error creating artist:', err);
            setError(err.response?.data?.error || 'Failed to create artist');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="page-container">
            <NavBar />
            <div className="content-wrap">
                <div className="form-container">
                    <h1>Create New Artist</h1>

                    {error && <div className="error-msg">{error}</div>}
                    {success && <div className="success-msg">{success}</div>}

                    <form onSubmit={handleSubmit} className="artist-form">
                        <div className="form-group">
                            <label htmlFor="name">Artist Name *</label>
                            <input
                                type="text"
                                id="name"
                                name="name"
                                value={formData.name}
                                onChange={handleChange}
                                placeholder="Enter artist name"
                                disabled={loading}
                            />
                        </div>

                        <div className="form-group">
                            <label htmlFor="biography">Biography *</label>
                            <textarea
                                id="biography"
                                name="biography"
                                value={formData.biography}
                                onChange={handleChange}
                                placeholder="Enter artist biography"
                                rows="5"
                                disabled={loading}
                            />
                        </div>

                        <div className="form-group">
                            <label htmlFor="genres">Genres * (comma-separated)</label>
                            <input
                                type="text"
                                id="genres"
                                name="genres"
                                value={formData.genres}
                                onChange={handleChange}
                                placeholder="e.g., Rock, Pop, Jazz"
                                disabled={loading}
                            />
                            <small className="form-hint">Enter genres separated by commas</small>
                        </div>

                        <div className="form-actions">
                            <button
                                type="button"
                                className="btn-secondary"
                                onClick={() => navigate('/artists')}
                                disabled={loading}
                            >
                                Cancel
                            </button>
                            <button
                                type="submit"
                                className="btn-primary"
                                disabled={loading}
                            >
                                {loading ? 'Creating...' : 'Create Artist'}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default ArtistCreate;
