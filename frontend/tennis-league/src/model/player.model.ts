export enum Sex {
  Male = 'M',
  Female = 'F',
}

export const SexLabels: Record<Sex, string> = {
  [Sex.Male]: 'Erkek',
  [Sex.Female]: 'Kadın',
};

export const SexOptions = [
  { label: 'Erkek', value: Sex.Male },
  { label: 'Kadın', value: Sex.Female },
];
export interface Player {
  id: number;
  name: string;
  surname: string;
  uuid: string;
  userId?: number;
  sex: Sex;
}

export interface CreatePlayerRequest {
  name: string;
  surname: string;
}
