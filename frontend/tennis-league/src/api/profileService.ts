
import axiosClient from './axiosClient';

import { ChangePasswordRequest } from '../model/user.model';

export const changeMyPassword = async (changePasswordRequest: ChangePasswordRequest): Promise<string> => {
    return axiosClient.patch<string>('/user/profile/change-password', changePasswordRequest);
};
