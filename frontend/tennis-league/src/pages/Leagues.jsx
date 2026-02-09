import React, { useEffect, useState } from 'react';
import { Card } from 'primereact/card';
import { DataScroller } from 'primereact/datascroller';
import { Button } from 'primereact/button';
import { getLeagues } from "../api/leagueService";


export default function Leagues() {
  
    const [leagues, setLeagues] = useState([]);
const [error, setError] = useState(null);
  
    useEffect(() => {
    getLeagues()
      .then(setLeagues)
      .catch((err) => {
        console.error("League fetch error:", err);
        setError("Ligler yüklenemedi");
      });
  }, []);
    
      const handleStandings = (league) => {
    console.log("Puan durumu:", league);
    // TODO: standings sayfasına yönlendirme
  };

  const handleFixtures = (league) => {
    console.log("Fikstür:", league);
    // TODO: fixtures sayfasına yönlendirme
  };

  return (
    <Card title="Ligler">
      {error && <p style={{ color: "red" }}>{error}</p>}

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
  );
}