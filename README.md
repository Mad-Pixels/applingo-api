# lingocards-api


graph TD
    A[Web Form] -->|Submit| B[API Gateway]
    B --> C[Lambda - Form Handler]
    C -->|Upload file| D[S3 Temp Bucket]
    D -->|Trigger| E[Lambda - File Processor]
    E -->|Convert to CSV| F[S3 Processed CSV Bucket]
    E -->|If successful| G[Lambda - Data Persister]
    G -->|Store metadata| H[DynamoDB]
    G -->|Move file| I[S3 Final Bucket]
    E -->|If failed| J[Delete from Temp Bucket]
    C -->|Return presigned URL| K[Client]
    K -->|Upload directly| D


```bash
curl -X POST http://localhost:4566/restapis/4663mz3v89/prod/_user_request_/api/v1/dictionary/file_presign \
  -d '{"content_type": "text/csv", "name": "file.csv"}' \
  -H "Content-Type: application/json"

curl -X POST http://localhost:4566/restapis/4663mz3v89/prod/_user_request_/api/v1/dictionary/data_put \
  -d '{"description": "description", "dictionary": "dictionary", "name": "name", "author": "author", "category_main": "category_main", "category_sub": "category_sub", "private": false}' \
  -H "Content-Type: application/json"

curl -X GET http://localhost:4566/restapis/4663mz3v89/prod/_user_request_/api/v1/category/get \
  -d '{"content_type": "text/csv", "name": "file.csv"}' \
  -H "Content-Type: application/json"
  
curl -X POST http://localhost:4566/restapis/4663mz3v89/prod/_user_request_/api/v1/dictionary/data_get \
  -H "Content-Type: application/json"   -d '{"author": "author", "is_private":false}'
```

```bash
cd /data/gen
go run dynamo_dictionary_table.go
```