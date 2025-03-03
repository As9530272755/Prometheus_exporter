# 使用官方的 Golang 基础镜像进行构建
FROM golang:alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制源代码到容器中
COPY exporter-demo .

# 编译 Golang 程序
RUN go build -o dial_exporter .

# 使用一个更小的基础镜像来运行应用
FROM alpine:latest

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/dial_exporter .

# 从构建阶段复制配置文件到容器中
COPY --from=builder /app/config.yaml .


# 暴露端口（根据你的应用需求）
EXPOSE 9120

# 启动应用
CMD ["./dial_exporter"]