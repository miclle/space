export enum PageStatus {
  draft      = 'Draft',
  published  = 'Published',
  offline    = 'Offline',
  deprecated = 'Deprecated',
}

export interface IPage {
  id:               number
  lang:             string
  version:          string
  status:           PageStatus
  title:            string
  short_title:      string
  body:             string
  html:             string

  created_at: number
  updated_at: number

  children_count: number
  children:       IPage[]
  parents?:       IPageParent[]
}

export interface IPageParent {
  id:             number
  parent_page_id: number
  lang:           string
  version:        string
  status:         PageStatus
  title:          string
  short_title:    string
}
