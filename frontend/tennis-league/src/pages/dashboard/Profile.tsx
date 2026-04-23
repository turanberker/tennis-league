import { Card } from "primereact/card";
import { useAuth } from "../../context/AuthContext";
import { Button } from "primereact/button";
import { useState } from "react";
import ChangePasswordDialog from "./ChangePasswordDialog";

interface ProfileCardProps extends DashboardProps {

}

export default function ProfileCard({ className = "col-12 md:col-4" }: ProfileCardProps) {

    const { user } = useAuth()
    const [changePasswordVisible, setChangePasswordVisible] = useState(false);
    const userHeader = (
        <div className="flex align-items-center gap-2">
            <i className="pi pi-user"></i>
            <span>Profil Bilgileri</span>
        </div>
    );

    const footer = (
        <div className="flex justify-content-end">
            <Button
                label="Şifre Değiştir"
                icon="pi pi-key"
                className="p-button-outlined p-button-sm p-button-warning"
                onClick={() => setChangePasswordVisible(true)}
            />
        </div>
    );

    return (<div className={className}>
        <Card title={userHeader} style={{ height: '100%' }} footer={footer}>
            <p><strong>İsim:</strong> {user?.name} {user?.surname}</p>
            <p><strong>Rol:</strong> {user?.role || 'Standart Kullanıcı'}</p>
        </Card>

        <ChangePasswordDialog
            visible={changePasswordVisible}
            onHide={() => setChangePasswordVisible(false)}
        />
    </div>)
}