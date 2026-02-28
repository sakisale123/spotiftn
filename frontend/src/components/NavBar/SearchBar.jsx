import React, { useState, useEffect, useRef } from 'react';
import { GENRES } from '../../utils/constants';
import './NavBar.css';

const SearchBar = ({ onSearch }) => {
    const [query, setQuery] = useState('');
    const [showFilters, setShowFilters] = useState(false);
    const [filters, setFilters] = useState({
        types: ['artist', 'album', 'song'],
        genres: []
    });

    const filterRef = useRef(null);
    const availableGenres = GENRES;

    useEffect(() => {
        const timer = setTimeout(() => {
            onSearch({ query, ...filters });
        }, 500);
        return () => clearTimeout(timer);
    }, [query, filters, onSearch]);

    useEffect(() => {
        const handleClickOutside = (event) => {
            if (filterRef.current && !filterRef.current.contains(event.target)) {
                setShowFilters(false);
            }
        };
        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    const handleTypeToggle = (type) => {
        setFilters(prev => ({
            ...prev,
            types: prev.types.includes(type)
                ? prev.types.filter(t => t !== type)
                : [...prev.types, type]
        }));
    };

    const handleGenreToggle = (genre) => {
        const lowerGenre = genre.toLowerCase();
        setFilters(prev => ({
            ...prev,
            genres: prev.genres.includes(lowerGenre)
                ? prev.genres.filter(g => g !== lowerGenre)
                : [...prev.genres, lowerGenre]
        }));
    };

    return (
        <div className="search-container">
            <div className="search-bar-wrapper">
                <input
                    type="text"
                    className="search-input"
                    placeholder="Search for artists, songs, or albums..."
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                />
                <button
                    className={`filter-toggle-btn ${showFilters ? 'active' : ''}`}
                    onClick={() => setShowFilters(!showFilters)}
                >
                    <svg viewBox="0 0 24 24" width="18" height="18" stroke="currentColor" strokeWidth="2" fill="none">
                        <line x1="4" y1="21" x2="4" y2="14" /><line x1="4" y1="10" x2="4" y2="3" />
                        <line x1="12" y1="21" x2="12" y2="12" /><line x1="12" y1="8" x2="12" y2="3" />
                        <line x1="20" y1="21" x2="20" y2="16" /><line x1="20" y1="12" x2="20" y2="3" />
                        <line x1="2" y1="14" x2="6" y2="14" /><line x1="10" y1="8" x2="14" y2="8" />
                        <line x1="18" y1="16" x2="22" y2="16" />
                    </svg>
                    <span>Filters</span>
                </button>
            </div>

            {showFilters && (
                <div className="filter-popup" ref={filterRef}>
                    <div className="filter-section">
                        <h4>Search In</h4>
                        <div className="filter-options">
                            {['artist', 'album', 'song'].map(type => (
                                <label key={type} className="filter-label">
                                    <input
                                        type="checkbox"
                                        checked={filters.types.includes(type)}
                                        onChange={() => handleTypeToggle(type)}
                                    />
                                    {type.charAt(0).toUpperCase() + type.slice(1)}s
                                </label>
                            ))}
                        </div>
                    </div>
                    <div className="filter-section">
                        <h4>Genres</h4>
                        <div className="filter-options genre-grid">
                            {availableGenres.map(genre => (
                                <label key={genre} className="filter-label">
                                    <input
                                        type="checkbox"
                                        checked={filters.genres.includes(genre.toLowerCase())}
                                        onChange={() => handleGenreToggle(genre)}
                                    />
                                    {genre}
                                </label>
                            ))}
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default SearchBar;
