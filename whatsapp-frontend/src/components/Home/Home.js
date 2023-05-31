import React, {Fragment, useEffect, useState} from "react";
import { useNavigate, useLocation } from "react-router-dom";
import jwtDecode from 'jwt-decode';
import Navbar from "../Navbar/Navbar";
import './Home.css';
import ChatSideBar from "../ChatSideBar/ChatSideBar";
import ChatHome from "../ChatHome/ChatHome";


const validateToken = (token) => {
    if (!token) {
      return false; // Token not provided
    }
  
    try {
      const decodedToken = jwtDecode(token);
      const currentTime = Date.now() / 1000; // Convert current time to seconds
  
      if (decodedToken.exp < currentTime) {
        return false; // Token has expired
      }
  
      // Token is valid and not expired
      return true;
    } catch (error) {
      return false; // Token decoding failed or invalid token format
    }
  };

const Home = ({ jwt, handleLogin, userId }) =>{
    const navigate = useNavigate();
    const {pathname} = useLocation();
    const [displayChatFields,setDisplayChatFields] = useState(false);
    const [selectedFriendEmail, setSelectedFriendEmail] = useState('');
    const [selectFriendId, setSelectedFriendId] = useState('');
    const [socket, setSocket] = useState(null);

    //on click handler for friend list
    const handleFriendClick = (friendId, friendEmail) => {
      setSelectedFriendId(friendId);
      setSelectedFriendEmail(friendEmail);
      setDisplayChatFields(true);
      //setSocket(socket);
      
    }

    
    

    useEffect( () => {
        const isValidToken = validateToken(jwt)
        if(!isValidToken){
            navigate("/login",{replace:true,state:{from: pathname}});
        }
        if(socket){
          socket.close();
        }
        if (selectFriendId){
          const newSocket = new WebSocket("ws://localhost:8080/ws?jwt="+jwt)
          setSocket(newSocket);
        }
    },[jwt,navigate,pathname,selectFriendId]);

    const handleAddFriend = (friendEmail) =>{
      console.log("Adding friend: ", friendEmail);
    }

    return (
        <Fragment>
            <Navbar handleLogin={handleLogin}/>
            <div className="home-container">
            <div className="chat-container">
                <div className="chat-sidebar">
                 <ChatSideBar onAddFriend={handleAddFriend} jwt={jwt} userId={userId} handleFriendClick={handleFriendClick}/> 
                {/* <h2>Your contacts</h2> */}
                </div>
                <div className="chat-content">
                {displayChatFields ? (
                  <Fragment>
                  {/* <h2>Chat Area Home</h2> */}
                  <ChatHome contactName={selectedFriendEmail} contactId={selectFriendId} contactStatus="online"
                            userId={userId} jwt={jwt} socket={socket}/>
                </Fragment>
                ):(
                  <div>
                    <p>Click on your friend list to start chatting!</p>
                  </div>
                )}
                
                </div>
            </div>
            </div>
        </Fragment>
      );
}

export default Home