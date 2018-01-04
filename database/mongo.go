package mongo

import (
	"fmt"

	"../config"
	"gopkg.in/mgo.v2"
)

var session *mgo.Session

func init() {
	var err error
	session, err = mgo.Dial(config.MongoURL)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
}

// Session for copy anthoer session for any request
func Session() *mgo.Session {
	return session.Copy()
}
