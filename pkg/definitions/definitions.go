package definitions

type ValidationErr struct {
	Field   string   `json:"field"`
	Reasons []string `json:"reasons"`
}

type Token struct {
	Reference string
	Token     string
}

type Map map[string]interface{}

type PaginationResult struct {
	Limit      int           `json:"limit"`
	Page       int           `json:"page"`
	Sort       string        `json:"sort,omitempty;query:sort"`
	TotalItems int64         `json:"totalItems"`
	Items      []interface{} `json:"items"`
}
