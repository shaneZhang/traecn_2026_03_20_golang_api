# Casdoor API 模块部署文档

## 1. 部署概述

本文档描述了重构后的 Casdoor API 模块的部署方案，包括环境要求、部署步骤、配置说明和运维指南。

## 2. 环境要求

### 2.1 硬件要求

| 组件 | 最低配置 | 推荐配置 |
|------|----------|----------|
| CPU | 4 核 | 8 核+ |
| 内存 | 8 GB | 16 GB+ |
| 磁盘 | 50 GB SSD | 100 GB SSD |
| 网络 | 100 Mbps | 1 Gbps |

### 2.2 软件要求

| 软件 | 版本要求 |
|------|----------|
| Go | 1.21+ |
| MySQL | 8.0+ |
| Nginx | 1.20+ |
| Docker | 20.10+（可选） |
| Kubernetes | 1.25+（可选） |

### 2.3 操作系统支持

- Ubuntu 22.04 LTS
- CentOS 8+
- Debian 11+
- macOS 12+

## 3. 部署架构

### 3.1 单机部署架构

```
┌─────────────────────────────────────────────────────┐
│                    用户请求                          │
└─────────────────────┬───────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────┐
│                   Nginx (反向代理)                   │
│                 监听端口: 80/443                     │
└─────────────────────┬───────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────┐
│                  Casdoor 应用服务                    │
│                 监听端口: 8000                       │
└─────────────────────┬───────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────┐
│                   MySQL 数据库                       │
│                 监听端口: 3306                       │
└─────────────────────────────────────────────────────┘
```

### 3.2 集群部署架构

```
                    ┌─────────────┐
                    │   负载均衡   │
                    │   (Nginx)   │
                    └──────┬──────┘
                           │
           ┌───────────────┼───────────────┐
           │               │               │
           ▼               ▼               ▼
    ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
    │  Casdoor 1  │ │  Casdoor 2  │ │  Casdoor 3  │
    │   :8000     │ │   :8000     │ │   :8000     │
    └──────┬──────┘ └──────┬──────┘ └──────┬──────┘
           │               │               │
           └───────────────┼───────────────┘
                           │
                    ┌──────▼──────┐
                    │ MySQL 主从   │
                    │ 读写分离     │
                    └─────────────┘
```

## 4. 部署步骤

### 4.1 准备工作

#### 4.1.1 安装依赖

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install -y git make gcc

# CentOS/RHEL
sudo yum install -y git make gcc
```

#### 4.1.2 安装 Go

```bash
# 下载 Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz

# 解压安装
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz

# 配置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
source ~/.bashrc

# 验证安装
go version
```

#### 4.1.3 安装 MySQL

```bash
# Ubuntu
sudo apt install -y mysql-server

# CentOS
sudo yum install -y mysql-server

# 启动服务
sudo systemctl start mysql
sudo systemctl enable mysql

# 安全配置
sudo mysql_secure_installation
```

### 4.2 获取代码

```bash
# 克隆代码
git clone https://github.com/casdoor/casdoor.git
cd casdoor

# 切换到重构分支
git checkout refactor/api-module
```

### 4.3 配置数据库

#### 4.3.1 创建数据库

```sql
-- 登录 MySQL
mysql -u root -p

-- 创建数据库
CREATE DATABASE casdoor CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建用户
CREATE USER 'casdoor'@'%' IDENTIFIED BY 'your_password';

-- 授权
GRANT ALL PRIVILEGES ON casdoor.* TO 'casdoor'@'%';
FLUSH PRIVILEGES;
```

#### 4.3.2 导入初始数据

```bash
mysql -u casdoor -p casdoor < sql/init.sql
```

### 4.4 配置应用

#### 4.4.1 配置文件

编辑 `conf/app.conf` 文件：

```ini
appname = casdoor
httpport = 8000
runmode = prod
copyrequestbody = true

# 数据库配置
driverName = mysql
dataSourceName = casdoor:your_password@tcp(localhost:3306)/casdoor?charset=utf8mb4&parseTime=true&loc=Local

# 缓存配置
cacheEnabled = true
cacheTTL = 600

# 日志配置
logLevel = info
logFile = logs/casdoor.log

# 安全配置
authState = "casdoor"
httpProxy = ""
```

#### 4.4.2 环境变量配置

创建 `.env` 文件：

```bash
# 数据库配置
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_NAME=casdoor
DB_USER=casdoor
DB_PASSWORD=your_password

# 缓存配置
CACHE_ENABLED=true
CACHE_TTL=600

# 日志配置
LOG_LEVEL=info
LOG_FILE=logs/casdoor.log

# 服务配置
SERVER_PORT=8000
RUN_MODE=prod
```

### 4.5 编译部署

#### 4.5.1 编译

```bash
# 安装依赖
go mod download

# 编译
go build -o casdoor .

