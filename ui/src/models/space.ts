export enum SpaceStatus {
  offline = 'Offline',
  online  = 'Online',
}

export interface ISpace {
  id:          number
  created_at:  number
  updated_at:  number

  name:          string
  key:           string
  lang:          string
  fallback_lang: string
  homepage_id:   number
  description:   string
  avatar:        string
  status:        SpaceStatus
}
