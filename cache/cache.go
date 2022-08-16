package cache

import (
	"123/database"
	"123/models"
	"encoding/json"
	"fmt"
)

var m = make(map[string][]byte)

func AppendCache(data []byte, db *database.DbConnect) error {
	var order models.Order
	err := json.Unmarshal(data, &order)
	if err != nil {
		return err
	}
	orderUid := order.OrderUid
	_, ok := m[orderUid]
	if ok {
		fmt.Println("Order is already cached!")
		return nil
	}
	if orderUid == "" {
		fmt.Println("Invalid data sent")
		return nil
	}
	m[orderUid] = data
	db.PostOrder(order)
	fmt.Println("Appended to cache ; db")
	return nil
}

func GetOrder(orderUid string) models.Order {
	var order models.Order
	_, ok := m[orderUid]
	if ok {
		json.Unmarshal(m[orderUid], &order)
		return order
	}
	return models.Order{}
}

func UpdateCache(order models.Order, data []byte) {
	m[order.OrderUid] = data
}
