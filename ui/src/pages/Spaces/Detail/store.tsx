import React from "react";
import { makeAutoObservable, runInAction } from "mobx";

import { IPageTree, ISpace } from "models";
import { Space } from "services";

export class SpaceStore {

  space: ISpace = {} as ISpace

  tree: IPageTree = []

  constructor() {
    makeAutoObservable(this);
  }

  async load(key: string) {
    const info = await Space.get(key);
    runInAction(() => {
      this.space = info;
    });
    return info
  }

  async update(key: string, data: Partial<ISpace>) {
    const info = await Space.update(key, data)
    runInAction(() => {
      this.space = info
    })
    return info
  }
}

export const SpaceContext = React.createContext<SpaceStore>({} as any);

export const useSpaceContext = () => React.useContext(SpaceContext);
