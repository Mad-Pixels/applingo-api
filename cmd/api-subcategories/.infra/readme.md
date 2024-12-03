# Description

Lambda for query dictionary categories and sub-categories.

# Examples
## Define variables

```bash
api="8sxucvpidu"
url="http://localhost:4566/restapis/${api}/prod/_user_request_/v1/subcategories"
```

## device/v1/category/query
```bash
timestamp=$(date -u +%s)
signature=$(echo -n "${timestamp}" | openssl dgst -sha256 -hmac "${token}" | sed 's/^.* //')
curl -X GET ${url} \
    -d "${body}" \
    -H "Content-Type: application/json" \
     -H "x-api-auth: ${timestamp}:::${signature}" 
```

## v1/category/query
```bash
curl -X GET ${url} -H "Content-Type: application/json"

curl -X POST ${url} \
-H "Content-Type: application/json" \
-d '{
  "side": "front",
  "code": "tn",
  "description": "nrJkOU^hrPF"
}'
curl -X DELETE "${url}?side=front&codes=tn"
```