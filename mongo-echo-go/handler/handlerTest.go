package handler

import (
	"context"
	"net/http"
	"os"

	"mongo-echo-go/modal"
	"mongo-echo-go/utils"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func GetTest(c echo.Context) error {
	// Get MongoDB collection
	collection := utils.GetDB().Database(utils.GetDBName()).Collection(utils.GetCollectionName())

	// Find document by name
	var doc modal.Test
	err := collection.FindOne(context.TODO(), bson.M{"name": "sample"}).Decode(&doc)
	if err != nil {
		c.Logger().SetOutput(os.Stdout)
		c.Logger().Errorf("Failed to find document: %v", err)
		return c.JSON(http.StatusInternalServerError, modal.Response{
			Code:    500,
			Message: "Failed to fetch data",
			Data:    nil,
		})
	}

	c.Logger().SetOutput(os.Stdout)
	c.Logger().Infof("Found document: %+v", doc)
	return c.JSON(http.StatusOK, modal.Response{
		Code:    200,
		Message: "",
		Data:    doc,
	})
}
