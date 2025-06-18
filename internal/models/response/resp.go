package response

type Response struct {
	Uuid   string `json:"uuid"`
	Amount int    `json:"amount"`
	Error  string `json:"error"`
}
