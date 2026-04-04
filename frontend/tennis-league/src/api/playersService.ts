import axiosClient from './axiosClient';
import { Player, CreatePlayerRequest, Sex, PlayerResponse } from '../model/player.model';

export const getPlayers = async (data?: { name?: string, sex?: Sex }): Promise<PlayerResponse[]> => {
  return axiosClient.get<PlayerResponse[]>('/player/list', {
    params: data,
  });
};

export const getPlayerByUuid = async (uuid: string) => {
  return axiosClient.get<Player>(`/player/${uuid}`);
};

export const createPlayer = async (payload: CreatePlayerRequest) => {
  return axiosClient.post<number>('/player', payload);
};

export const getUnassignedPlayers = async (sex: Sex) => {
  return axiosClient.get<Player[]>('/player/unassigned-players', {
    params: { sex: sex },
  });
};

export const assignPlayerToUser = async (playerId: string, params: { userId: string }) => {
  return axiosClient.put<Player[]>(`/player/${playerId}/assign-to-user`, null, {
    params: params
  });
};