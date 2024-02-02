import { useState } from "react"
import ChatPreview from "../../components/chatPreview/chatPreview"
import ChatBox from "../../components/chatBox/chatBox"

import "./chats.css"

function Chats() {
    const fakeChatData = [
        {
            uname: "mamad",
            messages: [
                {
                    uname: "mamad",
                    text: "Hi",
                    isUser: "false",
                    time: "Diruz",

                },
                {
                    uname: "ali",
                    text: "Hi!!",
                    isUser: "true",
                    time: "Diruz",

                },
                {
                    uname: "mamad",
                    text: "aleyk",
                    isUser: "false",
                    time: "Diruz",

                },
            ]
        },
        {
            uname: "dubidubi",
            messages: [
                {
                    uname: "ali",
                    text: "slm",
                    isUser: "true",
                    time: "Pariruz",

                },
                {
                    uname: "dubidubi",
                    text: "askjasdokj",
                    isUser: "false",
                    time: "Diruz",

                },
                {
                    uname: "ali",
                    text: "Va.",
                    isUser: "true",
                    time: "Emruz",

                },
            ]
        },
    ]
    const [selectedChat, setSelectedChat] = useState(0)
    const [chatList, setChatList] = useState(fakeChatData)
    const [searchQuery, setSearchQuery] = useState("")
    let selectedChats = []
    if (searchQuery != ""){
        selectedChats = chatList
    }
    else{
        selectedChats = chatList.filter((chat) => chat.uname.search(searchQuery) != -1)
    }

    const onChatClick = (id) => {
        setSelectedChat(id)
    }

    return (
        <div className="chatPageContainer">
            <input
                value={searchQuery}
                type="text"
                placeholder="Type query"
                onChange={ev => setSearchQuery(ev.target.value)}
                className={"searchBox"} />
            <div className="chatList">{
                selectedChats.map((chat, i) => {
                    return (
                        <ChatPreview
                            name={chat.uname} time={chat.time}
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
                <ChatBox chat={chatList[selectedChat]}/>
            </div>

        </div>
    )
}

export default Chats