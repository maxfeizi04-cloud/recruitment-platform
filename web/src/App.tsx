import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import { AuthProvider, useAuth } from './stores/auth';
import MainLayout from './layouts/MainLayout';
import LoginPage from './pages/auth/LoginPage';

// Candidate pages
import JobListPage from './pages/candidate/JobListPage';
import JobDetailPage from './pages/candidate/JobDetailPage';
import ResumeListPage from './pages/candidate/ResumeListPage';
import ApplicationListPage from './pages/candidate/ApplicationListPage';
import InterviewListPage from './pages/candidate/InterviewListPage';

// HR pages
import JobManagePage from './pages/hr/JobManagePage';
import JobEditPage from './pages/hr/JobEditPage';
import CandidateListPage from './pages/hr/CandidateListPage';

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const auth = useAuth();
  if (!auth.isAuthenticated) return <Navigate to="/login" replace />;
  return <>{children}</>;
}

function HomePage() {
  const auth = useAuth();
  return (
    <div style={{ textAlign: 'center', padding: 100 }}>
      <h2>欢迎使用招聘平台</h2>
      <p>{auth.isHR ? '从左侧菜单管理职位和候选人' : '从左侧菜单搜索职位并投递简历'}</p>
    </div>
  );
}

function App() {
  return (
    <ConfigProvider locale={zhCN}>
      <AuthProvider>
        <BrowserRouter>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/" element={<ProtectedRoute><MainLayout /></ProtectedRoute>}>
              <Route index element={<HomePage />} />
              <Route path="jobs" element={<JobListPage />} />
              <Route path="jobs/:id" element={<JobDetailPage />} />
              <Route path="jobs/manage" element={<JobManagePage />} />
              <Route path="jobs/:id/edit" element={<JobEditPage />} />
              <Route path="resumes" element={<ResumeListPage />} />
              <Route path="applications" element={<ApplicationListPage />} />
              <Route path="candidates" element={<CandidateListPage />} />
              <Route path="interviews" element={<InterviewListPage />} />
            </Route>
          </Routes>
        </BrowserRouter>
      </AuthProvider>
    </ConfigProvider>
  );
}

export default App;
