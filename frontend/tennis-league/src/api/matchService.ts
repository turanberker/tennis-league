import { MatchScore, MatchScoreResponse } from '../model/match.model';
import axiosClient from './axiosClient';

export const updateDate = async (id: string, data: { 'match-date': Date }) => {
  return await axiosClient.put(`match/${id}/update-date`, null, {
    params: data,
  });
};

export const updateMatchScore = async (
  id: string,
  data: MatchScore,
): Promise<MatchScoreResponse>  => {
  return await axiosClient.put<MatchScoreResponse>(`match/${id}/score`, data);
};
