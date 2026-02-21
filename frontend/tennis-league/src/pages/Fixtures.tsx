import React, { useEffect, useRef, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Card } from 'primereact/card';
import { getFixture } from '../api/leagueService';
import {
  LeagueFixtureMatchResponse,
  MatchScore,
  MatchStatusLabels,
  Status,
  TeamRefResponse,
} from '../model/match.model';
import { DataTable } from 'primereact/datatable';
import { Column } from 'primereact/column';
import { Button } from 'primereact/button';
import { OverlayPanel } from 'primereact/overlaypanel';
import { FloatLabel } from 'primereact/floatlabel';
import { Calendar } from 'primereact/calendar';
import {
  getSetScores,
  updateDate,
  updateMatchScore,
} from '../api/matchService';
import { Toast } from 'primereact/toast';
import { Sidebar } from 'primereact/sidebar';
import * as yup from 'yup';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';

/**
 * Ortak score alanı (gte=0,lte=99)
 */
const baseScoreSchema = yup.object().shape({
  team1Score: yup
    .number()
    .typeError('Sayı giriniz')
    .required('Zorunlu')
    .min(0, "0'dan küçük olamaz")
    .max(99, "99'dan büyük olamaz"),

  team2Score: yup
    .number()
    .typeError('Sayı giriniz')
    .required('Zorunlu')
    .min(0, "0'dan küçük olamaz")
    .max(99, "99'dan büyük olamaz"),
});

/**
 * tennis_set validation (backend: tennis_set)
 */
const tennisSetSchema = baseScoreSchema.test(
  'tennis-set',
  'Geçerli set skoru giriniz (6-0..4, 7-5, 7-6)',
  (value) => {
    if (!value) return false;

    const { team1Score, team2Score } = value;

    const max = Math.max(team1Score, team2Score);
    const min = Math.min(team1Score, team2Score);
    const diff = max - min;

    if (max === 6 && min <= 4) return true;
    if (max === 7 && min === 5) return true;
    if (max === 7 && min === 6) return true;

    return false;
  },
);
/**
 * super_tie validation (backend: super_tie)
 */
const superTieSchema = yup
  .object({
    team1Score: yup
      .number()
      .transform((v, o) => (o === '' ? undefined : v))
      .required('Zorunlu')
      .min(0)
      .max(99),

    team2Score: yup
      .number()
      .transform((v, o) => (o === '' ? undefined : v))
      .required('Zorunlu')
      .min(0)
      .max(99),
  })
  .test('super-tie', 'SuperTie min 10 ve 2 fark olmalı', (value) => {
    if (!value) return true;

    const max = Math.max(value.team1Score, value.team2Score);
    const diff = Math.abs(value.team1Score - value.team2Score);

    return max >= 10 && diff >= 2;
  })
  .nullable()
  .default(null);

export const matchScoreSchema = yup.object().shape({
  set1: tennisSetSchema.required(),
  set2: tennisSetSchema.required(),
  superTie: superTieSchema,
});

const initialValues = {
  set1: { team1Score: 0, team2Score: 0 },
  set2: { team1Score: 0, team2Score: 0 },
  superTie: null,
};

