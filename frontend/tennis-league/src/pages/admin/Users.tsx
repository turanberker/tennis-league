import React, { useEffect, useRef, useState } from 'react';
import { Card } from 'primereact/card';
import { Toast } from 'primereact/toast';
import { DataTable } from 'primereact/datatable';
import { Role, RoleLabels, User } from '../../model/user.model';
import { Button } from 'primereact/button';
import { Column } from 'primereact/column';
import { getUsers } from '../../api/userService';
import { OverlayPanel } from 'primereact/overlaypanel';
import { Dropdown } from 'primereact/dropdown';
import Players from '../Players';
import { Player, Sex, SexOptions } from '../../model/player.model';
import { getUnassignedPlayers } from '../../api/playersService';

export default function Users() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [selectedUser, setSelectedUser] = useState<User>();
  const [unassignedPlayers, setUnassignedPlayes] = useState<Player[]>([]);
  const [selectedPlayer, setSelectedPlayer] = useState<Player | null>();
  const [sexFilter, setSexFilter] = useState<Sex | null>(null);
  const toast = useRef<Toast>(null);

  const op = useRef<OverlayPanel>(null);
  useEffect(() => {
    loadUnassignedPlayers();
  }, [sexFilter]);

  useEffect(() => {
    loadUsers();
  }, []);
  // Oyuncuları yükle
  const loadUsers = async () => {
    try {
      console.log('1');
      setLoading(true);
      const res = await getUsers();
      setUsers(res);
      setError(null);
    } catch (err: any) {
      setError(err.message || 'Kullanıcılar yüklenemedi.');
      toast.current?.show({
        severity: 'error',
        summary: 'Hata',
        detail: err.message || 'Kullanıcılar yüklenemedi.',
        life: 3000,
      });
    } finally {
      setLoading(false);
    }
  };

  const loadUnassignedPlayers = async () => {
    if (sexFilter === null) return;
    else {
      const res = await getUnassignedPlayers(sexFilter);
      setUnassignedPlayes(res);
    }
  };

  const playerLabelItemTemplate = (option: Player) => {
    return (
      option ? option.name + ' ' + option.surname : 'Oyuncu seçin'
    ) as string;
  };

  const header = () => {
    return (
      <div className="flex justify-content-end">
        <Button
          label="Oyuncu ile eşle"
          tooltip={selectedUser?.playerId ? 'Bu oyuncu zaten eşleşmiş' : ''}
          tooltipOptions={{ showOnDisabled: true }}
          disabled={!selectedUser || !!selectedUser.playerId}
          icon="pi pi-plus"
          size="small"
          onClick={(e) => {
            setUnassignedPlayes([]);
            setSelectedPlayer(null);
            op.current?.toggle(e);
          }} // Tıklayınca Overlay açılır
        />

        <OverlayPanel ref={op} showCloseIcon>
          <div className="flex flex-column gap-3" style={{ width: '250px' }}>
            <label className="font-bold">Eşleşecek Oyuncuyunun Cinsiyeti</label>

            <Dropdown
              value={sexFilter}
              onChange={(e) => setSexFilter(e.value)}
              options={SexOptions}
              optionLabel="label"
              placeholder="Cinsiyet Listesi"
              className="w-full"
            />
            <Dropdown
              value={selectedPlayer}
              onChange={(e) => setSelectedPlayer(e.value)}
              options={unassignedPlayers}
              filterBy="name,surname"
              itemTemplate={playerLabelItemTemplate}
              optionLabel="name"
              placeholder="Oyuncu listesi"
              className="w-full"
              filter // Arama özelliği (isteğe bağlı)
            />

            <Button
              label="Eşleşmeyi Onayla"
              icon="pi pi-check"
              size="small"
              disabled={!selectedPlayer}
              onClick={() => {
                // Eşleme mantığınız buraya
                console.log(
                  `${selectedUser?.name} ile ${selectedPlayer?.name} eşleşti.`,
                );
                op.current?.hide(); // İşlem bitince kapat
              }}
            />
          </div>
        </OverlayPanel>
      </div>
    );
  };

  return (
    <>
      {' '}
      <Toast ref={toast} />
      <Card title="Kullanıcılar">
        {error && <p style={{ color: 'red' }}>{error}</p>}

        <DataTable
          value={users}
          selection={selectedUser!}
          onSelectionChange={(e) => setSelectedUser(e.value as User)} 
          header={header}
          dataKey="id"
          selectionMode="single"
          emptyMessage="Oyuncu bulunamadı"
          loading={loading}
          tableStyle={{ minWidth: '50rem' }}
        >
          <Column
            selectionMode="single"
            headerStyle={{ width: '3rem' }}
          ></Column>
          <Column field="name" header="Ad" />
          <Column field="surname" header="Soyad" />
          <Column
            field="role"
            header="Rolü"
            body={(rowData) => RoleLabels[rowData.role as Role]}
          />
          <Column
            field="playerId"
            header="Oyuncu Kaydı"
            body={(rowData) => (rowData.playerId ? 'Evet' : 'Hayır')}
          />
        </DataTable>
      </Card>
    </>
  );
}
