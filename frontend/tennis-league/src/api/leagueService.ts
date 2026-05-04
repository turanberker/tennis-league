import { CreateTeamRequest, LeagueTeamResponse } from "../model/team.model";
import { LeagueFixtureMatchResponse, MatchScore, MatchScoreResponse } from "../model/match.model";
import axiosClient from "./axiosClient";
import { ScoreBoardResponse } from "../model/standing.model";
import { League, LeagueListResponse, PersistLeagueRequest } from "../model/league.model";

export const getLeagues = async (params?: { name?: string }): Promise<LeagueListResponse[]> => {
  return await axiosClient.get("leagues/list", { params });
};

export const getLeagueById = async (id: string): Promise<League> => {
  return await axiosClient.get(`leagues/${id}`);
};

export const saveLeague = async (
  data: PersistLeagueRequest,
): Promise<string> => {
  return await axiosClient.post("leagues", data);
};

export const getTeams = async (leagueId: string): Promise<LeagueTeamResponse[]> => {
  const res = await axiosClient.get<LeagueTeamResponse[]>(
    `leagues/${leagueId}/teams`,
  );
  return res;
};

export const createTeam = async (
  leagueId: string,
  team: CreateTeamRequest,
): Promise<{ teamId: string, totalAttendanceCount: number }> => {
  return await axiosClient.post<{ teamId: string, totalAttendanceCount: number }>(`leagues/${leagueId}/teams`, team);
};

export const createFixture = async (leagueId: string) => {
  return await axiosClient.post(`leagues/${leagueId}/create-fixture`);
};

export const getFixture = async (
  leagueId: string,
  param?: { teamId: string }
): Promise<LeagueFixtureMatchResponse[]> => {
  return await axiosClient.get<LeagueFixtureMatchResponse[]>(
    `leagues/${leagueId}/fixture`, { params: param }
  );
};

export const getStandings = async (
  leagueId: string,
): Promise<ScoreBoardResponse[]> => {
  return await axiosClient.get<ScoreBoardResponse[]>(
    `leagues/${leagueId}/standings`,
  );
};

export const assignCoordinator = async (leagueId: string, data: { userId: string }): Promise<Boolean> => {
  return await axiosClient.post<Boolean>(`leagues/${leagueId}/coordinator`, null, { params: data })
}

export const updateMatchDate = async (leagueId: string, matchId: string, data: { 'match-date': Date }): Promise<Date> => {
  return await axiosClient.put(`leagues/${leagueId}/match/${matchId}/update-date`, null, {
    params: data,
  });
};

export const approveMatchResult = async (leagueId: string, matchId: string): Promise<MatchScore> => {
  return await axiosClient.put<MatchScore>(`leagues/${leagueId}/match/${matchId}/approve`);
};

export const updateLeagueMatchScore = async (leagueId: string, matchId: string, data: MatchScore): Promise<MatchScoreResponse> => {
  return await axiosClient.put<MatchScoreResponse>(`leagues/${leagueId}/match/${matchId}/update-score`, data);
};
