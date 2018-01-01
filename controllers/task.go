package controllers

import (
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"

	"../models"
	"github.com/gin-gonic/gin"
)

type Result struct {
	Success bool        `json:"success"`
	Code    int8        `json:"code"`
	Data    interface{} `json:"data"`
}

// CreateTask for create task
func CreateTask(c *gin.Context) {
	nowTime := time.Now().Unix() * 1000
	id := bson.NewObjectId()

	task := models.Task{
		ID:         id,
		Name:       "",
		Type:       1,
		Difficult:  2,
		Deadline:   time.Now().Format("2006-01-02 15:04:05"),
		Hours:      8,
		Finish:     false,
		Memo:       "",
		CreateTime: nowTime,
		UpdateTime: nowTime,
	}
	c.ShouldBindJSON(&task)

	if task.Name == "" {
		c.JSON(200, Result{
			Success: false,
			Code:    1,
			Data:    nil,
		})
		return
	}

	models.CreateTask(task, func(success bool) {
		if success {
			c.JSON(200, Result{
				Success: true,
				Code:    0,
				Data:    id.Hex(),
			})
		} else {
			c.JSON(200, Result{
				Success: false,
				Code:    1,
				Data:    nil,
			})
		}
	})
}

// DeleteTask for delete task
func DeleteTask(c *gin.Context) {
	var idList = make([]string, 1)
	id := c.Param("id")
	if id != "" {
		idList[0] = id
	} else {
		c.ShouldBindJSON(&idList)
	}
	models.DeleteTask(idList, func(success bool) {
		if success {
			c.JSON(200, Result{
				Success: true,
				Code:    0,
				Data:    nil,
			})
		} else {
			c.JSON(200, Result{
				Success: false,
				Code:    1,
				Data:    nil,
			})
		}
	})
}

// UpdateTask for update task
func UpdateTask(c *gin.Context) {
	c.String(200, "Update task success")
}

// GetTask returns the task list
func GetTask(c *gin.Context) {
	id := c.Param("id")
	// find one
	if id != "" {
		models.GetOneTask(id, func(task []models.Task) {
			if task != nil {
				c.JSON(200, Result{
					Success: true,
					Code:    0,
					Data:    task[0],
				})
			} else {
				c.JSON(200, Result{
					Success: true,
					Code:    0,
					Data:    nil,
				})
			}
		})
	} else {
		option := make(map[string]interface{})
		finishString := c.Query("finish")
		typeString := c.Query("type")
		difficultString := c.Query("difficult")
		minHoursString := c.Query("minHours")
		maxHoursString := c.Query("maxHours")

		if finishString != "" {
			option["finish"], _ = strconv.ParseBool(finishString)
		}
		if typeString != "" {
			option["type"], _ = strconv.ParseInt(typeString, 10, 8)
		}
		if difficultString != "" {
			option["difficult"], _ = strconv.ParseInt(difficultString, 10, 8)
		}

		if minHoursString != "" && maxHoursString != "" {
			minHour, _ := strconv.ParseInt(minHoursString, 10, 32)
			maxHour, _ := strconv.ParseInt(maxHoursString, 10, 32)

			option["hours"] = map[string]int{"$gte": int(minHour), "$lte": int(maxHour)}
		} else if minHoursString != "" && maxHoursString == "" {
			minHour, _ := strconv.ParseInt(minHoursString, 10, 32)

			option["hours"] = map[string]int{"$gte": int(minHour)}
		} else if minHoursString == "" && maxHoursString != "" {
			maxHour, _ := strconv.ParseInt(maxHoursString, 10, 32)

			option["hours"] = map[string]int{"$lte": int(maxHour)}
		}

		models.GetTask(option, func(task []models.Task) {
			c.JSON(200, Result{
				Success: true,
				Code:    0,
				Data:    task,
			})
		})
	}
}
