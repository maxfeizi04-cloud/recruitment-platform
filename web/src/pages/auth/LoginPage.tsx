import React, { useState } from 'react';
import { Card, Form, Input, Button, Radio, message, Typography } from 'antd';
import { MobileOutlined, LockOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { sendCode, login } from '../../api/auth';
import { useAuth } from '../../stores/auth';

const { Title } = Typography;

export default function LoginPage() {
  const [loading, setLoading] = useState(false);
  const [counting, setCounting] = useState(0);
  const auth = useAuth();
  const navigate = useNavigate();
  const [form] = Form.useForm();

  const handleSendCode = async () => {
    const phone = form.getFieldValue('phone');
    if (!phone || phone.length !== 11) {
      message.error('请输入正确的手机号');
      return;
    }
    try {
      await sendCode(phone);
      message.success('验证码已发送');
      setCounting(60);
      const timer = setInterval(() => {
        setCounting((c) => {
          if (c <= 1) { clearInterval(timer); return 0; }
          return c - 1;
        });
      }, 1000);
    } catch { /* error handled by interceptor */ }
  };

  const handleLogin = async (values: any) => {
    setLoading(true);
    try {
      const result = await login(values.phone, values.code, values.role);
      auth.login(result.token, result.user);
      message.success('登录成功');
      navigate('/');
    } catch { /* error handled by interceptor */ }
    finally { setLoading(false); }
  };

  if (auth.isAuthenticated) {
    navigate('/');
    return null;
  }

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh', background: '#f0f2f5' }}>
      <Card style={{ width: 400 }}>
        <Title level={3} style={{ textAlign: 'center' }}>招聘平台</Title>
        <Form form={form} onFinish={handleLogin} initialValues={{ role: 'candidate' }}>
          <Form.Item name="phone" rules={[{ required: true, message: '请输入手机号' }, { len: 11, message: '手机号为11位' }]}>
            <Input prefix={<MobileOutlined />} placeholder="手机号" maxLength={11} />
          </Form.Item>
          <Form.Item name="code" rules={[{ required: true, message: '请输入验证码' }, { len: 6, message: '验证码为6位' }]}>
            <Input prefix={<LockOutlined />} placeholder="验证码" maxLength={6}
              suffix={<Button type="link" size="small" disabled={counting > 0} onClick={handleSendCode}>{counting > 0 ? `${counting}s` : '获取验证码'}</Button>}
            />
          </Form.Item>
          <Form.Item name="role">
            <Radio.Group>
              <Radio.Button value="candidate">我是求职者</Radio.Button>
              <Radio.Button value="hr">我是招聘者</Radio.Button>
            </Radio.Group>
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} block>登录 / 注册</Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
}
