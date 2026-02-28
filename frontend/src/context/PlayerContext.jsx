import React, { createContext, useContext, useState, useRef, useEffect } from 'react';

const PlayerContext = createContext();

export const usePlayer = () => useContext(PlayerContext);

export const PlayerProvider = ({ children }) => {
    const [currentSong, setCurrentSong] = useState(null);
    const [isPlaying, setIsPlaying] = useState(false);
    const [volume, setVolume] = useState(0.5);
    const [progress, setProgress] = useState(0);
    const [duration, setDuration] = useState(0);
    const [isMinimized, setIsMinimized] = useState(false);
    const audioRef = useRef(new Audio());

    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';

    useEffect(() => {
        const audio = audioRef.current;

        const updateProgress = () => {
            setProgress((audio.currentTime / audio.duration) * 100 || 0);
        };

        const onLoadedMetadata = () => {
            setDuration(audio.duration);
        };

        const onEnded = () => {
            setIsPlaying(false);
            setProgress(0);
        };

        audio.addEventListener('timeupdate', updateProgress);
        audio.addEventListener('loadedmetadata', onLoadedMetadata);
        audio.addEventListener('ended', onEnded);

        return () => {
            audio.removeEventListener('timeupdate', updateProgress);
            audio.removeEventListener('loadedmetadata', onLoadedMetadata);
            audio.removeEventListener('ended', onEnded);
        };
    }, []);

    useEffect(() => {
        audioRef.current.volume = volume;
    }, [volume]);

    const playSong = (song) => {
        if (currentSong?.id === song.id) {
            togglePlay();
            return;
        }

        setCurrentSong(song);
        audioRef.current.src = `${apiUrl}/api/content/songs/${song.id}/stream`;
        audioRef.current.play()
            .then(() => setIsPlaying(true))
            .catch(err => console.error("Playback failed:", err));
    };

    const togglePlay = () => {
        if (!currentSong) return;

        if (isPlaying) {
            audioRef.current.pause();
        } else {
            audioRef.current.play()
                .catch(err => console.error("Playback failed:", err));
        }
        setIsPlaying(!isPlaying);
    };

    const seek = (percent) => {
        const time = (percent / 100) * audioRef.current.duration;
        audioRef.current.currentTime = time;
        setProgress(percent);
    };

    const toggleMinimized = () => {
        setIsMinimized(!isMinimized);
    };

    return (
        <PlayerContext.Provider value={{
            currentSong,
            isPlaying,
            volume,
            progress,
            duration,
            isMinimized,
            setVolume,
            playSong,
            togglePlay,
            seek,
            toggleMinimized
        }}>
            {children}
        </PlayerContext.Provider>
    );
};
