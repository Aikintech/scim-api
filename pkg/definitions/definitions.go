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
