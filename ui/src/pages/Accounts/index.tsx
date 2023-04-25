import React from "react";
import { observer } from "mobx-react-lite";
import { useQuery } from "@tanstack/react-query";
import dayjs from "dayjs";
import { Col, Form, Input, Layout, Row, Skeleton, Table } from "antd";
import { Account, IPagination, PaginationDefault } from "services";
import { IAccount } from "models";
import { ColumnsType } from "antd/es/table";
import { PageHeader } from "@ant-design/pro-components";
import { DecodedValueMap, NumberParam, StringParam, useQueryParams, withDefault } from "use-query-params";
import { debounce } from "lodash";

const Accounts = observer(() => {

  const [query, setQuery] = useQueryParams({
    page: withDefault(NumberParam, 1),
    page_size: withDefault(NumberParam, 25),

    q: StringParam,
  });

  const search = debounce((query: DecodedValueMap<any>) => setQuery({ ...query, page: 1 }), 500);

  const {
    isLoading,
    data: pagination,
    isFetching,
    refetch
  } = useQuery<IPagination<IAccount>>(['accounts.list', query], () => Account.list(query), {
    keepPreviousData: false,
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
    <Layout.Content>
      <div className="container">
        <PageHeader ghost={false} title="All accounts">
          <Row>
            <Col flex="auto">
              <Form layout="inline">
                <Form.Item name="q" initialValue={query.q}>
                  <Input.Search
                    onChange={(e) => search({ ...query, q: e.target.value || undefined })}
                    onSearch={(q) => refetch()}
                    allowClear
                    style={{ width: 280 }}
                  />
                </Form.Item>
              </Form>
            </Col>
            <Col />
          </Row>
        </PageHeader>

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
      </div>
    </Layout.Content>
  );
})

export default Accounts