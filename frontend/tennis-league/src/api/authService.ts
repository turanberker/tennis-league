import axiosClient from './axiosClient';

/* ============================= */
/*            TYPES              */
/* ============================= */

export interface LoginRequest {
  email: string;
  password: string;
}

export interface CurrentUser {
  userID: number;
  name: string;
  surname: string;
  role: string;
}

export interface AuthResponse {
  token: string;
  currentUser: CurrentUser;
}

export interface RegisterRequest {
  email: string;
  name: string;
  surname: string;
  password: string;
}

/* ============================= */
/*        SERVICE METHODS        */
/* ============================= */

export const login = async (payload: LoginRequest): Promise<AuthResponse> => {
  return axiosClient.post<AuthResponse>('/auth/login', payload);
};

export const register = async (
  payload: RegisterRequest,
): Promise<AuthResponse> => {
  return axiosClient.post<AuthResponse>('/auth/register', payload);
};
