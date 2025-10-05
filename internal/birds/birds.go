package birds

type BirdsJson struct {
	Birds    []Bird `json:"entities"`
	Total    int    `json:"total"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}
type Bird struct {
	Images      []string `json:"images"`
	LengthMin   string   `json:"lengthMin"`
	LengthMax   string   `json:"lengthMax"`
	Name        string   `json:"name"`
	ID          int      `json:"id"`
	SciName     string   `json:"sciName"`
	Region      []string `json:"region"`
	Family      string   `json:"family"`
	Order       string   `json:"order"`
	Status      string   `json:"status"`
	WingspanMin string   `json:"wingspanMin,omitempty"`
	WingspanMax string   `json:"wingspanMax,omitempty"`
}
