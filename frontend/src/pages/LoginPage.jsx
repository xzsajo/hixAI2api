import { useState } from "react";
import { verifyAuth } from "../services/api";
import Hix2ApiLogo from "../components/AppleLogo";
import "../styles/LoginPage.css";

const LoginPage = ({ onLoginSuccess }) => {
  const [secret, setSecret] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!secret.trim()) {
      setError("请输入管理密钥");
      return;
    }

    setIsLoading(true);
    setError("");

    try {
      await verifyAuth(secret);
      onLoginSuccess(secret);
    } catch (err) {
      if (err.response && err.response.status === 401) {
        setError("管理密钥无效，请重试");
      } else {
        setError("登录失败，请稍后重试");
        console.error("登录错误:", err);
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="login-container">
      <div className="login-card">
        <div className="login-header">
          <Hix2ApiLogo className="login-logo" />
          <h1>HIX2API 管理后台</h1>
        </div>

        <form className="login-form" onSubmit={handleSubmit}>
          {error && <div className="login-error">{error}</div>}

          <div className="form-group">
            <label htmlFor="secret">管理密钥</label>
            <input
              type="password"
              id="secret"
              value={secret}
              onChange={(e) => setSecret(e.target.value)}
              placeholder="请输入您的管理密钥"
              disabled={isLoading}
              autoFocus
            />
          </div>

          <button type="submit" className="login-button" disabled={isLoading}>
            {isLoading ? "登录中..." : "登录"}
          </button>
        </form>
      </div>
    </div>
  );
};

export default LoginPage;
