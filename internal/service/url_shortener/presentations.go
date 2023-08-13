package urlshortener

type ShortenURLReq struct {
	UserID    string
	OriginURL string
}

type ShortenURLRes struct {
	ID        uint64 `json:"id"`
	ShortCode string `json:"short_code"`
}

type ShortURLRes struct {
	OriginURL string `json:"origin_url"`
	ShortURL  string `json:"short_url"`
}

type ShortURLUpdateReq struct {
	OriginURL string `json:"origin_url"`
}

type StatisticsRes struct {
	VisitCount uint64 `json:"visit_count"`
}
