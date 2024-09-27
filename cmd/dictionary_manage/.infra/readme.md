# Examples (localstack)
## Define variables

```bash
api="ej6xoo4l3y"
url="http://localhost:4566/restapis/${api}/prod/_user_request_"

api_path_put="api/v1/dictionary/manage/put"
api_path_delete="api/v1/dictionary/manage/delete"
api_path_presign="api/v1/dictionary/manage/presign"
```

## api/v1/dictionary/manage/put
```bash
curl -X POST ${url}/${api_path_put} \
    -d '{"description": "description", "code":"", "dictionary": "dictionary", "name": "name", "author": "author", "category_main": "category_main", "category_sub": "category_sub", "private": false}' \
    -H "Content-Type: application/json" 
```
 
## api/v1/dictionary/manage/presign
```bash
curl -X POST ${url}/${api_path_presign} \
  -d '{"content_type": "text/csv", "name": "file.csv"}' \
  -H "Content-Type: application/json"
```

## api/v1/dictionary/manage/delete
```bash
curl -X POST ${url}/${api_path_delete} \
  -d '{"author": "author", "name": "name"}' \
  -H "Content-Type: application/json"
```