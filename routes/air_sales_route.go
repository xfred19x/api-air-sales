package routes

import (
	"api-air-sales/controllers" //add this
	"github.com/gin-gonic/gin"
)

const BasePath = "air-sales"

// SetupRoutes configura todas las rutas para la aplicación Air Sales
func SetupRoutes(router *gin.Engine) {
	v1 := router.Group(BasePath)

	routesGroup := v1.Group("/routes")
	{
		routesGroup.POST("/initiate", controllers.CreateRoute())
		routesGroup.GET("/retrieve", controllers.GetRouteFilters())
		routesGroup.GET("/retrieve/:routeId", controllers.GetRouteId())
	}

	reservationsGroup := v1.Group("/reservations")
	{
		reservationsGroup.POST("/initiate", controllers.CreateReservation())
		reservationsGroup.GET("/retrieve/:reserveId", controllers.GetReservationsId())
		reservationsGroup.PUT("/update/:reserveId", controllers.CancelReservation())
	}

	paymentsGroup := v1.Group("/payments")
	{
		paymentsGroup.POST("/initiate", controllers.CreatePayment())
		paymentsGroup.GET("/retrieve/:paymentId", controllers.GetPaymentsId())
		paymentsGroup.PUT("/update/:paymentId", controllers.CancelPayment())
	}

	ticketsGroup := v1.Group("/tickets")
	{
		ticketsGroup.POST("/initiate", controllers.CreateTicket())
		ticketsGroup.GET("/retrieve/:ticketId", controllers.GetTicketId())
	}
}
