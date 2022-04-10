package models

type Book struct {
	Id        string  `json:"id"`
	Title     string  `json:"title"`
	Author    string  `json:"author"`
	Price     float64 `json:"price"`
	Available bool    `json:"available"`
	Date      string  `json:"date"`
}

type UserBook struct {
	Id        string  `json:"id"`
	Title     string  `json:"title"`
	Author    string  `json:"author"`
	Price     float64 `json:"price"`
	Available bool    `json:"available"`
	Date      string  `json:"date"`
	Username  string  `json:"username"`
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

func GetBookFromUser(userBook *UserBook) (book Book) {
	book.Id = userBook.Id
	book.Title = userBook.Title
	book.Author = userBook.Author
	book.Price = userBook.Price
	book.Available = userBook.Available
	book.Date = userBook.Date
	return book
}
