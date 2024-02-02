//TODO: STYLE
import pfp from "../../logo192.png"
import "./profile.css"

//TODO: get this stuff with useEffect
function Profile({img=null, fname, lname, phone, username, bio}){
    //<img src={URL.createObjectURL(img)}/>
    return (
        <div className="profile-container">
        <img src={pfp} width="50" height="50"/>
        <h1>{username}</h1>
        <h2>{"Name: " + fname + " " + lname}</h2>
        <h2>{"Phone: " + phone}</h2>
        <p>
            {bio}
        </p>
        </div>
    )
}

export default Profile