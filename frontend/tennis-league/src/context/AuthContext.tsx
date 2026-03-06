import React, {
  createContext,
  useContext,
  useEffect,
  useState,
  ReactNode,
} from 'react';

/* ============================= */
/*            TYPES              */
/* ============================= */

export interface AuthUser {
  userID: number;
  name: string;
  surname: string;
  role: string;
}

interface AuthContextType {
  user: AuthUser | null;
  login: (userData: AuthUser) => void;
  logout: () => void;
  isAuthenticated: boolean;
  isLoading: boolean;
}

/* ============================= */
/*        CONTEXT CREATE         */
/* ============================= */

const AuthContext = createContext<AuthContextType | undefined>(undefined);

/* ============================= */
/*        PROVIDER               */
/* ============================= */

interface AuthProviderProps {
  children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [isLoading, setIsLoading] = useState(true); // Yüklenme durumu eklendi

  useEffect(() => {
    const stored = localStorage.getItem('user');
    if (stored) {
      setUser(JSON.parse(stored));
    }
    setIsLoading(false); // Okuma bitti
  }, []);

  const login = (userData: AuthUser) => {
    localStorage.setItem('user', JSON.stringify(userData));
    setUser(userData);
  };

  const logout = () => {
    localStorage.removeItem('user');
    localStorage.removeItem('token');
    setUser(null);
  };

  const value: AuthContextType = {
    user,
    login,
    logout,
    isAuthenticated: !!user,
    isLoading, // Bunu da dışarı aktaralım
  };

  return (
    <AuthContext.Provider value={value}>
      {!isLoading && children} {/* Yükleme bitene kadar çocukları basma */}
    </AuthContext.Provider>
  );
}

/* ============================= */
/*         CUSTOM HOOK           */
/* ============================= */

export function useAuth(): AuthContextType {
  const context = useContext(AuthContext);

  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }

  return context;
}
