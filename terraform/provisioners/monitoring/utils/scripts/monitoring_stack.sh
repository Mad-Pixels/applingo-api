#!/bin/bash

set -euo pipefail

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

fetch_metadata() {
  local path="$1"
  local token

  token=$(curl -X PUT "http://169.254.169.254/latest/api/token" \
    -H "X-aws-ec2-metadata-token-ttl-seconds: 21600" -s)

  curl -H "X-aws-ec2-metadata-token: $token" -s "http://169.254.169.254/latest/$path"
}

fetch_tag() {
  local tag_key="$1"
  curl -H "X-aws-ec2-metadata-token: $TOKEN" -s "http://169.254.169.254/latest/meta-data/tags/instance/$tag_key" || true
}



# =============================================== #
# ---------------- VARIABLES -------------------- #
# =============================================== #

TOKEN=$(curl -X PUT "http://169.254.169.254/latest/api/token" \
  -H "X-aws-ec2-metadata-token-ttl-seconds: 21600" -s)

REGION=$(curl -H "X-aws-ec2-metadata-token: $TOKEN" -s "http://169.254.169.254/latest/dynamic/instance-identity/document" | jq -r .region)
INSTANCE=$(curl -H "X-aws-ec2-metadata-token: $TOKEN" -s "http://169.254.169.254/latest/meta-data/instance-id")
ENVIRONMENT=$(fetch_tag Environment)
NAME=$(fetch_tag Name)


USER="ec2-user"
MONITORING_DIR="/home/${USER}/monitoring"

PROMETHEUS_DIR="${MONITORING_DIR}/prometheus"
EXPORTERS_DIR="${MONITORING_DIR}/exporters"
GRAFANA_DIR="${MONITORING_DIR}/grafana"
NGINX_DIR="${MONITORING_DIR}/nginx"

GRAFANA_PROVISION_DATASOURCE_DIR=${GRAFANA_DIR}/provisioning/datasources
GRAFANA_PROVISION_DASHBOARDS_DIR=${GRAFANA_DIR}/provisioning/dashboards
EXPORTER_CLOUDWATCH_DIR=${EXPORTERS_DIR}/cloudwatch



# =============================================== #
# ----------------- SHUTDOWN -------------------- #
# =============================================== #

# --- SHUTDOWN OLD STACK ---
log_block blue "Stopping previous monitoring stack (if any)"
cd ${MONITORING_DIR}
if [ -f docker-compose.yml ]; then
  docker-compose down --remove-orphans > /dev/null || true
fi

# --- BACKUP PROMETHEUS DATA ---
log_block green "Backup prometheus data"
if [ -n "$NAME" ] && [ -n "$ENVIRONMENT" ] && [ -d "${PROMETHEUS_DIR}/data" ]; then
  S3_BUCKET="${NAME}-${ENVIRONMENT}"

  if aws s3 ls "s3://${S3_BUCKET}" > /dev/null 2>&1; then
    echo ">>> Creating backup..."

    BACKUP_FILE="prometheus-backup.tar.gz"
    tar czf "/tmp/${BACKUP_FILE}" -C "${PROMETHEUS_DIR}/data" .

    echo ">>> Removing previous backup if exists..."
    aws s3 rm "s3://${S3_BUCKET}/backups/${BACKUP_FILE}" || true

    echo ">>> Uploading new backup..."
    aws s3 cp "/tmp/${BACKUP_FILE}" "s3://${S3_BUCKET}/backups/${BACKUP_FILE}" --storage-class STANDARD_IA

    echo ">>> Cleaning up local backup file..."
    rm -f "/tmp/${BACKUP_FILE}"
  else
    echo ">>> [INFO] S3 bucket ${S3_BUCKET} does not exist, skipping backup."
  fi 
else 
  echo ">>> [INFO] Missing NAME or ENVIRONMENT, skipping backup."
fi

