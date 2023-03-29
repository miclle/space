import { useState } from "react";
import { observer } from "mobx-react-lite";
import { Link, useNavigate } from "react-router-dom";
import { AxiosResponse } from "axios";
import { map } from "lodash";
import { Alert, Button, Form, Input, Layout, notification, Typography } from "antd";
import { AiOutlineLock, AiOutlineUser } from 'react-icons/ai';
import { Account } from "services";
import { IErrorMessage } from "services";

const Signin = observer(() => {

  const navigate = useNavigate();
  const [alert, setAlert] = useState('');

  const onFinish = (values: any) => {
    Account.signin(values.login, values.password)
      .then(() => {
        notification.success({ message: '登录成功' });
        navigate('/');
      })
      .catch((resp: AxiosResponse<IErrorMessage>) => {
        setAlert(map(resp.data.message, (value, key) => value).join('\n'))
      })
  };

  return (
    <Layout className="layout-session">
      <Layout.Content>
        <div className="session-section-wrapper">

          <Typography.Title level={1} style={{ textAlign: 'center' }}>Sign in to Space</Typography.Title>

          {
            alert &&
            <Alert type="error" closable style={{ marginBottom: 20 }} message={alert} onClose={() => setAlert('')} />
          }

          <Form
            name="signin"
            size="large"
            layout="vertical"
            initialValues={{ subscribe: true }}
            onFinish={onFinish}
          >
            <Form.Item
              name="login"
              label="Username / Email"
              rules={[{ required: true, message: 'Please input your login username/email!' }]}
              extra="You can sign in using a login username or an authenticated email address"
            >
              <Input prefix={<AiOutlineUser />} placeholder="eg. thomas" />
            </Form.Item>

            <Form.Item name="password" label="Password" rules={[{ required: true, message: 'Please input your login password!' }]}>
              <Input.Password prefix={<AiOutlineLock />} autoComplete="off" />
            </Form.Item>

            <Form.Item>
              <Button type="primary" htmlType="submit" block>Sign in</Button>
            </Form.Item>

            <Typography.Paragraph>
              New to Space? <Link to="/signup">Create an account</Link>.
            </Typography.Paragraph>
          </Form>
        </div>
      </Layout.Content>
    </Layout>
  );
})

export default Signin
