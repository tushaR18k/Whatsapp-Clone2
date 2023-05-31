import React, {useState} from 'react';
import './ChatInput.css';



const ChatInput = ({messages,onSendMessage, jwt, userId,selectedFriendId}) =>{
    const [message, setMessage] = useState('');

    const handleInputChange = (e) =>{
        setMessage(e.target.value);
    }

    const handleSendMessage = () => {
        if(message.trim() !== ''){
            fetch("http://localhost:8080/api/messages",{
                method: 'POST',
                headers:{
                    'Authorization': `Bearer ${jwt}`,
                    'Content-Type':'application/json',
                },
                body: JSON.stringify({"sender_id":userId,"receiver_id":selectedFriendId,"message_type":"text","context":message})
            })
                .then(response => response.json())
                .then(data => {
                    console.log(data)
                })
                .catch(error=>{
                    console.log(error);
                })
            onSendMessage({id:(messages.length === 0)?1:messages[messages.length - 1]['id']+1,sender_id:userId,receiver_id:selectedFriendId,message_type:"text",context:message})
            setMessage('');
        }
    };

    const handleKeyDown = (e) =>{
        if(e.key === 'Enter'){
            handleSendMessage();
        }
    };

    return (
        <div className='chat-input'>
            <input
                type="text"
                placeholder="Type your message..."
                value={message}
                onChange={handleInputChange}
                onKeyDown={handleKeyDown} />
                <button onClick={handleSendMessage}>Send</button>
        </div>
    );
}

export default ChatInput;