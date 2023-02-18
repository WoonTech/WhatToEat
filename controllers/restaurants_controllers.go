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
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "Failed creating restaurant", Content: err.Error()})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&res); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "Failed creating restaurant", Content: validationErr.Error()})
			return
		}

		err := resCollection.FindOne(ctx, bson.M{"name": res.Name}).Err()
		if err == nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "Failed creating restaurant", Content: "Duplicate restaurant name found"})
			return
		}

		resCount, _ := resCollection.CountDocuments(ctx, bson.M{})
		resId := int(resCount) + 1
		newRes := models.Restaurant{
			CreatedAt:     timeNow,
			UpdatedAt:     timeNow,
			Id:            resId,
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

		result, err := resCollection.InsertOne(ctx, bson.M{
			"id":             resId,
			"created_at":     newRes.CreatedAt,
			"updated_at":     newRes.UpdatedAt,
			"name":           newRes.Name,
			"type":           newRes.Type,
			"contact":        newRes.ContactNumber,
			"service_option": newRes.ServiceOption,
			"hours":          newRes.OpenHours,
			"website":        newRes.Website,
			"address":        newRes.Address,
			"rating":         newRes.Rating,
			"menu":           newRes.Menu,
			"status":         newRes.Status,
		})
		_ = result
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "Failed creating restaurant", Content: err.Error()})
			return
		}

		c.JSON(http.StatusCreated, models.Response{Status: http.StatusCreated, Message: "Successfully inserted restaurant", Content: newRes.Id})
	}
}

func GetRes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		resId := c.Param("id")
		var res models.Restaurant
		defer cancel()

		//convert string to integer
		id, err := strconv.Atoi(resId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "Failed retrieving restaurant", Content: "Restaurant ID should be an integer"})
			return
		}

		err1 := resCollection.FindOne(ctx, bson.M{"id": id}).Decode(&res)
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "Failed retrieving restaurant", Content: "Restaurant with specified ID not found"})
			return
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: "Successfully retrieved restaurant", Content: res})
	}
}

