package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"mongo-echo-go/modal"
	"mongo-echo-go/utils"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func GetTestSum(c echo.Context) error {
	collection := utils.GetDB().Database(utils.GetDBName()).Collection(utils.GetCollectionName())

	// Get and validate type parameter
	timeType := c.QueryParam("type")
	if timeType == "" {
		timeType = "week"
	}
	if timeType != "day" && timeType != "week" && timeType != "month" {
		return c.JSON(http.StatusBadRequest, modal.Response{
			Code:    400,
			Message: "Invalid type parameter, must be day/week/month",
			Data:    nil,
		})
	}

	// Set date format based on type
	dateFormat := "%Y-%m-%d"
	if timeType == "week" {
		dateFormat = "%Y-%U"
	} else if timeType == "month" {
		dateFormat = "%Y-%m"
	}

	// Parse time range parameter
	timeRange := c.QueryParam("time")
	var startTime, endTime time.Time
	if timeRange != "" {
		parts := strings.Split(strings.Trim(timeRange, "[]"), ",")
		if len(parts) != 2 {
			return c.JSON(http.StatusBadRequest, modal.Response{
				Code:    400,
				Message: "Invalid time range format, expected (start,end]",
				Data:    nil,
			})
		}
		start, err1 := strconv.ParseInt(parts[0], 10, 64)
		end, err2 := strconv.ParseInt(parts[1], 10, 64)
		if err1 != nil || err2 != nil {
			return c.JSON(http.StatusBadRequest, modal.Response{
				Code:    400,
				Message: "Invalid timestamp in time range",
				Data:    nil,
			})
		}
		startTime = time.Unix(start, 0)
		endTime = time.Unix(end, 0)
	} else {
		// Default to current week
		now := time.Now()
		startTime = now.AddDate(0, 0, -int(now.Weekday()))
		endTime = startTime.AddDate(0, 0, 7)
	}

	pipeline := []bson.M{
		{"$match": bson.M{
			"type": "1",
			"time": bson.M{
				"$gte": strconv.FormatInt(startTime.Unix(), 10),
				"$lte": strconv.FormatInt(endTime.Unix(), 10),
			},
		}},
		{"$addFields": bson.M{
			"date": bson.M{
				"$cond": bson.M{
					"if": bson.M{"$regexMatch": bson.M{
						"input": "$time",
						"regex": "^\\d+$",
					}},
					"then": bson.M{"$toDate": bson.M{
						"$multiply": bson.A{
							bson.M{"$toLong": "$time"},
							1000,
						},
					}},
					"else": bson.M{"$toDate": "$time"},
				},
			},
		}},
		{"$setWindowFields": bson.M{
			"sortBy": bson.M{"date": 1},
			"output": bson.M{
				"cumulativeCount": bson.M{
					"$sum": 1,
					"window": bson.M{
						"documents": []interface{}{"unbounded", "current"},
					},
				},
			},
		}},
		{"$group": bson.M{
			"_id": bson.M{
				"date": bson.M{"$dateToString": bson.M{
					"format": dateFormat,
					"date":   "$date",
				}},
				"status": "$status",
			},
			"count": bson.M{"$last": "$cumulativeCount"},
			"type":  bson.M{"$first": "$type"},
		}},
		{"$project": bson.M{
			"Time":   "$_id.date",
			"status": "$_id.status",
			"Type":   "$type",
			"count":  1,
			"_id":    0,
		}},
		{"$sort": bson.D{
			{Key: "Time", Value: -1},
			{Key: "count", Value: -1},
		}},
	}

	// Log the pipeline for debugging
	pipelineWrapper := bson.M{"pipeline": pipeline}
	pipelineJSON, err := bson.MarshalExtJSON(pipelineWrapper, false, false)
	if err != nil {
		c.Logger().SetOutput(os.Stdout)
		c.Logger().Errorf("Failed to marshal pipeline: %v", err)
	} else {
		fmt.Printf("Executing pipeline: %s\n", pipelineJSON) // Print to stdout directly
		c.Logger().SetOutput(os.Stdout)
		c.Logger().Infof("Executing pipeline: %s", pipelineJSON)
	}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		c.Logger().SetOutput(os.Stdout)
		c.Logger().Errorf("Aggregation failed: %v", err)
		return c.JSON(http.StatusInternalServerError, modal.Response{
			Code:    500,
			Message: "Failed to aggregate data",
			Data:    nil,
		})
	}
	defer cursor.Close(context.TODO())

	var results []bson.M
	if err := cursor.All(context.TODO(), &results); err != nil {
		c.Logger().SetOutput(os.Stdout)
		c.Logger().Errorf("Failed to decode results: %v", err)
		return c.JSON(http.StatusInternalServerError, modal.Response{
			Code:    500,
			Message: "Failed to process data",
			Data:    nil,
		})
	}

	c.Logger().SetOutput(os.Stdout)
	c.Logger().Infof("Returning %d aggregated documents", len(results))
	return c.JSON(http.StatusOK, modal.Response{
		Code:    200,
		Message: "",
		Data:    results,
	})
}

