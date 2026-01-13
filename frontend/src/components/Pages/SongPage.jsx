import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useParams, useNavigate } from 'react-router-dom';
import NavBar from '../NavBar/NavBar';
import './Pages.css';

const SongPage = () => {
    const { albumId } = useParams();
    const navigate = useNavigate();
    const [songs, setSongs] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchSongs = async () => {
            try {
                const token = localStorage.getItem('token');
                const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';

                const response = await axios.get(`${apiUrl}/api/content/albums/${albumId}/songs`, {
                    headers: {
                        'Authorization': `Bearer ${token}`
                    }
                });
                setSongs(response.data);
            } catch (err) {
                console.error("Error fetching songs:", err);
                setError('Failed to load songs.');
            } finally {
                setLoading(false);
            }
        };

        fetchSongs();
    }, [albumId]);

    const formatDuration = (seconds) => {
        const minutes = Math.floor(seconds / 60);
        const remainingSeconds = seconds % 60;
        return `${minutes}:${remainingSeconds < 10 ? '0' : ''}${remainingSeconds}`;
    };

    return (
        <div className="page-container">
            <NavBar />
            <div className="content-wrap">
                <button onClick={() => navigate(-1)} className="back-btn">← Back to Albums</button>
                <h1>Songs</h1>
                {loading && <p>Loading songs...</p>}
                {error && <div className="error-msg">{error}</div>}

                <div className="song-list">
                    {Array.isArray(songs) && songs.map((song, index) => (
                        <div key={song.id} className="song-item">
                            <span className="song-number">{index + 1}</span>
                            <div className="song-info">
                                <h3>{song.title}</h3>
                                <p>{song.genre || 'Unknown Genre'}</p>
                            </div>
                            <span className="song-duration">{formatDuration(song.duration)}</span>
                            <button className="play-btn">▶</button>
                        </div>
                    ))}
                    {!loading && songs.length === 0 && <p>No songs found for this album.</p>}
                </div>
            </div>
        </div>
    );
};

export default SongPage;
