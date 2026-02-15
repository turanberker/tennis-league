export enum Status {
  PERDING = 'PENDING',
  COMPLETED = 'COMPLETED',
  CANCELLED = 'CANCELLED',
}

export const MatchStatusLabels: Record<Status, string> = {
  [Status.PERDING]: 'Beklemede',
  [Status.COMPLETED]: 'Tamamlandı',
  [Status.CANCELLED]: 'İptal Edildi',
};

export interface LeagueFixtureMatchResponse {
  id: string;
  team1: TeamRefResponse;
  team2: TeamRefResponse;
  status: Status;
  matchDate: Date;
}

interface TeamRefResponse {
  id: string;
  name: string;
}
