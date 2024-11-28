# Description

Lambda for query dictionary categories and sub-categories.

# Examples
## Define variables

```bash
api="op3r1apk2i"
url="http://localhost:4566/restapis/${api}/prod/_user_request_/v1/categories"
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
curl -X GET ${url} -H "Content-Type: application/json"
```