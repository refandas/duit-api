package controller

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type SpendingController interface {
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindByUserId(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
