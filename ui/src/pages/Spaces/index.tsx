import { observer } from "mobx-react-lite";
import { PageHeader } from "@ant-design/pro-components";
import { useQuery } from "@tanstack/react-query";
import { DecodedValueMap, NumberParam, StringParam, useQueryParams, withDefault } from "use-query-params";
import { debounce } from "lodash";
import { Button, Col, Form, Input, Layout, Row, Skeleton, Table } from "antd";
import { AiOutlineSearch, AiOutlineSetting } from "react-icons/ai";

import { ISpace } from "models";
import { IPagination, PaginationDefault, Space } from "services";
import { Link } from "react-router-dom";

const Spaces = observer(() => {

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
  } = useQuery<IPagination<ISpace>>(['spaces.list', query], () => Space.list(query), {
    keepPreviousData: true,
    initialData: PaginationDefault,
  })

  const columns = [
    {
      title: '空间',
      dataIndex: 'name',
      key: 'name',
      render: (name: string, space: ISpace) => <Link to={`/spaces/${space.key}`}>{name}</Link>
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
    },
    {
      title: <AiOutlineSetting />,
      key: 'settings',
      width: 50,
      render: (space: ISpace) => <Link to={`/spaces/${space.key}/setting/profile`}><AiOutlineSetting /></Link>
    }
  ];

  if (isLoading) {
    return <Skeleton active />;
  }

  return (
    <Layout.Content>
      <div className="container">
        <PageHeader ghost={false} title="空间列表">
          <Row>
            <Col flex="auto">
              <Form layout="inline">

                <Form.Item name="q" initialValue={query.q}>
                  <Input allowClear style={{ width: 280 }} onChange={(e) => search({ ...query, q: e.target.value || undefined })} />
                </Form.Item>

                <Form.Item style={{ marginRight: 0 }}>
                  <Button onClick={() => refetch()}><AiOutlineSearch /></Button>
                </Form.Item>
              </Form>
            </Col>
            <Col>
              <Button type="primary" href="/spaces/new">创建空间</Button>
            </Col>
          </Row>
        </PageHeader>

        <Table
          columns={columns}
          loading={isFetching}
          dataSource={pagination.items}
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

export default Spaces