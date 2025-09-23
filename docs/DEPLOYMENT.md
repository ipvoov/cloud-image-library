# 部署指南

## 部署方式

本项目支持多种部署方式，包括Docker部署、Kubernetes部署和传统服务器部署。

## 1. Docker 部署

### 1.1 使用 Docker Compose（推荐）

**步骤1: 准备配置文件**

创建 `docker-compose.yml` 文件：

```yaml
version: '3.8'

services:
  # MySQL 数据库
  mysql:
    image: mysql:8.0
    container_name: cloud-picture-mysql
    environment:
      MYSQL_ROOT_PASSWORD: root123
      MYSQL_DATABASE: cloud_picture
      MYSQL_USER: cloud_user
      MYSQL_PASSWORD: cloud_pass
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./cloud_picture.sql:/docker-entrypoint-initdb.d/init.sql
    command: --default-authentication-plugin=mysql_native_password

  # Redis 缓存
  redis:
    image: redis:7-alpine
    container_name: cloud-picture-redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  # 后端服务
  backend:
    build:
      context: .
      dockerfile: manifest/docker/Dockerfile
    container_name: cloud-picture-backend
    ports:
      - "8123:8123"
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_NAME=cloud_picture
      - DB_USER=cloud_user
      - DB_PASSWORD=cloud_pass
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    depends_on:
      - mysql
      - redis
    volumes:
      - ./manifest/config:/app/config

  # 前端服务
  frontend:
    build:
      context: ./picture-frontend
      dockerfile: Dockerfile
    container_name: cloud-picture-frontend
    ports:
      - "80:80"
    depends_on:
      - backend

volumes:
  mysql_data:
  redis_data:
```

**步骤2: 启动服务**

```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f backend
```

### 1.2 单独构建镜像

**构建后端镜像:**

```dockerfile
# manifest/docker/Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/manifest/config ./config

EXPOSE 8123
CMD ["./main"]
```

**构建前端镜像:**

```dockerfile
# picture-frontend/Dockerfile
FROM node:18-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

## 2. Kubernetes 部署

### 2.1 创建命名空间

```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: cloud-picture
```

### 2.2 配置 ConfigMap

```yaml
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cloud-picture-config
  namespace: cloud-picture
data:
  config.yaml: |
    server:
      address: ":8123"
    
    database:
      default:
        link: "mysql:cloud_user:cloud_pass@tcp(mysql-service:3306)/cloud_picture"
    
    redis:
      default:
        address: redis-service:6379
        db: 0
```

### 2.3 部署 MySQL

```yaml
# k8s/mysql.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mysql
  namespace: cloud-picture
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mysql
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
      - name: mysql
        image: mysql:8.0
        env:
        - name: MYSQL_ROOT_PASSWORD
          value: "root123"
        - name: MYSQL_DATABASE
          value: "cloud_picture"
        - name: MYSQL_USER
          value: "cloud_user"
        - name: MYSQL_PASSWORD
          value: "cloud_pass"
        ports:
        - containerPort: 3306
        volumeMounts:
        - name: mysql-storage
          mountPath: /var/lib/mysql
      volumes:
      - name: mysql-storage
        persistentVolumeClaim:
          claimName: mysql-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: mysql-service
  namespace: cloud-picture
spec:
  selector:
    app: mysql
  ports:
  - port: 3306
    targetPort: 3306
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mysql-pvc
  namespace: cloud-picture
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
```

### 2.4 部署 Redis

```yaml
# k8s/redis.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis
  namespace: cloud-picture
spec:
  replicas: 1
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
      - name: redis
        image: redis:7-alpine
        ports:
        - containerPort: 6379
---
apiVersion: v1
kind: Service
metadata:
  name: redis-service
  namespace: cloud-picture
spec:
  selector:
    app: redis
  ports:
  - port: 6379
    targetPort: 6379
```

### 2.5 部署后端服务

```yaml
# k8s/backend.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend
  namespace: cloud-picture
spec:
  replicas: 2
  selector:
    matchLabels:
      app: backend
  template:
    metadata:
      labels:
        app: backend
    spec:
      containers:
      - name: backend
        image: cloud-picture-backend:latest
        ports:
        - containerPort: 8123
        volumeMounts:
        - name: config
          mountPath: /app/config
      volumes:
      - name: config
        configMap:
          name: cloud-picture-config
---
apiVersion: v1
kind: Service
metadata:
  name: backend-service
  namespace: cloud-picture
spec:
  selector:
    app: backend
  ports:
  - port: 8123
    targetPort: 8123
  type: ClusterIP
```

### 2.6 部署前端服务

```yaml
# k8s/frontend.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: frontend
  namespace: cloud-picture
spec:
  replicas: 2
  selector:
    matchLabels:
      app: frontend
  template:
    metadata:
      labels:
        app: frontend
    spec:
      containers:
      - name: frontend
        image: cloud-picture-frontend:latest
        ports:
        - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: frontend-service
  namespace: cloud-picture
spec:
  selector:
    app: frontend
  ports:
  - port: 80
    targetPort: 80
  type: LoadBalancer
