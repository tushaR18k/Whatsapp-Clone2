import React, {useState} from 'react';
import './ChatInput.css';



const ChatInput = ({messages,onSendMessage, jwt, userId,selectedFriendId}) =>{
    const [message, setMessage] = useState('');
    const [selectedFile, setSelectedFile] = useState(null);

    const handleInputChange = (e) =>{
        setMessage(e.target.value);
    }

    const handleFileChange = (e)=>{
        setSelectedFile(e.target.files[0]);
    }

    const handleSendMessage = () => {
        // let messageBody;
        // let fileMessage = 0;
        // if(selectedFile){
        //     // messageBody = JSON.stringify({"sender_id":userId,"receiver_id":selectedFriendId,
        //     //             "message_type":"document","context":message,"file_name":selectedFile.name,"file_size":selectedFile.size,
        //     //             "file_type":"document"})
        //     messageBody = JSON.stringify({"sender_id":userId,"receiver_id":selectedFriendId,"message_type":"text","context":message,
        //                         "file_name":selectedFile.name.toString(),"file_size":selectedFile.size.toString(),
        //                         "file_type":"document"})
        //     fileMessage = 1
        // }else{
        //     messageBody = JSON.stringify({"sender_id":userId,"receiver_id":selectedFriendId,"message_type":"text","context":message})
        // }
        if(message.trim() !== '' || selectedFile){
            const formData = new FormData();
            formData.append("sender_id",userId)
            formData.append("receiver_id",selectedFriendId)
            formData.append("message_type",'text')
            formData.append("context",message)
            let file_location="";
            if (selectedFile){
                formData.append('file', selectedFile);
            }
            fetch("http://localhost:8080/api/messages",{
                method: 'POST',
                headers:{
                    'Authorization': `Bearer ${jwt}`,
                },
                body: formData
            })
                .then(response => response.json())
                .then(data => {
                    file_location = data.file_location
                    let textMsg;
                    if(selectedFile){
                        textMsg = {id:(messages.length === 0)?1:messages[messages.length - 1]['id']+1,sender_id:userId,
                                receiver_id:selectedFriendId,message_type:"document",context:message,file_location:file_location}
                    }else{
                        textMsg = {id:(messages.length === 0)?1:messages[messages.length - 1]['id']+1,sender_id:userId,receiver_id:selectedFriendId,message_type:"text",context:message}
                    }
                    onSendMessage(textMsg)
                    
                })
                .catch(error=>{
                    console.log(error);
                })

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
            <div className='file-input-wrapper'>
                <input type='file'
                onChange={handleFileChange} 
                className='file-input'
                />
                <button className='file-input-button'>Select File</button>
            </div>
            
            <input
                type="text"
                className='text-input'
                placeholder="Type your message..."
                value={message}
                onChange={handleInputChange}
                onKeyDown={handleKeyDown} />
                <button onClick={handleSendMessage}>Send</button>
        </div>
    );
}

export default ChatInput;