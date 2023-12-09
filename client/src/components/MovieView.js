import React, { useState, useEffect } from 'react';
import ReactStars from "react-rating-stars-component";
import { useParams, useNavigate } from 'react-router-dom';
import '../App.css';

export default function MovieView() {
    const [movie, setMovie] = useState({});
    const [reviews, setReview] = useState([]);
    const { title } = useParams();
    const modifiedTitle = title.replace(/ /g, '-');
    const [inputValue, setInputValue] = useState('');
    const [rating, setRating] = useState(0);
    const isAuthenticated = localStorage.getItem('loggedIn') === 'true';

    const handleTextChange = (e) => {
        setInputValue(e.target.value);
    };

    const ratingChanged = (newRating) => {
        setRating(newRating);
    };

    const handleSubmit = (event) => {
        event.preventDefault();

        fetch(`http://localhost:1313/api/movies/add-review/${modifiedTitle}/${rating}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                reviewText: inputValue,
            }),
        })
            .then((response) => response.json())
            .then((data) => {
                console.log('Review and rating submitted successfully:', data);
                setInputValue('');
                setRating(0);
            })
            .catch((error) => {
                console.error('Error submitting review and rating:', error);
            });
    };

    const addedToLiked = (e) => {
        const out = e.target.checked ? 1 : 0;
        console.log(out);
    };

    const addedToWatchlist = (e) => {
        const out = e.target.checked ? 1 : 0;
        console.log(out);
    };

    useEffect(() => {
        fetch(`http://localhost:1313/api/movies/${modifiedTitle}`)
            .then(response => response.json())
            .then(data => setMovie(data))
            .catch(error => console.error('Error fetching movie items:', error));
    }, [modifiedTitle]);

    useEffect(() => {
        fetch(`http://localhost:1313/api/movies/reviews/${modifiedTitle}`)
            .then(response => response.json())
            .then(data => setReview(data))
            .catch(error => console.error('Error fetching review items:', error));
    }, [modifiedTitle]);

    const navigate = useNavigate();

    const goBack = () => {
        navigate(-1);
    };

    return (
        <div>
            <div className="m-4">
                <button onClick={goBack} className="hover:no-underline">
                    <div className="bg-blue-500 hover:bg-blue-700 text-white font-bold px-4 py-2 rounded text-xl">
                        <span>Go Back</span>
                    </div>
                </button>
            </div>

            <div className="m-4 bg-gray-100 text-gray-800 rounded-lg overflow-hidden shadow-xl border border-gray-700">
                <div className="grid grid-cols-2 gap-4">
                    <div className="px-6 py-3">
                        <div className="font-bold text-2xl">
                            <h1 className="text-gray-800 transition-colors duration-300">
                                {movie.Title}
                            </h1>
                        </div>
                        <p className="text-gray-900 text-2xl">Average Rating: {movie.AvgRating}/5</p>
                        <p className="text-gray-900 text-2xl">Πρώτη προβολή: {movie.ReleaseDate}</p>
                        <p className="text-gray-900 text-2xl">Είδος ταινίας: {movie.Genre}</p>
                        <br />
                        {isAuthenticated && (
                            <div>
                                <span> Add to watchlist? </span>
                                <input onChange={addedToWatchlist} className="mx-2" type="checkbox" />
                                <br />
                                <span> Add to liked? </span>
                                <input onChange={addedToLiked} className="mx-2" type="checkbox" />
                                <form onSubmit={handleSubmit}>
                                    <div>
                                        <label>
                                            Enter review:
                                            <input
                                                className="mx-2 border border-black"
                                                type="text"
                                                value={inputValue}
                                                onChange={handleTextChange}
                                            />
                                        </label>
                                    </div>
                                    <div>
                                        <ReactStars
                                            count={5}
                                            onChange={ratingChanged}
                                            size={24}
                                            activeColor="#ffd700"
                                            value={rating}
                                        />
                                    </div>
                                    <button className="border border-black px-2 py-1 hover:bg-red-100" type="submit">Submit</button>
                                </form>
                            </div>
                        )}
                    </div>
                    <div className="px-6 pt-3">
                        <div className="font-bold text-2xl">
                            <h1 className="text-gray-800 transition-colors duration-300">
                                Reviews
                            </h1>
                        </div>
                        <div className="grid grid-cols-3 gap-4 border border-black p-4 mb-4">
                            {reviews.map((review, index) => (
                                <div key={index}>
                                    <p>Rating Stars: {review.RatingStars}</p>
                                    <p>Review Text: {review.ReviewText}</p>
                                    <p>Date Posted: {review.DatePosted}</p>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>
            </div>
        </div >
    );
};
