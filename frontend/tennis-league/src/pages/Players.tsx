import { useEffect, useState, useRef } from 'react';
import { Card } from 'primereact/card';
import { Button } from 'primereact/button';
import { Toast } from 'primereact/toast';
import { getPlayers, createPlayer } from '../api/playersService';
import { useNavigate } from 'react-router-dom';
import * as yup from 'yup';
import { InputText } from 'primereact/inputtext';
import { Dropdown } from 'primereact/dropdown';
import { useForm, Controller, FormProvider } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { classNames } from 'primereact/utils';
import { Player, Sex, SexLabels, SexOptions } from '../model/player.model';
import { DataTable } from 'primereact/datatable';
import { Column } from 'primereact/column';
import { Sidebar } from 'primereact/sidebar';
import Guard from '../helper/Guard';
import { Role } from '../model/user.model';
import FormItem from '../components/FormItem';
import { isFieldRequired } from '../helper/form.helper';

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

  const methods
    = useForm<FormData>({
      resolver: yupResolver(schema),
      defaultValues: { name: '', surname: '', sex: undefined as any },
    });
  const [sexFilter, setSexFilter] = useState<Sex>()
  const [selectedPlayer, setSelectedPlayer] = useState<Player>()
  const [players, setPlayers] = useState<Player[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [createVisible, setCreateVisible] = useState<boolean>(false);

  const toast = useRef<Toast>(null);
  const navigate = useNavigate();

  useEffect(() => { loadPlayers() }, [sexFilter])

  // Oyuncuları yükle
  const loadPlayers = async () => {

    setLoading(true);
    const res = await getPlayers({ sex: sexFilter });
    setPlayers(res);
    setLoading(false);
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

      methods.reset();
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
      <>
        <div className="flex flex-column md:flex-row md:justify-content-between md:align-items-center gap-3 p-2">
          {/* SOL TARAF: Başlık ve Filtreler */}
          <div className="flex align-items-center flex-1">

            <Dropdown
              value={sexFilter}
              onChange={(e) => setSexFilter(e.value)}
              options={SexOptions}
              optionLabel="label"
              showClear
              placeholder="Cinsiyet"
              className="w-full md:w-14rem" // Sabit genişlik filtrede daha iyi durur
            />
          </div>

          {/* SAĞ TARAF: Aksiyon Butonları */}
          <div className="flex align-items-center gap-2">
            <Guard allowedRoles={[Role.ADMIN, Role.COORDINATOR]}>
              <Button
                label="Yeni Oyuncu"
                icon="pi pi-plus"
                size="small"
                onClick={() => setCreateVisible(true)}
              />
            </Guard>
            <Button
              severity="info"
              label="Detay"
              icon="pi pi-search"
              disabled={!selectedPlayer}
              outlined // Detay butonu ikincil olduğu için outlined daha şık durabilir
              size="small"
              onClick={() => handlePlayerDetail(selectedPlayer!.id)}
            />


          </div>
        </div>

      </>
    );
  };
  return (
    <>
      <Toast ref={toast} />

      <Card
        title="Oyuncular"
        subTitle="Mevcut oyuncuları görüntüleyebilir veya yeni oyuncu ekleyebilirsiniz."
      >
        <DataTable
          value={players}
          header={header}
          key="id"
          dataKey="id"
          selectionMode="single"
          selection={selectedPlayer!}
          onSelectionChange={(e) => setSelectedPlayer(e.value)}
          emptyMessage="Oyuncu bulunamadı"
          loading={loading}
          tableStyle={{ minWidth: '50rem' }}
        >
          <Column
            selectionMode="single"
            headerStyle={{ width: "3rem" }}
          ></Column>
          <Column field="name" header="Ad" />
          <Column field="surname" header="Soyad" />
          <Column
            field="sex"
            header="Cinsiyet"
            body={(rowData) => SexLabels[rowData.sex as Sex]}
          />

        </DataTable>
      </Card>

      <Sidebar
        visible={createVisible}
        position="right"
        onHide={() => setCreateVisible(false)}

      >
        <FormProvider {...methods}>
          <form onSubmit={methods.handleSubmit(onSubmit)} className="p-fluid">
            <FormItem
              label='Adı'
              name='name'
              required={isFieldRequired(schema, "name")}>
              <InputText
                id="name"
                placeholder="Örn: Novak"
                className={classNames({ 'p-invalid': methods.formState.errors.name })}
                {...methods.register('name')}
              />
            </FormItem>
            <FormItem
              label='Soyadı'
              name='surname'
              required={isFieldRequired(schema, "surname")}>
              <InputText
                placeholder="Örn: Djokovic"
                id="surname"
                className={classNames({ 'p-invalid': methods.formState.errors.surname })}
                {...methods.register('surname')}
              />
            </FormItem>
            <FormItem
              label='Cinsiyet'
              name='sex'
              required={isFieldRequired(schema, "sex")}>

              <Controller
                name="sex"
                control={methods.control}
                render={({ field }) => (
                  <Dropdown
                    {...field}
                    options={SexOptions}
                    placeholder="Cinsiyet seçiniz"
                    className={classNames({ 'p-invalid': methods.formState.errors.sex })}
                  />
                )}
              />
            </FormItem>
            <div className="mt-4 flex gap-2">
              <Button
                type="button"
                label="İptal"
                icon="pi pi-times"
                outlined
                severity="secondary"
                onClick={() => setCreateVisible(false)}
                className="w-full"
              />
              <Button
                type="submit"
                label="Kaydet"
                icon="pi pi-check"
                loading={methods.formState.isSubmitting}
                className="w-full"
              />
            </div>
          </form>

        </FormProvider>

      </Sidebar>
      {/* CREATE DIALOG */}
    </>
  );
}
