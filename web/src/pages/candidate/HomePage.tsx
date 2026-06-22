import { useState } from 'react';
import { Layout, Menu, Progress, Tag, Input, Badge, Button, Tabs } from 'antd';
import {
  SearchOutlined, BellOutlined, MessageOutlined, UserOutlined,
  SendOutlined, CalendarOutlined, StarOutlined, EyeOutlined,
  FileTextOutlined, HeartOutlined, ToolOutlined, ExperimentOutlined,
  SettingOutlined, HomeOutlined, TrophyOutlined, ReloadOutlined,
  EnvironmentOutlined, DollarOutlined, BulbOutlined,
  SafetyCertificateOutlined, RocketOutlined, TeamOutlined,
  RightOutlined, CaretUpOutlined,
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';

const { Header, Sider, Content } = Layout;

// ── Mock Data ──

const statsCards = [
  { title: '投递记录', count: 12, change: '+3', icon: <SendOutlined />, color: 'bg-blue-50 text-blue-500' },
  { title: '面试邀请', count: 5, change: '+2', icon: <CalendarOutlined />, color: 'bg-green-50 text-green-500' },
  { title: '收藏职位', count: 28, change: '+5', icon: <StarOutlined />, color: 'bg-orange-50 text-orange-500' },
];

const recommendedJobs = [
  {
    id: '1', title: '高级前端工程师', salary: '20K-35K',
    company: '字节跳动', logo: '字', city: '北京·朝阳区',
    tags: ['React', 'TypeScript', 'Node.js', 'Webpack'],
  },
  {
    id: '2', title: '全栈开发工程师', salary: '25K-40K',
    company: '腾讯科技', logo: '腾', city: '深圳·南山区',
    tags: ['Go', 'React', 'PostgreSQL', 'K8s'],
  },
  {
    id: '3', title: '前端架构师', salary: '30K-50K',
    company: '阿里巴巴', logo: '阿', city: '杭州·余杭区',
    tags: ['React', 'Micro-FE', 'TypeScript', 'Vite'],
  },
];

const recommendedCompanies = [
  { id: '1', name: '字节跳动', logo: '字', city: '北京', jobs: 156, type: '互联网' },
  { id: '2', name: '腾讯科技', logo: '腾', city: '深圳', jobs: 203, type: '互联网' },
  { id: '3', name: '华为技术', logo: '华', city: '上海', jobs: 89, type: '通信' },
];

const tips = [
  { icon: <FileTextOutlined />, bg: 'bg-blue-50 text-blue-500', title: '如何写好一份技术简历？', reads: '12.3k 阅读' },
  { icon: <BulbOutlined />, bg: 'bg-amber-50 text-amber-500', title: '面试前必须准备的 10 个问题', reads: '8.6k 阅读' },
  { icon: <RocketOutlined />, bg: 'bg-green-50 text-green-500', title: '如何谈薪资？资深 HR 教你这几招', reads: '15.2k 阅读' },
  { icon: <SafetyCertificateOutlined />, bg: 'bg-purple-50 text-purple-500', title: '签劳动合同要注意哪些坑', reads: '9.1k 阅读' },
  { icon: <TeamOutlined />, bg: 'bg-rose-50 text-rose-500', title: '技术面常见算法题汇总', reads: '22.7k 阅读' },
];

// ── Component ──

export default function HomePage() {
  const [selectedMenu, setSelectedMenu] = useState('home');
  const [activeTab, setActiveTab] = useState('jobs');
  const navigate = useNavigate();

  const menuItems = [
    { key: 'home', icon: <HomeOutlined />, label: '首页' },
    { key: 'resumes', icon: <FileTextOutlined />, label: '我的简历' },
    { key: 'applications', icon: <SendOutlined />, label: '投递记录' },
    { key: 'favorites', icon: <StarOutlined />, label: '我的收藏' },
    { key: 'interviews', icon: <CalendarOutlined />, label: '面试邀请' },
    { key: 'follows', icon: <EyeOutlined />, label: '我的关注' },
    { key: 'services', icon: <ToolOutlined />, label: '求职服务' },
    { key: 'assessment', icon: <ExperimentOutlined />, label: '职业测评' },
    { key: 'settings', icon: <SettingOutlined />, label: '账号设置' },
  ];

  const tabItems = [
    {
      key: 'jobs',
      label: '推荐职位',
      children: (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {recommendedJobs.map((job) => (
            <div
              key={job.id}
              className="border border-slate-100 rounded-xl p-4 hover:shadow-md hover:border-blue-100 transition-all duration-200 cursor-pointer"
              onClick={() => navigate(`/jobs/${job.id}`)}
            >
              <div className="flex justify-between items-start mb-3">
                <span className="font-semibold text-slate-800 text-sm">{job.title}</span>
                <span className="text-orange-500 font-medium text-sm whitespace-nowrap">{job.salary}</span>
              </div>
              <div className="flex items-center gap-2 text-xs text-slate-400 mb-3">
                <div className="w-6 h-6 rounded-md bg-[#1677FF] flex items-center justify-center text-white text-[10px] font-bold">{job.logo}</div>
                <span>{job.company}</span>
                <EnvironmentOutlined className="text-[10px]" />
                <span>{job.city}</span>
              </div>
              <div className="flex flex-wrap gap-1.5">
                {job.tags.map((tag) => (
                  <Tag key={tag} className="text-[10px] px-1.5 py-0 m-0 bg-slate-50 border-0 text-slate-500">{tag}</Tag>
                ))}
              </div>
            </div>
          ))}
        </div>
      ),
    },
    {
      key: 'companies',
      label: '推荐企业',
      children: (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {recommendedCompanies.map((c) => (
            <div key={c.id} className="border border-slate-100 rounded-xl p-4 hover:shadow-md hover:border-blue-100 transition-all duration-200 cursor-pointer">
              <div className="flex items-center gap-3 mb-3">
                <div className="w-10 h-10 rounded-xl bg-[#1677FF] flex items-center justify-center text-white font-bold">{c.logo}</div>
                <div>
                  <div className="font-semibold text-slate-800 text-sm">{c.name}</div>
                  <div className="text-xs text-slate-400">{c.type}</div>
                </div>
              </div>
              <div className="flex gap-4 text-xs text-slate-500">
                <span><EnvironmentOutlined className="mr-0.5" />{c.city}</span>
                <span>在招 {c.jobs} 个职位</span>
              </div>
            </div>
          ))}
        </div>
      ),
    },
  ];

  return (
    <Layout className="min-h-screen">
      {/* ===== Header ===== */}
      <Header className="bg-white border-b border-slate-100 flex items-center justify-between px-6 h-18 sticky top-0 z-50">
        {/* Left */}
        <div className="flex items-center gap-8">
          <a href="/" className="flex items-center gap-2 flex-shrink-0">
            <div className="w-10 h-10 rounded-lg bg-[#1677FF] flex items-center justify-center">
              <span className="text-white font-bold text-base">放</span>
            </div>
            <span className="text-lg font-bold text-slate-800 tracking-tight">放心</span>
          </a>
          <nav className="hidden lg:flex items-center gap-6 text-sm">
            <span className="text-[#1677FF] font-medium cursor-pointer">首页</span>
            <span className="text-slate-500 hover:text-[#1677FF] cursor-pointer transition-colors">找工作</span>
            <span className="text-slate-500 hover:text-[#1677FF] cursor-pointer transition-colors">企业</span>
          </nav>
        </div>

        {/* Right */}
        <div className="flex items-center gap-5">
          <Input.Search placeholder="搜索职位/公司" prefix={<SearchOutlined className="text-slate-300" />} className="w-72 [&_.ant-input]:bg-slate-50 [&_.ant-input]:border-slate-200 [&_.ant-input]:rounded-lg" size="middle" />
          <div className="flex items-center gap-4 text-lg text-slate-400">
            <MessageOutlined className="hover:text-[#1677FF] cursor-pointer transition-colors" />
            <Badge dot offset={[-2, 2]}>
              <BellOutlined className="hover:text-[#1677FF] cursor-pointer transition-colors" />
            </Badge>
          </div>
          <div className="flex items-center gap-2 cursor-pointer pl-4 border-l border-slate-100">
            <div className="w-8 h-8 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center">
              <UserOutlined className="text-white text-xs" />
            </div>
            <span className="text-sm text-slate-700 font-medium">你好，张一一</span>
          </div>
        </div>
      </Header>

      <Layout>
        {/* ===== Sider ===== */}
        <Sider width={240} className="bg-white border-r border-slate-100 pt-4">
          <Menu
            mode="inline"
            selectedKeys={[selectedMenu]}
            onClick={({ key }) => {
              setSelectedMenu(key);
              if (key === 'home') navigate('/');
              else if (key === 'resumes') navigate('/app/resumes');
              else if (key === 'applications') navigate('/app/applications');
              else if (key === 'interviews') navigate('/app/interviews');
              else if (key === 'favorites') navigate('/app/jobs');
            }}
            items={menuItems}
            className="border-0 px-2"
            style={{ background: 'transparent' }}
          />
        </Sider>

        {/* ===== Content ===== */}
        <Content className="bg-[#f4f7fc] p-8">
          <div className="flex flex-col gap-6">
            {/* Top Row: Welcome + Stats */}
            <div className="grid grid-cols-12 gap-6">
              {/* Welcome Card */}
              <div className="col-span-12 lg:col-span-7 bg-white rounded-xl p-6 flex flex-col justify-between">
                <div>
                  <h1 className="text-2xl font-bold text-slate-800">
                    下午好，张一一 <span className="inline-block animate-bounce">👋</span>
                  </h1>
                  <p className="text-sm text-slate-400 mt-2">完善简历可以提高求职成功率哦</p>
                </div>
                <div className="flex items-end justify-between mt-6">
                  <div className="flex-1 max-w-sm">
                    <div className="flex justify-between text-sm mb-1">
                      <span className="text-slate-500">简历完整度</span>
                      <span className="text-[#1677FF] font-medium">80%</span>
                    </div>
                    <Progress percent={80} showInfo={false} strokeColor="#1677FF" trailColor="#f1f5f9" size="small" />
                  </div>
                  <Button type="primary" className="bg-[#1677FF] border-0 rounded-lg font-medium h-9">
                    完善简历
                  </Button>
                </div>
              </div>

              {/* Stats Cards */}
              <div className="col-span-12 lg:col-span-5 grid grid-cols-3 gap-4">
                {statsCards.map((s) => (
                  <div key={s.title} className="bg-white rounded-xl p-4 relative overflow-hidden flex flex-col justify-between hover:shadow-sm transition-shadow cursor-pointer">
                    <span className="text-xs text-slate-400">{s.title}</span>
                    <span className="text-4xl font-bold text-slate-800 mt-1">{s.count}</span>
                    <div className="flex items-center gap-1 mt-1">
                      <CaretUpOutlined className="text-green-500 text-[10px]" />
                      <span className="text-xs text-green-500 font-medium">{s.change}</span>
                      <span className="text-[10px] text-slate-400 ml-0.5">较昨日</span>
                    </div>
                    <div className={`absolute -right-3 -bottom-3 w-20 h-20 rounded-full flex items-center justify-center opacity-15 ${s.color}`}>
                      <span className="text-3xl">{s.icon}</span>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Bottom Row: Recommendations + Tips */}
            <div className="grid grid-cols-12 gap-6">
              {/* Recommendations */}
              <div className="col-span-12 lg:col-span-8 bg-white rounded-xl p-6">
                <div className="flex items-center justify-between mb-5">
                  <Tabs
                    activeKey={activeTab}
                    onChange={setActiveTab}
                    items={tabItems.map((t) => ({ key: t.key, label: t.label }))}
                    className="[&_.ant-tabs-nav]:mb-0 [&_.ant-tabs-tab]:pb-3"
                  />
                  <Button type="text" icon={<ReloadOutlined />} className="text-slate-400 hover:text-[#1677FF] text-sm">换一批</Button>
                </div>
                {tabItems.find((t) => t.key === activeTab)?.children}
              </div>

              {/* Tips */}
              <div className="col-span-12 lg:col-span-4 bg-white rounded-xl p-6">
                <div className="flex items-center justify-between mb-5">
                  <h3 className="font-semibold text-slate-800">求职小贴士</h3>
                  <span className="text-xs text-slate-400 cursor-pointer hover:text-[#1677FF] transition-colors">更多 <RightOutlined className="text-[10px]" /></span>
                </div>
                <div className="flex flex-col gap-5">
                  {tips.map((tip, idx) => (
                    <div key={idx} className="flex items-start gap-3 cursor-pointer group">
                      <div className={`w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0 ${tip.bg}`}>
                        <span className="text-lg">{tip.icon}</span>
                      </div>
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium text-slate-700 group-hover:text-[#1677FF] transition-colors truncate">{tip.title}</p>
                        <p className="text-xs text-slate-400 mt-0.5">{tip.reads}</p>
                      </div>
                    </div>
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
