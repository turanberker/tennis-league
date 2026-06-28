
import  {mainClient} from '../axiosClient';

import { ChangePasswordRequest } from '../../model/user.model';

export const changeMyPassword = async (changePasswordRequest: ChangePasswordRequest): Promise<string> => {
    return mainClient.patch<string>('/user/profile/change-password', changePasswordRequest);
};
