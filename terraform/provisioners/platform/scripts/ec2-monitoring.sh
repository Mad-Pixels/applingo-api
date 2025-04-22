#!/bin/bash
set -euxo pipefail

# Update the system
yum update -y

# Install Docker
amazon-linux-extras install docker -y
systemctl enable docker
systemctl start docker

# Add EC2-user to the docker group (note: should use ec2-user instead of ec2)
usermod -aG docker ec2-user

# Install necessary tools
yum install -y amazon-cloudwatch-agent jq awslogs

# Set up Docker to start on boot
systemctl enable docker.service

# Install Docker Compose
curl -L "https://github.com/docker/compose/releases/download/v2.18.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose

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