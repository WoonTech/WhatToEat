package Controllers

import (
	"gin-mongo-api/configs"
	"net/http"
	"time"
	"what-to-eat/models"
	"what-to-eat/responses"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

var resCollection *mongo.Collection = configs.GetCollection(configs.DB, "restaurants")
var validate = validator.New()

func CreateRes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var res models.Restaurant
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&res); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
	}
}
