package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type Genre string

const (
	COMEDY    Genre = "Comedy"
	DRAMA           = "Drama"
	NOVEL           = "Novel"
	ACTION          = "Action"
	ADVENTURE       = "Adventure"
	FICTION         = "Fiction"
	HORROR          = "Horror"
	TALE            = "Fairy tale"
)

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
}

func (u *User) update(n User) {
	if len(n.Name) > 0 {
		u.Name = n.Name
	}
	if len(n.Surname) > 0 {
		u.Surname = n.Surname
	}
	if len(n.Email) > 0 {
		u.Email = n.Email
	}
}

type Rental struct {
	ID       int    `json:"id"`
	Reader   User   `json:"reader"`
	Book     Book   `json:"book"`
	RentDate string `json:"rent_date"`
	DueDate  string `json:"due_date"`
	Fine     Fine   `json:"fine"`
}

type Fine struct {
	Amount float32 `json:"amount"`
	Date   string  `json:"date"`
	Status bool    `json:"status"`
}

func (f *Fine) update(n Fine) {
	if n.Amount > 0 {
		f.Amount = n.Amount
	}
	f.Status = n.Status
}

type Book struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Authors     []string `json:"authors"`
	Date        string   `json:"date"`
	PagesCount  int16    `json:"pages_count"`
	Genres      []Genre  `json:"genres"`
	IsAvailable bool     `json:"is_available"`
}

func (b *Book) update(u Book) {
	if len(u.Title) > 0 {
		b.Title = u.Title
	}
	if len(u.Authors) > 0 {
		b.Authors = u.Authors
	}
	if len(u.Date) > 0 {
		b.Date = u.Date
	}
	if len(u.Genres) > 0 {
		b.Genres = u.Genres
	}
	if u.PagesCount > 0 {
		b.PagesCount = u.PagesCount
	}
	b.IsAvailable = u.IsAvailable
}

var users = []User{
	{ID: 1, Name: "Валериан", Surname: "Волков", Email: "hlazarev@example.org"},
	{ID: 2, Name: "Полина", Surname: "Лапина", Email: "zkrylova@example.org"},
	{ID: 3, Name: "Валентин", Surname: "Зиновьев", Email: "belaev.german@example.net"},
	{ID: 4, Name: "Захар", Surname: "Захаров", Email: "masnikov.florentina@example.com"},
	{ID: 5, Name: "Игорь", Surname: "Фомин", Email: "vrybakov@example.com"},
	{ID: 6, Name: "Марат", Surname: "Лукин", Email: "dara05@example.net"},
	{ID: 7, Name: "Вероника", Surname: "Сысоева", Email: "olga44@example.org"},
	{ID: 8, Name: "Ефим", Surname: "Виноградов", Email: "kozlov.antonina@example.com"},
	{ID: 9, Name: "Виталий", Surname: "Суворов", Email: "gblohina@example.com"},
	{ID: 10, Name: "Вера", Surname: "Кириллова", Email: "hrozkov@example.com"},
}

var books = []Book{
	{ID: 1, Title: "Сто лет одиночества", Authors: []string{"Габриэль Гарсиа Маркес"}, Date: "1967", PagesCount: 480, Genres: []Genre{DRAMA, NOVEL}, IsAvailable: true},
	{ID: 2, Title: "451 градус по Фаренгейту", Authors: []string{"Рэй Брэдбери"}, Date: "1967", PagesCount: 200, Genres: []Genre{FICTION, NOVEL}, IsAvailable: true},
	{ID: 3, Title: "Повелитель мух", Authors: []string{"Уильям Голдинг"}, Date: "1954", PagesCount: 190, Genres: []Genre{ADVENTURE, NOVEL}, IsAvailable: true},
	{ID: 4, Title: "Мастер и Маргарита", Authors: []string{"Михаил Булгаков"}, Date: "1967", PagesCount: 470, Genres: []Genre{NOVEL}, IsAvailable: true},
	{ID: 5, Title: "На Западном фронте без перемен", Authors: []string{"Эрих Мария Ремарк"}, Date: "1929", PagesCount: 200, Genres: []Genre{NOVEL}, IsAvailable: true},
	{ID: 6, Title: "Солярис", Authors: []string{"Станислав Лем"}, Date: "1961", PagesCount: 210, Genres: []Genre{FICTION, NOVEL}, IsAvailable: true},
	{ID: 7, Title: "Война и мир", Authors: []string{"Лев Толстой"}, Date: "1868", PagesCount: 1979, Genres: []Genre{NOVEL}, IsAvailable: true},
	{ID: 8, Title: "Маленький принц", Authors: []string{"Антуан де Сент-Экзюпери"}, Date: "1943", PagesCount: 78, Genres: []Genre{TALE}, IsAvailable: true},
	{ID: 9, Title: "Над пропастью во ржи", Authors: []string{"Дж. Д. Сэлинджер"}, Date: "1951", PagesCount: 240, Genres: []Genre{NOVEL}, IsAvailable: true},
	{ID: 10, Title: "Анна Каренина", Authors: []string{"Лев Толстой"}, Date: "1875", PagesCount: 1081, Genres: []Genre{NOVEL}, IsAvailable: true},
}

