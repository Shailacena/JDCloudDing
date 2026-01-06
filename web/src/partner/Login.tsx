import React, { useEffect } from 'react';
import { LockOutlined, UserOutlined, VerifiedOutlined } from '@ant-design/icons';
import { Button, Form, Input, message, FormProps } from 'antd';
import { useNavigate } from 'react-router-dom';
import { useAppContext } from '../AppProvider';
import { getRouteConfig } from './RouteConfigs';
import { useApis } from '../api/api';
import { PartnerLoginReq } from '../api/types';

const Login: React.FC = () => {

  let navigate = useNavigate();
  let ctx = useAppContext();
  let { partner1Login: partnerLogin } = useApis()

  useEffect(() => {
    if (ctx.cookie.token) {
      goHome();
    }
  }, [ctx.cookie]);

  const goHome = () => {
    navigate(getRouteConfig()[0].path, { replace: true });
  }

  const onFinish: FormProps<PartnerLoginReq>['onFinish'] = async (value) => {
    try {
      const { data } = await partnerLogin(value)
      ctx.auth.partnerSignin(data, value.username, () => {
        message.success('登录成功');
        setTimeout(() => {
          goHome();
        }, 500);
      });
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <div style={{ maxWidth: 300, margin: '0 auto', paddingTop: 100 }}>
      <h2 style={{ textAlign: 'center' }}>合作商管理后台</h2>
      <Form
        name="login"
        style={{ maxWidth: 360 }}
        onFinish={onFinish}
      >
        <Form.Item
          name="username"
          rules={[{ required: true, message: '请输入帐号' }]}
        >
          <Input prefix={<UserOutlined />} placeholder="帐号" />
        </Form.Item>
        <Form.Item
          name="password"
          rules={[{ required: true, message: '请输入密码' }]}
        >
          <Input prefix={<LockOutlined />} type="password" placeholder="密码" />
        </Form.Item>

        <Form.Item<PartnerLoginReq>
          name="verifiCode"
          rules={[{ required: true, message: '请输入谷歌验证码' }]}
        >
          <Input prefix={<VerifiedOutlined />} placeholder="谷歌验证码" />
        </Form.Item>

        <Form.Item>
          <Button block type="primary" htmlType="submit">
            登陆
          </Button>
        </Form.Item>
      </Form>
    </div>
  );
};

export default Login;