import { GET, POST } from "./lib/http";
import { Nullish } from "./lib/types";
import { IPagination, IPaginationQuery } from "./pagination";

import { IAccount } from "models";

export function overview(): Promise<IAccount> {
  return GET<IAccount>('/accounts/overview');
}

/**
 * 用户注册
 * @param args
 * @returns
 */
export function signup(args: Partial<IAccount>): Promise<IAccount> {
  return POST<IAccount>('/accounts/signup', args);
}

/**
 * 用户登录
 * @param login
 * @param password
 * @returns
 */
export function signin(login: string, password: string): Promise<IAccount> {
  return POST<IAccount>('/accounts/signin', { login, password });
}

export interface IListTicketsArgs extends IPaginationQuery {
  q?: string | Nullish
}

export function list(params?: IListTicketsArgs): Promise<IPagination<IAccount>> {
  return GET('/accounts', { params })
}