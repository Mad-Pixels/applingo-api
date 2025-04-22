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

# Install Docker Compose using pip with правильными версиями зависимостей
yum install -y python3-pip
pip3 install urllib3==1.26.16
pip3 install docker-compose

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

echo "EC2 instance setup completed successfully"