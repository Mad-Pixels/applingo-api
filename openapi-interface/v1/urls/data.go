package urls

const (
	ExpiresIn = 30
)

type Operation string

const (
	OperationUpload   Operation = "upload"
	OperationDownload Operation = "download"
)

type ContentType string

const (
	ContentTypeCSV  ContentType = "text/csv"
	ContentTypeXLSX ContentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	ContentTypeXLS  ContentType = "application/vnd.ms-excel"
)

type PostRequest struct {
	Operation   Operation   `json:"operation"`
	ContentType ContentType `json:"content_type,omitempty"`
	Name        string      `json:"name,omitempty"`
}

type PostResponse struct {
	URL       string `json:"url"`
	ExpiresIn int    `json:"expires_in"`
}
