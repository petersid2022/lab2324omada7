import React from "react";
import '../App.css';

export default class LoginPage extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            username: "",
            password: "",
            wrongpass: false,
            nouser: false,
        };
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleSubmit(e) {
        e.preventDefault();
        // console.log(JSON.stringify({
        //     username: this.state.username,
        //     password: this.state.password,
        // }));

        fetch("http://localhost:1313/login", {
            method: "POST",
            crossDomain: true,
            headers: {
                "Content-Type": "application/json",
                Accept: "application/json",
                "Access-Control-Allow-Origin": "*",
            },
            body: JSON.stringify({
                username: this.state.username,
                password: this.state.password,
            }),
        })
            .then((res) => res.json())
            .then((data) => {
                if (data.status === "ok") {
                    const { user, token } = data.data;
                    window.localStorage.setItem("token", token);
                    window.localStorage.setItem("loggedIn", true);
                    window.localStorage.setItem("username", user.Username);
                    window.localStorage.setItem("email", user.Email);
                    window.location.href = "/";
                }

                if (data.status === "passNoMatch") {
                    this.setState({ wrongpass: true });
                    this.setState({ nouser: false });
                }

                if (data.status === "usernotfound") {
                    this.setState({ nouser: true });
                    this.setState({ wrongpass: false });
                }
            }).catch((err) => {
                console.log(err);
            });
    }

    render() {
        return (

            <section>
                <div className="flex flex-col items-center justify-center px-6 py-8 mx-auto md:h-screen lg:py-0">
                    <div className="w-full bg-white rounded-lg shadow dark:border md:mt-0 sm:max-w-md xl:p-0 dark:bg-gray-800 dark:border-gray-700">
                        <div className="p-6 space-y-4 md:space-y-6 sm:p-8">
                            <h1 className="text-xl font-bold leading-tight tracking-tight text-gray-900 md:text-4xl dark:text-white">Sign in to your account</h1>
                            <form className="space-y-4 md:space-y-6" onSubmit={this.handleSubmit}>
                                <div>
                                    <label htmlFor="username" className="block mb-2 text-xl font-medium text-gray-900 dark:text-white">Username</label>
                                    <input type="username" name="username" id="username" className="bg-gray-50 border border-gray-300 text-gray-900 sm:text-xl rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" placeholder="johnsmith13" value={this.state.username} onChange={(e) => this.setState({ username: e.target.value })} required="" />
                                </div>
                                <div>
                                    <label htmlFor="password" className="block mb-2 text-xl font-medium text-gray-900 dark:text-white">Password</label>
                                    <input type="password" name="password" id="password" placeholder="••••••••" className="bg-gray-50 border border-gray-300 text-gray-900 sm:text-xl rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" value={this.state.password} onChange={(e) => this.setState({ password: e.target.value })} required="" />
                                </div>
                                {(this.state.nouser ? <div id="toast-warning" className="flex items-center w-full max-w-xs p-4 text-gray-500 bg-white rounded-lg shadow-lg dark:text-gray-400 dark:bg-gray-800" role="alert">
                                    <div className="inline-flex items-center justify-center flex-shrink-0 w-8 h-8 text-orange-500 bg-orange-100 rounded-lg dark:bg-orange-700 dark:text-orange-200">
                                        <svg aria-hidden="true" className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"> </path></svg>
                                        <span className="sr-only">Warning icon</span>
                                    </div>
                                    <div className="ml-3 text-sm font-normal">No user found.</div>
                                </div> : null)}
                                {(this.state.wrongpass ? <div id="toast-warning" className="flex items-center w-full max-w-xs p-4 text-gray-500 bg-white rounded-lg shadow-lg dark:text-gray-400 dark:bg-gray-800" role="alert">
                                    <div className="inline-flex items-center justify-center flex-shrink-0 w-8 h-8 text-orange-500 bg-orange-100 rounded-lg dark:bg-orange-700 dark:text-orange-200">
                                        <svg aria-hidden="true" className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"> </path></svg>
                                        <span className="sr-only">Warning icon</span>
                                    </div>
                                    <div className="w-full ml-3 text-sm font-bold">Wrong password.</div>
                                </div> : null)}
                                <button
                                    type="submit"
                                    className="w-full text-black bg-white border-2 border-primary-900 hover:bg-indigo-100 focus:ring-4 focus:outline-none focus:ring-primary-300 font-medium rounded-lg text-xl px-5 py-2.5 text-center dark:bg-primary-600 dark:hover:bg-primary-700 dark:focus:ring-primary-800"
                                >
                                    Sign in
                                </button>

                                <p className="text-xl font-light text-gray-500 dark:text-gray-400">
                                    Don’t have an account yet? <a href="/register" className="font-medium text-primary-600 hover:underline dark:text-primary-500">Sign up</a>
                                </p>
                            </form>
                        </div>
                    </div>
                </div>
            </section>
        );
    };
}
