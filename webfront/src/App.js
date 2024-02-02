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

const ProtectedRoute = ({ children }) => {
    const { token } = useAuth();

    if (!token) {
       return <Navigate to="/" replace />;
    }

    return children;
};

const Navigation = () => {
    const { token } = useAuth()
    const { onLogout } = useAuth()
    return (
        <nav>
            {!token && (
                <NavLink to="/">Login</NavLink>
            )}
            {token && (
                <button type="button" onClick={onLogout}>
                    Sign Out
                </button>
            )}
            <NavLink to="/dashboard">Dashboard</NavLink>

            
        </nav>
    );
};

function App() {
    const [token, setToken] = useState(null);

    return (
        <>
            <AuthProvider>
                <BrowserRouter>
                    <Navigation />
                    <Routes>
                        <Route path="/" element={<Login />} />
                        {/* <Route path="/main" element={
                            <ProtectedRoute><Dashboard/></ProtectedRoute>
                        }/> */}
                        <Route path="/signup" element={<SignUp />} />
                    </Routes>
                </BrowserRouter>
            </AuthProvider>
        </>
    );
}

export default App;
