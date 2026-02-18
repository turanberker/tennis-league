import axiosClient from './axiosClient';

export const updateDate = async (id: string, data: { 'match-date': Date }) => {
  return await axiosClient.put(`match/${id}/update-date`, null, {
    params: data,
  });
};
