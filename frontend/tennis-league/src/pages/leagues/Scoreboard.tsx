import React, { useEffect, useState } from 'react';
import { Card } from 'primereact/card';
import { useParams } from 'react-router-dom';
import { ScoreBoardResponse } from '../../model/standing.model';
import { getStandings } from '../../api/leagueService';
import { DataTable } from 'primereact/datatable';
import { Column } from 'primereact/column';

export default function Scoreboard() {
  const { id } = useParams();
  const [loading, setLoading] = useState<boolean>(false);
  const [board, setBoard] = useState<ScoreBoardResponse[]>([]);
  useEffect(() => {
    loadStandings();
  }, [id]);

  const loadStandings = async () => {
    setLoading(true);
    const response = await getStandings(id!);
    setBoard(response);
    setLoading(false);
  };
  return (
    <Card title="Puan Durumu">
      <DataTable
        value={board}
        loading={loading}
        emptyMessage="Puan Durumu bulunamadı"
        tableStyle={{ minWidth: '50rem' }}
        key="id"
      >
        <Column field="order" header="Sıra" />
        <Column field="name" header="Takım Adı" />
        <Column field="played" header="Oynadığı Maç" />
        <Column field="score" header="Toplam Puan" />

        <Column field="won" header="Kazandığı Maç" />
        <Column field="lost" header="Kaybettiği Maç" />
        <Column field="wonSets" header="Kazandığı Set" />
        <Column field="lostSets" header="Kaybettiği Set" />

        <Column field="wonGames" header="Kazandığı Oyun" />
        <Column field="lostGames" header="Kaybettiği Oyun" />
      </DataTable>
    </Card>
  );
}