var rentals = []Rental{
	{ID: 1, Reader: users[0], Book: books[3], RentDate: "13.10.2024", DueDate: "30.11.2024"},
	{ID: 2, Reader: users[8], Book: books[2], RentDate: "05.10.2024", DueDate: "14.11.2024"},
	{ID: 3, Reader: users[9], Book: books[7], RentDate: "01.11.2024", DueDate: "22.12.2024"},
	{ID: 4, Reader: users[3], Book: books[5], RentDate: "20.10.2024", DueDate: "13.12.2024"},
	{ID: 5, Reader: users[2], Book: books[9], RentDate: "01.10.2024", DueDate: "23.11.2024"},
}

func findBookIdxById(id int) int {
	for i, book := range books {
		if book.ID == id {
			return i
		}
	}
	return -1
}

func findUserIdxById(id int) int {
	for i, user := range users {
		if user.ID == id {
			return i
		}
	}
	return -1
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.POST("/books", addBooks)
	router.GET("/books/:id", getBooksById)
	router.PUT("/books/:id", updateBook)
	router.DELETE("/books/:id", deleteBook)

	router.GET("/users", getUsers)
	router.POST("/users", addUser)
	router.GET("/users/:id", findUserById)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)

	router.GET("/rentals", getRentals)
	router.POST("/rentals", addRental)
	router.GET("/rentals/:id", findRentalById)
	router.GET("/find_rental", findRentalByContent)
	router.DELETE("/rentals/:id", deleteRental)

	router.POST("/rentals/:id/fine", addFine)
	router.PUT("/rentals/:id/fine", payFine)
	router.DELETE("/rentals/:id/fine", removeFine)
	router.Run(":8080")
}

func getBooks(c *gin.Context) {
	c.JSON(http.StatusOK, books)
}

func getBooksById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, book := range books {
		if book.ID == id {
			c.JSON(http.StatusOK, book)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "book not found"})
}

func addBooks(c *gin.Context) {
	var newBook Book
	if err := c.BindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	books = append(books, newBook)
	c.JSON(http.StatusCreated, newBook)
}

func updateBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var newData Book
	if err := c.BindJSON(&newData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	println(newData.IsAvailable)
	for i := 0; i < len(books); i++ {
		if books[i].ID == id {
			books[i].update(newData)
			c.JSON(http.StatusOK, books[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "book not found"})
}

func deleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "book deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "book not found"})
}

func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users)
}

func findUserById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, user := range users {
		if user.ID == id {
			c.JSON(http.StatusOK, user)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

func addUser(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	users = append(users, newUser)
	c.JSON(http.StatusCreated, newUser)
}

func updateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	for i := 0; i < len(users); i++ {
		if users[i].ID == id {
			users[i].update(newUser)
			c.JSON(http.StatusOK, users[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

func deleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

func getRentals(c *gin.Context) {
	c.JSON(http.StatusOK, rentals)
}

func findRentalById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, rental := range rentals {
		if rental.ID == id {
			c.JSON(http.StatusOK, rental)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "rental not found"})
}

func findRentalByContent(c *gin.Context) {
	userId, err1 := strconv.Atoi(c.Query("userId"))
	bookId, err2 := strconv.Atoi(c.Query("bookId"))
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong parameters"})
		return
	}
	for _, rental := range rentals {
		if rental.Book.ID == bookId && rental.Reader.ID == userId {
			c.JSON(http.StatusOK, rental)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "rental not found"})
}

func addRental(c *gin.Context) {
	userId, err1 := strconv.Atoi(c.Query("userId"))
	bookId, err2 := strconv.Atoi(c.Query("bookId"))
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong parameters"})
		return
	}
	bookIdx := findBookIdxById(bookId)
	userIdx := findUserIdxById(userId)
	if bookId < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book not found"})
		return
	}
	if books[bookIdx].IsAvailable == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book not available"})
		return
	}
	if userId < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	var newRental Rental
	newRental.ID = rentals[len(rentals)-1].ID + 1
	newRental.Reader = users[userIdx]
	books[bookIdx].IsAvailable = false
	newRental.Book = books[bookIdx]
	newRental.RentDate = time.Now().String()
	newRental.DueDate = time.Now().AddDate(0, 1, 0).String()
	rentals = append(rentals, newRental)
	c.JSON(http.StatusCreated, newRental)
}

func deleteRental(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, rental := range rentals {
		if rental.ID == id {
			rental.Book.IsAvailable = false
			rentals = append(rentals[:i], rentals[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "rental deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "rental not found"})
}

func addFine(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i := 0; i < len(rentals); i++ {
		if rentals[i].ID == id {
			rentals[i].Fine = Fine{Amount: 1000.0, Date: time.Now().String(), Status: false}
			c.JSON(http.StatusOK, rentals[i].Fine)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "rental not found"})
}

func payFine(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i := 0; i < len(rentals); i++ {
		if rentals[i].ID == id {
			if rentals[i].Fine.Date == "" {
				c.JSON(http.StatusNotFound, gin.H{"message": "fine not found"})
				return
			}
			rentals[i].Fine.Status = true
			c.JSON(http.StatusOK, rentals[i].Fine)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "rental not found"})
}

func removeFine(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i := 0; i < len(rentals); i++ {
		if rentals[i].ID == id {
			if rentals[i].Fine.Date == "" {
				c.JSON(http.StatusNotFound, gin.H{"message": "fine not found"})
				return
			}
			rentals[i].Fine = Fine{}
			c.JSON(http.StatusOK, rentals[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "rental not found"})
}
