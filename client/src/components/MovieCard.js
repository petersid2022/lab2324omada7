import React from 'react';
import { Link } from 'react-router-dom';
import '../App.css';

const MovieCard = (props) => {
    const movie = props.movie;

    const movieTitle = movie.Title;
    const modifiedTitle = movieTitle.replace(/ /g, '-');

    return (
        <div className="m-4 hover:bg-gray-200 bg-gray-100 text-gray-800 max-w-xs rounded-lg overflow-hidden shadow-xl">
            <Link
                to={`/movies/${modifiedTitle}`}
                className="text-gray-800 hover:text-indigo-500 transition-colors duration-300"
                style={{ textDecoration: "none" }}
            >

                <img
                    src="https://images.unsplash.com/photo-1595769816263-9b910be24d5f?q=80&w=2079&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
                    alt="Movies"
                    height="200"
                />
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
                </div>
            </Link>
        </div>
    );
};

export default MovieCard;
