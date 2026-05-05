import { TeamRef } from './team.model';

export enum Status {
  PERDING = 'PENDING',
  COMPLETED = 'COMPLETED',
  SCORE_APPROVED = 'SCORE_APPROVED',
  CANCELLED = 'CANCELLED',
}

export const MatchStatusLabels: Record<Status, string> = {
  [Status.PERDING]: 'Beklemede',
  [Status.COMPLETED]: 'Tamamlandı',
  [Status.SCORE_APPROVED]: 'Skor Onaylandı',
  [Status.CANCELLED]: 'İptal Edildi',
};


export enum MatchType {
  SINGLE = 'SINGLE',
  DOUBLE = 'DOUBLE',
}

export const MatchTypeLabels: Record<MatchType, string> = {
  [MatchType.SINGLE]: 'Tekler',
  [MatchType.DOUBLE]: 'Çiftler',
};


export enum MatchSource {
  FRIENDLY = "FRIENDLY",
  TOURNAMENT = "TOURNAMENT",
  LEAGUE = "LEAGUE"
}

export const MatchSourceLabels: Record<MatchSource, string> = {
  [MatchSource.FRIENDLY]: 'Dostluk Maçı',
  [MatchSource.TOURNAMENT]: 'Turnuva',
  [MatchSource.LEAGUE]: 'Lig',
};

export interface LeagueFixtureMatchResponse {
  id: string;
  team1: TeamRefResponse;
  team2: TeamRefResponse;
  status: Status;
  matchDate?: Date | null;
}

export interface TeamRefResponse extends TeamRef {
  score?: number;
  winner?: boolean;
}

export interface MatchSetScoreResponse {
  matchInfo: MatchInfo
  setScore: MatchScore

}

interface MatchInfo {
  matchDate?: Date
  source: MatchSource
  sourceId?: string
  type: MatchType
  status: Status
  side1: string
  side2: string
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

export interface UpdateScoreRequest extends MatchScore {
  matchDate: Date
}
