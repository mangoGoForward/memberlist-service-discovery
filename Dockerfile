# 使用官方的 Go 镜像作为基础镜像
FROM golang:1.20-alpine AS builder

# 设置工作目录
WORKDIR /app

# 将源代码复制到容器中
COPY . .

# 构建应用程序
RUN CGO_ENABLED=0 GOOS=linux go build -o memberlist ./src

# 使用一个轻量级的基础镜像作为运行时镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/memberlist .

# 暴露服务端口
EXPOSE 8080 6789

# 启动应用程序
CMD ["./memberlist"]