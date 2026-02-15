import React, { useEffect, useRef, useState } from 'react';
import { BrowserRouter as Router } from 'react-router-dom';
import { Button } from 'primereact/button';
import { Sidebar } from 'primereact/sidebar';
import { Avatar } from 'primereact/avatar';
import { Menu } from 'primereact/menu';
import { Toast } from 'primereact/toast';

import 'primeflex/primeflex.css';
import LoginDialog from './components/auth/LoginDialog';
import RegisterDialog, {
  RegisterForm,
} from './components/auth/RegisterDialog';
import { SidebarLinks, AppRoutes } from './router/AppRouter';
import { AuthProvider, useAuth, AuthUser } from './context/AuthContext';
import { login as loginApi, register as registerApi } from './api/authService';
import { registerToast } from './api/toastService';

function Layout() {
  const [sidebarVisible, setSidebarVisible] = useState(false);
  const [loginDialogVisible, setLoginDialogVisible] = useState(false);
  const [registerDialogVisible, setRegisterDialogVisible] = useState(false);
  const menuRef = useRef<Menu>(null);
  const toast = useRef<Toast>(null);

useEffect(() => {
  const handler = (e: any) => {
    toast.current?.show({
      severity: 'error',
      summary: 'Hata',
      detail: e.detail,
      life: 3000,
    });
  };

  window.addEventListener('api-error', handler);
  return () => window.removeEventListener('api-error', handler);
}, []);
  const { user, login, logout, isAuthenticated } = useAuth();

  const profileItems = [
    {
      label: 'Profil',
      icon: 'pi pi-user',
      command: () => alert('Profil sayfası'),
    },
    { label: 'Çıkış Yap', icon: 'pi pi-sign-out', command: logout },
  ];

  /* -----------------------------
      LOGIN HANDLER
  ----------------------------- */
  const handleLogin = (data: { token: string; currentUser: AuthUser }) => {
    localStorage.setItem('token', data.token);
    login(data.currentUser);
    setLoginDialogVisible(false);

    toast.current?.show({
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
        style={{ height: 64, borderBottom: '1px solid #e5e7eb', background: '#fff' }}
      >
        <div className="flex align-items-center gap-3">
          <Button
            icon="pi pi-bars"
            text
            rounded
            onClick={() => setSidebarVisible(true)}
          />
          <span style={{ fontSize: 20, fontWeight: 600 }}>🎾 Tennis League</span>
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
                label={user?.name?.[0] ?? 'U'}
                shape="circle"
                className="cursor-pointer"
                onClick={(e) => menuRef.current?.toggle(e)}
              />
              <Menu model={profileItems} popup ref={menuRef} />
            </>
          )}
        </div>
      </header>

      {/* DIALOGS */}
      <LoginDialog
        visible={loginDialogVisible}
        onHide={() => setLoginDialogVisible(false)}
        onLogin={handleLogin}
        onShowRegister={() => {
          setLoginDialogVisible(false);
          setRegisterDialogVisible(true);
        }}
      />

      <RegisterDialog
        visible={registerDialogVisible}
        onHide={() => setRegisterDialogVisible(false)}
        onRegister={async (form: RegisterForm) => {
          try {
            const { captchaInput, passwordRepeat, ...payload } = form;

            const data = await registerApi(payload);

            localStorage.setItem('token', data.token);
            login(data.currentUser);
            setRegisterDialogVisible(false);

            toast.current?.show({
              severity: 'success',
              summary: 'Kayıt başarılı',
              detail: `Hoş geldin ${data.currentUser.name}`,
              life: 3000,
            });
          } catch (err: any) {
            toast.current?.show({
              severity: 'error',
              summary: 'Hata',
              detail: err.message || 'Kayıt başarısız',
              life: 3000,
            });
          }
        }}
      />

      {/* BODY */}
      <div className="flex flex-1">
        <Sidebar
          visible={sidebarVisible}
          onHide={() => setSidebarVisible(false)}
          style={{ width: 260 }}
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
