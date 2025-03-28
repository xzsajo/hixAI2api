import { useState } from "react";
import { Routes, Route, NavLink } from "react-router-dom";
import AdminLogo from "../components/AdminLogo";
import CookieManagement from "./CookieManagement";
import ApiKeyManagement from "./ApiKeyManagement";
import Hix2ApiLogo from "../components/AppleLogo";
import "../styles/DashboardPage.css";

const DashboardPage = ({ onLogout }) => {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  const toggleMobileMenu = () => {
    setIsMobileMenuOpen(!isMobileMenuOpen);
  };

  const closeMobileMenu = () => {
    setIsMobileMenuOpen(false);
  };

  return (
    <div className="dashboard-container">
      {/* 移动端菜单按钮 */}
      <button
        className="mobile-menu-button"
        onClick={toggleMobileMenu}
        aria-label={isMobileMenuOpen ? "关闭菜单" : "打开菜单"}
      >
        <span></span>
        <span></span>
        <span></span>
      </button>

      {/* 侧边栏 */}
      <aside className={`sidebar ${isMobileMenuOpen ? "open" : ""}`}>
        <div className="sidebar-header">
          <Hix2ApiLogo className="sidebar-logo" />
          <h2>HIX2API 管理</h2>
        </div>

        <nav className="sidebar-nav">
          <NavLink
            to="/dashboard/cookies"
            className={({ isActive }) =>
              isActive ? "nav-link active" : "nav-link"
            }
            onClick={closeMobileMenu}
          >
            Cookie 管理
          </NavLink>
          <NavLink
            to="/dashboard/api-keys"
            className={({ isActive }) =>
              isActive ? "nav-link active" : "nav-link"
            }
            onClick={closeMobileMenu}
          >
            API Key 管理
          </NavLink>
        </nav>

        <div className="sidebar-footer">
          <button onClick={onLogout} className="logout-button">
            退出登录
          </button>
        </div>
      </aside>

      {/* 内容区 */}
      <main className="dashboard-content">
        <Routes>
          <Route path="cookies" element={<CookieManagement />} />
          <Route path="api-keys" element={<ApiKeyManagement />} />
          <Route path="*" element={<CookieManagement />} />
        </Routes>
      </main>

      {/* 移动端菜单遮罩 */}
      {isMobileMenuOpen && (
        <div className="sidebar-backdrop" onClick={closeMobileMenu}></div>
      )}
    </div>
  );
};

export default DashboardPage;
