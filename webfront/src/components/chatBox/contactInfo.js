import pfp from "../../logo192.png"
import "./contactInfo.css"

function ContactInfo({img=null, name}){
    return(
        <div className="ContactInfo">
            <img src={pfp} width="25" height="25"/>
            {name}
        </div>
    )
}

export default ContactInfo