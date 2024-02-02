import pfp from "../../logo192.png"

import "./chatPreview.css"

function ChatPreview({ img = null, name, time, lastMsg, onClick, id, isActive }) {
    return (
        <div className={isActive == "true" ? "chatContainer active" : "chatContainer"} onClick={onClick}>
            <div className="pfp">
                <img src={pfp} width="25" height="25" />
            </div>
            <div className="name">{name}</div>
            <div className="time">{time}</div>
            <div className="lastMsg">{lastMsg}</div>
        </div>
    )
}

export default ChatPreview