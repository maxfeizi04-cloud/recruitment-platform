import {
  HomeOutlined, FileTextOutlined, SendOutlined, StarOutlined,
  CalendarOutlined, HeartOutlined, ToolOutlined, ExperimentOutlined, SettingOutlined,
  IdcardOutlined, SearchOutlined, TrophyOutlined, BankOutlined, BarChartOutlined,
  UserOutlined, AuditOutlined, FileSearchOutlined, SafetyCertificateOutlined,
  TeamOutlined, ApartmentOutlined, CheckCircleOutlined, FormOutlined, DatabaseOutlined,
  BellOutlined, LockOutlined, QuestionCircleOutlined, GlobalOutlined,
} from '@ant-design/icons';
import type { ReactNode } from 'react';

// ── Type ──

export interface MenuItem {
  key: string;
  label: string;
  icon: ReactNode;
}

// ── Preset Menus ──

export const candidateMenu: MenuItem[] = [
  { key: 'home', label: '首页', icon: <HomeOutlined /> },
  { key: 'resumes', label: '我的简历', icon: <FileTextOutlined /> },
  { key: 'applications', label: '投递记录', icon: <SendOutlined /> },
  { key: 'favorites', label: '我的收藏', icon: <StarOutlined /> },
  { key: 'interviews', label: '面试邀请', icon: <CalendarOutlined /> },
  { key: 'follows', label: '我的关注', icon: <HeartOutlined /> },
  { key: 'services', label: '求职服务', icon: <ToolOutlined /> },
  { key: 'assessment', label: '职业测评', icon: <ExperimentOutlined /> },
  { key: 'settings', label: '账号设置', icon: <SettingOutlined /> },
];

export const hrMenu: MenuItem[] = [
  { key: 'home', label: '首页', icon: <HomeOutlined /> },
  { key: 'jobs', label: '职位管理', icon: <IdcardOutlined /> },
  { key: 'resumes', label: '简历管理', icon: <FileTextOutlined /> },
  { key: 'talents', label: '人才推荐', icon: <SearchOutlined /> },
  { key: 'interviews', label: '面试管理', icon: <CalendarOutlined /> },
  { key: 'offers', label: 'Offer管理', icon: <TrophyOutlined /> },
  { key: 'company', label: '企业信息', icon: <BankOutlined /> },
  { key: 'stats', label: '数据统计', icon: <BarChartOutlined /> },
  { key: 'settings', label: '账号设置', icon: <SettingOutlined /> },
];

export const adminMenu: MenuItem[] = [
  { key: 'home', label: '首页', icon: <HomeOutlined /> },
  { key: 'users', label: '用户管理', icon: <UserOutlined /> },
  { key: 'companies', label: '企业管理', icon: <BankOutlined /> },
  { key: 'jobs', label: '职位管理', icon: <IdcardOutlined /> },
  { key: 'content', label: '内容管理', icon: <FileTextOutlined /> },
  { key: 'audit', label: '审核管理', icon: <AuditOutlined /> },
  { key: 'stats', label: '数据统计', icon: <BarChartOutlined /> },
  { key: 'config', label: '系统配置', icon: <SettingOutlined /> },
  { key: 'logs', label: '日志管理', icon: <FileSearchOutlined /> },
  { key: 'permissions', label: '权限管理', icon: <SafetyCertificateOutlined /> },
];

// ── Props ──

interface Props {
  items: MenuItem[];
  activeKey: string;
  onSelect: (key: string) => void;
  onBack?: () => void;
  showBack?: boolean;
  tagline?: string;
}

// ── Component ──

export default function LightSidebar({ items, activeKey, onSelect, onBack, showBack = false, tagline = '让求职招聘更放心' }: Props) {
  return (
    <aside className="w-64 h-screen bg-white border-r border-gray-100 flex flex-col flex-shrink-0">
      {/* Logo */}
      <div className="p-6">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl bg-[#1677FF] flex items-center justify-center flex-shrink-0">
            <svg viewBox="0 0 24 24" className="w-5 h-5 text-white" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
              <polyline points="20 6 9 17 4 12" />
            </svg>
          </div>
          <div className="leading-tight">
            <span className="text-lg font-bold text-gray-900 tracking-tight">放心</span>
            <p className="text-[10px] text-gray-400 leading-none mt-0.5">{tagline}</p>
          </div>
        </div>
      </div>

      {/* Back button */}
      {showBack && (
        <button
          onClick={onBack}
          className="mx-4 mb-2 flex items-center gap-1.5 text-sm text-gray-400 hover:text-[#1677ff] transition-colors"
        >
          <svg viewBox="0 0 24 24" className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <polyline points="15 18 9 12 15 6" />
          </svg>
          返回
        </button>
      )}

      {/* Menu List */}
      <nav className="flex-1 px-4 space-y-1 overflow-y-auto">
        {items.map((item) => {
          const isActive = activeKey === item.key;
          return (
            <button
              key={item.key}
              onClick={() => onSelect(item.key)}
              className={`w-full flex items-center gap-3 px-4 py-3 rounded-xl text-sm transition-all duration-200 ${
                isActive
                  ? 'bg-[#e6f4ff]/60 text-[#1677ff] font-semibold'
                  : 'text-gray-600 font-medium hover:bg-gray-50 hover:text-[#1677ff]'
              }`}
            >
              <span className="text-lg flex-shrink-0">{item.icon}</span>
              <span>{item.label}</span>
            </button>
          );
        })}
      </nav>
    </aside>
  );
}
