# Examples (localstack)
## Define variables

```bash
api="lg2iwpskaq"
url="http://localhost:4566/restapis/${api}/prod/_user_request_"
token="000XXX000"

device_path_query="device/v1/dictionary/query"
api_path_query="api/v1/dictionary/query"
```

## Define body
```bash
# Query by author
body='{"author": "author"}'

# Query by category_main
body='{"category_main": "category_main"}'

# Query by category_sub
body='{"category_sub": "category_sub"}'

# Query by is_private (true)
body='{"is_private": true}'

# Query by is_publish (true)
body='{"is_publish": true}'

# Query by author and is_private
body='{"author": "author", "is_private": true}'

# Query by category_main and is_private
body='{"category_main": "category_main", "is_private": false}'

# Query by category_sub and is_publish
body='{"category_sub": "category_sub", "is_publish": true}'
```

## device/v1/dictionary/query
```bash
arn_get="arn:aws:execute-api:us-east-1:000000000000:${api}/prod/POST/${device_path_query}"
timestamp=$(date -u +%s)

signature=$(echo -n "${timestamp}${arn_get}" | openssl dgst -sha256 -hmac "${token}" | sed 's/^.* //')
curl -X POST ${url}/${device_path_query} \
    -d "${body}" \
    -H "Content-Type: application/json" \
    -H "X-Timestamp: ${timestamp}" \
    -H "X-Signature: ${signature}"
```

## api/v1/dictionary/query
```bash
curl -X POST ${url}/${api_path_query} \
    -d "${body}" \
    -H "Content-Type: application/json"
```