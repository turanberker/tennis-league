import React, {
  createContext,
  useContext,
  useEffect,
  useState,
  ReactNode,
  useCallback,
} from 'react';
import { Role } from '../model/user.model';
import { useNavigate } from 'react-router-dom';
import { logout as logoutApi } from '../api/user/authService'
/* ============================= */
/*            TYPES              */
/* ============================= */

export interface AuthUser {
  userID: string;
  name: string;
  surname: string;
  role: Role;
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
  const navigate = useNavigate(); // Yönlendirme için hook
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

  // useCallback ile sarmaladık ki re-render döngüsüne girmesin
  const logout = useCallback(async () => {
    // 1. Verileri temizle
    localStorage.removeItem('user');
    await logoutApi();
    setUser(null);

    // 2. Kök dizine yönlendir
    // replace: true sayesinde geçmiş (history) temizlenir, geri basınca eski sayfaya dönmez.
    navigate('/', { replace: true });
  }, [navigate]);

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