```

### 2.7 部署 Ingress

```yaml
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: cloud-picture-ingress
  namespace: cloud-picture
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - host: cloud-picture.example.com
    http:
      paths:
      - path: /api
        pathType: Prefix
        backend:
          service:
            name: backend-service
            port:
              number: 8123
      - path: /
        pathType: Prefix
        backend:
          service:
            name: frontend-service
            port:
              number: 80
```

## 3. 传统服务器部署

### 3.1 环境准备

```bash
# 安装 Go
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装 Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# 安装 MySQL
sudo apt-get update
sudo apt-get install mysql-server

# 安装 Redis
sudo apt-get install redis-server
```

### 3.2 数据库初始化

```bash
# 创建数据库
mysql -u root -p
CREATE DATABASE cloud_picture CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'cloud_user'@'localhost' IDENTIFIED BY 'cloud_pass';
GRANT ALL PRIVILEGES ON cloud_picture.* TO 'cloud_user'@'localhost';
FLUSH PRIVILEGES;

# 导入数据库结构
mysql -u cloud_user -p cloud_picture < cloud_picture.sql
```

### 3.3 后端部署

```bash
# 克隆代码
git clone https://github.com/your-username/cloud-picture-library.git
cd cloud-picture-library

# 安装依赖
go mod tidy

# 配置环境变量
export DB_HOST=localhost
export DB_PORT=3306
export DB_NAME=cloud_picture
export DB_USER=cloud_user
export DB_PASSWORD=cloud_pass
export REDIS_HOST=localhost
export REDIS_PORT=6379

# 构建并运行
go build -o cloud-picture main.go
./cloud-picture
```

### 3.4 前端部署

```bash
# 进入前端目录
cd picture-frontend

# 安装依赖
npm install

# 构建生产版本
npm run build

# 使用 Nginx 部署
sudo apt-get install nginx
sudo cp -r dist/* /var/www/html/
sudo systemctl restart nginx
```

### 3.5 Nginx 配置

```nginx
# /etc/nginx/sites-available/cloud-picture
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态文件
    location / {
        root /var/www/html;
        try_files $uri $uri/ /index.html;
    }

    # API 代理
    location /api {
        proxy_pass http://localhost:8123;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # WebSocket 代理
    location /ws {
        proxy_pass http://localhost:8123;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## 4. 生产环境优化

### 4.1 性能优化

**后端优化:**
- 启用数据库连接池
- 配置 Redis 缓存
- 使用 CDN 加速静态资源
- 启用 Gzip 压缩

**前端优化:**
- 启用代码分割
- 使用 CDN 加载第三方库
- 启用浏览器缓存
- 压缩静态资源

### 4.2 安全配置

**HTTPS 配置:**
```bash
# 使用 Let's Encrypt 免费证书
sudo apt-get install certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

**防火墙配置:**
```bash
# 只开放必要端口
sudo ufw allow 22    # SSH
sudo ufw allow 80    # HTTP
sudo ufw allow 443   # HTTPS
sudo ufw enable
```

### 4.3 监控和日志

**系统监控:**
```bash
# 安装监控工具
sudo apt-get install htop iotop nethogs

# 配置日志轮转
sudo nano /etc/logrotate.d/cloud-picture
```

**应用监控:**
- 使用 Prometheus + Grafana 监控
- 配置 ELK 日志分析
- 设置告警通知

## 5. 故障排除

### 5.1 常见问题

**数据库连接失败:**
```bash
# 检查 MySQL 服务状态
sudo systemctl status mysql

# 检查连接配置
mysql -h localhost -u cloud_user -p cloud_picture
```

**Redis 连接失败:**
```bash
# 检查 Redis 服务状态
sudo systemctl status redis

# 测试连接
redis-cli ping
```

**端口冲突:**
```bash
# 查看端口占用
sudo netstat -tlnp | grep :8123
sudo lsof -i :8123
```

### 5.2 日志查看

```bash
# 查看应用日志
tail -f /var/log/cloud-picture/app.log

# 查看 Nginx 日志
tail -f /var/log/nginx/access.log
tail -f /var/log/nginx/error.log

# 查看系统日志
journalctl -u cloud-picture -f
```

## 6. 备份和恢复

### 6.1 数据库备份

```bash
# 创建备份脚本
#!/bin/bash
BACKUP_DIR="/backup/mysql"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p $BACKUP_DIR

mysqldump -u cloud_user -p cloud_picture > $BACKUP_DIR/cloud_picture_$DATE.sql

# 保留最近7天的备份
find $BACKUP_DIR -name "*.sql" -mtime +7 -delete
```

### 6.2 文件备份

```bash
# 备份上传的图片文件
rsync -av /var/www/uploads/ /backup/uploads/

# 备份配置文件
cp -r /app/config /backup/config_$(date +%Y%m%d)
```

### 6.3 恢复流程

```bash
# 恢复数据库
mysql -u cloud_user -p cloud_picture < /backup/mysql/cloud_picture_20240101_120000.sql

# 恢复文件
rsync -av /backup/uploads/ /var/www/uploads/
```
