#!/bin/bash
# Cloud-init EC2 monitoring stack installer
set -euo pipefail

# --- LOG UI BLOCKS ---
log_block() {
  COLOR="$1"
  MSG="$2"
  case "$COLOR" in
    green)  echo -e "\033[1;32m>>> [ OK ] $MSG\033[0m" ;;
    blue)   echo -e "\033[1;34m>>> [INFO] $MSG\033[0m" ;;
    red)    echo -e "\033[1;31m>>> [FAIL] $MSG\033[0m" ;;
    *)      echo -e ">>> $MSG" ;;
  esac
}

# --- FIX LOCALE ---
log_block blue "Configuring locale"
export LC_ALL=en_US.UTF-8
export LANG=en_US.UTF-8
export LANGUAGE=en_US.UTF-8

# --- INSTALL BASE PACKAGES ---
log_block green "Updating system and installing dependencies"
yum update -y > /dev/null
amazon-linux-extras install docker -y > /dev/null
yum install -y amazon-cloudwatch-agent jq awslogs python3-pip > /dev/null
pip3 install --quiet urllib3==1.26.16 docker-compose
ln -sf /usr/local/bin/docker-compose /usr/bin/docker-compose

# --- ENABLE DOCKER ---
log_block green "Starting Docker service"
systemctl enable docker > /dev/null
systemctl start docker
usermod -aG docker ec2-user

# --- SYSCTL TUNING ---
log_block green "Applying sysctl limits"
{
  echo "vm.swappiness=10"
  echo "vm.max_map_count=262144"
} >> /etc/sysctl.conf
sysctl -p > /dev/null

# --- LOGROTATE FOR DOCKER ---
log_block green "Configuring Docker log rotation"
cat > /etc/logrotate.d/docker <<'EOF'
/var/lib/docker/containers/*/*.log {
    rotate 7
    daily
    compress
    size=10M
    missingok
    delaycompress
    copytruncate
}
EOF

# --- CREATE DIRECTORIES ---
log_block green "Preparing monitoring folders"
mkdir -p /home/ec2-user/monitoring/{grafana,prometheus,nginx}
chown -R ec2-user:ec2-user /home/ec2-user/monitoring

cd /home/ec2-user/monitoring

# --- SHUTDOWN OLD STACK ---
log_block blue "Stopping previous monitoring stack (if any)"
if [ -f docker-compose.yml ]; then
  docker-compose down --remove-orphans > /dev/null || true
fi

# --- WRITE PROMETHEUS CONFIG ---
log_block green "Writing Prometheus config"
cat > prometheus/prometheus.yml <<'EOF'
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
  - job_name: 'node-exporter'
    static_configs:
      - targets: ['node-exporter:9100']
EOF

# --- WRITE NGINX CONFIG ---
log_block green "Writing Nginx config"
cat > nginx/nginx.conf <<'EOF'
events {}
http {
  server {
    listen 80 default_server;
    server_name _;

    location = /grafana {
      return 301 /grafana/;
    }

    location / {
      proxy_pass         http://grafana:3000/;
      proxy_http_version 1.1;
      proxy_set_header Host $host;
      proxy_set_header X-Real-IP $remote_addr;
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
      proxy_set_header X-Forwarded-Proto $scheme;
    }
  }
}
EOF

# --- WRITE DOCKER COMPOSE STACK ---
log_block green "Writing docker-compose.yml"
cat > docker-compose.yml <<'EOF'
version: '3'
services:
  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    restart: unless-stopped
    volumes:
      - ./prometheus:/etc/prometheus
      - prometheus_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.enable-lifecycle'
    ports:
      - "9090:9090"
    networks: [monitoring]

  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    restart: unless-stopped
    volumes:
      - grafana_data:/var/lib/grafana
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
      - GF_USERS_ALLOW_SIGN_UP=false
      - GF_SERVER_ROOT_URL=/grafana/
      - GF_SERVER_SERVE_FROM_SUB_PATH=true
    ports:
      - "3000:3000"
    networks: [monitoring]
    depends_on: [prometheus]

  node-exporter:
    image: prom/node-exporter:latest
    container_name: node-exporter
    restart: unless-stopped
    volumes:
      - /proc:/host/proc:ro
      - /sys:/host/sys:ro
      - /:/rootfs:ro
    command:
      - '--path.procfs=/host/proc'
      - '--path.rootfs=/rootfs'
      - '--path.sysfs=/host/sys'
      - '--collector.filesystem.mount-points-exclude=^/(sys|proc|dev|host|etc)($$|/)'
    ports:
      - "9100:9100"
    networks: [monitoring]

  nginx:
    image: nginx:stable-alpine
    container_name: nginx
    restart: unless-stopped
    ports:
      - "80:80"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
    depends_on:
      - grafana
    networks:
      - monitoring

networks:
  monitoring:
    driver: bridge

volumes:
  prometheus_data:
  grafana_data:
EOF

# --- START STACK ---
log_block blue "Starting monitoring stack"
docker-compose up -d > /dev/null

# --- SYSTEMD UNIT ---
log_block green "Creating systemd unit"
cat > /etc/systemd/system/monitoring.service <<'EOF'
[Unit]
Description=Docker Compose Monitoring Stack
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/home/ec2-user/monitoring
ExecStart=/usr/bin/docker-compose up -d
ExecStop=/usr/bin/docker-compose down

[Install]
WantedBy=multi-user.target
EOF

systemctl enable monitoring.service > /dev/null

log_block green "EC2 monitoring stack setup completed successfully."