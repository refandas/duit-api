package controller

import (
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/refandas/duit-api/helper"
	"github.com/refandas/duit-api/model/web"
	"github.com/refandas/duit-api/service"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type UserControllerImpl struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &UserControllerImpl{UserService: userService}
}

func (controller *UserControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userCreateRequest := web.UserCreateRequest{}
	helper.ReadFromRequestBody(request, &userCreateRequest)

	userId, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	userCreateRequest.Id = userId.String()

	password := []byte(userCreateRequest.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	userCreateRequest.Password = string(hashedPassword)
	userCreateRequest.CreatedAt = time.Now().UnixMilli()

	userResponse := controller.UserService.Create(request.Context(), userCreateRequest)
	webResponse := web.WebResponse{
		Code:   http.StatusCreated,
		Status: "CREATED",
		Data:   userResponse,
	}
	helper.WriteToResponseBody(writer, &webResponse)
}

func (controller *UserControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userUpdateRequest := web.UserUpdateRequest{}
	helper.ReadFromRequestBody(request, &userUpdateRequest)

	userId := params.ByName("userId")
	userUpdateRequest.Id = userId

	password := []byte(userUpdateRequest.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	userUpdateRequest.Password = string(hashedPassword)

	userResponse := controller.UserService.Update(request.Context(), userUpdateRequest)
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   userResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UserControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userId := params.ByName("userId")

	controller.UserService.Delete(request.Context(), userId)
	webResponse := web.WebResponse{
		Code:   http.StatusNoContent,
		Status: "DELETED",
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UserControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userId := params.ByName("userId")

	userResponse := controller.UserService.FindById(request.Context(), userId)
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   userResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}
