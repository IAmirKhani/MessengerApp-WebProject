import { useEffect, useState, createContext } from 'react';
import {
    NavLink,
    Routes,
    Route,
    BrowserRouter,
    Navigate,
} from "react-router-dom";

import { AuthProvider, Auth, useAuth } from './Auth.js'

import SignUp from './pages/signup/signup.js';
import Login from './pages/login/login.js';
import Profile from './pages/profile/profile.js';
import Navigation from './components/navBar/navBar.js';
import Chats from './pages/chats/chats.js';

import "./App.css"

const ProtectedRoute = ({ children }) => {
    const { token } = useAuth();

    if (!token) {
        return <Navigate to="/" replace />;
    }

    return children;
};

function App() {
    const [token, setToken] = useState(null);

    return (
        <>
            <BrowserRouter>
                <AuthProvider>
                    <Navigation />
                    <Routes>
                        <Route path="/" element={<Login />} />
                        {<Route path="/main" element={
                            <ProtectedRoute><Chats /></ProtectedRoute>
                        } />}
                        <Route path="/signup" element={<SignUp />} />
                        <Route path="/profile" element={<Profile fname={"Amir"} lname={"Poolad"} phone={"0912229999"} username={"HDxC"} bio={"lmao"} />} />
                    </Routes>
                </AuthProvider>
            </BrowserRouter>
        </>
    );
}

export default App;
