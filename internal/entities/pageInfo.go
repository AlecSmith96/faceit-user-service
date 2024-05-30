package entities

type PageInfo struct {
	NextPageToken string `json:"next_page_token"`
	PageSize      int    `json:"page_size"`
}
