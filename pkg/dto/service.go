package dto

type UpdateResponse struct {
	Result bool   `json:"result"`
	Info   string `json:"info"`
	Code   int    `json:"code"`
}

type StateResponse struct {
	Result bool   `json:"result"`
	Info   string `json:"info"`
}

type NameResponse struct {
	UID       int
	FirstName string
	LastName  string
}
