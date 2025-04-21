#!/bin/bash
set -euxo pipefail

um update -y
amazon-linux-extras install docker -y
systemctl enable docker
systemctl start docker
usermod -aG docker ec2

DOCKER_COMPOSE_VERSION="2.24.0"
curl -SL https://github.com/docker/compose/releases/download/v${DOCKER_COMPOSE_VERSION}/docker-compose-linux-aarch64 -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose

mkdir -p /opt/grafana
cat > /opt/grafana/docker-compose.yml <<EOF
version: "3"
services:
  grafana:
    image: grafana/grafana-oss:latest
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana-storage:/var/lib/grafana
volumes:
  grafana-storage:
EOF

cd /opt/grafana
docker-compose up -d