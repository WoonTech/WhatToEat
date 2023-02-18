package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"
	"what-to-eat/models"
	"what-to-eat/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		timeNow := time.Now().UTC()
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: utils.FailCreatedUser, Content: err.Error()})
			return
		}

		if validationErr := validate.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: utils.FailCreatedUser, Content: validationErr.Error()})
			return
		}

		if err := userCollection.FindOne(ctx, bson.M{"name": user.Name}).Err(); err == nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: utils.FailCreatedUser, Content: "Duplicate user name found"})
			return
		}

		userCount, _ := userCollection.CountDocuments(ctx, bson.M{})
		userId := int(userCount) + 1
		newUser := models.User{
			Id:               userId,
			GroupId:          user.GroupId,
			Type:             user.Type,
			CreatedAt:        timeNow,
			UpdatedAt:        timeNow,
			OwnRestaurant:    user.OwnRestaurant,
			ContactNumber:    user.ContactNumber,
			Name:             user.Name,
			CurrentLocation:  user.CurrentLocation,
			Email:            user.Email,
			ChatLogs:         user.ChatLogs,
			PinnedRestaurant: user.PinnedRestaurant,
			Status:           utils.StatusActive,
		}

		result, err := userCollection.InsertOne(ctx, newUser)
		_ = result
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailCreatedUser, Content: err.Error()})
			return
		}

		c.JSON(http.StatusCreated, models.Response{Status: http.StatusCreated, Message: utils.SuccessCreatedUser, Content: userId})
	}
}

func GetAllUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var users []models.User
		defer cancel()

		results, err := userCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailRetrievedUser, Content: err.Error()})
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var user models.User
			if err = results.Decode(&user); err != nil {
				c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailRetrievedUser, Content: err.Error()})
			}

			users = append(users, user)
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: utils.SuccessRetrievedUser, Content: users})

	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("id")
		var user models.User
		defer cancel()

		//convert string to integer
		id, err := strconv.Atoi(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailRetrievedUser, Content: "User ID should be an integer"})
			return
		}

		if err := userCollection.FindOne(ctx, bson.M{"id": id}).Decode(&user); err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailRetrievedUser, Content: "User with specified ID not found"})
			return
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: utils.SuccessRetrievedUser, Content: user})

	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("id")
		timeNow := time.Now().UTC()
		defer cancel()

		//convert string to integer
		id, err := strconv.Atoi(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailDeletedUser, Content: "User ID should be an integer"})
			return
		}

		//convert the restaurant status to deleted
		update := bson.D{
			{Key: "$set", Value: bson.D{{Key: "updated_at", Value: timeNow}, {Key: "status", Value: utils.StatusDeleted}}},
		}
		result, err := userCollection.UpdateOne(ctx, bson.M{"id": id}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailDeletedUser, Content: err.Error()})
			return
		}
		var updatedUser models.User
		if result.MatchedCount == 1 {
			err := userCollection.FindOne(ctx, bson.M{"id": id}).Decode(&updatedUser)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailDeletedUser, Content: err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusNotFound, models.Response{Status: http.StatusNotFound, Message: utils.FailDeletedUser, Content: "User with specified ID not found"})
			return
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: utils.SuccessDeletedUser, Content: id})

	}
}

func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10^time.Second)
		userId := c.Param("id")
		timeNow := time.Now().UTC()
		var user models.User
		defer cancel()

		//convert string to integer
		id, err := strconv.Atoi(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailUpdatedUser, Content: "User ID should be an integer"})
			return
		}

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: utils.FailUpdatedUser, Content: err.Error()})
			return
		}

		if validationErr := validate.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: utils.FailUpdatedUser, Content: validationErr.Error()})
			return
		}
		update := bson.D{
			{Key: "$set", Value: bson.D{{Key: "updated_at", Value: timeNow},
				{Key: "name", Value: user.Name},
				{Key: "type", Value: user.Type},
				{Key: "status", Value: user.Status},
				{Key: "group_id", Value: user.GroupId},
				{Key: "restaurant", Value: user.OwnRestaurant},
				{Key: "contact", Value: user.ContactNumber},
				{Key: "location", Value: user.CurrentLocation},
				{Key: "email", Value: user.Email},
				{Key: "chat_logs", Value: user.ChatLogs},
				{Key: "pinned_restaurant", Value: user.PinnedRestaurant}}},
		}
		result, err := userCollection.UpdateOne(ctx, bson.M{"id": id}, update)

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailUpdatedUser, Content: err.Error()})
			return
		}

		var updatedUser models.User
		if result.MatchedCount == 1 {
			err := userCollection.FindOne(ctx, bson.M{"id": id}).Decode(&updatedUser)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailUpdatedUser, Content: err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: utils.SuccessUpdatedUser, Content: updatedUser})

	}
}
