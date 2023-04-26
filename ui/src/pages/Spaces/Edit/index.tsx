import { useState } from "react";
import { observer } from "mobx-react-lite";
import { Link, useParams } from "react-router-dom";
import { map } from "lodash";
import { Button, Form, Input, notification, Radio, Select, Space as AntSpace, Switch } from "antd";
import { PageHeader } from '@ant-design/pro-components';

import { ISpace } from 'models';
import { AxiosResponse, IErrorMessage } from "services";

import { useSpaceContext } from "../Detail/store";

const EditSpace = observer(() => {
  const store = useSpaceContext();
  const { space } = store;

  const { key } = useParams() as { key: string };

  const [form] = Form.useForm();
  const [multilingual, setMultilingual] = useState(space.multilingual);

  const handleFormFinish = async (values: any) => {

    const update = { ...values }

    store.update(key, update)
      .then((space: ISpace) => {
        notification.success({ message: 'Update space success.' });
      })
      .catch((resp: AxiosResponse<IErrorMessage>) => {
        notification.error({
          key: 'update-space-error',
          message: 'Update space failure.',
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
            { title: <Link to={`/spaces/${space.key}`}>Space</Link> },
            { title: 'Settings' },
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
          initialValue={space.name}
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="description"
          label="Description"
          initialValue={space.description}
        >
          <Input.TextArea autoSize={{ minRows: 3 }} />
        </Form.Item>

        <Form.Item name="multilingual" label="Enable Multilingual" valuePropName="checked" initialValue={space.multilingual}>
          <Switch onChange={(value) => setMultilingual(value)} />
        </Form.Item>

        <Form.Item name="lang" label="Default Language" initialValue={space.lang}>
          <Select disabled={!multilingual} style={{ width: 200 }}>
            <Select.Option value="en-US">English</Select.Option>
            <Select.Option value="zh-CN">简体中文</Select.Option>
          </Select>
        </Form.Item>

        <Form.Item name="fallback_lang" label="Fallback Language" initialValue={space.fallback_lang}>
          <Select disabled={!multilingual} style={{ width: 200 }}>
            <Select.Option value="en-US">English</Select.Option>
            <Select.Option value="zh-CN">简体中文</Select.Option>
          </Select>
        </Form.Item>

        <Form.Item name="status" label="Status" initialValue={space.status}>
          <Radio.Group>
            <Radio value="offline">Offline</Radio>
            <Radio value="online">Online</Radio>
          </Radio.Group>
        </Form.Item>

        <Form.Item wrapperCol={{ offset: 4, span: 18 }}>
          <AntSpace>
            <Button type="primary" htmlType="submit">Submit</Button>
            <Link to={`/spaces/${space.key}`}><Button>Cancel</Button></Link>
          </AntSpace>
        </Form.Item>
      </Form>
    </>
  );
})

export default EditSpace