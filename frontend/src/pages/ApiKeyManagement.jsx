import { useState, useEffect } from "react";
import {
  getAllApiKeys,
  saveApiKey,
  updateApiKey,
  deleteApiKey,
} from "../services/api";
import DataTable from "../components/DataTable";
import Modal from "../components/Modal";
import LoadingIndicator from "../components/LoadingIndicator";
import "../styles/Management.css";

const ApiKeyManagement = () => {
  const [apiKeys, setApiKeys] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isEditMode, setIsEditMode] = useState(false);
  const [currentApiKey, setCurrentApiKey] = useState({
    id: "",
    apiKey: "",
    remark: "",
  });
  const [formError, setFormError] = useState("");

  // 列定义
  const columns = [
    { key: "apiKey", title: "API Key", sortable: true },
    { key: "createTime", title: "创建时间", sortable: true },
    { key: "remark", title: "备注", sortable: false },
  ];

  // 加载数据
  const loadApiKeys = async () => {
    try {
      setLoading(true);
      setError("");
      const response = await getAllApiKeys();
      if (response && response.code === 0) {
        setApiKeys(response.data || []);
      } else {
        setError(response.msg || "获取API Key失败");
      }
    } catch (err) {
      console.error("加载API Key错误:", err);
      setError("加载API Key失败，请稍后重试");
    } finally {
      setLoading(false);
    }
  };

  // 初始加载
  useEffect(() => {
    loadApiKeys();
  }, []);

  // 打开添加模态框
  const handleAddClick = () => {
    setIsEditMode(false);
    setCurrentApiKey({ id: "", apiKey: "", remark: "" });
    setFormError("");
    setIsModalOpen(true);
  };

  // 打开编辑模态框
  const handleEditClick = (apiKey) => {
    setIsEditMode(true);
    setCurrentApiKey({
      id: apiKey.id,
      apiKey: apiKey.apiKey,
      remark: apiKey.remark || "",
    });
    setFormError("");
    setIsModalOpen(true);
  };

  // 处理删除
  const handleDeleteClick = async (apiKey) => {
    if (window.confirm(`确定要删除 API Key "${apiKey.apiKey}" 吗？`)) {
      try {
        const response = await deleteApiKey(apiKey.id);
        if (response && response.code === 0) {
          loadApiKeys();
        } else {
          setError(response.msg || "删除API Key失败");
        }
      } catch (err) {
        console.error("删除API Key错误:", err);
        setError("删除API Key失败，请稍后重试");
      }
    }
  };

  // 处理表单输入变化
  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setCurrentApiKey({
      ...currentApiKey,
      [name]: value,
    });
  };

  // 处理表单提交
  const handleFormSubmit = async () => {
    if (!currentApiKey.apiKey.trim()) {
      setFormError("请输入API Key");
      return;
    }

    try {
      let response;
      if (isEditMode) {
        response = await updateApiKey(currentApiKey);
      } else {
        response = await saveApiKey(currentApiKey);
      }

      if (response && response.code === 0) {
        setIsModalOpen(false);
        loadApiKeys();
      } else {
        setFormError(
          response.msg || `${isEditMode ? "更新" : "添加"}API Key失败`
        );
      }
    } catch (err) {
      console.error(`${isEditMode ? "更新" : "添加"}API Key错误:`, err);
      setFormError(`${isEditMode ? "更新" : "添加"}API Key失败，请稍后重试`);
    }
  };

  return (
    <div className="management-container">
      <div className="management-header">
        <h1>API Key 管理</h1>
        <button className="add-button" onClick={handleAddClick}>
          添加 API Key
        </button>
      </div>

      {error && <div className="error-message">{error}</div>}

      {loading ? (
        <LoadingIndicator text="正在加载 API Key 数据..." />
      ) : (
        <DataTable
          columns={columns}
          data={apiKeys}
          onEdit={handleEditClick}
          onDelete={handleDeleteClick}
          emptyMessage="暂无 API Key 数据"
        />
      )}

      <Modal
        title={isEditMode ? "编辑 API Key" : "添加 API Key"}
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleFormSubmit}
        submitText={isEditMode ? "更新" : "添加"}
      >
        <div className="form-container">
          {formError && <div className="form-error">{formError}</div>}

          <div className="form-group">
            <label htmlFor="apiKey">API Key</label>
            <input
              type="text"
              id="apiKey"
              name="apiKey"
              value={currentApiKey.apiKey}
              onChange={handleInputChange}
              placeholder="请输入API Key"
            />
          </div>

          <div className="form-group">
            <label htmlFor="remark">备注</label>
            <input
              type="text"
              id="remark"
              name="remark"
              value={currentApiKey.remark}
              onChange={handleInputChange}
              placeholder="请输入备注（可选）"
            />
          </div>
        </div>
      </Modal>
    </div>
  );
};

export default ApiKeyManagement;
