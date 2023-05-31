import React, {useState} from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import './Login.css';

const Login = ( {handleLogin, addUserId}) => {
    const navigate = useNavigate();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');

    const handleSubmit = async (e) =>{
        e.preventDefault();

        try{    
            const response = await axios.post('http://localhost:8080/api/login',{email, password});
            const token = response.data.token;
            const id = response.data.id;

            localStorage.setItem('token',token);
            handleLogin(token);
            addUserId(id);

            navigate('/home');
        }catch(error){
            setError("Login Failed, check your credentials");
        }
    };

    return (
        <div className="login-container">
            <div className="login-form-container">
                <h2>Login</h2>
                {error && <div className="error">{error}</div>}
                <form className="login-form" onSubmit={handleSubmit}>
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
                        <button type="submit">Login</button>
                    </div>
                </form>
            </div>
        </div>
    );
}

export default Login