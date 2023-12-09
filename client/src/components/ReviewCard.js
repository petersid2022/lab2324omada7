import React from 'react';

const ReviewCard = ({ text, stars, date }) => {
    return (
        <div className="bg-gradient-to-br from-white via-white to-dde6d5 px-4 py-2 rounded border border-gray-300 shadow-lg w-64 m-4">
            <div>
                <h1 className="text-xl font-semibold">{stars}/5 ⭐</h1>
                <p className="mt-2 text-gray-600 overflow-hidden w-full">
                    "{text}"
                </p>
                <p className="mt-2 text-gray-600 overflow-hidden w-full">
                    posted on {date}
                </p>
            </div>
        </div>
    );
};

export default ReviewCard;
