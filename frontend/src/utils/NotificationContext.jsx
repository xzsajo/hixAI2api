import { createContext, useState, useContext } from "react";
import Notification from "../components/Notification";

// 创建通知上下文
const NotificationContext = createContext();

// 通知提供者组件
export const NotificationProvider = ({ children }) => {
  const [notification, setNotification] = useState({
    visible: false,
    type: "success",
    message: "",
  });

  // 显示成功通知
  const showSuccess = (message) => {
    setNotification({
      visible: true,
      type: "success",
      message,
    });
  };

  // 显示错误通知
  const showError = (message) => {
    setNotification({
      visible: true,
      type: "error",
      message,
    });
  };

  // 关闭通知
  const closeNotification = () => {
    setNotification({
      ...notification,
      visible: false,
    });
  };

  return (
    <NotificationContext.Provider
      value={{
        showSuccess,
        showError,
      }}
    >
      {children}
      <Notification
        visible={notification.visible}
        type={notification.type}
        message={notification.message}
        onClose={closeNotification}
      />
    </NotificationContext.Provider>
  );
};

// 自定义钩子，用于访问通知上下文
export const useNotification = () => {
  const context = useContext(NotificationContext);
  if (!context) {
    throw new Error(
      "useNotification must be used within a NotificationProvider"
    );
  }
  return context;
};
