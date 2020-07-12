package upload

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"packform/utils/injector"
	"strconv"
	"sync"
	"time"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UploadOrderItems(context *injector.DepContainer) gin.HandlerFunc{

	return func(c *gin.Context){

		var wg sync.WaitGroup

		orders_list := make(chan []string)
		file, _, _ := c.Request.FormFile("file");

		
		wg.Add(1)

		go ReadFromCsv(file, orders_list)

		go func(list chan[]string, wg *sync.WaitGroup,context *injector.DepContainer){
			defer wg.Done()

			db := context.GetDbContext() 
			query := `INSERT INTO order_item(order_id,price_per_unit,quantity, product) VALUES($1,$2,$3,$4)`

			tx, err := db.Begin()
			
			txStmt, err := tx.Prepare(query)

			if err != nil {
				tx.Rollback()
				fmt.Println("Err", err)
				return
			}

			for item := range list {
				order_id, err := strconv.Atoi(item[1])
				if err != nil {
					fmt.Println(err)
				}

				ppu, err := decimal.NewFromString(item[2])
				if err != nil {
					fmt.Println(err)
					ppu = decimal.NewFromInt(0)
				}

				qty, err := strconv.Atoi(item[3])
				if err != nil {
					fmt.Println(err)
				}

				product  := item[4]

				_ , err = txStmt.Exec(order_id, ppu, qty, product)
				if err != nil {
					tx.Rollback()
					fmt.Println(err)
					return
				}
			}

			err = tx.Commit()
			if err != nil {
				fmt.Println(err)
				return
			}

		}(orders_list, &wg, context)
				
		wg.Wait()

		c.JSON(200, gin.H{
			"message": "Successfully Added",
		})
	}
}

func UploadDeliveries(context *injector.DepContainer) gin.HandlerFunc{

	return func(c *gin.Context){

		var wg sync.WaitGroup

		orders_list := make(chan []string)
		file, _, _ := c.Request.FormFile("file");

		
		wg.Add(1)

		go ReadFromCsv(file, orders_list)



		go func(list chan[]string, wg *sync.WaitGroup,context *injector.DepContainer){
			defer wg.Done()

			db := context.GetDbContext() 
			query := `
				INSERT INTO deliveries(order_item_id,delivered_qty) 
				VALUES($1,$2)
			`

			tx, err := db.Begin()
			
			txStmt, err := tx.Prepare(query)

			if err != nil {
				tx.Rollback()
				fmt.Println("Err", err)
				return
			}

			for item := range list {
				order_item_id, err := strconv.Atoi(item[1])
				if err != nil{
					fmt.Println(err)
				} 

				delivered_qty, err := strconv.Atoi(item[2])
				if err != nil{
					fmt.Println(err)
				} 

				_ , err = txStmt.Exec(order_item_id, delivered_qty)
				if err != nil {
					tx.Rollback()
					fmt.Println(err)
					return
				}
			}

			err = tx.Commit()
			if err != nil {
				fmt.Println(err)
				return
			}
		}(orders_list, &wg, context)
				
		wg.Wait()

		c.JSON(200, gin.H{
			"message": "Successfully Added",
		})
	}
}

func UploadCustomers(context *injector.DepContainer) gin.HandlerFunc{

	return func(c *gin.Context){

		var wg sync.WaitGroup

		orders_list := make(chan []string)
		file, _, _ := c.Request.FormFile("file");

		
		wg.Add(1)

		go ReadFromCsv(file, orders_list)


		go func(list chan[]string, wg *sync.WaitGroup,context *injector.DepContainer){
			defer wg.Done()

			ctx, cancel := createNewCtx()
			defer cancel()
			
			mongoSession, err:= context.GetMongoClient().StartSession()
			if err != nil{
				fmt.Println(err)
				return
			}

			err = mongoSession.StartTransaction()
			if err != nil{
				fmt.Println(err)
				return
			}

			err = mongo.WithSession(ctx, mongoSession, func(sc mongo.SessionContext) error {
				customerCollection := context.GetMongoDbContext().Collection("customers");

				for item := range list {
					fmt.Println(item)

					iCompanyID,err := strconv.Atoi(item[4])
					if err != nil{
						mongoSession.AbortTransaction(ctx)
						fmt.Println(err) 
					}

					_ , err = customerCollection.InsertOne(ctx, bson.D{
						primitive.E{Key: "user_id", Value: item[0]},
						primitive.E{Key: "login", Value: item[1]},
						primitive.E{Key: "password", Value: item[2]},
						primitive.E{Key: "name", Value: item[3]},
						primitive.E{Key: "company", Value: iCompanyID},
						primitive.E{Key: "credit_cards", Value: strings.Split(item[5], ",")},
					})

					if err != nil{
						mongoSession.AbortTransaction(ctx)
						fmt.Println(err)
					}
				}

				mongoSession.CommitTransaction(ctx)
				return err
			})
			if err != nil{
				fmt.Println(err)
				return
			}
	

		}(orders_list, &wg, context)
				
		wg.Wait()

		c.JSON(200, gin.H{
			"message": "Successfully Added",
		})
	}
}


func UploadCompanies(context *injector.DepContainer) gin.HandlerFunc{

	return func(c *gin.Context){

		var wg sync.WaitGroup

		orders_list := make(chan []string)
		file, _, _ := c.Request.FormFile("file");
		
		wg.Add(1)

		go ReadFromCsv(file, orders_list)

		go func(list chan[]string, wg *sync.WaitGroup,context *injector.DepContainer){
			defer wg.Done()

			ctx, cancel := createNewCtx()
			defer cancel()
			
			mongoSession, err:= context.GetMongoClient().StartSession()
			if err != nil{
				fmt.Println(err)
				return
			}

			err = mongoSession.StartTransaction()
			if err != nil{
				fmt.Println(err)
				return
			}

			err = mongo.WithSession(ctx, mongoSession, func(sc mongo.SessionContext) error {
				companiesCollection := context.GetMongoDbContext().Collection("customer_companies");

				for item := range list {
					fmt.Println(item)

					iCompanyID,err := strconv.Atoi(item[0])
					if err != nil{
						mongoSession.AbortTransaction(ctx)
						fmt.Println(err) 
					}

					_ , err = companiesCollection.InsertOne(ctx, bson.D{
						primitive.E{Key: "company_id", Value: iCompanyID},
						primitive.E{Key: "company_name", Value: item[1]},
					})

					if err != nil{
						mongoSession.AbortTransaction(ctx)
						fmt.Println(err)
					}
				}

				mongoSession.CommitTransaction(ctx)
				return err
			})
			if err != nil{
				fmt.Println(err)
				return
			}
	

		}(orders_list, &wg, context)
				
		wg.Wait()

		c.JSON(200, gin.H{
			"message": "Successfully Added",
		})
	}
}


func createNewCtx() (context.Context, context.CancelFunc){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx,cancel
}

func ReadFromCsv(file multipart.File, list chan []string){
	reader := csv.NewReader(file)
	defer close(list)

	// Skip the first line which is header
	_ , err := reader.Read()
	if(err != nil){
		fmt.Println("err", err)
	}

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break;
		}else if err != nil {
			log.Fatal(err)
		}
		list <- line
	}
}