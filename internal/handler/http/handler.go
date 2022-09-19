package http

import(
	"log"
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"

	"github.com/go-rest-api/internal/adapter/contract"
	"github.com/go-rest-api/internal/model"
)

type HttpBalanceAdapter struct {
	metrics 	contract.MetricsServiceAdapterPort
	repository 	contract.BalanceRepositoryAdapterPort
	service 	contract.BalanceServiceAdapterPort
}

func NewBalanceHttpAdapter(metrics contract.MetricsServiceAdapterPort ,service contract.BalanceServiceAdapterPort, repository contract.BalanceRepositoryAdapterPort) *HttpBalanceAdapter {
	return &HttpBalanceAdapter{
		metrics: metrics,
		repository: repository,
		service: service,
	}
}

func (h *HttpBalanceAdapter) ListBalance(rw http.ResponseWriter, req *http.Request) {
	log.Printf("get/balance/list")
	rw.Header().Set("Content-Type", "application/json")

	result, err := h.service.ListBalance()
	if err != nil{
		json.NewEncoder(rw).Encode(err.Error())
		return
	}

	json.NewEncoder(rw).Encode(result)
	return
}

func (h *HttpBalanceAdapter) ListBalanceById(rw http.ResponseWriter, req *http.Request) {
	log.Printf("get/balance/list")
	rw.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)

	result, err := h.service.ListBalanceById(vars["id"],vars["sk"])
	if err != nil{
		json.NewEncoder(rw).Encode(err.Error())
		return
	}

	json.NewEncoder(rw).Encode(result)
	return
}

func (h *HttpBalanceAdapter) GetBalance(rw http.ResponseWriter, req *http.Request) {
	log.Printf("get/balance/{id}")
	rw.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)

	result, err := h.service.GetBalance(vars["id"])
	if err != nil{
		rw.WriteHeader(http.StatusNotFound)
		json.NewEncoder(rw).Encode(err.Error())
		return
	}

	json.NewEncoder(rw).Encode(result)
	return
}

func (h *HttpBalanceAdapter) AddBalance(rw http.ResponseWriter, req *http.Request) {
	log.Printf("/balance/save")
	rw.Header().Set("Content-Type", "application/json")

	balance := model.Balance{}
	err := json.NewDecoder(req.Body).Decode(&balance)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusBadRequest)
        return
    }

	res, err := h.service.AddBalance(balance)
	if err != nil{
		rw.WriteHeader(http.StatusNotFound)
		json.NewEncoder(rw).Encode(err.Error())
		return
	}

	json.NewEncoder(rw).Encode(res)
	return
}

func (h *HttpBalanceAdapter) Health(rw http.ResponseWriter, req *http.Request) {
	log.Printf("/health")
	rw.Header().Set("Content-Type", "application/json")

	res := h.metrics.Health()
	health := model.ManagerHealthDB{ Status: res }
	if (res == false){
		rw.WriteHeader(http.StatusBadRequest)
	}
	
	json.NewEncoder(rw).Encode(health)
	return
}

var x = 0
func (h *HttpBalanceAdapter) GetCount(rw http.ResponseWriter, req *http.Request) {
	log.Printf("/count/{id}")
	rw.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	x++

	count := model.Count{Id: vars["id"] , Count: strconv.Itoa(x)}

	json.NewEncoder(rw).Encode(count)
	return
}

func (h *HttpBalanceAdapter) StressCPU(rw http.ResponseWriter, req *http.Request) {
	log.Printf("post/stressCPU")
	rw.Header().Set("Content-Type", "application/json")
	
	setup := model.Setup{}
	err := json.NewDecoder(req.Body).Decode(&setup)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusBadRequest)
        return
    }
	res := h.metrics.StressCPU(setup.Count)

	json.NewEncoder(rw).Encode(res)
	return
}

// no longer use !!!!
/*func (h *HttpBalanceAdapter) SetUp(rw http.ResponseWriter, req *http.Request) {
	log.Printf("post/setup")
	rw.Header().Set("Content-Type", "application/json")
	
	setup := model.Setup{}
	err := json.NewDecoder(req.Body).Decode(&setup)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusBadRequest)
        return
    }

	json.NewEncoder(rw).Encode(setup)
	return
}*/