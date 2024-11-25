# Description

Lambda for getting bucket urls.

# Examples
## Define variables

```bash
api="ea9oxs8lq6"
url="http://localhost:4566/restapis/${api}/prod/_user_request_/v1/urls"
```

## Define body
```bash
body='{
  "operation": "download",
  "name": "1.csv"
}'

curl -X POST "${url}" -d "${body}" -H "Content-Type: application/json"
```
