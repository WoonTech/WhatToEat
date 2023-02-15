package controllers

import (
	"net/http"
	"strconv"
	"time"
	"what-to-eat/models"
	"what-to-eat/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
)

func CreateRes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var res models.Restaurant
		timeNow := time.Now().UTC()
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&res); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&res); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		err := resCollection.FindOne(ctx, bson.M{"name": res.Name}).Err()
		if err == nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Duplicate restaurant name"}})
			return
		}

		resCount, _ := resCollection.CountDocuments(ctx, bson.M{})
		newRes := models.Restaurant{
			CreatedAt:     timeNow,
			UpdatedAt:     timeNow,
			Id:            int(resCount) + 1,
			Name:          res.Name,
			Type:          res.Type,
			ContactNumber: res.ContactNumber,
			ServiceOption: res.ServiceOption,
			OpenHours:     res.OpenHours,
			Website:       res.Website,
			Address:       res.Address,
			Rating:        res.Rating,
			Menu:          res.Menu,
			Status:        utils.StatusActive,
		}

		result, err := resCollection.InsertOne(ctx, newRes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, models.Response{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": "Successfully inserted restaurant", "_id": result, "id": newRes.Id}})
	}
}

func GetRes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		resId := c.Param("id")
		var res models.Restaurant
		defer cancel()

		if err := c.BindJSON(&res); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//convert string to integer
		id, err := strconv.Atoi(resId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Restaurant ID should be a integer"}})
			return
		}

		err1 := resCollection.FindOne(ctx, bson.M{"id": id}).Decode(&res)
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Restaurant with specified ID not found"}})
			return
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": res}})
	}
}

func GetAllRes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var ress []models.Restaurant
		defer cancel()

		results, err := resCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var res models.Restaurant
			if err = results.Decode(&res); err != nil {
				c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			ress = append(ress, res)
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": ress}})

	}
}

func DeleteRes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		resId := c.Param("id")
		timeNow := time.Now().UTC()
		defer cancel()

		//convert string to integer
		id, err := strconv.Atoi(resId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Restaurant ID should be a integer"}})
			return
		}

		//convert the restaurant status to deleted
		result, err := resCollection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.D{{"updated_at", timeNow}, {"status", utils.StatusDeleted}}})
		var updatedRes models.Restaurant
		if result.MatchedCount == 1 {
			err := resCollection.FindOne(ctx, bson.M{"id": id}).Decode(&updatedRes)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		} else {
			c.JSON(http.StatusNotFound, models.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Restaurant with specified ID not found"}})
			return
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "Restaurant successfully deleted"}})

	}
}

func UpdateRes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		resId := c.Param("id")
		var res models.Restaurant
		timeNow := time.Now().UTC()
		defer cancel()

		//convert string to integer
		id, err := strconv.Atoi(resId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Restaurant ID should be a integer"}})
			return
		}

		//validate the request body
		if err := c.BindJSON(&res); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&res); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{
			"updated_at":    timeNow,
			"id":            id,
			"name":          res.Name,
			"type":          res.Type,
			"contact":       res.ContactNumber,
			"serviceoption": res.ServiceOption,
			"hours":         res.OpenHours,
			"website":       res.Website,
			"address":       res.Address,
			"status":        res.Status,
			"rating":        res.Rating,
			"menu":          res.Menu,
		}
		result, err := resCollection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": update})

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		var updatedRes models.Restaurant
		if result.MatchedCount == 1 {
			err := resCollection.FindOne(ctx, bson.M{"id": id}).Decode(&updatedRes)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		} else {
			c.JSON(http.StatusNotFound, models.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Restaurant with specified ID not found"}})
			return
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedRes}})

	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var menu models.Menu
		timeNow := time.Now().UTC()
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&menu); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		menuCount, err := menuCollection.CountDocuments(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		newMenu := models.Menu{
			Id:             int(menuCount) + 1,
			RestaurantName: menu.RestaurantName,
			RestaurantId:   menu.RestaurantId,
			CreatedAt:      timeNow,
			UpdatedAt:      timeNow,
			Menu:           menu.Menu,
		}

		menuResult, err := resCollection.InsertOne(ctx, newMenu)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//insert the menu into the restaurant
		resResult, err := resCollection.UpdateOne(ctx, bson.M{"id": newMenu.RestaurantId}, bson.M{"$set": bson.D{{"updated_at", timeNow}, {"menu", newMenu}}})
		_ = resResult
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, models.Response{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": menuResult, "name": newMenu.Id}})
	}
}

func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		menuId := c.Param("id")
		timeNow := time.Now().UTC()
		var menu models.Menu
		defer cancel()

		//convert string to integer
		id, err := strconv.Atoi(menuId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Restaurant ID should be a integer"}})
			return
		}

		//validate the request body
		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&menu); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		//update menu
		update := bson.M{
			"id":           id,
			"updated_at":   timeNow,
			"menu_details": menu.Menu,
		}
		result, err := resCollection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		var newMenu models.Menu
		if result.MatchedCount == 1 {
			err := resCollection.FindOne(ctx, bson.M{"id": id}).Decode(&newMenu)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		} else {
			c.JSON(http.StatusNotFound, models.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Menu with specified ID not found"}})
			return
		}

		//insert the menu into the restaurant
		resResult, err := resCollection.UpdateOne(ctx, bson.M{"id": newMenu.RestaurantId}, bson.M{"$set": bson.D{{"updated_at", timeNow}, {"menu", newMenu}}})
		_ = resResult
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": newMenu}})

	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		menuId := c.Param("id")
		var menu models.Menu
		defer cancel()

		id, err := strconv.Atoi(menuId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Restaurant ID should be a integer"}})
			return
		}

		err1 := resCollection.FindOne(ctx, bson.M{"id": id}).Decode(&menu)
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Menu with specified ID not found"}})
			return
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": menu}})
	}
}

// straight away delete the menu, dont need to convert it to delete
func DeleteMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		menuId := c.Param("id")
		timeNow := time.Now().UTC()
		defer cancel()

		id, err := strconv.Atoi(menuId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Restaurant ID should be a integer"}})
			return
		}

		result, err := resCollection.DeleteOne(ctx, bson.M{"id": id})

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, models.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Menu with specified ID not found"}})
			return
		}

		//delete menu from restaurant
		var emptyMenu models.Menu
		resResult, err := resCollection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.D{{"updated_at", timeNow}, {"menu", emptyMenu}}})
		_ = resResult
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "Menu successfully deleted"}})

	}
}
