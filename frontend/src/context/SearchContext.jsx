import React, { createContext, useState, useContext, useCallback } from 'react';
import axios from 'axios';

const SearchContext = createContext();

export const useSearch = () => useContext(SearchContext);

export const SearchProvider = ({ children }) => {
    const [query, setQuery] = useState('');
    const [results, setResults] = useState(null);
    const [loading, setLoading] = useState(false);
    const [isSearching, setIsSearching] = useState(false);

    const performSearch = useCallback(async (searchParams) => {
        const { query, types, genres } = searchParams;
        setQuery(query);

        if (!query && genres.length === 0) {
            setResults(null);
            setIsSearching(false);
            return;
        }

        setLoading(true);
        setIsSearching(true);

        try {
            const token = localStorage.getItem('token');
            const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';

            let url = `${apiUrl}/api/content/search?q=${encodeURIComponent(query)}`;

            if (types && types.length > 0) {
                types.forEach(t => url += `&type=${t}`);
            }
            if (genres && genres.length > 0) {
                genres.forEach(g => url += `&genre=${g}`);
            }

            const response = await axios.get(url, {
                headers: { 'Authorization': `Bearer ${token}` }
            });

            setResults(response.data);
        } catch (err) {
            console.error("Search failed:", err);
            setResults({ artists: [], albums: [], songs: [] });
        } finally {
            setLoading(false);
        }
    }, []);

    const clearSearch = () => {
        setQuery('');
        setResults(null);
        setIsSearching(false);
    };

    return (
        <SearchContext.Provider value={{ query, results, loading, isSearching, performSearch, clearSearch }}>
            {children}
        </SearchContext.Provider>
    );
};
