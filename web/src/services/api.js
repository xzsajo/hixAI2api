import axios from "axios";
import Cookies from "js-cookie";

let notificationCallbacks = {
  success: () => {},
  error: () => {},
};

// 记录当前活跃请求数
let activeRequestsCount = 0;
// 加载状态改变回调
let loadingStateCallbacks = [];

// 注册加载状态变化的回调函数
export const registerLoadingCallback = (callback) => {
  loadingStateCallbacks.push(callback);
};

// 更新加载状态
const updateLoadingState = (isLoading) => {
  loadingStateCallbacks.forEach((callback) => callback(isLoading));
};

// 初始化通知回调函数
export const initNotifications = (successCallback, errorCallback) => {
  notificationCallbacks.success = successCallback;
  notificationCallbacks.error = errorCallback;
};

// 创建一个 axios 实例
const api = axios.create({
  baseURL: "/api",
  timeout: 10000,
  headers: {
    "Content-Type": "application/json",
  },
});

// 请求拦截器 - 添加认证头
api.interceptors.request.use(
  (config) => {
    const backendSecret = Cookies.get("BACKEND_SECRET");
    if (backendSecret) {
      config.headers["Authorization"] = backendSecret;
    }

    // 增加活跃请求计数
    activeRequestsCount++;
    if (activeRequestsCount === 1) {
      // 第一个请求开始时，触发加载状态
      updateLoadingState(true);
    }

    return config;
  },
  (error) => {
    // 请求错误，减少计数
    activeRequestsCount = Math.max(0, activeRequestsCount - 1);
    if (activeRequestsCount === 0) {
      updateLoadingState(false);
    }

    return Promise.reject(error);
  }
);

// 响应拦截器 - 处理认证错误和通知
api.interceptors.response.use(
  (response) => {
    // 减少活跃请求计数
    activeRequestsCount = Math.max(0, activeRequestsCount - 1);
    if (activeRequestsCount === 0) {
      // 所有请求完成，关闭加载状态
      setTimeout(() => updateLoadingState(false), 300); // 添加短暂延迟，让过渡更平滑
    }

    // 如果响应成功，显示成功通知（只针对PUT、POST和DELETE请求）
    if (
      ["PUT", "POST", "DELETE"].includes(response.config.method.toUpperCase())
    ) {
      notificationCallbacks.success("操作成功");
    }
    return response;
  },
  (error) => {
    // 减少活跃请求计数
    activeRequestsCount = Math.max(0, activeRequestsCount - 1);
    if (activeRequestsCount === 0) {
      // 所有请求完成，关闭加载状态
      setTimeout(() => updateLoadingState(false), 300);
    }

    if (error.response) {
      // 服务器返回的错误状态码
      if (error.response.status === 401) {
        // 如果是未授权，清除cookie并重定向到登录页
        Cookies.remove("BACKEND_SECRET");
        window.location.href = "/login";
      } else {
        // 显示错误通知
        const errorMessage =
          error.response.data?.message || "请求失败，请稍后重试";
        notificationCallbacks.error(errorMessage);
      }
    } else if (error.request) {
      // 请求发送但没有收到响应
      notificationCallbacks.error("服务器无响应，请检查网络连接");
    } else {
      // 设置请求时发生错误
      notificationCallbacks.error("请求错误: " + error.message);
    }
    return Promise.reject(error);
  }
);

// API 认证函数
export const verifyAuth = async (secret) => {
  try {
    const response = await axios.post(
      "/api/auth/verify",
      {},
      {
        headers: {
          Authorization: secret,
        },
      }
    );
    return response.data;
  } catch (error) {
    throw error;
  }
};

// API Key 相关函数
export const getAllApiKeys = async () => {
  const response = await api.get("/key/all");
  return response.data;
};

export const saveApiKey = async (keyData) => {
  const response = await api.put("/key", keyData);
  return response.data;
};

export const updateApiKey = async (keyData) => {
  const response = await api.post("/key/update", keyData);
  return response.data;
};

export const deleteApiKey = async (id) => {
  const response = await api.delete(`/key/${id}`);
  return response.data;
};

// Cookie 相关函数
export const getAllCookies = async () => {
  const response = await api.get("/cookie/all");
  return response.data;
};

export const saveCookie = async (cookieData) => {
  const response = await api.put("/cookie", cookieData);
  return response.data;
};

export const updateCookie = async (cookieData) => {
  const response = await api.post("/cookie/update", cookieData);
  return response.data;
};

export const deleteCookie = async (id) => {
  const response = await api.delete(`/cookie/${id}`);
  return response.data;
};

export const refreshCookieCredit = async () => {
  const response = await api.post("/cookie/credit/refresh");
  return response.data;
};

export default api;
