import React, { useEffect, useState, useRef } from 'react';
import { Card } from 'primereact/card';
import { Button } from 'primereact/button';
import { Toast } from 'primereact/toast';
import { getPlayers, Player } from '../api/playersService';
import { useNavigate } from 'react-router-dom';

export default function Players() {
  const [players, setPlayers] = useState<Player[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const toast = useRef<Toast>(null);
  const navigate = useNavigate();

  // Oyuncuları yükle
  const loadPlayers = async () => {
    try {
      setLoading(true);
        const res = await getPlayers();
        console.log('Oyuncular yüklendi:', res);
      setPlayers(res);
      setError(null);
    } catch (err: any) {
      console.error(err);
      setError(err.message || 'Oyuncular yüklenemedi.');
      toast.current?.show({
        severity: 'error',
        summary: 'Hata',
        detail: err.message || 'Oyuncular yüklenemedi.',
        life: 3000,
      });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadPlayers();
  }, []);

  const handleAddPlayer = () => {
    navigate('/players/new'); // Yeni player sayfasına yönlendir
  };

  const handlePlayerDetail = (uuid: string) => {
    navigate(`/players/${uuid}`); // Player detay sayfasına yönlendir
  };

  return (
    <>
      <Toast ref={toast} />
      <Card
        title="Players"
        subTitle="Mevcut oyuncuları görüntüleyebilir veya yeni oyuncu ekleyebilirsiniz."
        header={
          <div className="flex justify-content-end p-3">
            <Button
              label="Yeni Player"
              icon="pi pi-plus"
              onClick={handleAddPlayer}
            />
          </div>
        }
      >
        {error && <p style={{ color: 'red' }}>{error}</p>}

        {loading ? (
          <p>Yükleniyor...</p>
        ) : players.length === 0 ? (
          <p>Oyuncu bulunamadı.</p>
        ) : (
          <div className="flex flex-column gap-3">
            {players.map((player) => (
              <div
                key={player.id}
                className="flex align-items-center justify-content-between p-3 border-round surface-border border-1"
              >
                <span className="font-medium">
                  {player.name} {player.surname}
                </span>

                <Button
                  label="Detay"
                  icon="pi pi-info-circle"
                  outlined
                  onClick={() => handlePlayerDetail(player.uuid)}
                />
              </div>
            ))}
          </div>
        )}
      </Card>
    </>
  );
}
