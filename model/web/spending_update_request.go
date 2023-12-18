package web

type SpendingUpdateRequest struct {
	Id          string  `validate:"required,uuid4" json:"id"`
	Title       string  `validate:"required,min=3" json:"title"`
	Description string  `validate:"" json:"description"`
	Amount      float64 `validate:"required,gte=0" json:"amount"`
	Date        int64   `validate:"required" json:"date"`
	Category    string  `validate:"lowercase" json:"category"`
}
