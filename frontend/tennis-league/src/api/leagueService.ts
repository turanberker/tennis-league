import { AxiosResponse } from "axios";
import { CreateTeamRequest, LeagueTeamResponse } from "../model/team.model";
import { LeagueFixtureMatchResponse } from "../model/match.model";
import axiosClient from "./axiosClient";
import { ScoreBoardResponse } from "../model/standing.model";
import { League, LeagueListResponse, PersistLeagueRequest } from "../model/league.model";

export const getLeagues = async (params?: { name?: String }): Promise<LeagueListResponse[]> => {
  return await axiosClient.get("leagues/list", { params });
};

export const getLeagueById = async (id: String): Promise<League> => {
  return await axiosClient.get(`leagues/${id}`);
};

export const saveLeague = async (
  data: PersistLeagueRequest,
): Promise<String> => {
  return await axiosClient.post("leagues", data);
};

export const getTeams = async (leagueId: String): Promise<LeagueTeamResponse[]> => {
  const res = await axiosClient.get<LeagueTeamResponse[]>(
    `leagues/${leagueId}/teams`,
  );
  return res;
};

export const createTeam = async (
  leagueId: String,
  team: CreateTeamRequest,
): Promise<{ teamId: String, totalAttendanceCount: number }> => {
  return await axiosClient.post<{ teamId: String, totalAttendanceCount: number }>(`leagues/${leagueId}/teams`, team);
};

export const createFixture = async (leagueId: String) => {
  return await axiosClient.post(`leagues/${leagueId}/create-fixture`);
};

export const getFixture = async (
  leagueId: String,
): Promise<LeagueFixtureMatchResponse[]> => {
  return await axiosClient.get<LeagueFixtureMatchResponse[]>(
    `leagues/${leagueId}/fixture`,
  );
};

export const getStandings = async (
  leagueId: String,
): Promise<ScoreBoardResponse[]> => {
  return await axiosClient.get<ScoreBoardResponse[]>(
    `leagues/${leagueId}/standings`,
  );
};

export const assignCoordinator = async (leagueId: String, data: { userId: String }): Promise<Boolean> => {

  return await axiosClient.post<Boolean>(`leagues/${leagueId}/coordinator`, null, { params: data })

}
