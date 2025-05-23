package http

import (
	"be-lotsanmateo-api/internal/adapter/http/handler/api"
	"be-lotsanmateo-api/internal/config"
	"be-lotsanmateo-api/pkg"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(env *config.Env) *gin.Engine {
	gin.ForceConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), cors.Default())

	test := api.NewTestHandler()

	customer := pkg.NewCustomerHandler(env)

	apiRest := r.Group("/api")
	{
		apiRest.GET("/test", test.TestHandler)
		apiRest.GET("/payment/plan/simulation")
		apiRest.POST("/v1/customer-onboarding", customer.CreateCustomer)
	}

	pdfReport := r.Group("/pdf")
	{
		pdfReport.GET("/payment/plan/simulation")
	}

	return r
}
