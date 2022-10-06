package http

import (
	"log"
	"net/http"
	"encoding/json"

	"github.com/go-rest-api/internal/error"
)

func MiddleWareHandlerToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("-------------------------------------------- \n")
		log.Println("MiddleWareHandlerToken (INICIO)")
		log.Println(r.Header.Get("jwt"))

		if len(r.Header.Get("jwt")) == 0 {
			erro.HandlerHttpError(w, erro.ErrUnauthorized)
			json.NewEncoder(w).Encode("erro.ErrUnauthorized")
			return
		}
	
		log.Println("MiddleWareHandlerToken (FIM)")
		log.Printf("-------------------------------------------- \n")
		
		next.ServeHTTP(w, r)
	})
}

func MiddleWareHandlerHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("-------------------------------------------- \n")
		log.Println("MiddleWareHandlerHeader (INICIO)")
		
		if reqHeadersBytes, err := json.Marshal(r.Header); err != nil {
			log.Println("Could not Marshal hhtp headers")
		} else {
			log.Println(string(reqHeadersBytes))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers","Content-Type,access-control-allow-origin, access-control-allow-headers")
	
		//log.Println(r.Header.Get("Host"))
		//log.Println(r.Header.Get("User-Agent"))
		//log.Println(r.Header.Get("X-Forwarded-For"))

		log.Println("MiddleWareHandlerHeader (FIM)")
		log.Printf("-------------------------------------------- \n")
		
		next.ServeHTTP(w, r)
	})
}