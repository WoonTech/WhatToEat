package controllers

import (
	"net/http"
	"strconv"
	"time"
	"what-to-eat/models"
	"what-to-eat/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

func CreatePoll() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var poll models.Poll
		timeNow := time.Now().UTC()
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&poll); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: utils.FailCreatedPoll, Content: err.Error()})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&poll); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: utils.FailCreatedPoll, Content: validationErr.Error()})
			return
		}

		//get last record instead of pollcount, because poll can be deleted
		userCount, _ := pollCollection.CountDocuments(ctx, bson.M{})
		var pollId int
		if int(userCount) == 0 {
			pollId = int(userCount) + 1
		} else {
			var lastrecord models.Poll
			opts := options.FindOne().SetSort(bson.M{"$natural": -1})
			if err := pollCollection.FindOne(ctx, bson.M{}, opts).Decode(&lastrecord); err != nil {
				c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: utils.FailCreatedPoll, Content: err.Error()})
				return
			}
			pollId = lastrecord.Id + 1
		}

		pollInserted := models.Poll{
			Id:             pollId,
			CreatedAt:      timeNow,
			UpdatedAt:      timeNow,
			ExpiredAt:      poll.ExpiredAt,
			Detail:         poll.Detail,
			ParticipantsNo: poll.ParticipantsNo,
		}
		result, err := pollCollection.InsertOne(ctx, pollInserted)
		_ = result
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailCreatedPoll, Content: err.Error()})
			return
		}

		c.JSON(http.StatusCreated, models.Response{Status: http.StatusCreated, Message: utils.SuccessCreatedPoll, Content: pollInserted})
	}
}

func GetPoll() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		pollId := c.Param("id")
		var poll models.Poll
		defer cancel()

		//convert string to integer
		id, err := strconv.Atoi(pollId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailRetrievedPoll, Content: "Poll ID should be an integer"})
			return
		}

		if err := pollCollection.FindOne(ctx, bson.M{"id": id}).Decode(&poll); err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailRetrievedPoll, Content: "Poll with specified ID not found"})
			return
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: utils.SuccessCreatedPoll, Content: poll})
	}
}

func DeletePoll() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		pollId := c.Param("id")
		defer cancel()

		//convert string to integer
		id, err := strconv.Atoi(pollId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailDeletedPoll, Content: "Poll ID should be an integer"})
			return
		}

		result, err := pollCollection.DeleteOne(ctx, bson.M{"id": id})

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailDeletedPoll, Content: err.Error()})
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, models.Response{Status: http.StatusNotFound, Message: utils.FailDeletedPoll, Content: "Poll with specified ID not found"})
			return
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: utils.SuccessDeletedPoll, Content: "Poll successfully deleted"})

	}
}

func GetAllPoll() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var polls []models.Poll
		defer cancel()

		results, err := pollCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailRetrievedPoll, Content: err.Error()})
			return
		}
		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var poll models.Poll
			if err = results.Decode(&poll); err != nil {
				c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailRetrievedPoll, Content: err.Error()})
			}

			polls = append(polls, poll)
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: utils.SuccessRetrievedPoll, Content: polls})

	}
}

func UpdatePoll() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		pollId := c.Param("id")
		timeNow := time.Now().UTC()
		var poll models.Poll
		defer cancel()

		//convert string to integer
		id, err := strconv.Atoi(pollId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailUpdatedPoll, Content: "Poll ID should be an integer"})
			return
		}

		//validate the request body
		if err := c.BindJSON(&poll); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: utils.FailUpdatedPoll, Content: err.Error()})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&poll); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: utils.FailUpdatedPoll, Content: validationErr.Error()})
			return
		}

		update := bson.D{
			{Key: "$set", Value: bson.D{{Key: "updated_at", Value: timeNow},
				{Key: "poll_details", Value: poll.Detail},
				{Key: "participants", Value: poll.ParticipantsNo}}},
		}
		result, err := pollCollection.UpdateOne(ctx, bson.M{"id": id}, update)

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailUpdatedPoll, Content: err.Error()})
			return
		}

		var updatedPoll models.Poll
		if result.MatchedCount == 1 {
			err := pollCollection.FindOne(ctx, bson.M{"id": id}).Decode(&updatedPoll)
			if err != nil {

			}
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: utils.SuccessUpdatedPoll, Content: updatedPoll})

	}
}
