package model

type RequestDataTable struct {
	Search    string   `json:"search"`
	ArrSeacrh []string `json:"arr_search"`
	Filter    []struct {
		Column   string `json:"column"`
		Input    string `json:"input"`
		Operator string `json:"operator"`
	} `json:"filter"`
	Start      int                    `json:"start"`
	Length     int                    `json:"length"`
	OrderBy    string                 `json:"order_by"`
	OrderDesc  string                 `json:"order_desc"`
	Additional map[string]interface{} `json:"additional"`
}
type ResultDataTable struct {
	Count int64                    `json:"count"`
	Data  []map[string]interface{} `json:"data"`
}
