import { useState, useEffect } from "react";
import {
  getAllCookies,
  saveCookie,
  updateCookie,
  deleteCookie,
  refreshCookieCredit,
} from "../services/api";
import DataContainer from "../components/DataContainer";
import Modal from "../components/Modal";
import "../styles/Management.css";

const CookieManagement = () => {
  const [cookies, setCookies] = useState([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [error, setError] = useState("");
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isEditMode, setIsEditMode] = useState(false);
  const [currentCookie, setCurrentCookie] = useState({
    id: "",
    cookie: "",
    remark: "",
  });
  const [formError, setFormError] = useState("");

  // 列定义
  const columns = [
    {
      key: "cookie",
      title: "Cookie",
      sortable: true,
      render: (row) => truncateText(row.cookie, 40),
    },
    {
      key: "credit",
      title: "额度",
      sortable: true,
      render: (row) => {
        const credit = row.credit || 0;
        let statusClass = "";
        let statusText = "";

        if (credit < 100 && credit > 0) {
          statusClass = "credit-low";
          statusText = "较少";
        } else if (credit >= 100 && credit < 200) {
          statusClass = "credit-medium";
          statusText = "较多";
        } else if (credit >= 200) {
          statusClass = "credit-high";
          statusText = "充足";
        } else {
          statusClass = "credit-empty";
          statusText = "无";
        }

        return (
          <div className="credit-container">
            <span className="credit-value">{credit}</span>
            <span className={`credit-status ${statusClass}`}>{statusText}</span>
          </div>
        );
      },
    },
    { key: "remark", title: "备注", sortable: false },
    { key: "createTime", title: "创建时间", sortable: true },
  ];

  // 文本截断函数
  const truncateText = (text, maxLength) => {
    if (!text) return "";
    return text.length > maxLength
      ? `${text.substring(0, maxLength)}...`
      : text;
  };

  // 加载数据
  const loadCookies = async () => {
    try {
      setLoading(true);
      setError("");
      const response = await getAllCookies();
      if (response && response.code === 0) {
        setCookies(response.data || []);
      } else {
        setError(response.msg || "获取Cookie失败");
      }
    } catch (err) {
      console.error("加载Cookie错误:", err);
      setError("加载Cookie失败，请稍后重试");
    } finally {
      setLoading(false);
    }
  };

  // 初始加载
  useEffect(() => {
    loadCookies();
  }, []);

  // 刷新全部Cookie额度
  const handleRefreshCredit = async () => {
    try {
      setRefreshing(true);
      setError("");
      const response = await refreshCookieCredit();
      if (response && response.code === 0) {
        loadCookies();
      } else {
        setError(response.msg || "刷新Cookie额度失败");
      }
    } catch (err) {
      console.error("刷新Cookie额度错误:", err);
      setError("刷新Cookie额度失败，请稍后重试");
    } finally {
      setRefreshing(false);
    }
  };

  // 打开添加模态框
  const handleAddClick = () => {
    setIsEditMode(false);
    setCurrentCookie({ id: "", cookie: "", remark: "" });
    setFormError("");
    setIsModalOpen(true);
  };

  // 打开编辑模态框
  const handleEditClick = (cookie) => {
    setIsEditMode(true);
    setCurrentCookie({
      id: cookie.id,
      cookie: cookie.cookie,
      remark: cookie.remark || "",
    });
    setFormError("");
    setIsModalOpen(true);
  };

  // 处理删除
  const handleDeleteClick = async (cookie) => {
    if (window.confirm(`确定要删除该 Cookie 吗？`)) {
      try {
        const response = await deleteCookie(cookie.id);
        if (response && response.code === 0) {
          loadCookies();
        } else {
          setError(response.msg || "删除Cookie失败");
        }
      } catch (err) {
        console.error("删除Cookie错误:", err);
        setError("删除Cookie失败，请稍后重试");
      }
    }
  };

  // 处理表单输入变化
  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setCurrentCookie({
      ...currentCookie,
      [name]: value,
    });
  };

  // 处理表单提交
  const handleFormSubmit = async () => {
    if (!currentCookie.cookie.trim()) {
      setFormError("请输入Cookie");
      return;
    }

    try {
      let response;
      if (isEditMode) {
        response = await updateCookie(currentCookie);
      } else {
        response = await saveCookie(currentCookie);
      }

      if (response && response.code === 0) {
        setIsModalOpen(false);
        loadCookies();
      } else {
        setFormError(
          response.msg || `${isEditMode ? "更新" : "添加"}Cookie失败`
        );
      }
    } catch (err) {
      console.error(`${isEditMode ? "更新" : "添加"}Cookie错误:`, err);
      setFormError(`${isEditMode ? "更新" : "添加"}Cookie失败，请稍后重试`);
    }
  };

  // 处理键盘事件
  const handleKeyDown = (e) => {
    if (e.key === "Enter" && !e.shiftKey) {
      // 如果按下回车且没有同时按下shift键，则提交表单
      e.preventDefault();
      handleFormSubmit();
    }
  };

  return (
    <div className="management-container">
      <div className="management-header">
        <h1>Cookie 管理</h1>
        <div className="management-actions">
          <button
            className="refresh-button"
            onClick={handleRefreshCredit}
            disabled={refreshing}
          >
            {refreshing ? "刷新中..." : "刷新额度"}
          </button>
          <button className="add-button" onClick={handleAddClick}>
            添加 Cookie
          </button>
        </div>
      </div>

      {error && <div className="error-message">{error}</div>}

      <DataContainer
        loading={loading}
        loadingText="正在加载 Cookie 数据..."
        columns={columns}
        data={cookies}
        onEdit={handleEditClick}
        onDelete={handleDeleteClick}
        emptyMessage="暂无 Cookie 数据"
      />

      <Modal
        title={isEditMode ? "编辑 Cookie" : "添加 Cookie"}
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleFormSubmit}
        submitText={isEditMode ? "更新" : "添加"}
      >
        <div className="form-container">
          {formError && <div className="form-error">{formError}</div>}

          <div className="form-group">
            <label htmlFor="cookie">Cookie</label>
            <textarea
              id="cookie"
              name="cookie"
              value={currentCookie.cookie}
              onChange={handleInputChange}
              onKeyDown={handleKeyDown}
              placeholder="请输入Cookie"
              rows={6}
            />
            <small className="form-help-text">
              请输入完整的Cookie值，系统会自动验证该Cookie是否有效
            </small>
          </div>

          <div className="form-group">
            <label htmlFor="remark">备注</label>
            <input
              type="text"
              id="remark"
              name="remark"
              value={currentCookie.remark}
              onChange={handleInputChange}
              onKeyDown={handleKeyDown}
              placeholder="请输入备注（可选）"
            />
          </div>
        </div>
      </Modal>
    </div>
  );
};

export default CookieManagement;
