import React from 'react';
import { Routes, Route, Navigate, useNavigate } from 'react-router-dom';
import { Button } from 'primereact/button';
import Dashboard from '../pages/Dashboard';
import Leagues from '../pages/Leagues';
import Players from '../pages/Players';
import Matches from '../pages/Matches';
import ProtectedRoute from './ProtectedRoute';
import Standings from '../pages/Standings';
import Fixtures from '../pages/Fixtures';
import Teams from '../pages/Teams';

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

      <Route
        path="/leagues"
        element={
          <ProtectedRoute>
            <Leagues />
          </ProtectedRoute>
        }
      />
      <Route
        path="/leagues/:id/standings"
        element={
          <ProtectedRoute>
            <Standings />
          </ProtectedRoute>
        }
      />
      <Route
        path="/leagues/:id/teams"
        element={
          <ProtectedRoute>
            <Teams />
          </ProtectedRoute>
        }
      />
      <Route
        path="/leagues/:id/fixtures"
        element={
          <ProtectedRoute>
            <Fixtures />
          </ProtectedRoute>
        }
      />
      <Route
        path="/players"
        element={
          <ProtectedRoute>
            <Players />
          </ProtectedRoute>
        }
      />

      <Route
        path="/matches"
        element={
          <ProtectedRoute>
            <Matches />
          </ProtectedRoute>
        }
      />

      <Route path="*" element={<Navigate to="/" />} />
    </Routes>
  );
}
