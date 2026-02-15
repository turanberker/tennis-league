import React, { useEffect, useState, useRef } from 'react';
import { Card } from 'primereact/card';
import { Button } from 'primereact/button';
import { Toast } from 'primereact/toast';
import { getPlayers, createPlayer } from '../api/playersService';

import { useNavigate } from 'react-router-dom';
import { Dialog } from 'primereact/dialog';
import * as yup from 'yup';
import { InputText } from 'primereact/inputtext';
import { Dropdown } from 'primereact/dropdown';
import { useForm, Controller } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { classNames } from 'primereact/utils';
import { Player, Sex, SexLabels, SexOptions } from '../model/player.model';
import { DataTable } from 'primereact/datatable';
import { Column } from 'primereact/column';

// VALIDATION
const schema = yup.object().shape({
  name: yup
    .string()
    .required('Ad zorunludur')
    .min(3, 'Ad en az 3 karakter olmalı')
    .max(75, 'Ad en fazla 75 karakter olabilir'),

  surname: yup
    .string()
    .required('Soyad zorunludur')
    .min(3, 'Soyad en az 3 karakter olmalı')
    .max(75, 'Soyad en fazla 75 karakter olabilir'),

  sex: yup
    .mixed<Sex>()
    .oneOf(Object.values(Sex), 'Cinsiyet seçiniz')
    .required('Cinsiyet zorunludur'),
});

interface FormData {
  name: string;
  surname: string;
  sex: Sex;
}

export default function Players() {
  const {
    register,
    handleSubmit,
    reset,
    control,
    formState: { errors, isSubmitting },
  } = useForm<FormData>({
    resolver: yupResolver(schema),
    defaultValues: { name: '', surname: '', sex: undefined as any },
  });

  const [players, setPlayers] = useState<Player[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [createVisible, setCreateVisible] = useState<boolean>(false);

  const toast = useRef<Toast>(null);
  const navigate = useNavigate();

  // Oyuncuları yükle
  const loadPlayers = async () => {
    try {
      setLoading(true);
      const res = await getPlayers();
      setPlayers(res);
      setError(null);
    } catch (err: any) {
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

  // CREATE
  const onSubmit = async (data: FormData) => {
    try {
      await createPlayer(data);

      toast.current?.show({
        severity: 'success',
        summary: 'Başarılı',
        detail: 'Oyuncu başarıyla oluşturuldu',
        life: 3000,
      });

      reset();
      setCreateVisible(false);
      loadPlayers();
    } catch (err: any) {
      toast.current?.show({
        severity: 'error',
        summary: 'Hata',
        detail: err.message || 'Oyuncu oluşturulamadı',
        life: 4000,
      });
    }
  };

  useEffect(() => {
    loadPlayers();
  }, []);

  const handlePlayerDetail = (uuid: string) => {
    navigate(`/players/${uuid}`);
  };

  const header = () => {
    return (
      <div className="flex justify-content-end">
        <Button
          label="Yeni Oyuncu"
          icon="pi pi-plus"
          onClick={() => setCreateVisible(true)}
        />
      </div>
    );
  };
  return (
    <>
      <Toast ref={toast} />

      <Card
        title="Oyuncular"
        subTitle="Mevcut oyuncuları görüntüleyebilir veya yeni oyuncu ekleyebilirsiniz."
      >
        {error && <p style={{ color: 'red' }}>{error}</p>}

        <DataTable
          value={players}
          header={header}
          key="id"
          emptyMessage="Oyuncu bulunamadı"
          loading={loading}
          tableStyle={{ minWidth: '50rem' }}
        >
          <Column field="name" header="Ad" />
          <Column field="surname" header="Soyad" />
          <Column
            field="sex"
            header="Cinsiyet"
            body={(rowData) => SexLabels[rowData.sex as Sex]}
          />
          <Column
            header="İşlem"
            body={(rowData) => (
              <Button
                icon="pi pi-info-circle"
                severity="info"
                rounded
                text
                onClick={() => handlePlayerDetail(rowData.id)}
              />
            )}
          />
        </DataTable>
      </Card>

      {/* CREATE DIALOG */}
      <Dialog
        header="Yeni Oyuncu Tanımla"
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
        {/* NAME */}
        <div className="flex flex-column gap-2">
          <label className="font-medium">Adı *</label>
          <InputText
            placeholder="Örn: Berker"
            className={classNames({ 'p-invalid': errors.name })}
            {...register('name')}
          />
          {errors.name && (
            <small className="p-error">{errors.name.message}</small>
          )}
        </div>

        {/* SURNAME */}
        <div className="flex flex-column gap-2 mt-3">
          <label className="font-medium">Soyadı *</label>
          <InputText
            placeholder="Örn: Turan"
            className={classNames({ 'p-invalid': errors.surname })}
            {...register('surname')}
          />
          {errors.surname && (
            <small className="p-error">{errors.surname.message}</small>
          )}
        </div>

        {/* SEX */}
        <div className="flex flex-column gap-2 mt-3">
          <label className="font-medium">Cinsiyet *</label>

          <Controller
            name="sex"
            control={control}
            render={({ field }) => (
              <Dropdown
                {...field}
                options={SexOptions}
                placeholder="Cinsiyet seçiniz"
                className={classNames({ 'p-invalid': errors.sex })}
              />
            )}
          />

          {errors.sex && (
            <small className="p-error">{errors.sex.message}</small>
          )}
        </div>
      </Dialog>
    </>
  );
}
