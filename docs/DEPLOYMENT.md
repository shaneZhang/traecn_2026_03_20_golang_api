# 部署文档

## 1. 环境要求

### 1.1 操作系统
| 操作系统 | 版本要求 | 推荐版本 |
|---------|---------|---------|
| Linux | Kernel 3.10+ | CentOS 7+/Ubuntu 18.04+/Debian 9+ |
| macOS | 10.14+ | macOS 13+ |
| Windows | 10+ (WSL2) | Windows 11 (WSL2) |

### 1.2 软件依赖
| 软件 | 版本要求 | 推荐版本 | 说明 |
|------|---------|---------|------|
| Go | 1.19+ | 1.21.x | 编译运行后端服务 |
| MySQL/MariaDB | 5.7+/10.2+ | MySQL 8.0 / MariaDB 10.6 | 数据存储 |
| Redis | 5.0+ | 7.x | 缓存、会话存储 |
| Nginx | 1.18+ | 1.24+ | 反向代理、静态文件 |
| Node.js | 16+ | 18.x | 前端构建（如需要） |

### 1.3 硬件配置

#### 开发环境
| 配置项 | 最低要求 | 推荐配置 |
|--------|---------|---------|
| CPU | 2 核 | 4 核 |
| 内存 | 4 GB | 8 GB |
| 磁盘 | 20 GB | 50 GB SSD |
| 网络 | - | 带宽 100 Mbps+ |

#### 生产环境
| 配置项 | 小型部署 | 中型部署 | 大型部署 |
|--------|---------|---------|---------|
| CPU | 4 核 | 8 核 | 16+ 核 |
| 内存 | 8 GB | 16 GB | 32+ GB |
| 磁盘 | 50 GB SSD | 100 GB SSD | 200+ GB SSD |
| 网络 | 带宽 100 Mbps | 带宽 500 Mbps | 带宽 1 Gbps+ |
| 实例数 | 1 | 2-3 | 5+ |

## 2. 部署架构

### 2.1 单机部署架构
```
┌─────────────────────────────────────────────────────────────────┐
│                     服务器 (Single Node)                        │
│                                                                 │
│  ┌──────────┐     ┌──────────┐     ┌──────────┐                │
│  │  Nginx   │────▶│  Casdoor│────▶│   MySQL  │                │
│  │ (Reverse │     │Backend   │     │ Database│                │
│  │  Proxy)  │     │ Service  │     │          │                │
│  └──────────┘     └──────────┘     └──────────┘                │
│                       │                                         │
│                       ▼                                         │
│  ┌──────────┐     ┌──────────┐                                  │
│  │ Frontend │────▶│   Redis  │                                  │
│  │ Static   │     │ (Cache)  │                                  │
│  │  Files   │     └──────────┘                                  │
│  └──────────┘                                                   │
└─────────────────────────────────────────────────────────────────┘
```

### 2.2 集群部署架构
```
┌─────────────────────────────────────────────────────────────────┐
│                     Load Balancer (LB)                          │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  Nginx / ALB / Cloud LB                                 │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┴─────────────────────┐
        │                                           │
        ▼                                           ▼
┌─────────────────────┐                   ┌─────────────────────┐
│   App Server 1      │                   │   App Server N      │
│  ┌───────────────┐  │                   │  ┌───────────────┐  │
│  │   Casdoor     │  │  . . .  N 台 . .  │  │   Casdoor     │  │
│  │  Backend      │  │                   │  │  Backend      │  │
│  └───────────────┘  │                   │  └───────────────┘  │
└─────────────────────┘                   └─────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
        ▼                     ▼                     ▼
┌───────────────┐   ┌───────────────┐   ┌───────────────┐
│   MySQL       │   │   Redis       │   │   File        │
│ (Master-Slave)│   │ (Cluster)     │   │   Storage     │
│   / Aurora    │   │ / Sentinel    │   │   (NAS/S3)    │
└───────────────┘   └───────────────┘   └───────────────┘
```

## 3. 部署步骤

### 3.1 前置准备

