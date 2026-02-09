import React, { useRef, useState } from 'react';
import { BrowserRouter as Router } from 'react-router-dom';
import { Button } from 'primereact/button';
import { Sidebar } from 'primereact/sidebar';
import { Avatar } from 'primereact/avatar';
import { Menu } from 'primereact/menu';
import { Toast } from 'primereact/toast';

import 'primeflex/primeflex.css';
import LoginDialog from './components/auth/LoginDialog';
import { AuthProvider, useAuth } from './context/AuthContext';
import { SidebarLinks, AppRoutes } from './router/AppRouter';
import { login as loginApi } from './api/authService'; // backend çağrısı

function Layout() {
  const [sidebarVisible, setSidebarVisible] = useState(false);
  const [loginDialogVisible, setLoginDialogVisible] = useState(false);
  const menuRef = useRef(null);
  const toast = useRef(null);

  const { user, login, logout, isAuthenticated } = useAuth();

  const profileItems = [
    {
      label: 'Profil',
      icon: 'pi pi-user',
      command: () => alert('Profil sayfası'),
    },
    { label: 'Çıkış Yap', icon: 'pi pi-sign-out', command: logout },
  ];

  // LoginDialog’dan çağrılacak fonksiyon
  const handleLogin = (data) => {
    // data = { token, currentUser }
    localStorage.setItem('token', data.token);
    login(data.currentUser); // AuthContext’e kaydet
    setLoginDialogVisible(false);

    toast.current.show({
      severity: 'success',
      summary: 'Giriş başarılı',
      detail: `Hoş geldin ${data.currentUser.name}`,
      life: 3000,
    });
  };

  return (
    <div className="min-h-screen flex flex-column">
      <Toast ref={toast} />
      {/* HEADER */}
      <header
        className="flex align-items-center justify-content-between px-4"
        style={{
          height: '64px',
          borderBottom: '1px solid #e5e7eb',
          background: '#ffffff',
        }}
      >
        <div className="flex align-items-center gap-3">
          <Button
            icon="pi pi-bars"
            text
            rounded
            onClick={() => setSidebarVisible(true)}
          />
          <span style={{ fontSize: '20px', fontWeight: 600 }}>
            🎾 Tennis League
          </span>
        </div>

        <div className="flex align-items-center gap-2">
          {!isAuthenticated ? (
            <Button
              label="Login"
              icon="pi pi-sign-in"
              onClick={() => setLoginDialogVisible(true)}
            />
          ) : (
            <>
              <Avatar
                label={user?.name?.[0] || 'U'}
                shape="circle"
                className="cursor-pointer"
                onClick={(e) => menuRef.current.toggle(e)}
              />
              <Menu model={profileItems} popup ref={menuRef} />
            </>
          )}
        </div>
      </header>

      <LoginDialog
        visible={loginDialogVisible}
        onHide={() => setLoginDialogVisible(false)}
        onLogin={handleLogin}
      />

      {/* BODY */}
      <div className="flex flex-1">
        <Sidebar
          visible={sidebarVisible}
          onHide={() => setSidebarVisible(false)}
          style={{ width: '260px' }}
        >
          <h3 className="mb-4">Menü</h3>
          <SidebarLinks />
        </Sidebar>

        <main className="flex-1 p-4" style={{ background: '#f9fafb' }}>
          <AppRoutes />
        </main>
      </div>
    </div>
  );
}

export default function App() {
  return (
    <Router>
      <AuthProvider>
        <Layout />
      </AuthProvider>
    </Router>
  );
}
