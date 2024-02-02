import { useState } from "react"
import ChatPreview from "../../components/chatPreview/chatPreview"
import ChatBox from "../../components/chatBox/chatBox"

function Chats() {
    [selectedChat, setSelectedChat] = useState(0)
    [chatList, setChatList] = useState([])
    [searchQuery, setSearchQuery] = useState("")

    const onChatClick = (id) => {
        setSelectedChat(id)
    }

    return (
        <>
            <input
                value={searchQuery}
                type="text"
                placeholder="Type query"
                onChange={ev => setSearchQuery(ev.target.value)}
                className={"searchBox"} />
            <div className="chatList">{
                chatList.map((chat, i) => {
                    (
                        <ChatPreview
                            name={chat.name} time={chat.time}
                            lastMsg={chat.lastMsg}
                            onClick={onChatClick}
                            id={i}
                            isActive={selectedChat == i ? "true" : "false"}
                        />
                    )
                }
                )
            }
            </div>
            <div className="chatBox">
                <ChatBox chat={chat}/>
            </div>

        </>
    )
}

export default Chats