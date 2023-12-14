import React from 'react';
import { Link } from 'react-router-dom';
import '../App.css';

const DirectorCard = (props) => {
    const director = props.director;

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
            to={`/directors/${director.director_id}`}
            className="m-4 px-3 hover:bg-gray-200 text-center bg-gray-50 text-gray-800 max-w-md rounded-lg overflow-hidden shadow-xl block"
            style={{ textDecoration: 'none' }}
        >
            <div className="pt-3">
                <div className="font-bold text-2xl">
                    <h1 className="text-gray-800 transition-colors duration-300">
                        {director.DirectorName}
                    </h1>
                </div>
            </div>
            <div className="pb-3">
                <p className="text-gray-900 text-2xl">Date of birth: {transformDateFormat(director.DateOfBirth)}</p>
                <p className="text-gray-900 text-2xl">Nationality: {director.Nationality}</p>
            </div>
        </Link>
    );
};

export default DirectorCard;
