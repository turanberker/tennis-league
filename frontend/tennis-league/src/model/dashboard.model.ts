import { MatchSource, MatchType } from "./match.model";

export interface PlayerStatisticsResponse {
    earnedSinglePoints: number;
    earnedDoublePoints: number;
    singlePoints: number;
    doublePoints: number;
}

export interface IncomingMatchResponse {
    matchId: string;
    matchDate: Date;
    matchType: MatchType;
    source: MatchSource;
    leagueId?: string;
    leagueName?: string;
    oppenentId?: string; // Rakip oyuncunun ID'si (tekler için) veya takım ID'si (çiftler için)
    oppenentName: string; // Rakip oyuncunun adı (tekler için) veya takım adı (çiftler için)

} 
