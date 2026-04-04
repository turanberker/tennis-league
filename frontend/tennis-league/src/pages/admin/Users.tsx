import { useCallback, useEffect, useRef, useState } from 'react';
import { Card } from 'primereact/card';
import { Toast } from 'primereact/toast';
import { DataTable } from 'primereact/datatable';
import { Role, RoleLabels, User } from '../../model/user.model';
import { Button } from 'primereact/button';
import { Column } from 'primereact/column';
import { getUsers } from '../../api/userService';
import { OverlayPanel } from 'primereact/overlaypanel';
import { Dropdown } from 'primereact/dropdown';
import { Player, Sex, SexOptions } from '../../model/player.model';
import { assignPlayerToUser, getUnassignedPlayers } from '../../api/playersService';

export default function Users() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState<boolean>(false);

  const [selectedUser, setSelectedUser] = useState<User>();
  const [unassignedPlayers, setUnassignedPlayes] = useState<Player[]>([]);
  const [selectedPlayer, setSelectedPlayer] = useState<Player | null>();
  const [sexFilter, setSexFilter] = useState<Sex | null>(null);
  const toast = useRef<Toast>(null);

  const op = useRef<OverlayPanel>(null);


  useEffect(() => {

    const loadUnassignedPlayers = async () => {
      if (sexFilter === null) {
        setUnassignedPlayes([]);
        return;
      }
      const res = await getUnassignedPlayers(sexFilter);
      setUnassignedPlayes(res);
    };

    loadUnassignedPlayers();
  }, [sexFilter]);

  const loadUsers = useCallback(async () => {

    setLoading(true);
    const res = await getUsers();
    if (res) {
      setUsers(res);

      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadUsers();
  }, [loadUsers]);

  const playerLabelItemTemplate = (option: Player) => {
    return (
      option ? option.name + ' ' + option.surname : 'Oyuncu seçin'
    ) as string;
  };

  const assignHandler = async () => {
    if (selectedPlayer && selectedUser) {
      const response = await assignPlayerToUser(selectedPlayer.id, { userId: selectedUser.id })
      if (response) {
        toast.current?.show({
          severity: 'success',
          summary: 'İşlem Başarılı',
          detail: `Oyuncu ile kullanıcı eşleşmesi gerçekleşti`,
          life: 3000,
        });
        loadUsers(); // Eşleme sonrası kullanıcıları güncelle
      }

    }

    op.current?.hide(); // İşlem bitince kapat
  }


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
            setSexFilter(null)
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
                assignHandler();
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
