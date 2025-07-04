package utils

type Pagination struct {
	Page         int32 `json:"page"`
	Limit        int32 `json:"limit"`
	TotalRecords int32 `json:"total_records"`
	TotalPages   int32 `json:"total_pages"`
	HasNext      bool  `json:"has_next"`
	HasPrev      bool  `json:"has_prev"`
}

func NewPagiantion(page, limit, totalRecords int32) *Pagination {
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limitInt := GetIntEnv("LIMIT_ITEM_ON_PER_PAGE", 10)
		limit = int32(limitInt)
	}

	totalPages := (totalRecords + limit - 1) / limit

	return &Pagination{
		Page:         page,
		Limit:        limit,
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		HasNext:      page < totalPages,
		HasPrev:      page > 1,
	}
}

func NewPagiantionResponse(data any, page, limit, totalRecords int32) map[string]any {
	return map[string]any{
		"data":       data,
		"pagination": NewPagiantion(page, limit, totalRecords),
	}
}