# 或者使用 Makefile
make build
```

#### 4.5.2 运行

```bash
# 直接运行
./casdoor

# 或使用 systemd 服务
sudo cp scripts/casdoor.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl start casdoor
sudo systemctl enable casdoor
```

### 4.6 配置 Nginx

#### 4.6.1 安装 Nginx

```bash
# Ubuntu
sudo apt install -y nginx

# CentOS
sudo yum install -y nginx
```

#### 4.6.2 配置反向代理

创建 `/etc/nginx/sites-available/casdoor`：

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 重定向到 HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    # SSL 配置
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # 日志配置
    access_log /var/log/nginx/casdoor_access.log;
    error_log /var/log/nginx/casdoor_error.log;

    # 反向代理配置
    location / {
        proxy_pass http://127.0.0.1:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket 支持
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";

        # 超时配置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # 静态文件缓存
    location /static/ {
        proxy_pass http://127.0.0.1:8000;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
}
```

#### 4.6.3 启用配置

```bash
# 创建软链接
sudo ln -s /etc/nginx/sites-available/casdoor /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重载 Nginx
sudo systemctl reload nginx
```

## 5. Docker 部署

### 5.1 Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o casdoor .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/casdoor .
COPY --from=builder /app/conf ./conf

EXPOSE 8000

CMD ["./casdoor"]
```

### 5.2 Docker Compose

```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: casdoor-mysql
    environment:
      MYSQL_ROOT_PASSWORD: root_password
      MYSQL_DATABASE: casdoor
      MYSQL_USER: casdoor
      MYSQL_PASSWORD: your_password
    volumes:
      - mysql_data:/var/lib/mysql
      - ./sql/init.sql:/docker-entrypoint-initdb.d/init.sql
    ports:
      - "3306:3306"
    networks:
      - casdoor-network

  casdoor:
    build: .
    container_name: casdoor-app
    depends_on:
      - mysql
    environment:
      DB_DRIVER: mysql
      DB_HOST: mysql
      DB_PORT: 3306
      DB_NAME: casdoor
      DB_USER: casdoor
      DB_PASSWORD: your_password
    ports:
      - "8000:8000"
    networks:
      - casdoor-network

  nginx:
    image: nginx:latest
    container_name: casdoor-nginx
    depends_on:
      - casdoor
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    ports:
      - "80:80"
      - "443:443"
    networks:
      - casdoor-network

volumes:
  mysql_data:

networks:
  casdoor-network:
    driver: bridge
```

### 5.3 部署命令

```bash
# 构建并启动
docker-compose up -d

# 查看日志
docker-compose logs -f casdoor

# 停止服务
docker-compose down
```

## 6. Kubernetes 部署

### 6.1 ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: casdoor-config
data:
  app.conf: |
    appname = casdoor
    httpport = 8000
    runmode = prod
    driverName = mysql
    dataSourceName = casdoor:$(DB_PASSWORD)@tcp($(DB_HOST):3306)/casdoor?charset=utf8mb4&parseTime=true&loc=Local
```

### 6.2 Secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: casdoor-secret
type: Opaque
stringData:
  db-password: your_password
```

### 6.3 Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: casdoor
spec:
  replicas: 3
  selector:
    matchLabels:
      app: casdoor
  template:
    metadata:
      labels:
        app: casdoor
    spec:
      containers:
      - name: casdoor
        image: casdoor:latest
        ports:
        - containerPort: 8000
        env:
        - name: DB_HOST
          value: mysql-service
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: casdoor-secret
              key: db-password
        resources:
          requests:
            cpu: 500m
            memory: 512Mi
          limits:
            cpu: 1000m
            memory: 1Gi
        livenessProbe:
          httpGet:
            path: /api/health
            port: 8000
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /api/health
            port: 8000
          initialDelaySeconds: 5
          periodSeconds: 5
```

### 6.4 Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: casdoor-service
spec:
  selector:
    app: casdoor
  ports:
  - port: 80
    targetPort: 8000
  type: LoadBalancer
```

## 7. 运维指南

### 7.1 服务管理

```bash
# 启动服务
sudo systemctl start casdoor

# 停止服务
sudo systemctl stop casdoor

# 重启服务
sudo systemctl restart casdoor

# 查看状态
sudo systemctl status casdoor

# 查看日志
sudo journalctl -u casdoor -f
```

### 7.2 日志管理

#### 7.2.1 日志配置

```ini
# conf/app.conf
logLevel = info
logFile = logs/casdoor.log
logMaxDays = 30
logMaxSize = 100  # MB
```

#### 7.2.2 日志轮转

创建 `/etc/logrotate.d/casdoor`：

```
/var/log/casdoor/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 0640 casdoor casdoor
    sharedscripts
    postrotate
        systemctl reload casdoor > /dev/null 2>&1 || true
    endscript
}
```

### 7.3 监控配置

#### 7.3.1 Prometheus 指标

应用已内置 Prometheus 指标，可通过 `/metrics` 端点访问。

#### 7.3.2 Grafana 仪表板

导入提供的 Grafana 仪表板配置文件 `monitoring/grafana-dashboard.json`。

### 7.4 备份策略

#### 7.4.1 数据库备份

```bash
#!/bin/bash
# backup.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backup/mysql"
MYSQL_USER="casdoor"
MYSQL_PASSWORD="your_password"
DATABASE="casdoor"

