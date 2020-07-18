package order

import (
	// "net/http"
	"packform/utils/injector"
	"packform/api/model/orders"
	"go.mongodb.org/mongo-driver/bson"
	"net/url"
	"sync"
	"fmt"
	"context"
	"time"
	"strconv"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/shopspring/decimal"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "github.com/tidwall/gjson"
)


/**
	Preemptively get required data from sql to do cross-db joins. 
	Since we need to get the customer's name only. We do not need
	to do computationally heavy stuff like sql joins
	just get DISTINCT names from orders and run a query to MongoDb
**/

type Result struct {
	Companies 		*NestedResult 			`bson:"companies"`
	Company_id		int									`bson:"company"`
	Credit_cards	[]string						`bson:"credit_cards"`
	Login 				string							`bson:"login"`
	Password 			string							`bson:"password"`
	User_id 			string							`bson:"user_id"`
}

type NestedResult struct {

	Company_id 		int 											`bson:"company_id"`
	Company_name	string 											`bson:"company_name"`
}

//Column in SQL declared as Numeric(16,4)
const DECIMAL_MAX = 999999999999.9999

func _joinMongoDb(dict chan map[string]string, context *injector.DepContainer,
		orderName string, custId string, page int) {

	defer close(dict)
	db := context.GetDbContext()
	var query string = `
		SELECT DISTINCT customer_id FROM orders
		WHERE order_name LIKE coalesce(NULLIF($1,''),'%')
		AND customer_id LIKE coalesce(NULLIF($2,''), '%')
		LIMIT 5 
		OFFSET coalesce(greatest($3,0),0) 
	`
	stmt, _ := db.Preparex(query)

	lookupStage := bson.D{
		{"$lookup", bson.D{
			{"from", "customer_companies"}, 
			{"localField","company"},
			{"foreignField","company_id"},
			{"as","companies"},
		}},
	}

	unwindStage := bson.D{
		{"$unwind","$companies"},
	}

	dictName := map[string]string{}

	var result []string

	/** Populate Dictionary with Names**/
	err := stmt.Select(&result,orderName,custId,page)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(result)

	for _, elem := range result {
		dictName[elem] = ""
	}


	customerCollection := context.GetMongoDbContext().Collection("customers");
	ctx, cancel := createNewCtx()
	defer cancel()
	
	loadedCursor , err :=  customerCollection.Aggregate(ctx,mongo.Pipeline{lookupStage,unwindStage})

	for loadedCursor.Next(ctx){
		var NestRes NestedResult

		Result := Result{
			Companies: &NestRes,
		}

		// fmt.Println(loadedCursor)

		err := loadedCursor.Decode(&Result)

		if err != nil {
				fmt.Println("cursor.Next() error:", err)
		} else {
				dictName[Result.User_id]= Result.Companies.Company_name
		}
	}

	dict <- dictName


}

/**
	Encode Query String in Client
**/

/** Prepared Statements parameters
	TODO: Structify
	$1: order_name
	$2: customer_id
	$3: created_at >
	$4: created_at <
	$5: page
	$6: total_amount >
	$7: total_amount < 
	$8: c_delivered_amount >
	$9: c_delivered_amount <
**/

