package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"lessons/lib_db"
	"lessons/models"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var database *sqlx.DB

type ResponseJson struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func main() {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "database.db"
	}
	dbPath := fmt.Sprintf("./%s", dbName)
	database, err := lib_db.DBInit(dbPath)
	if err != nil {
		log.Println(err)
	}

	router := gin.Default()
	router.POST("/api/v1/products", writeDB(database))
	router.GET("/api/v1/products", listProducts(database))
	router.DELETE("/api/v1/products/:id", deleteProduct(database))
	router.Run(":" + httpPort)
	defer database.Close()
}

func listProducts(db *sqlx.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		pp, err := lib_db.ListProducts(db)
		if err != nil {
			log.Println(err)
		}
		apiOkResponseList(db, &pp, c)
	}
}

func deleteProduct(db *sqlx.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		count, err := lib_db.IfExistsTitleFromDB(db, id)
		if err != nil {
			log.Println(err)
		}
		if count == 1 {
			err := lib_db.DeleteProductDB(db, id)
			if err != nil {
				log.Println(err)
			}
			apiErrorResponse(205, "JSON already deleted", c)
		} else {
			apiErrorResponse(408, "JSON Title not in DB", c)
		}
	}
}

func writeDB(db *sqlx.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var jsonRequest models.Product
		err := c.ShouldBindJSON(&jsonRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		titleNoNil := jsonRequest.Title != nil
		tagsLen4More := len(jsonRequest.Tags) >= 4
		titleLenLess50 := len([]rune(*jsonRequest.Title)) < 50

		switch {
		case !titleNoNil && tagsLen4More:
			apiErrorResponse(401, "JSON must have Title", c)
		case titleNoNil && !tagsLen4More:
			apiErrorResponse(402, "JSON must have 4 or more Tags", c)
		case !titleNoNil && !tagsLen4More:
			apiErrorResponse(403, "JSON must have Title and 4 or more Tags", c)
		case !titleLenLess50 && tagsLen4More:
			apiErrorResponse(404, "JSON must have Title 50 or less symbols", c)
		case !tagsLen4More && titleNoNil && titleLenLess50:
			apiErrorResponse(405, "JSON must have 4 or more Tags", c)
		case !tagsLen4More && titleLenLess50:
			apiErrorResponse(406, "JSON must have Title 50 or less symbols and 4 or more Tags", c)
		default:
			apiOkResponse(db, &jsonRequest, c)
		}
	}
}


func apiErrorResponse(code int, message string, c *gin.Context) {
	jsonError := ResponseJson{
		Code:    code,
		Message: message,
	}
	c.JSON(code, jsonError)
	return
}

func apiOkResponse(db *sqlx.DB, jsonRequest *models.Product, cont *gin.Context) {
	if lib_db.CheckTitleExists(db, jsonRequest) {
		apiErrorResponse(407, "JSON Title already exists", cont)
	} else {
		lib_db.InsertProductDB(db, jsonRequest)
		apiErrorResponse(200, "JSON saved to DB", cont)
	}
}

func apiOkResponseList(db *sqlx.DB, jsonRequest *[]models.Product, cont *gin.Context) {
	apiErrorResponse(203, "List of Products", cont)
	cont.JSON(203, jsonRequest)
}
