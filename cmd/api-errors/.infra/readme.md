# Description

Lambda for aggregate errors.

# Examples
## Define variables

```bash
api="xw66q4bfqv"

# localstack
url="http://localhost:4566/restapis/${api}/prod/_user_request_"
token="000XXX000"

device_errors_put="device/v1/errors/put"
```

## Define body
```bash
body='{"app_version": "1.0", "device": "iPhone", "error_message": "Failed to load remote dictionaries", "error_original": "httpError(statusCode: 404)", "error_type": "api", "os_version": "17.5", "timestamp": 1731240938}'
```

## device/v1/errors/put
```bash
timestamp=$(date -u +%s)
signature=$(echo -n "${timestamp}${arn_get}" | openssl dgst -sha256 -hmac "${token}" | sed 's/^.* //')
curl -X POST ${url}/${device_errors_put} \
    -d "${body}" \
    -H "Content-Type: application/json" \
    -H "x-timestamp: ${timestamp}" \
    -H "x-signature: ${signature}"
```
