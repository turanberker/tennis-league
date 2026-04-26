import { Role } from '../model/user.model';
import axiosClient from './axiosClient';

/* ============================= */
/*            TYPES              */
/* ============================= */

export interface LoginRequest {
  email: String;
  password: String;
}

export interface CurrentUser {
  userID: String;
  name: String;
  surname: String;
  role: Role;
}

export interface AuthResponse {
  currentUser: CurrentUser;
}

export interface RegisterRequest {
  email: String;
  name: String;
  surname: String;
  password: String;
}

/* ============================= */
/*        SERVICE METHODS        */
/* ============================= */

export const login = async (payload: LoginRequest): Promise<AuthResponse> => {
  return axiosClient.post<AuthResponse>('/auth/login', payload);
};

export const logout = async (): Promise<void> => {
  return axiosClient.post<void>('/auth/logout');
};


export const register = async (
  payload: RegisterRequest,
): Promise<AuthResponse> => {
  return axiosClient.post<AuthResponse>('/auth/register', payload);
};
