import "./chatBubble.css"

function ChatBubble({uname, message, isUser, time}){
    const containerClass = isUser == "true" ? "messageContainer user" : "messageContainer other"
    return(
        <div className={containerClass}>
            {uname}
            <br/>
            {message}
            <br />
            {time}
        </div>
    )
}

export default ChatBubble