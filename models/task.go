package models

import (
	"fmt"
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	MONGO_URI = "mongodb://localhost:27017/ly"
)

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

// GetTask that get all task list from db.
func GetTask(option map[string]interface{}, cb func(task []Task)) {

	fmt.Println("查询条件:", option)
	session, err := mgo.Dial(MONGO_URI)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("ly").C("task")

	var task []Task
	err = c.Find(option).All(&task)
	if err != nil {
		log.Fatal(err)
	}

	cb(task)
}

// GetOneTask that get one task by id for db.
func GetOneTask(id string, cb func(task []Task)) {
	session, err := mgo.Dial(MONGO_URI)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("ly").C("task")

	var task Task
	defer func() {
		if err := recover(); err != nil {
			cb(nil)
		}
	}()

	err = c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&task)
	if err != nil {
		panic(err)
	}

	cb([]Task{task})
}

// CreateTask that create one task
func CreateTask(task Task, cb func(success bool)) {
	session, err := mgo.Dial(MONGO_URI)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("ly").C("task")

	err = c.Insert(task)
	if err != nil {
		log.Fatal(err)
		cb(false)
		return
	}

	cb(true)
}

// DeleteTask that delete the task by id from db.
func DeleteTask(idList []string, cb func(success bool)) {
	session, err := mgo.Dial(MONGO_URI)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("ly").C("task")

	errCount := 0
	idListLen := len(idList)
	channel := make(chan bool, idListLen)
	for _, id := range idList {

		go func(id string) {
			defer func() {
				if err := recover(); err != nil {
					errCount++
					channel <- true
				}
			}()

			err := c.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
			if err != nil {
				panic(err)
			}
			channel <- true
		}(id)
	}
	for i := 0; i < idListLen; i++ {
		<-channel
	}

	if errCount > 0 {
		cb(false)
	} else {
		cb(true)
	}
}
