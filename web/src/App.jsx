import { useEffect, useState } from "react";
import { Routes, Route, Navigate, useNavigate } from "react-router-dom";
import Cookies from "js-cookie";
import LoginPage from "./pages/LoginPage";
import DashboardPage from "./pages/DashboardPage";
import GlobalLoadingOverlay from "./components/GlobalLoadingOverlay";
import {
  NotificationProvider,
  useNotification,
} from "./utils/NotificationContext";
import { initNotifications } from "./services/api";
import "./App.css";

// 带有通知功能的App内容组件
function AppContent() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const navigate = useNavigate();
  const { showSuccess, showError } = useNotification();

  useEffect(() => {
    // 初始化API通知回调
    initNotifications(showSuccess, showError);

    // 检查是否已认证
    const backendSecret = Cookies.get("BACKEND_SECRET");
    if (backendSecret) {
      setIsAuthenticated(true);
    }
  }, [showSuccess, showError]);

  // 登出函数
  const handleLogout = () => {
    Cookies.remove("BACKEND_SECRET");
    setIsAuthenticated(false);
    navigate("/login");
  };

  // 登录成功回调
  const handleLoginSuccess = (secret) => {
    Cookies.set("BACKEND_SECRET", secret, { expires: 7 }); // 保存7天
    setIsAuthenticated(true);
    navigate("/dashboard");
  };

  return (
    <div className="app-container">
      <GlobalLoadingOverlay />
      <Routes>
        <Route
          path="/login"
          element={
            isAuthenticated ? (
              <Navigate to="/dashboard" />
            ) : (
              <LoginPage onLoginSuccess={handleLoginSuccess} />
            )
          }
        />
        <Route
          path="/dashboard/*"
          element={
            isAuthenticated ? (
              <DashboardPage onLogout={handleLogout} />
            ) : (
              <Navigate to="/login" />
            )
          }
        />
        <Route
          path="*"
          element={<Navigate to={isAuthenticated ? "/dashboard" : "/login"} />}
        />
      </Routes>
    </div>
  );
}

// 主App组件
function App() {
  return (
    <NotificationProvider>
      <AppContent />
    </NotificationProvider>
  );
}

export default App;