func GetAllRes() gin.HandlerFunc {
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

func DeleteRes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		resId := c.Param("id")
		timeNow := time.Now().UTC()
		defer cancel()

		//convert string to integer
		id, err := strconv.Atoi(resId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "Failed deleting restaurant", Content: "Restaurant ID should be an integer"})
			return
		}

		//convert the restaurant status to deleted
		update := bson.D{
			{Key: "$set", Value: bson.D{{Key: "updated_at", Value: timeNow}, {Key: "status", Value: utils.StatusDeleted}}},
		}
		result, err := resCollection.UpdateOne(ctx, bson.M{"id": id}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "Failed deleting restaurant", Content: err.Error()})
			return
		}

		var updatedRes models.Restaurant
		if result.MatchedCount == 1 {
			err := resCollection.FindOne(ctx, bson.M{"id": id}).Decode(&updatedRes)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "Failed deleting restaurant", Content: err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusNotFound, models.Response{Status: http.StatusNotFound, Message: "Failed deleting restaurant", Content: "Restaurant with specified ID not found"})
			return
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: "Successfully deleted restaurant", Content: resId})

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
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "Failed updating restaurant", Content: "Restaurant ID should be an integer"})
			return
		}

		//validate the request body
		if err := c.BindJSON(&res); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "Failed updating restaurant", Content: err.Error()})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&res); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "Failed updating restaurant", Content: validationErr.Error()})
			return
		}

		update := bson.D{
			{Key: "$set", Value: bson.D{{Key: "updated_at", Value: timeNow},
				{Key: "name", Value: res.Name},
				{Key: "type", Value: res.Type},
				{Key: "status", Value: res.Status},
				{Key: "service_option", Value: res.ServiceOption},
				{Key: "hours", Value: res.OpenHours},
				{Key: "contact", Value: res.ContactNumber},
				{Key: "website", Value: res.Website},
				{Key: "address", Value: res.Address},
				{Key: "rating", Value: res.Rating},
				{Key: "menu", Value: res.Menu}}},
		}
		result, err := resCollection.UpdateOne(ctx, bson.M{"id": id}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "Failed updating restaurant", Content: err.Error()})
			return
		}

		var updatedRes models.Restaurant
		if result.MatchedCount == 1 {
			err := resCollection.FindOne(ctx, bson.M{"id": id}).Decode(&updatedRes)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "Failed updating restaurant", Content: err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusNotFound, models.Response{Status: http.StatusNotFound, Message: "Failed updating restaurant", Content: "Restaurant with specified ID not found"})
			return
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: "Successfully updated restaurant", Content: updatedRes})

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
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "Failed creating menu", Content: err.Error()})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&menu); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "Failed creating menu", Content: validationErr.Error()})
			return
		}

		//return if same id and restaurant name
		err := menuCollection.FindOne(ctx, bson.M{"restaurant_name": menu.RestaurantName, "restaurant_id": menu.RestaurantId}).Err()
		if err == nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "Failed creating menu", Content: "Duplicate restaurant name and ID found"})
			return
		}

		menuCount, err := menuCollection.CountDocuments(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "Failed creating menu", Content: err.Error()})
			return
		}
		menuId := int(menuCount) + 1
		newMenu := bson.M{
			"id":              menuId,
			"created_at":      timeNow,
			"updated_at":      timeNow,
			"restaurant_name": menu.RestaurantName,
			"restaurant_id":   menu.RestaurantId,
			"menu_details":    menu.Menu,
		}
		menuResult, err := menuCollection.InsertOne(ctx, newMenu)

		_ = menuResult
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "Failed creating menu", Content: err.Error()})
			return
		}
		//insert the menu into the restaurant
		update := bson.M{
			"updated_at": timeNow,
			"menu":       newMenu,
		}
		resResult, err := resCollection.UpdateOne(ctx, bson.M{"id": menu.RestaurantId}, bson.M{"$set": update})
		_ = resResult
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "Failed creating menu", Content: err.Error()})
			return
		}

		c.JSON(http.StatusCreated, models.Response{Status: http.StatusCreated, Message: "Successfully created menu", Content: menuId})
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
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "Failed updating menu", Content: "Menu ID should be an integer"})
			return
		}

		//validate the request body
		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "Failed updating menu", Content: err.Error()})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&menu); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "Failed updating menu", Content: validationErr.Error()})
			return
		}

		//updateMenu menu
		updateMenu := bson.M{
			"updated_at":   timeNow,
			"menu_details": menu.Menu,
		}
		result, err := menuCollection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": updateMenu})
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "Failed updating menu", Content: err.Error()})
			return
		}

		var newMenu models.Menu
		if result.MatchedCount == 1 {
			err := menuCollection.FindOne(ctx, bson.M{"id": id}).Decode(&newMenu)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "Failed updating menu", Content: err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusNotFound, models.Response{Status: http.StatusNotFound, Message: "Failed updating menu", Content: "Menu with specified ID not found"})
			return
		}

		//insert the menu into the restaurant
		resResult, err := resCollection.UpdateOne(ctx, bson.M{"id": newMenu.RestaurantId}, bson.M{"$set": bson.D{{Key: "updated_at", Value: timeNow}, {Key: "menu", Value: newMenu}}})
		_ = resResult
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "Failed updating menu", Content: err.Error()})
			return
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: "Successfully updated menu", Content: newMenu.Id})

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
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "Failed retrieving menu", Content: "Menu ID should be an integer"})
			return
		}

		err1 := menuCollection.FindOne(ctx, bson.M{"id": id}).Decode(&menu)
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "Failed retrieving menu", Content: "Menu ID not found"})
			return
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: "Successfully retrieved menu", Content: menu})
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
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "Failed deleting menu", Content: "Restaurant ID should be an integer"})
			return
		}

		result, err := menuCollection.DeleteOne(ctx, bson.M{"id": id})

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "Failed deleting menu", Content: err.Error()})
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, models.Response{Status: http.StatusNotFound, Message: "Failed deleting menu", Content: "Menu with specified ID not found"})
			return
		}

		//delete menu from restaurant
		//resResult, err := resCollection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.D{Key: "updatedat", Value: timeNow},"$unset" :"menu"}})
		var emptyMenu models.Menu
		update := bson.D{
			{Key: "$set", Value: bson.D{{Key: "updated_at", Value: timeNow}}},
			{Key: "$unset", Value: bson.D{{Key: "menu", Value: emptyMenu}}},
		}
		resResult, err := resCollection.UpdateOne(ctx, bson.M{"id": id}, update)

		_ = resResult
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "Failed deleting menu", Content: err.Error()})
			return
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: "Successfully deleted menu", Content: "Menu successfully deleted"})

	}
}
