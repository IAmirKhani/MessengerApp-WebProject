import { NavLink } from "react-router-dom";
import { useAuth } from "../../Auth";
import "./navBar.css"

function Navigation({ pfp = null }) {
    const { token } = useAuth()
    const { onLogout } = useAuth()
    const onDelete = () => { }
    return (
        <nav class="topnav">

            {!token && (
                <NavLink to="/">Login</NavLink>
            )}

            {token && (
                <>
                    <NavLink to="/main">Chats</NavLink>
                    <NavLink to="/profile">Profile</NavLink>

                    <button type="button" onClick={onLogout}>
                        Sign Out
                    </button>

                    <button type="button" onClick={onDelete}>
                        Delete Account
                    </button>

                </>
            )}


        </nav>
    );
};

export default Navigation