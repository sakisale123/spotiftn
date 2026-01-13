import React from 'react';
import './BottomPlayer.css';

const BottomPlayer = () => {
    return (
        <div className="bottom-player">
            <div className="now-playing">
                <div className="track-info">


                </div>
            </div>

            <div className="player-controls">
                <button className="control-btn" title="Previous">⏮</button>
                <button className="control-btn play-pause" title="Play/Pause">⏯</button>
                <button className="control-btn" title="Next">⏭</button>
            </div>

            <div className="volume-controls">
                <span>🔊</span>
                <div className="volume-bar">
                    <div className="volume-level" style={{ width: '60%' }}></div>
                </div>
            </div>
        </div>
    );
};

export default BottomPlayer;
