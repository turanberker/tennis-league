import React, { useEffect, useState, useRef } from 'react';
import { Card } from 'primereact/card';
import { InputText } from 'primereact/inputtext';
import { Button } from 'primereact/button';
import { Toast } from 'primereact/toast';
import { createFixture, getLeagues, saveLeague } from '../api/leagueService';
import * as yup from 'yup';
import { get, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { classNames } from 'primereact/utils';
import { useNavigate } from 'react-router-dom';
import { DataTable } from 'primereact/datatable';
import { Column } from 'primereact/column';
import { Sidebar } from 'primereact/sidebar';
import { formatDate } from '../helper/date.helper';

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
  const [selectedLeague, setSelectedLeague] = useState<League | null>();
  const [loading, setLoading] = useState<boolean>(false);
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

  const loadLeagues = async () => {
    setLoading(true);
    // Hata olursa 'data' null gelecek ve alt satırlar patlamayacak
    const data = await getLeagues();
    if (data) {
      setLeagues(data);
    }
    setLoading(false); // Hata ol
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
          size="small"
          onClick={() => setCreateVisible(true)}
        />
        <Button
          disabled={!selectedLeague}
          rounded
          text
          label="Takımlar & Oyuncular"
          icon="pi pi-chart-bar"
          outlined
          size="small"
          onClick={() => handleTeams()}
        />

        <Button
          disabled={!selectedLeague || !selectedLeague.fixtureCreatedDate}
          rounded
          text
          label="Fikstürü Gör"
          icon="pi pi-calendar"
          outlined
          size="small"
          onClick={() => handleFixtures()}
        />
        <Button
          rounded
          disabled={!selectedLeague || !selectedLeague.fixtureCreatedDate}
          text
          label="Puan Durumu"
          icon="pi pi-chart-bar"
          outlined
          size="small"
          onClick={() => handleStandings()}
        />
        <Button
          rounded
          disabled={!selectedLeague || !!selectedLeague.fixtureCreatedDate}
          text
          label="Fikstür Oluştur"
          icon="pi pi-plus-circle"
          severity="success"
          size="small"
          onClick={() => handleCreateFixture()}
        />
      </div>
    );
  };

  const handleStandings = () => {
    navigate(`/leagues/${selectedLeague!.id}/standings`);
  };

  const handleFixtures = () => {
    navigate(`/leagues/${selectedLeague!.id}/fixtures`);
  };

  const handleTeams = () => {
    navigate(`/leagues/${selectedLeague!.id}/teams`);
  };

  const handleCreateFixture = async () => {
    const data = await createFixture(selectedLeague!.id);
    if (data) {
      toast.current?.show({
        severity: 'success',
        summary: 'Başarılı',
        detail: 'Fikstür başarıyla oluşturuldu',
        life: 3000,
      });
      loadLeagues();
    }
  };

  const onSubmit = async (data: FormData) => {
    try {
      await saveLeague(data);
      reset();
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

  const customIcons = (
    <React.Fragment>
      <Button
        className="p-sidebar-icon p-link mr-2"
        icon="pi pi-check"
        tooltip="Kaydet"
        onClick={handleSubmit(onSubmit)}
        loading={isSubmitting}
        tooltipOptions={{ position: 'left' }}
      ></Button>
    </React.Fragment>
  );
  return (
    <>
      <Toast ref={toast} />
      <Card
        title="Ligler"
        subTitle="Mevcut ligleri görüntüleyebilir veya yeni lig tanımlayabilirsiniz."
      >
        <DataTable
          value={leagues}
          header={header}
          key="id"
          emptyMessage="Lig bulunamadı"
          loading={loading}
          tableStyle={{ minWidth: '50rem' }}
          dataKey="id"
          selectionMode="single"
          selection={selectedLeague!}
          onSelectionChange={(e) => setSelectedLeague(e.value as League)}
        >
          <Column
            selectionMode="single"
            headerStyle={{ width: '3rem' }}
          ></Column>
          <Column field="name" header="Lig Adı" />
          <Column  body={(league:League)=>formatDate(league.fixtureCreatedDate) } header="Lig Başlangıç Tarihi"/>
          <Column  body={(league:League)=>league.coordinators.join(",")} header="Koordinatörler"/>
        </DataTable>
      </Card>

      <Sidebar
        header="Yeni Lig Tanımla"
        visible={createVisible}
        position="right"
        onHide={() => setCreateVisible(false)}
        icons={() => customIcons}
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
      </Sidebar>
    </>
  );
}
