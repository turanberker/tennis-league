import React, { useEffect, useState, useRef } from 'react';
import { Card } from 'primereact/card';
import { InputText } from 'primereact/inputtext';
import { Button } from 'primereact/button';
import { Dialog } from 'primereact/dialog';
import { Toast } from 'primereact/toast';
import { getLeagues, saveLeague } from '../api/leagueService.ts';
import * as yup from 'yup';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { classNames } from 'primereact/utils';

// ================= VALIDATION SCHEMA =================
const schema = yup.object({
  name: yup
    .string()
    .required('Lig adı zorunludur.')
    .min(3, 'Lig adı en az 3 karakter olmalıdır.')
    .max(75, 'Lig adı en fazla 75 karakter olabilir.'),
});

export default function Leagues() {
  const [leagues, setLeagues] = useState([]);
  const [error, setError] = useState(null);
  const [createVisible, setCreateVisible] = useState(false);
  const toast = useRef(null);
  // ================= REACT HOOK FORM =================
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm({
    resolver: yupResolver(schema),
    defaultValues: { name: '' },
  });

  const loadLeagues = () => {
    getLeagues()
      .then(setLeagues)
      .catch((err) => {
        console.error('League fetch error:', err);
        setError('Ligler yüklenemedi');
      });
  };

  useEffect(() => {
    loadLeagues();
  }, []);

  const handleStandings = (league) => {
    console.log('Puan durumu:', league);
  };

  const handleFixtures = (league) => {
    console.log('Fikstür:', league);
  };

  const handleCreateLeague = () => {
    reset();
    setCreateVisible(true);
  };

  const onSubmit = async (data) => {
    try {
      await saveLeague(data);

      toast.current.show({
        severity: 'success',
        summary: 'Başarılı',
        detail: 'Lig başarıyla oluşturuldu',
        life: 3000,
      });

      setCreateVisible(false);
      loadLeagues(); // listeyi yenile
    } catch (err) {
      console.error(err);

      toast.current.show({
        severity: 'error',
        summary: 'Hata',
        detail: err.message || 'Lig oluşturulamadı',
        life: 4000,
      });
    }
  };

  return (
    <>
      <Toast ref={toast} />
      <Card
        title="Ligler"
        subTitle="Mevcut ligleri görüntüleyebilir veya yeni lig tanımlayabilirsiniz."
        header={
          <div className="flex justify-content-end p-3">
            <Button
              label="Yeni Lig Tanımla"
              icon="pi pi-plus"
              onClick={handleCreateLeague}
            />
          </div>
        }
      >
        {error && <p style={{ color: 'red' }}>{error}</p>}

        {leagues.length === 0 && !error ? (
          <p>Lig bulunamadı.</p>
        ) : (
          <div className="flex flex-column gap-3">
            {leagues.map((league) => (
              <div
                key={league.id}
                className="flex align-items-center justify-content-between p-3 border-round surface-border border-1"
              >
                <span className="font-medium">{league.name}</span>

                <div className="flex gap-2">
                  <Button
                    label="Puan Durumu"
                    icon="pi pi-chart-bar"
                    outlined
                    onClick={() => handleStandings(league)}
                  />

                  <Button
                    label="Fikstür"
                    icon="pi pi-calendar"
                    outlined
                    onClick={() => handleFixtures(league)}
                  />
                </div>
              </div>
            ))}
          </div>
        )}
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
