import { Button, Typography, Avatar } from 'antd';
import {
  LogoutOutlined, UserOutlined, BankOutlined,
} from '@ant-design/icons';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../stores/auth';
import LightSidebar from '../components/LightSidebar';
import type { MenuItem } from '../components/LightSidebar';

const { Text } = Typography;

const candidateMenuItems: MenuItem[] = [
  { key: '/jobs', label: '职位搜索', icon: <UserOutlined /> },
  { key: '/resumes', label: '我的简历', icon: <UserOutlined /> },
  { key: '/applications', label: '投递记录', icon: <UserOutlined /> },
  { key: '/interviews', label: '面试邀约', icon: <UserOutlined /> },
];

const hrMenuItems: MenuItem[] = [
  { key: '/jobs', label: '职位搜索', icon: <UserOutlined /> },
  { key: '/jobs/manage', label: '职位管理', icon: <UserOutlined /> },
  { key: '/candidates', label: '候选人管理', icon: <UserOutlined /> },
  { key: '/interviews', label: '面试管理', icon: <UserOutlined /> },
  { key: '/app/company', label: '企业信息', icon: <BankOutlined /> },
];

export default function MainLayout() {
  const auth = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const menuItems = auth.isHR ? hrMenuItems : candidateMenuItems;

  const handleLogout = () => {
    auth.logout();
    navigate('/login');
  };

  return (
    <div className="flex h-screen bg-white">
      <LightSidebar
        items={menuItems}
        activeKey={location.pathname}
        onSelect={(key) => navigate(key)}
      />

      <div className="flex-1 flex flex-col overflow-auto">
        <header style={{ background: '#fff', padding: '0 24px', display: 'flex', justifyContent: 'flex-end', alignItems: 'center', height: 64, borderBottom: '1px solid #f1f5f9' }}>
          <Avatar icon={<UserOutlined />} style={{ marginRight: 8 }} />
          <Text style={{ marginRight: 16 }}>{auth.user?.name || auth.user?.phone}</Text>
          <Button type="text" icon={<LogoutOutlined />} onClick={handleLogout}>退出</Button>
        </header>
        <main style={{ margin: 32, padding: 24, background: '#fff', borderRadius: 8 }}>
          <Outlet />
        </main>
      </div>
    </div>
  );
}
