package models

type Book struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	Author      string  `json:"author"`
	Available   bool    `json:"available"`
	Date        string  `json:"date,omitempty"`
	PublishDate string  `json:"publish,date,omitempty"`
}

type RHit struct {
	Index  string  `json:"_index"`
	Type   string  `json:"_type"`
	Id     string  `json:"_id"`
	Score  float64 `json:"_score"`
	Source Book    `json:"_source"`
}

type StoreResponse struct {
	Count       int64 `json:"books_num"`
	DistinctAut int   `json:"distinct_authors_num"`
}

type PutResponse struct {
	Id string `json:"id"`
}
