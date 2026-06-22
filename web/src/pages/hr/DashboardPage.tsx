import { useState } from 'react';
import { Dropdown, Badge, Table, DatePicker, Button, Avatar } from 'antd';
import {
  CaretDownOutlined, MailOutlined, QuestionCircleOutlined,
  LockOutlined, UserOutlined, PlusCircleOutlined, TeamOutlined,
  SendOutlined, RocketOutlined, RightOutlined, WarningFilled,
  AlertFilled, CaretUpOutlined, RiseOutlined,
  IdcardOutlined, FileTextOutlined, CalendarOutlined, TrophyOutlined,
  SearchOutlined,
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../stores/auth';
import LightSidebar, { hrMenu } from '../../components/LightSidebar';

const { RangePicker } = DatePicker;

// ── Mock Data ──

const statCards = [
  { title: '在招职位', count: 32, change: '+4', icon: <IdcardOutlined />, bg: 'bg-blue-50 text-blue-500' },
  { title: '收到简历', count: 256, change: '+18', icon: <FileTextOutlined />, bg: 'bg-green-50 text-green-500' },
  { title: '面试安排', count: 48, change: '+6', icon: <CalendarOutlined />, bg: 'bg-purple-50 text-purple-500' },
  { title: '待入职', count: 12, change: '+2', icon: <TrophyOutlined />, bg: 'bg-orange-50 text-orange-500' },
];

const top5Columns = [
  { title: '职位名称', dataIndex: 'title', key: 'title', render: (t: string) => <span className="text-sm text-slate-700">{t}</span> },
  { title: '浏览量', dataIndex: 'views', key: 'views', render: (v: number) => <span className="text-sm text-slate-500">{v}</span> },
  { title: '投递量', dataIndex: 'applications', key: 'applications', render: (a: number) => <span className="text-sm text-slate-500">{a}</span> },
  { title: '转化率', dataIndex: 'rate', key: 'rate', render: (r: string) => <span className="text-sm font-semibold text-[#1677FF]">{r}</span> },
];

const top5Data = [
  { key: '1', title: '产品经理', views: 1520, applications: 86, rate: '5.66%' },
  { key: '2', title: 'Java 开发工程师', views: 2341, applications: 102, rate: '4.36%' },
  { key: '3', title: '前端开发工程师', views: 1892, applications: 78, rate: '4.12%' },
  { key: '4', title: 'UI/UX 设计师', views: 987, applications: 35, rate: '3.55%' },
  { key: '5', title: '数据分析师', views: 756, applications: 24, rate: '3.17%' },
];

const quickActions = [
  { icon: <PlusCircleOutlined />, label: '发布职位', color: 'bg-blue-50 text-blue-500' },
  { icon: <SearchOutlined />, label: '搜索简历', color: 'bg-green-50 text-green-500' },
  { icon: <SendOutlined />, label: '邀请面试', color: 'bg-purple-50 text-purple-500' },
  { icon: <RocketOutlined />, label: '人才推荐', color: 'bg-orange-50 text-orange-500' },
];

const todos = [
  { icon: <WarningFilled className="text-amber-500" />, title: '12 份简历待筛选', color: 'bg-amber-50' },
  { icon: <AlertFilled className="text-orange-500" />, title: '5 场面试待安排', color: 'bg-orange-50' },
];

// ── SVG Chart Component ──

function ResumeChart() {
  const w = 560; const h = 200; const pad = { top: 20, right: 20, bottom: 30, left: 40 };
  const days = ['05-20', '05-21', '05-22', '05-23', '05-24', '05-25', '05-26'];
  const received = [42, 58, 35, 72, 65, 88, 95];
  const passed = [18, 25, 15, 38, 32, 45, 52];
  const maxVal = 100;
  const px = (i: number) => pad.left + (i / (days.length - 1)) * (w - pad.left - pad.right);
  const py = (v: number) => pad.top + (1 - v / maxVal) * (h - pad.top - pad.bottom);

  const mkPath = (data: number[]) =>
    data.map((v, i) => `${i === 0 ? 'M' : 'L'}${px(i).toFixed(1)} ${py(v).toFixed(1)}`).join(' ');

  return (
    <svg viewBox={`0 0 ${w} ${h}`} className="w-full h-auto" style={{ maxWidth: 560 }}>
      {/* Grid lines */}
      {[0, 25, 50, 75, 100].map((v) => (
        <g key={v}>
          <line x1={pad.left} y1={py(v)} x2={w - pad.right} y2={py(v)} stroke="#f1f5f9" strokeWidth="1" />
          <text x={pad.left - 8} y={py(v) + 4} textAnchor="end" className="text-[10px]" fill="#94a3b8">{v}</text>
        </g>
      ))}
      {/* X axis labels */}
      {days.map((d, i) => (
        <text key={d} x={px(i)} y={h - 8} textAnchor="middle" className="text-[10px]" fill="#94a3b8">{d}</text>
      ))}
      {/* Received line */}
      <path d={mkPath(received)} fill="none" stroke="#1677FF" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round" />
      {received.map((v, i) => (
        <circle key={`r${i}`} cx={px(i)} cy={py(v)} r="3.5" fill="#fff" stroke="#1677FF" strokeWidth="2" />
      ))}
      {/* Passed line */}
      <path d={mkPath(passed)} fill="none" stroke="#93c5fd" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round" strokeDasharray="6,3" />
      {passed.map((v, i) => (
        <circle key={`p${i}`} cx={px(i)} cy={py(v)} r="3" fill="#fff" stroke="#93c5fd" strokeWidth="2" />
      ))}
    </svg>
  );
}

// ── Component ──

export default function DashboardPage() {
  const [selectedMenu, setSelectedMenu] = useState('home');
  const auth = useAuth();
  const navigate = useNavigate();

  return (
    <div className="flex h-screen bg-white">
      <LightSidebar
        items={hrMenu}
        activeKey={selectedMenu}
        onSelect={(key) => {
          setSelectedMenu(key);
          if (key === 'jobs') navigate('/app/jobs/manage');
          if (key === 'resumes') navigate('/app/candidates');
          if (key === 'interviews') navigate('/app/interviews');
        }}
      />

      <div className="flex-1 flex flex-col overflow-auto">
        {/* ===== Header ===== */}
        <header className="bg-white border-b border-slate-100 flex items-center justify-between px-6 h-18">
          {/* Left */}
          <Dropdown menu={{ items: [
            { key: '1', label: '放心科技有限公司' },
            { key: '2', label: '放心教育分公司' },
          ]}}>
            <div className="flex items-center gap-1.5 cursor-pointer text-slate-700 font-medium text-sm">
              <span>放心科技有限公司</span>
              <CaretDownOutlined className="text-[10px] text-slate-400" />
            </div>
          </Dropdown>

          {/* Right */}
          <div className="flex items-center gap-5">
            <Badge count={12} size="small" offset={[-2, 2]}>
              <MailOutlined className="text-lg text-slate-400 hover:text-[#1677FF] cursor-pointer transition-colors" />
            </Badge>
            <QuestionCircleOutlined className="text-lg text-slate-400 hover:text-[#1677FF] cursor-pointer transition-colors" />
            <LockOutlined className="text-lg text-slate-400 hover:text-[#1677FF] cursor-pointer transition-colors" />

            <Dropdown menu={{ items: [
              { key: 'profile', label: '个人中心', onClick: () => navigate('/app/profile') },
              { key: 'logout', label: '退出登录', onClick: () => { auth.logout(); navigate('/login'); } },
            ] }}>
              <div className="flex items-center gap-2 cursor-pointer pl-4 border-l border-slate-100">
                <Avatar size={32} icon={<UserOutlined />} className="bg-gradient-to-br from-blue-400 to-blue-600" />
                <div className="hidden sm:block text-left leading-tight">
                  <div className="text-sm font-medium text-slate-700">李经理</div>
                  <div className="text-[11px] text-slate-400">招聘者</div>
                </div>
                <CaretDownOutlined className="text-[10px] text-slate-400" />
              </div>
            </Dropdown>
          </div>
        </header>

        {/* ===== Content ===== */}
        <main className="flex-1 p-10">
          <div className="flex flex-col gap-6">
            {/* ── 模块一：数据概览 ── */}
            <div>
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-lg font-bold text-slate-800">数据概览</h2>
                <RangePicker size="small" className="rounded-lg" />
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
                    <div className="text-4xl font-bold text-slate-800 mb-1">{s.count}</div>
                    <div className="flex items-center gap-1">
                      <CaretUpOutlined className="text-green-500 text-[10px]" />
                      <span className="text-xs text-green-500 font-medium">{s.change}</span>
                      <span className="text-[10px] text-slate-400 ml-0.5">较昨日</span>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* ── 模块二：图表 + TOP5 ── */}
            <div className="grid grid-cols-12 gap-6">
              {/* Resume Chart */}
              <div className="col-span-12 lg:col-span-7 bg-white rounded-xl p-6">
                <div className="flex items-center justify-between mb-4">
                  <h3 className="font-semibold text-slate-800">简历数据</h3>
                  <div className="flex items-center gap-4 text-xs">
                    <span className="flex items-center gap-1.5">
                      <span className="w-3 h-0.5 bg-[#1677FF] rounded inline-block" /> 收到简历
                    </span>
                    <span className="flex items-center gap-1.5">
                      <span className="w-3 h-0.5 bg-blue-300 rounded inline-block border-dashed" /> 通过筛选
                    </span>
                  </div>
                </div>
                <div className="flex justify-center mt-2">
                  <ResumeChart />
                </div>
              </div>

              {/* Top 5 Table */}
              <div className="col-span-12 lg:col-span-5 bg-white rounded-xl p-6">
                <div className="flex items-center justify-between mb-4">
                  <h3 className="font-semibold text-slate-800">职位发布效果 TOP5</h3>
                  <RiseOutlined className="text-green-500 text-sm" />
                </div>
                <Table
                  columns={top5Columns}
                  dataSource={top5Data}
                  pagination={false}
                  size="small"
                  showHeader={true}
                  className="[&_.ant-table-thead>tr>th]:!bg-transparent [&_.ant-table-thead>tr>th]:!text-slate-400 [&_.ant-table-thead>tr>th]:!text-xs [&_.ant-table-thead>tr>th]:!border-b [&_.ant-table-tbody>tr>td]:!border-b [&_.ant-table-tbody>tr>td]:!border-slate-50 [&_.ant-table]:!bg-transparent"
                />
              </div>
            </div>

            {/* ── 模块三：快捷操作 + 待办 ── */}
            <div className="grid grid-cols-12 gap-6">
              {/* Quick Actions */}
              <div className="col-span-12 lg:col-span-8 bg-white rounded-xl p-6">
                <h3 className="font-semibold text-slate-800 mb-5">快捷操作</h3>
                <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
                  {quickActions.map((a) => (
                    <button
                      key={a.label}
                      onClick={() => {
                        if (a.label === '发布职位') navigate('/app/jobs/new/edit');
                        if (a.label === '搜索简历') navigate('/app/candidates');
                        if (a.label === '邀请面试') navigate('/app/interviews');
                      }}
                      className="flex flex-col items-center gap-3 p-4 rounded-xl hover:bg-slate-50 transition-colors cursor-pointer border-0 bg-transparent"
                    >
                      <div className={`w-12 h-12 rounded-full flex items-center justify-center text-xl ${a.color}`}>
                        {a.icon}
                      </div>
                      <span className="text-sm text-slate-600 font-medium">{a.label}</span>
                    </button>
                  ))}
                </div>
              </div>

              {/* Todo List */}
              <div className="col-span-12 lg:col-span-4 bg-white rounded-xl p-6">
                <h3 className="font-semibold text-slate-800 mb-5">待办事项</h3>
                <div className="flex flex-col gap-3">
                  {todos.map((todo, idx) => (
                    <div
                      key={idx}
                      className={`flex items-center gap-3 p-4 rounded-xl cursor-pointer hover:shadow-sm transition-all group ${todo.color}`}
                    >
                      <span className="text-lg">{todo.icon}</span>
                      <span className="flex-1 text-sm font-medium text-slate-700">{todo.title}</span>
                      <RightOutlined className="text-xs text-slate-300 group-hover:text-slate-500 group-hover:translate-x-0.5 transition-all" />
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </div>
        </main>
      </div>
    </div>
  );
}
