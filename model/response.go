package model

type Response struct {
	Data       interface{} `json:"data"`
	Debug      interface{} `json:"debug"`
	Error      interface{} `json:"error"`
	Pagination interface{} `json:"pagination"`
}
