import React, { Fragment, useState, useEffect } from 'react';
import ChatHeader from "../ChatHeader/ChatHeader";
import ChatMessages from "../ChatMessages/ChatMessages";
import ChatInput from "../ChatInput/ChatInput";
import {MessageEvent, SendMessageEvent, RouteEvent} from '../Event';

const ChatHome = ({contactName,contactId,contactStatus,userId,jwt,socket}) =>{
    const [messages,setMessages] = useState([]);

    const sendEvent = (eventName,payload) => {
        const event = new MessageEvent(eventName,payload)

        socket.send(JSON.stringify(event))
    }

    const onSendMessage = (msg) => {
        // setMessages((prevMessages) => {
        //     return [...prevMessages, msg];
        // });

        //socket.send(msg.context)
        let outgoingEvent = new SendMessageEvent(msg.context,userId,contactId,
                    userId.toString()+contactId.toString(),
                    msg.message_type,msg.file_location)
        console.log(outgoingEvent)
        sendEvent("send_message", outgoingEvent);

    }

    useEffect(()=>{
        if(socket){
            //listener for the onmessage event
            socket.onmessage = function(evt) {
            console.log(evt)
            //parsing message as JSON
            const eventData = JSON.parse(evt.data);

            const event = Object.assign(new MessageEvent,eventData);
            //managing the message
            RouteEvent(event,setMessages);
    }
        }
    },[socket])



    return(

        <Fragment>
            <ChatHeader contactName={contactName} contactId={contactId} contactStatus={contactStatus}/>
            <ChatMessages userId={userId} selectedFriendId={contactId} jwt={jwt} messages={messages} setMessages={setMessages}/>
            <ChatInput messages={messages} onSendMessage={onSendMessage} jwt={jwt} userId={userId} selectedFriendId={contactId}/>

        </Fragment>
    );
}

export default ChatHome;