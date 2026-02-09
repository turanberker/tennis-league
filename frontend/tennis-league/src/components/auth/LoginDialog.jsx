import React, { useState } from 'react';
import { Dialog } from 'primereact/dialog';
import { InputText } from 'primereact/inputtext';
import { Password } from 'primereact/password';
import { Button } from 'primereact/button';
import { Toast } from 'primereact/toast';
import { login } from '../../api/authService';

export default function LoginDialog({ visible, onHide, onLogin }) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);

  const toast = React.useRef(null);

  const handleLogin = async () => {
    try {
      setLoading(true);
      const res = await login({ email, password });
      if (res) {
        onLogin(res); // parent’a gönderiyoruz
        toast.current.show({
          severity: 'success',
          summary: 'Giriş başarılı',
          detail: `Hoş geldin ${res.currentUser.name}`,
          life: 3000,
        });
        onHide();
      } else {
        toast.current.show({
          severity: 'error',
          summary: 'Hata',
          detail: res.errorDetail || 'Giriş başarısız',
          life: 3000,
        });
      }
    } catch (err) {
      console.error(err);
      toast.current.show({
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
      <Button label="Kayıt Ol" text />
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
}
