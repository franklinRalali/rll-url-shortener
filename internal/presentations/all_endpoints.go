package presentations

type Endpoint struct {
	Path        string `json:"path"`
	Method      string `json:"method"`
	Description string `json:"description"`
}

type AllEndpoints struct {
	Endpoints []Endpoint `json:"endpoints"`
}
