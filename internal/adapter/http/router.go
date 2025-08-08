package http

import (
	"be-lotsanmateo-api/internal/adapter/http/handler/api"
	"be-lotsanmateo-api/internal/adapter/http/handler/middleware"
	"be-lotsanmateo-api/internal/config"
	"be-lotsanmateo-api/pkg"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(env *config.Env) *gin.Engine {
	gin.ForceConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	cors.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AddAllowHeaders("x-request-id")
	corsConfig.AddAllowHeaders("x-language")
	corsConfig.AddAllowHeaders("x-user")
	corsConfig.AllowAllOrigins = true
	defaultConfig := cors.New(corsConfig)

	authMiddleware := middleware.NewAuthMiddleware(env)

	r.Use(gin.Logger(), gin.Recovery(), defaultConfig, authMiddleware.AuthMiddleware())

	test := api.NewTestHandler()

	customer := pkg.NewCustomerHandler(env)
	paymentPlan := pkg.NewPaymentPlanHandler()
	loanPaymentHandler := pkg.NewLoanPaymentHandler(env)
	activeFinancing := pkg.NewFinancingHandler(env)

	apiRest := r.Group("/api")
	{
		apiRest.GET("/test", test.TestHandler)
		apiRest.POST("/payment/plan/simulation", paymentPlan.HandleRequest)
		apiRest.POST("/v1/loan", customer.CreateCustomer)
		apiRest.GET("/v1/loan/payment", loanPaymentHandler.HandleRequest)
		apiRest.POST("/v1/customer-onboarding", customer.CreateCustomer)
		apiRest.POST("/v1/financing/active", activeFinancing.HandleRequest)
	}

	paymentPlanPdf := pkg.NewCalculatePlanPdfHandler(env)
	promissoryNoteHandler := pkg.NewPromissoryNoteHandler(env)
	invoiceHandler := pkg.NewInvoiceHandler(env)

	pdfReport := r.Group("/pdf")
	{
		pdfReport.GET("/payment/plan/simulation", paymentPlanPdf.GeneratePDF)
		pdfReport.GET("/promissory/note", promissoryNoteHandler.GeneratePDF)
		pdfReport.GET("/invoice/payment", invoiceHandler.GeneratePayment)
		pdfReport.GET("/invoice/down/payment", invoiceHandler.GenerateDownPayment)
		pdfReport.GET("/invoice/reservation", invoiceHandler.GenerateReservation)
	}

	return r
}
