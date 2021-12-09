package main

import (
	"github.com/Baja-KS/Webshop-API/AuthenticationService/internal/database"
	"github.com/Baja-KS/Webshop-API/AuthenticationService/internal/service"
	"github.com/Baja-KS/Webshop-API/AuthenticationService/internal/service/endpoints"
	"github.com/Baja-KS/Webshop-API/AuthenticationService/internal/service/middlewares"
	"github.com/Baja-KS/Webshop-API/AuthenticationService/internal/service/transport"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/rs/cors"
	"net/http"
	"os"
)

func main() {
	logger:=log.NewLogfmtLogger(os.Stderr)
	var err error
	db,err:=database.NewDatabase()
	if err != nil {
		log.With(logger,"err",err)
	}
	err=database.Migrate(db)
	if err != nil {
		log.With(logger,"err",err)
	}
	requestCount:=kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace:   "bajaks_webshop_api",
		Subsystem:   "authentication_service",
		Name:        "request_count",
		Help:        "Number of request received",
	},[]string{"method","error"})
	requestLatency:=kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace:   "bajaks_webshop_api",
		Subsystem:   "authentication_service",
		Name:        "request_latency",
		Help:        "Total duration of requests in microseconds",
	},[]string{"method","error"})

	//svc:=service.AuthenticationService{DB: db}
	//svc=middlewares.LoggingMiddleware{}
	var svc service.Service
	svc= &service.AuthenticationService{DB: db}
	svc= &middlewares.AuthorizationMiddleware{DB: db,Next: svc}
	svc= &middlewares.AuthenticationMiddleware{DB: db,Next: svc}
	svc= &middlewares.LoggingMiddleware{Logger: logger, Next: svc}
	svc= &middlewares.InstrumentingMiddleware{
		RequestCount: requestCount,
		RequestLatency: requestLatency,
		Next: svc,
	}

	ep:=endpoints.NewEndpointSet(svc)
	err = http.ListenAndServe(":8080",cors.AllowAll().Handler(transport.NewHTTPHandler(ep)))
	if err != nil {
		log.With(logger,"err",err)
	}
	log.With(logger,"msg","Listening to port")
}
