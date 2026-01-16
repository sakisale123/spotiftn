import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate, useParams } from 'react-router-dom';
import NavBar from '../NavBar/NavBar';
import './Pages.css';

const AlbumCreate = () => {
    const { artistId } = useParams(); 
    const [formData, setFormData] = useState({
        title: '',
        date: '',
        genre: '',
        artist_ids: []
    });
    const [artists, setArtists] = useState([]);
    const [selectedArtists, setSelectedArtists] = useState([]);
    const [loading, setLoading] = useState(false);
    const [loadingArtists, setLoadingArtists] = useState(true);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');
    const navigate = useNavigate();

    useEffect(() => {
        const fetchArtists = async () => {
            try {
                const token = localStorage.getItem('token');
                const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';

                const response = await axios.get(`${apiUrl}/api/content/artists`, {
                    headers: {
                        'Authorization': `Bearer ${token}`
                    }
                });
                setArtists(response.data || []);

                if (artistId) {
                    setSelectedArtists([artistId]);
                }
            } catch (err) {
                console.error('Error fetching artists:', err);
                setError('Failed to load artists');
            } finally {
                setLoadingArtists(false);
            }
        };

        fetchArtists();
    }, [artistId]);

    const handleChange = (e) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value
        });
    };

    const handleArtistToggle = (artistId) => {
        if (selectedArtists.includes(artistId)) {
            setSelectedArtists(selectedArtists.filter(id => id !== artistId));
        } else {
            setSelectedArtists([...selectedArtists, artistId]);
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setSuccess('');

        if (!formData.title.trim()) {
            setError('Album title is required');
            return;
        }
        if (!formData.date) {
            setError('Release date is required');
            return;
        }
        if (!formData.genre.trim()) {
            setError('Genre is required');
            return;
        }
        if (selectedArtists.length === 0) {
            setError('At least one artist must be selected');
            return;
        }

        setLoading(true);

        try {
            const token = localStorage.getItem('token');
            const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';

            const payload = {
                title: formData.title.trim(),
                date: new Date(formData.date).toISOString(),
                genre: formData.genre.trim(),
                artist_ids: selectedArtists
            };

            await axios.post(`${apiUrl}/api/content/albums`, payload, {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            });

            setSuccess('Album created successfully!');
            setTimeout(() => {
                if (artistId) {
                    navigate(`/artists/${artistId}/albums`);
                } else {
                    navigate('/artists');
                }
            }, 1500);
        } catch (err) {
            console.error('Error creating album:', err);
            setError(err.response?.data?.error || 'Failed to create album');
        } finally {
            setLoading(false);
        }
    };

    if (loadingArtists) {
        return (
            <div className="page-container">
                <NavBar />
                <div className="content-wrap">
                    <p>Loading...</p>
                </div>
            </div>
        );
    }

    return (
        <div className="page-container">
            <NavBar />
            <div className="content-wrap">
                <div className="form-container">
                    <h1>Create New Album</h1>

                    {error && <div className="error-msg">{error}</div>}
                    {success && <div className="success-msg">{success}</div>}

                    <form onSubmit={handleSubmit} className="album-form">
                        <div className="form-group">
                            <label htmlFor="title">Album Title *</label>
                            <input
                                type="text"
                                id="title"
                                name="title"
                                value={formData.title}
                                onChange={handleChange}
                                placeholder="Enter album title"
                                disabled={loading}
                            />
                        </div>

                        <div className="form-group">
                            <label htmlFor="date">Release Date *</label>
                            <input
                                type="date"
                                id="date"
                                name="date"
                                value={formData.date}
                                onChange={handleChange}
                                disabled={loading}
                            />
                        </div>

                        <div className="form-group">
                            <label htmlFor="genre">Genre *</label>
                            <input
                                type="text"
                                id="genre"
                                name="genre"
                                value={formData.genre}
                                onChange={handleChange}
                                placeholder="e.g., Rock, Pop, Jazz"
                                disabled={loading}
                            />
                        </div>

                        <div className="form-group">
                            <label>Artists * (select at least one)</label>
                            <div className="artist-selection">
                                {artists.map(artist => (
                                    <div key={artist.id} className="checkbox-item">
                                        <input
                                            type="checkbox"
                                            id={`artist-${artist.id}`}
                                            checked={selectedArtists.includes(artist.id)}
                                            onChange={() => handleArtistToggle(artist.id)}
                                            disabled={loading}
                                        />
                                        <label htmlFor={`artist-${artist.id}`}>{artist.name}</label>
                                    </div>
                                ))}
                            </div>
                            {artists.length === 0 && (
                                <p className="form-hint">No artists available. Please create an artist first.</p>
                            )}
                        </div>

                        <div className="form-actions">
                            <button
                                type="button"
                                className="btn-secondary"
                                onClick={() => navigate(-1)}
                                disabled={loading}
                            >
                                Cancel
                            </button>
                            <button
                                type="submit"
                                className="btn-primary"
                                disabled={loading}
                            >
                                {loading ? 'Creating...' : 'Create Album'}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default AlbumCreate;
