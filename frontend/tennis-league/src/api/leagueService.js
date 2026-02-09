import axiosClient from './axiosClient';

export const getLeagues = async () => {
  return await axiosClient.get('leagues/list');
};
