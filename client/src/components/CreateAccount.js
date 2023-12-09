import React from 'react';
import { useState } from 'react';
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

export default function CreateAccount() {

    const navigate = useNavigate();
    const [username, setUsername] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [accountCreated, setAccountCreated] = useState(false);
    const [remainingSeconds, setRemainingSeconds] = useState(5);

    const handleSubmit = (e) => {
        e.preventDefault();
        console.log(JSON.stringify({ username, email, password }));
        fetch('http://localhost:1313/create-account', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                Accept: 'application/json',
                "Access-Control-Allow-Origin": "*",
            },
            body: JSON.stringify({ username, email, password }),
        }).then((response) => response.json()).then((data) => {
            if (data.status === "ok") {
                setAccountCreated(true);
            } else {
                alert("User already exists");
            }
        });
    }

    useEffect(() => {
        if (accountCreated) {
            const redirectTimeout = setTimeout(() => {
                navigate('/login');
            }, remainingSeconds * 1000);

            const interval = setInterval(() => {
                setRemainingSeconds((prevSeconds) => prevSeconds - 1);
            }, 1000);

            return () => {
                clearTimeout(redirectTimeout);
                clearInterval(interval);
            };
        }
    }, [accountCreated, remainingSeconds, navigate]);


    return (
        <section>
            <div className="flex flex-col items-center justify-center px-6 py-8 mx-auto md:h-screen lg:py-0">
                <div className="w-full bg-white rounded-lg shadow dark:border md:mt-0 sm:max-w-md xl:p-0 dark:bg-gray-800 dark:border-gray-700">
                    <div className="p-6 space-y-4 md:space-y-6 sm:p-8">
                        <h1 className="text-xl font-bold leading-tight tracking-tight text-gray-900 md:text-4xl dark:text-white">
                            Sign up
                        </h1>

                        <form className="space-y-4 md:space-y-6" onSubmit={handleSubmit}>
                            <div>
                                <label htmlFor="username" className="block mb-2 text-xl font-medium text-gray-900 dark:text-white">Username</label>
                                <input type="text" name="username" id="username" className="bg-gray-50 border border-gray-300 text-gray-900 sm:text-xl rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" placeholder="johnsmith13" onChange={(e) => setUsername(e.target.value)} required="" />
                            </div>
                            <div>
                                <label htmlFor="email" className="block mb-2 text-xl font-medium text-gray-900 dark:text-white">Email</label>
                                <input type="email" name="email" id="email" className="bg-gray-50 border border-gray-300 text-gray-900 sm:text-xl rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" placeholder="name@company.com" onChange={(e) => setEmail(e.target.value)} required="" />
                            </div>
                            <div>
                                <label htmlFor="password" className="block mb-2 text-xl font-medium text-gray-900 dark:text-white">Password</label>
                                <input type="password" name="password" id="password" placeholder="••••••••" className="bg-gray-50 border border-gray-300 text-gray-900 sm:text-xl rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" onChange={(e) => setPassword(e.target.value)} required="" />
                            </div>
                            <button
                                type="submit"
                                className="w-full text-black bg-white border-2 border-primary-900 hover:bg-indigo-100 focus:ring-4 focus:outline-none focus:ring-primary-300 font-medium rounded-lg text-xl px-5 py-2.5 text-center dark:bg-primary-600 dark:hover:bg-primary-700 dark:focus:ring-primary-800"
                            >
                                Create an account
                            </button>

                            {accountCreated && (
                                <p className="animate-bounce text-xl font-medium text-center text-gray-400">
                                    Account created successfully :) Redirecting in {remainingSeconds} seconds...
                                </p>
                            )}
                            <p className="text-xl font-light text-gray-500 dark:text-gray-400">
                                Already have an account? <a href="/login" className="font-medium text-primary-600 hover:underline dark:text-primary-500">Login here</a>
                            </p>
                        </form>
                    </div>
                </div>
            </div>
        </section>
    );
}
