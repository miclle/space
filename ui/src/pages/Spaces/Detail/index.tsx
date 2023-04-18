import { useEffect, useMemo, useState } from "react";
import { observer } from "mobx-react-lite";
import { useQuery } from "@tanstack/react-query";
import { Link, Outlet, useLocation, useParams } from "react-router-dom";
import { StringParam, useQueryParams, withDefault } from "use-query-params";
import { Avatar, Empty, Layout, Menu, Select, Skeleton, Tree } from "antd";
import { ItemType } from "antd/es/menu/hooks/useItems";
import { AiOutlinePlusSquare, AiOutlineSetting } from "react-icons/ai";
import { MdKeyboardArrowDown } from "react-icons/md";
import { BsBoxSeam } from "react-icons/bs";

import { IPage, ISpace } from "models";
import { Page } from "services";

import { SpaceStore, SpaceContext } from "./store";

const Spaces = observer(() => {

  const { key, page_id } = useParams() as { key: string, page_id: string };
  const location = useLocation();

  const store = useMemo<SpaceStore>(() => new SpaceStore(), []);
  const { space } = store;

  const [menuItems, setMenuItems] = useState<ItemType[]>([]);
  const [menuSelectedKeys, setMenuSelectedKeys] = useState<string[]>([]);
  const [expandedKeys, setExpandedKeys] = useState<React.Key[]>([]);

  const [query, setQuery] = useQueryParams({
    lang: withDefault(StringParam, space.lang),
  });

  const {
    isLoading,
  } = useQuery<ISpace>(['spaces.get', key], () => store.load(key), {
    enabled: key !== ''
  })

  const {
    data: pages,
    isLoading: pageTreeIsLoading,
  } = useQuery<IPage[]>(['spaces.pages', key, query], () => Page.list(key, query), {
    enabled: key !== '',
    initialData: [],
  })

  const onExpandHandler = (expandedKeysValue: React.Key[]) => {
    console.log(expandedKeysValue);
    setExpandedKeys(expandedKeysValue)
  }

  useEffect(() => {
    const items: ItemType[] = []

    items.push(
      { type: 'divider' },
      {
        key: `/spaces/${space.key}/setting/profile`,
        icon: <AiOutlineSetting />,
        label: <Link to={`/spaces/${space.key}/setting/profile`}>Space Settings</Link>
      },
      { type: 'divider' },
      {
        key: `/spaces/${space.key}/pages/new`,
        icon: <AiOutlinePlusSquare />,
        label: <>
          {
            page_id && page_id !== `${space.homepage_id}`
              ? <Link to={`/spaces/${space.key}/pages/new?parent_id=${page_id}`}>Create a page</Link>
              : <Link to={`/spaces/${space.key}/pages/new`}>Create a page</Link>
          }
        </>
      },
      { type: 'divider' },
    )

    setMenuItems(items)
  }, [space, page_id]);

  useEffect(() => {
    setExpandedKeys(store.expandedKeys)
  }, [store.expandedKeys]);

  useEffect(() => {
    setMenuSelectedKeys([location.pathname]);
  }, [location.pathname]);

  if (isLoading) {
    return <Skeleton active />;
  }

  if (!isLoading && space === undefined) {
    return <Empty />;
  }

  return (
    <SpaceContext.Provider value={store}>
      <Layout.Sider
        width={250}
        theme="light"
        className="app-sider-navigation app-sider-navigation-fixed"
      >
        <div className="sider-navigation-brand">
          <Avatar shape="square" size={24} icon={<BsBoxSeam />} />
          <Link to={`/spaces/${space.key}`}>{space.name}</Link>
        </div>

        <div className="sider-navigation">
          <Menu
            className="sider-navigation-menu"
            mode="inline"
            inlineIndent={8}
            defaultSelectedKeys={[location.pathname]}
            selectedKeys={menuSelectedKeys}
            items={menuItems}
          />

          <div className="side-navigation-tree">
            {
              pageTreeIsLoading && <Skeleton active />
            }

            {
              !pageTreeIsLoading && <>
                <Select
                  style={{ width: '100%', marginBottom: 12 }}
                  defaultValue={query.lang}
                  onChange={(lang) => setQuery({ ...query, lang })}
                >
                  <Select.Option value="en-US">English</Select.Option>
                  <Select.Option value="zh-CN">简体中文</Select.Option>
                </Select>

                <Tree
                  showLine={true}
                  fieldNames={{ title: 'title', key: 'id' }}
                  treeData={pages as any}
                  draggable={{ icon: false }}
                  blockNode
                  expandedKeys={expandedKeys}
                  onExpand={onExpandHandler}
                  switcherIcon={
                    <span className="anticon anticon-down app-tree-switcher-icon" style={{ fontSize: 14 }}>
                      <MdKeyboardArrowDown size={14} />
                    </span>
                  }
                  titleRender={(node: any) =>
                    <Link to={`/spaces/${space.key}/pages/${node.id}`}>{node.short_title}</Link>
                  }
                />
              </>
            }

          </div>
        </div>
      </Layout.Sider>

      <Layout.Content style={{ marginLeft: 250, padding: 25 }}>
        <Outlet />
      </Layout.Content>
    </SpaceContext.Provider>
  );
})

export default Spaces