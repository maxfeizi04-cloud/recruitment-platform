import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Tabs, Checkbox, message } from 'antd';
import { sendCode, login } from '../../api/auth';
import { useAuth } from '../../stores/auth';
import {
  PhoneOutlined, LockOutlined, EyeOutlined, EyeInvisibleOutlined,
  WechatOutlined, AlipayCircleOutlined, MessageOutlined,
  AuditOutlined, GlobalOutlined, RocketOutlined, SafetyCertificateOutlined,
  ThunderboltOutlined, CaretDownOutlined,
} from '@ant-design/icons';

// ── Static Constants ──

const navLinks = [
  { label: '首页', active: true, hasDropdown: false },
  { label: '找工作', active: false, hasDropdown: true },
  { label: '招人才', active: false, hasDropdown: true },
  { label: '企业服务', active: false, hasDropdown: true },
  { label: '关于我们', active: false, hasDropdown: false },
];

const features = [
  { icon: <AuditOutlined />, title: '严格审核', desc: '企业职位人工审核' },
  { icon: <SafetyCertificateOutlined />, title: '信息安全', desc: '多重保护隐私安全' },
  { icon: <RocketOutlined />, title: '高效匹配', desc: '精准推荐合适职位' },
];

const ribbonItems = [
  { icon: <GlobalOutlined />, title: '海量优质企业', desc: '5000+ 认证企业在线招聘' },
  { icon: <AuditOutlined />, title: '真实职位信息', desc: 'AI 审核 + 人工复核保障' },
  { icon: <ThunderboltOutlined />, title: '精准职位推荐', desc: '智能算法匹配最适合的你' },
  { icon: <SafetyCertificateOutlined />, title: '全程安心服务', desc: '从投递到入职全程保障' },
];

const footerLinks = ['关于我们', '帮助中心', '用户协议', '隐私政策', '联系我们'];

// ── Component ──

