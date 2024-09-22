# Curl

```bash
# localstack
url="http://localhost:4566/restapis/4663mz3v89/prod/_user_request_"

curl -X GET ${url}/api/v1/category/get \
  -d '{}' \
  -H "Content-Type: application/json"
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