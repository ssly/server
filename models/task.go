package models

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"../config"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const MongoURL = config.MongoURL

// Task is the response data for task.
type Task struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	Name       string        `json:"name"`
	Type       int8          `json:"type"`
	Difficult  int8          `json:"difficult"`
	Deadline   string        `json:"deadline"`
	Hours      int           `json:"hours"`
	Finish     bool          `json:"finish"`
	Memo       string        `json:"memo"`
	CreateTime int64         `json:"createTime"`
	UpdateTime int64         `json:"updateTime"`
}

// CreateTask that create one task
func CreateTask(task *Task, cb func(success bool)) {
	defer func() {
		if err := recover(); err != nil {
			cb(false)
		}
	}()

	session, err := mgo.Dial(MongoURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("ly").C("task")

	err = c.Insert(task)
	if err != nil {
		panic(err)
	}

	cb(true)
}

// UpdateTask that modify one task from db.
func UpdateTask(item map[string]interface{}, cb func(success bool)) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			cb(false)
		}
	}()

	id := item["id"].(string)
	delete(item, "id")

	session, err := mgo.Dial(MongoURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("ly").C("task")

	err = c.Update(bson.M{"_id": bson.ObjectIdHex(id)}, item)
	if err != nil {
		panic(err)
	}
	cb(true)
}

// GetOptionForGetTask for get options for retrieve task
// contain finish / type / difficult / minHours / maxHours
func GetOptionForGetTask(c *gin.Context, option map[string]interface{}) {
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
}

// DeleteTaskMany for delete many task in datebase
func DeleteTaskMany(channel chan string, c *mgo.Collection, idList []string) {
	for _, id := range idList {
		go func(id string) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("Task: Call DeleteTaskMany error,", err)
					channel <- id
				}
			}()

			err := c.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
			if err != nil {
				fmt.Println("Task: Call DeleteTask error, this id is", id)
				panic(err)
			}

			channel <- ""
		}(id)
	}
}
