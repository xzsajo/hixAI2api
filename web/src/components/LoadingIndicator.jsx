import React from "react";
import "../styles/LoadingIndicator.css";

const LoadingIndicator = ({ text = "加载中..." }) => {
  return (
    <div className="loading-container">
      <div className="loading-spinner">
        <div className="spinner-circle"></div>
        <div className="spinner-circle"></div>
        <div className="spinner-circle"></div>
      </div>
      <div className="loading-text">{text}</div>
    </div>
  );
};

export default LoadingIndicator;
