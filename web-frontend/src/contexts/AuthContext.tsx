import { createContext, useContext, useState, ReactNode, useEffect } from 'react';
import { apiClient } from '@/lib/api';

interface AuthContextType {
  apiKey: string | null;
  isAuthenticated: boolean;
  loading: boolean;
  setApiKey: (key: string | null) => void;
  signOut: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [apiKey, setApiKeyState] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Load API key from localStorage on mount
    const storedKey = localStorage.getItem('driftlock_api_key');
    if (storedKey) {
      setApiKeyState(storedKey);
      apiClient.setApiKey(storedKey);
    }
    setLoading(false);
  }, []);

  const setApiKey = (key: string | null) => {
    setApiKeyState(key);
    apiClient.setApiKey(key);
  };

  const signOut = () => {
    setApiKey(null);
    apiClient.setApiKey(null);
  };

  return (
    <AuthContext.Provider
      value={{
        apiKey,
        isAuthenticated: !!apiKey,
        loading,
        setApiKey,
        signOut,
      }}
    >
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
