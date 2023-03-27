import React from "react";
import { makeAutoObservable, runInAction } from "mobx";
import { IAccount } from "models";
import { Account } from "services";

export class GlobalStore {

  account: IAccount = {} as IAccount

  constructor() {
    makeAutoObservable(this);
  }

  async loadOverview() {
    const info = await Account.overview();
    runInAction(() => {
      this.account = info;
    });
    return info
  }

}

export const GlobalContext = React.createContext<GlobalStore>({} as any);

export const useGlobalContext = () => React.useContext(GlobalContext);
