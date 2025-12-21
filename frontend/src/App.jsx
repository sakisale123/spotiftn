import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import RegisterPage from './components/Auth/RegisterPage';
import LoginPage from './components/Auth/LoginPage';
import ActivationPage from './components/Auth/ActivationPage';
import ForgotPasswordPage from './components/Auth/ForgotPasswordPage';
import ResetPasswordPage from './components/Auth/ResetPasswordPage';
import ChangePasswordPage from './components/Auth/ChangePasswordPage';
import ArtistPage from './components/Pages/ArtistPage';
import AlbumPage from './components/Pages/AlbumPage';
import SongPage from './components/Pages/SongPage';
import ProtectedRoute from './components/Auth/ProtectedRoute';
import './App.css';

function App() {
  return (
    <Router>
      <div className="App">
        <Routes>
          <Route path="/register" element={<RegisterPage />} />

          <Route path="/login" element={<LoginPage />} />

          <Route path="/activate" element={<ActivationPage />} />
          <Route path="/forgot-password" element={<ForgotPasswordPage />} />
          <Route path="/reset-password" element={<ResetPasswordPage />} />

          <Route path="/" element={<RegisterPage />} />

          <Route element={<ProtectedRoute />}>
            <Route path="/artists" element={<ArtistPage />} />
            <Route path="/artists/:artistId/albums" element={<AlbumPage />} />
            <Route path="/albums/:albumId/songs" element={<SongPage />} />
            <Route path="/change-password" element={<ChangePasswordPage />} />
          </Route>
        </Routes>
      </div>
    </Router>
  );
}

export default App;