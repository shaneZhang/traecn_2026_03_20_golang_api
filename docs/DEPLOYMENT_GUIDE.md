# 部署文档

## 1. 概述

本文档提供了重构后系统的部署指南，包括环境准备、配置、部署步骤和运维监控。

## 2. 环境要求

### 2.1 系统要求

| 组件 | 最低配置 | 推荐配置 |
|------|----------|----------|
| CPU | 2 核 | 4 核+ |
| 内存 | 4GB | 8GB+ |
| 磁盘 | 20GB | 100GB+ SSD |
| 网络 | 10Mbps | 100Mbps+ |

### 2.2 软件依赖

| 软件 | 版本 | 说明 |
|------|------|------|
| Go | 1.21+ | 运行环境 |
| MySQL | 8.0+ | 数据库 |
| Redis | 7.0+ | 缓存（可选） |
| Nginx | 1.20+ | 反向代理 |

## 3. 环境准备

### 3.1 安装 Go

```bash
# 下载并安装 Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# 配置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPROXY=https://goproxy.cn,direct' >> ~/.bashrc
source ~/.bashrc

# 验证安装
go version
```

### 3.2 安装 MySQL

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install mysql-server-8.0

# CentOS/RHEL
sudo yum install mysql-server

# 启动服务
sudo systemctl start mysql
sudo systemctl enable mysql

# 安全配置
sudo mysql_secure_installation
```

### 3.3 创建数据库

```bash
# 登录 MySQL
mysql -u root -p

# 创建数据库
CREATE DATABASE casdoor CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 创建用户
CREATE USER 'casdoor'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON casdoor.* TO 'casdoor'@'localhost';
FLUSH PRIVILEGES;

# 退出
EXIT;
```

### 3.4 安装 Redis（可选）

```bash
# Ubuntu/Debian
sudo apt-get install redis-server

# CentOS/RHEL
sudo yum install redis

# 启动服务
sudo systemctl start redis
sudo systemctl enable redis
```

## 4. 应用部署

### 4.1 代码部署

```bash
# 克隆代码
git clone https://github.com/casdoor/casdoor.git
cd casdoor

# 切换到重构分支（如果有）
git checkout refactor/api-modules

# 下载依赖
go mod download
go mod verify

# 编译
GOOS=linux GOARCH=amd64 go build -o casdoor main.go
```

### 4.2 配置文件

创建 `conf/app.conf`：

```ini
appname = casdoor
httpport = 8000
runmode = prod

# 数据库配置
driverName = mysql
dataSourceName = casdoor:your_password@tcp(localhost:3306)/casdoor

# Redis 配置（可选）
redisEndpoint = localhost:6379

# 会话配置
sessionConfig = ""
sessionProvider = ""
sessionSavePath = ""

# 日志配置
logConfig = {"filename": "logs/casdoor.log", "level": "info", "maxdays": 30}

# 静态文件
staticBaseUrl = ""

# 其他配置
isDemoMode = false
authConfig = ""
```

### 4.3 数据库迁移

```bash
# 运行数据库优化脚本
mysql -u casdoor -p casdoor < scripts/database_optimization.sql

# 或使用 Xorm 自动迁移
go run main.go migrate
```

### 4.4 启动服务

#### 4.4.1 直接启动

```bash
# 开发模式
./casdoor

# 生产模式（后台运行）
nohup ./casdoor > logs/casdoor.log 2>&1 &

# 查看日志
tail -f logs/casdoor.log
```

#### 4.4.2 Systemd 服务

创建 `/etc/systemd/system/casdoor.service`：

```ini
[Unit]
Description=Casdoor Identity Platform
After=network.target mysql.service redis.service

[Service]
Type=simple
User=casdoor
Group=casdoor
WorkingDirectory=/opt/casdoor
ExecStart=/opt/casdoor/casdoor
Restart=always
RestartSec=5
Environment="GO_ENV=production"

# 资源限制
LimitNOFILE=65536
LimitNPROC=4096

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
# 创建用户
sudo useradd -r -s /bin/false casdoor

# 设置权限
sudo mkdir -p /opt/casdoor
sudo cp -r . /opt/casdoor
sudo chown -R casdoor:casdoor /opt/casdoor

# 启动服务
sudo systemctl daemon-reload
sudo systemctl enable casdoor
sudo systemctl start casdoor

# 查看状态
sudo systemctl status casdoor
sudo journalctl -u casdoor -f
```

## 5. Nginx 配置

### 5.1 反向代理配置

创建 `/etc/nginx/sites-available/casdoor`：

```nginx
upstream casdoor {
    server 127.0.0.1:8000;
    keepalive 32;
}

