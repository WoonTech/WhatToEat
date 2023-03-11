package controllers

import (
	"net/http"
	"time"
	"what-to-eat/models"
	"what-to-eat/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

var sessions = map[string]models.Session{}

func isExpired(s models.Session) bool {
	return s.ExpiredAt.Before(time.Now())
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var cred models.Credentials
		timeNow := time.Now().UTC()
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&cred); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: utils.FailSignedUp, Content: err.Error()})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&cred); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: utils.FailSignedUp, Content: validationErr.Error()})
			return
		}

		//TO-DO
		//checkpassword regex and emailvalidation
		//send verification message
		//hashpassword

		credCount, _ := credCollection.CountDocuments(ctx, bson.M{})
		credId := int(credCount) + 1
		cred.Detail.Status = utils.StatusActive
		cred.Detail.LastLoginAt = timeNow
		newCred := models.Credentials{
			Id:            credId,
			CreatedAt:     timeNow,
			UpdatedAt:     timeNow,
			Username:      cred.Username,
			Password:      cred.Password,
			ContactNumber: cred.ContactNumber,
			Email:         cred.Email,
			Detail:        cred.Detail,
		}

		result, err := credCollection.InsertOne(ctx, newCred)
		_ = result
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailSignedUp, Content: err.Error()})
			return
		}

		var insertResult models.Credentials
		if err := credCollection.FindOne(ctx, bson.M{"id": credId}).Decode(&insertResult); err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailRetrievedUser, Content: err.Error()})
			return
		}
		c.JSON(http.StatusCreated, models.Response{Status: http.StatusCreated, Message: utils.SuccessSignedUp, Content: insertResult})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var creds models.Credentials
		timeNow := time.Now().UTC()
		defer cancel()

		//To-DO
		//check if cookie session existed

		if err := c.BindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: utils.FailLogin, Content: err.Error()})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&creds); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: utils.FailLogin, Content: validationErr.Error()})
			return
		}

		var userCreds models.Credentials
		err := credCollection.FindOne(ctx, bson.M{"username": creds.Username}).Decode(&userCreds)
		if err != nil {
			c.JSON(http.StatusNotFound, models.Response{Status: http.StatusInternalServerError, Message: utils.FailLogin, Content: "Incorrect username or password"})
			return
		}

		if creds.Password != userCreds.Password {
			c.JSON(http.StatusNotFound, models.Response{Status: http.StatusInternalServerError, Message: utils.FailLogin, Content: "Incorrect username or password"})
			return
		}

		//create cookie session
		sessionToken := uuid.NewString()
		expiredSecond := 5000
		expiredAt := time.Now().Add(time.Duration(expiredSecond) * time.Second).UTC()

		//add this sessionToken to db
		opts := options.FindOne().SetSort((bson.D{{Key: "age", Value: 1}}))
		if err := sessionCollection.FindOne(ctx, bson.M{"username": creds.Username}, opts).Err(); err != nil {
			if err == mongo.ErrNoDocuments {
				sessionCount, _ := sessionCollection.CountDocuments(ctx, bson.M{})
				sessionId := int(sessionCount) + 1
				session := models.Session{
					Id:        sessionId,
					CreatedAt: timeNow,
					UpdatedAt: timeNow,
					Username:  creds.Username,
					SessionId: sessionToken,
					ExpiredAt: expiredAt,
					Status:    utils.StatusActive,
				}
				result, err := sessionCollection.InsertOne(ctx, session)
				_ = result
				if err != nil {
					c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailCreatedSessionToken, Content: err.Error()})
					return
				}
			}
		} else {
			updatedSession := bson.D{
				{Key: "$set", Value: bson.D{{Key: "updated_at", Value: timeNow},
					{Key: "session_id", Value: sessionToken},
					{Key: "expired_at", Value: expiredAt}}},
			}
			result, err := sessionCollection.UpdateOne(ctx, bson.M{"username": creds.Username}, updatedSession)
			_ = result
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailUpdatedSessionToken, Content: err.Error()})
				return
			}
		}

		//set cookie as session token
		c.SetCookie("session_token", sessionToken, expiredSecond, "/", "localhost", false, true)
		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: utils.SuccessLogin, Content: sessionToken})
	}
}

/*func Refresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		sessionToken, err := c.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				c.JSON(http.StatusUnauthorized, models.Response{Status: http.StatusInternalServerError, Message: "error", Content: map[string]interface{}{"data": "Unauthorized cookie"}})
				return
			}
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusInternalServerError, Message: "error", Content: map[string]interface{}{"data": "Fail to retrieve cookie"}})
		}

		userSession, isExisted := sessions[sessionToken]
		if !isExisted {
			c.JSON(http.StatusUnauthorized, models.Response{Status: http.StatusInternalServerError, Message: "error", Content: map[string]interface{}{"data": "Unauthorized cookie"}})
			return
		}

		if isExpired(userSession) {
			delete(sessions, sessionToken)
			c.JSON(http.StatusUnauthorized, models.Response{Status: http.StatusInternalServerError, Message: "error", Content: map[string]interface{}{"data": "User session has expired"}})
			return
		}

		newSessionToken := uuid.NewString()
		expiredAt := time.Now().Add(120 * time.Second)

		sessions[newSessionToken] = models.Session{
			Username:  userSession.Username,
			ExpiredAt: expiredAt,
		}

		delete(sessions, sessionToken)

		c.SetCookie("session_token", sessionToken, 120, "/", "localhost", false, true)
		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: "success", Content: map[string]interface{}{"data": "Refresh successfully"}})

	}
}*/

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var creds models.Credentials
		timeNow := time.Now().UTC()
		defer cancel()

		if err := c.BindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: utils.FailloggedOut, Content: err.Error()})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&creds); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: utils.FailloggedOut, Content: validationErr.Error()})
			return
		}

		opts := options.FindOne().SetSort((bson.D{{Key: "age", Value: 1}}))
		if err := sessionCollection.FindOne(ctx, bson.M{"username": creds.Username}, opts).Err(); err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusNoContent, Message: utils.FailloggedOut, Content: "Username not found"})
				return
			}
		}

		updatedSession := bson.D{
			{Key: "$set", Value: bson.D{{Key: "updated_at", Value: timeNow},
				{Key: "session_id", Value: ""},
				{Key: "expired_at", Value: ""},
				{Key: "status", Value: utils.StatusDeleted}}},
		}

		result, err := sessionCollection.UpdateOne(ctx, bson.M{"username": creds.Username}, updatedSession)
		_ = result
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailloggedOut, Content: err.Error()})
			return
		}

		c.SetCookie("session_token", "", 0, "/", "localhost", false, true)
		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: utils.SuccessLoggedOut, Content: "Logout successfully"})

	}
}