#### 3.1.1 系统配置
```bash
# 1. 更新系统包
sudo yum update -y          # CentOS/RHEL
sudo apt update && sudo apt upgrade -y  # Ubuntu/Debian

# 2. 安装基础工具
sudo yum install -y wget curl tar unzip git vim  # CentOS/RHEL
sudo apt install -y wget curl tar unzip git vim  # Ubuntu/Debian

# 3. 配置时区
sudo timedatectl set-timezone Asia/Shanghai
timedatectl status

# 4. 优化文件句柄限制
sudo tee /etc/security/limits.conf << EOF
* soft nofile 65536
* hard nofile 65536
* soft nproc 65536
* hard nproc 65536
EOF

# 5. 关闭防火墙（或配置规则）
sudo systemctl stop firewalld
sudo systemctl disable firewalld

# 或配置防火墙规则（推荐）：
# sudo firewall-cmd --add-port=80/tcp --permanent
# sudo firewall-cmd --add-port=443/tcp --permanent
# sudo firewall-cmd --reload
```

#### 3.1.2 安装 Go 环境
```bash
# 1. 下载 Go（请使用最新版本，以下为示例）
GO_VERSION="1.21.5"
wget https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz

# 2. 安装
sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz

# 3. 配置环境变量
tee -a ~/.bashrc << 'EOF'
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct
EOF

source ~/.bashrc

# 4. 验证安装
go version
go env
```

### 3.2 数据库配置

#### 3.2.1 MySQL 安装配置
```bash
# CentOS/RHEL 安装 MySQL 8.0
sudo yum install -y https://dev.mysql.com/get/mysql80-community-release-el7-3.noarch.rpm
sudo yum install -y mysql-community-server
sudo systemctl start mysqld
sudo systemctl enable mysqld

# Ubuntu/Debian 安装 MySQL
sudo apt install -y mysql-server

# 查看初始密码
sudo grep 'temporary password' /var/log/mysqld.log

# 安全配置
sudo mysql_secure_installation

# 登录 MySQL
mysql -u root -p

# 创建数据库和用户
CREATE DATABASE casdoor CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'casdoor'@'%' IDENTIFIED BY 'YourStrongPassword123!';
GRANT ALL PRIVILEGES ON casdoor.* TO 'casdoor'@'%';
FLUSH PRIVILEGES;

# 配置 MySQL 参数（生产环境）
sudo tee /etc/my.cnf.d/casdoor.cnf << EOF
[mysqld]
character-set-server = utf8mb4
collation-server = utf8mb4_unicode_ci
init_connect='SET NAMES utf8mb4'
max_connections = 1000
max_connect_errors = 100000
table_open_cache = 2048
sort_buffer_size = 2M
read_buffer_size = 2M
read_rnd_buffer_size = 8M
join_buffer_size = 8M
thread_cache_size = 128
query_cache_size = 0
query_cache_type = 0
tmp_table_size = 128M
max_heap_table_size = 128M
slow_query_log = 1
slow_query_log_file = /var/log/mysql-slow.log
long_query_time = 2
EOF

sudo systemctl restart mysqld
```

#### 3.2.2 Redis 安装配置（可选但推荐）
```bash
# 安装 Redis
sudo yum install -y redis  # CentOS
sudo apt install -y redis-server  # Ubuntu

# 配置 Redis
sudo tee /etc/redis.conf << EOF
bind 0.0.0.0
port 6379
requirepass YourRedisPassword
maxmemory 2gb
maxmemory-policy allkeys-lru
appendonly yes
appendfsync everysec
EOF

sudo systemctl start redis
sudo systemctl enable redis

# 验证
redis-cli -a YourRedisPassword ping
# 返回: PONG
```

### 3.3 应用部署

#### 3.3.1 代码获取与编译
```bash
# 1. 克隆代码
mkdir -p /opt/casdoor
cd /opt/casdoor
git clone https://github.com/casdoor/casdoor.git .

# 或使用您的私有仓库
# git clone git@your-repo-url:casdoor/casdoor.git .

# 2. 切换到目标版本（如需要）
# git checkout v1.0.0

# 3. 下载依赖
go mod download
go mod verify

# 4. 编译
go build -o casdoor .

# 或带版本信息编译
BUILD_TIME=$(date "+%Y-%m-%d %H:%M:%S")
COMMIT_ID=$(git rev-parse --short HEAD)
go build -ldflags "-X main.buildTime=${BUILD_TIME} -X main.commitId=${COMMIT_ID}" -o casdoor .

# 5. 验证编译结果
./casdoor --version
ls -lh casdoor
```

