package routes

import (
	"github.com/gin-gonic/gin"
	"jamo/backend/internal/adapter/api"
	"jamo/backend/internal/core/helper"
	"jamo/backend/internal/core/middleware"
	"jamo/backend/internal/core/services"
	ports "jamo/backend/internal/port"
)

func SetupRouter(repository ports.Repository) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware)
	backendService := services.NewService(repository)

	handler := api.NewHTTPHandler(backendService)

	helper.LogEvent("INFO", "Configuring Routes!")

	//router.Use(middleware.SetHeaders)

	router.Group("/product")
	{
		router.GET("/product/amount/:amount", handler.GetProduct)
		router.POST("/product", handler.CreateProduct)
		router.GET("/product/ref/:ref", handler.GetProductById)
		router.POST("/product/cart-items", handler.GetCartItems)
		router.POST("/product/order", handler.CreateOrder)
		router.PUT("/product/order/update-payment/:id", handler.UpdateOrderPayment)
	}

	router.Group("/customer")
	{
		router.POST("/customer/subscribe", handler.SubscribeToNewLetter)
		router.POST("/customer/send-contact-mail", handler.ContactAdmin)
	}

	router.Group("/admin")
	{
		router.GET("/admin/orders/page/:page", handler.GetOrders)
		router.GET("/admin/order/get_dashboard_values", handler.GetDashBoardValues)
		router.GET("/admin/messages/page/:page", handler.GetAdminMsgs)
		router.GET("/admin/order/:id", handler.GetOrderById)
		router.PUT("/admin/order/update-delivery/:id", handler.UpdateDeliveryStatus)
	}

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{"error": "matching no route error"})
	})
	return router
}
