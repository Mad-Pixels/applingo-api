#!/bin/bash
# /var/log/cloud-init-output.log
# /var/log/prometheus-backup.log

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
yum install -y amazon-cloudwatch-agent jq awslogs python3-pip unzip > /dev/null
pip3 install --quiet urllib3==1.26.16 docker-compose
ln -sf /usr/local/bin/docker-compose /usr/bin/docker-compose

# --- FETCH INSTANCE METADATA (IMDSv2 COMPATIBLE) ---
log_block blue "Fetching EC2 instance metadata..."

fetch_metadata() {
  local path="$1"
  local token

  token=$(curl -X PUT "http://169.254.169.254/latest/api/token" \
    -H "X-aws-ec2-metadata-token-ttl-seconds: 21600" -s)

  curl -H "X-aws-ec2-metadata-token: $token" -s "http://169.254.169.254/latest/$path"
}

TOKEN=$(curl -X PUT "http://169.254.169.254/latest/api/token" \
  -H "X-aws-ec2-metadata-token-ttl-seconds: 21600" -s)

INSTANCE_ID=$(curl -H "X-aws-ec2-metadata-token: $TOKEN" -s "http://169.254.169.254/latest/meta-data/instance-id")
REGION=$(curl -H "X-aws-ec2-metadata-token: $TOKEN" -s "http://169.254.169.254/latest/dynamic/instance-identity/document" | jq -r .region)

# --- FETCH EC2 INSTANCE TAGS ---
log_block blue "Fetching EC2 instance tags..."

fetch_tag() {
  local tag_key="$1"
  curl -H "X-aws-ec2-metadata-token: $TOKEN" -s "http://169.254.169.254/latest/meta-data/tags/instance/$tag_key" || true
}

ENVIRONMENT=$(fetch_tag Environment)
NAME=$(fetch_tag Name)

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

# --- SWAP SETUP ---
log_block green "Creating swap file"

if ! grep -q '/swapfile' /etc/fstab; then
  if ! swapon --show | grep -q '/swapfile'; then
    rm -f /swapfile
    dd if=/dev/zero of=/swapfile bs=1M count=1024 status=progress
    chmod 600 /swapfile
    mkswap /swapfile > /dev/null
    swapon /swapfile
    echo '/swapfile swap swap defaults 0 0' >> /etc/fstab
    log_block green "Swap file created and activated"
  else
    log_block blue "Swap file already active, skipping creation"
  fi
else
  log_block blue "Swap file already configured in fstab"
fi

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
mkdir -p /home/ec2-user/monitoring/{grafana,prometheus,nginx,cloudwatch,provisioning}
mkdir -p /home/ec2-user/monitoring/provisioning/datasources
mkdir -p /home/ec2-user/monitoring/data/prometheus
chown -R ec2-user:ec2-user /home/ec2-user/monitoring

mkdir -p /home/ec2-user/.aws
chown -R ec2-user:ec2-user /home/ec2-user/.aws

cd /home/ec2-user/monitoring

# --- SHUTDOWN OLD STACK ---
log_block blue "Stopping previous monitoring stack (if any)"
if [ -f docker-compose.yml ]; then
  docker-compose down --remove-orphans > /dev/null || true
fi

# --- WRITE AWS CONFIG ---
log_block green "Writing AWS config"
cat > /home/ec2-user/.aws/config <<EOF
[default]
region = ${REGION}
sts_regional_endpoints = regional
EOF

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
  - job_name: 'cloudwatch'
    static_configs:
      - targets: ['localhost:9106']
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

# --- WRITE CLOUDWATCH EXPORTER CONFIG ---
log_block green "Writing CloudWatch Exporter config"
cat > /home/ec2-user/monitoring/cloudwatch/cloudwatch-exporter.yml <<EOF
apiVersion: v1alpha1
discovery:
  exportedTagsOnMetrics:
    AWS/Lambda: ["FunctionName"]
  jobs:
  - type: AWS/Lambda
    regions:
      - ${REGION}
    metrics:
      - name: Invocations
        statistics:
          - Sum
        period: 300
        length: 600
EOF

# --- WRITE GRAFANA PROVISIONING CONFIG ---
log_block green "Writing Grafana provisioning config"
cat > /home/ec2-user/monitoring/provisioning/datasources/prometheus.yml <<EOF
apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true
    editable: true
EOF

# --- WRITE DOCKER COMPOSE STACK ---
log_block green "Writing docker-compose.yml"
cat > docker-compose.yml <<EOF
version: '3'
services:
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
      - '--storage.tsdb.retention.time=60d'
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
      - ./provisioning:/etc/grafana/provisioning:ro
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
  
  cloudwatch-exporter:
    image: prometheuscommunity/yet-another-cloudwatch-exporter:latest
    container_name: cloudwatch-exporter
    restart: unless-stopped
    volumes:
      - ./cloudwatch/cloudwatch-exporter.yml:/tmp/config.yml
      - ~/.aws/config:/root/.aws/config:ro
    environment:
      - AWS_STS_REGIONAL_ENDPOINTS=regional
      - AWS_SDK_LOAD_CONFIG=true
      - AWS_REGION=${REGION}
    ports:
      - "9106:9106"
    networks: [monitoring]
    command: ["--config.file=/tmp/config.yml", "--listen-address=:9106"]

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

