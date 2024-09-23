# Curl

```bash
# localstack
url="http://localhost:4566/restapis/4663mz3v89/prod/_user_request_"

secret_token="your_secret_token"
path="/api/v1/dictionary/data_delete"
timestamp=$(date -u +%s)
signature=$(echo -n "${timestamp}${path}" | openssl dgst -sha256 -hmac "${secret_token}" | sed 's/^.* //')

curl -X POST ${url}${path} \
  -d '{"name": "name", "author": "author"}' \
  -H "Content-Type: application/json" \
  -H "X-Timestamp: ${timestamp}" \
  -H "X-Signature: ${signature}"
  
curl -X POST ${url}${path} \
  -d '{"id": "id"}' \
  -H "Content-Type: application/json" \
  -H "X-Timestamp: ${timestamp}" \
  -H "X-Signature: ${signature}"
```

# Request

```json
{
  "id": "dictionary id",
  "name": "dictionary name",
  "author": "author name"
}
```

# Response

```json
{
  "data": {
    "msg":"OK"
  }
}
```