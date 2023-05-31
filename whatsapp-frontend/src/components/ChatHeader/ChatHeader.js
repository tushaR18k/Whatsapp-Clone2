import React from 'react'

import './ChatHeader.css';



const ChatHeader = ({contactName, contactStatus,contactId}) => {
    return (
        <div className='chat-header'>
            <div className='profile-picture'>
                {/* Render the profile pic here */}
            </div>
            <div className='contact-info'>
                <h3 className='contact-name'>{contactName}</h3>
                <p className='contact-status'>{contactStatus}</p>
            </div>
        </div>
    )
}

export default ChatHeader;