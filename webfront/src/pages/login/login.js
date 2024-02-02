import React, { useState } from "react";
import { useAuth } from "../../Auth";
import { useNavigate} from 'react-router-dom';

import InputForm from "../../components/inputForm/inputForm";

function Login() {
    const navigate = useNavigate()
    const {onLogin} = useAuth()
    const [username, setUsername] = useState("")
    const [usernameError, setUsernameError] = useState("")
    const [password, setPassword] = useState("")
    const [passwordError, setPasswordError] = useState("")
    const onButtonClick = async () => {
        setPasswordError("")
        if ("" === password) {
            setPasswordError("Please enter a password")
            return
        }
        await onLogin(username, password)
        navigate('/profile')
    }


    return <div className={"mainContainer"}>
        <div className={"titleContainer"}>
            <h1>Login</h1>
        </div>

        <InputForm
                label={"Username"}
                value={username}
                placeholder={"Enter your username here"}
                setValue={setUsername}
                errorVal={usernameError} />

            <InputForm
                label={"Password"}
                value={password}
                placeholder={"Enter your password here"}
                setValue={setPassword}
                errorVal={passwordError}
                type="password" />

        <div className={"inputContainer"}>
            <input
                className={"inputButton"}
                type="button"
                onClick={onButtonClick}
                value={"Log in"} />
        </div>

        <a href="signup">Don't have an account?</a>
    </div>
}


export default Login