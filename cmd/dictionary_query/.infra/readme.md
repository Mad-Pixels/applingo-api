# Examples (localstack)
## Define variables

```bash
api="ej6xoo4l3y"
url="http://localhost:4566/restapis/${api}/prod/_user_request_"
token="000XXX000"

device_path_query="device/v1/dictionary/query"
api_path_query="api/v1/dictionary/query"
```

## Define body
```bash
# Query by category_main (public)
body='{"category_main": "category_main"}'

# Query by category_main (private)
body='{"category_main": "category_main", "code": "666"}'

# Query by category_sub (public)
body='{"category_sub": "category_sub"}'

# Query by category_sub (private)
body='{"category_sub": "category_sub", "code": "666"}'

# Query by public raws
body='{"is_public": true}'
body='{}'

# Query by private raws
body='{"is_public": false, "code": "666"}'

# Query by name, author
body='{"name": "name","author": "author"}'

# Query by name
body='{"name": "name"}'

# Query by author
body='{"author": "author"}'
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