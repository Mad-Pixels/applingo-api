# Description

Lambda for query dictionary categories and sub-categories.

# Examples
## Define variables

```bash
api="uv9cga7kcs"

# localstack
url="http://localhost:4566/restapis/${api}/prod/_user_request_"
token="000XXX000"

device_path_query="device/v1/category/query"
api_path_query="v1/category/query"
```

## Define body
```bash
# empty
body='{}'
```

## device/v1/category/query
```bash
timestamp=$(date -u +%s)
signature=$(echo -n "${timestamp}" | openssl dgst -sha256 -hmac "${token}" | sed 's/^.* //')
curl -X POST ${url}/${device_path_query} \
    -d "${body}" \
    -H "Content-Type: application/json" \
    -H "x-timestamp: ${timestamp}" \
    -H "x-signature: ${signature}"
```

## v1/category/query
```bash
curl -X POST ${url}/${api_path_query} \
    -d "${body}" \
    -H "Content-Type: application/json"
```