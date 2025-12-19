import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import NavBar from '../NavBar/NavBar';
import './Pages.css';

const ArtistPage = () => {
    const [artists, setArtists] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const navigate = useNavigate();

    useEffect(() => {
        const fetchArtists = async () => {
            try {
                const token = localStorage.getItem('token');
                const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';

                // Note: The /api/content/artists endpoint might need Auth header depending on backend implementation
                // Adding it just in case, though Gateway strips it usually unless configured
                const response = await axios.get(`${apiUrl}/api/content/artists`, {
                    headers: {
                        'Authorization': `Bearer ${token}`
                    }
                });
                setArtists(response.data);
            } catch (err) {
                console.error("Error fetching artists:", err);
                setError('Failed to load artists.');
            } finally {
                setLoading(false);
            }
        };

        fetchArtists();
    }, []);

    return (
        <div className="page-container">
            <NavBar />
            <div className="content-wrap">
                <h1>Artists</h1>
                {loading && <p>Loading artists...</p>}
                {error && <div className="error-msg">{error}</div>}

                <div className="artist-grid">
                    {artists.map(artist => (
                        <div key={artist.id} className="artist-card" onClick={() => navigate(`/artists/${artist.id}/albums`)}>
                            <div className="artist-placeholder">🎵</div>
                            <h3>{artist.name}</h3>
                            <p>{artist.genres && artist.genres.length > 0 ? artist.genres.join(', ') : 'N/A'}</p>
                        </div>
                    ))}
                    {!loading && artists.length === 0 && <p>No artists found.</p>}
                </div>
            </div>
        </div>
    );
};

export default ArtistPage;
