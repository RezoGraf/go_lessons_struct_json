package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Tag struct {
	Id  int
	Tag string
}

type Object struct {
	Title   string
	Comment string
}

type RequestJson struct {
	//Убрал binding что бы обрабатывать и выводить ошибки согласно заданию
	// Title          string `json:"title" binding:"required"`
	Title *string `json:"title"`
	// Tags           []Tag  `json:"tags" binding:"required"`
	Tags           []Tag  `json:"tags"`
	Description    string `json:"description,omitempty"`
	Price          string `json:"price,omitempty"`
	Additionalinfo Object `json:"additionalinfo,omitempty"`
}

type ResponseJson struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func main() {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}
	router := gin.Default()
	router.POST("/api/v1/products", jsonRequestFunc)
	router.Run(":" + httpPort)
}

func jsonRequestFunc(c *gin.Context) {
	var jsonRequest RequestJson
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
		apiOkResponse(&jsonRequest, c)
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

func apiOkResponse(jsonRequest *RequestJson, cont *gin.Context) {
	var RequestSave []RequestJson
	RequestSave = append(RequestSave, *jsonRequest)
	byteBuffer, err := ioutil.ReadFile("db.json")
	if err != nil {
		log.Println(err)
		toFileSave, err := json.MarshalIndent(RequestSave, "", " ")
		if err != nil {
			log.Println(err)
		}
		err = ioutil.WriteFile("db.json", toFileSave, 0644)
		if err != nil {
			log.Println(err)
		}
		apiErrorResponse(200, "File created, JSON saved", cont)
		return
	}
	var byteRead []RequestJson
	err = json.Unmarshal(byteBuffer, &byteRead)
	if err != nil {
		log.Println(err)
	}
	byteRead = append(byteRead, *jsonRequest)
	toFileSave, err := json.MarshalIndent(byteRead, "", " ")
	if err != nil {
		log.Println(err)
	}
	err = ioutil.WriteFile("db.json", toFileSave, 0644)
	if err != nil {
		log.Println(err)
	}
	jsonError := ResponseJson{
		Code:    202,
		Message: "JSON appended to file",
	}
	cont.JSON(202, jsonError)
}
