import React, { useEffect, useRef, useState } from 'react';
import { Card } from 'primereact/card';
import { Toast } from 'primereact/toast';
import { DataTable } from 'primereact/datatable';
import { Role, RoleLabels, User } from '../../model/user.model';
import { Button } from 'primereact/button';
import { Column } from 'primereact/column';
import { getUsers } from '../../api/userService';

export default function Users() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const toast = useRef<Toast>(null);

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
  const header = () => {
    return (
      <div className="flex justify-content-end">
        <Button label="Yeni Oyuncu" icon="pi pi-plus" />
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
          header={header}
          dataKey="id"
          emptyMessage="Oyuncu bulunamadı"
          loading={loading}
          tableStyle={{ minWidth: '50rem' }}
        >
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
