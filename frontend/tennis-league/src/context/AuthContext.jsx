import React, { createContext, useContext, useEffect, useState } from 'react';


const AuthContext = createContext();


export function AuthProvider({ children }) {
const [user, setUser] = useState(null);


useEffect(() => {
const stored = localStorage.getItem('user');
if (stored) setUser(JSON.parse(stored));
}, []);


const login = (userData) => {
localStorage.setItem('user', JSON.stringify(userData));
setUser(userData);
};


const logout = () => {
localStorage.removeItem('user');
setUser(null);
};


return (
<AuthContext.Provider value={{ user, login, logout, isAuthenticated: !!user }}>
{children}
</AuthContext.Provider>
);
}


export function useAuth() {
return useContext(AuthContext);
}