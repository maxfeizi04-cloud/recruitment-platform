import React from 'react';
import { Layout, Menu, Button, Typography, Avatar } from 'antd';
import {
  SearchOutlined, FileTextOutlined, SendOutlined, CalendarOutlined,
  PlusCircleOutlined, TeamOutlined, LogoutOutlined, UserOutlined,
} from '@ant-design/icons';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../stores/auth';

const { Header, Sider, Content } = Layout;
const { Text } = Typography;

const candidateMenu = [
  { key: '/jobs', icon: <SearchOutlined />, label: '职位搜索' },
  { key: '/resumes', icon: <FileTextOutlined />, label: '我的简历' },
  { key: '/applications', icon: <SendOutlined />, label: '投递记录' },
  { key: '/interviews', icon: <CalendarOutlined />, label: '面试邀约' },
];

const hrMenu = [
  { key: '/jobs', icon: <SearchOutlined />, label: '职位搜索' },
  { key: '/jobs/manage', icon: <PlusCircleOutlined />, label: '职位管理' },
  { key: '/candidates', icon: <TeamOutlined />, label: '候选人管理' },
  { key: '/interviews', icon: <CalendarOutlined />, label: '面试管理' },
];

export default function MainLayout() {
  const auth = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const menuItems = auth.isHR ? hrMenu : candidateMenu;

  const handleLogout = () => {
    auth.logout();
    navigate('/login');
  };

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider breakpoint="lg" collapsedWidth="0">
        <div style={{ padding: '16px', textAlign: 'center' }}>
          <Text strong style={{ color: '#fff', fontSize: 16 }}>招聘平台</Text>
        </div>
        <Menu
          theme="dark" mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>
      <Layout>
        <Header style={{ background: '#fff', padding: '0 24px', display: 'flex', justifyContent: 'flex-end', alignItems: 'center' }}>
          <Avatar icon={<UserOutlined />} style={{ marginRight: 8 }} />
          <Text style={{ marginRight: 16 }}>{auth.user?.name || auth.user?.phone}</Text>
          <Button type="text" icon={<LogoutOutlined />} onClick={handleLogout}>退出</Button>
        </Header>
        <Content style={{ margin: 24, padding: 24, background: '#fff', borderRadius: 8 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
}
