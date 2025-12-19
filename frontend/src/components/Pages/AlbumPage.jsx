import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useParams, useNavigate } from 'react-router-dom';
import NavBar from '../NavBar/NavBar';
import './Pages.css';

const AlbumPage = () => {
    const { artistId } = useParams();
    const navigate = useNavigate();
    const [albums, setAlbums] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchAlbums = async () => {
            try {
                const token = localStorage.getItem('token');
                const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';

                const response = await axios.get(`${apiUrl}/api/content/artists/${artistId}/albums`, {
                    headers: {
                        'Authorization': `Bearer ${token}`
                    }
                });
                setAlbums(response.data);
            } catch (err) {
                console.error("Error fetching albums:", err);
                setError('Failed to load albums.');
            } finally {
                setLoading(false);
            }
        };

        fetchAlbums();
    }, [artistId]);

    return (
        <div className="page-container">
            <NavBar />
            <div className="content-wrap">
                <button onClick={() => navigate(-1)} className="back-btn">← Back to Artists</button>
                <h1>Albums</h1>
                {loading && <p>Loading albums...</p>}
                {error && <div className="error-msg">{error}</div>}

                <div className="artist-grid">
                    {albums.map(album => (
                        <div key={album.id} className="artist-card" onClick={() => navigate(`/albums/${album.id}/songs`)}>
                            <div className="artist-placeholder">💿</div>
                            <h3>{album.title}</h3>
                            <p>{album.date ? new Date(album.date).getFullYear() : 'N/A'}</p>
                            <p className="sub-text">{album.genre || 'N/A'}</p>
                        </div>
                    ))}
                    {!loading && albums.length === 0 && <p>No albums found for this artist.</p>}
                </div>
            </div>
        </div>
    );
};

export default AlbumPage;
