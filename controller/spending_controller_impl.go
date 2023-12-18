package controller

import (
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/refandas/duit-api/helper"
	"github.com/refandas/duit-api/model/web"
	"github.com/refandas/duit-api/service"
	"net/http"
	"time"
)

type SpendingControllerImpl struct {
	SpendingService service.SpendingService
}

func NewSpendingController(spendingService service.SpendingService) SpendingController {
	return &SpendingControllerImpl{SpendingService: spendingService}
}

func (controller *SpendingControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	spendingCreateRequest := web.SpendingCreateRequest{}
	helper.ReadFromRequestBody(request, &spendingCreateRequest)

	spendingId, _ := uuid.NewRandom()
	spendingCreateRequest.Id = spendingId.String()
	spendingCreateRequest.CreatedAt = time.Now().UnixMilli()

	spendingResponse := controller.SpendingService.Create(request.Context(), spendingCreateRequest)
	webResponse := web.WebResponse{
		Code:   http.StatusCreated,
		Status: "CREATED",
		Data:   spendingResponse,
	}
	helper.WriteToResponseBody(writer, &webResponse)
}

func (controller *SpendingControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	spendingUpdateRequest := web.SpendingUpdateRequest{}
	helper.ReadFromRequestBody(request, &spendingUpdateRequest)

	spendingId := params.ByName("spendingId")
	spendingUpdateRequest.Id = spendingId

	spendingResponse := controller.SpendingService.Update(request.Context(), spendingUpdateRequest)
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   spendingResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *SpendingControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	spendingId := params.ByName("spendingId")

	controller.SpendingService.Delete(request.Context(), spendingId)
	webResponse := web.WebResponse{
		Code:   http.StatusNoContent,
		Status: "DELETED",
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *SpendingControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	spendingId := params.ByName("spendingId")

	spendingResponse := controller.SpendingService.FindById(request.Context(), spendingId)
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   spendingResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *SpendingControllerImpl) FindByUserId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userId := params.ByName("userId")

	spendingResponse := controller.SpendingService.FindByUserId(request.Context(), userId)
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   spendingResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
