import React, { useState } from "react";
import { useLocation } from "react-router-dom";
import { observer } from "mobx-react-lite";
import { useQuery } from "@tanstack/react-query";
import dayjs from "dayjs";
import { Layout, Menu, Skeleton, Table } from "antd";
import { ItemType } from "antd/es/menu/hooks/useItems";
import { Account, IPagination, PaginationDefault } from "services";
import { IAccount } from "models";
import { ColumnsType } from "antd/es/table";

const Accounts = observer(() => {

  const location = useLocation();

  const [menuItems, setMenuItems] = useState<ItemType[]>([]);
  const [menuSelectedKeys, setMenuSelectedKeys] = useState<string[]>([]);

  const {
    isLoading,
    data: pagination,
    isFetching,
  } = useQuery<IPagination<IAccount>>(['accounts.list'], () => Account.list(), {
    keepPreviousData: true,
    initialData: PaginationDefault,
  })

  const columns: ColumnsType<IAccount> = [
    {
      title: 'Login',
      dataIndex: 'login',
      key: 'login',
    },
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: 'Email',
      dataIndex: 'email',
      key: 'email',
      ellipsis: true,
    },
    {
      title: 'Created',
      dataIndex: 'created_at',
      width: 250,
      ellipsis: true,
      render: (timestamp: number) => dayjs.unix(timestamp).format('LLL')
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      width: 100
    },
  ];

  if (isLoading) {
    return <Skeleton active />;
  }

  return (
    <>
      <Layout.Sider
        width={250}
        theme="light"
        className="app-sider-navigation app-sider-navigation-fixed"
      >
        <Menu
          className="sider-navigation-menu"
          mode="inline"
          inlineIndent={8}
          defaultSelectedKeys={[location.pathname]}
          selectedKeys={menuSelectedKeys}
          items={menuItems}
        />
      </Layout.Sider>

      <Layout.Content style={{ marginLeft: 250, padding: 25 }}>
        <Table
          columns={columns}
          loading={isFetching}
          dataSource={pagination.items}
          rowKey="id"
          pagination={{
            size: 'default',
            current: pagination?.page,
            pageSize: pagination?.page_size,
            total: pagination?.total,
            showTotal: (total, range) => `${range[0]}-${range[1]} of ${total} items`,
          }}
        />
      </Layout.Content>
    </>
  );
})

export default Accounts