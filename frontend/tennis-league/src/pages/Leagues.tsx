import React, { useEffect, useState, useRef } from 'react';
import { Card } from 'primereact/card';
import { InputText } from 'primereact/inputtext';
import { Button } from 'primereact/button';
import { Dialog } from 'primereact/dialog';
import { Toast } from 'primereact/toast';
import { createFixture, getLeagues, saveLeague } from '../api/leagueService';
import * as yup from 'yup';
import { get, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { classNames } from 'primereact/utils';
import { useNavigate } from 'react-router-dom';
import { DataTable } from 'primereact/datatable';
import { Column } from 'primereact/column';

// ================= TYPES =================

interface FormData {
  name: string;
}

// ================= VALIDATION SCHEMA =================
const schema = yup.object({
  name: yup
    .string()
    .required('Lig adı zorunludur.')
    .min(3, 'Lig adı en az 3 karakter olmalıdır.')
    .max(75, 'Lig adı en fazla 75 karakter olabilir.'),
});

export default function Leagues() {
  const navigate = useNavigate();
  const [leagues, setLeagues] = useState<League[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [createVisible, setCreateVisible] = useState<boolean>(false);
  const toast = useRef<Toast>(null);

  // ================= REACT HOOK FORM =================
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<FormData>({
    resolver: yupResolver(schema),
    defaultValues: { name: '' },
  });

  const loadLeagues = () => {
    setLoading(true);
    getLeagues()
      .then((data: League[]) => {
        setLeagues(data);
        setLoading(false);
      })
      .catch((err) => {
        console.error('League fetch error:', err);
        setError('Ligler yüklenemedi');
      });
  };

  useEffect(() => {
    loadLeagues();
  }, []);

  const header = () => {
    return (
      <div className="flex justify-content-end">
        <Button
          label="Yeni Lig Tanımla"
          icon="pi pi-plus"
          onClick={() => setCreateVisible(true)}
        />
      </div>
    );
  };

  const handleStandings = (league: League) => {
    navigate(`/leagues/${league.id}/standings`);
  };

  const handleFixtures = (league: League) => {
    navigate(`/leagues/${league.id}/fixtures`);
  };

  const handleTeams = (league: League) => {
    navigate(`/leagues/${league.id}/teams`);
  };

  const handleCreateLeague = () => {
    reset();
    setCreateVisible(true);
  };

  const handleCreateFixture = async (league: League) => {
    await createFixture(league.id);
    toast.current?.show({
      severity: 'success',
      summary: 'Başarılı',
      detail: 'Fikstür başarıyla oluşturuldu',
      life: 3000,
    });
    loadLeagues();
  };

  const onSubmit = async (data: FormData) => {
    try {
      await saveLeague(data);

      toast.current?.show({
        severity: 'success',
        summary: 'Başarılı',
        detail: 'Lig başarıyla oluşturuldu',
        life: 3000,
      });

      setCreateVisible(false);
      loadLeagues(); // listeyi yenile
    } catch (err: any) {
      console.error(err);

      toast.current?.show({
        severity: 'error',
        summary: 'Hata',
        detail: err.message || 'Lig oluşturulamadı',
        life: 4000,
      });
    }
  };

  const getButtons = (league: League) => {
    const hasFixture = !!league.fixtureCreatedDate;

    return (
      <div className="flex gap-2">
        <Button
          rounded
          text
          label="Takımlar & Oyuncular"
          icon="pi pi-chart-bar"
          outlined
          onClick={() => handleTeams(league)}
        />
        {hasFixture ? (
          <>
            <Button
              rounded
              text
              label="Fikstürü Gör"
              icon="pi pi-calendar"
              outlined
              onClick={() => handleFixtures(league)}
            />
            <Button
              rounded
              text
              label="Puan Durumu"
              icon="pi pi-chart-bar"
              outlined
              onClick={() => handleStandings(league)}
            />
          </>
        ) : (
          <Button
            rounded
            text
            label="Fikstür Oluştur"
            icon="pi pi-plus-circle"
            severity="success"
            onClick={() => handleCreateFixture(league)}
          />
        )}
      </div>
    );
  };

  return (
    <>
      <Toast ref={toast} />
      <Card
        title="Ligler"
        subTitle="Mevcut ligleri görüntüleyebilir veya yeni lig tanımlayabilirsiniz."
      >
        {error && <p style={{ color: 'red' }}>{error}</p>}

        <DataTable
          value={leagues}
          header={header}
          key="id"
          emptyMessage="Lig bulunamadı"
          loading={loading}
          tableStyle={{ minWidth: '50rem' }}
        >
          <Column field="name" header="Lig Adı" />
          <Column
            header="İşlem"
            body={(rowData) => getButtons(rowData)}
          ></Column>
        </DataTable>
      </Card>

      <Dialog
        header="Yeni Lig Tanımla"
        visible={createVisible}
        style={{ width: '400px' }}
        modal
        onHide={() => setCreateVisible(false)}
        footer={
          <div className="flex justify-content-end gap-2">
            <Button
              label="İptal"
              icon="pi pi-times"
              text
              onClick={() => setCreateVisible(false)}
            />
            <Button
              label="Kaydet"
              icon="pi pi-check"
              onClick={handleSubmit(onSubmit)}
              loading={isSubmitting}
            />
          </div>
        }
      >
        <div className="flex flex-column gap-2">
          <label htmlFor="name" className="font-medium">
            Lig Adı *
          </label>

          <InputText
            id="name"
            placeholder="Örn: Süper Lig"
            className={classNames({ 'p-invalid': errors.name })}
            {...register('name')}
          />

          {errors.name && (
            <small className="p-error">{errors.name.message}</small>
          )}
        </div>
      </Dialog>
    </>
  );
}
