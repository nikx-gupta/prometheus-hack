package main

import (
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"prometheus-hack/postgresapi"
	"prometheus-hack/sqlapi"
)

func main() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.Handle("/pg/connect", postgresapi.ConnectHandler(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Connected"))
	}))

	http.Handle("/customer/get", postgresapi.GetCustomer())
	http.Handle("/sql/customer/get", sqlapi.GetStudent())


	fmt.Println("Starting listening on 2112")
	err := http.ListenAndServe(":2112", nil)
	if err != nil {
		panic(err)
	}
}
