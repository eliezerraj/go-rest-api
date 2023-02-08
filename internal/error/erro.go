package erro

import (
	"errors"
	"net/http"
)

var (
	ErrGetRate = errors.New("Erro no serviço RATE")
	ErrListNotAllowed = errors.New("Lista (SCAN) não permitida para o DynamoDB")
	ErrSaveDatabase = errors.New("Erro no UPSERT")
	ErrOpenDatabase = errors.New("Erro na abertura do DB")
	ErrNotFound = errors.New("Item não encontrado")
	ErrFunctionNotImpl = errors.New("Função não implementada")
	ErrInsert = errors.New("Erro na inserção do dado")
	ErrUnmarshal = errors.New("Erro na conversão do JSON")
	ErrUnauthorized = errors.New("Erro de autorização")
	ErrConvertion = errors.New("Erro de conversão de String para Inteiro")
)

func HandlerHttpError(w http.ResponseWriter, err error) { 
	switch err {
		case ErrUnauthorized:
			w.WriteHeader(http.StatusUnauthorized)	
		default:
			w.WriteHeader(http.StatusInternalServerError)
	}
}