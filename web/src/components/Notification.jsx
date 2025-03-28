import { useEffect } from "react";
import "../styles/Notification.css";

/**
 * 通知组件，用于显示操作成功或失败的提示信息
 * @param {string} type - 通知类型: 'success' 或 'error'
 * @param {string} message - 通知消息内容
 * @param {boolean} visible - 是否显示通知
 * @param {function} onClose - 关闭通知的回调函数
 * @param {number} duration - 自动关闭的时间(毫秒)，默认 3000ms
 */
const Notification = ({ type, message, visible, onClose, duration = 3000 }) => {
  // 自动关闭通知
  useEffect(() => {
    if (visible) {
      const timer = setTimeout(() => {
        onClose();
      }, duration);

      return () => clearTimeout(timer);
    }
  }, [visible, onClose, duration]);

  if (!visible) return null;

  return (
    <div className={`notification ${type} ${visible ? "visible" : ""}`}>
      <div className="notification-content">
        <div className="notification-icon">
          {type === "success" ? (
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="24"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
              <polyline points="22 4 12 14.01 9 11.01"></polyline>
            </svg>
          ) : (
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="24"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="8" x2="12" y2="12"></line>
              <line x1="12" y1="16" x2="12.01" y2="16"></line>
            </svg>
          )}
        </div>
        <div className="notification-message">{message}</div>
        <button className="notification-close" onClick={onClose}>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>
    </div>
  );
};

export default Notification;
