import axiosClient from './axiosClient'; // senin axiosClient setup

export const login = async ({ email, password }) => {
  const res = await axiosClient.post('/auth/login', { email, password });
  return res;
};

export const register = async (data) => {
  const res = await axiosClient.post(`/auth/register`, data);
  return res;
};
