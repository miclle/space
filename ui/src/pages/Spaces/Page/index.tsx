import { observer } from "mobx-react-lite";
import { Link, useParams } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { StringParam, useQueryParams, withDefault } from "use-query-params";
import dayjs from "dayjs";
import { Breadcrumb, Button, Col, Empty, Row, Select, Skeleton, Space, Tag, Typography } from "antd";
import { PageHeader } from "@ant-design/pro-components";

import { IPage } from 'models';
import { Page } from "services";

import { useSpaceContext } from "../Detail/store";
import { AiOutlineEdit, AiOutlineFileAdd } from "react-icons/ai";

const PageDetail = observer(() => {

  const store = useSpaceContext();
  const { space } = store;

  const { page_id } = useParams() as { page_id: string };

  const [query] = useQueryParams({
    lang: withDefault(StringParam, space.lang),
  });

  const {
    isLoading,
    data: page,
  } = useQuery<IPage>(['spaces.pages.get', page_id, query], () => Page.get(space.key, page_id, query), {
    enabled: !!page_id,
    onSuccess(data) {
      store.setCurrentPage(data);
    },
  })

  if (isLoading) {
    return <Skeleton active />;
  }

  if (!isLoading && page === undefined) {
    return <Empty />;
  }

  return (
    <>
      <Row>
        <Col flex="auto">
          <Breadcrumb
            items={[
              { title: <Link to={`/spaces/${space.key}`}>{space.name}</Link> },
              ...(page.parents?.map((parent) => ({ title: <Link to={`/spaces/${space.key}/pages/${parent.id}`}>{parent.short_title}</Link> })) || []),
              { title: page.title },
            ]}
          />
        </Col>
        <Col>
          <Space>
            <Link to={`/spaces/${space.key}/pages/${page.id}/edit`}><Button>Edit</Button></Link>
            <Select defaultValue={page.lang}>
              <Select.Option value="en-US">English</Select.Option>
              <Select.Option value="zh-CN">简体中文</Select.Option>
            </Select>
          </Space>
        </Col>
      </Row>

      <PageHeader
        ghost={false}
        title={page.title}
        tags={<Tag color="blue">{page.lang}</Tag>}
      >
        <Space>
          <Typography.Text type="secondary"><AiOutlineFileAdd />{dayjs.unix(page.created_at).format('YYYY-MM-DD HH:mm:ss')}</Typography.Text>
          <Typography.Text type="secondary"><AiOutlineEdit />{dayjs.unix(page.updated_at).format('YYYY-MM-DD HH:mm:ss')}</Typography.Text>
        </Space>
      </PageHeader>

      <div className="page-content" style={{ padding: '0 16px' }} dangerouslySetInnerHTML={{ __html: page.html || '' }} />
    </>
  );
})

export default PageDetail