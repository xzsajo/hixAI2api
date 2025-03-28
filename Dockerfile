FROM --platform=$BUILDPLATFORM node:16-bullseye AS frontend-builder

WORKDIR /app

# 首先只复制依赖相关文件，利用Docker缓存层
COPY ./frontend/package*.json ./frontend/
# 安装所有依赖（包括开发依赖，因为构建工具通常是开发依赖）
RUN cd ./frontend && npm ci

# 复制前端源代码
COPY ./frontend ./frontend
RUN mkdir -p /app/frontend/dist

# 构建前端项目
RUN cd ./frontend && NODE_ENV=production npm run build

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
# 从common/constants.go中提取版本号
RUN grep -oP 'var Version = "\K[^"]+' ./common/constants.go > VERSION

# 复制源代码
COPY . .
# 复制前端构建产物
COPY --from=frontend-builder /app/frontend/dist /build/frontend/dist

# 版本号处理逻辑
RUN if [ ! -s VERSION ]; then \
        if [ -d .git ]; then \
            git describe --tags > VERSION || echo "v1.0.0" > VERSION; \
        else \
            echo "v1.0.0" > VERSION; \
        fi; \
    fi

# 使用缓存优化和并行构建
RUN go build -trimpath -ldflags "-s -w -X 'hixai2api/common.Version=$(cat VERSION)' -linkmode external -extldflags '-static'" -o /app/hixai2api

# 最终运行镜像：使用Alpine
FROM alpine:latest

# 添加非root用户
RUN adduser -D -u 1000 appuser && \
    apk add --no-cache \
    ca-certificates \
    tzdata

# 复制二进制文件
COPY --from=builder /app/hixai2api /hixai2api

# 创建并设置数据目录权限
RUN mkdir -p /app/hixai2api/data && \
    chown -R appuser:appuser /app/hixai2api

# 切换到非root用户
USER appuser

# 配置容器
EXPOSE 7044
WORKDIR /app/hixai2api/data

# 健康检查
#HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 CMD wget -q --spider http://localhost:7044/swagger/index.html || exit 1

# 启动应用
ENTRYPOINT ["/hixai2api"]
