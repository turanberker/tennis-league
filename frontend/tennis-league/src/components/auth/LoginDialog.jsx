import React, { useState } from "react";
import { Dialog } from "primereact/dialog";
import { InputText } from "primereact/inputtext";
import { Password } from "primereact/password";
import { Button } from "primereact/button";

// 🔐 Login Dialog Component
export default function LoginDialog({ visible, onHide, onLogin }) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const footer = (
    <div className="flex justify-content-between w-full">
      <Button label="Kayıt Ol" text />
      <Button label="Giriş Yap" icon="pi pi-sign-in" onClick={onLogin} />
    </div>
  );

  return (
    <Dialog
      header="Giriş Yap"
      visible={visible}
      style={{ width: "400px" }}
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
  );
}

