package reports

type PostRequest struct {
	AppVersion     string         `json:"app_version"`
	OsVersion      string         `json:"os_version"`
	Device         string         `json:"device"`
	ErrorMessage   string         `json:"error_message"`
	ErrorOriginal  string         `json:"error_original"`
	ErrorType      string         `json:"error_type"`
	Timestamp      int            `json:"timestamp"`
	ReplicaID      string         `json:"replica_id"`
	AdditionalInfo map[string]any `json:"additional_info,omitempty"`
}
