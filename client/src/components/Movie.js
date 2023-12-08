import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import '../App.css';

export default function Movie() {
    const [movie, setMovie] = useState({});
    const { title } = useParams();
    const modifiedTitle = title.replace(/ /g, '-');

    useEffect(() => {
        fetch(`http://localhost:1313/api/movies/${modifiedTitle}`)
            .then(response => response.json())
            .then(data => setMovie(data))
            .catch(error => console.error('Error fetching items:', error));
    }, [modifiedTitle]);

    const navigate = useNavigate();

    const goBack = () => {
        navigate(-1);
    };


    return (
        <div>
            <div className="m-4">
                <button onClick={goBack} className="hover:no-underline">
                    <div className="text-xl bg-gray-100 hover:bg-gray-300 text-gray-800 font-semibold px-1 rounded shadow">
                        <span>Go Back</span>
                    </div>
                </button>
            </div>

            <div className="m-4 bg-gray-100 text-gray-800 rounded-lg overflow-hidden shadow-xl border border-gray-700">
                <div className="px-6 pt-3">
                    <div className="font-bold text-2xl">
                        <h1 className="text-gray-800 transition-colors duration-300">
                            {movie.Title}
                        </h1>
                    </div>
                </div>
                <div className="px-6 pb-3">
                    <p className="text-gray-900 text-2xl">Average Rating: {movie.AvgRating}/5</p>
                    <p className="text-gray-900 text-2xl">Πρώτη προβολή: {movie.ReleaseDate}</p>
                    <p className="text-gray-900 text-2xl">Είδος ταινίας: {movie.Genre}</p>
                    <br />
                    <span> Add to watchlist? </span>
                    <input className="mx-2" type="checkbox" />
                </div>
            </div>
        </div>
    );
};
