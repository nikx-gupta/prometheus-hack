package postgresapi

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var db *sql.DB
var (
	db_connection_time_gauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "postgresapi_db_connection_time",
		Help: "Db Connection time taken",
	})
)
var (
	customer_fetch_time_gauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "postgresapi_customer_fetch_sec",
		Help: "Fetching 10 Customers Time",
	}, []string{"page_index", "page_size"})
)

func GetCustomer() http.Handler{
	return ConnectHandler(func(w http.ResponseWriter, r *http.Request) {
		pi, _ := strconv.Atoi(r.URL.Query().Get("pi"))
		ps, _ := strconv.Atoi(r.URL.Query().Get("ps"))
		tm := TimeOperation(fmt.Sprintf("Fetch %d Customers", pi * ps), func() {
			repo := &CustomerRepo{db: db, output: w}
			err := repo.Get(pi, ps)
			if err != nil {
				panic(err)
			}
		})

		customer_fetch_time_gauge.WithLabelValues(fmt.Sprint(pi), fmt.Sprint(ps)).Set(tm.Seconds())
	})
}

func ConnectHandler(handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if db != nil {
			handler(w, r)
			return
		}

		tm := TimeOperation("Opening Connection", func() {
			var psqlconn string
			if psqlconn = os.Getenv("postgres_conn"); psqlconn == "" {
				psqlconn = "host=10.1.232.147 user=postgres password=demo@123 dbname=postgres sslmode=disable"
			}

			fmt.Printf("Connection: %s", psqlconn)
			var err error
			db, err = sql.Open("postgres", psqlconn)
			if err != nil {
				fmt.Println("Error in opening Connection")
				panic(err)
			}
			err = db.Ping()
			if err != nil {
				fmt.Println("Error in Ping")
				panic(err)
			}
		})

		db_connection_time_gauge.Set(float64(tm.Milliseconds()))
		handler(w, r)
	})
}

func TimeOperation(msg string, targetFunc func()) time.Duration {
	st := time.Now()
	targetFunc()
	et := time.Since(st)
	fmt.Printf("%s: Time: %s\n", msg, et)

	return et
}
