import React, { useEffect, useRef, useState } from 'react';
import { useForm, Controller, Resolver } from 'react-hook-form';
import * as yup from 'yup';
import { yupResolver } from '@hookform/resolvers/yup';
import { Card } from 'primereact/card';
import { Button } from 'primereact/button';
import { Dialog } from 'primereact/dialog';
import { Toast } from 'primereact/toast';
import { useParams } from 'react-router-dom';
import { CreateTeamRequest, TeamResponse } from '../model/team.model';
import { createTeam, getTeams } from '../api/leagueService';
import { Player } from '../model/player.model';
import { getPlayers } from '../api/playersService';
import { Dropdown } from 'primereact/dropdown';
import { InputText } from 'primereact/inputtext';

type CreateTeamForm = {
  name: string;
  player1: Player | null;
  player2: Player | null;
};

const schema = yup.object({
  name: yup
    .string()
    .required('Takım adı zorunludur')
    .min(5, 'Takım adı en az 5 karakter olmalı')
    .max(75, 'Takım adı en fazla 75 karakter olabilir'),
  player1: yup.mixed<Player>().nullable().required('Birinci oyuncuyu seçin'),
  player2: yup
    .mixed<Player>()
    .nullable()
    .required('İkinci oyuncuyu seçin')
    .test(
      'different-player',
      'İki oyuncu birbirinden farklı olmalıdır',
      function (value) {
        const { player1 } = this.parent;
        if (!value || !player1) return true;
        return (value as Player).id !== (player1 as Player).id;
      },
    ),
}) as yup.Schema<CreateTeamForm>; // TypeScript tip uyumu için

const Teams: React.FC = () => {
  const { id } = useParams<{ id: string }>();

  const [teams, setTeams] = useState<TeamResponse[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [createDialogVisible, setCreateDialogVisible] =
    useState<boolean>(false);

  const toast = useRef<Toast>(null);

  const {
    control,
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<CreateTeamForm>({
    resolver: yupResolver(schema as any),
    defaultValues: { name: '', player1: null, player2: null },
  });

  const [players, setPlayers] = useState<Player[]>([]);

  // Oyuncuları yükle
  useEffect(() => {
    const loadPlayers = async () => {
      try {
        const res = await getPlayers();
        setPlayers(res);
      } catch (err: any) {
        console.error('Oyuncular yüklenirken hata oluştu', err);
        toast.current?.show({
          severity: 'error',
          summary: 'Hata',
          detail: err.message || 'Oyuncular yüklenemedi',
          life: 3000,
        });
      }
    };
    loadPlayers();
  }, []);

  // Takım oluştur
  const onSubmit = async (data: CreateTeamForm) => {
    if (!id) return;

    const payload: CreateTeamRequest = {
      name: data.name,
      playerIds: [data.player1!.id, data.player2!.id],
    };
    const teamId = await createTeam(id, payload);
    setCreateDialogVisible(false);
    reset();
    loadTeams();
  };

  // Lig takımlarını yükle
  const loadTeams = async (): Promise<void> => {
    if (!id) return;

    try {
      setLoading(true);
      console.log('Loading teams for league:', id);
      const res: TeamResponse[] = await getTeams(id);
      setTeams(res);
      setError(null);
    } catch (err: any) {
      console.error('Takımlar yüklenirken hata oluştu:', err);
      setError(err.message || 'Takımlar yüklenemedi');
      toast.current?.show({
        severity: 'error',
        summary: 'Hata',
        detail: err.message || 'Takımlar yüklenemedi',
        life: 3000,
      });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadTeams();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id]);

  const playerLabelItemTemplate = (option: Player) => {
    return (
      option ? option.name + ' ' + option.surname : 'Oyuncu seçin'
    ) as string;
  };

  return (
    <>
      <Toast ref={toast} />

      <Card
        title="Takımlar & Oyuncular"
        header={
          <div className="flex justify-content-end p-3">
            <Button
              label="Yeni Takım"
              icon="pi pi-plus"
              onClick={() => setCreateDialogVisible(true)}
            />
          </div>
        }
      >
        {error && <p style={{ color: 'red' }}>{error}</p>}

        {loading ? (
          <p>Yükleniyor...</p>
        ) : teams.length === 0 ? (
          <p>Takım bulunamadı.</p>
        ) : (
          <div className="flex flex-column gap-3">
            {teams.map((team) => (
              <div
                key={team.id}
                className="p-3 border-1 surface-border border-round"
              >
                <span className="font-medium">{team.name}</span>
              </div>
            ))}
          </div>
        )}
      </Card>

      {/* Yeni Takım Dialog */}
      <Dialog
        header="Yeni Takım Oluştur"
        visible={createDialogVisible}
        style={{ width: '500px' }}
        modal
        onHide={() => setCreateDialogVisible(false)}
        footer={
          <div className="flex justify-content-end gap-2">
            <Button
              label="İptal"
              icon="pi pi-times"
              text
              onClick={() => setCreateDialogVisible(false)}
            />
            <Button
              label="Kaydet"
              icon="pi pi-check"
              onClick={handleSubmit(onSubmit)}
            />
          </div>
        }
      >
        <div className="flex flex-column gap-3">
          <label>Takım Adı *</label>
          <InputText
            {...register('name')}
            className={errors.name ? 'p-invalid' : ''}
            placeholder="Takım adı girin"
          />
          {errors.name && (
            <small className="p-error">{errors.name.message}</small>
          )}

          <label>Oyuncu 1 *</label>
          <Controller
            name="player1"
            control={control}
            render={({ field }) => (
              <Dropdown
                {...field}
                onChange={(e) => field.onChange(e.value)}
                options={players}
                optionLabel="name"
                dataKey="id"
                itemTemplate={playerLabelItemTemplate}
                valueTemplate={playerLabelItemTemplate}
                placeholder="Oyuncu 1 seçin"
                className={errors.player1 ? 'p-invalid' : ''}
              />
            )}
          />
          {errors.player1 && (
            <small className="p-error">{errors.player1.message}</small>
          )}

          <label>Oyuncu 2 *</label>
          <Controller
            name="player2"
            control={control}
            render={({ field }) => (
              <Dropdown
                {...field}
                onChange={(e) => field.onChange(e.value)}
                options={players}
                optionLabel="name"
                dataKey="id"
                itemTemplate={playerLabelItemTemplate}
                valueTemplate={playerLabelItemTemplate}
                placeholder="Oyuncu 2 seçin"
                className={errors.player1 ? 'p-invalid' : ''}
              />
            )}
          />
          {errors.player2 && (
            <small className="p-error">{errors.player2.message}</small>
          )}
        </div>
      </Dialog>
    </>
  );
};

export default Teams;
