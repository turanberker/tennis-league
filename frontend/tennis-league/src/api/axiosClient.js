import axios from 'axios';

const axiosClient = axios.create({
  baseURL: 'http://localhost:8500',
  headers: {
    'Content-Type': 'application/json',
  },
});

// İleride JWT eklemek için interceptor
axiosClient.interceptors.request.use((config) => {
  const user = localStorage.getItem('user');

  if (user) {
    const token = JSON.parse(user)?.token;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
  }

  return config;
});

// Response interceptor → backend response unwrap + global error
axiosClient.interceptors.response.use(
  (response) => {
    const apiResponse = response.data;

    // Beklenen format: { data, success, errorDetail }
    if (apiResponse?.success) {
      return apiResponse.data;
    }

    // success false ise hata fırlat
    const error = new Error(apiResponse?.errorDetail || 'Bilinmeyen hata');
    return Promise.reject(error);
  },
  (error) => {
    // 401 → otomatik logout yapılabilir (ileride AuthContext bağlanacak)
    if (error.response?.status === 401) {
      localStorage.removeItem('user');
      window.location.href = '/';
    }

    return Promise.reject(error);
  },
);

export default axiosClient;
