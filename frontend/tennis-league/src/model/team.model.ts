export interface LeagueTeamResponse {
  id: string;
  name: string;
  power: number
}

export interface CreateTeamRequest {
  name: string;
  playerIds: string[];
}

export interface TeamRef {
  id: string
  name: string
}