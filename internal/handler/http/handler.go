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
	log.Printf("/balance/list")
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
	log.Printf("/balance/list")
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
	log.Printf("/balance/{id}")
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
		log.Panic(err)
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

func (h *HttpBalanceAdapter) UpdateBalance(rw http.ResponseWriter, req *http.Request) {
	log.Printf("/balance/update")
	rw.Header().Set("Content-Type", "application/json")

	balance := model.Balance{}
	err := json.NewDecoder(req.Body).Decode(&balance)
    if err != nil {
		log.Panic(err)
        http.Error(rw, err.Error(), http.StatusBadRequest)
        return
    }

	res, err := h.service.UpdateBalance(balance)
	if err != nil{
		rw.WriteHeader(http.StatusNotFound)
		json.NewEncoder(rw).Encode(err.Error())
		return
	}

	json.NewEncoder(rw).Encode(res)
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
	log.Printf("/stressCPU")
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

func (h *HttpBalanceAdapter) ShowHeader(rw http.ResponseWriter, req *http.Request) {
	log.Printf("/header")
	
 	/*result := ""

	if reqHeadersBytes, err := json.Marshal(req.Header); err != nil {
		log.Println("Could not Marshal http headers")
		result = "Could not Marshal http headers"
	} else {
		log.Println(string(reqHeadersBytes))
		result = string(reqHeadersBytes)
	}

	responseBody := `{"textfield":"I'm a text.","num":1234,"list":[1,2,3]}`
    var data map[string]interface{}
    err := json.Unmarshal([]byte(responseBody), &data)
    if err != nil {
        panic(err)
    }
    log.Println(data["list"])
	log.Println(result)*/

	json.NewEncoder(rw).Encode(req.Header)
	return
}
