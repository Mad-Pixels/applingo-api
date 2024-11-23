package sort

// QueryType implement "sort by" field.
type QueryType string

const (
	QueryByDate   = "date"
	QueryByRating = "rating"
)

// IsValid check QueryType value.
func (q QueryType) IsValid() bool {
	switch q {
	case QueryByDate, QueryByRating:
		return true
	}
	return false
}

// String ...
func (q QueryType) String() string {
	return string(q)
}
