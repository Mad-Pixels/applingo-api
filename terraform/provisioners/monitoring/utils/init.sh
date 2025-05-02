#!/bin/bash

set -euo pipefail

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

TOKEN=$(curl -X PUT "http://169.254.169.254/latest/api/token" \
  -H "X-aws-ec2-metadata-token-ttl-seconds: 21600" -s)

ENVIRONMENT=$(fetch_tag Environment)
NAME=$(fetch_tag Name)

rm -rf /home/ec2-user/init-instance.sh
if [ -n "$NAME" ] && [ -n "$ENVIRONMENT" ]; then
    S3_BUCKET="${NAME}-${ENVIRONMENT}"

    if aws s3 ls "s3://${S3_BUCKET}/scriptis/init-instance.sh" > /dev/null 2>&1; then
        aws s3 cp "s3://${S3_BUCKET}/scriptis/init-instance.sh" "/home/ec2-user/init-instance.sh"
        chmod +x /home/ec2-user/init-instance.sh

        /home/ec2-user/init-instance.sh
    else
        echo "main action not found in bucket"
    fi
fi
