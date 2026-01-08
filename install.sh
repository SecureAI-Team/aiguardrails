#!/usr/bin/env bash
set -euo pipefail

#######################################################################
# AI GuardRails 一键部署脚本
# 适用于: Ubuntu 24.04 LTS (阿里云 ECS)
# 用法: curl -sSL https://raw.githubusercontent.com/.../install.sh | bash
#       或者: bash install.sh
#######################################################################

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# 配置变量
PROJECT_NAME="aiguardrails"
INSTALL_DIR="${INSTALL_DIR:-/opt/aiguardrails}"
REPO_URL="${REPO_URL:-}"
DOMAIN="${DOMAIN:-}"

# 检查是否为root
check_root() {
    if [[ $EUID -eq 0 ]]; then
        log_warn "建议使用非root用户运行，将自动使用sudo"
    fi
}

# 检查系统
check_system() {
    log_info "检查系统环境..."
    if [[ ! -f /etc/os-release ]]; then
        log_error "不支持的操作系统"
        exit 1
    fi
    source /etc/os-release
    if [[ "$ID" != "ubuntu" ]] || [[ "${VERSION_ID}" < "22.04" ]]; then
        log_warn "推荐使用 Ubuntu 22.04/24.04 LTS"
    fi
    log_info "系统: $PRETTY_NAME"
}

# 安装Docker
install_docker() {
    if command -v docker &> /dev/null; then
        log_info "Docker 已安装: $(docker --version)"
        return
    fi
    
    log_info "安装 Docker..."
    sudo apt-get update
    sudo apt-get install -y ca-certificates curl gnupg
    
    sudo install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    sudo chmod a+r /etc/apt/keyrings/docker.gpg
    
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
      $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
      sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    
    sudo apt-get update
    sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    
    # 添加当前用户到docker组
    sudo usermod -aG docker $USER
    log_info "Docker 安装完成"
}

# 准备项目目录
prepare_project() {
    log_info "准备项目目录..."
    sudo mkdir -p "$INSTALL_DIR"
    sudo chown -R $USER:$USER "$INSTALL_DIR"
    
    if [[ -n "$REPO_URL" ]]; then
        if [[ -d "$INSTALL_DIR/.git" ]]; then
            log_info "更新代码..."
            cd "$INSTALL_DIR" && git pull
        else
            log_info "克隆代码..."
            git clone "$REPO_URL" "$INSTALL_DIR"
        fi
    else
        log_warn "本地部署模式，请确保代码已复制到 $INSTALL_DIR"
    fi
    
    cd "$INSTALL_DIR"
}

# 生成配置
generate_config() {
    log_info "生成配置文件..."
    
    # 生成随机密码
    ADMIN_TOKEN=${ADMIN_TOKEN:-$(openssl rand -hex 24)}
    ADMIN_JWT_SECRET=${ADMIN_JWT_SECRET:-$(openssl rand -hex 32)}
    POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-$(openssl rand -hex 16)}
    REDIS_PASSWORD=${REDIS_PASSWORD:-$(openssl rand -hex 16)}
    ADMIN_BOOT_PASSWORD=${ADMIN_BOOT_PASSWORD:-$(openssl rand -base64 12)}
    
    cat > "$INSTALL_DIR/.env" <<EOF
# ================================================
# AI GuardRails 环境配置
# 生成时间: $(date '+%Y-%m-%d %H:%M:%S')
# ================================================

# 数据库
POSTGRES_DB=aiguardrails
POSTGRES_USER=app
POSTGRES_PASSWORD=$POSTGRES_PASSWORD

# Redis
REDIS_PASSWORD=$REDIS_PASSWORD
REDIS_NAMESPACE=aiguardrails

# 管理员
ADMIN_TOKEN=$ADMIN_TOKEN
ADMIN_JWT_SECRET=$ADMIN_JWT_SECRET
ADMIN_BOOT_USER=admin@example.com
ADMIN_BOOT_PASSWORD=$ADMIN_BOOT_PASSWORD

# 端口
HTTP_PORT=80
HTTPS_PORT=443

# 域名 (可选，用于SSL证书)
DOMAIN=${DOMAIN:-}

# CORS
ALLOWED_ORIGINS=*

# OPA策略引擎
OPA_ENABLED=true

# =========== 以下按需配置 ===========

# 通义千问内容审核 (可选)
# QWEN_API_TOKEN=your-qwen-token

# 微信登录 (可选)
# WECHAT_APP_ID=
# WECHAT_APP_SECRET=

# 支付宝登录 (可选)
# ALIPAY_APP_ID=
# ALIPAY_PRIVATE_KEY=
# ALIPAY_PUBLIC_KEY=

# 短信服务 (可选)
# SMS_PROVIDER=aliyun
# SMS_ACCESS_KEY=
# SMS_SECRET_KEY=
# SMS_SIGN_NAME=
# SMS_TEMPLATE_CODE=

