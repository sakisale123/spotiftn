import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import RegisterPage from './components/Auth/RegisterPage';
import LoginPage from './components/Auth/LoginPage';
import ArtistPage from './components/Pages/ArtistPage';
import ProtectedRoute from './components/Auth/ProtectedRoute';
import './App.css';

function App() {
  return (
    <Router>
      <div className="App">
        <Routes>
          {/* Ruta za registraciju */}
          <Route path="/register" element={<RegisterPage />} />

          {/* Ruta za login */}
          <Route path="/login" element={<LoginPage />} />

          {/* Ruta za pocetnu - preusmerava na registraciju za sada */}
          <Route path="/" element={<RegisterPage />} />

          {/* Protected Routes */}
          <Route element={<ProtectedRoute />}>
            <Route path="/artists" element={<ArtistPage />} />
          </Route>
        </Routes>
      </div>
    </Router>
  );
}

export default App;