import { useEffect, useRef } from "react";
import "../styles/Modal.css";

const Modal = ({
  title,
  children,
  isOpen,
  onClose,
  onSubmit,
  submitText = "保存",
  cancelText = "取消",
  disableSubmit = false,
}) => {
  const modalRef = useRef(null);

  useEffect(() => {
    // 管理滚动锁定
    if (isOpen) {
      document.body.style.overflow = "hidden";
    } else {
      document.body.style.overflow = "";
    }

    // 清理函数
    return () => {
      document.body.style.overflow = "";
    };
  }, [isOpen]);

  // 处理点击外部关闭
  const handleBackdropClick = (e) => {
    if (modalRef.current && !modalRef.current.contains(e.target)) {
      onClose();
    }
  };

  // ESC关闭模态框
  useEffect(() => {
    const handleEsc = (e) => {
      if (e.key === "Escape" && isOpen) {
        onClose();
      }
    };

    window.addEventListener("keydown", handleEsc);
    return () => {
      window.removeEventListener("keydown", handleEsc);
    };
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  return (
    <div className="modal-backdrop" onClick={handleBackdropClick}>
      <div className="modal-container" ref={modalRef}>
        <div className="modal-header">
          <h3>{title}</h3>
          <button
            className="modal-close-button"
            onClick={onClose}
            aria-label="关闭"
          >
            ✕
          </button>
        </div>
        <div className="modal-content">{children}</div>
        <div className="modal-footer">
          <button className="modal-cancel-button" onClick={onClose}>
            {cancelText}
          </button>
          <button
            className="modal-submit-button"
            onClick={onSubmit}
            disabled={disableSubmit}
          >
            {submitText}
          </button>
        </div>
      </div>
    </div>
  );
};

export default Modal;
