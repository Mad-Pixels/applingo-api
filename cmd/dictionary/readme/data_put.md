# Curl

```bash
# localstack
url="http://localhost:4566/restapis/1ysz0unl0u/prod/_user_request_"

secret_token="your_secret_token"
path="/api/v1/dictionary/data_put"
timestamp=$(date -u +%s)
signature=$(echo -n "${timestamp}${path}" | openssl dgst -sha256 -hmac "${secret_token}" | sed 's/^.* //')

curl -X POST ${url}${path} \
  -d '{"description": "description", "dictionary": "dictionary", "name": "name", "author": "author", "category_main": "category_main", "category_sub": "category_sub", "private": false}' \
  -H "Content-Type: application/json" \
  -H "X-Timestamp: ${timestamp}" \
  -H "X-Signature: ${signature}"
```

# Request

```json
{
  "description": "dictionary description",
  "dictionary": "{{ s3_url_for_download }}",
  "name": "dictionary name",
  "author": "author",
  "category_main": "main dictionary category",
  "category_sub": "dictionary subcategory",
  "private": "is dictionary is private?"
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