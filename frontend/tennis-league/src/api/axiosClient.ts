import axios, { AxiosRequestConfig } from 'axios';

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

    if (apiResponse?.success) {
      return apiResponse.data;
    }

    return Promise.reject(
      new Error(apiResponse?.errorDetail || 'Bilinmeyen hata')
    );
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('user');
      window.location.href = '/';
    }

    return Promise.reject(error);
  }
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