#### 3.3.2 配置文件
```bash
# 复制配置文件模板
cp conf/app.conf.example conf/app.conf

# 编辑配置文件
vim conf/app.conf
```

**关键配置项说明 (`conf/app.conf`)**:
```ini
# 应用基本配置
appname = casdoor
httpport = 8000
runmode = prod           # 生产环境设为 prod
autorender = false
copyrequestbody = true
sessionon = true
sessionname = casdoor-session
sessionprovider = file  # 生产环境推荐 redis
sessiongcmaxlifetime = 3600
sessionproviderconfig = "127.0.0.1:6379,100,YourRedisPassword"

# 数据库配置
driverName = mysql
dataSourceName = root:123456@tcp(localhost:3306)/casdoor?charset=utf8mb4&parseTime=true&loc=Local
dbName = casdoor
tableNamePrefix =
showSql = false          # 生产环境设为 false
maxIdleConns = 20
maxOpenConns = 100
connMaxLifetime = 300

# Redis 配置
redisEndpoint = 
redisUsername = 
redisPassword = 

# 日志配置
logPath = ./logs
logLevel = info         # 生产环境设为 info 或 warn

# 安全配置
tokenSecret = "YourJwtSecretKeyHereMakeItLongAndRandom"
corsDomains = "*"       # 生产环境限制为具体域名

# 邮箱配置（用于发送邮件）
smtpHost = smtp.example.com
smtpPort = 587
smtpUsername = no-reply@example.com
smtpPassword = your-smtp-password
smtpFrom = Casdoor <no-reply@example.com>

# 短信配置（用于发送短信）
smsProvider = 
smsAccessKeyId = 
smsAccessKeySecret = 
smsSign = 
smsTemplateCode = 

# 存储配置（用于头像、文件等）
storageProvider = 
storageEndpoint = 
storageAccessKeyId = 
storageAccessKeySecret = 
storageBucket = 
storagePathPrefix = 
storageHost = 
storageRegion = 
```

#### 3.3.3 初始化数据库
```bash
# 1. 检查数据库连接
./casdoor check-db

# 2. 初始化表结构
./casdoor init-db

# 3. 导入初始数据（如需要）
# ./casdoor import-data init-data.sql
```

#### 3.3.4 配置系统服务（Systemd）
```bash
# 创建服务用户
sudo useradd -r -s /sbin/nologin casdoor
sudo chown -R casdoor:casdoor /opt/casdoor

# 创建服务文件
sudo tee /etc/systemd/system/casdoor.service << EOF
[Unit]
Description=Casdoor Identity and Access Management Service
After=network.target mysqld.service redis.service
Wants=mysqld.service redis.service

[Service]
Type=simple
User=casdoor
Group=casdoor
WorkingDirectory=/opt/casdoor
ExecStart=/opt/casdoor/casdoor
Restart=on-failure
RestartSec=5
LimitNOFILE=65536
LimitNPROC=65536

# 环境变量（可选）
# Environment="GIN_MODE=release"
# Environment="CASDOOR_CONFIG_PATH=/etc/casdoor/app.conf"

# 安全加固
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/casdoor

[Install]
WantedBy=multi-user.target
EOF

# 重新加载 systemd
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start casdoor

# 查看状态
sudo systemctl status casdoor

# 查看日志
journalctl -u casdoor -f -n 100

# 设置开机自启
sudo systemctl enable casdoor
```

#### 3.3.5 验证服务启动
```bash
# 检查进程
ps -ef | grep casdoor

# 检查端口
netstat -tulpn | grep 8000
# 或
ss -tulpn | grep 8000

# 健康检查
curl -i http://localhost:8000/health

# 应返回类似:
# HTTP/1.1 200 OK
# {"status":"ok"}

# 查看日志
tail -f /opt/casdoor/logs/casdoor.log
```

### 3.4 Nginx 反向代理配置

#### 3.4.1 安装 Nginx
```bash
# CentOS/RHEL
sudo yum install -y nginx

# Ubuntu/Debian
sudo apt install -y nginx

# 启动并设置开机自启
sudo systemctl start nginx
sudo systemctl enable nginx
```

