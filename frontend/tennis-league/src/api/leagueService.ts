import { AxiosResponse } from 'axios';
import { CreateTeamRequest, TeamResponse } from '../model/team.model';
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

export const getTeams = async (leagueId: string): Promise<TeamResponse[]> => {
  const res = await axiosClient.get<TeamResponse[]>(
    `leagues/${leagueId}/teams`,
  );
  return res;
};

export const createTeam = async (
  leagueId: string,
  team: CreateTeamRequest,
): Promise<string> => {
  return await axiosClient.post<string>(`leagues/${leagueId}/teams`, team);
};
