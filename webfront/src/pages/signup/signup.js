import React, { useState } from "react";
import { useNavigate, redirect } from 'react-router-dom';
import InputForm from "../../components/inputForm/inputForm.js";

function SignUp() {
    //const navigate = useNavigate()

    const [fname, setFname] = useState("")
    const [fnameError, setFnameError] = useState("") //non empty

    const [lname, setLname] = useState("")
    const [lnameError, setLnameError] = useState("") //non empty

    const [phone, setPhone] = useState("")
    const [phoneError, setPhoneError] = useState("") //non empty, uniqeness check

    const [username, setUsername] = useState("")
    const [usernameError, setUsernameError] = useState("") //non empty, uniqeness check

    const [password, setPassword] = useState("")
    const [passwordError, setPasswordError] = useState("") //non empty

    const [image, setImage] = useState()
    const [bio, setBio] = useState("")


    const onButtonClick = async () => {

        //const data = { name: username, email: email, password: password, passwordConfirm: confirmPassword }
        //TODO: new api
        // response = await fetch("localhost:8000/api/auth/register", {
        //     mode: 'no-cors',
        //     method: 'POST',
        //     headers: {
        //         "Content-Type": "application/json",
        //     },
        //     body: JSON.stringify(data),
        //     credentials: 'include'
        // })
        //     .then((resp) => (resp.json()))
        //     .then((stuff) => {
        //         stuff
        //     })
        //TODO: if reg success
        //navigate('/main')
        //TODO: if reg unsuccess
        //SHOW ERRORS

    }

    return (
        <>
            <h1>Sign Up</h1>
            <InputForm
                label={"First Name"}
                value={fname}
                placeholder={"Enter your first name here"}
                setValue={setFname}
                errorVal={fnameError} />

            <InputForm
                label={"Last Name"}
                value={lname}
                placeholder={"Enter your last name here"}
                setValue={setLname}
                errorVal={lnameError} />

            <InputForm
                label={"Phone Name"}
                value={phone}
                placeholder={"Enter your phone here"}
                setValue={setPhone}
                errorVal={phoneError} />

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
                Image:
                <input
                    type="file"
                    name="image"
                    accept="image/*"
                    onChange={(event) =>
                        setImage(event.target.files[0])
                    }
                />
            </div>

            <div className={"inputContainer"}>
                Bio:
                <textarea
                    value={bio}
                    placeholder="Enter your bio here"
                    onChange={ev => setBio(ev.target.value)}
                    className={"inputBox"} />
            </div>

            <div className={"inputContainer"}>
                <input
                    className={"inputButton"}
                    type="button"
                    onClick={onButtonClick}
                    value={"Sign Up"} />
            </div>
        </>

    );
}

export default SignUp