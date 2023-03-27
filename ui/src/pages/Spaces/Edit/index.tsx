import { observer } from "mobx-react-lite";
import { Link, useParams } from "react-router-dom";
import { map } from "lodash";
import { Breadcrumb, Button, Form, Input, notification, Radio, Select, Space as AntSpace } from "antd";
import { PageHeader } from '@ant-design/pro-components';

import { ISpace } from 'models';
import { AxiosResponse, IErrorMessage } from "services";

import { useSpaceContext } from "../Detail/store";

const EditSpace = observer(() => {
  const store = useSpaceContext();
  const { space } = store;

  const { key } = useParams() as { key: string };

  const [form] = Form.useForm();

  const handleFormFinish = async (values: any) => {

    const update = { ...values }

    store.update(key, update)
      .then((space: ISpace) => {
        notification.success({ message: '更新空间成功' });
      })
      .catch((resp: AxiosResponse<IErrorMessage>) => {
        notification.error({
          key: 'update-space-error',
          message: '更新空间失败',
          description: map(resp.data.message, (value, key) => value).join('\n')
        });
      })
  }

  return (
    <>
      <PageHeader
        ghost={false}
        breadcrumb={
          <Breadcrumb>
            <Breadcrumb.Item><Link to={`/spaces/${space.key}`}>空间首页</Link></Breadcrumb.Item>
            <Breadcrumb.Item>空间设置</Breadcrumb.Item>
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
          initialValue={space.name}
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="description"
          label="空间描述"
          initialValue={space.description}
        >
          <Input.TextArea autoSize={{ minRows: 3 }} />
        </Form.Item>

        <Form.Item
          name="status"
          label="空间状态"
          initialValue={space.status}
        >
          <Radio.Group>
            <Radio value="offline">Offline</Radio>
            <Radio value="online">Online</Radio>
          </Radio.Group>
        </Form.Item>

        <Form.Item name="lang" label="默认语言" initialValue={space.lang}>
          <Select style={{ width: 200 }}>
            <Select.Option value="en-US">English</Select.Option>
            <Select.Option value="zh-CN">简体中文</Select.Option>
          </Select>
        </Form.Item>

        <Form.Item name="fallback_lang" label="备用语言" initialValue={space.fallback_lang}>
          <Select style={{ width: 200 }}>
            <Select.Option value="en-US">English</Select.Option>
            <Select.Option value="zh-CN">简体中文</Select.Option>
          </Select>
        </Form.Item>

        <Form.Item wrapperCol={{ offset: 4, span: 18 }}>
          <AntSpace>
            <Button type="primary" htmlType="submit">保存</Button>
            <Link to={`/spaces/${space.key}`}><Button>取消</Button></Link>
          </AntSpace>
        </Form.Item>
      </Form>
    </>
  );
})

export default EditSpace