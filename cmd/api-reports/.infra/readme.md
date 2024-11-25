# Description

Lambda for aggregate errors.

# Examples
## Define variables

```bash
api="ea9oxs8lq6"
url="http://localhost:4566/restapis/${api}/prod/_user_request_/v1/reports"
```

## Define body
```bash
body='{
  "app_version": "1.0",
  "device": "iPhone",
  "error_message": "Failed to load remote dictionaries",
  "error_original": "httpError(statusCode: 404)",
  "error_type": "api",
  "os_version": "17.5",
  "timestamp": 1731240938
}'

curl -X POST "${url}" -d "${body}" -H "Content-Type: application/json"
```
