package middleware

import (
	mongo "../database"
	"github.com/gin-gonic/gin"
)

// ConnectDB copy mgo session for any request
func ConnectDB(c *gin.Context) {
	session := mongo.Session()
	defer session.Close()

	c.Set("DB", session.DB("ly"))

	c.Next()
}
