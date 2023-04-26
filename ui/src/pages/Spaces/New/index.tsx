import React, { useState } from 'react';
import { observer } from 'mobx-react-lite';
import { Link, useNavigate } from 'react-router-dom';
import { map } from 'lodash';
import { Button, Form, Input, Layout, notification, Radio, Select, Space as AntSpace, Switch } from 'antd';
import { PageHeader } from '@ant-design/pro-components';

import { ISpace } from 'models';
import { AxiosResponse, Space, IErrorMessage } from "services";

const NewSpace = observer(() => {
  const navigate = useNavigate();

  const [form] = Form.useForm();
  const [multilingual, setMultilingual] = useState(false);

  const handleFormFinish = async (values: Partial<ISpace>) => {
    Space.create(values)
      .then((space: ISpace) => {
        navigate(`/spaces/${space.key}`);
        notification.success({ message: 'Space created successfully.' });
      })
      .catch((resp: AxiosResponse<IErrorMessage>) => {
        notification.error({
          key: 'create-space-error',
          message: 'Space creation failure.',
          description: map(resp.data.message, (value, key) => value).join('\n')
        });
      })
  }

  return (
    <Layout.Content>
      <div className="container">
        <PageHeader
          ghost={false}
          title="Create Space"
          breadcrumb={{
            items: [
              { title: <Link to="/">Home</Link> },
              { title: <Link to="/spaces">Spaces</Link> },
              { title: 'Create Space' },
            ]
          }}
        />

        <Form<Partial<ISpace>>
          name="space-form"
          form={form}
          preserve={false}
          layout="horizontal"
          labelCol={{ span: 4 }}
          wrapperCol={{ span: 18 }}
          onFinish={handleFormFinish}
        >
          <Form.Item
            name="name"
            label="Space Name"
            rules={[{ required: true }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
            name="key"
            label="Space Key"
            extra="This key cannot be modified after the space is created. It will appear as a unique identifier on the URL of the space."
            rules={[{ required: true }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
            name="description"
            label="Description"
          >
            <Input.TextArea autoSize={{ minRows: 3 }} />
          </Form.Item>

          <Form.Item name="multilingual" label="Enable Multilingual" valuePropName="checked" initialValue={multilingual}>
            <Switch onChange={(value) => setMultilingual(value)} />
          </Form.Item>

          <Form.Item name="lang" label="Default Language" initialValue="en-US">
            <Select disabled={!multilingual} style={{ width: 200 }}>
              <Select.Option value="en-US">English</Select.Option>
              <Select.Option value="zh-CN">简体中文</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item name="fallback_lang" label="Fallback Language" initialValue="en-US">
            <Select disabled={!multilingual} style={{ width: 200 }}>
              <Select.Option value="en-US">English</Select.Option>
              <Select.Option value="zh-CN">简体中文</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item name="status" label="Status" initialValue="offline">
            <Radio.Group>
              <Radio value="offline">Offline</Radio>
              <Radio value="online">Online</Radio>
            </Radio.Group>
          </Form.Item>

          <Form.Item wrapperCol={{ offset: 4, span: 18 }}>
            <AntSpace>
              <Button type="primary" htmlType="submit">Submit</Button>
              <Link to="/spaces"><Button>Cancel</Button></Link>
            </AntSpace>
          </Form.Item>
        </Form>
      </div>
    </Layout.Content>
  );
})

export default NewSpace;
