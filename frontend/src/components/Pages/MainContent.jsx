import React from 'react';
import { useSearch } from '../../context/SearchContext';
import { useNavigate } from 'react-router-dom';
import NavBar from '../NavBar/NavBar';
import './Pages.css';

const MainContent = ({ children }) => {
    const { results, loading, isSearching, clearSearch } = useSearch();
    const navigate = useNavigate();

    if (!isSearching) {
        return <>{children}</>;
    }

    return (
        <div className="content-wrap">
            <div className="page-header">
                <h1>Search Results</h1>
                <button className="btn-secondary" onClick={clearSearch}>Clear Search</button>
            </div>

            {loading && <p>Searching...</p>}

            {results && (
                <div className="search-results-container">
                    {results.artists && results.artists.length > 0 && (
                        <section className="result-section">
                            <h2>Artists</h2>
                            <div className="artist-grid">
                                {results.artists.map(artist => (
                                    <div
                                        key={artist.id}
                                        className="artist-card"
                                        onClick={() => {
                                            navigate(`/artists/${artist.id}/albums`);
                                            clearSearch();
                                        }}
                                    >
                                        <div className="artist-placeholder">🎵</div>
                                        <h3>{artist.name}</h3>
                                        <p>{artist.genres?.join(', ')}</p>
                                    </div>
                                ))}
                            </div>
                        </section>
                    )}

                    {results.albums && results.albums.length > 0 && (
                        <section className="result-section">
                            <h2>Albums</h2>
                            <div className="album-grid">
                                {results.albums.map(album => (
                                    <div
                                        key={album.id}
                                        className="album-card"
                                        onClick={() => {
                                            navigate(`/albums/${album.id}/songs`);
                                            clearSearch();
                                        }}
                                    >
                                        <div className="album-placeholder">💿</div>
                                        <h3>{album.title}</h3>
                                        <p>{album.genre}</p>
                                    </div>
                                ))}
                            </div>
                        </section>
                    )}

                    {results.songs && results.songs.length > 0 && (
                        <section className="result-section">
                            <h2>Songs</h2>
                            <div className="song-list-container">
                                <div className="song-grid-header">
                                    <span>#</span>
                                    <span>Title</span>
                                    <span>Genre</span>
                                    <span>Duration</span>
                                </div>
                                {results.songs.map((song, index) => (
                                    <div
                                        key={song.id}
                                        className="song-item"
                                        onClick={() => {
                                            navigate(`/albums/${song.album_id}/songs`);
                                            clearSearch();
                                        }}
                                    >
                                        <span className="song-index">{index + 1}</span>
                                        <div className="song-info">
                                            <span className="song-title">{song.title}</span>
                                        </div>
                                        <span className="song-genre">{song.genre}</span>
                                        <span className="song-duration">{Math.floor(song.duration / 60)}:{(song.duration % 60).toString().padStart(2, '0')}</span>
                                    </div>
                                ))}
                            </div>
                        </section>
                    )}

                    {(!results.artists || results.artists.length === 0) &&
                        (!results.albums || results.albums.length === 0) &&
                        (!results.songs || results.songs.length === 0) && (
                            <p>No results found for your query.</p>
                        )}
                </div>
            )}
        </div>
    );
};

export default MainContent;