# --- CHECK AND RESTORE PROMETHEUS DATA ---
log_block blue "Checking Prometheus data availability..."
if [ -n "${NAME:-}" ] && [ -n "${ENVIRONMENT:-}" ]; then
  ENDPOINT="applingo-monitoring-s3-endpoint"
  BUCKET_NAME="${NAME}-${ENVIRONMENT}"

  if [ ! -d "/home/ec2-user/monitoring/data/prometheus" ] || [ -z "$(ls -A /home/ec2-user/monitoring/data/prometheus 2>/dev/null)" ]; then
    log_block blue "Prometheus data directory is empty, attempting to restore from S3..."

    if aws --endpoint-url ${ENDPOINT} s3 ls "s3://${BUCKET_NAME}/prometheus-backup.tar.gz" > /dev/null 2>&1; then
      mkdir -p /home/ec2-user/monitoring/data
      aws --endpoint-url ${ENDPOINT} s3 cp "s3://${BUCKET_NAME}/prometheus-backup.tar.gz" /tmp/prometheus-backup.tar.gz
      tar -xzf /tmp/prometheus-backup.tar.gz -C /home/ec2-user/monitoring/data
      rm /tmp/prometheus-backup.tar.gz
      log_block green "Prometheus data restored from S3 backup."
    else
      log_block blue "No backup found in S3, continuing with empty data."
    fi
  else 
    log_block green "Prometheus data directory is not empty, skipping restore."
  fi
else
  log_block blue "Skipping Prometheus restore because NAME or ENVIRONMENT is empty."
fi

# --- BACKUP SCRIPT ---
log_block green "Creating Prometheus backup script..."
cat > /home/ec2-user/monitoring/backup.sh <<'EOF'
#!/bin/bash
set -euo pipefail

ENDPOINT_URL="http://s3.us-east-2.amazonaws.com"
BACKUP_FILE="/tmp/prometheus-backup.tar.gz"
DATA_DIR="/home/ec2-user/monitoring/data/prometheus"

# --- Fetch metadata token ---
TOKEN=$(curl -sX PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600")

# --- Get tags directly from IMDS ---
NAME=$(curl -s -H "X-aws-ec2-metadata-token: $TOKEN" http://169.254.169.254/latest/meta-data/tags/instance/Name || echo "")
ENVIRONMENT=$(curl -s -H "X-aws-ec2-metadata-token: $TOKEN" http://169.254.169.254/latest/meta-data/tags/instance/Environment || echo "")

if [ -n "$NAME" ] && [ -n "$ENVIRONMENT" ]; then
  S3_BUCKET="${NAME}-${ENVIRONMENT}"

  # Check if bucket exists
  if aws --endpoint-url "$ENDPOINT_URL" s3 ls "s3://${S3_BUCKET}" > /dev/null 2>&1; then
    echo ">>> Creating backup..."
    tar czf "$BACKUP_FILE" -C "$DATA_DIR" .

    echo ">>> Removing previous backup if exists..."
    aws --endpoint-url "$ENDPOINT_URL" s3 rm "s3://${S3_BUCKET}/prometheus-backup.tar.gz" || true

    echo ">>> Uploading new backup..."
    aws --endpoint-url "$ENDPOINT_URL" s3 cp "$BACKUP_FILE" "s3://${S3_BUCKET}/prometheus-backup.tar.gz" --storage-class STANDARD_IA

    echo ">>> Cleaning up local backup file..."
    rm -f "$BACKUP_FILE"
  else
    echo ">>> [INFO] S3 bucket ${S3_BUCKET} does not exist, skipping backup."
  fi
else
  echo ">>> [INFO] Missing NAME or ENVIRONMENT, skipping backup."
fi
EOF

chmod +x /home/ec2-user/monitoring/backup.sh
chown ec2-user:ec2-user /home/ec2-user/monitoring/backup.sh

# --- CRON JOB SETUP ---
log_block green "Setting up daily cron job for Prometheus backup..."
touch /var/log/prometheus-backup.log
chown ec2-user:ec2-user /var/log/prometheus-backup.log
chmod 644 /var/log/prometheus-backup.log

cat > /etc/cron.d/prometheus-backup <<'EOF'
SHELL=/bin/bash
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
0 3 * * * ec2-user /home/ec2-user/monitoring/backup.sh >> /var/log/prometheus-backup.log 2>&1
EOF
chmod 644 /etc/cron.d/prometheus-backup
systemctl restart crond
log_block green "Backup system configured successfully."

log_block green "EC2 monitoring stack setup completed successfully."