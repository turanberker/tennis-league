// hooks/useLeague.ts
import { useQuery, UseQueryResult } from '@tanstack/react-query';
import { getLeagueById } from '../api/leagueService';
import { League } from '../model/league.model';

export const useLeague = (leagueId: string | undefined): UseQueryResult<League, Error> => {
    return useQuery<League, Error>({
        queryKey: ['league', leagueId],
        queryFn: async () => {
            // getLeagueById fonksiyonunun AxiosResponse<League> döndüğünü varsayıyorum
            const response = await getLeagueById(leagueId!);
            return response;
        },
        enabled: !!leagueId, // leagueId varsa çalışır
        staleTime: 1000 * 60 * 5, // 5 dakika boyunca veriyi taze kabul et
    });
};