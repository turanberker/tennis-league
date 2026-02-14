import React, { useEffect, useState, useRef } from 'react';
import { Card } from 'primereact/card';
import { Button } from 'primereact/button';
import { Toast } from 'primereact/toast';
import { getPlayers, Player,createPlayer } from '../api/playersService';
import { useNavigate } from 'react-router-dom';
import { Dialog } from 'primereact/dialog';
import * as yup from 'yup';
import { InputText } from 'primereact/inputtext';
import { set, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { classNames } from 'primereact/utils';

const schema = yup.object().shape({
  name: yup
    .string()
    .required("Ad zorunludur")
    .min(3, "Ad en az 3 karakter olmalı")
    .max(75, "Ad en fazla 75 karakter olabilir"),

  surname: yup
    .string()
    .required("Soyad zorunludur")
    .min(3, "Soyad en az 3 karakter olmalı")
    .max(75, "Soyad en fazla 75 karakter olabilir"),
});

interface FormData {
  name: string;
  surname: string;
}

export default function Players() {
const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<FormData>({
    resolver: yupResolver(schema),
    defaultValues: { name: '',surname: '' },
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

   const onSubmit = async (data: FormData) => {
      try {
        await createPlayer(data);
  
        toast.current?.show({
          severity: 'success',
          summary: 'Başarılı',
          detail: 'Lig başarıyla oluşturuldu',
          life: 3000,
        });
  
        setCreateVisible(false);
        loadPlayers(); // listeyi yenile
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
 
  useEffect(() => {
    loadPlayers();
  }, []);


  const handlePlayerDetail = (uuid: string) => {
    navigate(`/players/${uuid}`); // Player detay sayfasına yönlendir
  };

  return (
    <>
      <Toast ref={toast} />
      <Card
        title="Oyuncular"
        subTitle="Mevcut oyuncuları görüntüleyebilir veya yeni oyuncu ekleyebilirsiniz."
        header={
          <div className="flex justify-content-end p-3">
            <Button
              label="Yeni Oyuncu"
              icon="pi pi-plus"
              onClick={()=>setCreateVisible(true)}
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
              <div className="flex flex-column gap-2">
                <label htmlFor="name" className="font-medium">
                  Adı *
                </label>
      
                <InputText
                  id="name"
                  placeholder="Örn: Berker"
                  className={classNames({ 'p-invalid': errors.name })}
                  {...register('name')}
                />
      
                {errors.name && (
                  <small className="p-error">{errors.name.message}</small>
                )}
        </div>
         <div className="flex flex-column gap-2">
                <label htmlFor="name" className="font-medium">
                  Soyadı *
                </label>
      
                <InputText
                  id="surname"
                  placeholder="Örn: Turan"
                  className={classNames({ 'p-invalid': errors.name })}
                  {...register('surname')}
                />
      
                {errors.name && (
                  <small className="p-error">{errors.name.message}</small>
                )}
              </div>
            </Dialog>
    </>
  );
}
