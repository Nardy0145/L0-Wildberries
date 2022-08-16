package database

import (
	"123/models"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DbConnect struct {
	Db *sqlx.DB
}

func InitDB() (*DbConnect, error) {
	db, err := sqlx.Connect("postgres", "user=root password=root host=localhost dbname=postgres sslmode=disable")
	if err != nil {
		return nil, err
	}
	return &DbConnect{Db: db}, nil
}

func (db DbConnect) GetOrder(orderUid string) (models.Order, error) {
	var rawOrder models.RawOrder
	var order models.Order
	err := db.Db.Get(&rawOrder, "SELECT data FROM wborders WHERE order_uid = $1", orderUid)
	if err != nil {
		return models.Order{}, err
	}
	tempJson := rawOrder.Data
	tempJson, _ = tempJson.MarshalJSON()
	err = json.Unmarshal(tempJson, &order)
	if err != nil {
		return models.Order{}, err
	}
	return order, nil
}

func (db DbConnect) PostOrder(data models.Order) {
	dataJson, _ := json.Marshal(data)
	_, err := db.Db.Exec("insert into  wborders (order_uid, data) values ($1, $2)", data.OrderUid, dataJson)
	if err != nil {
		fmt.Println("DB: POST ORDER: ", err)
	}
}

func (db DbConnect) CreateTable() {
	_, err := db.Db.Exec(`
	CREATE TABLE wborders(
	    id BIGSERIAL PRIMARY KEY,
	    order_uid VARCHAR UNIQUE NOT NULL,
	    data JSON NOT NULL
	);
	CREATE INDEX idx_order_uid ON wborders USING HASH (order_uid);
	`)
	if err != nil {
		fmt.Println(err)
	}
}

func (db DbConnect) ResetTable() {
	db.Db.Exec("DROP TABLE wborders")
	db.CreateTable()
	db.Db.Exec(`
ALTER SEQUENCE wborders_id_seq RESTART WITH 1;
UPDATE wborders SET id=nextval('wborders_id_seq');
`)
}
