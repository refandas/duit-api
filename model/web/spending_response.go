package web

type SpendingResponse struct {
	Id          string  `json:"id"`
	UserId      string  `json:"user_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Date        int64   `json:"date"`
	Category    string  `json:"category"`
	CreatedAt   int64   `json:"created_at"`
}
