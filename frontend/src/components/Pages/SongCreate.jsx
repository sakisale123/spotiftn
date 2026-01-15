import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate, useParams } from 'react-router-dom';
import NavBar from '../NavBar/NavBar';
import './Pages.css';

const SongCreate = () => {
    const { albumId } = useParams(); // Optional: if creating from album page
    const [formData, setFormData] = useState({
        title: '',
        duration: '',
        genre: '',
        album_id: albumId || '',
        artist_ids: []
    });
    const [albums, setAlbums] = useState([]);
    const [artists, setArtists] = useState([]);
    const [selectedArtists, setSelectedArtists] = useState([]);
    const [loading, setLoading] = useState(false);
    const [loadingData, setLoadingData] = useState(true);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');
    const navigate = useNavigate();

    useEffect(() => {
        const fetchData = async () => {
            try {
                const token = localStorage.getItem('token');
                const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';

                // Fetch artists
                const artistsResponse = await axios.get(`${apiUrl}/api/content/artists`, {
                    headers: { 'Authorization': `Bearer ${token}` }
                });
                setArtists(artistsResponse.data || []);

                // Fetch all albums - we need to get them from all artists
                const albumsData = [];
                for (const artist of artistsResponse.data || []) {
                    try {
                        const albumsResponse = await axios.get(
                            `${apiUrl}/api/content/artists/${artist.id}/albums`,
                            { headers: { 'Authorization': `Bearer ${token}` } }
                        );
                        if (Array.isArray(albumsResponse.data)) {
                            albumsData.push(...albumsResponse.data);
                        }
                    } catch (err) {
                        console.error(`Error fetching albums for artist ${artist.id}:`, err);
                    }
                }

                // Remove duplicates based on album id
                const uniqueAlbums = Array.from(
                    new Map(albumsData.map(album => [album.id, album])).values()
                );
                setAlbums(uniqueAlbums);

            } catch (err) {
                console.error('Error fetching data:', err);
                setError('Failed to load data');
            } finally {
                setLoadingData(false);
            }
        };

        fetchData();
    }, []);

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

        // Validation
        if (!formData.title.trim()) {
            setError('Song title is required');
            return;
        }
        if (!formData.duration || formData.duration <= 0) {
            setError('Valid duration is required');
            return;
        }
        if (!formData.genre.trim()) {
            setError('Genre is required');
            return;
        }
        if (!formData.album_id) {
            setError('Album must be selected');
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
                duration: parseInt(formData.duration),
                genre: formData.genre.trim(),
                album_id: formData.album_id,
                artist_ids: selectedArtists,
                audio_url: '' // For now, empty - will be handled later with file upload
            };

            await axios.post(`${apiUrl}/api/content/songs`, payload, {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            });

            setSuccess('Song created successfully!');
            setTimeout(() => {
                navigate(`/albums/${formData.album_id}/songs`);
            }, 1500);
        } catch (err) {
            console.error('Error creating song:', err);
            setError(err.response?.data?.error || 'Failed to create song');
        } finally {
            setLoading(false);
        }
    };

    if (loadingData) {
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
                    <h1>Create New Song</h1>

                    {error && <div className="error-msg">{error}</div>}
                    {success && <div className="success-msg">{success}</div>}

                    <form onSubmit={handleSubmit} className="song-form">
                        <div className="form-group">
                            <label htmlFor="title">Song Title *</label>
                            <input
                                type="text"
                                id="title"
                                name="title"
                                value={formData.title}
                                onChange={handleChange}
                                placeholder="Enter song title"
                                disabled={loading}
                            />
                        </div>

                        <div className="form-group">
                            <label htmlFor="duration">Duration (seconds) *</label>
                            <input
                                type="number"
                                id="duration"
                                name="duration"
                                value={formData.duration}
                                onChange={handleChange}
                                placeholder="e.g., 180"
                                min="1"
                                disabled={loading}
                            />
                            <small className="form-hint">Duration in seconds (e.g., 180 for 3 minutes)</small>
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
                            <label htmlFor="album_id">Album *</label>
                            <select
                                id="album_id"
                                name="album_id"
                                value={formData.album_id}
                                onChange={handleChange}
                                disabled={loading}
                            >
                                <option value="">Select an album</option>
                                {albums.map(album => (
                                    <option key={album.id} value={album.id}>
                                        {album.title}
                                    </option>
                                ))}
                            </select>
                            {albums.length === 0 && (
                                <p className="form-hint">No albums available. Please create an album first.</p>
                            )}
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
                                {loading ? 'Creating...' : 'Create Song'}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default SongCreate;