#### 3.4.2 配置反向代理
```bash
# 创建配置文件
sudo tee /etc/nginx/conf.d/casdoor.conf << 'EOF'
upstream casdoor_backend {
    server 127.0.0.1:8000;
    # 集群部署时添加多个后端:
    # server 192.168.1.101:8000;
    # server 192.168.1.102:8000;
    
    keepalive 32;
}

server {
    listen 80;
    server_name your-domain.com;  # 替换为你的域名
    
    # 安全头
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    add_header Content-Security-Policy "default-src 'self'";
    
    # Gzip 压缩
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;
    
    # 代理配置
    location / {
        proxy_pass http://casdoor_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket 支持（如需要）
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        
        # 超时配置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
        send_timeout 60s;
        
        # 缓冲配置
        proxy_buffering on;
        proxy_buffer_size 4k;
        proxy_buffers 8 4k;
    }
    
    # 静态文件（如前端部署在同一服务器）
    # location /static/ {
    #     root /opt/casdoor/web;
    #     expires 30d;
    #     add_header Cache-Control "public, max-age=2592000";
    # }
    
    # 健康检查路径
    location /health {
        access_log off;
        return 200 '{"status":"ok"}';
        add_header Content-Type application/json;
    }
    
    # 禁止访问隐藏文件
    location ~ /\.ht {
        deny all;
    }
}
EOF

# 测试配置
sudo nginx -t

# 重载配置
sudo systemctl reload nginx
```

#### 3.4.3 HTTPS 配置（推荐）
```bash
# 安装 certbot（用于 Let's Encrypt 证书）
sudo yum install -y certbot python3-certbot-nginx  # CentOS
sudo apt install -y certbot python3-certbot-nginx  # Ubuntu

# 获取并安装证书（自动配置）
sudo certbot --nginx -d your-domain.com

# 或手动配置（已有证书）
sudo tee /etc/nginx/conf.d/casdoor-ssl.conf << 'EOF'
server {
    listen 443 ssl;
    server_name your-domain.com;
    
    # SSL 证书路径
    ssl_certificate /path/to/your/fullchain.pem;
    ssl_certificate_key /path/to/your/privkey.pem;
    
    # SSL 配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    
    # HSTS
    add_header Strict-Transport-Security "max-age=63072000" always;
    
    # 其他配置与 HTTP 相同...
    location / {
        proxy_pass http://casdoor_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# HTTP 重定向到 HTTPS
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$host$request_uri;
}
EOF
```

## 4. 集群部署（可选）

### 4.1 会话共享配置
```ini
# conf/app.conf
sessionprovider = redis
sessionproviderconfig = "redis-host:6379,100,your-redis-password,casdoor:sess:"
```

### 4.2 负载均衡配置
在 Nginx 配置中添加多个 upstream 节点：
```nginx
upstream casdoor_backend {
    least_conn;  # 负载均衡策略: least_conn / ip_hash / round_robin
    
    server 192.168.1.101:8000 max_fails=3 fail_timeout=30s;
    server 192.168.1.102:8000 max_fails=3 fail_timeout=30s;
    server 192.168.1.103:8000 max_fails=3 fail_timeout=30s;
    
    keepalive 64;
}
```

## 5. Docker 部署

### 5.1 Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /build
COPY . .

RUN go mod download
RUN go build -o casdoor .

FROM alpine:3.18

RUN apk add --no-cache ca-certificates tzdata && \
    update-ca-certificates

WORKDIR /app

COPY --from=builder /build/casdoor .
COPY conf ./conf
COPY web ./web

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -q -O - http://localhost:8000/health || exit 1

EXPOSE 8000

CMD ["./casdoor"]
```

### 5.2 Docker Compose
```yaml
version: '3.8'

services:
  casdoor:
    build: .
    ports:
      - "8000:8000"
    environment:
      - RUNMODE=prod
      - DRIVERNAME=mysql
      - DATASOURCENAME=casdoor:casdoor123@tcp(db:3306)/casdoor?charset=utf8mb4&parseTime=true&loc=Local
    volumes:
      - ./conf:/app/conf
      - ./logs:/app/logs
    depends_on:
      - db
      - redis
    restart: unless-stopped

  db:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root123
      MYSQL_DATABASE: casdoor
      MYSQL_USER: casdoor
      MYSQL_PASSWORD: casdoor123
    volumes:
      - mysql_data:/var/lib/mysql
    command: --default-authentication-plugin=mysql_native_password --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass redis123
    volumes:
      - redis_data:/data
    restart: unless-stopped

