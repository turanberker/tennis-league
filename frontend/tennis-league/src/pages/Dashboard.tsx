import React from 'react';
import { Card } from 'primereact/card';
import { useAuth } from '../context/AuthContext';
import { Button } from 'primereact/button';
import ProfileCard from './dashboard/Profile';

import IncomingMatchesCard from './dashboard/IncomingMatches';
import StatisticsCard from './dashboard/Statistics';


export default function Dashboard() {
    const { user, isAuthenticated } = useAuth();

    const authenticatedContext = () => {
        return (<div className="p-4">
            <Card title={`Hoş geldin, ${user?.name || 'Kullanıcı'}!`} className="mb-4">
                <p className="m-0">Sistemdeki güncel durumuna aşağıdan göz atabilirsin.</p>
            </Card>

            <div className="grid">
                {/* 1. Widget: Kullanıcı Özeti */}
                <ProfileCard />

                {/* 2. Widget: İstatistik/Durum */}
                <StatisticsCard />

                {/* 3. Widget: Hızlı Aksiyonlar */}
                {/* <div className="col-12 md:col-4">
                    <Card title="Hızlı İşlemler" style={{ height: '100%' }}>
                        <div className="flex flex-column gap-2">
                            <Button label="Yeni Rapor Oluştur" icon="pi pi-plus" className="p-button-sm" />
                            <Button label="Ayarlara Git" icon="pi pi-cog" className="p-button-sm p-button-outlined" />
                        </div>
                    </Card>
                </div> */}

                {/* 2. Widget: Yaklaşan Maçlar */}
                <IncomingMatchesCard className="col-12 md:col-8" />


            </div>
        </div >)
    }

    return (
        <>
            {isAuthenticated ? (
                authenticatedContext()
            ) : (
                <Card title="Hoş geldiniz!" className="m-4">
                    <p className="m-0">Lütfen giriş yaparak kişiselleştirilmiş içeriğe erişin.</p>
                </Card>
            )}
        </>
    );
}