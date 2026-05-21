import React, { useState } from 'react';


import { PlayerResponse, Sex } from '../model/player.model';
import { getTeamMembers } from '../api/doubleTeamService';
import { ProgressSpinner } from 'primereact/progressspinner';
import { Dialog } from 'primereact/dialog';
import { Button } from 'primereact/button';

interface TeamInfoProps {
    teamName: string;
    teamId: string;
    nameLeftAddon?: React.ReactNode
}

export default function TeamInfo({ teamName, teamId, nameLeftAddon }: TeamInfoProps) {
    const [visible, setVisible] = useState<boolean>(false);
    const [players, setPlayers] = useState<PlayerResponse[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);

    const handleOpenDialog = async () => {
        // Servis çağrısını sadece dialog açıldığında yapıyoruz
        setVisible(true);
        if (players && players.length > 0) return;
        const fetchTeamPlayers = async () => {
            setLoading(true);
            setError(null);
            const players = await getTeamMembers(teamId)
            if (players == null) {
                setError("Oyuncu bilgileri alınırken bir hata oluştu.")
                setLoading(false);
                return;
            }
            setPlayers(players);
            setLoading(false);
        };

        fetchTeamPlayers();
    };

    // Dialog'un alt kısmındaki kapat butonu şablonu
    const dialogFooter = (
        <Button
            label="Kapat"
            icon="pi pi-times"
            onClick={() => setVisible(false)}
            className="p-button-text"
        />
    );

    return (
        <>
            <div className="inline-flex items-center gap-2 font-medium">
                {/* Takım İsmi */}
                {nameLeftAddon}      <span style={{ fontSize: '1.125rem', color: '#1f2937', display: 'inline-flex', alignItems: 'center', }}>{teamName}</span>

                {/* PrimeReact Yerleşik İkonlu Küçük Bilgi Butonu */}
                <Button
                    icon="pi pi-info-circle"
                    rounded
                    text
                    severity="secondary"
                    onClick={handleOpenDialog}
                    tooltip="Oyuncu Bilgileri"
                    tooltipOptions={{ position: 'top' }}
                    style={{
                        width: '2rem',
                        height: '2rem',
                        verticalAlign: 'middle',
                        display: 'inline-flex',
                        alignItems: 'center',
                        justifyContent: 'center'
                    }}
                />
            </div>
            {/* PrimeReact Dialog (Modal) */}
            <Dialog
                header={`${teamName} Kadrosu`}
                visible={visible}
                style={{ width: '320px' }}
                onHide={() => setVisible(false)}
                footer={dialogFooter}
                draggable={false}
                resizable={false}
                breakpoints={{ '960px': '75vw', '641px': '90vw' }}
            >
                <div style={{ minHeight: '100px', display: 'flex', flexDirection: 'column', justifyContent: 'center', gap: '0.75rem' }}>

                    {loading ? (
                        // Yükleniyor Durumu (PrimeReact Spinner)
                        <div className="flex flex-column align-items-center justify-content-center py-4" style={{ gap: '0.5rem', textAlign: 'center' }}>
                            <ProgressSpinner style={{ width: '30px', height: '30px' }} strokeWidth="4" />
                            <span style={{ fontSize: '0.75rem', color: '#6b7280' }}>Oyuncular getiriliyor...</span>
                        </div>
                    ) : error ? (
                        // Hata Durumu
                        <p style={{ fontSize: '0.75rem', color: '#ef4444', textAlign: 'center' }}>{error}</p>
                    ) : players.length > 0 ? (
                        // Oyuncu Listesi (Maksimum 2 Oyuncu)
                        players.map((player, index) => (
                            <div
                                key={index}
                                style={{
                                    display: 'flex',
                                    alignItems: 'center',
                                    justifyContent: 'between',
                                    padding: '0.625rem',
                                    backgroundColor: '#f9fafb',
                                    borderRadius: '8px',
                                    border: '1px solid #f3f4f6'
                                }}
                                className="flex align-items-center justify-content-between"
                            >
                                {/* Sol Taraf: İkon ve İsim */}
                                <div className="flex align-items-center" style={{ gap: '0.5rem' }}>
                                    <i className="pi pi-user" style={{ color: '#9ca3af', fontSize: '0.9rem' }}></i>
                                    <span style={{ fontSize: '0.875rem', color: '#374151', fontWeight: 500 }}>{player.name} {player.surname}</span>
                                    <small className="text-600">Puan: {player.doublePoints}</small>
                                </div>

                                {/* Sağ Taraf: Cinsiyet Badge'i */}
                                <span
                                    style={{
                                        fontSize: '0.75rem',
                                        padding: '0.125rem 0.5rem',
                                        borderRadius: '20px',
                                        fontWeight: 600,
                                        backgroundColor: player.sex === Sex.Female ? '#fdf2f8' : '#eff6ff',
                                        color: player.sex === Sex.Female ? '#db2777' : '#2563eb',
                                        border: `1px solid ${player.sex === Sex.Female ? '#fbcfe8' : '#bfdbfe'}`
                                    }}
                                >
                                    {player.sex}
                                </span>
                            </div>
                        ))
                    ) : (
                        // Boş Liste Durumu
                        <p style={{ fontSize: '0.75rem', color: '#6b7280', textAlign: 'center' }}>Bu takıma ait oyuncu bulunamadı.</p>
                    )}

                </div>
            </Dialog>
        </>
    );
}