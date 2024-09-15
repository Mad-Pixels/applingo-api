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
curl -X POST http://localhost:4566/restapis/l60qosmxou/prod/_user_request_/api/v1/dictionary/file_presign \
  -d '{"content_type": "text/csv", "name": "file.csv"}' \
  -H "Content-Type: application/json"

curl -X POST http://localhost:4566/restapis/l60qosmxou/prod/_user_request_/api/v1/dictionary/data_put \
  -d '{"description": "description", "dictionary": "dictionary", "name": "name", "author": "author", "category": "category", "sub_category": "sub_category", "private": false}' \
  -H "Content-Type: application/json"

curl -X GET http://localhost:4566/restapis/l60qosmxou/prod/_user_request_/api/v1/category/get \
  -d '{"content_type": "text/csv", "name": "file.csv"}' \
  -H "Content-Type: application/json"
```