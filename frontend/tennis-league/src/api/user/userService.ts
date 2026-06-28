

import {User} from '../../model/user.model';
import {userClient} from "../axiosClient";

export const getUsers = async (): Promise<User[]> => {
    return userClient.get<User[]>(`/user/list`);
};
