// 分页配置
export interface MyPaginationConfig {
  current: number;
  pageSize: number;
  total: number;
  onChange: (current: number, pageSize: number) => void;
}

export enum PAGE_SIZE {
  TEN = 10,
  FIFTEEN = 15,
  TWENTY = 20,
}

export const PAGE_DEFAULT_INDEX = 1;