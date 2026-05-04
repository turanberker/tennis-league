import { MatchScore, MatchScoreResponse, MatchSetScoreResponse } from '../model/match.model';
import axiosClient from './axiosClient';

export const updateDate = async (id: string, data: { 'match-date': Date }): Promise<Date> => {
  return await axiosClient.put(`match/${id}/update-date`, null, {
    params: data,
  });
};

export const updateFriendlyMatchScore = async (
  id: string,
  data: MatchScore,
): Promise<MatchScoreResponse> => {
  return await axiosClient.put<MatchScoreResponse>(`match/${id}/score`, data);
};

export const getMatchInfo = async (id: string): Promise<MatchSetScoreResponse> => {
  return await axiosClient.get<MatchSetScoreResponse>(`match/${id}/match-info`);
};

export const approve = async (id: string): Promise<MatchScore> => {
  return await axiosClient.put<MatchScore>(`match/${id}/approve`);
};
