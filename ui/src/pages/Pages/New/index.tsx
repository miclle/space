import { observer } from "mobx-react-lite";
import { Link, useNavigate } from "react-router-dom";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { NumberParam, StringParam, useQueryParams } from "use-query-params";
import { map } from "lodash";
import { Button, Col, Form, Input, notification, Row, Typography } from "antd";
import { PageHeader } from "@ant-design/pro-components";

import { AxiosResponse, IErrorMessage, Page } from "services";

import Editor from "components/Editor";
import { IPage } from "models";
import { useSpaceContext } from "pages/Spaces/Detail/store";

const NewPage = observer(() => {

  const navigate = useNavigate();
  const client = useQueryClient()

  const store = useSpaceContext();
  const { space } = store;

  const [form] = Form.useForm();

  const [query] = useQueryParams({
    parent_id: NumberParam,
    lang: StringParam,
  });

  const {
    data: parentPage,
  } = useQuery<IPage>(['spaces.pages.get', query.parent_id, query], () => Page.get(space.key, `${query.parent_id}`, query), {
    enabled: !!query.parent_id,
  })

  const submit = (status: string) => {
    form.setFieldValue('status', status)
    form.submit()
  }

  const handleFormFinish = async (values: Page.ICreatePageArgs) => {
    Page.create(space.key, values)
      .then((page: IPage) => {
        client.prefetchQuery(['spaces.pages', space.key, {lang: page.lang}]);
        navigate(`/spaces/${space.key}/pages/${page.id}?lang=${page.lang}`);
        notification.success({ message: 'Page created successfully' });
      })
      .catch((resp: AxiosResponse<IErrorMessage>) => {
        notification.error({
          key: 'create-page-error',
          message: 'Page creation failure',
          description: map(resp.data.message, (value, key) => value).join('\n')
        });
      })
  }

  return (
    <>
      <PageHeader
        ghost={false}
        breadcrumb={{
          items: [
            { title: <Link to={`/spaces/${space.key}`}>{space.name}</Link> },
            ...(parentPage?.parents?.map((parent) => ({ title: <Link to={`/spaces/${space.key}/pages/${parent.id}`}>{parent.short_title}</Link> })) || []),
            ...(parentPage ? [{ title: <Link to={`/spaces/${space.key}/pages/${parentPage.id}`}>{parentPage.short_title}</Link> }] : []),
            { title: 'Add page' }
          ]
        }}
      />

      <Form
        name="page-form"
        form={form}
        preserve={false}
        layout="vertical"
        onFinish={handleFormFinish}
        style={{ padding: 16 }}
      >
        {
          parentPage &&
          <Form.Item name="parent_id" noStyle initialValue={parentPage.id}>
            <Input type="hidden" readOnly />
          </Form.Item>
        }

        <Form.Item name="title" label="Title" rules={[{ required: true }]}>
          <Input placeholder="Give this page a title" />
        </Form.Item>

        <Row gutter={24}>
          <Col span={12}>
            <Form.Item
              name="short_title"
              label={<>Short Title <Typography.Text type="secondary" style={{ marginLeft: 5 }}>(The short title will be used in the navigation menu)</Typography.Text></>}
            >
              <Input placeholder="Blank indicates the same as the page title" />
            </Form.Item>
          </Col>
          <Col span={12}>

          </Col>
        </Row>

        <Form.Item name="body" rules={[{ required: true }]}>
          <Editor />
        </Form.Item>

        <Form.Item name="status" noStyle initialValue="published">
          <Input type="hidden" readOnly />
        </Form.Item>

        <Row gutter={12}>
          <Col flex="auto">
            <Link to={`/spaces/${space.key}`}><Button>Cancel</Button></Link>
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

export default NewPage