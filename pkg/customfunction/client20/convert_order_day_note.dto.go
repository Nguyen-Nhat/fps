package funcClient20

type OrderDay struct {
	Weekly   string `json:"weekly"`
	Monthly  string `json:"monthly"`
	EveryDay bool   `json:"everyday"`
}
