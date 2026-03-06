export interface User {
  id: string;
  name: string;
  surname: string;
  email: string;
  approved: boolean;
  role: Role;
  playerId?: string;
}

export enum Role {
  ADMIN = 'ADMIN',
  PLAYER = 'PLAYER',
  COORDINATOR = 'COORDINATOR',
}

export const RoleLabels: Record<Role, string> = {
  [Role.ADMIN]: 'Admin',
  [Role.PLAYER]: 'Oyuncu',
  [Role.COORDINATOR]: 'Koordinatör',
};

export const RoleOptions = [
  { label: 'Admin', value: Role.ADMIN },
  { label: 'Oyuncu', value: Role.PLAYER },
  { label: 'Koordinatör', value: Role.COORDINATOR },
];
