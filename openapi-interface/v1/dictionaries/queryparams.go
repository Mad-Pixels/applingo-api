package dictionaries

import (
	"fmt"

	"github.com/Mad-Pixels/applingo-api/openapi-interface"
)

type QueryParams struct {
	openapi.QueryParams
	sortBy        string
	subcategory   string
	lastEvaluated string
	isPublic      *bool
}

// NewQueryParams теперь принимает openapi.QueryParams вместо map[string]string
func NewQueryParams(base openapi.QueryParams) (*QueryParams, error) {
	qp := &QueryParams{QueryParams: base}

	// Используем GetStringDefault для параметров, которые могут отсутствовать
	qp.sortBy = base.GetStringDefault("sort_by", "")
	qp.subcategory = base.GetStringDefault("subcategory", "")
	qp.lastEvaluated = base.GetStringDefault("last_evaluated", "")

	// Парсим isPublic только если он указан
	if base.Has("is_public") {
		isPublic, err := base.GetBool("is_public")
		if err != nil {
			return nil, fmt.Errorf("invalid is_public parameter: %w", err)
		}
		qp.isPublic = &isPublic
	}

	return qp, nil
}

func (q *QueryParams) SortBy() string        { return q.sortBy }
func (q *QueryParams) Subcategory() string   { return q.subcategory }
func (q *QueryParams) LastEvaluated() string { return q.lastEvaluated }
func (q *QueryParams) IsPublic() *bool       { return q.isPublic }
