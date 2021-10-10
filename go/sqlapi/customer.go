package sqlapi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
)

type Student struct {
	StudentId uuid.UUID    `json:"studentId"`
	Name      sql.NullString `json:"name"`
	Email     sql.NullString `json:"email"`
}

type StudentRepo struct {
	db     *sql.DB
	output io.Writer
}

func (c *StudentRepo) Get(pageIndex int, pageSize int) {
	query := AddPaging("select StudentId,Name,Email from dbo.Student order by Name", pageIndex, pageSize)
	data, err := c.db.Query(query)
	if err != nil {
		fmt.Println("Error in Querying Db")
		panic(err)
	}

	defer data.Close()
	var customers []Student
	for data.Next() {
		var cust Student

		err := data.Scan(&cust.StudentId, &cust.Name, &cust.Email)
		if err != nil {
			fmt.Println("Error in getting records")
			panic(err)
		}

		customers = append(customers, cust)
	}

	json.NewEncoder(c.output).Encode(customers)
}

func AddPaging(s string, pageIndex int, pageSize int) string {
	if pageSize <= 0 {
		pageSize = 10
	}

	if pageIndex <= 0 {
		pageIndex = 1
	}

	query := fmt.Sprintf("%s OFFSET %d rows fetch next %d rows only", s, (pageIndex-1)*pageSize, pageSize)
	fmt.Printf("Query: %s", query)
	return query
}
