#!/bin/bash
# /var/log/cloud-init-output.log
set -euxo pipefail

# -------- common setup --------
yum update -y
amazon-linux-extras install docker -y
systemctl enable docker
systemctl start docker
usermod -aG docker ec2-user
yum install -y amazon-cloudwatch-agent jq awslogs
systemctl enable docker.service
yum install -y python3-pip
pip3 install urllib3==1.26.16 docker-compose
ln -sf /usr/local/bin/docker-compose /usr/bin/docker-compose

echo "vm.swappiness=10"    >> /etc/sysctl.conf
echo "vm.max_map_count=262144" >> /etc/sysctl.conf
sysctl -p

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

# -------- monitoring dirs --------
mkdir -p /home/ec2-user/monitoring/{grafana,prometheus,nginx}
chown -R ec2-user:ec2-user /home/ec2-user/monitoring

# -------- prometheus.yml --------
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

# -------- nginx.conf --------
cat > /home/ec2-user/monitoring/nginx/nginx.conf << 'EOF'
events {}
http {
  server {
    listen 80 default_server;
    server_name _;

    location / {
      proxy_pass         http://grafana:3000/;
      proxy_set_header   Host              $host;
      proxy_set_header   X-Real-IP         $remote_addr;
      proxy_set_header   X-Forwarded-For   $proxy_add_x_forwarded_for;
      proxy_set_header   X-Forwarded-Proto $scheme;
    }
  }
}
EOF

# -------- docker-compose.yml --------
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
      - GF_SERVER_ROOT_URL=/
      - GF_SERVER_SERVE_FROM_SUB_PATH=false
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

# -------- start stack & systemd unit --------
cd /home/ec2-user/monitoring
/usr/bin/docker-compose up -d

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

systemctl enable monitoring.service

echo "EC2 instance setup completed successfully with monitoring stack"
