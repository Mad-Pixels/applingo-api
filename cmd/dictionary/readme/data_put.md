# Curl

```bash
# localstack
url="http://localhost:4566/restapis/4663mz3v89/prod/_user_request_"

curl -X POST ${url}/api/v1/dictionary/data_put \
  -d '{"description": "description", "dictionary": "dictionary", "name": "name", "author": "author", "category_main": "category_main", "category_sub": "category_sub", "private": false}' \
  -H "Content-Type: application/json"
```

# Request

```json
{
  "description": "dictionary description",
  "dictionary": "{{ s3_url_for_download }}",
  "name": "dictionary name",
  "author": "author",
  "category_main": "main dictionary category",
  "category_sub": "dictionary subcategory",
  "private": "is dictionary is private?"
}
```

# Response

```json
{
  "data": {
    "msg":"OK"
  }
}
```