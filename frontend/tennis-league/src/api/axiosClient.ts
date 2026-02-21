import axios, { AxiosError, AxiosRequestConfig } from 'axios';
import { showGlobalError } from './toastService';
import { ApiError } from '../model/apiResponse.model';

const instance = axios.create({
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
    const apiResponse = response.data;

    if (!apiResponse.success) {
      const message =
        apiResponse.error || apiResponse.errorDetail || 'Bilinmeyen hata';

      showGlobalError(message); // 🔥 burada toast

      throw new ApiError(message, response.status);
    }

    return apiResponse.data;
  },
  (error: AxiosError<any>) => {
    if (error.response) {
      const status = error.response.status;
      const message =
        error.response.data?.error ||
        error.response.data?.message ||
        'Sunucu hatası';

      if (status === 401) {
        localStorage.removeItem('user');
        window.location.href = '/';
      }

      showGlobalError(message); // 🔥 burada toast

      throw new ApiError(message, status);
    }

    showGlobalError('Network hatası');
    throw new ApiError('Network hatası');
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
