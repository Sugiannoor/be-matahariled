package helpers

type PaginationData struct {
	Search string `json:"search,omitempty"`
	Page   int    `json:"page,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Sort   string `json:"sort,omitempty"`
	SortBy string `json:"sort_by,omitempty"`
}
