package controllers

import (
	"net/http"
	"time"
	"what-to-eat/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&cred); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "error", Content: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&cred); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "error", Content: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		//checkpassword regex and emailvalidation
		//send verification message
		//hashpassword

		cred.Id = primitive.NewObjectID()

		result, err := credCollection.InsertOne(ctx, cred)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Content: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, models.Response{Status: http.StatusCreated, Message: "success", Content: map[string]interface{}{"data": result}})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var creds models.Credentials
		defer cancel()

		//check if cookie session existed

		if err := c.BindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "error", Content: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&creds); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "error", Content: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		var expectedCreds models.Credentials
		err := credCollection.FindOne(ctx, bson.M{"username": creds.Username}).Decode(&expectedCreds)
		if err != nil {
			c.JSON(http.StatusNotFound, models.Response{Status: http.StatusInternalServerError, Message: "error", Content: map[string]interface{}{"data": "The username or password you entered is incorrect"}})
			return
		}

		if creds.Password != expectedCreds.Password {
			c.JSON(http.StatusNotFound, models.Response{Status: http.StatusInternalServerError, Message: "error", Content: map[string]interface{}{"data": "The username or password you entered is incorrect"}})
			return
		}

		//create cookie session
		sessionToken := uuid.NewString()
		expiredAt := time.Now().Add(120 * time.Second)
		sessions[sessionToken] = models.Session{
			Username:  creds.Username,
			ExpiredAt: expiredAt,
		}

		//set cookie as session token
		c.SetCookie("session_token", sessionToken, 120, "/", "localhost", false, true)
		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: "success", Content: map[string]interface{}{"data": "Login successfully"}})
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
