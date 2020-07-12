package customers

import (
	"fmt"
)

type Orders struct{
	Id 					int  									`json:"order_id" db:"order_id"`
	Created_At  string 								`json:"created_at" db:"order_date"`
	Order_Name	string								`json:"order_name" db:"order_name"`
	Customer_Id string								`json:"customer_id" db:"customer_id"`
}

type OrderCollections struct {
	OrderList []Orders
}

func NewCollections() *OrderCollections{
	return &OrderCollections{
		OrderList: []Orders{} ,
	}
}

func New() *Orders{
	return &Orders{}
}