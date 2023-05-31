import React from "react";
import { useNavigate } from "react-router-dom";

import './Navbar.css';

const Navbar = ({handleLogin}) => {
    const navigate = useNavigate();

    const handleLogout = () => {
        handleLogin('');
        navigate('/login');
    }

    return (
        <nav className="navbar">
            <div className="navbar-logo">My Logo</div>
            <button className="logout-button" onClick={handleLogout}>Logout</button>
        </nav>
    );
};

export default Navbar;