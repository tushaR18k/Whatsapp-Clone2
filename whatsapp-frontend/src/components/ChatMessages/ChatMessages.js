import React, {useEffect, useState} from 'react'
import './ChatMessages.css';



const fetchMessages = async (userId,selectedFriendId,pageNumber,jwt,setIsLoading,setHasMoreMessages,setMessages) => {
            try {
                setIsLoading(true);
                const response = await fetch(`http://localhost:8080/api/messages/${userId}/${selectedFriendId}?page=${pageNumber}`,{
                    method: 'GET',
                    headers:{
                        'Authorization': `Bearer ${jwt}`,
                        'Content-Type': 'application/json',
                    }
                });
                
                const data = await response.json();
                console.log(data);
        
                if (data.length === 0) {
                    setHasMoreMessages(false);
                } else {
                    setMessages((prevMessages) => {
                        const uniqueMessages = data.filter(
                            (newMessage) =>
                                !prevMessages.some((message) =>message.id === newMessage.id)
                        );
                        return [...prevMessages, ...uniqueMessages];
                    });
                }
            } catch (error) {
                console.error('Error fetching message:', error);
            } finally{
                setIsLoading(false);
            }
        };

const ChatMessages = ({userId,selectedFriendId,jwt,messages,setMessages}) => {
    const [pageNumber, setPageNumber] = useState(1);
    const [hasMoreMessages, setHasMoreMessages] = useState(true);
    const [isLoading, setIsLoading] = useState(false);

    useEffect(()=>{
        //Reset the states
        setMessages([]);
        setPageNumber(1);
        setHasMoreMessages(true);
        
    },[selectedFriendId])

    useEffect(()=>{
        //Fetching the messages
        if(selectedFriendId){
            fetchMessages(userId,selectedFriendId,pageNumber,jwt,setIsLoading,setHasMoreMessages,setMessages);
        }
    },[selectedFriendId,pageNumber,jwt])

    const loadMoreMessages = () => {
        setPageNumber((prevPageNumber)=>prevPageNumber+1);
    };

    return (
        <div className='chat-messages'>
            {/* Rendering the messages */}
            {messages.length === 0 ? (
                <div className='no-messages'>
                    <p>No messages between you. Start Chatting!</p>
                </div>
            ):(
                <div>
                    {messages.map((message) => (
                    <div key={message.id} className={`message ${message.sender_id === userId ? 'sender' : 'receiver'}`}>
                        <div className={`message-content ${message.sender_id == userId ? 'align-right': 'align-left'}`}>
                            <p>{message.context}</p>
                            {message.message_type === 'document' && <a className="download-link" href={`http://localhost:8080/api/files/${message.file_path}`}>Download</a>}
                            <span>{message.sender}</span>
                        </div>
                    </div>
                ))}
                </div>
                
            )}
            
            {/* Loading more messages button */}
            {messages.length !== 0 && hasMoreMessages && !isLoading && (
                <button className='load-more-button' onClick={loadMoreMessages}>
                    Load More..
                </button>
            )}

            {isLoading && <p>Loading messages....</p>}
            
        </div>
    )


}

export default ChatMessages;