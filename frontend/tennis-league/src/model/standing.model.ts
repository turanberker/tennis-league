import { TeamRef } from './team.model';

export interface ScoreBoardResponse extends TeamRef {
  order: number;
  played: number;
  won: number;
  lost: number;
  wonSets: number;
  lostSets: number;
  wonGames: number;
  lostGames: number;
  score: number;
}
