import React, {Fragment, useState, useEffect} from 'react';
import { useNavigate, useLocation } from "react-router-dom";
import {MessageEvent,RouteEvent} from '../Event';

import './ChatSideBar.css';
import jwtDecode from 'jwt-decode';

const ChatSideBar = ({onAddFriend,jwt, userId, handleFriendClick}) => {
    const navigate = useNavigate();
    const {pathname} = useLocation();
    const [showAddFriendForm, setShowAddFriendForm] = useState(false);
    const [friendEmail, setFriendEmail] = useState('');
    const [errorMessage, setErrorMessage] = useState('');
    const [friends, setFriends] = useState(null) 
    const [friendError, setFriendError] = useState(false);

    const handleAddFriendClick = () =>{
        setShowAddFriendForm(true);
    };

    const handleFriendEmailChange = (e) =>{
        setFriendEmail(e.target.value);
        setErrorMessage('');
    }

    const handleFriendListItemClick = (id, email) =>{

        // //Establish websocket connection
        // const socket = new WebSocket("ws://localhost:8080/ws?jwt="+jwt)

        // //Handling WebSocket events
        // socket.onopen = () =>{
        //     console.log('WebSocket connection opened');
        // };

        // socket.onclose = () =>{
        //     console.log('WebSocket connection opened');
        // }

        //handleFriendClick(id, email,socket)
        handleFriendClick(id,email)
    }

    //fetching friends from backend
    const fetchFriends = () => {
        fetch("http://localhost:8080/api/friends/"+userId,{
                    method: 'GET',
                    headers:{
                        'Authorization': `Bearer ${jwt}`,
                        'Content-Type': 'application/json',
                    },
                })
                .then(response => response.json())
                .then(data => {
                    const friendsData = data;
                    console.log("FriendsData: "+friendsData);
                    setFriends(friendsData)
                })
            .catch(error => {
                setFriendError(true);
                console.log(error);
            })
    }

    useEffect(() => {
            if(userId){
                fetchFriends();
            }
                
            console.log(friends);
    },[jwt,navigate,pathname]);

    const handleAddFriendSubmit = (e) =>{
        e.preventDefault();
        console.log(friendEmail);
        console.log(userId);
        fetch("http://localhost:8080/api/add-friend",{
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${jwt}`,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({"userId":userId,"friendEmail": friendEmail}),
        })
            .then(response => response.json())
            .then(data => {
                console.log(data);
                if(data["error"] === "Unauthorized" || data["error"] === "Invalid Token"){
                    navigate("/login")
                }else if(data["error"]){
                    setErrorMessage(data["error"]);
                }else{
                    setShowAddFriendForm(false);
                    fetchFriends();
                }
            })
            .catch(error =>{
                console.error(error);
            });

        setFriendEmail('');
    }

    const handleCancelAddFriend = () =>{
        setFriendEmail('');
        setErrorMessage('');
        setShowAddFriendForm(false);
    }

    return (
        <div className='chat-sidebar-inner'>
            <button className='add-friend-button' onClick={handleAddFriendClick}>
                Add Friend
            </button>
            {showAddFriendForm && (
                <Fragment>
                    <div className="backdrop-overlay open" onClick={handleCancelAddFriend} />
                    <div className='add-friend-popup open'>
                        <form className='add-friend-form' onSubmit={handleAddFriendSubmit}>
                            <h3>Add Friend</h3>
                            <input
                                type='email'
                                placeholder="Enter friend's email"
                                value={friendEmail}
                                onChange={handleFriendEmailChange}
                            />
                            <button type='submit'>Add</button>
                            <button type="button" onClick={handleCancelAddFriend}>
                                Cancel
                            </button>
                            {errorMessage && <div className='error-message'>{errorMessage}</div>}
                        </form>
                    </div>  
                </Fragment> 
            )}
            <div className='friend-list'>
                <h2 className='friend-list-title'>Friends</h2>
                {friendError && (
                    <div className='friend-list-error'>Error fetching friend records!</div>
                )}
                {friends && (
                    <ul className='friend-list-items'>
                    {friends.map((friend) => (
                        <li className="friend-list-item" key={friend.id} onClick={() => handleFriendListItemClick(friend.id, friend.email)}>{friend.email}</li>
                    ))}
                </ul>
                )}
                
            </div>
        </div>
    );
}


export default ChatSideBar;