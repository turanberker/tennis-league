import React, { useRef } from 'react';
import { Toast } from 'primereact/toast';
import { useParams } from 'react-router-dom';
import { LeagueCard } from '../../components/LeagueCard';
import { useLeague } from '../../hooks/useLeague';
import { ProgressSpinner } from 'primereact/progressspinner';
import { LEAGUE_FORMAT } from '../../model/league.model';
import { LeaguePlayers } from '../../components/LeaguePlayers';
import { LeagueTeams } from '../../components/LeagueTeams';


const Teams: React.FC = () => {
  const { id } = useParams<{ id: string }>();

  const { data: league, isLoading } = useLeague(id);


  const toast = useRef<Toast>(null);



  if (isLoading) return <ProgressSpinner />;

  return (
    <>
      <Toast ref={toast} />

      <LeagueCard id={id!}></LeagueCard>

      {league?.format === LEAGUE_FORMAT.Double && (
        <LeagueTeams leagueId={id!} />
      )}

      {league?.format === LEAGUE_FORMAT.Single && (
        <LeaguePlayers leagueId={id!} />
      )}
    </>
  );
};

export default Teams;
