package database

import "github.com/gin-gonic/gin"

// Session is interface of mgo.Session and further databasse
// session type
type Session interface {
	Close()
}

// Connection is interface handle database connecting
// for manage session
type Connection interface {
	Connect() gin.HandlerFunc
}
