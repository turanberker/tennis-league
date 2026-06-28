import { useEffect, useRef, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Card } from 'primereact/card';
import { approveMatchResult, getFixture, getTeams, updateLeagueMatchScore, updateMatchDate } from '../../api/mail/leagueService';
import {
  LeagueFixtureMatchResponse,
  MatchScore,
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
import { LeagueTeamResponse } from '../../model/team.model';
import { Dropdown } from 'primereact/dropdown';
import Guard from '../../helper/Guard';
import { Role } from '../../model/user.model';
import TeamInfo from '../../components/DoubleTeamInfo';

export default function Fixtures() {
  const { id } = useParams();
  const calendarRef = useRef<Calendar>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [matches, setMatches] = useState<LeagueFixtureMatchResponse[]>([]);
  const [updateScoreVisible, setUpdateScoreVisible] = useState<boolean>(false);
  const [selectedMatch, setSelectedMatch] = useState<LeagueFixtureMatchResponse>();
  const dateOP = useRef<OverlayPanel>(null);
  const toast = useRef<Toast>(null);
  const [tempDate, setTempDate] = useState<Date | null>(null);

  const [teams, setTeams] = useState<LeagueTeamResponse[]>([]); // Takımlar için state
  const [selectedTeamId, setSelectedTeamId] = useState<string | null>(null); // Filtre için state

  const { data: league } = useLeague(id);

  useEffect(() => {

    const initializeFixture = async () => {
      const data = await getFixture(id!);
      setMatches(data);
      await loadTeams(id!);
      setLoading(false); // İşlem bitince yükleniyor durumunu kapat
    }

    if (id) {
      initializeFixture();
    }
  }, [id]);


  const loadTeams = async (leagueId: string) => {
    const data = await getTeams(leagueId);
    setTeams(data);
  };

  const loadFixture = async (id: string) => {
    setLoading(true);
    const data = await getFixture(id);


    if (selectedMatch) {

      let selectedMatchId = selectedMatch.id;
      setSelectedMatch(undefined)
      const selected = data.find(e => e.id === selectedMatchId)

      setSelectedMatch(selected)
    }

    setMatches(data);
    setLoading(false);
  };

  const handleMatchDate = async (forcedDate?: Date | null) => {
    calendarRef.current?.hide()

    const dateToSend = forcedDate !== undefined ? forcedDate : tempDate;

    if (selectedMatch && league) {
      const res = await updateMatchDate(league.id, selectedMatch.id, { 'match-date': dateToSend });

      if (res) {
        // Listeyi güncelle
        setMatches((prev) =>
          prev.map((m) =>
            m.id === selectedMatch.id ? { ...m, matchDate: dateToSend } : m,
          ),
        );

        // Seçili objeyi güncelle (Senkronizasyon için önemli)
        setSelectedMatch(prev => prev ? { ...prev, matchDate: dateToSend } : prev);

        toast.current?.show({
          severity: 'success',
          summary: 'İşlem Başarılı',
          detail: `Maç Tarihi Güncellendi`,
          life: 3000,
        });

        dateOP.current?.hide(); // İşlem bitince paneli kapat
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
      <TeamInfo
        teamName={team.name}
        teamId={team.id}
        nameLeftAddon={team.winner && (
          <i className="pi pi-check-circle" style={{ display: 'inline-flex', alignItems: 'center', color: '#22c55e', fontSize: '1rem' }}></i>
        )}
      />
    );
  };

  const handleSubmitScore = async (matchId: string, score: MatchScore): Promise<boolean> => {

    let res = await updateLeagueMatchScore(id!, matchId, score)

    return res != null ? true : false;
  }

  // --- Filtreleme Mantığı ---
  // Eğer bir takım seçiliyse, o takımın team1 veya team2 olduğu maçları getirir.
  const filteredMatches = selectedTeamId
    ? matches.filter(m => m.team1.id === selectedTeamId || m.team2.id === selectedTeamId)
    : matches;

  const header = () => {
    return (
      <div className="flex justify-content-between align-items-center">
        {/* Sol Tarafta Filtreleme Dropdown */}
        <div className="flex align-items-center gap-2">
          <span className="text-sm font-semibold">Takım Filtrele:</span>
          <Dropdown
            value={selectedTeamId}
            options={teams}
            optionLabel="name"
            optionValue="id"
            placeholder="Tüm Takımlar"
            showClear
            filter
            filterBy="name"
            className="p-inputtext-sm w-full md:w-15rem"
            onChange={(e) => setSelectedTeamId(e.value)}
          />
        </div>
        {/* Sağ Tarafta Aksiyon Butonları */}
        <div className="flex gap-2">
          <Guard allowedRoles={[Role.ADMIN, Role.COORDINATOR]}>
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
                setTempDate(selectedMatch?.matchDate ? new Date(selectedMatch.matchDate) : null);
              }}
            />
          </Guard>
          <Button
            rounded
            text
            size='small'
            label="Maç Skoru"
            disabled={!selectedMatch}
            icon="pi pi-pencil"
            outlined
            onClick={() => handleMatchScore()}
          />
          <Guard allowedRoles={[Role.ADMIN, Role.COORDINATOR]}>
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
          </Guard>
        </div>
      </div>

    )
  }



  return (
    <>
      <Toast ref={toast} />
      <LeagueCard id={id!}></LeagueCard>
      <Card title="Fikstür">
        <DataTable
          value={filteredMatches}
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
        <div className="flex flex-column gap-3">
          <FloatLabel>
            <Calendar
              appendTo="self"
              ref={calendarRef}
              showButtonBar
              showTime
              hourFormat="24"
              inputId="match_date_input"
              footerTemplate={() => (
                <div className="flex justify-content-end p-2">
                  <Button type="button" label="Tarihi Güncelle" className="p-button-sm" onClick={() => handleMatchDate()} />
                </div>
              )}
              onClearButtonClick={() => { setTempDate(null); handleMatchDate(null) }}
              value={tempDate}
              onChange={(e) => setTempDate(e.value as Date)}
            />
            <label htmlFor="match_date_input">Maç Tarihi Seçin</label>
          </FloatLabel>


        </div>
      </OverlayPanel>
      <MatchScoreSidebar visible={updateScoreVisible} matchId={selectedMatch?.id}
        onHide={() => setUpdateScoreVisible(false)} submitMatchScore={handleSubmitScore}
        onSuccess={() => loadFixture(id!)} />

    </>
  );
}
