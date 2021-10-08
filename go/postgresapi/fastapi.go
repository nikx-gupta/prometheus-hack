package postgresapi

import (
	"database/sql"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func RunFastHttp() {
	// http.Handle("/metrics", promhttp.Handler())

	//http.Handle("/connect", ConnectHandler(func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("Connected"))
	//}))
	r := router.New()
	r.GET("/", func(c *fasthttp.RequestCtx) {
		fmt.Fprintf(c, "Hello, world!")
	})

	r.GET("/customer/get",FastLogTimeHandler(FatConnectHandler(func(c *fasthttp.RequestCtx) {
		repo := &CustomerRepo{db: db, output: c.Response.BodyWriter()}
		err := repo.Get(-1, -1)
		if err != nil {
			panic(err)
		}
	})))

	fmt.Println("Starting listening on 2112")
	err := 	fasthttp.ListenAndServe(":2112", r.Handler)
	if err != nil {
		panic(err)
	}
}

func FastLogTimeHandler(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func (ctx *fasthttp.RequestCtx) {
		uuid:=uuid.New().String()
		ctx.Request.Header.Set("id", uuid)
		ctx.Response.Header.Set("id",  uuid)
		st := time.Now()
		next(ctx)
		et := time.Now()
		fmt.Printf("Time For Request: %s\n", et.Sub(st))
	}
}

func FatConnectHandler(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	if db != nil {
		return handler
	}

	return func(ctx *fasthttp.RequestCtx) {
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

		handler(ctx)
	}
}
