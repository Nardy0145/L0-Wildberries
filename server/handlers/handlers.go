package handlers

import (
	"123/cache"
	"123/database"
	"123/models"
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"title": "Orders",
	})
}

func GetOrder(c *gin.Context) {
	db := c.MustGet("db").(database.DbConnect)
	orderUid := c.Query("order_uid")
	order := cache.GetOrder(orderUid)
	if order.OrderUid == "" {
		data, err := db.GetOrder(orderUid)
		if err != nil {
			c.String(400, "GET ERROR: ", err)
			return
		}
		c.JSON(200, data)
	}
}

func PostOrder(c *gin.Context) {
	db := c.MustGet("db").(database.DbConnect)
	var data models.Order
	err := c.BindJSON(&data)
	if err != nil {
		c.String(400, "Post ERROR: \n\n", err)
		return
	}
	c.JSON(200, data)
	db.PostOrder(data)
}

func ResetTable(c *gin.Context) {
	db := c.MustGet("db").(database.DbConnect)
	db.ResetTable()
	c.String(200, "Table has been reset")
}