func GetOrdersCollection(context *injector.DepContainer) gin.HandlerFunc{
	db := context.GetDbContext()
	query := `
	SELECT 
			order_id,
			TO_CHAR(MIN(created_at), 'Mon DDth, HH:MM AM ') AS order_date,
			MIN(customer_id) AS customer_id,
			MIN(order_name) AS order_name,
			MIN(created_at) AS order_date,
			SUM(price_per_unit * quantity) AS total_amount,
			COALESCE (SUM(price_per_unit * delivered_qty), 0.0000) AS c_delivered_amount
		FROM (
			SELECT * FROM orders
			WHERE order_name LIKE coalesce(NULLIF($1,''),'%')
			AND customer_id LIKE coalesce(NULLIF($2,''), '%')
			AND created_at::date > coalesce($3, '1970-01-01')::date
			AND created_at::date < coalesce($4, '9999-01-01')::date
			OFFSET coalesce(greatest($5*5,0),0) 
		) AS filtered_orders
		LEFT JOIN order_item ON filtered_orders.id = order_item.order_id
		LEFT JOIN deliveries ON order_item.id = deliveries.order_item_id
		GROUP BY order_id
		HAVING coalesce(SUM(price_per_unit * quantity),0) >= coalesce($6::Numeric,0)
		AND coalesce(SUM(price_per_unit * quantity),0) <= coalesce($7::Numeric,0)
 		AND coalesce(SUM(price_per_unit * delivered_qty),0) >= coalesce($8::Numeric,0)
		AND coalesce(SUM(price_per_unit * delivered_qty),0) <= coalesce($9::Numeric,0)
	`

	stmt, err := db.Preparex(query)
	if err != nil {
		panic(err)
	}
	

	return func(c *gin.Context){

		var qPage string = c.Query("page")
		var qOrderName string = c.Query("order_name")
		var qCustId string = c.Query("cust_id")

		//Query Strings
		var qDeliveredMin string = c.Query("delivered_min")
		var qDeliveredMax string = c.Query("delivered_max")
		var qTotalMin string = c.Query("total_min")
		var qTotalMax string = c.Query("total_max")

		var qStartDate string = c.Query("startDate")
		var qEndDate string = c.Query("endDate")
		//


		orderName, err := url.QueryUnescape(qOrderName)
		custId , err:= url.QueryUnescape(qCustId)

		delivered_min, err := url.QueryUnescape(qDeliveredMin)
		delivered_max , err:= url.QueryUnescape(qDeliveredMax)

		total_min, err := url.QueryUnescape(qTotalMin)
		total_max , err:= url.QueryUnescape(qTotalMax)

		
		// cust_company, err := url.QueryUnescape(qCustCompany)

		start_date , err:= url.QueryUnescape(qStartDate)
		if start_date == "" {
			start_date = "1970-01-01"
		}

		end_date , err:= url.QueryUnescape(qEndDate)
		if end_date == "" {
			end_date = "9999-01-01"
		}

		ipage, err := strconv.Atoi(qPage)
		if err != nil {
			ipage = 0
		}
		
		iDelivered_min , err := strconv.Atoi(delivered_min)
		if err!= nil {
			iDelivered_min = 0
		}

		iDelivered_max , err := decimal.NewFromString(delivered_max)
		if err!= nil {
			iDelivered_max = decimal.NewFromFloat(DECIMAL_MAX)
		}

		iTotal_min , err := strconv.Atoi(total_min)
		if err!= nil {
			iTotal_min = 0
		}

		iTotal_max , err := decimal.NewFromString(total_max)
		if err!= nil {
			iTotal_max = decimal.NewFromFloat(DECIMAL_MAX)
		}
		

		fmt.Println(orderName)
		fmt.Println(custId)
		fmt.Println(ipage)
		fmt.Println(iDelivered_min)
		fmt.Println(iDelivered_max)
		fmt.Println(iTotal_min)
		fmt.Println(iTotal_max)
		fmt.Println("Sdate",start_date)
		fmt.Println("Edatec",end_date)


		orderCollection := orders.NewCollections().OrderList


		/** MongoDB**/
		ch_dict := make(chan map[string]string)
		go _joinMongoDb(ch_dict, context, orderName, custId,ipage)
		
		var result_dict map[string]string
		result_dict = <- ch_dict

		// for k,v := range result_dict {
		// 	println(k,v)
		// }

		// fmt.Println(stmt)

		// iTotal_max, err = decimal.NewFromString("99999.99")
		// iDelivered_max, err = decimal.NewFromString("99999.99")

		if err != nil {
			fmt.Println(err)
		}

		// sql statement takes float as []uint_8 -> cast in psql
		errsql := stmt.Select(&orderCollection,
			orderName,
			custId,
			start_date,
			end_date,
			ipage,
			iTotal_min,
			iTotal_max,
			iDelivered_min,
			iDelivered_max,
		)

		if errsql != nil{
			fmt.Println(errsql)
			c.AbortWithStatusJSON(500, gin.H{"status": false, "message": errsql.Error()})
			return
		}

		for i, order := range orderCollection {
			orderCollection[i].Customer_company = result_dict[order.Customer_Id]
		}
		c.JSON(200, gin.H{
			"message": "OrderCollection" ,
			"result":orderCollection,
		})

		
	}
}



func TestMongo(context *injector.DepContainer) gin.HandlerFunc{
	db := context.GetDbContext()
	var query string = `
		SELECT DISTINCT customer_id FROM orders
		WHERE order_name LIKE coalesce(NULLIF($1,''),'%')
		AND customer_id LIKE coalesce(NULLIF($2,''), '%')
		LIMIT 5 
		OFFSET coalesce($3,0) 
	`
	stmt, _ := db.Preparex(query)

	lookupStage := bson.D{
		{"$lookup", bson.D{
			{"from", "customer_companies"}, 
			{"localField","company"},
			{"foreignField","company_id"},
			{"as","companies"},
		}},
	}

	unwindStage := bson.D{
		{"$unwind","$companies"},
	}

	dictName := map[string]string{}


	return func(c *gin.Context){

		var result []string

		/** Populate Dictionary with Names**/
		err := stmt.Select(&result,"","",0)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, elem := range result {
			dictName[elem] = ""
		}


		customerCollection := context.GetMongoDbContext().Collection("customers");
		ctx, cancel := createNewCtx()
		defer cancel()
	
		loadedCursor , err :=  customerCollection.Aggregate(ctx,mongo.Pipeline{lookupStage,unwindStage})
	
		for loadedCursor.Next(ctx){
			var NestRes NestedResult
			Result := Result{
				Companies: &NestRes,
			}

			err := loadedCursor.Decode(&Result)

			if err != nil {
					fmt.Println("cursor.Next() error:", err)
					return
			} else {
					dictName[Result.User_id]= Result.Companies.Company_name
			}
		}

		c.JSON(200, gin.H{
			"message": "OrderCollection" ,
		})
	}
}



func AddOrder(context *injector.DepContainer) gin.HandlerFunc{

	return func(c *gin.Context){
		var wg sync.WaitGroup

		orders_list := make(chan []string)
		file, _, _ := c.Request.FormFile("file");

		
		wg.Add(1)
		go ReadFromCsv(file,orders_list)
		go AddOrdersToDb(orders_list, &wg, context)
		wg.Wait()

		c.JSON(200, gin.H{
			"message": "Add",
		})
	}
}


func createNewCtx() (context.Context, context.CancelFunc){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx,cancel
}