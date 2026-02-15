import axiosClient from './axiosClient';

export const getLeagues = async (params?: { name?: string }) => {
  return await axiosClient.get('leagues/list', { params });
};

export const getLeague = async (id: string) => {
  return await axiosClient.get(`leagues/${id}`);
};

export const saveLeague = async (data: { name: string }) => {
  return await axiosClient.post('leagues', data);
};
