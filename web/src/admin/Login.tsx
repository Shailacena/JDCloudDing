import { Flex, Button, Card, Form, Input, message } from 'antd';
import bg from '../assets/bg.jpg';
import { useAppContext } from '../AppProvider';
import { useNavigate } from 'react-router-dom';
import type { FormProps } from 'antd';
import { useApis } from '../api/api';
import { useEffect } from 'react';
import { AdminLoginReq } from '../api/types';
import { LockOutlined, UserOutlined, VerifiedOutlined } from '@ant-design/icons';

function Login() {
  let navigate = useNavigate();
  let ctx = useAppContext();
  let { adminLogin } = useApis()

  useEffect(() => {
    if (ctx.cookie.token) {
      goHome();
    }
  }, []);

  const goHome = () => {
    navigate('/admin/home', { replace: true });
  }

  const onFinish: FormProps<AdminLoginReq>['onFinish'] = async (value) => {
    try {
      const { data } = await adminLogin(value)
      ctx.auth.adminSignin(data, value.username, () => {
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
    <>
      <Flex style={{ height: "100%", backgroundImage: `url(${bg})`, backgroundSize: "cover" }} >
        <Card className="login" title={<h2>管理后台</h2>}>
          <Form
            name="basic"
            autoComplete="off"
            onFinish={onFinish}
          >
            <Form.Item<AdminLoginReq>
              name="username"
              rules={[{ required: true, message: '请输入用户名' }]}
            >
              <Input prefix={<UserOutlined />} placeholder="用户名" />
            </Form.Item>

            <Form.Item<AdminLoginReq>
              name="password"
              rules={[{ required: true, message: '请输入密码' }]}
            >
              <Input.Password prefix={<LockOutlined />} type="password" placeholder="密码" />
            </Form.Item>

            <Form.Item<AdminLoginReq>
              name="verifiCode"
              rules={[{ required: true, message: '请输入谷歌验证码' }]}
            >
              <Input prefix={<VerifiedOutlined />} placeholder="谷歌验证码" />
            </Form.Item>

            <Form.Item>
              <Button size="large" block type="primary" htmlType="submit">
                登录
              </Button>
            </Form.Item>
          </Form >
        </Card>
      </Flex>
    </>
  )
}

export default Login
