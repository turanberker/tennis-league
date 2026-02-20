import { MatchScore } from '../model/match.model';
import axiosClient from './axiosClient';

export const updateDate = async (id: string, data: { 'match-date': Date }) => {
  return await axiosClient.put(`match/${id}/update-date`, null, {
    params: data,
  });
};

export const updateMatchScore = async (id: string, data: MatchScore) => {
  return await axiosClient.put(`match/${id}/score`, data);
};
