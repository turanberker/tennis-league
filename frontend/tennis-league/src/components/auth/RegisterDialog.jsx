import React, { useState } from 'react';
import { Dialog } from 'primereact/dialog';
import { InputText } from 'primereact/inputtext';
import { Password } from 'primereact/password';
import { Button } from 'primereact/button';
import Captcha from '../Captcha';

export default function RegisterDialog({ visible, onHide, onRegister }) {
  const [form, setForm] = useState({
    email: '',
    name: '',
    surname: '',
    password: '',
    passwordRepeat: '',
    captcha: '',
  });

  const handleChange = (field, value) => {
    setForm((prev) => ({ ...prev, [field]: value }));
  };

  const footer = (
    <Button
      label="Kayıt Ol"
      icon="pi pi-user-plus"
      onClick={() => onRegister(form)}
    />
  );

  return (
    <Dialog
      header="Kayıt Ol"
      visible={visible}
      style={{ width: '420px' }}
      modal
      onHide={onHide}
      footer={footer}
    >
      <div className="flex flex-column gap-3">
        <span className="p-float-label">
          <InputText
            id="email"
            value={form.email}
            onChange={(e) => handleChange('email', e.target.value)}
            className="w-full"
          />
          <label htmlFor="email">Email</label>
        </span>

        <span className="p-float-label">
          <InputText
            id="name"
            value={form.name}
            onChange={(e) => handleChange('name', e.target.value)}
            className="w-full"
          />
          <label htmlFor="name">Ad</label>
        </span>

        <span className="p-float-label">
          <InputText
            id="surname"
            value={form.surname}
            onChange={(e) => handleChange('surname', e.target.value)}
            className="w-full"
          />
          <label htmlFor="surname">Soyad</label>
        </span>

        <span className="p-float-label">
          <Password
            id="password"
            value={form.password}
            onChange={(e) => handleChange('password', e.target.value)}
            toggleMask
            feedback={false}
            className="w-full"
          />
          <label htmlFor="password">Şifre</label>
        </span>

        <span className="p-float-label">
          <Password
            id="passwordRepeat"
            value={form.passwordRepeat}
            onChange={(e) => handleChange('passwordRepeat', e.target.value)}
            toggleMask
            feedback={false}
            className="w-full"
          />
          <label htmlFor="passwordRepeat">Şifre Tekrar</label>
        </span>

        {/* Basit captcha (placeholder) */}
        <Captcha
          value={form.captchaInput}
          onChange={(val) => handleChange('captchaInput', val)}
        />
      </div>
    </Dialog>
  );
}
