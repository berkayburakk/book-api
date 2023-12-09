package controllers

import (
	"book-api/initiliazers"
	"book-api/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"time"
)

// CreateBook godoc
// @Summary Create Book
// @Description Save book data in db.
// @Accept json
// @Produce json
// @Param book body models.Book.ID true "Book Create Data"
// @Success 201 {object} models.Book
// @Failure 400 {object} map[string]interface{}
// @Router /books [post]
func CreateBook(c *gin.Context) {
	var body struct {
		BookName string
		Barcode  string
		Author   string
		Category string
	}

	err := c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if checkMissingFields(c, body, errorMessages) {
		return
	}

	if barcodeExists(body.Barcode) {
		c.JSON(http.StatusBadRequest, gin.H{"error": BarcodeAlreadyExist})
		return
	}

	book := models.Book{BookName: body.BookName, Barcode: body.Barcode, Author: body.Author, Category: body.Category, CreatedDate: time.Now(), UpdatedDate: time.Now()}

	result := initiliazers.DB.Create(&book)

	if result.Error != nil {
		c.Status(400)
		return
	}

	response := gin.H{
		"id":       book.ID,
		"BookName": book.BookName,
		"Barcode":  book.Barcode,
		"Author":   book.Author,
		"Category": book.Category,
	}

	c.JSON(http.StatusCreated, response)
}

// GetBooks godoc
// @Summary Get all books
// @Description Get all books from database
// @Produce json
// @Success 200 {array} models.Book
// @Router /books [get]
func GetBooks(c *gin.Context) {
	//time.Sleep(2 * time.Second)

	var books []models.Book

	initiliazers.DB.Find(&books)

	if len(books) == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"books": books,
	})
}

// GetBookById godoc
// @Summary Get a book by ID
// @Description Get a book from database by ID
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} models.Book
// @Router /books/{id} [get]
func GetBookById(c *gin.Context) {
	var book models.Book

	id := c.Param("id")

	if err := initiliazers.DB.First(&book, id).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	response := gin.H{
		"BookName": book.BookName,
		"Barcode":  book.Barcode,
		"Author":   book.Author,
		"Category": book.Category,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateBook godoc
// @Summary Update a book by ID
// @Description Update a book in database by ID
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body models.Book true "Book Update Data"
// @Success 200 {object} models.Book
// @Router /books/{id} [put]
func UpdateBook(c *gin.Context) {

	id := c.Param("id")

	var body struct {
		BookName string
		Barcode  string
		Author   string
		Category string
	}

	err := c.Bind(&body)
	if err != nil {
		return
	}

	var book models.Book
	if err := initiliazers.DB.First(&book, id).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if isBarcodeExistsWithDifferentID(body.Barcode, id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": BarcodeAlreadyExist})
		return
	}

	if checkMissingFields(c, body, errorMessages) {
		return
	}

	initiliazers.DB.Model(&book).Updates(models.Book{BookName: body.BookName, Barcode: body.Barcode, Author: body.Author, Category: body.Category, UpdatedDate: time.Now()})

	response := gin.H{
		"BookName": book.BookName,
		"Barcode":  book.Barcode,
		"Author":   book.Author,
		"Category": book.Category,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteBook godoc
// @Summary Delete a book by ID
// @Description Delete a book from database by ID
// @Produce json
// @Param id path string true "Book ID"
// @Success 204 "No Content"
// @Router /books/{id} [delete]
func DeleteBook(c *gin.Context) {

	id := c.Param("id")

	var book models.Book
	result := initiliazers.DB.Unscoped().Delete(&book, id)

	if result.RowsAffected == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	c.Status(http.StatusNoContent)
}

func findMissingFields(data interface{}) []string {
	v := reflect.ValueOf(data)
	var missingFields []string

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldName := v.Type().Field(i).Name

		if fieldValue.Interface() == "" {
			missingFields = append(missingFields, fieldName)
		}
	}

	return missingFields
}

func checkMissingFields(c *gin.Context, body interface{}, errorMessages map[string]string) bool {
	missingFields := findMissingFields(body)
	if len(missingFields) > 0 {
		for _, field := range missingFields {
			errorMsg := errorMessages[field]
			c.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
			return true
		}
	}
	return false
}

func barcodeExists(barcode string) bool {
	existingBook := models.Book{}
	err := initiliazers.DB.Where("barcode = ?", barcode).First(&existingBook).Error
	return err == nil
}

func isBarcodeExistsWithDifferentID(barcode string, id string) bool {
	var existingBook models.Book
	if err := initiliazers.DB.Where("barcode = ? AND id != ?", barcode, id).First(&existingBook).Error; err == nil {
		return true
	}
	return false
}
