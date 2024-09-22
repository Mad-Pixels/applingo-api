# Curl

```bash
# localstack
url="http://localhost:4566/restapis/4663mz3v89/prod/_user_request_"

curl -X POST ${url}/api/v1/dictionary/file_presign \
  -d '{"content_type": "text/csv", "name": "file.csv"}' \
  -H "Content-Type: application/json"
```

# Request

```json
{
  "content_type": "text/csv",
  "name": "filename"
}
```

# Response

```json
{
  "data":{
    "url":"http://172.17.0.2:4566/lingocards-processing/file.csv?X-Amz-Algorithm=AWS4-HMAC-SHA256..."
  }
}
```