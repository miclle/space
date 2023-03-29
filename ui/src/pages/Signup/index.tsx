import React, { useState } from "react";
import { observer } from "mobx-react-lite";
import { useNavigate } from "react-router-dom";
import { Alert, Button, Checkbox, Form, Input, Layout, notification, Typography } from "antd";
import { AiOutlineLock, AiOutlineMail, AiOutlineUser } from 'react-icons/ai';
import { Account } from 'services';
import { AxiosResponse } from "axios";
import { IErrorMessage } from "services";
import { map } from "lodash";
import { IAccount } from "models";

const Signup = observer(() => {

  const navigate = useNavigate();
  const [alert, setAlert] = useState('');

  const onFinish = (values: Partial<IAccount>) => {
    Account.signup(values)
      .then(() => {
        notification.success({ message: '注册成功' });
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

          <Typography.Title level={1} style={{ textAlign: 'center' }}>Join Space today</Typography.Title>

          {
            alert &&
            <Alert type="error" closable style={{ marginBottom: 20 }} message={alert} onClose={() => setAlert('')} />
          }

          <Form
            name="signup"
            size="large"
            layout="vertical"
            initialValues={{ subscribe: true }}
            onFinish={onFinish}
          >
            <Form.Item name="name" label="Name" rules={[{ required: true, message: 'Please input your name!' }]}>
              <Input prefix={<AiOutlineUser />} autoFocus placeholder="eg. Thomas" />
            </Form.Item>

            <Form.Item name="login" label="Username" rules={[{ required: true, message: 'Please input your username!' }]}>
              <Input prefix={<AiOutlineUser />} placeholder="eg. thomas" />
            </Form.Item>

            <Form.Item name="email" label="Email" rules={[{ required: true, message: 'Please input your email!' }]} extra="Users retrieve accounts and subscribe to updates">
              <Input type="email" prefix={<AiOutlineMail />} placeholder="eg. email@domain.com" />
            </Form.Item>

            <Form.Item name="password" label="Password" rules={[{ required: true, message: 'Please input your login password!' }]}>
              <Input.Password prefix={<AiOutlineLock />} autoComplete="off" />
            </Form.Item>

            <Form.Item name="subscribe" valuePropName="checked">
              <Checkbox>Send me product updates and announcements.</Checkbox>
            </Form.Item>

            <Form.Item>
              <Button type="primary" htmlType="submit" block>Sign up</Button>
            </Form.Item>

            <Typography.Paragraph>
              By registering, you agree to these <a href="/terms-of-service">Terms of Service</a> and <a href="/privacy-policy">Privacy Policy</a>, including the use of <a href="/cookie-use">cookies</a>.
              Others can reach you by the email or phone number provided.
              If you do not agree with the terms, do not submit for registration!
            </Typography.Paragraph>
          </Form>
        </div>
      </Layout.Content>
    </Layout>
  );
})

export default Signup
