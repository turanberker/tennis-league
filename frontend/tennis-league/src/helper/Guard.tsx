import React from 'react';
import { Role, User } from '../model/user.model';
import { AuthUser, useAuth } from '../context/AuthContext';


interface GuardProps {
    children: React.ReactNode;
    allowedRoles: Role[];
    condition?: (user: AuthUser) => boolean; // Opsiyonel: Özel fonksiyon (true dönerse gösterir)
    fallback?: React.ReactNode; // Yetkisi yoksa yerine ne görünsün? (Opsiyonel)
}

export default function Guard({ children, allowedRoles, condition, fallback = null }: GuardProps) {
    const { user, isAuthenticated } = useAuth();

    // 1. Giriş yapmamışsa zaten gösterme
    if (!isAuthenticated || !user) {
        return <>{fallback}</>;
    }

    // 2. Rol kontrolü (Eğer allowedRoles dizisi verilmişse)
    const hasRole = allowedRoles
        ? allowedRoles.includes(user.role as Role)
        : true; // Belirtilmemişse rol engel teşkil etmez

    // 3. Custom fonksiyon kontrolü (Eğer condition verilmişse)
    const satisfiesCondition = condition
        ? condition(user)
        : true; // Belirtilmemişse fonksiyon engel teşkil etmez

    // İki şart da sağlanıyorsa çocukları göster
    if (hasRole && satisfiesCondition) {
        return <>{children}</>;
    }

    return <>{fallback}</>;

}