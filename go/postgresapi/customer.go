package postgresapi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
)

type Customer struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Category  string `json:"category"`
}

type CustomerRepo struct {
	db *sql.DB
	output io.Writer
}

func (c *CustomerRepo) Get(pageIndex int, pageSize int) error {
	query := AddPaging("select id,category,first_name,last_name from db_customer", pageIndex, pageSize)
	data, err := c.db.Query(query)
	if err != nil {
		fmt.Println("Error in Querying Db")
		panic(err)
	}

	defer data.Close()
	var customers []Customer
	for data.Next() {
		var cust Customer
		err := data.Scan(&cust.Id, &cust.Category, &cust.FirstName, &cust.LastName)
		if err != nil {
			fmt.Println("Error in getting records")
			panic(err)
		}

		customers = append(customers, cust)
	}

	json.NewEncoder(c.output).Encode(customers)

	return nil
}

func AddPaging(s string, pageIndex int, pageSize int) string {
	if  pageSize <= 0 {
		pageSize = 10
	}

	if pageIndex <= 0 {
		pageIndex = 1
	}

	return fmt.Sprintf("%s OFFSET %d rows fetch next %d rows only",s, (pageIndex - 1) * pageSize, pageSize)
}
