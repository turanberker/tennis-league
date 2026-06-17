import axiosClient from '../axiosClient';
import { Player, CreatePlayerRequest, Sex, PlayerResponse } from '../../model/player.model';


const USER_API_URL = process.env.REACT_APP_USER_URL || 'http://localhost:8000';

export const getPlayers = async (data?: { name?: string, sex?: Sex }): Promise<PlayerResponse[]> => {
  return axiosClient.get<PlayerResponse[]>(`${USER_API_URL}/player/list`, {
    params: data,
  });
};

export const getPlayerByUuid = async (uuid: string) => {
  return axiosClient.get<Player>(`${USER_API_URL}/player/${uuid}`);
};

export const createPlayer = async (payload: CreatePlayerRequest) => {
  return axiosClient.post<number>(`${USER_API_URL}/player`, payload);
};

export const getUnassignedPlayers = async (sex: Sex) => {
  return axiosClient.get<Player[]>(`${USER_API_URL}/player/unassigned-players`, {
    params: { sex: sex },
  });
};

export const assignPlayerToUser = async (playerId: string, params: { userId: string }) => {
  return axiosClient.put<Player[]>(`${USER_API_URL}/player/${playerId}/assign-to-user`, null, {
    params: params
  });
};