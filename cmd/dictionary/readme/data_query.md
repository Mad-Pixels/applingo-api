# Curl

```bash
# localstack
url="http://localhost:4566/restapis/4663mz3v89/prod/_user_request_"

secret_token="your_secret_token"
path="/api/v1/dictionary/data_query"
timestamp=$(date -u +%s)
signature=$(echo -n "${timestamp}${path}" | openssl dgst -sha256 -hmac "${secret_token}" | sed 's/^.* //')

# Query by author
curl -X POST ${url}${path} \
  -d '{"author": "author_name"}' \
  -H "Content-Type: application/json" \
  -H "X-Timestamp: ${timestamp}" \
  -H "X-Signature: ${signature}"

# Query by category_main
curl -X POST ${url}${path} \
  -d '{"category_main": "main_category"}' \
  -H "Content-Type: application/json" \
  -H "X-Timestamp: ${timestamp}" \
  -H "X-Signature: ${signature}"

# Query by category_sub
curl -X POST ${url}${path} \
  -d '{"category_sub": "subcategory"}' \
  -H "Content-Type: application/json" \
  -H "X-Timestamp: ${timestamp}" \
  -H "X-Signature: ${signature}"

# Query by is_private (true)
curl -X POST ${url}${path} \
  -d '{"is_private": true}' \
  -H "Content-Type: application/json" \
  -H "X-Timestamp: ${timestamp}" \
  -H "X-Signature: ${signature}"

# Query by is_publish (true)
curl -X POST ${url}${path} \
  -d '{"is_publish": true}' \
  -H "Content-Type: application/json" \
  -H "X-Timestamp: ${timestamp}" \
  -H "X-Signature: ${signature}"

# Query by author and is_private
curl -X POST ${url}${path} \
  -d '{"author": "author_name", "is_private": true}' \
  -H "Content-Type: application/json" \
  -H "X-Timestamp: ${timestamp}" \
  -H "X-Signature: ${signature}"

# Query by category_main and is_private
curl -X POST ${url}${path} \
  -d '{"category_main": "main_category", "is_private": false}' \
  -H "Content-Type: application/json" \
  -H "X-Timestamp: ${timestamp}" \
  -H "X-Signature: ${signature}"

# Query by category_sub and is_publish
curl -X POST ${url}${path} \
  -d '{"category_sub": "subcategory", "is_publish": true}' \
  -H "Content-Type: application/json" \
  -H "X-Timestamp: ${timestamp}" \
  -H "X-Signature: ${signature}"

# Query with pagination using last_evaluated
curl -X POST ${url}${path} \
  -d '{"author": "author_name", "last_evaluated": "eyJhdXRob3IiOnsiVmFsdWUiOiJhdXRob3IifSwiaWQiOnsiVmFsd..."}' \
  -H "Content-Type: application/json" \
  -H "X-Timestamp: ${timestamp}" \
  -H "X-Signature: ${signature}"
```

# Request

```json
{
  "id": "dictionary id",
  "name": "dictionary name",
  "author": "author name",
  "category_main": "main category",
  "category_sub": "subcategory",
  "is_private": true,
  "is_publish": true,
  "last_evaluated": "last evaluated key for pagination"
}
```

# Response

```json
{
  "data": {
    "items": [
      {
        "id": "38b1b42b56acb4f502034aefd4c467ac",
        "name": "name",
        "author": "author",
        "category_main": "category_main",
        "category_sub": "category_sub",
        "is_private": 1,
        "is_publish": 1
      },
      // ... more items
    ],
    "last_evaluated": "eyJhdXRob3IiOnsiVmFsdWUiOiJhdXRob3IifSwiaWQiOnsiVmFsd..."
  }
}
```