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
cd /data/gen
go run dynamo_dictionary_table.go
```