mkdir -p $BACKUP_DIR

mysqldump -u$MYSQL_USER -p$MYSQL_PASSWORD $DATABASE | gzip > $BACKUP_DIR/casdoor_$DATE.sql.gz

# 删除 7 天前的备份
find $BACKUP_DIR -name "*.sql.gz" -mtime +7 -delete
```

#### 7.4.2 定时备份

```bash
# 添加到 crontab
crontab -e

# 每天凌晨 2 点执行备份
0 2 * * * /path/to/backup.sh
```

### 7.5 性能调优

#### 7.5.1 数据库优化

```sql
-- MySQL 配置优化
SET GLOBAL innodb_buffer_pool_size = 4G;
SET GLOBAL innodb_log_file_size = 512M;
SET GLOBAL max_connections = 500;
SET GLOBAL query_cache_size = 256M;
```

#### 7.5.2 应用优化

```ini
# conf/app.conf

# 连接池配置
dbMaxIdleConns = 20
dbMaxOpenConns = 100
dbConnMaxLifetime = 3600

# 缓存配置
cacheEnabled = true
cacheTTL = 600
cacheMaxSize = 10000
```

## 8. 故障排查

### 8.1 常见问题

#### 8.1.1 服务无法启动

```bash
# 检查端口占用
netstat -tlnp | grep 8000

# 检查日志
tail -f logs/casdoor.log

# 检查配置
./casdoor -check-config
```

#### 8.1.2 数据库连接失败

```bash
# 测试数据库连接
mysql -h localhost -u casdoor -p casdoor

# 检查防火墙
sudo ufw status
sudo ufw allow 3306
```

#### 8.1.3 性能问题

```bash
# 查看系统资源
top
htop

# 查看数据库慢查询
mysql -e "SHOW VARIABLES LIKE 'slow_query%';"

# 分析性能瓶颈
go tool pprof http://localhost:8000/debug/pprof/profile
```

### 8.2 健康检查

```bash
# API 健康检查
curl http://localhost:8000/api/health

# 数据库健康检查
curl http://localhost:8000/api/health/db

# 缓存健康检查
curl http://localhost:8000/api/health/cache
```

## 9. 升级指南

### 9.1 升级步骤

```bash
# 1. 备份数据库
mysqldump -u casdoor -p casdoor > backup.sql

# 2. 停止服务
sudo systemctl stop casdoor

# 3. 更新代码
git pull origin main

# 4. 更新依赖
go mod download

# 5. 数据库迁移
mysql -u casdoor -p casdoor < sql/migration.sql

# 6. 重新编译
go build -o casdoor .

# 7. 启动服务
sudo systemctl start casdoor

# 8. 验证升级
curl http://localhost:8000/api/version
```

### 9.2 回滚方案

```bash
# 1. 停止服务
sudo systemctl stop casdoor

# 2. 恢复数据库
mysql -u casdoor -p casdoor < backup.sql

# 3. 切换到旧版本
git checkout v1.0.0

# 4. 重新编译
go build -o casdoor .

# 5. 启动服务
sudo systemctl start casdoor
```

## 10. 安全加固

### 10.1 系统安全

```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 配置防火墙
sudo ufw enable
sudo ufw allow 22
sudo ufw allow 80
sudo ufw allow 443

# 禁用 root 登录
sudo sed -i 's/PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config
sudo systemctl reload sshd
```

### 10.2 应用安全

```ini
# conf/app.conf

# 禁用调试模式
runmode = prod

# 配置 CORS
corsAllowOrigins = https://your-domain.com
corsAllowMethods = GET,POST,PUT,DELETE
corsAllowHeaders = Content-Type,Authorization

# 配置安全头
securityHeaders = true
xssProtection = true
contentTypeNosniff = true
frameDeny = true
```

### 10.3 数据库安全

```sql
-- 删除匿名用户
DELETE FROM mysql.user WHERE User='';

-- 禁止 root 远程登录
DELETE FROM mysql.user WHERE User='root' AND Host NOT IN ('localhost', '127.0.0.1', '::1');

-- 删除测试数据库
DROP DATABASE IF EXISTS test;

-- 刷新权限
FLUSH PRIVILEGES;
```

## 11. 联系支持

如有问题，请联系：
- 技术支持邮箱：support@example.com
- 文档地址：https://docs.example.com
- GitHub Issues：https://github.com/casdoor/casdoor/issues
