import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import './App.css';
import Movies from './components/Movies.js';
import MovieView from './components/MovieView.js';
import CreateAccount from './components/CreateAccount';
import LoginPage from './components/LoginPage';
import Directors from './components/Directors';
import DirectorView from './components/DirectorView';
import Actors from './components/Actors';
import ActorView from './components/ActorView';

function App() {
    return (
        <Router>
            <Routes>
                <Route path='/' element={<Movies />} />
                <Route path='/login' element={<LoginPage />} />
                <Route path='/register' element={<CreateAccount />} />
                <Route path='/movies' element={<Movies />} />
                <Route path='/actors' element={<Actors />} />
                <Route path='/directors' element={<Directors />} />
                <Route path='/movies/:title' element={<MovieView />} />
                <Route path='/directors/:id' element={<DirectorView />} />
                <Route path='/actors/:id' element={<ActorView />} />
            </Routes>
        </Router>
    );
}

export default App;
