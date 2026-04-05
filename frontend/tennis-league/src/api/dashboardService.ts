import { PlayerStatisticsResponse } from "../model/dashboard.model";
import axiosClient from "./axiosClient";

export const getStatistics = async (params?: { limit?: number }): Promise<PlayerStatisticsResponse> => {
    return await axiosClient.get("me/statistics", { params });
};