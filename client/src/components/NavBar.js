import React, { useState, useEffect } from 'react';
import "../App.css";

export default function NavBar({ searchTerm, handleChange }) {
    const isAuthenticated = localStorage.getItem('loggedIn') === 'true';
    const username = localStorage.getItem('username');
    const [activeLink, setActiveLink] = useState('');

    useEffect(() => {
        const currentPath = window.location.pathname;
        const link = currentPath === '/' ? 'home' : currentPath.substring(1);
        setActiveLink(link);
    }, []);

    const handleClick = (link) => {
        setActiveLink(link);
    };


    const logout = () => {
        window.localStorage.clear();
        window.location.href = '/';
    };

    const login = () => {
        window.location.href = '/login';
    };

    return (
        <header className="z-50 bg-white shadow-lg h-20 hidden md:flex top-0 fixed w-full">
            <nav className="header-links contents font-semibold text-base lg:text-lg">
                <ul className="flex items-center mr-auto text-2xl">
                    <li className={`p-3 xl:p-6 ${activeLink === 'home' ? 'active' : ''}`}>
                        <a href="/" onClick={() => handleClick('movies')} className="relative">
                            <span>Movies</span>
                            {activeLink === 'home' && <div className="absolute bottom-0 left-0 w-full h-1 bg-purple-800" />}
                        </a>
                    </li>
                    <li className={`p-3 xl:p-6 ${activeLink === 'actors' ? 'active' : ''}`}>
                        <a href="/actors" onClick={() => handleClick('actors')} className="relative">
                            <span>Actors</span>
                            {activeLink === 'actors' && <div className="absolute bottom-0 left-0 w-full h-1 bg-purple-800" />}
                        </a>
                    </li>
                    <li className={`p-3 xl:p-6 ${activeLink === 'directors' ? 'active' : ''}`}>
                        <a href="/directors" onClick={() => handleClick('directors')} className="relative">
                            <span>Directors</span>
                            {activeLink === 'directors' && <div className="absolute bottom-0 left-0 w-full h-1 bg-purple-800" />}
                        </a>
                    </li>
                </ul>
            </nav>
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
                        <h2 className="mr-auto text-2xl"> Welcome back, {username}!</h2>
                        <button onClick={logout} className="bg-red-500 hover:bg-red-700 text-white font-bold px-4 py-2 rounded text-2xl mx-5">Logout</button>
                    </div>
                ) : (
                    <button onClick={login} className="bg-blue-500 hover:bg-blue-700 text-white font-bold px-4 py-2 rounded text-2xl mx-5">Login</button>
                )}
            </div>
        </header>
    );
}
