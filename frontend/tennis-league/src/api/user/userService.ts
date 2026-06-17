import axiosClient from '../axiosClient';

import {User} from '../../model/user.model';

const USER_API_URL = process.env.REACT_APP_USER_URL || 'http://localhost:8000';

export const getUsers = async (): Promise<User[]> => {
    return axiosClient.get<User[]>(`${USER_API_URL}/user/list`);
};
