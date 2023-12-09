import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import './App.css';
import Movies from './components/Movies.js';
import MovieView from './components/MovieView.js';
import CreateAccount from './components/CreateAccount';
import LoginPage from './components/LoginPage';
//import UserDetails from './components/UserDetails';

function App() {
    //const isLoggedIn = window.localStorage.getItem("loggedIn");
    return (
        <Router>
            <Routes>
                {/*
                <Route
                    exact
                    path="/"
                    element={isLoggedIn === "true" ? <UserDetails /> : <LoginPage />}
                />
                */}
                <Route path='/' element={<Movies />} />
                <Route path='/login' element={<LoginPage />} />
                <Route path='/register' element={<CreateAccount />} />
                <Route path='/movies' element={<Movies />} />
                <Route path='/movies/:title' element={<MovieView />} />
            </Routes>
        </Router>
    );
}

export default App;
