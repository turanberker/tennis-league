import axios, { AxiosError, AxiosInstance, AxiosRequestConfig } from 'axios';
import { sendLogoutEvent, showGlobalError } from './toastService';
import { ApiResponse } from '../model/apiResponse.model';

// 1. INSTANCE'LARI OLUŞTURUYORUZ
export const mainInstance = axios.create({
  baseURL: process.env.REACT_APP_API_URL,
  withCredentials: true,
});

export const userInstance = axios.create({
  baseURL: process.env.REACT_APP_USER_URL,
  withCredentials: true,
});

// Refresh token sürecini yöneten ortak state'ler
let isRefreshing = false;
let failedQueue: any[] = [];

const processQueue = (error: any, token: string | null = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
  failedQueue = [];
};

// 2. INTERCEPTOR'LARI TEK BİR FONKSİYONDA TOPLUYORUZ
const attachInterceptors = (targetInstance: AxiosInstance) => {
  targetInstance.interceptors.response.use(
      (response) => {
        const apiResponse = response.data as ApiResponse<any>;
        if (apiResponse.success) {
          return apiResponse.data;
        }

        showGlobalError(apiResponse.error?.message || 'İşlem başarısız');
        return null;
      },
      async (error: AxiosError<any>) => {
        const originalRequest = error.config as any;
        const message = error.response?.data?.error?.message || 'Sistem hatası';
        const status = error.response?.status;

        // 401 ve Token Yenileme Süreci
        if (status === 401 && originalRequest && !originalRequest._retry) {
          if (isRefreshing) {
            return new Promise((resolve, reject) => {
              failedQueue.push({ resolve, reject });
            })
                .then(() => targetInstance(originalRequest)) // Hangi instance hata aldıysa onunla tekrarla
                .catch(() => null);
          }

          originalRequest._retry = true;
          isRefreshing = true;

          try {
            // Backend'e refresh isteği at (Hangi instance üzerinden gittiği önemli değil, cookie tabanlı)
            await axios.post(`${process.env.REACT_APP_USER_URL}/auth/refresh`, {}, { withCredentials: true });

            isRefreshing = false;
            processQueue(null);

            return targetInstance(originalRequest); // Asıl isteği ilgili instance ile tekrarla
          } catch (refreshError) {
            isRefreshing = false;
            processQueue(refreshError);

            sendLogoutEvent("Oturum süreniz doldu, lütfen tekrar giriş yapın.");
            return null;
          }
        }

        // Diğer Tüm Hatalar
        if (status !== 401) {
          showGlobalError(message);
        }

        return null;
      },
  );
};

// 3. AYNI INTERCEPTOR'I İKİ INSTANCE'A DA BAĞLIYORUZ
attachInterceptors(mainInstance);
attachInterceptors(userInstance);

// 4. TYPED WRAPPER FONKSİYONU (Her iki instance için de dinamik sarmalayıcı üreten yardımcı fonksiyon)
const createClientWrapper = (targetInstance: AxiosInstance) => ({
  get: <T = any>(url: string, config?: AxiosRequestConfig) =>
      targetInstance.get<any, T>(url, config),

  post: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) =>
      targetInstance.post<any, T>(url, data, config),

  put: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) =>
      targetInstance.put<any, T>(url, data, config),

  patch: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) =>
      targetInstance.patch<any, T>(url, data, config),

  delete: <T = any>(url: string, config?: AxiosRequestConfig) =>
      targetInstance.delete<any, T>(url, config),
});

// 5. DIŞA AKTARMALAR
export const mainClient = createClientWrapper(mainInstance);
export const userClient = createClientWrapper(userInstance);
