# 构建阶段：使用 Alpine 镜像确保 musl libc 兼容性
FROM golang:alpine AS builder

# 安装编译依赖（SQLite + CGO 必需）
RUN apk add --no-cache \
    gcc \
    musl-dev \
    sqlite-dev

# 启用 CGO 并配置环境
ENV CGO_ENABLED=1 \
    GO111MODULE=on \
    GOOS=linux

WORKDIR /build

# 复制依赖文件（利用 Docker 缓存层加速构建）
COPY go.mod go.sum ./
RUN go mod download

# 复制源码并静态编译
COPY . .
RUN go build -trimpath -ldflags "-s -w -linkmode external -extldflags '-static'" -o /app/hixai2api

# ----------------------------
# 运行时阶段：最小化 Alpine 镜像
FROM alpine:latest

# 安装运行时基础依赖
RUN apk add --no-cache \
    ca-certificates

# 从构建阶段复制二进制文件
COPY --from=builder /app/hixai2api /hixai2api

# 配置容器
EXPOSE 7044
WORKDIR /app/hixai2api/data
ENTRYPOINT ["/hixai2api"]
