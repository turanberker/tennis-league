import { Card } from "primereact/card";
import { useAuth } from "../../context/AuthContext";

export default function ProfileCard() {

    const { user } = useAuth()

    const userHeader = (
        <div className="flex align-items-center gap-2">
            <i className="pi pi-user"></i>
            <span>Profil Bilgileri</span>
        </div>
    );

    return (<div className="col-12 md:col-4">
        <Card title={userHeader} style={{ height: '100%' }}>
            <p><strong>Tam İsim:</strong> {user?.name} {user?.surname}</p>
            <p><strong>Rol:</strong> {user?.role || 'Standart Kullanıcı'}</p>
        </Card>
    </div>)
}