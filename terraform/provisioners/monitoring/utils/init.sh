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
# ------------------- INIT ---------------------- #
# =============================================== #

TOKEN=$(curl -X PUT "http://169.254.169.254/latest/api/token" \
  -H "X-aws-ec2-metadata-token-ttl-seconds: 21600" -s)

ENVIRONMENT=$(fetch_tag Environment)
NAME=$(fetch_tag Name)
USER="ec2-user"

# --- DEFINE LOCALE ---
log_block green "Configuring locale"
export LANGUAGE=en_US.UTF-8
export LC_ALL=en_US.UTF-8
export LANG=en_US.UTF-8

# --- INSTALL BASE PACKAGES ---
log_block green "Updating system and installing dependencies"
yum update -y > /dev/null
amazon-linux-extras install docker -y > /dev/null
yum install -y amazon-cloudwatch-agent jq awslogs python3-pip unzip sqlite3 > /dev/null
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

# =============================================== #
# ----------------- MONITORING ------------------ #
# =============================================== #

SCRIPT="monitoring_stack.sh"

log_block blue "MONITORING part"
rm -rf /home/${USER}/${SCRIPT}
if [ -n "$NAME" ] && [ -n "$ENVIRONMENT" ]; then
    S3_BUCKET="${NAME}-${ENVIRONMENT}"

    if aws s3 ls "s3://${S3_BUCKET}/scriptis/${SCRIPT}" > /dev/null 2>&1; then
        aws s3 cp "s3://${S3_BUCKET}/scriptis/${SCRIPT}" "/home/${USER}/${SCRIPT}"
        chmod +x /home/${USER}/${SCRIPT}

        log_block green "Invoke monitoring stack"
        /home/${USER}/${SCRIPT}
    else
        log_block red "Monitoring action not found in bucket"
    fi
fi
