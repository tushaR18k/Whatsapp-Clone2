import React, {useState} from "react";
import {useNavigate} from 'react-router-dom';
import axios from 'axios';

import './Signup.css';


const Signup = () =>{
    const navigate = useNavigate();
    const [username, setUsername] = useState('');
    const [email,setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');



    const handleSubmit = async(e) =>{

        e.preventDefault();

        try{
            const response = await axios.post('http://localhost:8080/api/signup',{username,email,password});
            const token = response.data.token;
            localStorage.setItem('token',token)

            navigate("/login");
        }catch(error){
            setError("Signup Failed. Please try again");
        }
    };

    return (
        <div className="signup-container">
            <h2>Sign Up</h2>
            {error && <div className="error">{error}</div>}
            <form className="signup-form" onSubmit={handleSubmit}>
                <div className="form-group">
                    <label htmlFor="username">Username:</label>
                    <input
                        type="text"
                        id="username"
                        value={username}
                        onChange={(e)=>setUsername(e.target.value)}
                    />
                </div>
                <div className="form-group">
                    <label htmlFor="email">Email:</label>
                    <input
                        type="email"
                        id="email"
                        value={email}
                        onChange={(e)=>setEmail(e.target.value)}
                    />
                </div>
                <div className="form-group">
                    <label htmlFor="password">Password:</label>
                    <input
                        type="password"
                        id="password"
                        value={password}
                        onChange={(e)=>setPassword(e.target.value)}
                    />
                </div>
                <div className="form-group">
                    <button type="submit">Sign Up</button>
                </div>
            </form>
        </div>
    );
}

export default Signup;