# OIDC (可选，企业SSO)
# OIDC_ISSUER=
# OIDC_AUDIENCE=
# OIDC_JWKS_URL=
EOF

    chmod 600 "$INSTALL_DIR/.env"
    log_info "配置文件已生成: $INSTALL_DIR/.env"
    log_warn "请保存以下凭证:"
    echo "-------------------------------------------"
    echo "管理员邮箱: admin@example.com"
    echo "管理员密码: $ADMIN_BOOT_PASSWORD"
    echo "API Token: $ADMIN_TOKEN"
    echo "-------------------------------------------"
}

# 创建Nginx配置
create_nginx_config() {
    log_info "创建 Nginx 配置..."
    mkdir -p "$INSTALL_DIR/deploy/nginx/ssl"
    mkdir -p "$INSTALL_DIR/deploy/nginx/logs"
    
    cat > "$INSTALL_DIR/deploy/nginx/nginx.conf" <<'EOF'
worker_processes auto;
error_log /var/log/nginx/error.log warn;
pid /var/run/nginx.pid;

events {
    worker_connections 2048;
    use epoll;
    multi_accept on;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;
    
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';
    access_log /var/log/nginx/access.log main;
    
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;
    client_max_body_size 10M;
    
    # Gzip
    gzip on;
    gzip_vary on;
    gzip_proxied any;
    gzip_comp_level 6;
    gzip_types text/plain text/css text/xml application/json application/javascript application/xml;
    
    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    limit_req_zone $binary_remote_addr zone=auth:10m rate=5r/s;
    
    upstream api {
        server api:8080;
        keepalive 32;
    }
    
    upstream frontend {
        server frontend:80;
        keepalive 16;
    }
    
    server {
        listen 80;
        server_name _;
        
        # Health check
        location /health {
            access_log off;
            return 200 'OK';
            add_header Content-Type text/plain;
        }
        
        # API
        location /v1/ {
            limit_req zone=api burst=20 nodelay;
            proxy_pass http://api;
            proxy_http_version 1.1;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_set_header Connection "";
            proxy_connect_timeout 10s;
            proxy_read_timeout 60s;
        }
        
        # Auth endpoints (stricter rate limit)
        location /v1/auth/ {
            limit_req zone=auth burst=5 nodelay;
            proxy_pass http://api;
            proxy_http_version 1.1;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        }
        
        # Frontend
        location / {
            proxy_pass http://frontend;
            proxy_http_version 1.1;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
        }
        
        # Static assets cache
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2)$ {
            proxy_pass http://frontend;
            expires 30d;
            add_header Cache-Control "public, immutable";
        }
    }
}
EOF
    log_info "Nginx 配置完成"
}

# 运行数据库迁移
run_migrations() {
    log_info "运行数据库迁移..."
    cd "$INSTALL_DIR"
    
    # 等待数据库启动
    docker compose -f docker-compose.prod.yml up -d db
    sleep 10
    
    # 执行迁移
    for sql_file in backend/migrations/*.sql; do
        if [[ -f "$sql_file" ]]; then
            log_info "执行迁移: $(basename $sql_file)"
            docker compose -f docker-compose.prod.yml exec -T db psql -U app -d aiguardrails < "$sql_file" 2>/dev/null || true
        fi
    done
    log_info "迁移完成"
}

# 构建并启动服务
start_services() {
    log_info "构建并启动服务..."
    cd "$INSTALL_DIR"
    
    # 构建镜像
    docker compose -f docker-compose.prod.yml build --no-cache
    
    # 启动所有服务
    docker compose -f docker-compose.prod.yml up -d
    
    log_info "等待服务启动..."
    sleep 15
    
    # 检查服务状态
    docker compose -f docker-compose.prod.yml ps
}

# 打印完成信息
print_success() {
    PUBLIC_IP=$(curl -s --connect-timeout 5 http://100.100.100.200/latest/meta-data/eip 2>/dev/null || curl -s ifconfig.me 2>/dev/null || echo "YOUR_IP")
    
    echo ""
    echo "========================================================"
    echo -e "${GREEN}✅ AI GuardRails 部署完成!${NC}"
    echo "========================================================"
    echo ""
    echo "访问地址:"
    echo "  - 前端: http://$PUBLIC_IP"
    echo "  - API:  http://$PUBLIC_IP/v1"
    echo "  - 健康检查: http://$PUBLIC_IP/health"
    echo ""
    echo "管理命令:"
    echo "  - 查看日志: docker compose -f docker-compose.prod.yml logs -f"
    echo "  - 重启服务: docker compose -f docker-compose.prod.yml restart"
    echo "  - 停止服务: docker compose -f docker-compose.prod.yml down"
    echo ""
    echo "配置文件: $INSTALL_DIR/.env"
    echo "========================================================"
}

# 主函数
main() {
    echo ""
    echo "========================================================"
    echo "  AI GuardRails 一键部署 (Ubuntu 24.04 / 阿里云 ECS)"
    echo "========================================================"
    echo ""
    
    check_root
    check_system
    install_docker
    prepare_project
    generate_config
    create_nginx_config
    run_migrations
    start_services
    print_success
}

main "$@"
