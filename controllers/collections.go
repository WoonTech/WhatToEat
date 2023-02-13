package controllers

import (
	"what-to-eat/configuration"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

var resCollection *mongo.Collection = configuration.GetCollection(configuration.DB, "restaurants")
var menuCollection *mongo.Collection = configuration.GetCollection(configuration.DB, "menus")
var credCollection *mongo.Collection = configuration.GetCollection(configuration.DB, "credential")
var pollCollection *mongo.Collection = configuration.GetCollection(configuration.DB, "authentication")
var userCollection *mongo.Collection = configuration.GetCollection(configuration.DB, "users")
var counterCollection *mongo.Collection = configuration.GetCollection(configuration.DB, "counters")
var validate = validator.New()
