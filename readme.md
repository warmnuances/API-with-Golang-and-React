Front End Repo : [https://github.com/warmnuances/test_frontend](https://github.com/warmnuances/test_frontend)
Back End Repo : [https://github.com/warmnuances/packform_test](https://github.com/warmnuances/packform_test)


Running the Backend Application: 
1. Set environment variables in `env.template.txt` to `.env`
2. Replace the key values in the `.env` (Note. the " . " is important)
4. go run main.go

Running the Frontend Application: 
1.npm install 
2.npm start

Get All Orders :
`host:  localhost:8080`

`schemes:`
`- http`

# **Paths:**

### **/orders:**
`GET`
`localhost:8080/orders/all?` <br />
summary:  Returns a list of users. <br />
Query String: (All Queries are optional)
1. cust_id : string
2. order_name : string
3. start_date : string `yyyy-mm-ddTHH:mm:12Z`
4. end_date : string `yyyy-mm-ddTHH:mm:12Z`
5. page: number
6. delivered_min:  number
7. delivered_max:  number
8. total_min:  number
9. total_max:  number

<br />

## **Uploading CSV:**
API endpoints are created for uploading specific CSVs. <br/>
Files must be uploaded as :
`POST`
1. Test task - Postgres - orders 
`/orders/bulk` <br/>
2. Test task - Postgres - order_items
`/add/orderitem` <br/>
3. Test task - Postgres - deliveries
`/add/delivery` <br/>
4. Test task - Mongo - customers
`/add/customers` <br/>
5. Test task - Mongo - customer_companies
`/add/companies` <br/>
