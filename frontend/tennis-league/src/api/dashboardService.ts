import { IncomingMatchResponse, PlayerStatisticsResponse } from "../model/dashboard.model";
import axiosClient from "./axiosClient";

export const getIncomingMathces = async (params?: { limit?: number }): Promise<IncomingMatchResponse[]> => {
    return await axiosClient.get("me/incoming-matches", { params });
}