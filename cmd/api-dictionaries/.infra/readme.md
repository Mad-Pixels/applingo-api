# Description

Lambda for manage dictionaries.

# Examples
## Define variables

```bash
token="000XXX000"
api="vp0yxvxfow"
url="http://localhost:4566/restapis/${api}/prod/_user_request_/v1/dictionaries"
```

```bash
timestamp=$(date -u +%s)
signature=$(echo -n "${timestamp}${arn_get}" | openssl dgst -sha256 -hmac "${token}" | sed 's/^.* //')

curl -X GET "${url}" -H "Content-Type: application/json"  \
    -H "x-api-auth: ${timestamp}:::${signature}" 

curl -X GET "${url}" 

timestamp=$(date -u +%s)
signature=$(echo -n "${timestamp}${arn_get}" | openssl dgst -sha256 -hmac "${token}" | sed 's/^.* //')
curl -X POST ${url} \
    -d '{"description": "description", "filename": "1.csv", "name": "testdictionary", "author": "author", "category": "language", "subcategory": "ru-il", "public": true, "level": "A1", "topic":"topic"}' \
    -H "Content-Type: application/json" \
    -H "x-timestamp: ${timestamp}" \
    -H "x-signature: ${signature}" 
```
