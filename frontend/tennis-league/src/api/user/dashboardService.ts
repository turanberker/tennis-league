import {PlayerStatisticsResponse} from "../../model/dashboard.model";
import axiosClient from "../axiosClient";


const USER_API_URL = process.env.REACT_APP_USER_URL || 'http://localhost:8000';

export const getStatistics = async (params?: { limit?: number }): Promise<PlayerStatisticsResponse> => {
    return await axiosClient.get(`${USER_API_URL}/me/statistics`, {params});
};

