import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Card } from 'primereact/card';
import { getFixture } from '../api/leagueService';
import {
  LeagueFixtureMatchResponse,
  MatchStatusLabels,
  Status,
} from '../model/match.model';
import { set } from 'react-hook-form';
import { DataTable } from 'primereact/datatable';
import { Column } from 'primereact/column';
import { Button } from 'primereact/button';

export default function Fixtures() {
  const { id } = useParams();
  const [loading, setLoading] = useState<boolean>(false);
  const [matches, setMatches] = useState<LeagueFixtureMatchResponse[]>([]);
  useEffect(() => {
    loadFixture(id!);
  }, [id]);

  const loadFixture = async (id: string) => {
    setLoading(true);
    const data = await getFixture(id);
    setMatches(data);
    setLoading(false);
  };

  const handleMatchDate = (match: LeagueFixtureMatchResponse) => {
    // Maç tarihini ayarlama işlemi burada yapılacak
    console.log('Tarih ayarla:', match);
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
              onClick={() => handleMatchDate(match)}
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
            onClick={() => handleMatchDate(match)}
          />
        );
    }
  };

  return (
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
        <Column field="matchDate" header="Maç Tarihi"></Column>
        <Column header="İşlem" body={(rowData) => getButtons(rowData)}></Column>
      </DataTable>
    </Card>
  );
}
