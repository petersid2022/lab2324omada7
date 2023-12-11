import React, { useState, useEffect } from 'react';
import ReactStars from "react-rating-stars-component";
import ReviewCard from './ReviewCard';
import { TiArrowBack } from 'react-icons/ti';
import { useParams, useNavigate } from 'react-router-dom';
import { Link } from 'react-router-dom';
import '../App.css';

export default function MovieView() {
    const [movie, setMovie] = useState({});
    const [staff, setStaff] = useState({});
    const [reviews, setReview] = useState([]);
    const { title } = useParams();
    const modifiedTitle = title.replace(/ /g, '-');
    const [inputValue, setInputValue] = useState('');
    const [rating, setRating] = useState(0);
    const isAuthenticated = localStorage.getItem('loggedIn') === 'true';
    const username = localStorage.getItem('username');
    const [watchlistStatus, setWatchlistStatus] = useState('');
    const [likedStatus, setLikedStatus] = useState('');
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        setLoading(true);

        fetch(`http://localhost:1313/watchlistStatus/${modifiedTitle}/${username}`)
            .then(response => response.json())
            .then(data => {
                const newWatchlistStatus = data.data;
                setWatchlistStatus(newWatchlistStatus);
            })
            .catch(error => {
                console.error('Error fetching watchlist status:', error);
            });

        fetch(`http://localhost:1313/likedStatus/${modifiedTitle}/${username}`)
            .then(response => response.json())
            .then(data => {
                const newLikedStatus = data.data;
                setLikedStatus(newLikedStatus);
            })
            .catch(error => {
                console.error('Error fetching liked status:', error);
            })
            .finally(() => {
                setLoading(false);
            });
    }, [modifiedTitle, username]);

    const handleWatchlistButtonClick = () => {
        const newWatchlistStatus = watchlistStatus === 'added' ? 'not added' : 'added';

        const watchlistPayload = {
            MovieID: modifiedTitle,
            Username: username,
        };

        fetch('http://localhost:1313/api/watchlist', {
            method: 'POST',
            crossDomain: true,
            headers: {
                'Content-Type': 'application/json',
                Accept: 'application/json',
                'Access-Control-Allow-Origin': '*',
            },
            body: JSON.stringify(watchlistPayload),
        })
            .then(response => response.json())
            .then(data => {
                console.log('Watchlist status updated:', data);
                setWatchlistStatus(newWatchlistStatus);
            })
            .catch(error => console.error('Error updating watchlist status:', error));
    };

    const handleLikedButtonClick = () => {
        const newLikedStatus = likedStatus === 'liked' ? 'not liked' : 'liked';

        const likedPayload = {
            MovieID: modifiedTitle,
            Username: username,
        };

        fetch('http://localhost:1313/api/liked', {
            method: 'POST',
            crossDomain: true,
            headers: {
                'Content-Type': 'application/json',
                Accept: 'application/json',
                'Access-Control-Allow-Origin': '*',
            },
            body: JSON.stringify(likedPayload),
        })
            .then(response => response.json())
            .then(data => {
                console.log('Liked status updated:', data);
                setLikedStatus(newLikedStatus);
            })
            .catch(error => console.error('Error updating liked status:', error));
    };

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
            crossDomain: true,
            headers: {
                "Content-Type": "application/json",
                Accept: "application/json",
                "Access-Control-Allow-Origin": "*",
            },
            body: JSON.stringify({
                reviewText: inputValue,
                userName: username,
            }),
        }).then((response) => response.json()).then((data) => {
            console.log('Review and rating submitted successfully:', data);
            setInputValue('');
            setRating(0);
        }).catch((error) => { console.error('Error submitting review and rating:', error); });
    };

    useEffect(() => {
        fetch(`http://localhost:1313/api/movies/staff/${modifiedTitle}`)
            .then(response => response.json())
            .then(data => setStaff(data))
            .catch(error => console.error('Error fetching movie items:', error));
    }, [modifiedTitle]);

    let members;
    if (staff && staff.length > 0) {
        members = staff;
    } else {
        console.log('staff object empty or undefined');
    }

    const generateRoleUrl = (role) => {
        if (role) {
            const lowercasedRole = role.toLowerCase();
            const roles = lowercasedRole + 's';
            return `/${roles}`;
        }
        return '/';
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

    if (loading) {
        return <p>Loading...</p>;
    }

    return (
        <div className="max-w-fit">
            <div className="mr-4 ml-4 mt-4">
                <button onClick={goBack} className="hover:no-underline">
                    <div className="flex fitems-center bg-white hover:bg-gray-100 text-gray-800 font-semibold py-1.5 px-3 rounded shadow text-xl">
                        <TiArrowBack className="mt-1 mr-1" />
                        <span>Go Back</span>
                    </div>
                </button>
            </div>

            <div className="flex justify-center w-screen">
                <div className="m-4 bg-white text-gray-800 rounded-lg overflow-hidden shadow-2xl">
                    <div className="px-6 py-3">
                        <div className="flex items-center mt-2">
                            <img
                                className="object-cover h-48 w-48 mr-4 rounded"
                                src="https://images.unsplash.com/photo-1595769816263-9b910be24d5f?q=80&w=2079&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
                                alt="Movie"
                            />
                            <div className="flex flex-col">
                                <div className="font-bold text-2xl">
                                    <h1 className="text-gray-800 transition-colors duration-300">
                                        {movie.Title}
                                    </h1>
                                </div>
                                {members
                                    ? members.map((member, index) => (
                                        <Link key={index} to={`${generateRoleUrl(member.role)}/${member.type_id}`}>
                                            <p className="text-gray-900 text-2xl">{member.role}: {member.name}</p>
                                        </Link>
                                    ))
                                    : <p className="text-gray-900 text-2xl">N/A</p>
                                }
                                <p className="text-gray-900 text-2xl">Μέση Βαθμολογία: {movie.AvgRating}/5</p>
                                <p className="text-gray-900 text-2xl">Ημερομηνία πρώτης προβολής: {movie.ReleaseDate}</p>
                                <p className="text-gray-900 text-2xl">Είδος ταινίας: {movie.Genre}</p>
                            </div>
                        </div>
                        <br />
                        {isAuthenticated && (
                            <div>
                                <div className="flex items-center">
                                    <input
                                        type="button"
                                        className="bg-white hover:bg-gray-100 text-gray-800 font-semibold py-1 px-2 border border-gray-400 rounded shadow"
                                        onClick={handleWatchlistButtonClick}
                                        value={watchlistStatus === 'added' ? 'Remove from Watchlist' : 'Add to Watchlist'}
                                    />
                                    <input
                                        type="button"
                                        className="ml-2 bg-white hover:bg-gray-100 text-gray-800 font-semibold py-1 px-2 border border-gray-400 rounded shadow"
                                        onClick={handleLikedButtonClick}
                                        value={likedStatus === 'liked' ? 'Unlike' : 'Like'}
                                    />
                                </div>
                                <form onSubmit={handleSubmit}>
                                    <div className="flex items-center mt-1">
                                        <label>
                                            Watched it? Why not give it a short review:
                                            <textarea
                                                className="mt-2 block p-2.5 w-full text-sm text-gray-900 bg-gray-50 rounded-lg border border-gray-300 focus:ring-blue-500 focus:border-blue-500"
                                                rows="4"
                                                maxLength={200}
                                                value={inputValue}
                                                onChange={handleTextChange}
                                            ></textarea>
                                        </label>
                                    </div>
                                    <label className="flex flex-row items-center">
                                        And a rating:
                                        <div className="w-1"></div>
                                        <ReactStars
                                            count={5}
                                            onChange={ratingChanged}
                                            size={24}
                                            activeColor="#ffd700"
                                            value={rating}
                                        />
                                    </label>
                                    <button className="bg-white hover:bg-gray-100 text-gray-800 font-semibold py-1 px-2 border border-gray-400 rounded shadow" type="submit">Submit</button>
                                </form>
                            </div>
                        )}
                        <h1 className="mt-4 text-center font-bold text-2xl text-gray-800 transition-colors duration-300">
                            Reviews
                        </h1>
                        <h1 className="text-center font-normal text-xl text-gray-800 transition-colors duration-300">
                            see what other people had to say about {movie.Title}
                        </h1>
                    </div>
                    <div className="max-w-4xl flex flex-wrap justify-evenly">
                        {reviews.map((review, index) => (
                            <ReviewCard
                                key={index}
                                text={review.ReviewText}
                                stars={review.RatingStars}
                                date={review.DatePosted}
                            />
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
};
