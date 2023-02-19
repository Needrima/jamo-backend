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

	productsRouter := router.Group("/product")
	{
		productsRouter.GET("/amount/:amount", handler.GetProduct)
		productsRouter.POST("", handler.CreateProduct)
		productsRouter.GET("/ref/:ref", handler.GetProductById)
		productsRouter.POST("/cart-items", handler.GetCartItems)
		productsRouter.POST("/order", handler.CreateOrder)
		productsRouter.PUT("/order/update-payment/:id", handler.UpdateOrderPayment)
	}

	customersRouter := router.Group("/customer")
	{
		customersRouter.POST("/subscribe", handler.SubscribeToNewLetter)
		customersRouter.POST("/send-contact-mail", handler.ContactAdmin)
	}

	adminRouter := router.Group("/admin")
	{
		adminRouter.GET("/orders/page/:page", handler.GetOrders)
		adminRouter.GET("/order/get_dashboard_values", handler.GetDashBoardValues)
		adminRouter.GET("/messages/page/:page", handler.GetAdminMsgs)
		adminRouter.GET("/order/:id", handler.GetOrderById)
		adminRouter.PUT("/order/update-delivery/:id", handler.UpdateDeliveryStatus)
	}

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{"error": "matching no route error"})
	})
	return router
}
