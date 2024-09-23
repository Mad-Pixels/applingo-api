# Curl

```bash
# localstack
url="http://localhost:4566/restapis/4663mz3v89/prod/_user_request_"

secret_token="your_secret_token"
path="/api/v1/dictionary/file_presign"
timestamp=$(date -u +%s)
signature=$(echo -n "${timestamp}${path}" | openssl dgst -sha256 -hmac "${secret_token}" | sed 's/^.* //')

curl -X POST ${url}${path} \
  -d '{"content_type": "text/csv", "name": "file.csv"}' \
  -H "Content-Type: application/json" \
  -H "X-Timestamp: ${timestamp}" \
  -H "X-Signature: ${signature}"
```

# Request

```json
{
  "content_type": "text/csv",
  "name": "filename"
}
```

# Response

```json
{
  "data":{
    "url":"http://172.17.0.2:4566/lingocards-processing/file.csv?X-Amz-Algorithm=AWS4-HMAC-SHA256..."
  }
}
```