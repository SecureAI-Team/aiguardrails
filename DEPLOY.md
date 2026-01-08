# 阿里云 ECS 部署指南

## 快速部署 (一键安装)

在 Ubuntu 24.04 / 22.04 LTS 服务器上执行：

```bash
# 方式1: 本地执行
git clone https://github.com/your-org/aiguardrails.git /opt/aiguardrails
cd /opt/aiguardrails
bash install.sh

# 方式2: 远程执行 (需要设置REPO_URL)
REPO_URL=https://github.com/your-org/aiguardrails.git bash install.sh
```

## 系统要求

| 项目 | 最低配置 | 推荐配置 |
|------|---------|---------|
| CPU | 2核 | 4核+ |
| 内存 | 4GB | 8GB+ |
| 磁盘 | 40GB | 100GB SSD |
| 系统 | Ubuntu 22.04 | Ubuntu 24.04 LTS |

## 手动部署步骤

### 1. 安装 Docker

```bash
sudo apt-get update
sudo apt-get install -y ca-certificates curl gnupg

sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list

sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

sudo usermod -aG docker $USER
```

### 2. 配置环境

```bash
cd /opt/aiguardrails
cp env.template .env
vim .env  # 编辑配置
```

### 3. 启动服务

```bash
docker compose -f docker-compose.prod.yml up -d
```

### 4. 验证部署

```bash
# 检查服务状态
docker compose -f docker-compose.prod.yml ps

# 查看日志
docker compose -f docker-compose.prod.yml logs -f

# 健康检查
curl http://localhost/health
```

## 配置说明

主要配置项 (`.env` 文件):

| 配置 | 说明 | 必填 |
|------|------|------|
| `POSTGRES_PASSWORD` | 数据库密码 | ✅ |
| `REDIS_PASSWORD` | Redis密码 | ✅ |
| `ADMIN_TOKEN` | API管理令牌 | ✅ |
| `ADMIN_JWT_SECRET` | JWT签名密钥 | ✅ |
| `WECHAT_APP_ID` | 微信AppID | 可选 |
| `ALIPAY_APP_ID` | 支付宝AppID | 可选 |
| `SMS_PROVIDER` | 短信服务商 | 可选 |

## 常用命令

```bash
# 重启服务
docker compose -f docker-compose.prod.yml restart

# 停止服务
docker compose -f docker-compose.prod.yml down

# 更新部署
git pull
docker compose -f docker-compose.prod.yml build
docker compose -f docker-compose.prod.yml up -d

# 查看日志
docker compose -f docker-compose.prod.yml logs -f api

# 进入数据库
docker compose -f docker-compose.prod.yml exec db psql -U app -d aiguardrails

# 备份数据库
docker compose -f docker-compose.prod.yml exec db pg_dump -U app aiguardrails > backup.sql
```

## 端口说明

| 端口 | 服务 |
|------|------|
| 80 | HTTP (Nginx) |
| 443 | HTTPS (Nginx) |
| 8080 | API (内部) |

## 阿里云安全组配置

在阿里云控制台配置安全组规则：

| 方向 | 端口 | 协议 | 授权对象 |
|------|------|------|---------|
| 入 | 80 | TCP | 0.0.0.0/0 |
| 入 | 443 | TCP | 0.0.0.0/0 |
| 入 | 22 | TCP | 您的IP |

## 故障排查

```bash
# 查看容器状态
docker compose -f docker-compose.prod.yml ps

# 查看API日志
docker compose -f docker-compose.prod.yml logs api

# 检查数据库连接
docker compose -f docker-compose.prod.yml exec api nc -zv db 5432

# 检查Redis连接
docker compose -f docker-compose.prod.yml exec api nc -zv redis 6379
```

## 支持

- 文档: `/docs`
- API: `/v1`
- 健康检查: `/health`
