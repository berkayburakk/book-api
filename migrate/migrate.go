package main

import (
	"book-api/initiliazers"
	"book-api/models"
)

func init() {

	initiliazers.LoadEnvVariables()
	initiliazers.ConnectToDB()
}

func main() {
	initiliazers.DB.AutoMigrate(&models.Book{})

}
