import { PlayerResponse } from "../model/player.model";
import  {mainClient} from "./axiosClient";

export const getTeamMembers = async (id: string): Promise<PlayerResponse[]> => {
    return await mainClient.get(`double-team/${id}/members`);
};
