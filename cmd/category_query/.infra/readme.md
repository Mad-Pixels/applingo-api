# Examples (localstack)
## Define variables

```bash
api="lg2iwpskaq"
url="http://localhost:4566/restapis/${api}/prod/_user_request_"
token="000XXX000"

device_path_query="device/v1/category/query"
api_path_query="api/v1/category/query"
```

## Define body
```bash
# empty
body='{}'
```

## device/v1/category/query
```bash
arn_get="arn:aws:execute-api:us-east-1:000000000000:${api}/prod/GET/${device_path_query}"
timestamp=$(date -u +%s)

signature=$(echo -n "${timestamp}${arn_get}" | openssl dgst -sha256 -hmac "${token}" | sed 's/^.* //')
curl -X GET ${url}/${device_path_query} \
    -d "${body}" \
    -H "Content-Type: application/json" \
    -H "X-Timestamp: ${timestamp}" \
    -H "X-Signature: ${signature}"
```

## api/v1/category/query
```bash
curl -X GET ${url}/${api_path_query} \
    -d "${body}" \
    -H "Content-Type: application/json"
```