import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

interface ProtectedRouteProps {
  children: React.ReactNode;
  allowedRoles?: string[]; // İzin verilen roller listesi// children tipini belirtiyoruz
}
export default function ProtectedRoute({
  children,
  allowedRoles,
}: ProtectedRouteProps) {
  const { user, isAuthenticated, isLoading } = useAuth();
  // Veri henüz localStorage'dan okunuyorsa hiçbir şey yapma veya Spinner dön
  if (isLoading) {
    return <div>Yükleniyor...</div>; // Veya PrimeReact ProgressSpinner
  }
  // 1. Giriş yapmamışsa ana sayfaya at
  if (!isAuthenticated) {
    return <Navigate to="/" replace />;
  }

  // 2. Rol kontrolü: Eğer roller belirtilmişse ve kullanıcının rolü listede yoksa
  if (allowedRoles && !allowedRoles.includes(user?.role || '')) {
    return <Navigate to="/unauthorized" replace />; // Yetkisiz erişim sayfası
  }

  return <>{children}</>;
}
