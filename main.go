package main

import (
	"os"
	"packform/api/controller/order"
	"packform/api/controller/upload"
	"packform/utils/injector"

	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

var context *injector.DepContainer

func Init(){
	context = injector.InitContainer()
	gotenv.Load()
	context.SetDbContext(os.Getenv("DB_HOST"), os.Getenv("CONN_STR"))
	context.InitTables()
	context.SetMongoContext(os.Getenv("MONGO_URI"))

}

func SetupRouter() *gin.Engine{
	Init()

	r := gin.Default()
	OrderRouter := r.Group("/orders")
	{
		OrderRouter.GET("/mongo", order.TestMongo(context))
		OrderRouter.GET("/all", order.GetOrdersCollection(context))
		OrderRouter.POST("/bulk", order.AddOrder(context))
	}
	AddRouter := r.Group("/add")
	{
		AddRouter.POST("/orderitem", upload.UploadOrderItems(context))
		AddRouter.POST("/delivery", upload.UploadDeliveries(context))
		AddRouter.POST("/customers", upload.UploadCustomers(context))
		AddRouter.POST("/companies", upload.UploadCompanies(context))
	}

	return r
}


func main()  {

	r := SetupRouter()
	r.Run(":8080")
}