volumes:
  mysql_data:
  redis_data:
```

## 6. 运维管理

### 6.1 服务管理
```bash
# 启动服务
sudo systemctl start casdoor

# 停止服务
sudo systemctl stop casdoor

# 重启服务
sudo systemctl restart casdoor

# 查看状态
sudo systemctl status casdoor

# 查看日志（实时）
journalctl -u casdoor -f

# 查看最近 100 行日志
journalctl -u casdoor -n 100

# 按时间范围查看日志
journalctl -u casdoor --since "2024-01-01" --until "2024-01-02"
```

### 6.2 日志管理
```bash
# 日志位置
/opt/casdoor/logs/

# 配置日志轮转
sudo tee /etc/logrotate.d/casdoor << EOF
/opt/casdoor/logs/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 casdoor casdoor
    postrotate
        systemctl reload casdoor > /dev/null 2>&1 || true
    endscript
}
EOF
```

### 6.3 备份与恢复

#### 数据库备份
```bash
# 全量备份
mysqldump -u casdoor -p'YourPassword' casdoor | gzip > /backup/casdoor-$(date +%Y%m%d_%H%M%S).sql.gz

# 定时备份（添加到 crontab）
0 2 * * * /usr/bin/mysqldump -u casdoor -p'YourPassword' casdoor | gzip > /backup/casdoor-$(date +\%Y\%m\%d_\%H\%M\%S).sql.gz

# 保留最近 30 天的备份
find /backup -name "casdoor-*.sql.gz" -mtime +30 -delete
```

#### 数据恢复
```bash
# 停止服务
sudo systemctl stop casdoor

# 恢复数据库
gunzip < /backup/casdoor-20240101_020000.sql.gz | mysql -u casdoor -p'YourPassword' casdoor

# 启动服务
sudo systemctl start casdoor
```

## 7. 性能调优

### 7.1 JVM/Go 运行时调优
```bash
# 设置 Go 环境变量（可在 service 文件中配置）
GOGC=100        # 触发 GC 的内存增长百分比
GOMAXPROCS=8    # 可使用的 CPU 核心数（通常设为物理核心数）
GODEBUG=gctrace=1  # 启用 GC 调试日志
```

### 7.2 数据库调优
参考 `3.2.1` 节的 MySQL 配置建议。

### 7.3 应用级调优
```ini
# conf/app.conf
maxIdleConns = 50        # 根据连接数调整
maxOpenConns = 200       # 根据并发数调整
connMaxLifetime = 300    # 连接生命周期（秒）
```

## 8. 故障排查

### 8.1 常见问题

#### 问题1：服务启动失败
**排查步骤**:
1. 检查配置文件格式：`./casdoor check-config`
2. 检查数据库连接：`mysql -u user -p -h host`
3. 查看日志：`journalctl -u casdoor -n 100`
4. 检查端口占用：`netstat -tulpn | grep 8000`

#### 问题2：数据库连接失败
**排查步骤**:
1. 检查 MySQL 服务状态：`systemctl status mysqld`
2. 验证连接参数：确认 `conf/app.conf` 中的 `dataSourceName`
3. 测试网络连通性：`telnet mysql-host 3306`
4. 检查用户权限和防火墙

#### 问题3：性能问题
**排查步骤**:
1. 查看系统资源：`top`, `htop`, `iostat`, `vmstat`
2. 检查慢查询：开启 MySQL 慢查询日志分析
3. 分析 Goroutine 堆栈：`curl http://localhost:8000/debug/pprof/goroutine?debug=2`
4. 启用 pprof 分析（开发环境）

### 8.2 紧急恢复
```bash
# 服务挂掉重启
sudo systemctl restart casdoor

# 回滚到上一版本（如已备份）
# 1. 停止服务
sudo systemctl stop casdoor

# 2. 恢复二进制和配置
cp /opt/backup/casdoor.old /opt/casdoor/casdoor
cp /opt/backup/app.conf.old /opt/casdoor/conf/app.conf

# 3. 启动服务
sudo systemctl start casdoor
```

---

**文档版本**: v1.0  
**最后更新**: 2024年  
**适用版本**: v2.x
