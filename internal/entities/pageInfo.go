package entities

type PageInfo struct {
	NextPageToken string `json:"next_page_token"`
	PageSize      int32  `json:"page_size"`
	SortBy        string `json:"sort_by"`
	Descending    bool   `json:"descending"`
}
