import { Role } from '../../model/user.model';
import axiosClient from '../axiosClient';

/* ============================= */
/*            TYPES              */
/* ============================= */

export interface LoginRequest {
  email: string;
  password: string;
}

export interface CurrentUser {
  userID: string;
  name: string;
  surname: string;
  role: Role;
}

export interface AuthResponse {
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
const USER_API_URL = process.env.REACT_APP_USER_URL || 'http://localhost:8000';

export const login = async (payload: LoginRequest): Promise<AuthResponse> => {
  return axiosClient.post<AuthResponse>(`${USER_API_URL}/auth/login`, payload);
};

export const logout = async (): Promise<void> => {
  return axiosClient.post<void>(`${USER_API_URL}/auth/logout`);
};


export const register = async (
  payload: RegisterRequest,
): Promise<AuthResponse> => {
  return axiosClient.post<AuthResponse>(`${USER_API_URL}/auth/register`, payload);
};
