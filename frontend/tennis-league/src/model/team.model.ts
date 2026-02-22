export interface TeamResponse {
  id: string;
  name: string;
}

export interface CreateTeamRequest {
  name: string;
  playerIds: string[];
}

export interface TeamRef{
  id: string
  name:string
}