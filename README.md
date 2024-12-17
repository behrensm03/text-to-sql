# text-to-sql

Running the app: `go run cmd/app.go` or using the 'Debug' launch configuration.

If you make changes to the database structure, you will have to run the 'Debug Destroy DB' launch configuration, or run `go run cmd/rebuild/rebuild.go`

Sample queries:
- Show me orders by john smith
- Show me all customers
- Show me the most recent order by john smith
- How many orders has each customer placed
- Show me the most expensive product
- Which product has the most orders
- Which customer has spent the most total money
