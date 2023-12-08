import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import './App.css';
import Movies from './components/Movies.js';
import Movie from './components/Movie.js';

function App() {
    return (
        <Router>
            <Routes>
                <Route path='/' element={<Movies />} />
                <Route path='/movies' element={<Movies />} />
                <Route path='/movies/:title' element={<Movie />} />
            </Routes>
        </Router>
    );
}

export default App;
