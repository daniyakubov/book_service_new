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

func GetBookFromUser(userBook *UserBook) (book Book) {
	book.Id = userBook.Id
	book.Title = userBook.Title
	book.Author = userBook.Author
	book.Price = userBook.Price
	book.Available = userBook.Available
	book.Date = userBook.Date
	return book
}
