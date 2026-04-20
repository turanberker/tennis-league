import { useEffect, useRef, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Card } from 'primereact/card';
import { approveMatchResult, getFixture, updateMatchDate } from '../../api/leagueService';
import {
  LeagueFixtureMatchResponse,
  MatchStatusLabels,
  Status,
  TeamRefResponse,
} from '../../model/match.model';
import { DataTable } from 'primereact/datatable';
import { Column } from 'primereact/column';
import { Button } from 'primereact/button';
import { OverlayPanel } from 'primereact/overlaypanel';
import { FloatLabel } from 'primereact/floatlabel';
import { Calendar } from 'primereact/calendar';
import { Toast } from 'primereact/toast';
import { useLeague } from '../../hooks/useLeague';
import { LeagueCard } from '../../components/LeagueCard';
import { MatchScoreSidebar } from '../../components/match/MatchScoreSidebar';

export default function Fixtures() {
  const { id } = useParams();
  const [loading, setLoading] = useState<boolean>(false);
  const [matches, setMatches] = useState<LeagueFixtureMatchResponse[]>([]);
  const [updateScoreVisible, setUpdateScoreVisible] = useState<boolean>(false);
  const [selectedMatch, setSelectedMatch] = useState<LeagueFixtureMatchResponse>();
  const dateOP = useRef<OverlayPanel>(null);
  const toast = useRef<Toast>(null);

  const { data: league } = useLeague(id);

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

    if (selectedMatch && league) {
      const res = await updateMatchDate(league?.id, selectedMatch?.id, { 'match-date': date });
      if (res) {
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
    }
  };

  const handleMatchScore = async () => {

    setUpdateScoreVisible(true);
  };

  const approveHandler = async () => {
    const response = await approveMatchResult(league!.id, selectedMatch!.id);
    if (response) {
      toast.current?.show({
        severity: 'success',
        summary: 'İşlem Başarılı',
        detail: `Maç sonucu onaylandı, skorbord en kısa sürede güncellenecektir.`,
        life: 3000,
      });
      loadFixture(id!);
    }

  };

  const generateTeamName = (team: TeamRefResponse) => {
    return (
      <span style={{ fontWeight: team.winner ? 'bold' : 'normal' }}>
        {team.name}
      </span>
    );
  };

  const header = () => {
    return (
      <div className="flex justify-content-end">
        <Button
          rounded
          text
          label="Tarih Ayarla"
          icon="pi pi-calendar"
          outlined
          disabled={!selectedMatch || selectedMatch.status === Status.COMPLETED}
          size='small'
          onClick={(e) => {
            dateOP.current?.toggle(e);
          }}
        />
        <Button
          rounded
          text
          size='small'
          label="Maç Skoru Gir"
          disabled={!selectedMatch}
          icon="pi pi-pencil"
          outlined
          onClick={() => handleMatchScore()}
        />
        <Button
          rounded
          text
          size='small'
          label='Onayla'
          icon="pi pi-check"
          disabled={!selectedMatch || selectedMatch.status !== Status.COMPLETED}
          outlined
          onClick={() => approveHandler()}
        />
      </div>)
  }

  return (
    <>
      <Toast ref={toast} />
      <LeagueCard id={id!}></LeagueCard>
      <Card title="Fikstür">
        <DataTable
          value={matches}
          dataKey="id"
          loading={loading}
          emptyMessage="Fikstür bulunamadı"
          tableStyle={{ minWidth: '50rem' }}
          selectionMode="single"
          selection={selectedMatch!}
          onSelectionChange={(e) => setSelectedMatch(e.value)}
          key="id"
          header={header}
        >
          <Column
            selectionMode="single"
            headerStyle={{ width: "3rem" }}
          ></Column>
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
              rowData.status === Status.SCORE_APPROVED ||
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
      <MatchScoreSidebar visible={updateScoreVisible} matchId={selectedMatch?.id} onHide={() => setUpdateScoreVisible(false)} onSuccess={() => loadFixture(id!)} />

    </>
  );
}
