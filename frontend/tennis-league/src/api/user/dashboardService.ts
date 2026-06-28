import {PlayerStatisticsResponse} from "../../model/dashboard.model";
import {userClient} from "../axiosClient";


export const getStatistics = async (params?: { limit?: number }): Promise<PlayerStatisticsResponse> => {
    return await userClient.get(`/me/statistics`, {params});
};

