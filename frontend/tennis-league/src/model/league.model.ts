export interface League {
  id: String;
  name: String;
  format: LEAGUE_FORMAT;
  category: LEAGUE_CATEGORY;
  processType: LEAGUE_PROCESS_TYPE,
  status: LEAGUE_STATUS,
  totalAttentance: number;
  fixtureCreatedDate?: Date;
  startedDate?: Date;
  endDate?: Date;
}

export interface LeagueListResponse {
  id: String;
  name: String;
  format: LEAGUE_FORMAT;
  category: LEAGUE_CATEGORY;
  processType: LEAGUE_PROCESS_TYPE,
  status: LEAGUE_STATUS,
  totalAttentance: number;
  coordinatorUserIds: String[];
}

export interface PersistLeagueRequest {
  name: string;
  format: LEAGUE_FORMAT;
  category: LEAGUE_CATEGORY;
  processType: LEAGUE_PROCESS_TYPE;
}

export enum LEAGUE_FORMAT {
  Single = "SINGLE",
  Double = "DOUBLE",
  Team = "TEAM",
}

export const League_Format_Labels: Record<LEAGUE_FORMAT, string> = {
  [LEAGUE_FORMAT.Single]: "Single",
  [LEAGUE_FORMAT.Double]: "Double",
  [LEAGUE_FORMAT.Team]: "Takım",
};

export const League_Format_Options = [
  { label: "Single", value: LEAGUE_FORMAT.Single },
  { label: "Double", value: LEAGUE_FORMAT.Double },
  { label: "Takım", value: LEAGUE_FORMAT.Team },
];

export enum LEAGUE_CATEGORY {
  MIX = "Mix",
  MALE = "Erkek",
  FEMALE = "Kadın",
}

export const League_Category_Labels: Record<LEAGUE_CATEGORY, string> = {
  [LEAGUE_CATEGORY.MIX]: "Mix",
  [LEAGUE_CATEGORY.MALE]: "Erkek",
  [LEAGUE_CATEGORY.FEMALE]: "Kadın",
};

export const League_Category_Options = [
  { label: "Mix", value: LEAGUE_CATEGORY.MIX },
  { label: "Erkek", value: LEAGUE_CATEGORY.MALE },
  { label: "Kadın", value: LEAGUE_CATEGORY.FEMALE },
];

export enum LEAGUE_PROCESS_TYPE {
  FIXTURE = "FIXTURE",
  DEFI = "DEFI",
}

export const League_Process_Type_Labels: Record<LEAGUE_PROCESS_TYPE, string> = {
  [LEAGUE_PROCESS_TYPE.FIXTURE]: "Fixture",
  [LEAGUE_PROCESS_TYPE.DEFI]: "Defi",
};

export const League_Process_Type_Options = [
  { label: "Fixture", value: LEAGUE_PROCESS_TYPE.FIXTURE },
  { label: "Defi", value: LEAGUE_PROCESS_TYPE.DEFI },
];

export enum LEAGUE_STATUS {
  DRAFT = "DRAFT",
  ACTIVE = "ACTIVE",
  COMPLETED = "COMPLETED",
}

export const League_Status_Labels: Record<LEAGUE_STATUS, string> = {
  [LEAGUE_STATUS.DRAFT]: "Taslak",
  [LEAGUE_STATUS.ACTIVE]: "Aktif",
  [LEAGUE_STATUS.COMPLETED]: "Tamamlandı",
};

export const League_Status_Options = [
  { label: "Taslak", value: LEAGUE_STATUS.DRAFT },
  { label: "Aktif", value: LEAGUE_STATUS.ACTIVE },
  { label: "Tamamlandı", value: LEAGUE_STATUS.COMPLETED },
];
