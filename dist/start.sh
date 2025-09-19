#!/bin/bash

# 健身训练应用启动脚本（Linux版本）
echo "正在启动健身训练应用..."

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "错误: 未找到Go环境，请先安装Go"
    echo "安装命令: wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz"
    echo "          sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz"
    echo "          export PATH=$PATH:/usr/local/go/bin"
    exit 1
fi

# 进入后端目录
cd "$(dirname "$0")/backend"

# 安装依赖
echo "安装Go依赖..."
go mod tidy

# 启动服务器
echo "启动服务器..."
echo "后台管理: http://localhost:8769"
echo "移动端训练: http://localhost:8769/mobile"
echo "按 Ctrl+C 停止服务器"

go run main.go
