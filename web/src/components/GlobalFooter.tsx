import {
  BankOutlined, IdcardOutlined, StarOutlined, SafetyCertificateOutlined,
} from '@ant-design/icons';

// ── Configurable Data ──

const ribbonItems = [
  { icon: <BankOutlined />, title: '海量优质企业', desc: '严选优质企业，靠谱有保障' },
  { icon: <IdcardOutlined />, title: '真实职位信息', desc: '人工审核职位，真实可靠' },
  { icon: <StarOutlined />, title: '精准职位推荐', desc: '智能匹配偏好，高效求职' },
  { icon: <SafetyCertificateOutlined />, title: '全程安心服务', desc: '求职有保障，服务更贴心' },
];

const footerLinks = ['关于我们', '帮助中心', '用户协议', '隐私政策', '联系我们'];

// ── Component ──

export default function GlobalFooter() {
  return (
    <>
      {/* ===== Ribbon Bar ===== */}
      <section className="bg-[#f8f9fa] border-t border-b border-gray-100">
        <div className="max-w-7xl mx-auto px-6 lg:px-8 py-8">
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-8">
            {ribbonItems.map((item) => (
              <div key={item.title} className="flex items-center gap-4">
                <div className="w-11 h-11 rounded-xl bg-blue-50 flex items-center justify-center flex-shrink-0">
                  <span className="text-xl text-[#1677ff]">{item.icon}</span>
                </div>
                <div>
                  <h4 className="text-sm font-bold text-[#1f1f1f]">{item.title}</h4>
                  <p className="text-xs text-gray-400 mt-1">{item.desc}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ===== Global Footer ===== */}
      <footer className="bg-white">
        <div className="max-w-7xl mx-auto px-6 lg:px-8 py-8">
          <div className="flex flex-col md:flex-row items-center justify-between gap-6">
            {/* Left: Logo */}
            <div className="flex items-center gap-2.5">
              <div className="w-9 h-9 rounded-xl bg-[#1677FF] flex items-center justify-center flex-shrink-0">
                <svg viewBox="0 0 24 24" className="w-5 h-5 text-white" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                  <polyline points="20 6 9 17 4 12" />
                </svg>
              </div>
              <div className="leading-tight">
                <span className="text-sm font-bold text-[#1f1f1f]">放心</span>
                <p className="text-[10px] text-[#8c8c8c] leading-none mt-0.5">让求职招聘更放心</p>
              </div>
            </div>

            {/* Center: Links */}
            <nav className="flex items-center gap-6 text-xs text-gray-500 font-medium">
              {footerLinks.map((item) => (
                <a key={item} href="#" className="hover:text-[#1677ff] transition-colors">{item}</a>
              ))}
            </nav>

            {/* Right: Copyright */}
            <div className="text-right text-xs text-gray-400 space-y-0.5">
              <p>© 2026 放心 All Rights Reserved.</p>
              <p>湘ICP备2023001234号-1</p>
            </div>
          </div>
        </div>
      </footer>
    </>
  );
}
