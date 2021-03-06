package controllers

import (
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"../models"
	"github.com/gin-gonic/gin"
)

type Result struct {
	Success bool        `json:"success"`
	Code    int8        `json:"code"`
	Data    interface{} `json:"data"`
}

// ErrorResponse contain error message Code and Message
// Code api in /api/code.md
type ErrorResponse struct {
	Code    int16    `json:"code"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}

// GetTask returns the task list
func GetTask(c *gin.Context) {
	id := c.Param("id")
	list := make([]models.Task, 1)

	var err error
	collection := c.MustGet("DB").(*mgo.Database).C("task")

	// find one
	if id != "" {
		defer func() {
			// ObjectIdHex error
			if err := recover(); err != nil {
				c.JSON(200, []interface{}{})
			}
		}()
		err = collection.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&list[0])
	} else {
		option := make(map[string]interface{})

		// set options for find
		models.GetOptionForGetTask(c, option)

		err = collection.Find(option).All(&list)
	}

	if err != nil {
		fmt.Println(err)
		c.JSON(500, &ErrorResponse{Code: 301, Message: "database.connect.error"})
		return
	}

	c.JSON(200, list)
}

// CreateTask for create task
func CreateTask(c *gin.Context) {
	nowTime := time.Now().Unix() * 1000
	id := bson.NewObjectId()

	// 传入格式校验
	item := &models.Task{
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
	c.ShouldBindJSON(item)
	// no name
	if item.Name == "" {
		c.JSON(400, &ErrorResponse{Code: 101, Message: "task.name.empty"})
		return
	}

	models.CreateTask(item, func(success bool) {
		if success {
			c.JSON(200, map[string]string{"id": id.Hex()})
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
	var idFromURL = c.Param("id")
	var object interface{}
	nowTime := time.Now().Unix() * 1000

	c.ShouldBindJSON(&object)

	// Determine the type of field to be modified
	item := object.(map[string]interface{})
	for k := range item {
		switch k {
		case "type", "difficult":
			item[k] = int8(item[k].(float64))
		case "name", "deadline", "memo":
			item[k] = item[k].(string)
		case "finish":
			item[k] = item[k].(bool)
		case "Hours":
			item[k] = int(item[k].(float64))
		}
	}

	if item["id"] == nil {
		if idFromURL != "" {
			item["id"] = idFromURL
		} else {
			c.JSON(200, &Result{Success: false, Code: 1, Data: nil})
			return
		}
	}

	// modify update time by now time
	item["updateTime"] = nowTime
	id := item["id"].(string)

	models.UpdateTask(item, func(success bool) {
		// fmt.Println("我想知道success到底是", success)
		var code int8
		var data interface{}
		if success {
			code = 0
			data = id
		} else {
			code = 1
			data = nil
		}
		fmt.Println(success)
		// fmt.Println(code)
		fmt.Println(data)
		c.JSON(200, &Result{
			Success: success,
			Code:    code,
			Data:    data,
		})
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

	collection := c.MustGet("DB").(*mgo.Database).C("task")

	idListLen := len(idList)
	channel := make(chan string, idListLen)

	go models.DeleteTaskMany(channel, collection, idList)

	errIDList := make([]string, 0, idListLen)
	for i := 0; i < idListLen; i++ {
		// If someone false, response error
		errID := <-channel
		if errID != "" {
			errIDList = append(errIDList, errID)
		}
	}

	errIDListLen := len(errIDList)
	if errIDListLen != 0 {
		if errIDListLen == idListLen {
			// delete all error
			c.JSON(400, &ErrorResponse{Code: 202, Message: "task.delete.all.error"})
		} else {
			c.JSON(400, &ErrorResponse{
				Code:    201,
				Message: "task.delete.some.error",
				Data:    errIDList,
			})
		}

		return
	}

	c.String(200, "")
}
