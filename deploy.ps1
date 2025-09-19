# 健身训练应用服务器部署脚本 (PowerShell版本)
# 使用方法: .\deploy.ps1 [服务器IP] [用户名]

param(
    [string]$ServerIP = "8.134.32.71",
    [string]$ServerUser = "root"
)

$APP_NAME = "workout-tracker"
$REMOTE_DIR = "/opt/$APP_NAME"
$SERVICE_NAME = "workout-tracker"

Write-Host "开始部署健身训练应用到服务器: $ServerUser@$ServerIP" -ForegroundColor Green

# 1. 打包应用
Write-Host "正在打包应用..." -ForegroundColor Yellow
if (Test-Path "dist") {
    Remove-Item -Recurse -Force "dist"
}
New-Item -ItemType Directory -Force -Path "dist" | Out-Null

# 复制文件
Copy-Item -Recurse -Path "backend" -Destination "dist/"
Copy-Item -Recurse -Path "frontend" -Destination "dist/"
Copy-Item -Recurse -Path "data" -Destination "dist/"
New-Item -ItemType Directory -Force -Path "dist/uploads" | Out-Null

# 复制配置文件
Copy-Item -Path "start.sh" -Destination "dist/"
Copy-Item -Path "$SERVICE_NAME.service" -Destination "dist/"

Write-Host "应用打包完成" -ForegroundColor Green

# 2. 上传到服务器
Write-Host "正在上传文件到服务器..." -ForegroundColor Yellow

# 检查是否安装了scp命令
if (!(Get-Command "scp" -ErrorAction SilentlyContinue)) {
    Write-Host "错误: 未找到scp命令。请安装OpenSSH客户端或使用WSL。" -ForegroundColor Red
    Write-Host "你可以通过以下方式安装:" -ForegroundColor Yellow
    Write-Host "1. Windows 设置 > 应用 > 可选功能 > 添加功能 > OpenSSH 客户端" -ForegroundColor Yellow
    Write-Host "2. 或者使用 choco install openssh" -ForegroundColor Yellow
    Write-Host "3. 或者使用WSL执行bash脚本" -ForegroundColor Yellow
    exit 1
}

# 创建远程目录
ssh "$ServerUser@$ServerIP" "sudo mkdir -p $REMOTE_DIR && sudo chown $ServerUser $REMOTE_DIR"

# 上传文件
scp -r dist/* "$ServerUser@$ServerIP`:$REMOTE_DIR/"

# 3. 在服务器上执行部署命令
Write-Host "正在服务器上配置应用..." -ForegroundColor Yellow

$deployCommands = @"
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
"@

ssh "$ServerUser@$ServerIP" $deployCommands

Write-Host "部署完成！" -ForegroundColor Green
Write-Host "访问地址: http://$ServerIP`:8769" -ForegroundColor Cyan
Write-Host "后台管理: http://$ServerIP`:8769" -ForegroundColor Cyan
Write-Host "移动端训练: http://$ServerIP`:8769/mobile" -ForegroundColor Cyan
