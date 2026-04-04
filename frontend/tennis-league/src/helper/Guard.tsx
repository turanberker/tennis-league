import React from 'react';
import { Role } from '../model/user.model';
import { AuthUser, useAuth } from '../context/AuthContext';


interface GuardProps {
    children: React.ReactNode;
    allowedRoles: Role[] | Role;
    condition?: (user: AuthUser) => boolean; // Opsiyonel: Özel fonksiyon (true dönerse gösterir)
    fallback?: React.ReactNode; // Yetkisi yoksa yerine ne görünsün? (Opsiyonel)
}

export default function Guard({ children, allowedRoles, condition, fallback = null }: GuardProps) {
    const { user, isAuthenticated } = useAuth();

    // 1. Giriş yapmamışsa zaten gösterme
    if (!isAuthenticated || !user) {
        return <>{fallback}</>;
    }

    // 2. Rol kontrolü (Refactor edilmiş kısım)
    let hasRole = true; // Varsayılan olarak izin ver (allowedRoles yoksa)

    if (allowedRoles) {
        if (Array.isArray(allowedRoles)) {
            // Eğer diziyse includes kullan
            hasRole = allowedRoles.includes(user.role as Role);
        } else {
            // Eğer tekil değerse direkt karşılaştır
            hasRole = user.role === allowedRoles;
        }
    }

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