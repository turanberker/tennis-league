import axios, { AxiosError, AxiosRequestConfig } from 'axios';
import { sendLogoutEvent, showGlobalError } from './toastService';
import { ApiResponse } from '../model/apiResponse.model';

export const instance = axios.create({
  baseURL: process.env.REACT_APP_API_URL,
  withCredentials: true,
});

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

// RESPONSE INTERCEPTOR
instance.interceptors.response.use(
  (response) => {
    const apiResponse = response.data as ApiResponse<any>;
    if (apiResponse.success) {
      return apiResponse.data;
    }

    showGlobalError(apiResponse.error?.message || 'İşlem başarısız');

    return null;
  },
  async (error: AxiosError<any>) => {

    // 1. Değişkenleri bir kez ve en üstte tanımlıyoruz
    const originalRequest = error.config as any;
    const message = error.response?.data?.error?.message || 'Sistem hatası';
    const status = error.response?.status;

    // 2. SENARYO: Unauthorized (401) ve Token Yenileme Süreci
    if (status === 401 && originalRequest && !originalRequest._retry) {

      // Eğer zaten bir yenileme isteği yoldaysa, bu isteği beklemeye al
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then(() => instance(originalRequest))
          .catch(() => null);
      }

      // İlk 401 hatası: Yenileme bayrağını kaldır ve süreci başlat
      originalRequest._retry = true;
      isRefreshing = true;

      try {
        // Backend'e refresh isteği at (Cookie otomatik gider)
        await axios.post(`${process.env.REACT_APP_USER_URL}/auth/refresh`, {}, { withCredentials: true });

        isRefreshing = false;
        processQueue(null); // Bekleyen diğer 401'li istekleri serbest bırak

        return instance(originalRequest); // Asıl isteği yeni token ile tekrarla
      } catch (refreshError) {
        isRefreshing = false;
        processQueue(refreshError); // Kuyruktakileri de iptal et

        // Refresh başarısızsa (örneğin 7 gün dolmuşsa) zorunlu logout
        sendLogoutEvent("Oturum süreniz doldu, lütfen tekrar giriş yapın.");
        return null;
      }
    }

    // 3. SENARYO: Diğer Tüm Hatalar (500, 403, 404 vb.)
    // 401 durumunda yukarıda refresh denediğimiz için hemen hata göstermiyoruz.
    // Ama status 401 değilse, kullanıcıya ne hata varsa (toast ile) gösteriyoruz.
    if (status !== 401) {
      showGlobalError(message);
    }

    // En son her durumda null dönerek .tsx tarafında catch yazma zorunluluğunu kaldırıyoruz
    return null;
  },
);

/* ============================= */
/*   TYPED WRAPPER METHODS       */
/* ============================= */

const axiosClient = {
  get: <T = any>(url: string, config?: AxiosRequestConfig) =>
    instance.get<any, T>(url, config),

  post: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) =>
    instance.post<any, T>(url, data, config),

  put: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) =>
    instance.put<any, T>(url, data, config),

  patch: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) =>
    instance.patch<any, T>(url, data, config),

  delete: <T = any>(url: string, config?: AxiosRequestConfig) =>
    instance.delete<any, T>(url, config),
};

export default axiosClient;
