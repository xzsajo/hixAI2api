FROM --platform=$BUILDPLATFORM node:16-bullseye AS frontend-builder

WORKDIR /app

# 首先只复制依赖相关文件，利用Docker缓存层
COPY web/package*.json ./web/

# 安装所有依赖（包括开发依赖，因为构建工具通常是开发依赖）
RUN cd ./web && npm ci

# 为Node.js提供crypto polyfill
RUN cd ./web && npm install --save-dev crypto-browserify

# 复制前端源代码
COPY web ./web

RUN mkdir -p /app/web/dist

# 创建一个临时polyfill文件来模拟crypto.getRandomValues
# 注意：这里将文件扩展名改为.cjs，使其被视为CommonJS模块
RUN echo "const crypto = require('crypto'); \
if (!crypto.getRandomValues) { \
  crypto.getRandomValues = function(array) { \
    const bytes = crypto.randomBytes(array.length); \
    array.set(bytes); \
    return array; \
  }; \
} \
global.crypto = crypto;" > /app/web/crypto-polyfill.cjs

# 使用node -r选项预加载crypto-polyfill.cjs
RUN cd ./web && NODE_ENV=production node -r /app/web/crypto-polyfill.cjs ./node_modules/vite/bin/vite.js build

# Go构建阶段：使用Alpine镜像
FROM golang:alpine AS builder

# 安装编译依赖
RUN apk add --no-cache \
    gcc \
    musl-dev \
    sqlite-dev \
    build-base \
    grep

# 启用CGO并配置环境
ENV CGO_ENABLED=1 \
    GO111MODULE=on \
    GOOS=linux

WORKDIR /build

# 先只复制go.mod和go.sum以利用缓存
COPY go.mod go.sum ./

RUN go mod download

# 复制版本文件并提取版本号
COPY ./common/constants.go ./common/constants.go

# 复制源代码
COPY . .

# 复制前端构建产物
COPY --from=frontend-builder /app/web/dist /build/web/dist

# 使用缓存优化和并行构建
RUN go build -trimpath -ldflags "-s -w -linkmode external -extldflags '-static'" -o /app/hixai2api


# 最终运行镜像：使用Alpine
FROM alpine:latest

# 添加非root用户
RUN adduser -D -u 1000 hixai2apiuser && \
    apk add --no-cache \
    ca-certificates \
    tzdata

# 复制二进制文件
COPY --from=builder /app/hixai2api /hixai2api

# 配置容器
EXPOSE 7044

WORKDIR /app/hixai2api/data

# 健康检查
#HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 CMD wget -q --spider http://localhost:7044/swagger/index.html || exit 1

# 启动应用
ENTRYPOINT ["/hixai2api"]
