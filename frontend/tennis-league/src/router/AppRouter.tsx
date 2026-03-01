import React from 'react';
import { Routes, Route, Navigate, useNavigate } from 'react-router-dom';
import { Button } from 'primereact/button';
import Dashboard from '../pages/Dashboard';
import Leagues from '../pages/Leagues';
import Players from '../pages/Players';
import Matches from '../pages/Matches';
import ProtectedRoute from './ProtectedRoute';
import Fixtures from '../pages/leagues/Fixtures';
import Teams from '../pages/Teams';
import NewPlayer from '../pages/player/NewPlayer';
import PlayerDetail from '../pages/player/PlayerDetail';
import Scoreboard from '../pages/leagues/Scoreboard';

export function SidebarLinks() {
  const navigate = useNavigate();

  return (
    <div className="flex flex-column gap-2">
      <Button
        label="Dashboard"
        text
        icon="pi pi-home"
        onClick={() => navigate('/')}
      />
      <Button
        label="Ligler"
        text
        icon="pi pi-sitemap"
        onClick={() => navigate('/leagues')}
      />
      <Button
        label="Oyuncular"
        text
        icon="pi pi-users"
        onClick={() => navigate('/players')}
      />
      <Button
        label="Maçlar"
        text
        icon="pi pi-calendar"
        onClick={() => navigate('/matches')}
      />
    </div>
  );
}

export function AppRoutes() {
  return (
    <Routes>
      <Route path="/" element={<Dashboard />} />

      <Route path="/leagues" element={<Leagues />} />
      <Route path="/leagues/:id/teams" element={<Teams />} />
      <Route path="/leagues/:id/fixtures" element={<Fixtures />} />
      <Route path="/leagues/:id/standings" element={<Scoreboard />} />
      <Route path="/players" element={<Players />} />

      <Route path="/players/:uuid" element={<PlayerDetail />} />

      <Route path="/matches" element={<Matches />} />

      <Route path="*" element={<Navigate to="/" />} />
    </Routes>
  );
}
