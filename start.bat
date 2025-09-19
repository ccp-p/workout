@echo off
echo 正在启动健身训练应用...
echo.

cd /d "%~dp0backend"

echo 检查Go环境...
go version >nul 2>&1
if errorlevel 1 (
    echo 错误: 未找到Go环境，请先安装Go
    echo 下载地址: https://golang.org/dl/
    pause
    exit /b 1
)

echo 安装依赖包...
go mod tidy
if errorlevel 1 (
    echo 错误: 安装依赖包失败
    pause
    exit /b 1
)

echo.
echo 启动服务器...
echo 后台管理: http://localhost:8080
echo 移动端训练: http://localhost:8080/mobile
echo.
echo 按 Ctrl+C 停止服务器
echo.

go run main.go

pause