func GetTest(c echo.Context) error {
	collection := utils.GetDB().Database(utils.GetDBName()).Collection(utils.GetCollectionName())

	// Get and validate type parameter
	timeType := c.QueryParam("type")
	if timeType == "" {
		timeType = "week"
	}
	if timeType != "day" && timeType != "week" && timeType != "month" {
		return c.JSON(http.StatusBadRequest, modal.Response{
			Code:    400,
			Message: "Invalid type parameter, must be day/week/month",
			Data:    nil,
		})
	}

	// Set date format based on type
	dateFormat := "%Y-%m-%d"
	if timeType == "week" {
		dateFormat = "%Y-%U"
	} else if timeType == "month" {
		dateFormat = "%Y-%m"
	}

	// Parse time range parameter
	timeRange := c.QueryParam("time")
	var startTime, endTime time.Time
	if timeRange != "" {
		parts := strings.Split(strings.Trim(timeRange, "[]"), ",")
		if len(parts) != 2 {
			return c.JSON(http.StatusBadRequest, modal.Response{
				Code:    400,
				Message: "Invalid time range format, expected [start,end)",
				Data:    nil,
			})
		}
		start, err1 := strconv.ParseInt(parts[0], 10, 64)
		end, err2 := strconv.ParseInt(parts[1], 10, 64)
		if err1 != nil || err2 != nil {
			return c.JSON(http.StatusBadRequest, modal.Response{
				Code:    400,
				Message: "Invalid timestamp in time range",
				Data:    nil,
			})
		}
		startTime = time.Unix(start, 0)
		endTime = time.Unix(end, 0)
	} else {
		// Default to current week
		now := time.Now()
		startTime = now.AddDate(0, 0, -int(now.Weekday()))
		endTime = startTime.AddDate(0, 0, 7)
	}

	pipeline := []bson.M{
		{"$match": bson.M{
			"type": "1",
			"time": bson.M{
				"$gte": strconv.FormatInt(startTime.Unix(), 10),
				"$lt":  strconv.FormatInt(endTime.Unix(), 10),
			},
		}},
		{"$group": bson.M{
			"_id": bson.M{
				"date": bson.M{
					"$dateToString": bson.M{
						"format": dateFormat,
						"date": bson.M{
							"$cond": bson.M{
								"if": bson.M{"$regexMatch": bson.M{
									"input": "$time",
									"regex": "^\\d+$",
								}},
								"then": bson.M{"$toDate": bson.M{
									"$multiply": bson.A{
										bson.M{"$toLong": "$time"},
										1000,
									},
								}},
								"else": bson.M{"$toDate": "$time"},
							},
						},
					},
				},
				"status": "$status",
			},
			"count": bson.M{"$sum": 1},
			"type":  bson.M{"$first": "$type"},
		}},
		{"$project": bson.M{
			"Time":   "$_id.date",
			"status": "$_id.status",
			"Type":   "$type",
			"count":  1,
			"_id":    0,
		}},
		{"$sort": bson.M{"Time": -1}},
	}

	// Log the pipeline for debugging
	pipelineWrapper := bson.M{"pipeline": pipeline}
	pipelineJSON, err := bson.MarshalExtJSON(pipelineWrapper, false, false)
	if err != nil {
		c.Logger().SetOutput(os.Stdout)
		c.Logger().Errorf("Failed to marshal pipeline: %v", err)
	} else {
		fmt.Printf("Executing pipeline: %s\n", pipelineJSON) // Print to stdout directly
		c.Logger().SetOutput(os.Stdout)
		c.Logger().Infof("Executing pipeline: %s", pipelineJSON)
	}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		c.Logger().SetOutput(os.Stdout)
		c.Logger().Errorf("Aggregation failed: %v", err)
		return c.JSON(http.StatusInternalServerError, modal.Response{
			Code:    500,
			Message: "Failed to aggregate data",
			Data:    nil,
		})
	}
	defer cursor.Close(context.TODO())

	var results []bson.M
	if err := cursor.All(context.TODO(), &results); err != nil {
		c.Logger().SetOutput(os.Stdout)
		c.Logger().Errorf("Failed to decode results: %v", err)
		return c.JSON(http.StatusInternalServerError, modal.Response{
			Code:    500,
			Message: "Failed to process data",
			Data:    nil,
		})
	}

	c.Logger().SetOutput(os.Stdout)
	c.Logger().Infof("Returning %d aggregated documents", len(results))
	return c.JSON(http.StatusOK, modal.Response{
		Code:    200,
		Message: "",
		Data:    results,
	})
}
