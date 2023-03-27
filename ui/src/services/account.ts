import { GET } from "./lib/http";

import { IAccount } from "models";

export function overview(): Promise<IAccount> {
  return GET<IAccount>('/overview');
}
