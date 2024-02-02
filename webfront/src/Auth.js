import { useEffect, useState, createContext, useContext } from 'react';
import { useNavigate } from 'react-router-dom';

const authContext = createContext(null)

export const useLocalStorage = (keyName, defaultValue) => {
    const [storedValue, setStoredValue] = useState(() => {
      try {
        const value = window.localStorage.getItem(keyName);
        if (value) {
          return JSON.parse(value);
        } else {
          window.localStorage.setItem(keyName, JSON.stringify(defaultValue));
          return defaultValue;
        }
      } catch (err) {
        return defaultValue;
      }
    });
    const setValue = (newValue) => {
      try {
        window.localStorage.setItem(keyName, JSON.stringify(newValue));
      } catch (err) {}
      setStoredValue(newValue);
    };
    return [storedValue, setValue];
  };



const Auth = async (email, password) => 
{
    const data = {email: email, password: password}
    console.log(data)
    //TODO: change api
    // return fetch("http://localhost:8000/api/auth/login", {
    //     method: 'POST',
    //     headers: {
    //             "Content-Type": "application/json",
    //     },
    //     body: JSON.stringify(data),
    //     credentials: 'include'
    // }).
    // then((resp) => resp.json()).
    // then( (data) =>
    //     ({
    //         data: data.access_token
    //     })
    //     )
    return 50
}

const AuthProvider = ({ children }) => {
    const nav = useNavigate()
    const [token, setToken] = useLocalStorage("token", null);
    const handleLogin = async (email, password) => {
        const newToken = await Auth(email, password);
        console.log(newToken)
        setToken(newToken);
    };

    const handleLogout = () => {
        //TODO: change api
        //fetch("http://localhost:8000/api/auth/logout", {method:"GET", credentials: 'include'})
        setToken(null);
        nav("/")
    };

    const value = {
        token,
        onLogin: handleLogin,
        onLogout: handleLogout,
    };

    return (
        <authContext.Provider value={value}>
            {children}
        </authContext.Provider>
    );
};

const useAuth = () => {
    return useContext(authContext);
};

export { Auth, AuthProvider, useAuth}