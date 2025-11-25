package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"fmt"
)

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
	Year   int    `json:"year" binding:"required,min=1000,max=2100"`
}

var (
	books   = make([]Book, 0)
	booksMu sync.Mutex
	nextID  = 1
)

func listBooks(c *gin.Context) {
	booksMu.Lock()
	defer booksMu.Unlock()
	c.JSON(http.StatusOK, books)
}

func getBook(c *gin.Context) {
	id := c.Param("id")
	booksMu.Lock()
	defer booksMu.Unlock()
	for _, b := range books {
		if b.ID == id {
			c.JSON(http.StatusOK, b)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
}

func createBook(c *gin.Context) {
	var input Book
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booksMu.Lock()
	defer booksMu.Unlock()
	input.ID = itoa(nextID)
	nextID++
	books = append(books, input)
	c.JSON(http.StatusCreated, input)
}

func updateBook(c *gin.Context) {
	id := c.Param("id")
	var input Book
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booksMu.Lock()
	defer booksMu.Unlock()
	for i, b := range books {
		if b.ID == id {
			input.ID = id
			books[i] = input
			c.JSON(http.StatusOK, input)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
}

func deleteBook(c *gin.Context) {
	id := c.Param("id")
	booksMu.Lock()
	defer booksMu.Unlock()
	for i, b := range books {
		if b.ID == id {
			books = append(books[:i], books[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
}

func itoa(i int) string {
	// simple int->string to avoid extra imports
	return fmt.Sprintf("%d", i)
}

func main() {
	router := gin.Default()

	booksGroup := router.Group("/books")
	{
		booksGroup.GET("", listBooks)
		booksGroup.GET("/:id", getBook)
		booksGroup.POST("", createBook)
		booksGroup.PUT("/:id", updateBook)
		booksGroup.DELETE("/:id", deleteBook)
	}

	router.Run(":8080")
}


