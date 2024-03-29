import { GET, PATCH, POST } from './lib/http';
import { Nullish } from './lib/types';

import { IPage, PageStatus } from 'models';

export interface ICreatePageArgs {
  parent_id?:        number
  page_id?:          number
  title:             string
  version:           string
  status:            PageStatus
  short_title?:      string
  body:              string
}

export function create(spaceKey: string, args: ICreatePageArgs): Promise<IPage> {
  return POST(`/spaces/${spaceKey}/pages`, args)
}
export interface IGetPageParams {
  lang?:    string | Nullish
  version?: string | Nullish
  depth?:   number | Nullish
}

export function list(spaceKey: string, params?: IGetPageParams): Promise<IPage[]> {
  return GET(`/spaces/${spaceKey}/pages`, { params })
}

export function get(spaceKey: string, id: string, params?: IGetPageParams): Promise<IPage> {
  return GET(`/spaces/${spaceKey}/pages/${id}`, { params })
}

export type IUpdatePageArgs = ICreatePageArgs;

export function update(spaceKey: string, id: string, args: IUpdatePageArgs): Promise<IPage> {
  return PATCH(`/spaces/${spaceKey}/pages/${id}`, args)
}
