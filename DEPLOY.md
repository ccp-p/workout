# 健身训练应用服务器部署指南

## 🚀 自动部署（推荐）

### 1. 准备部署脚本
```bash
# 赋予执行权限
chmod +x deploy.sh

# 执行部署（替换为你的服务器信息）
./deploy.sh 你的服务器IP 用户名
```

### 2. 部署示例
```bash
# 示例：部署到IP为 192.168.1.100 的服务器，用户名为 root
./deploy.sh 192.168.1.100 root

# 示例：部署到域名服务器
./deploy.sh your-domain.com ubuntu
```

## 🔧 手动部署

### 1. 服务器环境准备

#### 安装Go环境（Ubuntu/Debian）
```bash
# 下载Go
wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz

# 安装Go
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# 配置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 验证安装
go version
```

#### 安装Go环境（CentOS/RHEL）
```bash
# 下载Go
wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz

# 安装Go
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# 配置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bash_profile
source ~/.bash_profile

# 验证安装
go version
```

### 2. 上传应用文件

#### 方法一：使用SCP
```bash
# 打包本地文件
tar -czf workout-tracker.tar.gz backend/ frontend/ data/ start.sh workout-tracker.service

# 上传到服务器
scp workout-tracker.tar.gz user@your-server:/opt/

# 在服务器上解压
ssh user@your-server
cd /opt
sudo tar -xzf workout-tracker.tar.gz
sudo mkdir -p workout-tracker/uploads
```

#### 方法二：使用Git
```bash
# 在服务器上克隆仓库
cd /opt
sudo git clone https://github.com/your-username/workout.git workout-tracker
cd workout-tracker
sudo mkdir -p uploads
```

### 3. 配置应用

#### 安装依赖
```bash
cd /opt/workout-tracker/backend
sudo go mod tidy
```

#### 设置权限
```bash
sudo chmod +x /opt/workout-tracker/start.sh
sudo chown -R www-data:www-data /opt/workout-tracker/
```

### 4. 配置系统服务

#### 创建systemd服务
```bash
# 复制服务文件
sudo cp /opt/workout-tracker/workout-tracker.service /etc/systemd/system/

# 重新加载服务
sudo systemctl daemon-reload

# 启用服务（开机自启）
sudo systemctl enable workout-tracker

# 启动服务
sudo systemctl start workout-tracker

# 查看服务状态
sudo systemctl status workout-tracker
```

#### 服务管理命令
```bash
# 启动服务
sudo systemctl start workout-tracker

# 停止服务
sudo systemctl stop workout-tracker

# 重启服务
sudo systemctl restart workout-tracker

# 查看日志
sudo journalctl -u workout-tracker -f

# 查看服务状态
sudo systemctl status workout-tracker
```

### 5. 配置防火墙

#### Ubuntu/Debian (UFW)
```bash
# 允许8769端口
sudo ufw allow 8769/tcp

# 启用防火墙
sudo ufw enable

# 查看防火墙状态
sudo ufw status
```

#### CentOS/RHEL (firewalld)
```bash
# 允许8769端口
sudo firewall-cmd --permanent --add-port=8769/tcp

# 重新加载防火墙
sudo firewall-cmd --reload

# 查看开放端口
sudo firewall-cmd --list-ports
```

### 6. 配置反向代理（可选）

#### 使用Nginx
```bash
# 安装Nginx
sudo apt update
sudo apt install nginx

# 创建配置文件
sudo nano /etc/nginx/sites-available/workout-tracker
```

**Nginx配置内容：**
```nginx
server {
    listen 80;
    server_name your-domain.com;  # 替换为你的域名
    
    location / {
        proxy_pass http://localhost:8769;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

```bash
# 启用站点
sudo ln -s /etc/nginx/sites-available/workout-tracker /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重启Nginx
sudo systemctl restart nginx
```

## 🌐 域名和SSL配置

### 1. 配置域名
- 在域名服务商处将域名A记录指向服务器IP
- 等待DNS解析生效（通常5-30分钟）

### 2. 配置SSL证书（使用Let's Encrypt）
```bash
# 安装Certbot
sudo apt install certbot python3-certbot-nginx

# 获取SSL证书
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo crontab -e
# 添加以下行：
0 12 * * * /usr/bin/certbot renew --quiet
```

## 📊 监控和维护

### 1. 查看应用日志
```bash
# 查看实时日志
sudo journalctl -u workout-tracker -f

# 查看最近日志
sudo journalctl -u workout-tracker --since "1 hour ago"

# 查看错误日志
sudo journalctl -u workout-tracker -p err
```

### 2. 数据备份
```bash
# 创建备份脚本
sudo nano /opt/backup-workout.sh
```

**备份脚本内容：**
```bash
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/opt/backups"
APP_DIR="/opt/workout-tracker"

mkdir -p $BACKUP_DIR

# 备份数据文件
tar -czf $BACKUP_DIR/workout-data-$DATE.tar.gz $APP_DIR/data/

# 删除7天前的备份
find $BACKUP_DIR -name "workout-data-*.tar.gz" -mtime +7 -delete

echo "备份完成: workout-data-$DATE.tar.gz"
```

```bash
# 设置定时备份
sudo chmod +x /opt/backup-workout.sh
sudo crontab -e
# 添加每日备份：
0 2 * * * /opt/backup-workout.sh
```

### 3. 更新应用
```bash
# 停止服务
sudo systemctl stop workout-tracker

# 更新代码
cd /opt/workout-tracker
sudo git pull

# 更新依赖
cd backend
sudo go mod tidy

# 启动服务
sudo systemctl start workout-tracker
```

## 🔍 故障排除

### 1. 常见问题

#### 服务无法启动
```bash
# 检查Go环境
go version

# 检查端口占用
sudo netstat -tlnp | grep 8769

# 查看详细错误
sudo journalctl -u workout-tracker -n 50
```

#### 无法访问应用
```bash
# 检查服务状态
sudo systemctl status workout-tracker

# 检查防火墙
sudo ufw status
sudo firewall-cmd --list-ports

# 检查端口监听
sudo ss -tlnp | grep 8769
```

#### 数据丢失
```bash
# 恢复备份
cd /opt/workout-tracker
sudo rm -rf data/
sudo tar -xzf /opt/backups/workout-data-YYYYMMDD_HHMMSS.tar.gz
sudo systemctl restart workout-tracker
```

### 2. 性能优化

#### 数据库优化（如果使用数据库）
- 定期清理旧数据
- 优化查询索引
- 监控数据库性能

#### 服务器优化
```bash
# 监控系统资源
htop
free -h
df -h

# 优化Go应用
export GOMEMLIMIT=512MiB
```

## 📞 支持

如遇到问题，请检查：
1. 服务日志：`sudo journalctl -u workout-tracker -f`
2. 系统资源：`htop`，`free -h`
3. 网络连接：`ping your-domain.com`
4. 防火墙设置：`sudo ufw status`

**访问地址：**
- 应用首页：`http://your-domain.com` 或 `http://your-server-ip:8769`
- 后台管理：`http://your-domain.com` 或 `http://your-server-ip:8769`
- 移动端训练：`http://your-domain.com/mobile` 或 `http://your-server-ip:8769/mobile`
