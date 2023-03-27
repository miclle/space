import { Nullish } from "./lib/types";

export interface IPagination<T> {
  items: Array<T>;
  total: number;
  page?: number;
  page_size?: number;
}

export const PaginationDefault: IPagination<any> = {
  items: [],
  total: 0,
  page: 1,
  page_size: 30,
}

export interface IPaginationQuery {
  page?: number  | Nullish;
  page_size?: number  | Nullish;
}
