import React from "react";
import { computed, makeAutoObservable, runInAction } from "mobx";

import { IPage, ISpace } from "models";
import { Space } from "services";

export class SpaceStore {

  space: ISpace = {} as ISpace
  currentPage?: IPage

  constructor() {
    makeAutoObservable(this, {
      expandedKeys: computed,
    });
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

  setCurrentPage(page: IPage) {
    runInAction(() => {
      this.currentPage = page
    })
  }

  get expandedKeys(): React.Key[] {
    const keys: React.Key[] = []

    if (this.currentPage) {
      this.currentPage.parents?.forEach((page) => {
        keys.push(page.id)
      })
      keys.push(this.currentPage.id)
    }

    return keys
  }
}

export const SpaceContext = React.createContext<SpaceStore>({} as any);

export const useSpaceContext = () => React.useContext(SpaceContext);
