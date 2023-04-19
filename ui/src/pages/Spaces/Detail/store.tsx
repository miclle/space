import React from "react";
import { makeAutoObservable, runInAction } from "mobx";
import { union } from "lodash";

import { IPage, ISpace } from "models";
import { Space } from "services";

export class SpaceStore {

  space: ISpace = {} as ISpace;

  currentPage?: IPage;

  expandedKeys: React.Key[] = [];

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

  setCurrentPage(page: IPage) {
    runInAction(() => {
      this.currentPage = page

      const keys: React.Key[] = [...this.expandedKeys]

      this.currentPage.parents?.forEach((page) => {
        keys.push(page.id)
      })

      if (page.children_count > 0){
        keys.push(this.currentPage.id)
      }

      this.expandedKeys = union(keys);
    })
  }

  setExpandedKeys(expandedKeys: React.Key[]) {
    console.log(expandedKeys);

    this.expandedKeys = expandedKeys
  }

}

export const SpaceContext = React.createContext<SpaceStore>({} as any);

export const useSpaceContext = () => React.useContext(SpaceContext);
