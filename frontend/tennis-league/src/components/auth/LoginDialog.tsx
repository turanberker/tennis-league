import React, { useState, useRef } from 'react';
import { Dialog } from 'primereact/dialog';
import { InputText } from 'primereact/inputtext';
import { Password } from 'primereact/password';
import { Button } from 'primereact/button';
import { Toast } from 'primereact/toast';
import { login } from '../../api/authService';

interface LoginDialogProps {
  visible: boolean;
  onHide: () => void;
  onLogin: (res: any) => void; // API tipine göre özelleştir
  onShowRegister: () => void;
}

const LoginDialog: React.FC<LoginDialogProps> = ({
  visible,
  onHide,
  onLogin,
  onShowRegister,
}) => {
  const [email, setEmail] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);

  const toast = useRef<Toast>(null);

  const handleLogin = async () => {
    try {
      setLoading(true);
      const res = await login({ email, password });
      if (res) {
        onLogin(res);
        toast.current?.show({
          severity: 'success',
          summary: 'Giriş başarılı',
          detail: `Hoş geldin ${res.currentUser.name}`,
          life: 3000,
        });
        onHide();
      } else {
        toast.current?.show({
          severity: 'error',
          summary: 'Hata',
          detail:  'Giriş başarısız',
          life: 3000,
        });
      }
    } catch (err: any) {
      console.error(err);
      toast.current?.show({
        severity: 'error',
        summary: 'Hata',
        detail: err.message || 'Giriş başarısız',
        life: 3000,
      });
    } finally {
      setLoading(false);
    }
  };

  const footer = (
    <div className="flex justify-content-between w-full">
      <Button label="Kayıt Ol" text onClick={onShowRegister} />
      <Button
        label="Giriş Yap"
        icon="pi pi-sign-in"
        onClick={handleLogin}
        loading={loading}
      />
    </div>
  );

  return (
    <>
      <Toast ref={toast} />
      <Dialog
        header="Giriş Yap"
        visible={visible}
        style={{ width: '400px' }}
        modal
        onHide={onHide}
        footer={footer}
      >
        <div className="flex flex-column gap-3">
          <span className="p-float-label">
            <InputText
              id="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full"
            />
            <label htmlFor="email">Email</label>
          </span>

          <span className="p-float-label">
            <Password
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              toggleMask
              feedback={false}
              className="w-full"
            />
            <label htmlFor="password">Şifre</label>
          </span>
        </div>
      </Dialog>
    </>
  );
};

export default LoginDialog;