# --- BACKUP GRAFANA DATA ---
log_block green "Backup grafana data"
if [ -n "$NAME" ] && [ -n "$ENVIRONMENT" ] && [ -f "${GRAFANA_DIR}/data/grafana.db" ]; then
  S3_BUCKET="${NAME}-${ENVIRONMENT}"

  if aws s3 ls "s3://${S3_BUCKET}" > /dev/null 2>&1; then
    echo ">>> Deleting current dashboards"
    sqlite3 ${GRAFANA_DIR}/data/grafana.db "DELETE FROM dashboard;"
    sqlite3 ${GRAFANA_DIR}/data/grafana.db "DELETE FROM dashboard_provisioning;"
    sqlite3 ${GRAFANA_DIR}/data/grafana.db "DELETE FROM dashboard_version;"
    sqlite3 ${GRAFANA_DIR}/data/grafana.db "DELETE FROM dashboard_public;"
    sqlite3 ${GRAFANA_DIR}/data/grafana.db "DELETE FROM dashboard_acl;"
    sqlite3 ${GRAFANA_DIR}/data/grafana.db "DELETE FROM dashboard_snapshot;"

    echo ">>> Removing previous backup if exists..."
    aws s3 rm "s3://${S3_BUCKET}/backups/grafana.db" || true

    echo ">>> Uploading new backup..."
    aws s3 cp "${GRAFANA_DIR}/data/grafana.db" "s3://${S3_BUCKET}/backups/grafana.db" --storage-class STANDARD_IA 
  else
    echo ">>> [INFO] S3 bucket ${S3_BUCKET} does not exist, skipping backup."
  fi
else 
  echo ">>> [INFO] Missing NAME or ENVIRONMENT, skipping backup."
fi

# --- CLEANUP INSTANCE ---
rm -rf ${MONITORING_DIR}
rm -rf /home/${USER}/.aws



# =============================================== #
# ----------------- PREPARE --------------------- #
# =============================================== #

# --- DIRECTORIES ---
log_block green "Preparing monitoring folders"

mkdir -p ${GRAFANA_PROVISION_DATASOURCE_DIR}
mkdir -p ${GRAFANA_PROVISION_DASHBOARDS_DIR}
mkdir -p ${EXPORTER_CLOUDWATCH_DIR}

mkdir -p ${PROMETHEUS_DIR}
mkdir -p ${EXPORTERS_DIR}
mkdir -p ${GRAFANA_DIR}
mkdir -p ${NGINX_DIR}
chown -R ${USER}:${USER} ${MONITORING_DIR}

mkdir -p /home/${USER}/.aws
chown -R ${USER}:${USER} /home/${USER}/.aws

mkdir -p ${PROMETHEUS_DIR}/data
chown -R 65534:65534 ${PROMETHEUS_DIR}/data

mkdir -p ${GRAFANA_DIR}/data/
chown -R 472:472 ${GRAFANA_DIR}/data

cd ${MONITORING_DIR}

# --- WRITE AWS CONFIG ---
log_block green "Writing AWS config"
cat > /home/${USER}/.aws/config <<EOF
[default]
region = ${REGION}
sts_regional_endpoints = regional
EOF

# --- WRITE PROMETHEUS CONFIG ---
log_block green "Writing Prometheus config"
cat > ${PROMETHEUS_DIR}/prometheus.yml <<'EOF'
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
      - targets: ['cloudwatch-exporter:9106']
EOF

# --- CHECK AND RESTORE PROMETHEUS DATA ---
log_block blue "Checking Prometheus data availability..."
if [ -n "${NAME:-}" ] && [ -n "${ENVIRONMENT:-}" ]; then
  S3_BUCKET="${NAME}-${ENVIRONMENT}"

  if [ ! -d "${PROMETHEUS_DIR}/data" ] || [ -z "$(ls -A ${PROMETHEUS_DIR}/data 2>/dev/null)" ]; then
    log_block blue "Prometheus data directory is empty, attempting to restore from S3..."
    BACKUP_FILE="prometheus-backup.tar.gz"

    if aws s3 ls "s3://${S3_BUCKET}/backups/${BACKUP_FILE}" > /dev/null 2>&1; then
      aws s3 cp "s3://${S3_BUCKET}/backups/${BACKUP_FILE}" /tmp/${BACKUP_FILE}
      tar -xzf /tmp/${BACKUP_FILE} -C ${PROMETHEUS_DIR}/data
      rm /tmp/${BACKUP_FILE}
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

# --- WRITE NGINX CONFIG ---
log_block green "Writing Nginx config"
cat > ${NGINX_DIR}/nginx.conf <<'EOF'
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
      proxy_set_header   Host $host;
      proxy_set_header   X-Real-IP $remote_addr;
      proxy_set_header   X-Forwarded-For $proxy_add_x_forwarded_for;
      proxy_set_header   X-Forwarded-Proto $scheme;
    }
  }
}
EOF

