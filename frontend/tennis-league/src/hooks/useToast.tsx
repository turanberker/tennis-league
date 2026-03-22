// context/ToastContext.tsx
import React, { createContext, useContext, useRef } from 'react';
import { Toast } from 'primereact/toast';
import { ToastMessage } from 'primereact/toast';

type ToastContextType = {
    show: (message: ToastMessage) => void;
};

const ToastContext = createContext<ToastContextType | null>(null);

export const ToastProvider = ({ children }: { children: React.ReactNode }) => {
    const toastRef = useRef<Toast>(null);

    const show = (message: ToastMessage) => {
        toastRef.current?.show(message);
    };

    return (
        <ToastContext.Provider value={{ show }}>
            {children}
            <Toast ref={toastRef} position="top-right" />
        </ToastContext.Provider>
    );
};

export const useToast = () => {
    const ctx = useContext(ToastContext);
    if (!ctx) throw new Error('ToastProvider yok');
    return ctx;
};