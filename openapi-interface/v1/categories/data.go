package categories

type Item struct {
	Name string `json:"name"`
}

type GetResponse struct {
	FrontCategory []Item `json:"front_category"`
	BackCategory  []Item `json:"back_category"`
}
