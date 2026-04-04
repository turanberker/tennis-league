import axios, { AxiosError, AxiosRequestConfig } from 'axios';
import { showGlobalError } from './toastService';
import { ApiResponse } from '../model/apiResponse.model';

export const instance = axios.create({
  baseURL: 'http://localhost:8500',
  headers: { 'Content-Type': 'application/json' },
});

// REQUEST INTERCEPTOR (token)
instance.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// RESPONSE INTERCEPTOR
instance.interceptors.response.use(
  (response) => {
    const apiResponse = response.data as ApiResponse<any>;
    if (apiResponse.success) {
      return apiResponse.data;
    }

    // Backend success: false döndü
    showGlobalError(apiResponse.error?.message || 'İşlem başarısız');

    // 🔥 KRİTİK: Reject fırlatmak yerine null dönüyoruz.
    // Böylece .tsx tarafında .catch yazmaya gerek kalmıyor.
    return null;
  },
  (error: AxiosError<any>) => {
    const message = error.response?.data?.error?.message || 'Sistem hatası';

    if (error.response?.data.code === "AUTH_102") {
      //burada logout u tetiklemek istiyorum. 
    }
    showGlobalError(message);

    // 🔥 HTTP hatalarında da (401, 500 vb.) sessiz kalıyoruz
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

  delete: <T = any>(url: string, config?: AxiosRequestConfig) =>
    instance.delete<any, T>(url, config),
};

export default axiosClient;
