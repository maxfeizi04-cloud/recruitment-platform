import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import { AuthProvider, useAuth } from './stores/auth';
import MainLayout from './layouts/MainLayout';
import LoginPage from './pages/auth/LoginPage';

// Candidate pages
import CandidateHomePage from './pages/candidate/HomePage';
import JobListPage from './pages/candidate/JobListPage';
import JobDetailPage from './pages/candidate/JobDetailPage';
import ResumeListPage from './pages/candidate/ResumeListPage';
import ApplicationListPage from './pages/candidate/ApplicationListPage';
import InterviewListPage from './pages/candidate/InterviewListPage';

// HR pages
import HRDashboardPage from './pages/hr/DashboardPage';
import JobManagePage from './pages/hr/JobManagePage';
import JobEditPage from './pages/hr/JobEditPage';
import CandidateListPage from './pages/hr/CandidateListPage';

// Admin pages
import AdminDashboardPage from './pages/admin/DashboardPage';

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const auth = useAuth();
  if (!auth.isAuthenticated) return <Navigate to="/login" replace />;
  return <>{children}</>;
}

function HomePage() {
  const auth = useAuth();
  if (auth.user?.role === 'admin') return <AdminDashboardPage />;
  if (auth.isHR) return <HRDashboardPage />;
  return <CandidateHomePage />;
}

function App() {
  return (
    <ConfigProvider locale={zhCN}>
      <AuthProvider>
        <BrowserRouter>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/" element={<ProtectedRoute><HomePage /></ProtectedRoute>} />
            <Route path="/app" element={<ProtectedRoute><MainLayout /></ProtectedRoute>}>
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
