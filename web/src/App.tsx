import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import { AuthProvider, useAuth } from './stores/auth';
import MainLayout from './layouts/MainLayout';
import LoginPage from './pages/auth/LoginPage';

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const auth = useAuth();
  if (!auth.isAuthenticated) return <Navigate to="/login" replace />;
  return <>{children}</>;
}

function App() {
  return (
    <ConfigProvider locale={zhCN}>
      <AuthProvider>
        <BrowserRouter>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/" element={<ProtectedRoute><MainLayout /></ProtectedRoute>}>
              {/* Placeholder routes — Phase 2 will add real pages */}
              <Route index element={<div style={{textAlign:'center',padding:100}}><h2>欢迎使用招聘平台</h2><p>请从左侧菜单选择功能</p></div>} />
              <Route path="jobs" element={<div>职位列表（建设中）</div>} />
              <Route path="jobs/manage" element={<div>职位管理（建设中）</div>} />
              <Route path="resumes" element={<div>简历管理（建设中）</div>} />
              <Route path="applications" element={<div>投递记录（建设中）</div>} />
              <Route path="candidates" element={<div>候选人管理（建设中）</div>} />
              <Route path="interviews" element={<div>面试管理（建设中）</div>} />
            </Route>
          </Routes>
        </BrowserRouter>
      </AuthProvider>
    </ConfigProvider>
  );
}

export default App;
