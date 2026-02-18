import { useEffect, useRef, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Card } from 'primereact/card';
import { getFixture } from '../api/leagueService';
import {
  LeagueFixtureMatchResponse,
  MatchStatusLabels,
  Status,
} from '../model/match.model';
import { DataTable } from 'primereact/datatable';
import { Column } from 'primereact/column';
import { Button } from 'primereact/button';
import { OverlayPanel } from 'primereact/overlaypanel';
import { FloatLabel } from 'primereact/floatlabel';
import { Calendar } from 'primereact/calendar';
import { updateDate } from '../api/matchService';
import { Toast } from 'primereact/toast';

export default function Fixtures() {
  const { id } = useParams();
  const [loading, setLoading] = useState<boolean>(false);
  const [matches, setMatches] = useState<LeagueFixtureMatchResponse[]>([]);

  const [selectedMatch, setSelectedMatch] = useState<
    LeagueFixtureMatchResponse | undefined
  >();
  const dateOP = useRef<OverlayPanel>(null);
  const toast = useRef<Toast>(null);
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

  const handleMatchScore = (match: LeagueFixtureMatchResponse) => {
    // Maç tarihini ayarlama işlemi burada yapılacak
    console.log('Maç skoru gir:', match);
  };

  const getButtons = (match: LeagueFixtureMatchResponse) => {
    const completed = match.status === Status.COMPLETED;

    switch (match.status) {
      case Status.PERDING:
        return (
          <>
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
            {match.matchDate ? (
              <Button
                rounded
                text
                tooltip="Maç Skoru Gir"
                icon="pi pi-pencil"
                outlined
                onClick={() => handleMatchScore(match)}
              />
            ) : null}
          </>
        );
      case Status.COMPLETED:
        return <> </>;
      case Status.CANCELLED:
        return (
          <Button
            rounded
            text
            label="Tarih Ayarla"
            icon="pi pi-calendar"
            outlined
            onClick={(e) => {
              setSelectedMatch(match);
              dateOP.current?.toggle(e);
            }}
          />
        );
    }
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
          <Column field="team1.name" header="1. Takım" />
          <Column field="team2.name" header="2. Takım" />
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
    </>
  );
}
