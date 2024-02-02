import ContactInfo from "./contactInfo"
import ChatBubble from "./chatBubble"

import "./chatBox.css"

function ChatBox({chat})
{
    //contact info on top
    //messages bottom

    return(
        <div>
            <div className="chatContainer">
            <ContactInfo name={chat.uname} img={null}/>
            </div>
            <div className="chatContainer">
                {
                chat.messages.map((message) => ((<ChatBubble user={message.uname} message={message.text} isUser={message.isUser} time={message.time}/>)) )
                }
            </div>

        </div>
    )
}

export default ChatBox