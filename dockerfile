# 使用官方 Golang 镜像作为构建环境
FROM golang:1.21 as builder

# 设置工作目录
WORKDIR /robot-bot

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

# 使用精简的 Alpine 镜像作为运行环境
FROM alpine:latest

# 安装必要的依赖（如需要）
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制可执行文件
COPY --from=builder /robot/myapp .

# 暴露端口（如果需要）
EXPOSE 8080

# 运行应用
CMD ["./myapp"]