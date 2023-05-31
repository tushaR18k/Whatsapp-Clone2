import React, { Fragment, useState } from 'react';
import {BrowserRouter as Router, Route, Routes, Navigate, Outlet} from 'react-router-dom';
import Signup from './components/Signup/Signup';
import Login from './components/Login/Login';
import Home from './components/Home/Home';

const App = () => {

  const [jwt, setJwt] = useState('');
  const [userId, setUserId] = useState('');

  const handleLogin = (recievedJwt) =>{
    setJwt(recievedJwt);
  }

  const addUserId = (id) =>{
    setUserId(id);
  }

  const isAuthenticated = !!jwt;
  return (
    <Router>
      <Routes>
        <Route path="/signup" element={<Signup />}/>
        <Route path="/login" element={<Login handleLogin={handleLogin} addUserId={addUserId}/>}/>
        <Route
          path="/home"
          element={isAuthenticated ? <Home jwt={jwt} handleLogin={handleLogin} userId={userId}/> : <Navigate to="/login" replace />}
        />
      </Routes>
    </Router>
  );
}

export default App;
