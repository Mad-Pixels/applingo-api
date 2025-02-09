# Description

Lambda for manage dictionaries.

# Examples
## Define variables

```bash
token="000XXX000"
api="xa716jsuuh"
url="http://localhost:4566/restapis/${api}/live/_user_request_/v1/dictionaries"
```

```bash
timestamp=$(date -u +%s)
signature=$(echo -n "${timestamp}${arn_get}" | openssl dgst -sha256 -hmac "${token}" | sed 's/^.* //')

curl -X PATCH "${url}?name=hello&author=author&subcategory=subcategory" -H "Content-Type: application/json"  \
    -H "x-api-auth: ${timestamp}:::${signature}" \
    -H "x-operation-name: patchStatisticDictionariesV1" \
    -d '{"downloads":"increase", "rating":"increase", }' 

timestamp=$(date -u +%s)
signature=$(echo -n "${timestamp}${arn_get}" | openssl dgst -sha256 -hmac "${token}" | sed 's/^.* //')

curl -X GET "${url}" -H "x-api-auth: ${timestamp}:::${signature}"

timestamp=$(date -u +%s)
signature=$(echo -n "${timestamp}${arn_get}" | openssl dgst -sha256 -hmac "${token}" | sed 's/^.* //')
curl -X POST ${url} \
    -d '{"description": "description", "filename": "1.csv", "name": "testdictionary", "author": "author", "category": "language", "subcategory": "ru-il", "public": true, "level": "A1", "topic":"topic"}' \
    -H "Content-Type: application/json" \
    -H "x-timestamp: ${timestamp}" \
    -H "x-signature: ${signature}" 
```
