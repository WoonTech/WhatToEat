package controllers

import (
	"net/http"
	"time"
	"what-to-eat/models"
	"what-to-eat/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
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
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailSignedUp, Content: err.Error()})
			return
		}

		c.JSON(http.StatusCreated, models.Response{Status: http.StatusCreated, Message: utils.SuccessSignedUp, Content: result})
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

		var expectedCreds models.Credentials
		err := credCollection.FindOne(ctx, bson.M{"username": creds.Username}).Decode(&expectedCreds)
		if err != nil {
			c.JSON(http.StatusNotFound, models.Response{Status: http.StatusInternalServerError, Message: utils.FailLogin, Content: "Incorrect username or password"})
			return
		}

		if creds.Password != expectedCreds.Password {
			c.JSON(http.StatusNotFound, models.Response{Status: http.StatusInternalServerError, Message: utils.FailLogin, Content: "Incorrect username or password"})
			return
		}

		//create cookie session
		sessionToken := uuid.NewString()
		expiredAt := time.Now().Add(120 * time.Second).UTC()
		session := models.Session{
			Id:        1,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			Username:  creds.Username,
			SessionId: sessionToken,
			ExpiredAt: expiredAt,
		}
		result, err := sessionCollection.InsertOne(ctx, session)
		_ = result
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: utils.FailCreatedSessionToken, Content: err.Error()})
			return
		}
		//add this sessionToken to db

		//set cookie as session token
		c.SetCookie("session_token", sessionToken, 120, "/", "localhost", false, true)
		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: utils.SuccessLogin, Content: "Successfully login"})
	}
}

func Refresh() gin.HandlerFunc {
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
}

func Logout() gin.HandlerFunc {
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

		delete(sessions, sessionToken)

		c.SetCookie("session_token", "", 0, "/", "localhost", false, true)
		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: "success", Content: map[string]interface{}{"data": "Logout successfully"}})

	}
}