export default function Fixtures() {
  const { id } = useParams();
  const [loading, setLoading] = useState<boolean>(false);
  const [matches, setMatches] = useState<LeagueFixtureMatchResponse[]>([]);
  const [updateScoreVisible, setUpdateScoreVisible] = useState<boolean>(false);
  const [showSuperTie, setShowSuperTie] = useState(false);
  const [selectedMatch, setSelectedMatch] = useState<
    LeagueFixtureMatchResponse | undefined
  >();
  const dateOP = useRef<OverlayPanel>(null);
  const toast = useRef<Toast>(null);

  const {
    register,
    handleSubmit,
    reset,
    setValue,
    getValues,
    formState: { errors, isSubmitting },
  } = useForm<MatchScore>({
    resolver: yupResolver(matchScoreSchema),
    defaultValues: initialValues,
  });

  useEffect(() => {
    loadFixture(id!);
  }, [id]);

  const loadFixture = async (id: string) => {
    setLoading(true);
    const data = await getFixture(id);
    setMatches(data);
    setLoading(false);
  };

  const handleMatchDate = async (date: Date) => {
    // Maç tarihini ayarlama işlemi burada yapılacak
    console.log('Tarih ayarla:', selectedMatch);
    if (selectedMatch) {
      await updateDate(selectedMatch?.id, { 'match-date': date });
      setMatches((prev) =>
        prev.map((m) =>
          m.id === selectedMatch.id ? { ...m, matchDate: date } : m,
        ),
      );
      toast.current?.show({
        severity: 'success',
        summary: 'İşlem Başarılı',
        detail: `Maç Tarihi Güncellendi`,
        life: 3000,
      });
    }
    setSelectedMatch(undefined);
  };

  const handleMatchScore = async (match: LeagueFixtureMatchResponse) => {
    console.log(match);
    const setScores = await getSetScores(match.id);
    if (setScores) {
      setValue('set1.team1Score', setScores.set1?.team1Score ?? null);
      setValue('set1.team2Score', setScores.set1?.team2Score ?? null);

      setValue('set2.team1Score', setScores.set2?.team1Score ?? null);
      setValue('set2.team2Score', setScores.set2?.team2Score ?? null);
      if (setScores.superTie) {
        setValue('superTie.team1Score', setScores.superTie?.team1Score ?? null);
        setValue('superTie.team2Score', setScores.superTie?.team2Score ?? null);
        setShowSuperTie(true);
      } else {
        setShowSuperTie(false);
      }
    }
    setUpdateScoreVisible(true);
    setSelectedMatch(match);
  };

  const getButtons = (match: LeagueFixtureMatchResponse) => {
    const completed = match.status === Status.COMPLETED;

    switch (match.status) {
      case Status.PERDING:
        return (
          <>
            {tarihAyarlaButton(match)}
            {match.matchDate ? skorButton(match) : null}
          </>
        );
      case Status.COMPLETED:
        return <> {skorButton(match)}</>;
      case Status.CANCELLED:
        return tarihAyarlaButton(match);
    }
  };

  const tarihAyarlaButton = (match: LeagueFixtureMatchResponse) => {
    return (
      <Button
        rounded
        text
        tooltip="Tarih Ayarla"
        icon="pi pi-calendar"
        outlined
        onClick={(e) => {
          setSelectedMatch(match);
          dateOP.current?.toggle(e);
        }}
      />
    );
  };

  const skorButton = (match: LeagueFixtureMatchResponse) => {
    return (
      <Button
        rounded
        text
        tooltip="Maç Skoru Gir"
        icon="pi pi-pencil"
        outlined
        onClick={() => handleMatchScore(match)}
      />
    );
  };

  const onSubmit = async (data: MatchScore) => {
    if (selectedMatch) {
      try {
        const score = await updateMatchScore(selectedMatch?.id, data);
        console.log(score);
        toast.current?.show({
          severity: 'success',
          summary: 'Başarılı',
          detail: 'Maç Skoru Kaydedilmiştir',
          life: 3000,
        });

        reset();
        setUpdateScoreVisible(false);
        loadFixture(id!);
      } catch (err: any) {
        toast.current?.show({
          severity: 'error',
          summary: 'Hata',
          detail: err.message || 'Skore Kaydedilemedi',
          life: 4000,
        });
      }
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

  const generateTeamName = (team: TeamRefResponse) => {
    return (
      <span style={{ fontWeight: team.winner ? 'bold' : 'normal' }}>
        {team.name}
      </span>
    );
  };

  return (
    <>
      <Toast ref={toast} />
      <Card
        title="Fikstür"
        subTitle="Ligdeki maçların fikstürü burada görüntülenecek."
      >
        <DataTable
          value={matches}
          dataKey={id}
          loading={loading}
          emptyMessage="Fikstür bulunamadı"
          tableStyle={{ minWidth: '50rem' }}
          key="id"
        >
          <Column
            header="1. Takım"
            body={(rowData: LeagueFixtureMatchResponse) =>
              generateTeamName(rowData.team1)
            }
          />
          <Column
            header="2. Takım"
            body={(rowData: LeagueFixtureMatchResponse) =>
              generateTeamName(rowData.team2)
            }
          />
          <Column
            header="Skor"
            body={(rowData: LeagueFixtureMatchResponse) =>
              rowData.status === Status.APPROVED ||
              rowData.status === Status.COMPLETED
                ? rowData.team1.score + '-' + rowData.team2.score
                : '-'
            }
          ></Column>
          <Column
            field="status"
            header="Durum"
            body={(rowData) => MatchStatusLabels[rowData.status as Status]}
          ></Column>
          <Column
            field="matchDate"
            header="Maç Tarihi"
            body={(rowData) =>
              rowData.matchDate
                ? new Date(rowData.matchDate).toLocaleString()
                : '-'
            }
          ></Column>

          <Column
            header="İşlem"
            body={(rowData) => getButtons(rowData)}
          ></Column>
        </DataTable>
      </Card>
      <OverlayPanel ref={dateOP}>
        <FloatLabel>
          <Calendar
            appendTo="self"
            showButtonBar
            showTime
            hourFormat="24"
            inputId="birth_date"
            value={selectedMatch?.matchDate}
            onChange={(e) => {
              handleMatchDate(e.value as Date);
              dateOP.current?.hide();
            }}
          />
          <label htmlFor="birth_date">Maç Tarihi</label>
        </FloatLabel>
      </OverlayPanel>

      <Sidebar
        header="Maç Skoru"
        visible={updateScoreVisible}
        position="right"
        onHide={() => setUpdateScoreVisible(false)}
        icons={() => customIcons}
      >
        <form onSubmit={handleSubmit(onSubmit)} className="p-fluid">
          {/* SET 1 */}
          <h3>Set 1</h3>
          <div className="p-grid">
            <div className="p-col-6">
              <label>{selectedMatch?.team1.name}</label>
              <input
                type="number"
                max={7}
                min={0}
                {...register('set1.team1Score', { valueAsNumber: true })}
                className="p-inputtext"
              />
            </div>

            <div className="p-col-6">
              <label>{selectedMatch?.team2.name}</label>
              <input
                type="number"
                max={7}
                min={0}
                {...register('set1.team2Score', { valueAsNumber: true })}
                className="p-inputtext"
              />
            </div>
          </div>

          {errors.set1?.message && (
            <small className="p-error">{errors.set1.message}</small>
          )}

          {/* SET 2 */}
          <h3 className="mt-4">Set 2</h3>
          <div className="p-grid">
            <div className="p-col-6">
              <label>{selectedMatch?.team1.name}</label>
              <input
                type="number"
                max={7}
                min={0}
                {...register('set2.team1Score', { valueAsNumber: true })}
                className="p-inputtext"
              />
            </div>

            <div className="p-col-6">
              <label>{selectedMatch?.team2.name}</label>
              <input
                type="number"
                max={7}
                min={0}
                {...register('set2.team2Score', { valueAsNumber: true })}
                className="p-inputtext"
              />
            </div>
          </div>
          {errors.set2?.message && (
            <small className="p-error">{errors.set2.message}</small>
          )}

          <div className="mt-4">
            <div className="flex align-items-center gap-2">
              <input
                type="checkbox"
                id="superTieToggle"
                checked={showSuperTie}
                onChange={(e) => {
                  const checked = e.target.checked;
                  setShowSuperTie(checked);

                  if (checked) {
                    setValue('superTie', {
                      team1Score: 0,
                      team2Score: 0,
                    });
                  } else {
                    setValue('superTie', null);
                  }
                }}
              />
              <label htmlFor="superTieToggle">Super Tie Oynandı</label>
            </div>
          </div>

          {/* SUPER TIE */}
          {showSuperTie && (
            <>
              <h3 className="mt-3">Super Tie</h3>

              <div className="p-grid">
                <div className="p-col-6">
                  <label>{selectedMatch?.team1.name}</label>
                  <input
                    type="number"
                    {...register('superTie.team1Score', {
                      valueAsNumber: true,
                    })}
                    className="p-inputtext"
                  />
                </div>

                <div className="p-col-6">
                  <label>{selectedMatch?.team2.name}</label>
                  <input
                    type="number"
                    {...register('superTie.team2Score', {
                      valueAsNumber: true,
                    })}
                    className="p-inputtext"
                  />
                </div>
              </div>

              {errors.superTie?.message && (
                <small className="p-error">{errors.superTie.message}</small>
              )}
            </>
          )}

          {/* BUTTON */}
          <div className="mt-4">
            <button
              type="submit"
              className="p-button p-component"
              disabled={isSubmitting}
            >
              Kaydet
            </button>
          </div>
        </form>
      </Sidebar>
    </>
  );
}
