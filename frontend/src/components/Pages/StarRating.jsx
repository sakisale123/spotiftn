import React, { useState } from 'react';
import './Pages.css';

const StarRating = ({ initialRating, onRate }) => {
    const [hover, setHover] = useState(0);

    return (
        <div className="star-rating">
            {[...Array(5)].map((star, index) => {
                const ratingValue = index + 1;
                return (
                    <button
                        type="button"
                        key={index}
                        className={ratingValue <= (hover || initialRating) ? "on" : "off"}
                        onClick={() => onRate(ratingValue)}
                        onMouseEnter={() => setHover(ratingValue)}
                        onMouseLeave={() => setHover(0)}
                    >
                        <span className="star">&#9733;</span>
                    </button>
                );
            })}
        </div>
    );
};

export default StarRating;
