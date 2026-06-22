import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { sendCode, login } from '../../api/auth';
import { useAuth } from '../../stores/auth';
import {
  PhoneOutlined, LockOutlined, EyeOutlined, EyeInvisibleOutlined,
  WechatOutlined, AlipayCircleOutlined, MessageOutlined,
  AuditOutlined, GlobalOutlined, RocketOutlined, SafetyCertificateOutlined,
  ThunderboltOutlined,
} from '@ant-design/icons';

export default function LoginPage() {
  const [activeTab, setActiveTab] = useState<'candidate' | 'hr' | 'admin'>('candidate');
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
    try {
      await sendCode(phone);
      setCounting(60);
      const timer = setInterval(() => setCounting((c) => {
        if (c <= 1) { clearInterval(timer); return 0; }
        return c - 1;
      }), 1000);
    } catch { setError('验证码发送失败'); }
  };

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!phone || !code) { setError('请填写手机号和验证码'); return; }
    setLoading(true); setError('');
    try {
      const result = await login(phone, code, activeTab);
      auth.login(result.token, result.user);
      navigate('/');
    } catch { setError('登录失败，请检查验证码'); }
    finally { setLoading(false); }
  };

  if (auth.isAuthenticated) { navigate('/'); return null; }

  return (
    <div className="min-h-screen bg-slate-50 font-sans">
      {/* ===== Header ===== */}
      <header className="sticky top-0 z-50 bg-white/95 backdrop-blur border-b border-slate-100">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-18">
            <a href="/" className="flex items-center gap-2 flex-shrink-0">
              <div className="w-10 h-10 rounded-lg bg-[#1677FF] flex items-center justify-center">
                <span className="text-white font-bold text-lg">放</span>
              </div>
              <span className="text-xl font-bold text-slate-800 tracking-tight">放心</span>
            </a>

            <nav className="hidden md:flex items-center gap-8">
              {['首页', '找工作', '招人才', '企业服务', '关于我们'].map((item) => (
                <a key={item} href="#" className="text-slate-600 hover:text-[#1677FF] transition-colors duration-200 text-base font-medium">{item}</a>
              ))}
            </nav>

            <div className="flex items-center gap-4">
              <button className="hidden sm:inline-flex items-center px-4 py-2 text-sm font-medium text-[#1677FF] border border-[#1677FF] rounded-lg hover:bg-blue-50 transition-colors duration-200">企业登录</button>
              <a href="#" className="text-sm text-slate-500 hover:text-[#1677FF] transition-colors duration-200">注册账号</a>
            </div>
          </div>
        </div>
      </header>

      {/* ===== Hero Section ===== */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20 lg:py-28">
        <div className="flex flex-col lg:flex-row items-center gap-12 lg:gap-16">
          {/* Left: Promo */}
          <div className="flex-1 lg:w-[60%] max-w-2xl">
            <h1 className="text-5xl sm:text-6xl lg:text-7xl font-extrabold text-slate-900 leading-tight tracking-tight">
              放心找工作 <span className="text-[#1677FF]">安心好未来</span>
            </h1>
            <p className="mt-4 text-lg text-slate-500 tracking-wide">真实 · 安全 · 高效 · 可靠</p>

            <div className="mt-10 space-y-6">
              {[
                { icon: <AuditOutlined />, title: '严格审核', desc: 'AI + 人工双重审核，确保企业资质真实有效' },
                { icon: <SafetyCertificateOutlined />, title: '信息安全', desc: '银行级数据加密，个人信息全程隐私保护' },
                { icon: <RocketOutlined />, title: '高效匹配', desc: '智能推荐算法，精准匹配你的理想职位' },
              ].map((item) => (
                <div key={item.title} className="flex items-start gap-4">
                  <div className="w-11 h-11 rounded-xl bg-blue-50 flex items-center justify-center flex-shrink-0">
                    <span className="text-xl text-[#1677FF]">{item.icon}</span>
                  </div>
                  <div>
                    <h3 className="font-semibold text-slate-800">{item.title}</h3>
                    <p className="mt-1 text-sm text-slate-500">{item.desc}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Right: Login Card */}
          <div className="flex-1 lg:w-[40%] w-full max-w-md">
            <div className="bg-white rounded-2xl shadow-xl shadow-slate-200/50 p-8">
              <div className="text-center mb-6">
                <h2 className="text-2xl font-bold text-slate-800">欢迎登录放心</h2>
                <p className="mt-2 text-sm text-slate-500">
                  还没有账号？<a href="#" className="text-[#1677FF] hover:underline font-medium">立即注册</a>
                </p>
              </div>

              {/* Tabs */}
              <div className="flex border-b border-slate-200 mb-6">
                <button
                  onClick={() => setActiveTab('candidate')}
                  className={`flex-1 pb-3 text-sm font-medium transition-all duration-200 border-b-2 ${activeTab === 'candidate' ? 'text-[#1677FF] border-[#1677FF]' : 'text-slate-400 border-transparent hover:text-slate-600'}`}
                >求职者登录</button>
                <button
                  onClick={() => setActiveTab('hr')}
                  className={`flex-1 pb-3 text-sm font-medium transition-all duration-200 border-b-2 ${activeTab === 'hr' ? 'text-[#1677FF] border-[#1677FF]' : 'text-slate-400 border-transparent hover:text-slate-600'}`}
                >企业登录</button>
                <button
                  onClick={() => setActiveTab('admin')}
                  className={`flex-1 pb-3 text-xs font-medium transition-all duration-200 border-b-2 ${activeTab === 'admin' ? 'text-slate-700 border-slate-700' : 'text-slate-300 border-transparent hover:text-slate-500'}`}
                >管理</button>
              </div>

              {/* Form */}
              <form onSubmit={handleLogin} className="space-y-5">
                {error && (
                  <div className="bg-red-50 border border-red-200 text-red-600 text-sm rounded-lg px-4 py-3">{error}</div>
                )}
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1.5">手机号</label>
                  <div className="relative">
                    <span className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400"><PhoneOutlined /></span>
                    <input type="tel" maxLength={11} value={phone} onChange={(e) => setPhone(e.target.value)} placeholder="请输入手机号" className="w-full pl-10 pr-4 py-2.5 text-sm border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#1677FF]/20 focus:border-[#1677FF] transition-all duration-200 bg-slate-50/50" />
                  </div>
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1.5">验证码</label>
                  <div className="relative">
                    <span className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400"><LockOutlined /></span>
                    <input type={showPwd ? 'text' : 'password'} maxLength={6} value={code} onChange={(e) => setCode(e.target.value)} placeholder="请输入验证码" className="w-full pl-10 pr-24 py-2.5 text-sm border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#1677FF]/20 focus:border-[#1677FF] transition-all duration-200 bg-slate-50/50" />
                    <button type="button" onClick={() => setShowPwd(!showPwd)} className="absolute right-20 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600">
                      {showPwd ? <EyeInvisibleOutlined /> : <EyeOutlined />}
                    </button>
                    <button type="button" onClick={handleSendCode} disabled={counting > 0} className="absolute right-2 top-1/2 -translate-y-1/2 text-xs font-medium text-[#1677FF] hover:text-blue-700 disabled:text-slate-400 disabled:cursor-not-allowed px-2 py-1 whitespace-nowrap">
                      {counting > 0 ? `${counting}s` : '获取验证码'}
                    </button>
                  </div>
                </div>
                <div className="flex items-center justify-between">
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input type="checkbox" checked={remember} onChange={(e) => setRemember(e.target.checked)} className="w-4 h-4 rounded border-slate-300 text-[#1677FF] focus:ring-[#1677FF]" />
                    <span className="text-sm text-slate-500">记住我</span>
                  </label>
                  <a href="#" className="text-sm text-[#1677FF] hover:underline">忘记密码？</a>
                </div>
                <button type="submit" disabled={loading} className="w-full py-2.5 bg-[#1677FF] hover:bg-blue-600 active:bg-blue-700 text-white font-medium rounded-lg transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed shadow-sm shadow-blue-200">
                  {loading ? '登录中...' : '登录'}
                </button>
              </form>

              <div className="relative my-6">
                <div className="absolute inset-0 flex items-center"><div className="w-full border-t border-slate-200"></div></div>
                <div className="relative flex justify-center"><span className="bg-white px-4 text-xs text-slate-400">其他登录方式</span></div>
              </div>

              <div className="flex justify-center gap-6">
                {[
                  { icon: <WechatOutlined />, color: 'hover:bg-green-50 hover:text-green-500', label: '微信' },
                  { icon: <AlipayCircleOutlined />, color: 'hover:bg-blue-50 hover:text-blue-500', label: '支付宝' },
                  { icon: <MessageOutlined />, color: 'hover:bg-orange-50 hover:text-orange-500', label: '短信' },
                ].map((s) => (
                  <button key={s.label} title={s.label} className={`w-11 h-11 rounded-full border border-slate-200 flex items-center justify-center text-xl text-slate-400 transition-all duration-200 ${s.color}`}>{s.icon}</button>
                ))}
              </div>
            </div>
          </div>
        </div>
      </main>

      {/* ===== Features Bar ===== */}
      <section className="bg-white border-t border-slate-100">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
            {[
              { icon: <GlobalOutlined />, title: '海量优质企业', desc: '5000+ 认证企业在线招聘' },
              { icon: <AuditOutlined />, title: '真实职位信息', desc: 'AI 审核 + 人工复核保障' },
              { icon: <ThunderboltOutlined />, title: '精准职位推荐', desc: '智能算法匹配最适合的你' },
              { icon: <SafetyCertificateOutlined />, title: '全程安心服务', desc: '从投递到入职全程保障' },
            ].map((item) => (
              <div key={item.title} className="flex items-start gap-4">
                <div className="w-14 h-14 rounded-xl bg-blue-50 flex items-center justify-center flex-shrink-0">
                  <span className="text-2xl text-[#1677FF]">{item.icon}</span>
                </div>
                <div>
                  <h3 className="font-semibold text-slate-800 text-sm">{item.title}</h3>
                  <p className="mt-1 text-xs text-slate-500">{item.desc}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ===== Footer ===== */}
      <footer className="bg-white border-t border-slate-100">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <div className="flex flex-col md:flex-row items-center justify-between gap-4">
            <div className="flex items-center gap-2">
              <div className="w-7 h-7 rounded-md bg-[#1677FF] flex items-center justify-center">
                <span className="text-white font-bold text-xs">放</span>
              </div>
              <span className="text-sm font-semibold text-slate-700">放心</span>
            </div>
            <nav className="flex items-center gap-6">
              {['关于我们', '帮助中心', '用户协议', '隐私政策', '联系我们'].map((item) => (
                <a key={item} href="#" className="text-sm text-slate-400 hover:text-[#1677FF] transition-colors duration-200">{item}</a>
              ))}
            </nav>
            <p className="text-xs text-slate-400">© 2024 放心 All Rights Reserved.</p>
          </div>
        </div>
      </footer>
    </div>
  );
}
