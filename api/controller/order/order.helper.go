package order

import (
	// "net/http"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"sync"
	"packform/utils/injector"
)



func AddOrdersToDb(list chan[]string, wg *sync.WaitGroup,context *injector.DepContainer) error{
	defer wg.Done()

	db := context.GetDbContext() 
	query := `INSERT INTO orders(created_at,order_name,customer_id) VALUES($1,$2,$3)`

	tx, err := db.Begin()
	txStmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}

	for item := range list {
		var datetime string = item[1]
		var order string = item[2]
		var customer string = item[3]
		_ , err = txStmt.Exec(datetime, order, customer)
		if err != nil {
			tx.Rollback()
			fmt.Println(err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

/** 
	Close channel here to prevent further queueing more data  
**/
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