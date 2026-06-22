import { useState } from 'react';
import { Layout, Menu, Badge, Avatar, DatePicker } from 'antd';
import {
  HomeOutlined, UserOutlined, BankOutlined, IdcardOutlined,
  FileTextOutlined, AuditOutlined, BarChartOutlined, SettingOutlined,
  FileSearchOutlined, SafetyCertificateOutlined, BellOutlined,
  FullscreenOutlined, QuestionCircleOutlined, CaretDownOutlined,
  TeamOutlined, ApartmentOutlined, CheckCircleOutlined,
  AlertOutlined, BlockOutlined, RiseOutlined, ArrowUpOutlined,
  AppstoreOutlined, SearchOutlined, FormOutlined, DatabaseOutlined,
  RightOutlined,
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';

const { Header, Sider, Content } = Layout;
const { RangePicker } = DatePicker;

// ── Mock Data ──

const menuItems = [
  { key: 'home', icon: <HomeOutlined />, label: '首页' },
  { key: 'users', icon: <UserOutlined />, label: '用户管理' },
  { key: 'companies', icon: <BankOutlined />, label: '企业管理' },
  { key: 'jobs', icon: <IdcardOutlined />, label: '职位管理' },
  { key: 'content', icon: <FileTextOutlined />, label: '内容管理' },
  { key: 'audit', icon: <AuditOutlined />, label: '审核管理' },
  { key: 'stats', icon: <BarChartOutlined />, label: '数据统计' },
  { key: 'config', icon: <SettingOutlined />, label: '系统配置' },
  { key: 'logs', icon: <FileSearchOutlined />, label: '日志管理' },
  { key: 'permissions', icon: <SafetyCertificateOutlined />, label: '权限管理' },
];

const statCards = [
  { title: '注册用户', count: '156,823', change: '+1,256', icon: <TeamOutlined />, bg: 'bg-blue-50 text-blue-500' },
  { title: '企业用户', count: '23,456', change: '+189', icon: <BankOutlined />, bg: 'bg-green-50 text-green-500' },
  { title: '职位数量', count: '78,901', change: '+432', icon: <IdcardOutlined />, bg: 'bg-purple-50 text-purple-500' },
  { title: '简历数量', count: '312,556', change: '+2,108', icon: <FileTextOutlined />, bg: 'bg-orange-50 text-orange-500' },
];

const alerts = [
  { icon: <AlertOutlined className="text-red-500" />, label: '企业认证待审核', count: '125 家', color: 'text-red-500', bg: 'bg-red-50' },
  { icon: <BlockOutlined className="text-amber-500" />, label: '职位审核待处理', count: '87 条', color: 'text-amber-500', bg: 'bg-amber-50' },
  { icon: <AlertOutlined className="text-orange-500" />, label: '用户举报待处理', count: '23 条', color: 'text-orange-500', bg: 'bg-orange-50' },
];

const quickEntries = [
  { icon: <TeamOutlined />, label: '用户管理', color: 'bg-blue-50 text-blue-500' },
  { icon: <BankOutlined />, label: '企业管理', color: 'bg-green-50 text-green-500' },
  { icon: <CheckCircleOutlined />, label: '职位审核', color: 'bg-purple-50 text-purple-500' },
  { icon: <FormOutlined />, label: '内容审核', color: 'bg-orange-50 text-orange-500' },
  { icon: <DatabaseOutlined />, label: '数据报表', color: 'bg-teal-50 text-teal-500' },
];

// ── SVG Line Chart: User Growth Trend ──
function UserGrowthChart() {
  const w = 560; const h = 200; const pad = { top: 20, right: 20, bottom: 30, left: 50 };
  const labels = ['01-01', '02-01', '03-01', '04-01', '05-01', '06-01'];
  const totalUsers = [89, 105, 118, 132, 143, 156];
  const newUsers  = [5.2, 6.8, 4.9, 8.1, 6.3, 9.5];
  const maxVal = 160;
  const px = (i: number) => pad.left + (i / (labels.length - 1)) * (w - pad.left - pad.right);
  const py = (v: number) => pad.top + (1 - v / maxVal) * (h - pad.top - pad.bottom);

  const mkPath = (data: number[]) =>
    data.map((v, i) => `${i === 0 ? 'M' : 'L'}${px(i).toFixed(1)} ${py(v).toFixed(1)}`).join(' ');

  return (
    <svg viewBox={`0 0 ${w} ${h}`} className="w-full h-auto" style={{ maxWidth: 560 }}>
      {[0, 40, 80, 120, 160].map((v) => (
        <g key={v}>
          <line x1={pad.left} y1={py(v)} x2={w - pad.right} y2={py(v)} stroke="#f1f5f9" strokeWidth="1" />
          <text x={pad.left - 8} y={py(v) + 4} textAnchor="end" className="text-[10px]" fill="#94a3b8">{v}K</text>
        </g>
      ))}
      {labels.map((d, i) => (
        <text key={d} x={px(i)} y={h - 8} textAnchor="middle" className="text-[10px]" fill="#94a3b8">{d}</text>
      ))}
      <path d={mkPath(totalUsers)} fill="none" stroke="#1677FF" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round" />
      {totalUsers.map((v, i) => (
        <circle key={`t${i}`} cx={px(i)} cy={py(v)} r="3.5" fill="#fff" stroke="#1677FF" strokeWidth="2" />
      ))}
      <path d={mkPath(newUsers)} fill="none" stroke="#93c5fd" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" strokeDasharray="5,3" />
      {newUsers.map((v, i) => (
        <circle key={`n${i}`} cx={px(i)} cy={py(v)} r="3" fill="#fff" stroke="#93c5fd" strokeWidth="2" />
      ))}
    </svg>
  );
}

// ── SVG Area Chart: Platform PV ──
function PVAreaChart() {
  const w = 400; const h = 200; const pad = { top: 10, right: 10, bottom: 30, left: 40 };
  const labels = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];
  const data = [28, 42, 35, 58, 45, 68, 55];
  const maxVal = 80;
  const px = (i: number) => pad.left + (i / (labels.length - 1)) * (w - pad.left - pad.right);
  const py = (v: number) => pad.top + (1 - v / maxVal) * (h - pad.top - pad.bottom);

  const linePath = data.map((v, i) => `${i === 0 ? 'M' : 'L'}${px(i).toFixed(1)} ${py(v).toFixed(1)}`).join(' ');
  const areaPath = `${linePath} L${px(data.length - 1).toFixed(1)} ${h - pad.bottom} L${px(0).toFixed(1)} ${h - pad.bottom} Z`;

  return (
    <svg viewBox={`0 0 ${w} ${h}`} className="w-full h-auto" style={{ maxWidth: 400 }}>
      <defs>
        <linearGradient id="pvGrad" x1="0" y1="0" x2="0" y2="1">
          <stop offset="0%" stopColor="#1677FF" stopOpacity="0.3" />
          <stop offset="100%" stopColor="#1677FF" stopOpacity="0.02" />
        </linearGradient>
      </defs>
      {[0, 20, 40, 60, 80].map((v) => (
        <g key={v}>
          <line x1={pad.left} y1={py(v)} x2={w - pad.right} y2={py(v)} stroke="#f1f5f9" strokeWidth="1" />
          <text x={pad.left - 5} y={py(v) + 4} textAnchor="end" className="text-[10px]" fill="#94a3b8">{v}K</text>
        </g>
      ))}
      {labels.map((d, i) => (
        <text key={d} x={px(i)} y={h - 8} textAnchor="middle" className="text-[10px]" fill="#94a3b8">{d}</text>
      ))}
      <path d={areaPath} fill="url(#pvGrad)" />
      <path d={linePath} fill="none" stroke="#1677FF" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
      {data.map((v, i) => (
        <circle key={i} cx={px(i)} cy={py(v)} r="3" fill="#fff" stroke="#1677FF" strokeWidth="2" />
      ))}
    </svg>
  );
}

