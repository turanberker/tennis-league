import axiosClient from './axiosClient';
import { Player, CreatePlayerRequest, Sex } from '../model/player.model';

export const getPlayers = async (data?: { name?: string, sex?: Sex }): Promise<Player[]> => {
  return axiosClient.get<Player[]>('/player/list', {
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
