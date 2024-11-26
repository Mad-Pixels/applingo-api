# Description

Lambda for manage dictionaries.

# Examples
## Define variables

```bash
api="op3r1apk2i"
url="http://localhost:4566/restapis/${api}/prod/_user_request_/v1/dictionaries"
```

```bash
timestamp=$(date -u +%s)
signature=$(echo -n "${timestamp}${arn_get}" | openssl dgst -sha256 -hmac "${token}" | sed 's/^.* //')

curl -X GET "${url}?is_public=true" -H "Content-Type: application/json"  \
    -H "x-timestamp: ${timestamp}" \
    -H "x-signature: ${signature}"

curl -X POST ${url} \
    -d '{"description": "description", "filename": "1.csv", "name": "test dictionary", "author": "author", "category": "language", "subcategory": "ru-he", "public": true}' \
    -H "Content-Type: application/json" 
```
