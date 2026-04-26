import React, { useEffect } from 'react';
import { Role } from '../model/user.model';
import { AuthUser, useAuth } from '../context/AuthContext';


interface GuardProps {
    children: React.ReactNode;
    allowedRoles?: Role[] | Role;
    condition?: (user: AuthUser) => boolean; // Opsiyonel: Özel fonksiyon (true dönerse gösterir)
    fallback?: React.ReactNode; // Yetkisi yoksa yerine ne görünsün? (Opsiyonel)
    onFail?: () => void;// Yetki olmadığında çalışacak metod
}

export default function Guard({ children, allowedRoles, condition, fallback = null, onFail }: GuardProps) {
    const { user, isAuthenticated } = useAuth();
    // Yetki kontrol mantığı aynı kalsın...
    // 1. Tüm yetki mantığını tek bir boolean değişkende topla (Render sırasında hiçbir fonksiyon çağırma!)
    const hasRole = !allowedRoles || (
        Array.isArray(allowedRoles)
            ? (allowedRoles.length === 0 || allowedRoles.includes(user?.role as Role))
            : user?.role === allowedRoles
    );

    const satisfiesCondition = !condition || (user ? condition(user) : false);

    // Nihai izin durumu
    const isAllowed = !!(isAuthenticated && user && hasRole && satisfiesCondition);

    // 2. YAN ETKİ (Side Effect): State güncelleme veya dış metod tetikleme sadece BURADA yapılır.
    useEffect(() => {
        if (!isAllowed && onFail) {
            onFail(); // Artık güvenli, çizim bittikten sonra çalışır.
        }
    }, [isAllowed, onFail]);

    // 3. RENDER: Sadece ne görüneceğine karar ver
    if (isAllowed) {
        return <>{children}</>;
    }

    return <>{fallback}</>;

}