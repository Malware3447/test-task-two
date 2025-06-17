package request

type Request struct {
	Uuid      string `json:"uuid"`
	Operation string `json:"operation"`
	Amount    int    `json:"amount"`
}
