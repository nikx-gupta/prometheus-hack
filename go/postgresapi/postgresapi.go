package postgresapi

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

var db *sql.DB

var(selectTime = promauto.NewGauge(prometheus.GaugeOpts{
	
}))

func Run() {
	http.Handle("/metrics", promhttp.Handler())

	http.Handle("/connect", ConnectHandler(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Connected"))
	}))

	http.Handle("/customer/get", LogTimeHandler(ConnectHandler(func(w http.ResponseWriter, r *http.Request) {
		repo := &CustomerRepo{db: db, output: w}
		err := repo.Get(-1, -1)
		if err != nil {
			panic(err)
		}
	})))

	fmt.Println("Starting listening on 2112")
	err := http.ListenAndServe(":2112", nil)
	if err != nil {
		panic(err)
	}
}

func LogTimeHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("id", uuid.New().String())
		w.Header().Set("id",r.Header.Get("id"))
		st := time.Now()
		next.ServeHTTP(w, r)
		et := time.Now()
		fmt.Printf("Time For Request: %s\n", et.Sub(st))

	})
}

func ConnectHandler(handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if db != nil {
			handler(w, r)
			return
		}

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

		handler(w, r)
	})
}
