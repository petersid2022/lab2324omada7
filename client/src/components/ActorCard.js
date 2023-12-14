import React from 'react';
import { Link } from 'react-router-dom';
import '../App.css';

const ActorCard = (props) => {
    const actor = props.actor;

    const transformDateFormat = (inputDate) => {
        var parts = inputDate.split("-");
        var year = parts[0];
        var month = parts[1];
        var day = parts[2];
        var dateObject = new Date(year, month - 1, day);
        var transformedDate = dateObject.getDate() + '/' + (dateObject.getMonth() + 1) + '/' + dateObject.getFullYear();
        return transformedDate;
    }


    return (
        <Link
            to={`/actors/${actor.actor_id}`}
            className="m-4 px-3 text-center hover:bg-gray-200 bg-gray-50 text-gray-800 max-w-md rounded-lg overflow-hidden shadow-xl block"
            style={{ textDecoration: 'none' }}
        >
            <div className="pt-3">
                <div className="font-bold text-2xl">
                    <h1 className="text-gray-800 transition-colors duration-300">
                        {actor.ActorName}
                    </h1>
                </div>
            </div>
            <div className="pb-3">
                <p className="text-gray-900 text-2xl">Date of birth: {transformDateFormat(actor.DateOfBirth)}</p>
                <p className="text-gray-900 text-2xl">Nationality: {actor.Nationality}</p>
            </div>
        </Link>
    );
};

export default ActorCard;
