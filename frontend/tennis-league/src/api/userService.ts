import axiosClient from './axiosClient';

import { User } from '../model/user.model';

export const getUsers = async (): Promise<User[]> => {
  return axiosClient.get<User[]>('/user/list');
};