// ── Component ──

export default function AdminDashboardPage() {
  const [selectedMenu, setSelectedMenu] = useState('home');

  return (
    <Layout className="min-h-screen">
      {/* ===== Sider ===== */}
      <Sider width={260} className="!bg-[#111c24] min-h-screen">
        {/* Logo */}
        <div className="px-5 py-5">
          <div className="flex items-center gap-2.5 mb-2">
            <div className="w-9 h-9 rounded-xl bg-gradient-to-br from-blue-500 to-emerald-400 flex items-center justify-center">
              <CheckCircleOutlined className="text-white text-lg" />
            </div>
            <span className="text-xl font-bold text-white tracking-tight">放心</span>
          </div>
          <p className="text-[10px] text-slate-500 tracking-wider pl-0.5">系统管理平台</p>
        </div>

        {/* Menu */}
        <Menu
          mode="inline"
          theme="dark"
          selectedKeys={[selectedMenu]}
          onClick={({ key }) => setSelectedMenu(key)}
          items={menuItems}
          className="bg-transparent border-0 px-2 [&_.ant-menu-item-selected]:!bg-[#1677FF] [&_.ant-menu-item]:!rounded-lg [&_.ant-menu-item]:!my-0.5"
          style={{ background: 'transparent', color: '#94a3b8' }}
        />
      </Sider>

      <Layout>
        {/* ===== Header ===== */}
        <Header className="bg-white border-b border-slate-100 flex items-center justify-between px-6 h-16">
          <h2 className="text-base font-bold text-slate-700 tracking-wide">管理控制台</h2>

          <div className="flex items-center gap-5">
            <Badge count={10} size="small" offset={[-2, 2]}>
              <BellOutlined className="text-lg text-slate-400 hover:text-[#1677FF] cursor-pointer transition-colors" />
            </Badge>
            <FullscreenOutlined className="text-lg text-slate-400 hover:text-[#1677FF] cursor-pointer transition-colors" />
            <QuestionCircleOutlined className="text-lg text-slate-400 hover:text-[#1677FF] cursor-pointer transition-colors" />

            <div className="flex items-center gap-2 cursor-pointer pl-4 border-l border-slate-100">
              <Avatar size={32} icon={<UserOutlined />} className="bg-gradient-to-br from-slate-600 to-slate-800" />
              <div className="hidden sm:block text-left leading-tight">
                <div className="text-sm font-medium text-slate-700">管理员</div>
                <div className="text-[11px] text-slate-400">系统管理员</div>
              </div>
              <CaretDownOutlined className="text-[10px] text-slate-400" />
            </div>
          </div>
        </Header>

        {/* ===== Content ===== */}
        <Content className="bg-white p-8">
          <div className="flex flex-col gap-6">
            {/* ── 模块一：平台数据概览 ── */}
            <div>
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-lg font-bold text-slate-800">平台数据概览</h2>
                <RangePicker size="small" className="rounded-lg" defaultValue={undefined} placeholder={['开始日期', '结束日期']} />
              </div>
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                {statCards.map((s) => (
                  <div key={s.title} className="bg-white rounded-xl p-5 relative overflow-hidden hover:shadow-sm transition-shadow cursor-pointer">
                    <div className="flex items-center justify-between mb-3">
                      <span className="text-sm text-slate-400">{s.title}</span>
                      <div className={`w-9 h-9 rounded-lg flex items-center justify-center ${s.bg} bg-opacity-20`}>
                        <span className="text-base">{s.icon}</span>
                      </div>
                    </div>
                    <div className="text-4xl font-bold text-slate-800 mb-1 tracking-tight">{s.count}</div>
                    <div className="flex items-center gap-1">
                      <ArrowUpOutlined className="text-green-500 text-[10px]" />
                      <span className="text-xs text-green-500 font-medium">{s.change}</span>
                      <span className="text-[10px] text-slate-400 ml-0.5">较昨日</span>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* ── 模块二：趋势图表监控 ── */}
            <div className="grid grid-cols-12 gap-6">
              {/* User Growth Line Chart */}
              <div className="col-span-12 lg:col-span-7 bg-white rounded-xl p-6">
                <div className="flex items-center justify-between mb-4">
                  <h3 className="font-semibold text-slate-800">用户增长趋势</h3>
                  <div className="flex items-center gap-4 text-xs">
                    <span className="flex items-center gap-1.5">
                      <span className="w-3 h-0.5 bg-[#1677FF] rounded inline-block" /> 用户总数
                    </span>
                    <span className="flex items-center gap-1.5">
                      <span className="w-3 h-0.5 bg-blue-300 inline-block" /> 新增用户
                    </span>
                  </div>
                </div>
                <div className="flex justify-center mt-2">
                  <UserGrowthChart />
                </div>
              </div>

              {/* PV Area Chart */}
              <div className="col-span-12 lg:col-span-5 bg-white rounded-xl p-6">
                <div className="flex items-center justify-between mb-4">
                  <h3 className="font-semibold text-slate-800">平台访问趋势</h3>
                  <div className="flex items-center gap-1.5 text-xs text-slate-400">
                    <span className="w-3 h-0.5 bg-[#1677FF] rounded inline-block" /> 访问量(PV)
                  </div>
                </div>
                <div className="flex justify-center mt-2">
                  <PVAreaChart />
                </div>
              </div>
            </div>

            {/* ── 模块三：系统预警 + 快捷入口 ── */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {/* System Alerts */}
              <div className="bg-white rounded-xl p-6">
                <h3 className="font-semibold text-slate-800 mb-5 flex items-center gap-2">
                  <AppstoreOutlined /> 系统预警
                </h3>
                <div className="flex flex-col gap-3">
                  {alerts.map((alert, idx) => (
                    <div
                      key={idx}
                      className={`flex items-center gap-3 p-4 rounded-xl cursor-pointer hover:shadow-sm transition-all group ${alert.bg}`}
                    >
                      <span className="text-lg">{alert.icon}</span>
                      <span className="flex-1 text-sm font-medium text-slate-700">{alert.label}</span>
                      <span className={`text-sm font-semibold ${alert.color} flex items-center gap-1`}>
                        {alert.count} <RightOutlined className="text-[10px] group-hover:translate-x-0.5 transition-transform" />
                      </span>
                    </div>
                  ))}
                </div>
              </div>

              {/* Quick Entries */}
              <div className="bg-white rounded-xl p-6">
                <h3 className="font-semibold text-slate-800 mb-5 flex items-center gap-2">
                  <RiseOutlined /> 快捷入口
                </h3>
                <div className="flex justify-around items-center h-[calc(100%-2.5rem)]">
                  {quickEntries.map((entry) => (
                    <button
                      key={entry.label}
                      className="flex flex-col items-center gap-3 p-3 rounded-xl hover:bg-slate-50 transition-colors cursor-pointer border-0 bg-transparent"
                    >
                      <div className={`w-12 h-12 rounded-full flex items-center justify-center text-xl ${entry.color}`}>
                        {entry.icon}
                      </div>
                      <span className="text-xs text-slate-500 font-medium">{entry.label}</span>
                    </button>
                  ))}
                </div>
              </div>
            </div>
          </div>
        </Content>
      </Layout>
    </Layout>
  );
}
