import { MatchScore, MatchScoreResponse, MatchSetScoreResponse } from '../model/match.model';
import {mainClient} from "./axiosClient";


export const updateDate = async (id: string, data: { 'match-date': Date }): Promise<Date> => {
  return await mainClient.put(`match/${id}/update-date`, null, {
    params: data,
  });
};

export const updateFriendlyMatchScore = async (
  id: string,
  data: MatchScore,
): Promise<MatchScoreResponse> => {
  return await mainClient.put<MatchScoreResponse>(`match/${id}/score`, data);
};

export const getMatchInfo = async (id: string): Promise<MatchSetScoreResponse> => {
  return await mainClient.get<MatchSetScoreResponse>(`match/${id}/match-info`);
};

export const approve = async (id: string): Promise<MatchScore> => {
  return await mainClient.put<MatchScore>(`match/${id}/approve`);
};
