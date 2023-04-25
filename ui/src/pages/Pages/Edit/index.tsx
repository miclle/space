import { observer } from "mobx-react-lite";
import { Link, useNavigate, useParams } from "react-router-dom";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { map } from "lodash";
import { Breadcrumb, Button, Col, Empty, Form, Input, notification, Row, Select, Skeleton, Typography } from "antd";

import CodeEditor from "components/CodeEditor";
import { IPage } from 'models';
import { AxiosResponse, IErrorMessage, Markdown, Page } from "services";
import { useSpaceContext } from "pages/Spaces/Detail/store";

const EditPage = observer(() => {

  const navigate = useNavigate();
  const client = useQueryClient()

  const store = useSpaceContext();
  const { space } = store;

  const { page_id } = useParams() as { key: string, page_id: string };

  const [form] = Form.useForm();

  const {
    isLoading,
    data: page,
  } = useQuery<IPage>(['spaces.pages.get', page_id], () => Page.get(space.key, page_id), {
    enabled: !!page_id,
  })

  const submit = (status: string) => {
    form.setFieldValue('status', status)
    form.submit()
  }

  const handleFormFinish = async (values: Page.IUpdatePageArgs) => {
    Page.update(space.key, page_id, values)
      .then((page: IPage) => {
        client.prefetchQuery(['spaces.pages', space.key, { lang: page.lang }]);
        client.setQueryData(['spaces.pages.get', page_id, { lang: page.lang }], page);
        navigate(`/spaces/${space.key}/pages/${page.id}?lang=${page.lang}`);
        notification.success({ message: 'Page updated successfully' });
      })
      .catch((resp: AxiosResponse<IErrorMessage>) => {
        notification.error({
          key: 'create-page-error',
          message: 'Page update failure',
          description: map(resp.data.message, (value, key) => value).join('\n')
        });
      })
  }

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
              { title: page.title }
            ]}
          />
        </Col>
      </Row >

      <Form
        name="page-form"
        form={form}
        preserve={false}
        layout="vertical"
        onFinish={handleFormFinish}
        style={{ paddingBlock: 16 }}
      >
        <Form.Item name="title" label="Title" rules={[{ required: true }]} initialValue={page.title}>
          <Input placeholder="Give this page a title" />
        </Form.Item>

        <Row gutter={24}>
          <Col span={12}>
            <Form.Item
              name="short_title"
              label={<>Short Title <Typography.Text type="secondary" style={{ marginLeft: 5 }}>(The short title will be used in the navigation menu)</Typography.Text></>}
              initialValue={page.short_title}
            >
              <Input placeholder="Blank indicates the same as the page title" />
            </Form.Item>
          </Col>
          <Col span={6}>
            <Form.Item name="lang" label="Lang" initialValue={page.lang}>
              <Select>
                <Select.Option value="en-US">English</Select.Option>
                <Select.Option value="zh-CN">简体中文</Select.Option>
              </Select>
            </Form.Item>
          </Col>
          <Col span={6} />
        </Row>

        <Form.Item name="body" rules={[{ required: true }]} initialValue={page.body}>
          <CodeEditor
            lang="markdown"
            preview={(value) => Markdown.preview(value)}
          />
        </Form.Item>

        <Form.Item name="status" noStyle initialValue="published">
          <Input type="hidden" readOnly />
        </Form.Item>

        <Row gutter={12}>
          <Col flex="auto">
            <Link to={`/spaces/${space.key}/pages/${page.id}`}><Button>Cancel</Button></Link>
          </Col>
          <Col>
            <Button onClick={() => submit('draft')}>Save as draft</Button>
          </Col>
          <Col>
            <Button type="primary" onClick={() => submit('published')}>Publish</Button>
          </Col>
        </Row>
      </Form>
    </>
  );
})

export default EditPage