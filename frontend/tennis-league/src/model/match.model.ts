import { TeamRef } from './team.model';

export enum Status {
  PERDING = 'PENDING',
  COMPLETED = 'COMPLETED',
  APPROVED = 'APPROVED',
  CANCELLED = 'CANCELLED',
}

export const MatchStatusLabels: Record<Status, string> = {
  [Status.PERDING]: 'Beklemede',
  [Status.COMPLETED]: 'Tamamlandı',
  [Status.APPROVED]: 'Onaylandı',
  [Status.CANCELLED]: 'İptal Edildi',
};

export interface LeagueFixtureMatchResponse {
  id: string;
  team1: TeamRefResponse;
  team2: TeamRefResponse;
  status: Status;
  matchDate: Date;
}

export interface TeamRefResponse extends TeamRef {
  score?: number;
  winner?: boolean;
}

export interface MatchScore {
  set1: SetScore;
  set2: SetScore;
  superTie: SetScore | null;
}

interface SetScore {
  team1Score: number;
  team2Score: number;
}

export interface MatchScoreResponse {
  team1: null;
  team2: null;
}
