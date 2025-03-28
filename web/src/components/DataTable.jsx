import { useState } from "react";
import "../styles/DataTable.css";

const DataTable = ({
  columns,
  data,
  onEdit,
  onDelete,
  emptyMessage = "暂无数据",
}) => {
  const [sortConfig, setSortConfig] = useState({ key: null, direction: "asc" });

  // 排序功能
  const sortedData = () => {
    if (sortConfig.key === null || !data || data.length === 0) {
      return data;
    }

    return [...data].sort((a, b) => {
      const aValue = a[sortConfig.key];
      const bValue = b[sortConfig.key];

      if (aValue === null || aValue === undefined) return 1;
      if (bValue === null || bValue === undefined) return -1;

      if (aValue < bValue) {
        return sortConfig.direction === "asc" ? -1 : 1;
      }
      if (aValue > bValue) {
        return sortConfig.direction === "asc" ? 1 : -1;
      }
      return 0;
    });
  };

  const requestSort = (key) => {
    let direction = "asc";
    if (sortConfig.key === key && sortConfig.direction === "asc") {
      direction = "desc";
    }
    setSortConfig({ key, direction });
  };

  const getSortIcon = (key) => {
    if (sortConfig.key !== key) {
      return null;
    }
    return sortConfig.direction === "asc" ? "↑" : "↓";
  };

  // 生成排序后的表格数据
  const sortedItems = sortedData();

  return (
    <div className="data-table-container">
      <table className="data-table">
        <thead>
          <tr>
            {columns.map((column) => (
              <th
                key={column.key}
                onClick={() => column.sortable && requestSort(column.key)}
                className={column.sortable ? "sortable" : ""}
              >
                {column.title} {column.sortable && getSortIcon(column.key)}
              </th>
            ))}
            <th className="action-column">操作</th>
          </tr>
        </thead>
        <tbody>
          {sortedItems && sortedItems.length > 0 ? (
            sortedItems.map((row, index) => (
              <tr key={row.id || index}>
                {columns.map((column) => (
                  <td key={`${row.id || index}-${column.key}`}>
                    {column.render ? column.render(row) : row[column.key]}
                  </td>
                ))}
                <td className="action-column">
                  <div className="action-buttons">
                    <button
                      className="edit-button"
                      onClick={() => onEdit(row)}
                      aria-label="编辑"
                    >
                      编辑
                    </button>
                    <button
                      className="delete-button"
                      onClick={() => onDelete(row)}
                      aria-label="删除"
                    >
                      删除
                    </button>
                  </div>
                </td>
              </tr>
            ))
          ) : (
            <tr>
              <td colSpan={columns.length + 1} className="empty-message">
                {emptyMessage}
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
};

export default DataTable;
