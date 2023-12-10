import React, { useState, useEffect } from 'react';
import { TiArrowBack } from 'react-icons/ti';
import MovieCardDirector from './MovieCardDirector';
import { useParams, useNavigate } from 'react-router-dom';
import '../App.css';

export default function DirectorView() {
    const [director, setDirector] = useState({});
    const { id } = useParams();

    useEffect(() => {
        fetch(`http://localhost:1313/api/directors/${id}`)
            .then(response => response.json())
            .then(data => setDirector(data))
            .catch(error => console.error('Error fetching director items:', error));
    }, [id]);

    const navigate = useNavigate();

    const goBack = () => {
        navigate(-1);
    };

    // const directorField = director !== null ? director[0] : null;
    const directorField = director[0];

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
                                className="object-cover h-32 w-32 mr-4 rounded"
                                src="https://images.unsplash.com/photo-1595769816263-9b910be24d5f?q=80&w=2079&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
                                alt="Movie"
                            />
                            <div className="flex flex-col">
                                <div className="font-bold text-2xl">
                                    <p className="text-gray-900">
                                        {directorField && directorField.DirectorName}
                                    </p>
                                    <p className="text-gray-900">
                                        Nationality: {directorField && directorField.Nationality}
                                    </p>
                                    <p className="text-gray-900">
                                        Date of birth: {directorField && directorField.DateOfBirth}
                                    </p>
                                </div>
                            </div>
                        </div>
                        <h1 className="mt-4 text-center font-bold text-2xl text-gray-800 transition-colors duration-300">
                            Movies directed:
                        </h1>
                        {/*
                        {directorField && (
                            <div className="border border-black">
                                <p className="text-gray-900 text-2xl">{directorField.Title}</p>
                                <p className="text-gray-900 text-2xl">{directorField.ReleaseDate}</p>
                                <p className="text-gray-900 text-2xl">{directorField.Genre}</p>
                                <p className="text-gray-900 text-2xl">{directorField.AvgRating}</p>
                            </div>
                        )}
                        */}
                        {directorField && (
                            <MovieCardDirector movie={directorField} />
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
};
