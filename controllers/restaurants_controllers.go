package controllers

import (
	"net/http"
	"time"
	collections "what-to-eat/controllers/collections"
	"what-to-eat/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
)

func CreateRes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var res models.Restaurant
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

		duplicateCount, _ := collections.counterCollection.FindOne(ctx, bson.M{}).resCollection.CountDocuments(ctx, bson.M{"name": res.Name})
		if duplicateCount > 0 {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Duplicate name found in database"}})
			return
		}

		resCount, _ := collections.ResCollection().CountDocuments(ctx, bson.M{})
		newRes := models.Restaurant{
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
		}

		result, err := resCollection.InsertOne(ctx, newRes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, models.Response{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result, "name": res.Name}})
	}
}

func GetRes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		resId := c.Param("id")
		var res models.Restaurant
		defer cancel()

		err := resCollection.FindOne(ctx, bson.M{"id": resId}).Decode(&res)
		if err != nil {
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
		defer cancel()

		result, err := resCollection.DeleteOne(ctx, bson.M{"id": resId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		}

		if result.DeletedCount < 1 {
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

		update := bson.M{
			"id":            resId,
			"name":          res.Name,
			"type":          res.Type,
			"contact":       res.ContactNumber,
			"serviceoption": res.ServiceOption,
			"hours":         res.OpenHours,
			"website":       res.Website,
			"address":       res.Address,
			"rating":        res.Rating,
			"menu":          res.Menu,
		}
		result, err := resCollection.UpdateOne(ctx, bson.M{"id": resId}, bson.M{"$set": update})

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		var updatedRes models.Restaurant
		if result.MatchedCount == 1 {
			err := resCollection.FindOne(ctx, bson.M{"id": resId}).Decode(&updatedRes)
			if err != nil {

			}
		}

		c.JSON(http.StatusOK, models.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedRes}})

	}
}

// func create menu return detail then use it in create res

/*func CreateOrUpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var menu models.MenuHeader
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
		/*duplicateCount, _ := resCollection.CountDocuments(ctx, bson.M{"name": res.Name})
		if duplicateCount > 0 {
			c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Duplicate name found in database"}})
			return
		}

		newRes := models.MenuHeader{
			Id:        int(menuCount) + 1,
			CreatedAt: menu.CreatedAt,
			UpdatedAt: menu.UpdatedAt,
			Menu:      menu.Menu,
		}

		result, err := resCollection.InsertOne(ctx, newRes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, models.Response{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result, "name": menu.Name}})
	}
}*/
/*func CreateMenu(menu []interface{}) models.MenuDetails {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	menuCount, err := menuCollection.CountDocuments(ctx, bson.M{})
	if duplicateCount > 0 {
		c.JSON(http.StatusBadRequest, models.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Duplicate name found in database"}})
		return
	}

	result, err := resCollection.InsertMany(ctx, menu)
	if err != nil {
		panic(err)
	}

}*/
