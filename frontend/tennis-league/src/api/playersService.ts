import axiosClient from './axiosClient';
import { Player, CreatePlayerRequest } from '../model/player.model';

export const getPlayers = async (name?: string): Promise<Player[]> => {
  return axiosClient.get<Player[]>('/player/list', {
    params: { name },
  });
};

export const getPlayerByUuid = async (uuid: string) => {
  return axiosClient.get<Player>(`/player/${uuid}`);
};

export const createPlayer = async (payload: CreatePlayerRequest) => {
  return axiosClient.post<number>('/player', payload);
};
