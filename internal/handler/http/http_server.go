package http

import (
	"time"
	"log"
	"net/http"
	"strconv"
	"syscall"
	"os"
	"os/signal"
	"context"
	_ "net/http/pprof"
	"github.com/gorilla/mux"
	"encoding/json"

	"github.com/go-rest-api/internal/model"
)

type DebugServer struct {
	*http.Server
}

func NewDebugServer(address string) *DebugServer {
	return &DebugServer{
		&http.Server{
			Addr:    address,
			Handler: http.DefaultServeMux,
		},
	}
}

type HttpServer struct {
	start time.Time
	http_server_setup model.ManagerInfo
}

func NewHttpServer(start time.Time, http_server_setup model.ManagerInfo) HttpServer {
	return HttpServer{	start: start, 
						http_server_setup: http_server_setup,
					}
}

func (s HttpServer) StartHttpServer(handler_balance *HttpBalanceAdapter) {
	duration := time.Since(s.start).Nanoseconds()
	log.Print("Server HTTP started", duration)

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/info", func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(s.http_server_setup)
	})

	list_balance := myRouter.Methods(http.MethodGet).Subrouter()
    list_balance.HandleFunc("/balance/list", handler_balance.ListBalance)
	//list_balance.Use(MiddleWareHandler)

	list_balance_id := myRouter.Methods(http.MethodGet).Subrouter()
    list_balance_id.HandleFunc("/balance/list_by_id/{id}&{sk}", handler_balance.ListBalanceById) 

	get_balance := myRouter.Methods(http.MethodGet).Subrouter()
    get_balance.HandleFunc("/balance/{id}", handler_balance.GetBalance) 

	add_balance := myRouter.Methods(http.MethodPost).Subrouter()
    add_balance.HandleFunc("/balance/save", handler_balance.AddBalance)

	health := myRouter.Methods(http.MethodGet).Subrouter()
    health.HandleFunc("/health", handler_balance.Health)

	get_count := myRouter.Methods(http.MethodGet).Subrouter()
    get_count.HandleFunc("/count/{id}", handler_balance.GetCount) 
	
	cpu_stress := myRouter.Methods(http.MethodPost).Subrouter()
    cpu_stress.HandleFunc("/stress/cpu", handler_balance.StressCPU)

	setup := myRouter.Methods(http.MethodPost).Subrouter()
    setup.HandleFunc("/setup", handler_balance.SetUp)

	srv := http.Server{
		Addr:         ":" +  strconv.Itoa(s.http_server_setup.Server.Port),      	
		Handler:      myRouter,                	          
		ReadTimeout:  time.Duration(s.http_server_setup.Server.ReadTimeout) * time.Second,   
		WriteTimeout: time.Duration(s.http_server_setup.Server.WriteTimeout) * time.Second,  
		IdleTimeout:  time.Duration(s.http_server_setup.Server.IdleTimeout) * time.Second, 
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Print("Internal error", err)
		}
	}()

	debugServer := NewDebugServer("127.0.0.1:6060")
	go func() {
		log.Print("Starting Server! http://localhost:6060/debug/pprof/ ")
		err := debugServer.ListenAndServe()
		if err != nil {
			log.Print("PPROF Internal error", err)
		}
	}()
	
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch

	log.Printf("Stopping Server")
	ctx , cancel := context.WithTimeout(context.Background(), time.Duration(s.http_server_setup.Server.CtxTimeout) * time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		log.Print("WARNING Dirty Shutdown", err)
		return
	}
	log.Printf("Stop DONE !")
}