# 智能云图库 (Cloud Picture Library)

一个基于 Go + Vue3 的现代化图片管理系统，支持图片上传、管理、协同编辑和空间管理等功能。

## 🌟 项目特色

- **🎨 智能图片管理**: 支持多种图片格式，自动提取图片信息和颜色特征
- **👥 团队协作**: 支持团队空间和实时协同编辑
- **🔐 权限管理**: 细粒度的用户权限控制
- **📊 数据分析**: 丰富的空间使用分析和统计功能
- **🌐 现代化UI**: 基于 Ant Design Vue 的美观界面
- **⚡ 高性能**: Go 后端 + Vue3 前端，响应迅速

## 🛠️ 技术栈

### 后端技术
- **Go 1.23+**: 主要开发语言
- **GoFrame v2**: Web框架和ORM
- **MySQL 8.0+**: 主数据库
- **Redis**: 缓存和会话存储
- **WebSocket**: 实时通信
- **腾讯云COS**: 对象存储服务
- **火山引擎**: AI服务集成
```

## 🚀 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/your-username/cloud-picture-library.git
cd cloud-picture-library
```

### 2. 数据库配置

```bash
# 创建数据库
mysql -u root -p
CREATE DATABASE cloud_picture CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 导入数据库结构
mysql -u root -p cloud_picture < cloud_picture.sql
```

### 3. 后端配置

```bash
# 复制配置文件
cp manifest/config/config.yaml.example manifest/config/config.yaml

# 编辑配置文件，修改数据库连接信息
vim manifest/config/config.yaml
```

配置文件示例：
```yaml
server:
  address: ":8123"

database:
  default:
    link: "mysql:root:password@tcp(localhost:3306)/cloud_picture"

redis:
  default:
    address: 127.0.0.1:6379
    db: 0
```

### 4. 启动后端服务

```bash
# 安装依赖
go mod tidy

# 启动服务
go run main.go
```

### 5. 启动前端服务

```bash
cd picture-frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```
```

## 📊 功能截图

### 首页
![首页](docs/images/homepage.png)

### 图片管理
![图片管理](docs/images/picture_manage.png)

### 空间分析
![空间分析](docs/images/spece.png)

### 协同编辑
![协同编辑](docs/images/edit.png)


⭐ 如果这个项目对您有帮助，请给我们一个星标！
