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

	switch {
	case jsonRequest.Title == nil && len(jsonRequest.Tags) >= 4:
		jsonError := ResponseJson{
			Code:    401,
			Message: "JSON must have Title",
		}
		c.JSON(401, jsonError)
		return
	case jsonRequest.Title != nil && len(jsonRequest.Tags) < 4:
		jsonError := ResponseJson{
			Code:    402,
			Message: "JSON must have 4 or more Tags",
		}
		c.JSON(402, jsonError)
		return
	case jsonRequest.Title == nil && len(jsonRequest.Tags) < 4:
		jsonError := ResponseJson{
			Code:    403,
			Message: "JSON must have Title and 4 or more Tags",
		}
		c.JSON(403, jsonError)
		return
	case len(*jsonRequest.Title) > 50 && len(jsonRequest.Tags) >= 4:
		jsonError := ResponseJson{
			Code:    404,
			Message: "JSON must have Title 50 or less symbols",
		}
		c.JSON(404, jsonError)
		return
	case len(jsonRequest.Tags) < 4 && jsonRequest.Title != nil && len(*jsonRequest.Title) < 50:
		jsonError := ResponseJson{
			Code:    405,
			Message: "JSON  must have 4 or more Tags",
		}
		c.JSON(405, jsonError)
		return
	case len(jsonRequest.Tags) < 4 && len(*jsonRequest.Title) > 50:
		jsonError := ResponseJson{
			Code:    406,
			Message: "JSON must have Title 50 or less symbols and 4 or more Tags",
		}
		c.JSON(406, jsonError)
		return
	default:
		var RequestSave []RequestJson
		RequestSave = append(RequestSave, jsonRequest)
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
			jsonOK := ResponseJson{
				Code:    200,
				Message: "File created, JSON saved",
			}
			c.JSON(http.StatusOK, jsonOK)
			return
		}
		var byteRead []RequestJson
		err = json.Unmarshal(byteBuffer, &byteRead)
		if err != nil {
			log.Println(err)
		}
		byteRead = append(byteRead, jsonRequest)
		toFileSave, err := json.MarshalIndent(byteRead, "", " ")
		if err != nil {
			log.Println(err)
		}
		err = ioutil.WriteFile("db.json", toFileSave, 0644)
		if err != nil {
			log.Println(err)
		}
		jsonOK := ResponseJson{
			Code:    201,
			Message: "JSON appended to file",
		}
		c.JSON(http.StatusOK, jsonOK)
	}
}
