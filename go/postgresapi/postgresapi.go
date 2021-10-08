package postgresapi

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var db *sql.DB
var (
	db_connection_time_gauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "postgresapi_db_connection_time",
		Help: "Db Connection time taken",
	})
)
var (
	customer10_fetch_time_gauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "postgresapi_customer10_fetch_time",
		Help: "Fetching 10 Customers Time",
	})
)

func Run() {
	http.Handle("/metrics", promhttp.Handler())

	http.Handle("/connect", ConnectHandler(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Connected"))
	}))

	http.Handle("/customer/get", ConnectHandler(func(w http.ResponseWriter, r *http.Request) {
		tm := TimeOperation("Fetch 10 Customers", func() {
			repo := &CustomerRepo{db: db, output: w}
			err := repo.Get(1, 10)
			if err != nil {
				panic(err)
			}
		})

		customer10_fetch_time_gauge.Set(tm.Seconds())
	}))

	fmt.Println("Starting listening on 2112")
	err := http.ListenAndServe(":2112", nil)
	if err != nil {
		panic(err)
	}
}

func ConnectHandler(handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if db != nil {
			handler(w, r)
			return
		}

		tm := TimeOperation("Opening Connection", func() {
			psqlconn := fmt.Sprintf("host=localhost user=nikx password=demo@123 dbname=postgres sslmode=disable")
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

		db_connection_time_gauge.Set(tm.Seconds())
		handler(w, r)
	})
}

func TimeOperation(msg string, targetFunc func()) time.Duration {
	st := time.Now()
	targetFunc()
	et := time.Now().Sub(st)
	fmt.Printf("%s: Time: %s\n", msg, et)

	return et
}
