import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate, useParams } from 'react-router-dom';
import { isAdmin } from '../../utils/auth';
import { GENRES } from '../../utils/constants';
import './Pages.css';

const ArtistEdit = () => {
    const { id } = useParams();
    const [formData, setFormData] = useState({
        name: '',
        biography: '',
        genres: []
    });
    const [loading, setLoading] = useState(true);
    const [submitting, setSubmitting] = useState(false);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');
    const navigate = useNavigate();

    useEffect(() => {
        const fetchArtist = async () => {
            try {
                const token = localStorage.getItem('token');
                const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';

                const response = await axios.get(`${apiUrl}/api/content/artists/${id}`, {
                    headers: {
                        'Authorization': `Bearer ${token}`
                    }
                });

                const artist = response.data;
                setFormData({
                    name: artist.name || '',
                    biography: artist.biography || '',
                    genres: Array.isArray(artist.genres) ? artist.genres.map(g => g.toLowerCase()) : []
                });
            } catch (err) {
                console.error('Error fetching artist:', err);
                setError('Failed to load artist data');
            } finally {
                setLoading(false);
            }
        };

        fetchArtist();
    }, [id]);

    const handleChange = (e) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value
        });
    };

    const handleGenreToggle = (genre) => {
        const lowerGenre = genre.toLowerCase();
        const currentGenres = [...formData.genres];
        if (currentGenres.includes(lowerGenre)) {
            setFormData({
                ...formData,
                genres: currentGenres.filter(g => g !== lowerGenre)
            });
        } else {
            setFormData({
                ...formData,
                genres: [...currentGenres, lowerGenre]
            });
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setSuccess('');

        // Validation
        if (!formData.name.trim()) {
            setError('Artist name is required');
            return;
        }
        if (!formData.biography.trim()) {
            setError('Biography is required');
            return;
        }
        if (formData.genres.length === 0) {
            setError('At least one genre is required');
            return;
        }

        setSubmitting(true);

        try {
            const token = localStorage.getItem('token');
            const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';

            const payload = {
                name: formData.name.trim(),
                biography: formData.biography.trim(),
                genres: formData.genres
            };

            await axios.put(`${apiUrl}/api/content/artists/${id}`, payload, {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            });

            setSuccess('Artist updated successfully!');
            setTimeout(() => {
                navigate('/artists');
            }, 1500);
        } catch (err) {
            console.error('Error updating artist:', err);
            setError(err.response?.data?.error || 'Failed to update artist');
        } finally {
            setSubmitting(false);
        }
    };

    if (loading) {
        return (
            <div className="page-container">
                <div className="content-wrap">
                    <p>Loading artist data...</p>
                </div>
            </div>
        );
    }

    return (
        <div className="page-container">
            <div className="content-wrap">
                <div className="form-container">
                    <h1>Edit Artist</h1>

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
                                disabled={submitting}
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
                                disabled={submitting}
                            />
                        </div>

                        <div className="form-group">
                            <label>Genres * (select all that apply)</label>
                            <div className="genre-selection">
                                {GENRES.map(genre => (
                                    <div key={genre} className="checkbox-item">
                                        <input
                                            type="checkbox"
                                            id={`genre-${genre}`}
                                            checked={formData.genres.includes(genre.toLowerCase())}
                                            onChange={() => handleGenreToggle(genre)}
                                            disabled={submitting}
                                        />
                                        <label htmlFor={`genre-${genre}`}>{genre}</label>
                                    </div>
                                ))}
                            </div>
                        </div>

                        <div className="form-actions">
                            <button
                                type="button"
                                className="btn-secondary"
                                onClick={() => navigate('/artists')}
                                disabled={submitting}
                            >
                                Cancel
                            </button>
                            <button
                                type="submit"
                                className="btn-primary"
                                disabled={submitting}
                            >
                                {submitting ? 'Updating...' : 'Update Artist'}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default ArtistEdit;
