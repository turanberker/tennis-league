import { AxiosResponse } from 'axios';
import { CreateTeamRequest, TeamResponse } from '../model/team.model';
import { LeagueFixtureMatchResponse } from '../model/match.model';
import axiosClient from './axiosClient';
import { ScoreBoardResponse } from '../model/standing.model';

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

export const createFixture = async (leagueId: string) => {
  return await axiosClient.post(`leagues/${leagueId}/create-fixture`);
};

export const getFixture = async (
  leagueId: string,
): Promise<LeagueFixtureMatchResponse[]> => {
  return await axiosClient.get<LeagueFixtureMatchResponse[]>(
    `leagues/${leagueId}/fixture`,
  );
};

export const getStandings = async (
  leagueId: string,
): Promise<ScoreBoardResponse[]> => {
  return await axiosClient.get<ScoreBoardResponse[]>(
    `leagues/${leagueId}/standings`,
  );
};
