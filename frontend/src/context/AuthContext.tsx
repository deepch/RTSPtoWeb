import React, { createContext, useContext, useState, useEffect } from 'react';
import client from '@/api/client';
import { useNavigate } from 'react-router-dom';

type Role = 'admin' | 'user' | null;

interface AuthContextType {
  user: { username: string; role: Role } | null;
  login: (username: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<{ username: string; role: Role } | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const navigate = useNavigate();

  const checkAuth = async () => {
    try {
      // Try to fetch streams to check if logged in
      await client.get('/streams');

      // If successful, check role by trying to access admin page
      try {
        await client.get('/pages/stream/add');
        setUser({ username: 'Admin', role: 'admin' });
      } catch (error) {
        setUser({ username: 'User', role: 'user' });
      }
    } catch (error) {
      setUser(null);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    checkAuth();
  }, []);

  const login = async (username: string, password: string) => {
    await client.post('/login', { username, password });
    await checkAuth();
    navigate('/');
  };

  const logout = async () => {
    await client.get('/logout');
    setUser(null);
    navigate('/login');
  };

  return (
    <AuthContext.Provider value={{ user, login, logout, isLoading }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
