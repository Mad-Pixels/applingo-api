package dictionaries

const (
	PageLimit = 40
)

type Item struct {
	Name        string `json:"name" dynamodbav:"name"`
	Category    string `json:"category" dynamodbav:"category"`
	Subcategory string `json:"subcategory" dynamodbav:"subcategory"`
	Author      string `json:"author" dynamodbav:"author"`
	Dictionary  string `json:"dictionary" dynamodbav:"dictionary"`
	Description string `json:"description" dynamodbav:"description"`
	CreatedAt   int    `json:"created_at" dynamodbav:"created_at"`
	Rating      int    `json:"rating" dynamodbav:"rating"`
	IsPublic    int    `json:"is_public" dynamodbav:"is_public"`
}

type GetResponse struct {
	Items         []Item `json:"items"`
	LastEvaluated string `json:"last_evaluated,omitempty"`
}

type PostRequest struct {
	Description string `json:"description"`
	Filename    string `json:"filename"`
	Name        string `json:"name"`
	Author      string `json:"author"`
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
	IsPublic    int    `json:"is_public"`
}
