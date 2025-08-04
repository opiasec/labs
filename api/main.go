package main

import (
	"appseclabs/config"
	"appseclabs/database"
	"appseclabs/k8s"
	"appseclabs/lab"
	labdefinitions "appseclabs/labDefinitions"
	"appseclabs/services/webhook"
	"net/http"

	"appseclabs/middlewares"

	"github.com/labstack/echo/v4"
	mid "github.com/labstack/echo/v4/middleware"
)

func main() {

	// Start Config
	config.CheckEnvVariables()

	// Start Logger
	logger := config.NewZapLogger()
	defer logger.Sync()
	sugar := logger.Sugar()

	// Initialize Clients
	db := database.NewDatabase()
	k8sClient := k8s.NewK8sClient(db)

	// Initialize Services
	webhookService := webhook.NewWebhookService()
	if webhookService.IsEnabled() {
		sugar.Info("Webhook service is enabled")
	} else {
		sugar.Info("Webhook service is disabled")
	}

	// Initialize Usecases
	labUsecase := lab.NewLabUsecase(k8sClient, db, sugar, webhookService)
	labDefinitionsUsecase := labdefinitions.NewLabDefinitionsUsecase(db)
	// Initialize Handlers
	labHandler := lab.NewLabHandler(labUsecase)
	labDefinitionsHandler := labdefinitions.NewLabDefinitionsHandler(labDefinitionsUsecase)
	// Create Echo Instance
	e := echo.New()

	// Middlewares
	e.Logger = &config.ZapLogger{Sugar: sugar}
	e.Use(middlewares.ZapLoggerMiddleware(logger))
	e.Use(mid.Recover())

	// Health Check
	e.GET("/healthcheck", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})

	// API Group
	gAPI := e.Group("/api")

	// Lab Routes
	labAPI := gAPI.Group("/labs")
	labAPI.POST("/", labHandler.CreateLabHandler)
	labAPI.DELETE("/:namespace", labHandler.DeleteLabHandler)
	labAPI.GET("/:namespace", labHandler.GetLabResultHandler)
	labAPI.POST("/:namespace/finish", labHandler.FinishLabHandler)

	// Lab Definitions Routes
	labDefinitionsAPI := gAPI.Group("/lab-definitions")
	labDefinitionsAPI.GET("/", labDefinitionsHandler.GetAllLabDefinitionsHandler)
	labDefinitionsAPI.POST("/", labDefinitionsHandler.CreateLabDefinitionHandler)
	labDefinitionsAPI.GET("/:slug", labDefinitionsHandler.GetLabDefinitionHandler)
	labDefinitionsAPI.PUT("/:slug", labDefinitionsHandler.UpdateLabDefinitionHandler)
	labDefinitionsAPI.DELETE("/:slug", labDefinitionsHandler.DeleteLabDefinitionHandler)

	// Lab Evaluators Routes
	evaluationAPI := gAPI.Group("/evaluator")
	evaluationAPI.GET("/", labDefinitionsHandler.GetEvaluators)

	e.Logger.Fatal(e.Start(":8084"))
}
