import React from "react";
import LoadingIndicator from "./LoadingIndicator";
import DataTable from "./DataTable";
import "../styles/DataContainer.css";

/**
 * 数据容器组件，用于统一封装数据表格和加载指示器的显示
 * 确保在加载状态和数据显示状态下有一致的宽度和样式
 */
const DataContainer = ({
  loading,
  loadingText,
  columns,
  data,
  onEdit,
  onDelete,
  emptyMessage,
}) => {
  return (
    <div className="data-container">
      {loading ? (
        <LoadingIndicator text={loadingText} />
      ) : (
        <DataTable
          columns={columns}
          data={data}
          onEdit={onEdit}
          onDelete={onDelete}
          emptyMessage={emptyMessage}
        />
      )}
    </div>
  );
};

export default DataContainer;
