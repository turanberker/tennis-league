import { PlayerResponse } from "../model/player.model";
import axiosClient from "./axiosClient";

export const getTeamMembers = async (id: string): Promise<PlayerResponse[]> => {
    return await axiosClient.get(`double-team/${id}/members`);
};
