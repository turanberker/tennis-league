// hooks/useLeague.ts
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { getLeagueById } from '../api/leagueService';
import { League } from '../model/league.model';

export const useLeague = (leagueId: string | undefined) => {
    const queryClient = useQueryClient();

    const query = useQuery<League, Error>({
        queryKey: ['league', leagueId],
        queryFn: async () => {
            const response = await getLeagueById(leagueId!);
            return response;
        },
        enabled: !!leagueId,
        staleTime: 1000 * 60 * 5,
    });

    // Cache'i manuel güncellemek için yardımcı fonksiyon
    const updateLeagueCache = (newData: Partial<League>) => {
        if (!leagueId) return;

        queryClient.setQueryData(['league', leagueId], (oldData: League | undefined) => {
            if (!oldData) return oldData;
            return {
                ...oldData,
                ...newData // Sadece gönderdiğin alanları günceller
            };
        });
    };

    return { ...query, updateLeagueCache };
};