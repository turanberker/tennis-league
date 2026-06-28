import { Role } from '../../model/user.model';
import {userClient} from "../axiosClient";


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

export const login = async (payload: LoginRequest): Promise<AuthResponse> => {
  return userClient.post<AuthResponse>(`/auth/login`, payload);
};

export const logout = async (): Promise<void> => {
  return userClient.post<void>(`/auth/logout`);
};


export const register = async (
  payload: RegisterRequest,
): Promise<AuthResponse> => {
  return userClient.post<AuthResponse>(`/auth/register`, payload);
};
