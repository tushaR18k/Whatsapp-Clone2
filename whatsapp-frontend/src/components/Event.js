
export class MessageEvent {
    constructor(type, payload){
        this.type = type;
        this.payload = payload;
    }
}

export class SendMessageEvent{
    constructor(message, userId,friendId,chatroom){
        this.message = message;
        this.from = userId;
        this.to = friendId;
        this.chatroom = chatroom
    }
}

export class NewMessageEvent{
    constructor(message,from,sent,to,chatroom){
        this.message = message
        this.from = from
        this.to = to
        this.chatroom = chatroom
        this.sent = sent
    }
}

export function RouteEvent(event,setMessages){
    if(event.type === undefined){
        alert("no 'type' field in event");
    }
    switch(event.type){
        case "new_message":
            console.log("new message");
            console.log(event.payload);
            const messageEvent = Object.assign(new NewMessageEvent, event.payload);
            setMessages((prevMessages) => {
                let newMessage = {id:(prevMessages.length === 0)?1:prevMessages[prevMessages.length - 1]['id']+1,sender_id:messageEvent.from,receiver_id:messageEvent.to,message_type:"text",context:messageEvent.message}
                return [...prevMessages, newMessage];
            });
            break;
        default:
            alert("unsupported message type")
            break;
    }
}