# --- WRITE CLOUDWATCH EXPORTER CONFIG ---
log_block green "Writing CloudWatch Exporter config"
cat > ${EXPORTER_CLOUDWATCH_DIR}/exporter.yml <<EOF
apiVersion: v1alpha1
sts-region: ${REGION}
discovery:
  exportedTagsOnMetrics:
    AWS/Lambda: 
      - "FunctionName"
      - "Environment"
      - "Project"
      - "Target"
  jobs:
  - type: AWS/Lambda
    regions:
      - ${REGION}
    searchTags:
      - key: Project
        value: applingo
    metrics:
      - name: Invocations
        statistics:
          - Sum
        period: 300
        length: 300
      
      - name: Errors
        statistics:
          - Sum
        period: 300
        length: 300
      
      - name: Duration
        statistics:
          - Average
          - Maximum
          - Minimum
          - p50
          - p90
          - p95
          - p99
        period: 300
        length: 300
      
      - name: Throttles
        statistics:
          - Sum
        period: 300
        length: 300
      
      - name: ConcurrentExecutions
        statistics:
          - Maximum
        period: 300
        length: 300
EOF

# --- WRITE GRAFANA PROVISIONING CONFIG ---
log_block green "Writing Grafana provisioning config"
cat > ${GRAFANA_PROVISION_DATASOURCE_DIR}/prometheus.yml <<EOF
apiVersion: 1

datasources:
  - name: Prometheus
    uid: prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true
    editable: true
EOF

# --- WRITE GRAFANA DASHBOARDS CONFIG ---
log_block green "Writing Grafana dashboards provisioning config"
cat > ${GRAFANA_PROVISION_DASHBOARDS_DIR}/dashboards.yml <<EOF
apiVersion: 1

providers:
  - name: 'default'
    orgId: 1
    folder: ''
    type: file
    disableDeletion: false
    editable: true
    options:
      path: /var/lib/grafana/dashboards
    updateIntervalSeconds: 30
EOF

# --- CHECK AND RESTORE GRAFANA DASHBOARDS ---
if [ -n "${NAME:-}" ] && [ -n "${ENVIRONMENT:-}" ]; then
  S3_BUCKET="${NAME}-${ENVIRONMENT}"

  log_block blue "Attempting to download Grafana dashboards from S3..."
  if aws s3 ls "s3://${S3_BUCKET}/backups/dashboards/" > /dev/null 2>&1; then
    aws s3 sync "s3://${S3_BUCKET}/backups/dashboards/" ${GRAFANA_DIR}/data/dashboards --delete
    log_block green "Grafana dashboards restored from S3."
  else 
    log_block blue "No dashboards found in S3, skipping Grafana restore."
  fi 
else 
  log_block blue "Skipping Grafana restore because NAME or ENVIRONMENT is empty."
fi

# --- CHECK AND RESTORE GRAFANA DATA ---
if [ -n "${NAME:-}" ] && [ -n "${ENVIRONMENT:-}" ]; then
  S3_BUCKET="${NAME}-${ENVIRONMENT}"

  log_block blue "Attempting to download Grafana data from S3..."
  if aws s3 ls  "s3://${S3_BUCKET}/backups/grafana.db" > /dev/null 2>&1; then 
    aws s3 cp "s3://${S3_BUCKET}/backups/grafana.db" "${GRAFANA_DIR}/data/grafana.db"
    chown 472:472 ${GRAFANA_DIR}/data/grafana.db
    chmod 660 ${GRAFANA_DIR}/data/grafana.db
    log_block green "Grafana data restored from S3."
  else 
    log_block blue "No data found in S3, skipping Grafana restore."
  fi
else 
  log_block blue "Skipping Grafana restore because NAME or ENVIRONMENT is empty."
fi

chown -R 472:472 ${GRAFANA_DIR}/data

