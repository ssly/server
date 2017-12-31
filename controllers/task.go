package controllers

import (
	"strconv"

	"../models"
	"github.com/gin-gonic/gin"
)

type Result struct {
	Success bool          `json:"success"`
	Code    int8          `json:"code"`
	Data    []models.Task `json:"data"`
}

// CreateTask for create task
func CreateTask(c *gin.Context) {
	c.String(200, "Create task uccess")
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
				Success: true,
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
				c.JSON(200, struct {
					Success bool        `json:"success"`
					Code    int8        `json:"code"`
					Data    models.Task `json:"data"`
				}{
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
