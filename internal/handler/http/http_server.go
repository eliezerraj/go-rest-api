package http

import (
	"time"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"syscall"
	"os"
	"os/signal"
	"context"
	_ "net/http/pprof"
	"encoding/json"

	"github.com/gorilla/mux"

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
	log.Print("Server HTTP started v.2", duration)

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/balance/test", func(rw http.ResponseWriter, req *http.Request) {
		log.Printf("/balance/test")
		
		log.Printf("%s %s", req.Method, req.URL.Path)

		if reqBytes, err := json.Marshal(req.Header); err != nil {
			log.Println("Could not Marshal header")
		} else {
			log.Print(string(reqBytes))
		}

		var payload interface{} 
		buffer, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Panic(err)
		}
		req.Body.Close()
		json.Unmarshal(buffer, &payload)
		m := payload.(map[string]interface{})
		log.Println("body : ",m)

		rw.Header().Set("Content-Type", "application/json")
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Headers","Content-Type,access-control-allow-origin, access-control-allow-headers")
		json.NewEncoder(rw).Encode("OK-tezt")
	})

	//---------------------------------------------------------------------
	show_header := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    show_header.HandleFunc("/header", handler_balance.ShowHeader)
	show_header.Use(MiddleWareHandlerHeader)

	myRouter.HandleFunc("/info", func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Headers","Content-Type,access-control-allow-origin, access-control-allow-headers")
		log.Printf("/info")
		json.NewEncoder(rw).Encode(s.http_server_setup)
	})

	health := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    health.HandleFunc("/health", handler_balance.Health)
	health.Use(MiddleWareHandlerHeader)
	
	myRouter.HandleFunc("/live", func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		log.Printf("get/live")

		liveness := model.ManagerHealth{ Liveness: s.http_server_setup.Setup.Liveness }
		if (!s.http_server_setup.Setup.Liveness){
			rw.WriteHeader(http.StatusBadRequest)
		}
		json.NewEncoder(rw).Encode(liveness.Liveness)
	})

	myRouter.HandleFunc("/ready", func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		log.Printf("get/ready")

		readiness := model.ManagerHealth{ Readiness: s.http_server_setup.Setup.Readiness }
		if (!s.http_server_setup.Setup.Readiness){
			rw.WriteHeader(http.StatusBadRequest)
		}

		res, err := handler_balance.repository.Ping()
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
		}
		if (!res) {
			rw.WriteHeader(http.StatusBadRequest)
		}
		//log.Printf("get/ready > handler_balance.repository", res)

		json.NewEncoder(rw).Encode(readiness.Readiness)
	})

	myRouter.HandleFunc("/setup", func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		log.Printf("post/setup")

		setup := model.Setup{}
		err := json.NewDecoder(req.Body).Decode(&setup)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		s.http_server_setup.Setup.Liveness = setup.Liveness
		s.http_server_setup.Setup.Readiness = setup.Readiness

		json.NewEncoder(rw).Encode(s.http_server_setup)
	})

	//---------------------------------------------------------------------
	list_balance := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    list_balance.HandleFunc("/balance/list", handler_balance.ListBalance)
	list_balance.Use(MiddleWareHandlerHeader)
	//list_balance.Use(MiddleWareHandlerToken)

	list_balance_id := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    list_balance_id.HandleFunc("/balance/list_by_id/{id}&{sk}", handler_balance.ListBalanceById) 
	list_balance_id.Use(MiddleWareHandlerHeader)

	get_balance := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    get_balance.HandleFunc("/balance/{id}", handler_balance.GetBalance) 
	get_balance.Use(MiddleWareHandlerHeader)

	add_balance := myRouter.Methods(http.MethodPost, http.MethodOptions).Subrouter()
    add_balance.HandleFunc("/balance/save", handler_balance.AddBalance)
	add_balance.Use(MiddleWareHandlerHeader)

	update_balance := myRouter.Methods(http.MethodPost, http.MethodOptions).Subrouter()
    update_balance.HandleFunc("/balance/update", handler_balance.UpdateBalance)
	update_balance.Use(MiddleWareHandlerHeader)

	//---------------------------------------------------------------------

	get_count := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    get_count.HandleFunc("/count/{id}", handler_balance.GetCount) 
	get_count.Use(MiddleWareHandlerHeader)

	cpu_stress := myRouter.Methods(http.MethodPost, http.MethodOptions).Subrouter()
    cpu_stress.HandleFunc("/stress/cpu", handler_balance.StressCPU)
	cpu_stress.Use(MiddleWareHandlerHeader)

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