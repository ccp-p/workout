#!/bin/bash

# 健身训练应用服务器部署脚本
# 使用方法: ./deploy.sh [服务器IP] [用户名]

SERVER_IP=${1:-"your-server-ip"}
SERVER_USER=${2:-"root"}
APP_NAME="workout-tracker"
REMOTE_DIR="/opt/$APP_NAME"
SERVICE_NAME="workout-tracker"

echo "开始部署健身训练应用到服务器: $SERVER_USER@$SERVER_IP"

# 1. 打包应用
echo "正在打包应用..."
rm -rf dist
mkdir -p dist

# 复制后端文件
cp -r backend dist/
cp -r frontend dist/
cp -r data dist/
mkdir -p dist/uploads

# 复制配置文件
cp start.sh dist/
cp $SERVICE_NAME.service dist/

echo "应用打包完成"

# 2. 上传到服务器
echo "正在上传文件到服务器..."
scp -r dist/* $SERVER_USER@$SERVER_IP:$REMOTE_DIR/

# 3. 在服务器上执行部署命令
echo "正在服务器上配置应用..."
ssh $SERVER_USER@$SERVER_IP << EOF
    # 进入应用目录
    cd $REMOTE_DIR

    # 安装Go依赖
    cd backend
    go mod tidy
    cd ..

    # 设置可执行权限
    chmod +x start.sh

    # 配置systemd服务
    sudo cp $SERVICE_NAME.service /etc/systemd/system/
    sudo systemctl daemon-reload
    sudo systemctl enable $SERVICE_NAME
    sudo systemctl restart $SERVICE_NAME

    # 配置防火墙（如果需要）
    sudo ufw allow 8769/tcp

    echo "部署完成！应用已启动在端口8769"
    sudo systemctl status $SERVICE_NAME
EOF

echo "部署完成！"
echo "访问地址: http://$SERVER_IP:8769"
echo "后台管理: http://$SERVER_IP:8769"
echo "移动端训练: http://$SERVER_IP:8769/mobile"
