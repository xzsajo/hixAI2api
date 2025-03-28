import { useState, useEffect } from "react";
import { registerLoadingCallback } from "../services/api";
import "../styles/GlobalLoadingOverlay.css";

const GlobalLoadingOverlay = () => {
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    // 注册全局加载状态回调
    registerLoadingCallback(setIsLoading);
  }, []);

  if (!isLoading) return null;

  return (
    <div className="global-loading-overlay">
      <div className="loading-dots">
        <div className="loading-dot"></div>
        <div className="loading-dot"></div>
        <div className="loading-dot"></div>
      </div>
    </div>
  );
};

export default GlobalLoadingOverlay;
