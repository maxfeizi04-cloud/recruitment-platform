import { useState } from 'react';
import { Avatar, Progress, Timeline, Tag, Dropdown, Button } from 'antd';
import {
  EditOutlined, EnvironmentOutlined, DollarOutlined,
  PhoneOutlined, MailOutlined, CalendarOutlined, HomeOutlined,
  BookOutlined, DownloadOutlined, EyeOutlined, MoreOutlined,
  TrophyOutlined, ProjectOutlined, FilePdfOutlined,
  UserOutlined, StarOutlined,
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../stores/auth';
import LightSidebar, { candidateMenu } from '../../components/LightSidebar';

// ── Mock Data ──

const traitTags = ['沟通能力强', '逻辑思维清晰', '团队协作', '用户体验', '数据分析'];

const basicInfo = [
  { icon: <PhoneOutlined />, label: '手机号', value: '138****8888' },
  { icon: <MailOutlined />, label: '邮箱', value: 'zhangyiyi@email.com' },
  { icon: <CalendarOutlined />, label: '出生年月', value: '1996.08' },
  { icon: <HomeOutlined />, label: '现居地', value: '浙江省杭州市' },
  { icon: <BookOutlined />, label: '最高学历', value: '本科' },
];

const workExperiences = [
  {
    company: '阿里巴巴集团',
    role: '高级产品经理',
    period: '2021.04 - 至今',
    color: '#1677FF',
    duties: [
      '负责电商中台商品域产品规划与落地，管理 15 人产品团队',
      '主导商品推荐引擎重构，转化率提升 23%，年度 GMV 增长 35 亿',
      '推动跨 BU 协作，建立数据驱动产品迭代机制，需求交付周期缩短 40%',
    ],
  },
  {
    company: '网易网络',
    role: '产品助理 → 产品经理',
    period: '2019.06 - 2021.03',
    color: '#94a3b8',
    duties: [
      '参与网易云课堂 B 端产品设计，服务 200+ 企业客户',
      '独立负责用户增长模块，通过 A/B 实验优化注册漏斗，转化率提升 18%',
      '荣获公司年度优秀新人奖',
    ],
  },
];

const skills = [
  { name: '产品设计', level: 95, label: '精通' },
  { name: '数据分析', level: 88, label: '熟练' },
  { name: '用户研究', level: 85, label: '熟练' },
  { name: '项目管理', level: 82, label: '熟练' },
  { name: 'Axure RP', level: 90, label: '精通' },
  { name: 'SQL/Python', level: 70, label: '良好' },
  { name: 'Figma', level: 88, label: '熟练' },
  { name: 'A/B Test', level: 80, label: '熟练' },
];

const certificates = [
  { title: 'PMP 项目管理认证', year: '2023.06', icon: <TrophyOutlined />, bg: 'bg-amber-50 text-amber-500' },
  { title: 'Google UX Design', year: '2022.12', icon: <StarOutlined />, bg: 'bg-blue-50 text-blue-500' },
  { title: 'AWS Cloud Practitioner', year: '2021.08', icon: <ProjectOutlined />, bg: 'bg-orange-50 text-orange-500' },
];

const projects = [
  {
    title: '智能学习平台 0-1 建设',
    role: '产品负责人',
    period: '2022.03 - 2022.09',
    result: '从零搭建智能学习平台，上线首月 DAU 突破 10 万，获公司年度最佳项目奖。',
  },
  {
    title: '电商推荐系统升级',
    role: '核心产品经理',
    period: '2023.01 - 2023.06',
    result: '引入多目标排序模型，人均点击率提升 31%，营收贡献超 2 亿元。',
  },
];

// ── Component ──

export default function ProfilePage() {
  const [selectedMenu, setSelectedMenu] = useState('resume');
  const auth = useAuth();
  const navigate = useNavigate();

  return (
    <div className="flex h-screen bg-white">
      <LightSidebar
        items={candidateMenu}
        activeKey={selectedMenu}
        onSelect={(key) => setSelectedMenu(key)}
        showBack={true}
        onBack={() => navigate('/')}
      />

      <div className="flex-1 flex flex-col overflow-auto">
        {/* ===== Header ===== */}
        <header className="bg-white border-b border-slate-100 flex items-center justify-between px-6 h-16 sticky top-0 z-50">
          <div className="flex items-center gap-8">
            <nav className="hidden md:flex items-center gap-6 text-sm">
              <span className="text-[#1677FF] font-medium cursor-pointer">首页</span>
              <span className="text-slate-500 hover:text-[#1677FF] cursor-pointer transition-colors">找工作</span>
              <span className="text-slate-500 hover:text-[#1677FF] cursor-pointer transition-colors">企业</span>
            </nav>
          </div>

          <Dropdown menu={{ items: [
            { key: 'profile', label: '个人中心', onClick: () => navigate('/app/profile') },
            { key: 'logout', label: '退出登录', onClick: () => { auth.logout(); navigate('/login'); } },
          ] }}>
            <div className="flex items-center gap-2 cursor-pointer">
              <div className="w-8 h-8 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center">
                <UserOutlined className="text-white text-xs" />
              </div>
              <span className="text-sm text-slate-700 font-medium hidden sm:inline">张一一</span>
            </div>
          </Dropdown>
        </header>

        {/* ===== Content ===== */}
        <main className="flex-1 p-8">
          <div className="grid grid-cols-12 gap-6">
            {/* ── Left: Main Resume Column (8 cols) ── */}
            <div className="col-span-12 lg:col-span-8 flex flex-col gap-6">
              {/* Card 1: Profile Hero */}
              <div className="bg-white rounded-xl shadow-sm p-6 relative">
                <div className="absolute top-6 right-6">
                  <button className="text-[#1677FF] text-sm font-medium flex items-center gap-1 hover:text-blue-700 transition-colors">
                    <EditOutlined /> 编辑简历
                  </button>
                </div>

                <div className="flex items-start gap-5">
                  {/* Avatar */}
                  <div className="relative flex-shrink-0">
                    <Avatar size={88} icon={<UserOutlined />} className="bg-gradient-to-br from-blue-400 to-purple-500" />
                    <span className="absolute -bottom-1 -right-1 w-6 h-6 rounded-full border-2 border-white bg-green-400 flex items-center justify-center">
                      <span className="w-2.5 h-2.5 rounded-full bg-white" />
                    </span>
                  </div>

                  {/* Info */}
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-3 mb-1">
                      <h1 className="text-2xl font-bold text-[#1f1f1f]">张一一</h1>
                      <span className="text-base text-slate-400">♀</span>
                      <Tag className="!text-[11px] !px-2 !py-0 !text-green-500 !bg-green-50 !border-green-100 !rounded-full">
                        在职 · 考虑机会
                      </Tag>
                    </div>
                    <div className="text-sm text-[#595959] mb-3">
                      产品经理 · 5年工作经验
                    </div>
                    <div className="flex items-center gap-6 text-sm text-[#8c8c8c] mb-4">
                      <span className="flex items-center gap-1"><EnvironmentOutlined /> 期望城市：杭州市</span>
                      <span className="flex items-center gap-1"><DollarOutlined /> 期望薪资：20K-30K</span>
                    </div>
                    <div className="flex flex-wrap gap-2">
                      {traitTags.map((tag) => (
                        <Tag key={tag} className="!text-xs !px-3 !py-0.5 !bg-slate-50 !border-slate-100 !text-slate-500 !rounded-md">
                          {tag}
                        </Tag>
                      ))}
                    </div>
                  </div>
                </div>
              </div>

              {/* Card 2: Basic Info */}
              <div className="bg-white rounded-xl shadow-sm p-6">
                <h2 className="text-base font-bold text-[#1f1f1f] mb-4">基本信息</h2>
                <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
                  {basicInfo.map((item) => (
                    <div key={item.label} className="flex flex-col gap-1">
                      <span className="text-xs text-[#8c8c8c] flex items-center gap-1">
                        <span className="text-[10px]">{item.icon}</span> {item.label}
                      </span>
                      <span className="text-sm font-medium text-[#1f1f1f]">{item.value}</span>
                    </div>
                  ))}
                </div>
              </div>

              {/* Card 3: Personal Summary */}
              <div className="bg-white rounded-xl shadow-sm p-6">
                <h2 className="text-base font-bold text-[#1f1f1f] mb-3">个人简介</h2>
                <p className="text-sm text-[#595959] leading-relaxed">
                  拥有 5 年互联网产品经理经验，擅长 B 端与 C 端产品设计，具备从 0 到 1 的产品搭建能力和大规模用户产品的迭代优化经验。
                  熟悉电商、教育行业，对用户增长、数据驱动决策有深入理解。具备优秀的跨部门沟通能力和团队管理经验，曾带领 15 人产品团队完成多个核心项目。
                  追求卓越的产品体验，致力于用技术驱动商业价值增长。
                </p>
              </div>

              {/* Card 4: Work Experience */}
              <div className="bg-white rounded-xl shadow-sm p-6">
                <h2 className="text-base font-bold text-[#1f1f1f] mb-5">工作经历</h2>
                <Timeline
                  items={workExperiences.map((exp) => ({
                    color: exp.color,
                    children: (
                      <div key={exp.company}>
                        <div className="flex items-center gap-3 mb-1 flex-wrap">
                          <span className="font-semibold text-[#1f1f1f] text-sm">{exp.company}</span>
                          <Tag className="!text-[11px] !px-2 !py-0 !bg-blue-50 !text-blue-600 !border-blue-100 !rounded-md">
                            {exp.role}
                          </Tag>
                          <span className="text-xs text-[#8c8c8c]">{exp.period}</span>
                        </div>
                        <ul className="mt-2 space-y-1 text-sm text-[#595959] list-disc list-inside">
                          {exp.duties.map((d, i) => (
                            <li key={i} className="leading-relaxed">{d}</li>
                          ))}
                        </ul>
                      </div>
                    ),
                  }))}
                />
              </div>

              {/* Card 5: Education */}
              <div className="bg-white rounded-xl shadow-sm p-6">
                <h2 className="text-base font-bold text-[#1f1f1f] mb-4">教育经历</h2>
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-xl bg-blue-50 flex items-center justify-center flex-shrink-0">
                    <BookOutlined className="text-[#1677FF]" />
                  </div>
                  <div>
                    <div className="font-semibold text-sm text-[#1f1f1f]">浙江大学</div>
                    <div className="text-xs text-[#8c8c8c] mt-0.5">本科 · 计算机科学与技术 · 2015.09 - 2019.06</div>
                  </div>
                </div>
              </div>
            </div>

            {/* ── Right: Sidebar Column (4 cols) ── */}
            <div className="col-span-12 lg:col-span-4 flex flex-col gap-6">
              {/* Action Buttons */}
              <div className="flex gap-2">
                <Button type="primary" icon={<EyeOutlined />} className="flex-1 bg-[#1677FF] border-0 rounded-lg h-10 font-medium">
                  预览简历
                </Button>
                <Button icon={<DownloadOutlined />} className="flex-1 border-[#1677FF] text-[#1677FF] rounded-lg h-10 font-medium">
                  下载简历
                </Button>
                <Dropdown menu={{ items: [{ key: '1', label: '刷新简历' }, { key: '2', label: '设为公开' }, { key: '3', label: '复制链接' }] }}>
                  <Button icon={<MoreOutlined />} className="w-10 border-slate-200 rounded-lg h-10 flex items-center justify-center" />
                </Dropdown>
              </div>

              {/* Card A: Skills */}
              <div className="bg-white rounded-xl shadow-sm p-6">
                <h2 className="text-base font-bold text-[#1f1f1f] mb-4">专业技能</h2>
                <div className="grid grid-cols-2 gap-x-6 gap-y-4">
                  {skills.map((skill) => (
                    <div key={skill.name}>
                      <div className="flex items-center justify-between mb-1.5">
                        <span className="text-sm font-medium text-[#1f1f1f]">{skill.name}</span>
                        <span className="text-[11px] text-[#8c8c8c]">{skill.label}</span>
                      </div>
                      <Progress percent={skill.level} showInfo={false} strokeColor="#1677FF" trailColor="#f1f5f9" size="small" />
                    </div>
                  ))}
                </div>
              </div>

              {/* Card B: Certificates */}
              <div className="bg-white rounded-xl shadow-sm p-6">
                <h2 className="text-base font-bold text-[#1f1f1f] mb-4">证书荣誉</h2>
                <div className="flex flex-col gap-4">
                  {certificates.map((cert) => (
                    <div key={cert.title} className="flex items-center gap-3">
                      <div className={`w-9 h-9 rounded-lg flex items-center justify-center flex-shrink-0 ${cert.bg}`}>
                        <span className="text-sm">{cert.icon}</span>
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="text-sm font-medium text-[#1f1f1f] truncate">{cert.title}</div>
                        <div className="text-xs text-[#8c8c8c] mt-0.5">{cert.year}</div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* Card C: Projects */}
              <div className="bg-white rounded-xl shadow-sm p-6">
                <h2 className="text-base font-bold text-[#1f1f1f] mb-4">项目经历</h2>
                <div className="flex flex-col gap-4">
                  {projects.map((proj) => (
                    <div key={proj.title} className="border border-slate-100 rounded-lg p-4 hover:border-blue-100 transition-colors">
                      <div className="flex items-center gap-2 mb-1 flex-wrap">
                        <span className="text-sm font-semibold text-[#1f1f1f]">{proj.title}</span>
                        <Tag className="!text-[10px] !px-1.5 !py-0 !bg-slate-50 !border-slate-100 !text-slate-500 !rounded">
                          {proj.role}
                        </Tag>
                      </div>
                      <div className="text-[11px] text-[#8c8c8c] mb-2">{proj.period}</div>
                      <p className="text-xs text-[#595959] leading-relaxed">{proj.result}</p>
                    </div>
                  ))}
                </div>
              </div>

              {/* Card D: Resume Attachment */}
              <div className="bg-white rounded-xl shadow-sm p-6">
                <h2 className="text-base font-bold text-[#1f1f1f] mb-4">简历附件</h2>
                <div className="bg-slate-50 rounded-lg p-4 flex items-center justify-between">
                  <div className="flex items-center gap-3 min-w-0">
                    <div className="w-9 h-9 rounded-lg bg-red-50 flex items-center justify-center flex-shrink-0">
                      <FilePdfOutlined className="text-red-500 text-lg" />
                    </div>
                    <div className="min-w-0">
                      <div className="text-sm font-medium text-[#1f1f1f] truncate">张一一_产品经理_简历.pdf</div>
                      <div className="text-xs text-[#8c8c8c] mt-0.5">2.4 MB</div>
                    </div>
                  </div>
                  <button className="text-[#1677FF] text-sm font-medium flex items-center gap-1 hover:text-blue-700 transition-colors flex-shrink-0 ml-3">
                    <DownloadOutlined /> 下载
                  </button>
                </div>
              </div>
            </div>
          </div>
        </main>
      </div>
    </div>
  );
}
