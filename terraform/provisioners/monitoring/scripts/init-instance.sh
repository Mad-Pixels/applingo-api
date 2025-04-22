#!/bin/bash
# /var/log/cloud-init-output.log
set -euxo pipefail

# Update the system
yum update -y

# Install Docker
amazon-linux-extras install docker -y
systemctl enable docker
systemctl start docker

# Add EC2-user to the docker group
usermod -aG docker ec2-user

# Install necessary tools
yum install -y amazon-cloudwatch-agent jq awslogs

# Set up Docker to start on boot
systemctl enable docker.service

# Install Docker Compose
yum install -y python3-pip
pip3 install urllib3==1.26.16
pip3 install docker-compose

ln -sf /usr/local/bin/docker-compose /usr/bin/docker-compose

# Configure memory and swap limits
echo "vm.swappiness=10" >> /etc/sysctl.conf
echo "vm.max_map_count=262144" >> /etc/sysctl.conf
sysctl -p

# Set up log rotation for Docker
cat > /etc/logrotate.d/docker << 'EOF'
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

# Create monitoring directory structure
mkdir -p /home/ec2-user/monitoring/grafana
mkdir -p /home/ec2-user/monitoring/prometheus
chown -R ec2-user:ec2-user /home/ec2-user/monitoring

# Create Prometheus configuration
cat > /home/ec2-user/monitoring/prometheus/prometheus.yml << 'EOF'
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

# Create docker-compose.yml configuration
cat > /home/ec2-user/monitoring/docker-compose.yml << 'EOF'
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
      - '--web.console.libraries=/etc/prometheus/console_libraries'
      - '--web.console.templates=/etc/prometheus/consoles'
      - '--web.enable-lifecycle'
    ports:
      - "9090:9090"
    networks:
      - monitoring

  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    restart: unless-stopped
    volumes:
      - grafana_data:/var/lib/grafana
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
      - GF_USERS_ALLOW_SIGN_UP=false
    ports:
      - "3000:3000"
    networks:
      - monitoring
    depends_on:
      - prometheus

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
    networks:
      - monitoring

networks:
  monitoring:
    driver: bridge

volumes:
  prometheus_data:
  grafana_data:
EOF

# Fix permissions
chown -R ec2-user:ec2-user /home/ec2-user/monitoring

# Wait for docker-compose to be available
sleep 5

# Start monitoring stack with absolute path
cd /home/ec2-user/monitoring
/usr/bin/docker-compose up -d

# Create systemd service to start monitoring on boot
cat > /etc/systemd/system/monitoring.service << 'EOF'
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

# Enable and start the monitoring service
systemctl enable monitoring.service

echo "EC2 instance setup completed successfully with monitoring stack"