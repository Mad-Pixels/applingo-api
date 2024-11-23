# Description

Lambda for query dictionaries.

# Examples
## Define variables

```bash
api="xw66q4bfqv"

# localstack
url="http://localhost:4566/restapis/${api}/prod/_user_request_"
token="000XXX000"

device_path_query="device/v1/dictionary/query"
api_path_query="v1/dictionary/query"

device_download_url="device/v1/dictionary/download_url"
api_download_url="v1/dictionary/download_url"
```

## Define body
```bash
# Query by public
body='{"is_public": true, "sort_by": "date"}'
body='{"is_public": true, "sort_by": "rating"}'
body='{"is_public": true, "sort_by": "rating", "last_evaluated": "...."}'

# Query by sub category
body='{"subcategory": "ru-en", "sort_by": "date"}'
body='{"subcategory": "en-he", "sort_by": "rating"}'
body='{"subcategory": "en-he", "sort_by": "rating", "last_evaluated": "...."}'

# Get S3 file url
body='{"dictionary": "my_dictionary"}'
```

## device/v1/dictionary/query
```bash
timestamp=$(date -u +%s)
signature=$(echo -n "${timestamp}${arn_get}" | openssl dgst -sha256 -hmac "${token}" | sed 's/^.* //')
curl -X POST ${url}/${device_path_query} \
    -d "${body}" \
    -H "Content-Type: application/json" \
    -H "x-timestamp: ${timestamp}" \
    -H "x-signature: ${signature}"
```

## v1/dictionary/query
```bash
curl -X POST ${url}/${api_path_query} \
    -d "${body}" \
    -H "Content-Type: application/json"
```

## device/v1/dictionary/download_url
```bash
timestamp=$(date -u +%s)
signature=$(echo -n "${timestamp}" | openssl dgst -sha256 -hmac "${token}" | sed 's/^.* //')
curl -X POST ${url}/${device_download_url} \
    -d "${body}" \
    -H "Content-Type: application/json" \
    -H "x-timestamp: ${timestamp}" \
    -H "x-signature: ${signature}"
```

## v1/dictionary/download_url
```bash
curl -X POST ${url}/${api_download_url} \
    -d "${body}" \
    -H "Content-Type: application/json"
```
