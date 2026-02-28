import React from 'react';
import { usePlayer } from '../../context/PlayerContext';
import './BottomPlayer.css';

const BottomPlayer = () => {
    const {
        currentSong,
        isPlaying,
        progress,
        duration,
        volume,
        isMinimized,
        setVolume,
        togglePlay,
        seek,
        toggleMinimized
    } = usePlayer();

    if (!currentSong) return null;

    if (isMinimized) {
        return (
            <div className="player-mini-trigger" onClick={toggleMinimized}>
                <span className="mini-icon">🎵</span>
                <span className="mini-text">Now Playing: {currentSong.title}</span>
                <span className="expand-icon">🔼</span>
            </div>
        );
    }

    const handleProgressChange = (e) => {
        seek(parseFloat(e.target.value));
    };

    const handleVolumeChange = (e) => {
        setVolume(parseFloat(e.target.value));
    };

    return (
        <div className="bottom-player">
            <div className="now-playing">
                <div className="artist-placeholder" style={{ width: '40px', height: '40px', fontSize: '1rem', marginRight: '10px' }}>🎵</div>
                <div className="track-info">
                    <span className="track-name">{currentSong.title}</span>
                    <span className="artist-name">{currentSong.genre || 'Unknown Genre'}</span>
                </div>
            </div>

            <div className="player-middleware">
                <div className="player-controls">
                    <button className="control-btn" title="Previous">⏮</button>
                    <button className="control-btn play-pause" onClick={togglePlay} title="Play/Pause">
                        {isPlaying ? '⏸' : '▶'}
                    </button>
                    <button className="control-btn" title="Next">⏭</button>
                </div>
                <div className="playback-bar-container">
                    <input
                        type="range"
                        className="playback-bar"
                        min="0"
                        max="100"
                        value={progress}
                        onChange={handleProgressChange}
                    />
                </div>
            </div>

            <div className="volume-controls">
                <button className="minimize-btn" onClick={toggleMinimized} title="Minimize">🔽</button>
                <span>{volume > 0 ? (volume > 0.5 ? '🔊' : '🔉') : '🔇'}</span>
                <input
                    type="range"
                    className="volume-slider"
                    min="0"
                    max="1"
                    step="0.01"
                    value={volume}
                    onChange={handleVolumeChange}
                />
            </div>
        </div>
    );
};

export default BottomPlayer;
