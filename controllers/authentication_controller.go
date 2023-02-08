package controllers

import (
	"net/http"
	"time"
	"what-to-eat/configuration"
	"what-to-eat/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

var authCollection *mongo.Collection = configuration.GetCollection(configuration.DB, "polls")

func CreateAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var poll models.Poll
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&poll); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&poll); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		poll.Id = primitive.NewObjectID()

		result, err := pollCollection.InsertOne(ctx, poll)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, models.Response{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetBasicAuthKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		authId := c.Param("id")
		authPassword := c.Param("password")
		var token models.Token
		var poll models.Poll
		defer cancel()

		objid, _ := primitive.ObjectIDFromHex(resId)

		err := pollCollection.FindOne(ctx, bson.M{"id": objid}).Decode(&poll)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Poll with specified ID not found"}})
			return
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": poll}})
	}
}

func GetGoogleAuthKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		resId := c.Param("id")
		var auth models.Auth
		var poll models.Poll
		defer cancel()

		objid, _ := primitive.ObjectIDFromHex(resId)

		err := pollCollection.FindOne(ctx, bson.M{"id": objid}).Decode(&poll)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Poll with specified ID not found"}})
			return
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": poll}})
	}
}

func GetFacebookAuthKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		resId := c.Param("id")
		var auth models.Auth
		var poll models.Poll
		defer cancel()

		objid, _ := primitive.ObjectIDFromHex(resId)

		err := pollCollection.FindOne(ctx, bson.M{"id": objid}).Decode(&poll)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Poll with specified ID not found"}})
			return
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": poll}})
	}
}

func DeletePoll() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		pollId := c.Param("id")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(pollId)

		result, err := pollCollection.DeleteOne(ctx, bson.M{"id": objId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, models.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Poll with specified ID not found"}})
			return
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "Poll successfully deleted"}})

	}
}

func UpdatePoll() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		pollId := c.Param("id")
		var poll models.Poll
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(pollId)

		//validate the request body
		if err := c.BindJSON(&poll); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&poll); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{
			"id":        poll.Id,
			"detail":    poll.Detail,
			"No":        poll.ParticipantsNo,
			"createdat": poll.CreatedAt,
			"updated":   poll.UpdatedAt,
			"expiredat": poll.ExpiredAt,
		}
		result, err := pollCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		var updatedPoll models.Poll
		if result.MatchedCount == 1 {
			err := pollCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedPoll)
			if err != nil {

			}
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedPoll}})

	}
}
