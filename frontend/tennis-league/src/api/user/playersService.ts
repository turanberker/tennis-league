
import { Player, CreatePlayerRequest, Sex, PlayerResponse } from '../../model/player.model';
import {userClient} from "../axiosClient";


export const getPlayers = async (data?: { name?: string, sex?: Sex }): Promise<PlayerResponse[]> => {
  return userClient.get<PlayerResponse[]>(`/player/list`, {
    params: data,
  });
};

export const getPlayerByUuid = async (uuid: string) => {
  return userClient.get<Player>(`/player/${uuid}`);
};

export const createPlayer = async (payload: CreatePlayerRequest) => {
  return userClient.post<number>(`/player`, payload);
};

export const getUnassignedPlayers = async (sex: Sex) => {
  return userClient.get<Player[]>(`/player/unassigned-players`, {
    params: { sex: sex },
  });
};

export const assignPlayerToUser = async (playerId: string, params: { userId: string }) => {
  return userClient.put<Player[]>(`/player/${playerId}/assign-to-user`, null, {
    params: params
  });
};