export default function LoginPage() {
  const [activeTab, setActiveTab] = useState<'candidate' | 'hr'>('candidate');
  const [phone, setPhone] = useState('');
  const [code, setCode] = useState('');
  const [showPwd, setShowPwd] = useState(false);
  const [remember, setRemember] = useState(false);
  const [loading, setLoading] = useState(false);
  const [counting, setCounting] = useState(0);
  const [error, setError] = useState('');

  const auth = useAuth();
  const navigate = useNavigate();

  const handleSendCode = async () => {
    if (!phone || phone.length !== 11) { setError('请输入正确的手机号'); return; }
    setError('');
    try { await sendCode(phone); setCounting(60);
      const t = setInterval(() => setCounting((c) => { if (c <= 1) { clearInterval(t); return 0; } return c - 1; }), 1000);
    } catch { setError('验证码发送失败'); }
  };

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!phone || !code) { setError('请填写手机号和验证码'); return; }
    setLoading(true); setError('');
    try { const r = await login(phone, code, activeTab); auth.login(r.token, r.user); navigate('/'); }
    catch { setError('登录失败，请检查验证码'); }
    finally { setLoading(false); }
  };

  if (auth.isAuthenticated) { navigate('/'); return null; }

  return (
    <div className="min-h-screen bg-white font-sans overflow-x-hidden">
      {/* ===== Header ===== */}
      <header className="sticky top-0 z-50 bg-white/95 backdrop-blur border-b border-gray-100">
        <div className="max-w-7xl mx-auto px-6 lg:px-8">
          <div className="flex items-center justify-between h-20">
            {/* Logo */}
            <a href="/" className="flex items-center gap-3 flex-shrink-0">
              <div className="w-11 h-11 rounded-xl bg-[#1677FF] flex items-center justify-center">
                <span className="text-white font-bold text-xl">放</span>
              </div>
              <div className="leading-tight">
                <span className="text-lg font-bold text-[#1f1f1f] tracking-tight">放心</span>
                <p className="text-[10px] text-[#8c8c8c] leading-none mt-0.5">让求职招聘更放心</p>
              </div>
            </a>

            {/* Nav */}
            <nav className="hidden lg:flex items-center gap-1">
              {navLinks.map((item) => (
                <a key={item.label} href="#" className={`flex items-center gap-1 px-4 py-2 text-base font-medium rounded-lg transition-colors ${item.active ? 'text-[#1677FF] bg-blue-50/50' : 'text-[#1f1f1f] hover:text-[#1677FF] hover:bg-slate-50'}`}>
                  {item.label}{item.hasDropdown && <CaretDownOutlined className="text-[10px] mt-px" />}
                </a>
              ))}
            </nav>

            {/* Right */}
            <div className="flex items-center gap-3">
              <button className="text-sm font-medium text-[#8c8c8c] hover:text-[#1677FF] transition-colors px-3 py-2">企业登录</button>
              <button className="text-sm font-medium text-[#1677FF] border border-[#1677FF] rounded-lg px-4 py-2 hover:bg-blue-50 transition-colors">注册账号</button>
            </div>
          </div>
        </div>
      </header>

      {/* ===== Hero Main Content ===== */}
      <main className="max-w-7xl mx-auto px-6 lg:px-8 py-20 lg:py-32 min-h-[calc(100vh-80px)] flex items-center">
        <div className="grid grid-cols-12 gap-8 items-center">
          {/* Left: Brand (col-span-3) */}
          <div className="col-span-12 lg:col-span-3">
            <h1 className="text-5xl lg:text-6xl font-extrabold text-[#1f1f1f] leading-tight tracking-tight">
              放心找工作
            </h1>
            <h2 className="text-5xl lg:text-6xl font-extrabold text-[#1677FF] leading-tight tracking-tight mt-1">
              安心好未来
            </h2>
            <p className="mt-4 text-lg text-[#8c8c8c] tracking-wide">真实 · 安全 · 高效 · 可靠</p>

            <div className="mt-10 space-y-6">
              {features.map((f) => (
                <div key={f.title} className="flex items-start gap-4 group cursor-pointer">
                  <div className="w-12 h-12 rounded-xl bg-blue-50 flex items-center justify-center flex-shrink-0 group-hover:bg-[#1677FF] transition-colors duration-300">
                    <span className="text-lg text-[#1677FF] group-hover:text-white transition-colors duration-300">{f.icon}</span>
                  </div>
                  <div>
                    <h3 className="font-semibold text-[#1f1f1f] text-sm">{f.title}</h3>
                    <p className="mt-0.5 text-xs text-[#8c8c8c]">{f.desc}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Center: 3D Illustration Area (col-span-5) */}
          <div className="hidden lg:block col-span-5">
            <div className="relative h-[520px] rounded-3xl bg-[radial-gradient(circle_at_40%_40%,#e6f4ff_0%,#f4f7fc_70%)] overflow-hidden">
              {/* City Skyline SVG Silhouette */}
              <div className="absolute bottom-0 left-0 right-0 opacity-20">
                <svg viewBox="0 0 600 200" className="w-full" preserveAspectRatio="none">
                  <rect x="20" y="100" width="40" height="100" fill="#1677FF" rx="2" />
                  <rect x="70" y="60" width="50" height="140" fill="#1677FF" rx="2" />
                  <rect x="130" y="80" width="35" height="120" fill="#1677FF" rx="2" />
                  <rect x="175" y="40" width="45" height="160" fill="#1677FF" rx="2" />
                  <rect x="230" y="90" width="55" height="110" fill="#1677FF" rx="2" />
                  <rect x="295" y="50" width="40" height="150" fill="#1677FF" rx="2" />
                  <rect x="345" y="70" width="30" height="130" fill="#1677FF" rx="2" />
                  <rect x="385" y="100" width="50" height="100" fill="#1677FF" rx="2" />
                  <rect x="445" y="55" width="45" height="145" fill="#1677FF" rx="2" />
                  <rect x="500" y="85" width="40" height="115" fill="#1677FF" rx="2" />
                  <rect x="550" y="65" width="50" height="135" fill="#1677FF" rx="2" />
                </svg>
              </div>

              {/* 3D Office Scene */}
              <div className="absolute inset-0 flex items-center justify-center">
                <div className="relative w-72 h-72 group cursor-pointer transition-transform duration-500 hover:scale-105">
                  {/* Backdrop glow */}
                  <div className="absolute inset-0 rounded-full bg-blue-200/30 blur-3xl" />

                  {/* Monitor */}
                  <div className="absolute top-8 left-1/2 -translate-x-1/2 w-48 h-32 bg-white/80 backdrop-blur-md border border-white/60 shadow-xl rounded-2xl flex items-center justify-center">
                    <div className="w-40 h-24 bg-gradient-to-br from-blue-50 to-white rounded-lg border border-blue-100 flex items-center justify-center">
                      <svg className="w-16 h-12" viewBox="0 0 64 48"><rect x="8" y="4" width="48" height="36" rx="4" fill="#e6f4ff" stroke="#1677FF" strokeWidth="1.5"/><rect x="24" y="44" width="16" height="3" rx="1" fill="#1677FF"/><rect x="20" y="47" width="24" height="2" rx="1" fill="#94a3b8"/><circle cx="32" cy="20" r="6" fill="#1677FF" opacity="0.3"/></svg>
                    </div>
                  </div>

                  {/* Office Chair */}
                  <div className="absolute bottom-12 left-8 w-20 h-24">
                    <div className="absolute bottom-0 left-1/2 -translate-x-1/2 w-14 h-3 bg-gradient-to-r from-blue-400 to-blue-600 rounded-full shadow-lg" />
                    <div className="absolute bottom-3 left-1/2 -translate-x-1/2 w-2.5 h-20 bg-slate-400 rounded-t-full" />
                    <div className="absolute bottom-16 left-1/2 -translate-x-1/2 w-16 h-10 bg-white/80 backdrop-blur-md border border-white/60 rounded-2xl shadow-md" />
                  </div>

                  {/* Briefcase */}
                  <div className="absolute bottom-16 right-8 w-16 h-12 bg-gradient-to-br from-blue-400 to-blue-600 rounded-xl shadow-lg flex items-center justify-center">
                    <div className="w-8 h-1.5 bg-blue-300 rounded-full" />
                  </div>

                  {/* Plant */}
                  <div className="absolute bottom-8 right-12 w-10 h-12">
                    <div className="absolute bottom-0 left-1/2 -translate-x-1/2 w-5 h-6 bg-amber-200 rounded-b-lg" />
                    <div className="absolute bottom-4 left-1/2 -translate-x-1/2 w-8 h-8 bg-green-400 rounded-full shadow-sm" />
                    <div className="absolute bottom-3 left-1/2 -translate-x-1/2 w-7 h-7 bg-green-300 rounded-full" />
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Right: Login Card (col-span-4) */}
          <div className="col-span-12 lg:col-span-4">
            <div className="bg-white/90 backdrop-blur rounded-3xl shadow-xl shadow-blue-500/5 border border-slate-100 p-10">
              <div className="text-center mb-6">
                <h2 className="text-3xl font-bold text-[#1f1f1f]">欢迎登录放心</h2>
                <p className="mt-2 text-sm text-[#8c8c8c]">
                  还没有账号？<a href="#" className="text-[#1677FF] hover:underline font-medium">立即注册</a>
                </p>
              </div>

              {/* Tabs */}
              <Tabs
                activeKey={activeTab}
                onChange={(k) => setActiveTab(k as 'candidate' | 'hr')}
                centered
                items={[
                  { key: 'candidate', label: '求职者登录' },
                  { key: 'hr', label: '企业登录' },
                ]}
                className="[&_.ant-tabs-ink-bar]:!bg-[#1677FF] [&_.ant-tabs-tab-active_.ant-tabs-tab-btn]:!text-[#1677FF] [&_.ant-tabs-tab]:!text-[#8c8c8c]"
              />

              {/* Form */}
              <form onSubmit={handleLogin} className="space-y-5 mt-2">
                {error && (
                  <div className="bg-red-50 border border-red-200 text-red-600 text-sm rounded-xl px-4 py-3">{error}</div>
                )}

                {/* Phone */}
                <div>
                  <label className="block text-sm font-medium text-[#1f1f1f] mb-1.5">手机号</label>
                  <div className="relative">
                    <span className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-300"><PhoneOutlined /></span>
                    <input type="tel" maxLength={11} value={phone} onChange={(e) => setPhone(e.target.value)}
                      placeholder="请输入手机号"
                      className="w-full pl-11 pr-4 py-3.5 text-sm border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#1677FF]/20 focus:border-[#1677FF] transition-all duration-200 bg-slate-50/50" />
                  </div>
                </div>

                {/* Code */}
                <div>
                  <label className="block text-sm font-medium text-[#1f1f1f] mb-1.5">验证码</label>
                  <div className="relative">
                    <span className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-300"><LockOutlined /></span>
                    <input type={showPwd ? 'text' : 'password'} maxLength={6} value={code} onChange={(e) => setCode(e.target.value)}
                      placeholder="请输入验证码"
                      className="w-full pl-11 pr-24 py-3.5 text-sm border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#1677FF]/20 focus:border-[#1677FF] transition-all duration-200 bg-slate-50/50" />
                    <button type="button" onClick={() => setShowPwd(!showPwd)}
                      className="absolute right-[88px] top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 p-1">
                      {showPwd ? <EyeInvisibleOutlined /> : <EyeOutlined />}
                    </button>
                    <button type="button" onClick={handleSendCode} disabled={counting > 0}
                      className="absolute right-3 top-1/2 -translate-y-1/2 text-xs font-semibold text-[#1677FF] hover:text-blue-700 disabled:text-slate-300 disabled:cursor-not-allowed whitespace-nowrap">
                      {counting > 0 ? `${counting}s后重发` : '获取验证码'}
                    </button>
                  </div>
                </div>

                {/* Remember + Forgot */}
                <div className="flex items-center justify-between">
                  <Checkbox checked={remember} onChange={(e) => setRemember(e.target.checked)}
                    className="text-sm text-[#8c8c8c] [&_.ant-checkbox-checked_.ant-checkbox-inner]:!bg-[#1677FF] [&_.ant-checkbox-checked_.ant-checkbox-inner]:!border-[#1677FF]">
                    记住我
                  </Checkbox>
                  <a href="#" className="text-sm text-[#1677FF] hover:underline">忘记密码？</a>
                </div>

                {/* Submit Button */}
                <button type="submit" disabled={loading}
                  className="w-full py-3.5 bg-[#1677FF] hover:bg-blue-600 active:bg-blue-700 text-white text-base font-semibold rounded-xl transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed shadow-lg shadow-blue-500/25 hover:shadow-blue-500/40">
                  {loading ? '登录中...' : '登 录'}
                </button>
              </form>

              {/* Divider + Social */}
              <div className="relative my-6">
                <div className="absolute inset-0 flex items-center"><div className="w-full border-t border-slate-100"></div></div>
                <div className="relative flex justify-center"><span className="bg-white px-4 text-xs text-slate-300 tracking-wider">其他登录方式</span></div>
              </div>

              <div className="flex justify-center gap-5">
                {[
                  { icon: <WechatOutlined />, hover: 'hover:bg-green-50 hover:text-green-500 hover:border-green-300', label: '微信登录' },
                  { icon: <AlipayCircleOutlined />, hover: 'hover:bg-blue-50 hover:text-blue-500 hover:border-blue-300', label: '支付宝登录' },
                  { icon: <MessageOutlined />, hover: 'hover:bg-orange-50 hover:text-orange-500 hover:border-orange-300', label: '短信登录' },
                ].map((s) => (
                  <button key={s.label} title={s.label}
                    className={`w-12 h-12 rounded-xl border border-slate-200 flex items-center justify-center text-xl text-slate-300 transition-all duration-200 hover:scale-110 ${s.hover}`}>
                    {s.icon}
                  </button>
                ))}
              </div>
            </div>
          </div>
        </div>
      </main>

      {/* ===== Ribbon Bar ===== */}
      <section className="bg-gray-50/60 border-y border-gray-100">
        <div className="mx-auto px-6 lg:px-8 py-16">
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
            {ribbonItems.map((item) => (
              <div key={item.title} className="flex items-start gap-4">
                <div className="w-14 h-14 rounded-xl bg-blue-50 flex items-center justify-center flex-shrink-0">
                  <span className="text-2xl text-[#1677FF]">{item.icon}</span>
                </div>
                <div>
                  <h3 className="font-semibold text-[#1f1f1f] text-sm">{item.title}</h3>
                  <p className="mt-1 text-xs text-[#8c8c8c]">{item.desc}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ===== Footer ===== */}
      <footer className="bg-white">
        <div className="mx-auto px-6 lg:px-8 py-14">
          <div className="flex flex-col md:flex-row items-center justify-between gap-6">
            {/* Logo */}
            <div className="flex items-center gap-2.5">
              <div className="w-9 h-9 rounded-lg bg-[#1677FF] flex items-center justify-center">
                <span className="text-white font-bold text-sm">放</span>
              </div>
              <div className="leading-tight">
                <span className="text-sm font-semibold text-[#1f1f1f]">放心</span>
                <p className="text-[10px] text-[#8c8c8c] leading-none mt-0.5">让求职招聘更放心</p>
              </div>
            </div>

            {/* Links */}
            <nav className="flex flex-wrap justify-center gap-6">
              {footerLinks.map((item) => (
                <a key={item} href="#" className="text-sm text-[#8c8c8c] hover:text-[#1677FF] transition-colors">{item}</a>
              ))}
            </nav>

            {/* Copyright */}
            <p className="text-sm text-[#8c8c8c] text-center">© 2026 放心 All Rights Reserved. 湘ICP备2023001234号-1</p>
          </div>
        </div>
      </footer>
    </div>
  );
}
