FROM --platform=$BUILDPLATFORM node:16 AS frontend-builder

WORKDIR /app
COPY ./frontend/package*.json ./frontend/
RUN cd ./frontend && npm install

COPY ./frontend ./frontend
# 从common/constants.go中获取版本号
RUN mkdir -p /app/frontend/dist

# 构建前端项目
RUN cd ./frontend && npm run build

# 构建阶段：使用 Alpine 镜像确保 musl libc 兼容性
FROM golang:alpine AS builder

# 安装编译依赖（SQLite + CGO 必需）
RUN apk add --no-cache \
    gcc \
    musl-dev \
    sqlite-dev \
    build-base

# 启用 CGO 并配置环境
ENV CGO_ENABLED=1 \
    GO111MODULE=on \
    GOOS=linux

WORKDIR /build

# 复制依赖文件（利用 Docker 缓存层加速构建）
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .
# 从前端构建阶段复制构建产物到正确的嵌入路径
COPY --from=frontend-builder /app/frontend/dist /build/frontend/dist

# 使用git tag作为版本号，如果没有则使用common/constants.go中的默认值
RUN if [ -d .git ]; then \
        git describe --tags > VERSION || echo "v1.3.1" > VERSION; \
    else \
        echo "v1.3.1" > VERSION; \
    fi

# 执行构建，添加版本号
RUN go build -trimpath -ldflags "-s -w -X 'hixai2api/common.Version=$(cat VERSION)' -linkmode external -extldflags '-static'" -o /app/hixai2api

# ----------------------------
# 运行时阶段：最小化 Alpine 镜像
FROM alpine:latest

# 安装运行时基础依赖
RUN apk add --no-cache \
    ca-certificates \
    tzdata

# 从构建阶段复制二进制文件
COPY --from=builder /app/hixai2api /hixai2api

# 配置容器
EXPOSE 7044
WORKDIR /app/hixai2api/data
ENTRYPOINT ["/hixai2api"]