server {
    listen 80;
    server_name auth.example.com;
    
    # 重定向到 HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name auth.example.com;
    
    # SSL 证书
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    
    # 日志
    access_log /var/log/nginx/casdoor_access.log;
    error_log /var/log/nginx/casdoor_error.log;
    
    # 静态文件
    location /static/ {
        alias /opt/casdoor/web/build/;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
    
    # API 代理
    location /api/ {
        proxy_pass http://casdoor;
        proxy_http_version 1.1;
        proxy_set_header Connection "";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 超时设置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
        
        # 缓冲区
        proxy_buffering on;
        proxy_buffer_size 4k;
        proxy_buffers 8 4k;
    }
    
    # 前端路由
    location / {
        root /opt/casdoor/web/build;
        try_files $uri $uri/ /index.html;
    }
}
```

启用配置：

```bash
sudo ln -s /etc/nginx/sites-available/casdoor /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 5.2 负载均衡配置（多实例）

```nginx
upstream casdoor {
    least_conn;
    
    server 10.0.1.10:8000 weight=5;
    server 10.0.1.11:8000 weight=5;
    server 10.0.1.12:8000 backup;
    
    keepalive 32;
}
```

## 6. Docker 部署

### 6.1 Dockerfile

```dockerfile
# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o casdoor main.go

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/casdoor .
COPY --from=builder /app/conf ./conf

EXPOSE 8000

CMD ["./casdoor"]
```

### 6.2 Docker Compose

创建 `docker-compose.yml`：

```yaml
version: '3.8'

services:
  casdoor:
    build: .
    ports:
      - "8000:8000"
    environment:
      - DB_DRIVER=mysql
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_NAME=casdoor
      - DB_USER=casdoor
      - DB_PASSWORD=your_password
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    depends_on:
      - mysql
      - redis
    volumes:
      - ./conf:/root/conf
      - ./logs:/root/logs
    restart: unless-stopped
    
  mysql:
    image: mysql:8.0
    environment:
      - MYSQL_ROOT_PASSWORD=root_password
      - MYSQL_DATABASE=casdoor
      - MYSQL_USER=casdoor
      - MYSQL_PASSWORD=your_password
    volumes:
      - mysql_data:/var/lib/mysql
      - ./scripts/database_optimization.sql:/docker-entrypoint-initdb.d/01-optimization.sql
    restart: unless-stopped
    
  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
    restart: unless-stopped
    
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - casdoor
    restart: unless-stopped

volumes:
  mysql_data:
  redis_data:
```

部署：

```bash
# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f casdoor

# 停止服务
docker-compose down

# 更新部署
docker-compose pull
docker-compose up -d
```

## 7. Kubernetes 部署

### 7.1 ConfigMap

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
    dataSourceName = casdoor:password@tcp(mysql:3306)/casdoor
```

### 7.2 Deployment

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
        image: casdoor/casdoor:latest
        ports:
        - containerPort: 8000
        env:
        - name: GO_ENV
          value: "production"
        volumeMounts:
        - name: config
          mountPath: /root/conf
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "2Gi"
            cpu: "2000m"
        livenessProbe:
          httpGet:
            path: /api/health
            port: 8000
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /api/ready
            port: 8000
          initialDelaySeconds: 5
          periodSeconds: 5
      volumes:
      - name: config
        configMap:
          name: casdoor-config
```

### 7.3 Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: casdoor
spec:
  selector:
    app: casdoor
  ports:
  - port: 80
    targetPort: 8000
  type: ClusterIP
```

### 7.4 Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: casdoor
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt
spec:
  tls:
  - hosts:
    - auth.example.com
    secretName: casdoor-tls
  rules:
  - host: auth.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: casdoor
            port:
              number: 80
```

部署：

```bash
kubectl apply -f k8s/
```

## 8. 监控与告警

### 8.1 Prometheus 监控

```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'casdoor'
    static_configs:
      - targets: ['localhost:8000']
    metrics_path: /metrics
```

### 8.2 Grafana 仪表板

关键指标：
- HTTP 请求率
- 响应时间分布
- 错误率
- 数据库连接池状态
- 缓存命中率

### 8.3 告警规则

```yaml
# alert.rules
groups:
- name: casdoor
  rules:
  - alert: HighErrorRate
    expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "High error rate detected"
      
  - alert: SlowResponse
    expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 0.5
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "Slow response time detected"
      
  - alert: DatabaseConnectionsHigh
    expr: db_connections_active / db_connections_max > 0.8
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "Database connection pool near limit"
```

## 9. 备份与恢复

### 9.1 数据库备份

```bash
#!/bin/bash
# backup.sh

BACKUP_DIR="/backup/casdoor"
DATE=$(date +%Y%m%d_%H%M%S)

# MySQL 备份
mysqldump -u casdoor -p casdoor > $BACKUP_DIR/casdoor_$DATE.sql

# 压缩
gzip $BACKUP_DIR/casdoor_$DATE.sql

# 保留最近 7 天
find $BACKUP_DIR -name "casdoor_*.sql.gz" -mtime +7 -delete
```

设置定时任务：

```bash
# 每天凌晨 2 点备份
0 2 * * * /opt/casdoor/scripts/backup.sh
```

### 9.2 数据恢复

```bash
# 解压备份
gunzip casdoor_20240101_020000.sql.gz

# 恢复数据
mysql -u casdoor -p casdoor < casdoor_20240101_020000.sql
```

## 10. 故障排查

### 10.1 常见问题

#### 10.1.1 服务无法启动

```bash
# 检查日志
journalctl -u casdoor -n 100

# 检查端口占用
netstat -tlnp | grep 8000

# 检查配置文件
cat conf/app.conf

# 检查数据库连接
mysql -u casdoor -p -h localhost casdoor
```

#### 10.1.2 数据库连接失败

```bash
# 检查 MySQL 状态
systemctl status mysql

# 检查连接数
mysql -e "SHOW STATUS LIKE 'Threads_connected';"

# 检查最大连接数
mysql -e "SHOW VARIABLES LIKE 'max_connections';"

# 增加连接数
mysql -e "SET GLOBAL max_connections = 500;"
```

#### 10.1.3 性能问题

```bash
# 查看慢查询日志
tail -f /var/log/mysql/slow.log

# 分析性能瓶颈
go tool pprof http://localhost:8000/debug/pprof/profile

# 查看 Goroutine
curl http://localhost:8000/debug/pprof/goroutine?debug=1
```

### 10.2 调试模式

```bash
# 启用调试日志
export LOG_LEVEL=debug
./casdoor

# 启用性能分析
export PPROF_ENABLED=true
./casdoor
```

## 11. 升级指南

### 11.1 升级前准备

```bash
# 备份数据
./scripts/backup.sh

# 检查当前版本
./casdoor version

# 查看更新日志
cat CHANGELOG.md
```

### 11.2 升级步骤

```bash
# 1. 停止服务
sudo systemctl stop casdoor

# 2. 备份当前版本
mv /opt/casdoor /opt/casdoor_backup

# 3. 部署新版本
git pull origin main
go build -o casdoor main.go
sudo cp -r . /opt/casdoor

# 4. 数据库迁移
./casdoor migrate

# 5. 启动服务
sudo systemctl start casdoor

# 6. 验证
./scripts/health_check.sh
```

### 11.3 回滚

```bash
# 停止服务
sudo systemctl stop casdoor

# 恢复备份
rm -rf /opt/casdoor
mv /opt/casdoor_backup /opt/casdoor

# 恢复数据库
mysql -u casdoor -p casdoor < backup/casdoor_xxx.sql

# 启动服务
sudo systemctl start casdoor
```

## 12. 安全加固

### 12.1 系统安全

```bash
# 防火墙配置
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

# SELinux（CentOS）
setsebool -P httpd_can_network_connect 1

# 文件权限
chmod 600 conf/app.conf
chmod 700 logs
```

### 12.2 应用安全

- 定期更新依赖包
- 启用 HTTPS
- 配置安全响应头
- 限制请求频率
- 启用审计日志

## 13. 附录

### 13.1 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| GO_ENV | 运行环境 | development |
| DB_DRIVER | 数据库驱动 | mysql |
| DB_HOST | 数据库主机 | localhost |
| DB_PORT | 数据库端口 | 3306 |
| REDIS_HOST | Redis 主机 | localhost |
| LOG_LEVEL | 日志级别 | info |
| PPROF_ENABLED | 启用性能分析 | false |

### 13.2 常用命令

```bash
# 查看版本
./casdoor version

# 数据库迁移
./casdoor migrate

# 生成配置文件模板
./casdoor config init

# 健康检查
./casdoor health

# 性能分析
./casdoor pprof
```

### 13.3 参考链接

- [官方文档](https://casdoor.org/)
- [API 文档](https://casdoor.org/docs/api/)
- [GitHub](https://github.com/casdoor/casdoor)
