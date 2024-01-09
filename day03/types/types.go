package types

type Place struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Address string   `json:"address"`
	Phone   string   `json:"phone"`
	Loc     Location `json:"location"`
}

type Location struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type Total struct {
	Value int `json:"value"`
}

type Hits struct {
	Hits  []Hit `json:"hits"`
	Total Total `json:"total"`
}

type Responses struct {
	Hits Hits `json: "hits"`
}

type Hit struct {
	Hit Place `json:"_source"`
}
