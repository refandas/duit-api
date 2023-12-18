package app

import (
	"github.com/julienschmidt/httprouter"
	"github.com/refandas/duit-api/controller"
	"github.com/refandas/duit-api/exception"
)

// Router is a struct representing an HTTP router and associated controllers.
type Router struct {

	// UserController represents the controller for user-related functionality.
	UserController controller.UserController

	// SpendingController represents the controller for user's spending-related functionality.
	SpendingController controller.SpendingController
}

// NewRouter creates and returns a new instance of httprouter.Router
//
// The NewRouter function is used to define and handle routes based
// on the defined controller.
func (controller Router) NewRouter() *httprouter.Router {
	router := httprouter.New()

	// The user handler will only be defined if the UserController is defined.
	if controller.UserController != nil {
		router.GET("/api/v1/users/:userId", controller.UserController.FindById)
		router.PUT("/api/v1/users/:userId", controller.UserController.Update)
		router.POST("/api/v1/users", controller.UserController.Create)
		router.DELETE("/api/v1/users/:userId", controller.UserController.Delete)
	}

	// The user's spending handler will only be defined if the SpendingController is defined.
	if controller.SpendingController != nil {
		router.GET("/api/v1/users/:userId/spendings", controller.SpendingController.FindByUserId)
		router.GET("/api/v1/spendings/:spendingId", controller.SpendingController.FindById)
		router.PUT("/api/v1/spendings/:spendingId", controller.SpendingController.Update)
		router.POST("/api/v1/spendings", controller.SpendingController.Create)
		router.DELETE("/api/v1/spendings/:spendingId", controller.SpendingController.Delete)
	}

	// Setting an error handler when panic occurs.
	router.PanicHandler = exception.ErrorHandler

	return router
}
