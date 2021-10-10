package sqlapi

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
	sql_connection_time_gauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "sql_connection_time",
		Help: "Sql Server Connection time taken",
	})
)
var (
	sql_fetch_time_gauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sql_fetch_sec",
		Help: "Fetching Time",
	}, []string{"table_name", "page_index", "page_size"})
)

func GetStudent() http.Handler {
	return ConnectHandler(func(w http.ResponseWriter, r *http.Request) {
		pi, _ := strconv.Atoi(r.URL.Query().Get("pi"))
		ps, _ := strconv.Atoi(r.URL.Query().Get("ps"))
		tm := TimeOperation(fmt.Sprintf("Fetch %d Student", ps), func() {
			repo := &StudentRepo{db: db, output: w}
			repo.Get(pi, ps)
		})

		sql_fetch_time_gauge.WithLabelValues("student", fmt.Sprint(pi), fmt.Sprint(ps)).Set(tm.Seconds())
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
				psqlconn = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
					"nikxtestdb.database.windows.net", "nikx", "Demo@123", 1433, "testconnlatency")
				// psqlconn = "sqlserver://weberit_apidbuserexam:Exambash%2321@wdb2.my-hosting-panel.com:1433?database=weberit_ebprod&connection+timeout=30"
			}

			fmt.Printf("Connection: %s", psqlconn)
			var err error
			db, err = sql.Open("sqlserver", psqlconn)
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

		sql_connection_time_gauge.Set(float64(tm.Milliseconds()))
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
