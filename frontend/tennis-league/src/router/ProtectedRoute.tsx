import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { Role } from '../model/user.model';

interface ProtectedRouteProps {
  children: React.ReactNode;
  allowedRoles?: Role[]; // İzin verilen roller listesi// children tipini belirtiyoruz

}
export default function ProtectedRoute({
  children,
  allowedRoles,
}: ProtectedRouteProps) {
  const { user, isAuthenticated, isLoading } = useAuth();
  // Veri henüz localStorage'dan okunuyorsa hiçbir şey yapma veya Spinner dön
  if (isLoading) {
    <div className="flex justify-center items-center h-screen">
      {/* Replace with your actual Spinner component */}
      <span>Yükleniyor...</span>
    </div>
  }
  // 1. Giriş yapmamışsa ana sayfaya at
  if (!isAuthenticated) {
    // We send the current path to the login page so we can redirect back later
    // return <Navigate to="/" state={{ from: location }} replace />;

    return <Navigate to="/" replace />;
  }
  const userRole = user?.role as Role;
  if (allowedRoles && !allowedRoles.includes(userRole)) {
    return <Navigate to="/unauthorized" replace />;
  }

  return <>{children}</>;
}
