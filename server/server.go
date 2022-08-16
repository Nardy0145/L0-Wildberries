package server

import (
	"123/database"
	"123/server/handlers"
	"github.com/gin-gonic/gin"
	"log"
)

func dbMiddleWare(db database.DbConnect) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

func initRoutes(r *gin.Engine) {
	r.GET("/", handlers.Index)
	r.GET("/get_order", handlers.GetOrder)
	r.POST("/post_order", handlers.PostOrder)
	r.GET("/reset_table", handlers.ResetTable)
}

func RunServer() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Run error: DB INIT ERR, ", err)
	}
	r.Use(dbMiddleWare(*db))
	r.LoadHTMLGlob("templates/*.html")
	initRoutes(r)
	err = r.Run(":8080")
	if err != nil {
		log.Fatal("Server running error: ", err)
	}
}
