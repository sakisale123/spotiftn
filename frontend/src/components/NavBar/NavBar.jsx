import React, { useState, useEffect, useRef } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import axios from 'axios';
import { getUserId } from '../../utils/auth';
import './NavBar.css';

const NavBar = () => {
    const navigate = useNavigate();
    const [isMenuOpen, setIsMenuOpen] = useState(false);
    const [notifications, setNotifications] = useState([]);
    const [showNotifications, setShowNotifications] = useState(false);
    const notifRef = useRef(null);

    const userId = getUserId();

    const handleLogout = () => {
        localStorage.removeItem('token');
        navigate('/login');
    };

    const toggleMenu = () => {
        setIsMenuOpen(!isMenuOpen);
    };

    useEffect(() => {
        if (!userId) return;

        const fetchNotifications = async () => {
            try {
                const token = localStorage.getItem('token');
                const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';

                const response = await axios.get(`${apiUrl}/api/notifications/${userId}`, {
                    headers: { 'Authorization': `Bearer ${token}` }
                });

                if (Array.isArray(response.data)) {
                    setNotifications(response.data);
                }
            } catch (err) {
                console.error("Failed to fetch notifications:", err);
            }
        };

        fetchNotifications();

        const interval = setInterval(fetchNotifications, 30000);
        return () => clearInterval(interval);
    }, [userId]);

    useEffect(() => {
        const handleClickOutside = (event) => {
            if (notifRef.current && !notifRef.current.contains(event.target)) {
                setShowNotifications(false);
            }
        };

        document.addEventListener('mousedown', handleClickOutside);
        return () => {
            document.removeEventListener('mousedown', handleClickOutside);
        };
    }, []);

    const toggleNotifications = () => {
        setShowNotifications(!showNotifications);
    };

    const formatDate = (dateString) => {
        const date = new Date(dateString);
        return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    };

    return (
        <nav className="navbar">
            <div className="navbar-logo">
                <Link to="/artists">Spotiftn</Link>
                <div className="hamburger" onClick={toggleMenu}>
                    <span></span>
                    <span></span>
                    <span></span>
                </div>
            </div>
            <ul className={`navbar-links ${isMenuOpen ? 'active' : ''}`}>
                <li onClick={() => setIsMenuOpen(false)}><Link to="/artists">Artists</Link></li>
            </ul>

            <div className="navbar-right">
                {userId && (
                    <div className="notification-container" ref={notifRef}>
                        <div className="notification-bell" onClick={toggleNotifications}>
                            <svg viewBox="0 0 24 24" width="24" height="24" stroke="currentColor" strokeWidth="2" fill="none" strokeLinecap="round" strokeLinejoin="round">
                                <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>
                                <path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
                            </svg>
                        </div>

                        {showNotifications && (
                            <div className="notification-dropdown">
                                <div className="notification-header">
                                    <h3>Notifications</h3>
                                </div>
                                <div className="notification-list">
                                    {notifications.length === 0 ? (
                                        <div className="notification-empty">No notifications</div>
                                    ) : (
                                        notifications.map((notif) => (
                                            <div key={notif.id} className={`notification-item ${!notif.is_read ? 'unread' : ''}`}>
                                                <div className="notification-icon">
                                                    {notif.type === 'concert' ? '🎫' :
                                                        notif.type === 'release' ? '🎵' :
                                                            notif.type === 'trending' ? '🔥' : '📢'}
                                                </div>
                                                <div className="notification-content">
                                                    <p className="notification-message">{notif.message}</p>
                                                    <span className="notification-time">{formatDate(notif.created_at)}</span>
                                                </div>
                                            </div>
                                        ))
                                    )}
                                </div>
                            </div>
                        )}
                    </div>
                )}

                <div className="navbar-auth desktop-only">
                    <Link to="/change-password" title="Change Password" style={{ textDecoration: 'none' }}>
                        <button className="change-password-btn">Change Password</button>
                    </Link>
                    <button onClick={handleLogout} className="logout-btn">Logout</button>
                </div>
            </div>
        </nav>
    );
};

export default NavBar;
