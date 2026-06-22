import React, { createContext, useContext, useState, useEffect, useCallback } from 'react';
import type { User } from '../api/auth';

interface AuthState {
  token: string | null;
  user: User | null;
  login: (token: string, user: User) => void;
  logout: () => void;
  isAuthenticated: boolean;
  isHR: boolean;
}

const AuthContext = createContext<AuthState>(null!);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [token, setToken] = useState<string | null>(() => localStorage.getItem('token'));
  const [user, setUser] = useState<User | null>(() => {
    const u = localStorage.getItem('user');
    return u ? JSON.parse(u) : null;
  });

  const loginFn = useCallback((t: string, u: User) => {
    localStorage.setItem('token', t);
    localStorage.setItem('user', JSON.stringify(u));
    setToken(t);
    setUser(u);
  }, []);

  const logout = useCallback(() => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setToken(null);
    setUser(null);
  }, []);

  return (
    <AuthContext.Provider value={{
      token, user,
      login: loginFn, logout,
      isAuthenticated: !!token,
      isHR: user?.role === 'hr',
    }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}
