// Account 账户
export interface IAccount {
  id:         number
	created_at: number
	updated_at: number

	login:  string
  name:   string
  email:  string
  status: string
  avatar: string
}