package controllers

import (
	"net/http"
	"time"
	"what-to-eat/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
)

func GetGoogleRes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var ress []models.Restaurant
		defer cancel()

		results, err := resCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "Failed retrieving restaurant", Content: err.Error()})
			return
		}
		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var res models.Restaurant
			if err = results.Decode(&res); err != nil {
				c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "Failed retrieving restaurant", Content: err.Error()})
				return
			}

			ress = append(ress, res)
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: "Successfully retrieved restaurant", Content: ress})

	}
}
