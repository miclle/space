import React from 'react';
import { observer } from 'mobx-react-lite';
import { Link, useNavigate } from 'react-router-dom';
import { map } from 'lodash';
import { Breadcrumb, Button, Form, Input, Layout, notification, Radio, Select, Space as AntSpace } from 'antd';
import { PageHeader } from '@ant-design/pro-components';

import { ISpace } from 'models';
import { AxiosResponse, Space, IErrorMessage } from "services";

const NewSpace = observer(() => {
  const navigate = useNavigate();

  const [form] = Form.useForm();

  const handleFormFinish = async (values: Partial<ISpace>) => {
    Space.create(values)
      .then((space: ISpace) => {
        navigate(`/spaces/${space.key}`);
        notification.success({ message: '创建空间成功' });
      })
      .catch((resp: AxiosResponse<IErrorMessage>) => {
        notification.error({
          key: 'create-space-error',
          message: '创建空间失败',
          description: map(resp.data.message, (value, key) => value).join('\n')
        });
      })
  }

  return (
    <Layout.Content>
      <div className="container">
        <PageHeader
          ghost={false}
          title="添加空间"
          breadcrumb={
            <Breadcrumb>
              <Breadcrumb.Item><Link to="/">首页</Link></Breadcrumb.Item>
              <Breadcrumb.Item><Link to="/spaces">空间列表</Link></Breadcrumb.Item>
              <Breadcrumb.Item>添加空间</Breadcrumb.Item>
            </Breadcrumb>
          }
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
            label="空间名称"
            rules={[{ required: true, message: '空间名称不能为空!' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
            name="key"
            label="空间键值 Key"
            extra="创建空间后，不能修改此键。它将作为唯一标识符出现在空间的 URL 上。"
            rules={[{ required: true, message: '空间键值 Key 不能为空!' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
            name="description"
            label="空间描述"
          >
            <Input.TextArea autoSize={{ minRows: 3 }} />
          </Form.Item>

          <Form.Item name="lang" label="默认语言" initialValue="en-US">
            <Select style={{ width: 200 }}>
              <Select.Option value="en-US">English</Select.Option>
              <Select.Option value="zh-CN">简体中文</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item name="fallback_lang" label="备用语言" initialValue="en-US">
            <Select style={{ width: 200 }}>
              <Select.Option value="en-US">English</Select.Option>
              <Select.Option value="zh-CN">简体中文</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item
            name="status"
            label="空间状态"
            initialValue="offline"
          >
            <Radio.Group>
              <Radio value="offline">Offline</Radio>
              <Radio value="online">Online</Radio>
            </Radio.Group>
          </Form.Item>

          <Form.Item wrapperCol={{ offset: 4, span: 18 }}>
            <AntSpace>
              <Button type="primary" htmlType="submit">提交</Button>
              <Link to="/spaces"><Button>取消</Button></Link>
            </AntSpace>
          </Form.Item>
        </Form>
      </div>
    </Layout.Content>
  );
})

export default NewSpace;
