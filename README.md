# Pika 探针监控系统

<div align="center">

一个基于 Go + PostgreSQL 的实时探针监控系统

[快速开始](#快速开始) • [功能特性](#功能特性) • [文档](#文档) • [架构](#架构)

</div>

## 简介

Pika 是一个轻量级的探针监控系统，支持实时数据采集、存储和查询。系统采用 WebSocket 进行探针与服务端的通信，使用 PostgreSQL 存储时序数据，提供完整的 RESTful API 和用户管理功能。

## 快速开始

### 环境要求

- Docker 20.10+
- Docker Compose 1.29+

### 一键部署

#### 1. 下载配置文件

```bash
# 下载 docker-compose.yml 配置文件
curl -O https://raw.githubusercontent.com/dushixiang/pika/main/docker-compose.yml

# 或使用 wget
wget https://raw.githubusercontent.com/dushixiang/pika/main/docker-compose.yml
```

#### 2. 修改配置（可选）

编辑 `docker-compose.yml` 文件，根据需要修改以下配置：

- **数据库密码**：`POSTGRES_PASSWORD` 和 `DATABASE_POSTGRES_PASSWORD`（生产环境建议修改）
- **JWT 密钥**：`APP_JWT_SECRET`（必须修改为至少 32 位的随机字符串）

生成随机密钥：
```bash
openssl rand -base64 32
```

#### 3. 启动服务

```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f pika
```

#### 4. 访问服务

服务启动后，访问 http://localhost:8080

#### 5. 停止服务

```bash
# 停止服务
docker-compose stop

# 停止并删除容器
docker-compose down

# 停止并删除容器及数据卷
docker-compose down -v
```

### 生产环境部署建议

#### 1. 安全配置

- 修改默认的数据库密码
- 设置强随机的 JWT 密钥
- 使用 HTTPS 反向代理（如 Nginx）
- 限制数据库端口仅允许内部访问

#### 2. 数据持久化

数据库数据默认存储在 `./data/postgresql` 目录，请定期备份：

```bash
# 备份数据库
docker-compose exec postgresql pg_dump -U pika pika > backup.sql

# 恢复数据库
docker-compose exec -T postgresql psql -U pika pika < backup.sql
```

#### 3. 反向代理配置（Nginx 示例）

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket 支持
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

### 故障排查

#### 服务无法启动

```bash
# 查看详细日志
docker-compose logs -f

# 检查容器状态
docker-compose ps

# 重启服务
docker-compose restart
```

#### 数据库连接失败

- 确认 PostgreSQL 容器已启动且健康检查通过
- 检查数据库配置是否正确
- 查看数据库日志：`docker-compose logs postgresql`

#### 端口冲突

如果 8080 或 5432 端口被占用，修改 `docker-compose.yml` 中的端口映射：

```yaml
ports:
  - "8081:8080"  # 将 8080 改为其他端口
```