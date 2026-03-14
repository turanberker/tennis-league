import axios, { AxiosError, AxiosRequestConfig } from 'axios';
import { showGlobalError } from './toastService';
import { ApiError, ApiResponse } from '../model/apiResponse.model';

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
    // DURUM 1: HTTP 200-299 arası bir kod döndü
    const apiResponse = response.data as ApiResponse<any>;

    // Backend 200 dönse bile kendi içinde success: false göndermiş olabilir
    if (!apiResponse.success) {
      const errorDetail = apiResponse.error || {
        code: 'UNKNOWN',
        message: 'İşlem başarısız',
      };
      showGlobalError(errorDetail.message);
      return Promise.reject(
        new ApiError(errorDetail.message, response.status, errorDetail.code),
      );
    }

    return apiResponse.data; // Componente sadece 'T' tipindeki data gider
  },
  (error: AxiosError<any>) => {
    // DURUM 2: HTTP 400, 500 vb. bir hata kodu döndü
    let message = 'Bir hata oluştu';
    let status = error.response?.status;
    let code = 'NETWORK_ERROR';

    if (error.response && error.response.data) {
      // Backend 400 dönerken senin ApiResponse formatını gönderiyor:
      // { success: false, error: { code: '...', message: '...' } }
      const apiResponse = error.response.data;

      if (apiResponse.error) {
        message = apiResponse.error.message;
        code = apiResponse.error.code;
      } else if ((apiResponse as any).message) {
        // Eğer backend bazen standart Error objesi dönerse (fallback)
        message = (apiResponse as any).message;
      }
    } else if (error.request) {
      message =
        'Sunucuya ulaşılamıyor. Lütfen internet bağlantınızı kontrol edin.';
    }

    // 🔥 Toast gösterimi
    showGlobalError(message);

    // Hataları her zaman Promise.reject ile fırlatmak asenkron akışı korur
    return Promise.reject(new ApiError(message, status, code));
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
