import "./inputForm.css"

function InputForm({label, value, placeholder, setValue, errorVal, type="text"})
{
    return(
        <div className={"inputContainer"}>
                {label}:
                <input
                    value={value}
                    type={type}
                    placeholder={placeholder}
                    onChange={ev => setValue(ev.target.value)}
                    className={"inputBox"} />
                {errorVal && (<br/>)}
                <label className="errorLabel">{errorVal}</label>
            </div>
    )
}

export default InputForm