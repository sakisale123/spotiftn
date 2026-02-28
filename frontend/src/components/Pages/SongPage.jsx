import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useParams, useNavigate } from 'react-router-dom';
import StarRating from './StarRating';
import { isAdmin, getUserId } from '../../utils/auth';
import { usePlayer } from '../../context/PlayerContext';
import './Pages.css';

const SongPage = () => {
    const { albumId } = useParams();
    const navigate = useNavigate();
    const [songs, setSongs] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const userIsAdmin = isAdmin();
    const userId = getUserId();
    const { playSong, currentSong, isPlaying } = usePlayer();
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';

    const fetchSongRatings = async (songList) => {
        const token = localStorage.getItem('token');
        const updatedSongs = await Promise.all(songList.map(async (song) => {
            try {
                const response = await axios.get(`${apiUrl}/api/ratings/user/song/${song.id}?userId=${userId}`, {
                    headers: { 'Authorization': `Bearer ${token}` }
                });
                return { ...song, userRating: response.data.score || 0, ratingId: response.data.id };
            } catch (err) {
                return { ...song, userRating: 0 };
            }
        }));
        setSongs(updatedSongs);
    };

    useEffect(() => {
        const fetchSongs = async () => {
            try {
                const token = localStorage.getItem('token');
                const response = await axios.get(`${apiUrl}/api/content/albums/${albumId}/songs`, {
                    headers: { 'Authorization': `Bearer ${token}` }
                });
                fetchSongRatings(response.data);
            } catch (err) {
                console.error("Error fetching songs:", err);
                setError('Failed to load songs.');
            } finally {
                setLoading(false);
            }
        };

        fetchSongs();
    }, [albumId]);

    const handleRate = async (songId, score, existingRatingId) => {
        try {
            const token = localStorage.getItem('token');
            if (existingRatingId) {
                await axios.put(`${apiUrl}/api/ratings/${existingRatingId}?userId=${userId}`, { score }, {
                    headers: { 'Authorization': `Bearer ${token}` }
                });
            } else {
                await axios.post(`${apiUrl}/api/ratings?userId=${userId}`, { songId, score }, {
                    headers: { 'Authorization': `Bearer ${token}` }
                });
            }
            // Refresh ratings
            fetchSongRatings(songs);
        } catch (err) {
            console.error("Error rating song:", err);
            alert("Failed to save rating. Check if you are logged in.");
        }
    };

    const formatDuration = (seconds) => {
        const minutes = Math.floor(seconds / 60);
        const remainingSeconds = seconds % 60;
        return `${minutes}:${remainingSeconds < 10 ? '0' : ''}${remainingSeconds}`;
    };

    return (
        <div className="page-container">
            <div className="content-wrap">
                <button onClick={() => navigate(-1)} className="back-btn">← Back to Albums</button>
                <div className="page-header">
                    <h1>Songs</h1>
                    {userIsAdmin && (
                        <button
                            className="btn-primary"
                            onClick={() => navigate(`/songs/create/${albumId}`)}
                        >
                            Create Song
                        </button>
                    )}
                </div>
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
                            <div className="song-actions">
                                <StarRating
                                    initialRating={song.userRating}
                                    onRate={(score) => handleRate(song.id, score, song.ratingId)}
                                />
                                <span className="song-duration">{formatDuration(song.duration)}</span>
                                <button
                                    className={`play-btn ${currentSong?.id === song.id && isPlaying ? 'playing' : ''}`}
                                    onClick={() => playSong(song)}
                                >
                                    {currentSong?.id === song.id && isPlaying ? '⏸' : '▶'}
                                </button>
                            </div>
                        </div>
                    ))}
                    {!loading && songs.length === 0 && <p>No songs found for this album.</p>}
                </div>
            </div>
        </div>
    );
};

export default SongPage;
