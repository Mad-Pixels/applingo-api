# Curl

```bash
# localstack
url="http://localhost:4566/restapis/4663mz3v89/prod/_user_request_"

secret_token="your_secret_token"
path="/api/v1/category/get"
timestamp=$(date -u +%s)
signature=$(echo -n "${timestamp}${path}" | openssl dgst -sha256 -hmac "${secret_token}" | sed 's/^.* //')

curl -X GET ${url}${path} \
  -d '{}' \
  -H "Content-Type: application/json" \
  -H "X-Timestamp: ${timestamp}" \
  -H "X-Signature: ${signature}"
```

# Response

```json
{
  "data":{
    "categories":[
      {
        "name":"language",
        "sub_categories": [
          "ru-en",
          "en-ru"
        ]
      }
    ]
  }
}
```