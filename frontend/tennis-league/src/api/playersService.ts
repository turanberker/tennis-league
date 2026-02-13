import axiosClient from './axiosClient';

export interface Player {
  id: number;
  name: string;
  surname: string;
  uuid: string;
  userId: number;
}

export interface CreatePlayerRequest {
  name: string;
  surname: string;
  userId: number;
}

export const getPlayers = async (name?: string) => {
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
