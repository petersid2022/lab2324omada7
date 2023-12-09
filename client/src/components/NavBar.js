import React from 'react';
import "../App.css";

export default function NavBar({ searchTerm, handleChange }) {
    const isAuthenticated = localStorage.getItem('loggedIn') === 'true';
    const userName = localStorage.getItem('username');

    const logout = () => {
        window.localStorage.clear();
        window.location.href = '/';
    };

    const login = () => {
        window.location.href = '/login';
    };

    return (
        <header className="z-50 bg-white shadow-lg h-20 hidden md:flex top-0 fixed w-full">
            <div className="flex justify-between items-center w-full">
                <input
                    type="text"
                    className="border-2 text-2xl border-gray-300 p-2 rounded text-black mx-5"
                    placeholder="Search"
                    value={searchTerm}
                    onChange={handleChange}
                />
                {isAuthenticated ? (
                    <div className="flex items-center">
                    <h2 className="text-right bg-gradient-to-r from-purple-800 via-violet-900 to-purple-800 bg-clip-text text-xl font-extrabold text-transparent">Hey {userName},<br />welcome back!</h2>
                    <button onClick={logout} className="bg-red-500 hover:bg-red-700 text-white font-bold px-4 py-2 rounded text-2xl mx-5">Logout</button>
                    </div>
                ) : (
                    <button onClick={login} className="bg-blue-500 hover:bg-blue-700 text-white font-bold px-4 py-2 rounded text-2xl mx-5">Login</button>
                )}
            </div>
        </header>
    );
}
