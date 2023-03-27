import { GET, PATCH, POST } from './lib/http';
import { Nullish } from './lib/types';
import { IPagination, IPaginationQuery } from './pagination';

import { ISpace } from 'models';

export interface IListSpacesArgs extends IPaginationQuery {
  q?: string | Nullish
}

export function list(params?: IListSpacesArgs): Promise<IPagination<ISpace>> {
  return GET('/spaces', { params })
}

export function create(args: Partial<ISpace>): Promise<ISpace> {
  return POST('/spaces', args)
}

export function get(key: string): Promise<ISpace> {
  return GET(`/spaces/${key}`)
}

export function update(key: string, args: Partial<ISpace>): Promise<ISpace> {
  return PATCH(`/spaces/${key}`, args)
}
