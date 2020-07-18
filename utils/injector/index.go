package injector

import (
	"fmt"
	"context"
	"time"
	"log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options" 
)

type DepContainer struct{
	db *sqlx.DB
	mongo *mongo.Client
}

func InitContainer() *DepContainer{
	instance := &DepContainer{}
	return instance
}

// TODO: get schema as args
func (depContainer *DepContainer) InitTables(){
	var schema  = `
		CREATE TABLE IF NOT EXISTS orders(
			id 			SERIAL PRIMARY KEY,
			order_name 	VARCHAR(64)	NOT NULL UNIQUE , 
			customer_id VARCHAR(64) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL default current_timestamp
	);
	
	CREATE TABLE IF NOT EXISTS order_item(
		id 			SERIAL PRIMARY KEY,
		order_id 	INT NOT NULL REFERENCES orders (id) ON UPDATE CASCADE ON DELETE CASCADE,
		price_per_unit NUMERIC(16,4) ,
		quantity	INT NOT NULL,
		product 	VARCHAR(64) NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS deliveries(
		id 			SERIAL PRIMARY KEY,
		order_item_id 	INT NOT NULL REFERENCES order_item (id) ON UPDATE CASCADE ON DELETE CASCADE,
		delivered_qty	INT NOT NULL
	);
	`	

	depContainer.db.MustExec(schema)
}

func (depContainer *DepContainer) SetDbContext(driverName string, dataSourceName string){

	fmt.Println("Added Database")
	db, err := sqlx.Connect(driverName,dataSourceName)

	ping_err := db.Ping()
	if ping_err != nil {
		fmt.Println(ping_err)
	}

	if err != nil {
		log.Panicln(err)
		log.Panic("Database Instance not initialised")
		depContainer.db = nil
	}else{
		depContainer.db = db
	}	
}

func (depContainer *DepContainer) GetDbContext() *sqlx.DB{
	return depContainer.db
}

func (depContainer *DepContainer) SetMongoContext(mongoUri string){
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUri))
	if err != nil {
		fmt.Println(err)
	}else{
		depContainer.mongo = client
		/** Initialise documents**/
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    err = client.Connect(ctx)
    if err != nil {
				log.Fatal(err)
		}
		defer cancel()

		err = client.Ping(ctx, nil)
		if err != nil{
			log.Fatal(err)
		}
	}
}

func (depContainer *DepContainer) GetMongoDbContext() (*mongo.Database){
	return depContainer.mongo.Database("packform")
}

func (depContainer *DepContainer) GetMongoClient() (*mongo.Client){
	return depContainer.mongo
}




