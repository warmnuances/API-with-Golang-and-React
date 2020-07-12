package orders

import (
	"github.com/shopspring/decimal"
)

type Orders struct{
	Id 								int  									`json:"order_id" db:"order_id"`
	Created_At  			string 								`json:"created_at" db:"order_date"`
	Order_Name				string								`json:"order_name" db:"order_name"`
	Customer_Id 			string								`json:"customer_id" db:"customer_id"`
	Total_amount 			*decimal.Decimal			`json:"total_amount" db:"total_amount"`
	Delivered_amount 	*decimal.Decimal			`json:"delivered_amount" db:"c_delivered_amount"`
	Customer_company 	string								`json:"customer_company"`
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