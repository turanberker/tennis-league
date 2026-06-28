import { IncomingMatchResponse, PlayerStatisticsResponse } from "../model/dashboard.model";
import {mainClient} from "./axiosClient";


export const getIncomingMathces = async (params?: { limit?: number }): Promise<IncomingMatchResponse[]> => {
    return await mainClient.get("me/incoming-matches", { params });
}