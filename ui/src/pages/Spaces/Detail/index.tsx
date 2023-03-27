import { useEffect, useMemo, useState } from "react";
import { observer } from "mobx-react-lite";
import { useQuery } from "@tanstack/react-query";
import { Link, Outlet, useLocation, useParams } from "react-router-dom";
import { StringParam, useQueryParams, withDefault } from "use-query-params";
import { Avatar, Button, Empty, Layout, Menu, Select, Skeleton } from "antd";
import { ItemType } from "antd/es/menu/hooks/useItems";
import { NodeRendererProps, Tree } from "react-arborist";
import { AiOutlinePlusSquare, AiOutlineSetting } from "react-icons/ai";
import { MdKeyboardArrowDown, MdKeyboardArrowRight } from "react-icons/md";
import { BsBoxSeam, BsDot } from "react-icons/bs";

import { IPageTree, IPageTreeNode, ISpace } from "models";
import { Page } from "services";

import { ClusterStore, SpaceContext } from "./store";

const Spaces = observer(() => {

  const { key, page_id } = useParams() as { key: string, page_id: string };
  const location = useLocation();

  const store = useMemo<ClusterStore>(() => new ClusterStore(), []);
  const { space } = store;

  const [menuItems, setMenuItems] = useState<ItemType[]>([]);
  const [menuSelectedKeys, setMenuSelectedKeys] = useState<string[]>([]);

  const [query, setQuery] = useQueryParams({
    lang: withDefault(StringParam, space.lang),
  });

  const [tree, setTree] = useState<IPageTree>([])

  const {
    isLoading,
  } = useQuery<ISpace>(['spaces.get', key], () => store.load(key), {
    enabled: key !== ''
  })

  const {
    isLoading: pageTreeIsLoading,
  } = useQuery<IPageTree>(['spaces.pages.tree', key, query.lang], () => Page.tree(key, query), {
    enabled: key !== '',
    onSuccess: (data) => {
      setTree(data)
    },
  })

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

            <Select
              style={{ width: '100%', marginBottom: 12 }}
              defaultValue={query.lang}
              onChange={(lang) => setQuery({ ...query, lang })}
            >
              <Select.Option value="en-US">English</Select.Option>
              <Select.Option value="zh-CN">简体中文</Select.Option>
            </Select>

            <Tree<IPageTreeNode>
              data={tree}
              width="100%"
            >
              {({ node, style }: NodeRendererProps<IPageTreeNode>) => {
                /* This node instance can do many things. See the API reference. */
                return (
                  <div style={style}>
                    <Button
                      type="text"
                      size="small"
                      icon={
                        node.isLeaf ? <BsDot /> :
                          node.isClosed ? <MdKeyboardArrowRight /> : <MdKeyboardArrowDown />
                      }
                      onClick={() => node.toggle()}
                    />
                    <Link to={`/spaces/${space.key}/pages/${node.data.id}`}>{node.data.short_title}</Link>
                  </div>
                );
              }}
            </Tree>
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