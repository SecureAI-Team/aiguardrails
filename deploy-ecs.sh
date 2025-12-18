#!/usr/bin/env bash
set -euo pipefail

# One-click deploy on Ubuntu 24.04 host with Docker/Compose

PROJECT_DIR="/opt/aiguardrails"
REPO_URL="${REPO_URL:-}"          # optional: if empty, assume code already present
API_PORT="${API_PORT:-8080}"
WEB_PORT="${WEB_PORT:-8081}"
PROXY_PORT="${PROXY_PORT:-8443}"

# Defaults (auto-generated if not provided)
DATABASE_URL="${DATABASE_URL:-postgres://app:app@db:5432/aiguardrails?sslmode=disable}"
REDIS_URL="${REDIS_URL:-redis://redis:6379}"
ADMIN_SECRET_KEY="${ADMIN_SECRET_KEY:-ADMIN_TOKEN}"
QWEN_SECRET_KEY="${QWEN_SECRET_KEY:-QWEN_API_TOKEN}"
ALLOWED_ORIGINS="${ALLOWED_ORIGINS:-*}"

# Generate tokens if not provided via secret keys
if [ -z "${ADMIN_TOKEN:-}" ]; then
  ADMIN_TOKEN=$(openssl rand -hex 24)
fi
if [ -z "${QWEN_API_TOKEN:-}" ]; then
  QWEN_API_TOKEN="placeholder-qwen-token"
fi

echo "[1/6] Install Docker & Compose"
if ! command -v docker >/dev/null 2>&1; then
  sudo apt-get update
  sudo apt-get install -y ca-certificates curl gnupg
  sudo install -m 0755 -d /etc/apt/keyrings
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
  echo \
    "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
    $(. /etc/os-release && echo \"$VERSION_CODENAME\") stable" | \
    sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
  sudo apt-get update
  sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
fi

echo "[2/6] Prepare project directory"
sudo mkdir -p "$PROJECT_DIR"
sudo chown "$USER" "$PROJECT_DIR"

if [ -n "$REPO_URL" ]; then
  echo "[3/6] Clone/Update repo"
  if [ ! -d "$PROJECT_DIR/.git" ]; then
    git clone "$REPO_URL" "$PROJECT_DIR"
  else
    (cd "$PROJECT_DIR" && git pull)
  fi
else
  echo "[3/6] Using existing workspace at $PROJECT_DIR (copy your code there)"
fi

cd "$PROJECT_DIR"

echo "[4/6] Create .env file for docker-compose"
cat > .env <<EOF
DATABASE_URL=$DATABASE_URL
REDIS_URL=$REDIS_URL
ADMIN_SECRET_KEY=$ADMIN_SECRET_KEY
QWEN_SECRET_KEY=$QWEN_SECRET_KEY
ALLOWED_ORIGINS=$ALLOWED_ORIGINS
API_PORT=$API_PORT
WEB_PORT=$WEB_PORT
PROXY_PORT=$PROXY_PORT
ADMIN_TOKEN=$ADMIN_TOKEN
QWEN_API_TOKEN=$QWEN_API_TOKEN
# Optional overrides:
# OIDC_ISSUER=
# OIDC_AUDIENCE=
# OIDC_JWKS_URL=
# REDIS_NAMESPACE=aiguardrails
# OUTPUT_MODE=mark
# QWEN_API_BASE=https://dashscope.aliyuncs.com/api/v1/services/aigc/text-moderation
# QWEN_MODEL=qwen-moderation
# QWEN_TIMEOUT_SEC=8
# QWEN_RPS=5
# QWEN_RETRIES=2
# LLM_CACHE_TTL_MIN=10
# LLM_WORKERS=2
# LLM_QUEUE_SIZE=64
EOF

echo "[5/6] Pull/build containers"
docker compose pull || true
docker compose build

echo "[6/6] Generate self-signed TLS cert for proxy (if missing)"
mkdir -p proxy/certs
if [ ! -f proxy/certs/proxy.crt ] || [ ! -f proxy/certs/proxy.key ]; then
  openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout proxy/certs/proxy.key -out proxy/certs/proxy.crt \
    -subj "/C=CN/ST=NA/L=NA/O=AI-Sec/OU=Proxy/CN=proxy.local"
fi

echo "Start services"
docker compose up -d

echo "Done. API on :$API_PORT, frontend on :$WEB_PORT, proxy on :$PROXY_PORT"
echo "Use 'docker compose logs -f' to view logs."

