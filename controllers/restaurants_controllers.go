package controllers

import (
	"net/http"
	"time"
	"what-to-eat/configuration"
	ctxResponse "what-to-eat/models/response"
	ctxRestaurant "what-to-eat/models/restaurant"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

var resCollection *mongo.Collection = configuration.GetCollection(configuration.DB, "restaurants")
var validate = validator.New()

func CreateRes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var res ctxRestaurant.Restaurant
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&res); err != nil {
			c.JSON(http.StatusBadRequest, ctxResponse.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&res); validationErr != nil {
			c.JSON(http.StatusBadRequest, ctxResponse.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newRes := ctxRestaurant.Restaurant{
			Id:   primitive.NewObjectID(),
			Name: res.Name,
		}

		result, err := resCollection.InsertOne(ctx, newRes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ctxResponse.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, ctxResponse.Response{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetRes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		resId := c.Param("id")
		var res ctxRestaurant.Restaurant
		defer cancel()

		objid, _ := primitive.ObjectIDFromHex(resId)

		err := resCollection.FindOne(ctx, bson.M{"id": objid}).Decode(&res)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ctxResponse.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Restaurant with specified ID not found"}})
			return
		}

		c.JSON(http.StatusOK, ctxResponse.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": res}})
	}
}

func GetAllRes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var ress []ctxRestaurant.Restaurant
		defer cancel()

		results, err := resCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, ctxResponse.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var res ctxRestaurant.Restaurant
			if err = results.Decode(&res); err != nil {
				c.JSON(http.StatusInternalServerError, ctxResponse.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			ress = append(ress, res)
		}

		c.JSON(http.StatusOK, ctxResponse.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": ress}})

	}
}

func DeleteRes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		resId := c.Param("id")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(resId)

		result, err := resCollection.DeleteOne(ctx, bson.M{"id": objId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, ctxResponse.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, ctxResponse.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Restaurant with specified ID not found"}})
			return
		}

		c.JSON(http.StatusOK, ctxResponse.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "Restaurant successfully deleted"}})

	}
}

func UpdateRes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		resId := c.Param("id")
		var res ctxRestaurant.Restaurant
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(resId)

		//validate the request body
		if err := c.BindJSON(&res); err != nil {
			c.JSON(http.StatusBadRequest, ctxResponse.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&res); validationErr != nil {
			c.JSON(http.StatusBadRequest, ctxResponse.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		//update := bson.M{"id": res.}//need to update whole json body
		update := bson.M{
			"name" = 
			Name               string             `json:"name,omitempty" validate:"requried"`
			Type               string             `json:"type,omitempty" validate:"requried"`
			ContactNumber      string             `json:"contact,omitempty"`
			ServiceOptionEntry uint8              `json:"serviceoption,omitempty"`
			OpenHours          string             `json:"hours,omitempty"`
			Website            string             `json:"website,omitempty"`
			Address            string             `json:"address,omitempty"`
			CommentEntry       uint8              `json:"comment,omitempty"`
			RatingEntry        uint8              `json:"rating,omitempty"`
			ItemsEntry 
		}
		result, err := resCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

		if err != nil {
			c.JSON(http.StatusInternalServerError, ctxResponse.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		var updatedRes ctxRestaurant.Restaurant
		if result.MatchedCount == 1 {
			err := resCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedRes)
			if err != nil {
				c.JSON(http.StatusInternalServerError, ctxResponse.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, ctxResponse.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedRes}})

	}
}
