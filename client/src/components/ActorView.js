import React, { useState, useEffect } from 'react';
import { TiArrowBack } from 'react-icons/ti';
import MovieCardDirector from './MovieCardDirector';
import { useParams, useNavigate } from 'react-router-dom';
import '../App.css';

export default function ActorView() {
    const [actor, setActor] = useState({});
    const { id } = useParams();

    const transformDateFormat = (inputDate) => {
        var parts = inputDate.split("-");
        var year = parts[0];
        var month = parts[1];
        var day = parts[2];
        var dateObject = new Date(year, month - 1, day);
        var transformedDate = dateObject.getDate() + '/' + (dateObject.getMonth() + 1) + '/' + dateObject.getFullYear();
        return transformedDate;
    }

    useEffect(() => {
        fetch(`http://localhost:1313/api/actors/${id}`)
            .then(response => response.json())
            .then(data => setActor(data))
            .catch(error => console.error('Error fetching actor items:', error));
    }, [id]);

    const navigate = useNavigate();

    const goBack = () => {
        navigate(-1);
    };

    const actorField = actor ? actor[0] : null;

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
                                className="object-cover h-24 w-24 mr-4 rounded"
                                src="https://images.unsplash.com/photo-1595769816263-9b910be24d5f?q=80&w=2079&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
                                alt="Movie"
                            />
                            <div className="flex flex-col">
                                <div className="font-bold text-2xl">
                                    <p className="text-gray-900">Name: {actorField ? actorField.ActorName : 'N/A'}</p>
                                    <p className="text-gray-900">Nationality: {actorField ? actorField.Nationality : 'N/A'}</p>
                                    <p className="text-gray-900">Date of birth: {actorField ? transformDateFormat(actorField.DateOfBirth) : 'N/A'}</p>
                                </div>
                            </div>
                        </div>
                        <div className="flex flex-col items-center">
                            <h1 className="mt-4 text-center font-bold text-2xl text-gray-800 transition-colors duration-300">
                                Movies he played in:
                            </h1>
                            {/*
                        {actorField && (
                            <div className="border border-black">
                                <p className="text-gray-900 text-2xl">{actorField.Title}</p>
                                <p className="text-gray-900 text-2xl">{actorField.ReleaseDate}</p>
                                <p className="text-gray-900 text-2xl">{actorField.Genre}</p>
                                <p className="text-gray-900 text-2xl">{actorField.AvgRating}</p>
                            </div>
                        )}
                        */}
                            {actorField ? (
                                <MovieCardDirector movie={actorField} />
                            ) :
                                'N/A'}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};
