package controllers

const BarcodeAlreadyExist = "Barcode already exist"

var errorMessages = map[string]string{
	"BookName": "BookName can not be empty",
	"Barcode":  "Barcode can not be empty",
	"Author":   "Author can not be empty",
	"Category": "Category can not be empty",
}
