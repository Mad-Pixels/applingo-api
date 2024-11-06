# Description

Lambda for manage dictionaries.

# Examples
## Define variables

```bash
api="g9a3jrpzt1"

# localstack
url="http://localhost:4566/restapis/${api}/prod/_user_request_"

api_path_put="v1/dictionary/manage/put"
api_path_delete="v1/dictionary/manage/delete"
api_path_presign="v1/dictionary/manage/upload_url"
```

## v1/dictionary/manage/put
```bash
# public object
curl -X POST ${url}/${api_path_put} \
    -d '{"description": "description", "code":"", "filename": "dictionary", "name": "name", "author": "author", "category_main": "category_main", "category_sub": "category_sub", "is_public": true}' \
    -H "Content-Type: application/json" 

curl -X POST ${url} \
    -d '{"description": "description", "filename": "revert.csv", "name": "name6", "author": "author", "category": "category_main", "subcategory": "ru-he", "is_public": true}' \
    -H "Content-Type: application/json" 


type handleDataPutRequest struct {
	Description string `json:"description" validate:"required"`
	Filename    string `json:"filename" validate:"required"`
	Name        string `json:"name" validate:"required,min=4,max=32"`
	Author      string `json:"author" validate:"required"`
	Category    string `json:"category" validate:"required"`
	Subcategory string `json:"subcategory" validate:"required"`
	IsPublic    bool   `json:"is_public" validate:"required"`
}

# private object
curl -X POST ${url}/${api_path_put} \
    -d '{"description": "description", "code":"666", "filename": "dictionary", "name": "name", "author": "author", "category_main": "category_main", "category_sub": "category_sub", "is_public": false}' \
    -H "Content-Type: application/json" 
```

## v1/dictionary/manage/presign
```bash
curl -X POST ${url}/${api_path_presign} \
  -d '{"content_type": "text/csv", "name": "file.csv"}' \
  -H "Content-Type: application/json"
```

## v1/dictionary/manage/delete
```bash
curl -X POST ${url}/${api_path_delete} \
  -d '{"author": "author", "name": "name"}' \
  -H "Content-Type: application/json"
```