# --- WRITE DOCKER COMPOSE STACK ---
log_block green "Writing docker-compose.yml"
cat > ${MONITORING_DIR}/docker-compose.yml <<EOF
version: '3'
services:
  nginx:
    image: nginx:stable-alpine
    container_name: nginx
    restart: unless-stopped
    ports:
      - "80:80"
    volumes:
      - ${NGINX_DIR}/nginx.conf:/etc/nginx/nginx.conf:ro
    depends_on:
      - grafana
    networks:
      - monitoring

  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    restart: unless-stopped
    volumes:
      - ${PROMETHEUS_DIR}:/etc/prometheus
      - ${PROMETHEUS_DIR}/data:/prometheus
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
      - ${GRAFANA_DIR}/provisioning:/etc/grafana/provisioning:ro
      - ${GRAFANA_DIR}/data:/var/lib/grafana
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
      - ${EXPORTER_CLOUDWATCH_DIR}/exporter.yml:/tmp/config.yml
      - /home/${USER}/.aws/config:/root/.aws/config:ro
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
EOF



# =============================================== #
# ------------------ START ---------------------- #
# =============================================== #

log_block blue "Starting monitoring stack"
docker-compose up -d > /dev/null

# --- SYSTEMD UNIT ---
log_block green "Creating systemd unit"
cat > /etc/systemd/system/monitoring.service <<EOF
[Unit]
Description=Docker Compose Monitoring Stack
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=${MONITORING_DIR}
ExecStart=/usr/bin/docker-compose up -d
ExecStop=/usr/bin/docker-compose down

[Install]
WantedBy=multi-user.target
EOF

systemctl enable monitoring.service > /dev/null
systemctl daemon-reload



# =============================================== #
# ------------------ COMMON --------------------- #
# =============================================== #

# --- BACKUP SCRIPT ---
log_block green "Creating Prometheus backup script..."
cat > /home/ec2-user/monitoring/backup.sh <<'EOF'
#!/bin/bash
set -euo pipefail

BACKUP_FILE="prometheus-backup.tar.gz"
DATA_DIR="/home/ec2-user/monitoring/prometheus/data"

TOKEN=$(curl -sX PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600")
NAME=$(curl -s -H "X-aws-ec2-metadata-token: $TOKEN" http://169.254.169.254/latest/meta-data/tags/instance/Name || echo "")
ENVIRONMENT=$(curl -s -H "X-aws-ec2-metadata-token: $TOKEN" http://169.254.169.254/latest/meta-data/tags/instance/Environment || echo "")

if [ -n "$NAME" ] && [ -n "$ENVIRONMENT" ]; then
  S3_BUCKET="${NAME}-${ENVIRONMENT}"

  if aws s3 ls "s3://${S3_BUCKET}" > /dev/null 2>&1; then
    echo ">>> Creating backup..."
    tar czf "/tmp/${BACKUP_FILE}" -C "${DATA_DIR}" .

    echo ">>> Removing previous backup if exists..."
    aws s3 rm "s3://${S3_BUCKET}/backups/${BACKUP_FILE}" || true

    echo ">>> Uploading new backup..."
    aws s3 cp "$BACKUP_FILE" "s3://${S3_BUCKET}/backups/${BACKUP_FILE}" --storage-class STANDARD_IA

    echo ">>> Cleaning up local backup file..."
    rm -f "$BACKUP_FILE"
  else
    echo ">>> [INFO] S3 bucket ${S3_BUCKET} does not exist, skipping backup."
  fi
else 
  echo ">>> [INFO] Missing NAME or ENVIRONMENT, skipping backup."
fi
EOF

chmod +x ${MONITORING_DIR}/backup.sh
chown ${USER}:${USER} ${MONITORING_DIR}/backup.sh



# =============================================== #
# ------------------- CRON ---------------------- #
# =============================================== #

# --- BACKUP JOB SETUP ---
log_block green "Setting up daily cron job for Prometheus backup..."
touch /var/log/prometheus-backup.log

chown ${USER}:${USER} /var/log/prometheus-backup.log
chmod 644 /var/log/prometheus-backup.log

cat > /etc/cron.d/prometheus-backup <<EOF
SHELL=/bin/bash
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
0 3 * * * ${USER} ${MONITORING_DIR}/backup.sh >> /var/log/prometheus-backup.log 2>&1
EOF

chmod 644 /etc/cron.d/prometheus-backup
systemctl restart crond
log_block green "Backup system configured successfully."

log_block green "EC2 monitoring stack setup completed successfully."