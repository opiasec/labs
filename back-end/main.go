package main

import (
	"appseclabsplataform/admin"
	"appseclabsplataform/auth"
	"appseclabsplataform/config"
	"appseclabsplataform/dashboard"
	"appseclabsplataform/database"
	lab "appseclabsplataform/lab"
	labdefinition "appseclabsplataform/labDefinition"
	"appseclabsplataform/middleware"
	labcluster "appseclabsplataform/services/labCluster"
	labide "appseclabsplataform/services/labIDE"
	labsession "appseclabsplataform/services/labSession"
	"appseclabsplataform/webhook"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	mid "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Check env vars
	if err := config.CheckEnvironmentVariables(); err != nil {
		slog.Error("failed to check environment variables", "error", err)
		os.Exit(1)
	}

	// Initialize Config
	config := config.NewConfig()

	// Start Services
	labCluster := labcluster.NewLabClusterService(config)
	LabIDEService := labide.NewLabIDEService(config)
	postgresDB := database.NewDatabase(config)
	labSession := labsession.NewLabSessionService(postgresDB, labCluster, config)
	go labSession.ListenForLabSessionExpiry()

	//Start Usecases
	labUsecase := lab.NewLabUsecase(labCluster, postgresDB, labSession, LabIDEService)
	labDefinitionUsecase := labdefinition.NewLabDefinitionUsecase(labCluster, postgresDB)
	adminUsecase := admin.NewAdminUsecase(postgresDB, labCluster)
	webhookUsecase := webhook.NewWebHookUsecase(postgresDB)
	dashboardUsecase := dashboard.NewDashboardUsecase(postgresDB)
	authUsecase := auth.NewAuthUsecase(postgresDB, config)

	// Start Handlers
	labHandler := lab.NewLabHandler(labUsecase)
	labDefinitionHandler := labdefinition.NewLabDefinitionHandler(labDefinitionUsecase)
	adminHandler := admin.NewAdminHandler(adminUsecase)
	webhookHandler := webhook.NewWebHookHandler(webhookUsecase)
	dashboardHandler := dashboard.NewDashboardHandler(dashboardUsecase)
	authHandler := auth.NewAuthHandler(authUsecase)

	// Echo Instance
	e := echo.New()

	// Setup Configuration
	if config.AuthConfig.PasswordEnabled {
		slog.Info("Password authentication is enabled")
		err := postgresDB.CreateRootAccount(config.AuthConfig)
		if err != nil {
			slog.Error("failed to create root account", "error", err)
			os.Exit(1)
		}
	}

	// Middleware
	e.Use(mid.Logger())
	e.Use(mid.Recover())
	e.Use(mid.CORSWithConfig(mid.CORSConfig{
		AllowOrigins: []string{"http://127.0.0.1:4200", "http://localhost:4200", "https://labs.opiasec.com"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Healthcheck
	e.GET("/healthcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	// Config Routes
	gConfig := e.Group("/api/config")
	gConfig.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, config.PublicConfig)
	})

	//Auth Routes
	gAuth := e.Group("/api/auth")
	gAuth.POST("/login", authHandler.Login)

	//API Routes
	gAPI := e.Group("/api")
	gAPI.Use(middleware.Auth(config))

	// Lab Routes
	gLabs := gAPI.Group("/labs")
	gLabs.POST("/", labHandler.CreateLab)
	gLabs.GET("/:namespace", labHandler.GetLabResult)
	gLabs.GET("/:namespace/status", labHandler.GetLabStatus)
	gLabs.POST("/:namespace/finish", labHandler.FinishLab)
	gLabs.POST("/:namespace/leave", labHandler.LeaveLab)
	gLabs.POST("/:namespace/feedback", labHandler.SendFeedback)
	gLabs.GET("/", labHandler.GetAllLabsByUserAndStatus)

	// Dashboard Routes
	gDashboard := gAPI.Group("/dashboard")
	gDashboard.GET("/data", dashboardHandler.GetDashboardData)

	// Lab Admin Routes - Lab Sessions
	gLabsAdmin := gAPI.Group("/admin/labs")
	gLabsAdmin.Use(middleware.Admin)
	gLabsAdmin.GET("", adminHandler.GetLabsSessions)
	gLabsAdmin.GET("/:namespace", adminHandler.GetLabSession)
	gLabsAdmin.GET("/status", adminHandler.GetPossiblesStatus)
	gLabsAdmin.POST("/:namespace/status", adminHandler.ChangeLabStatus)

	// Lab Admin Routes - Lab Definition CRUD
	gLabsAdminDef := gAPI.Group("/admin/lab-definition")
	gLabsAdminDef.Use(middleware.Admin)
	gLabsAdminDef.GET("/:slug", adminHandler.GetLabDefinition)
	gLabsAdminDef.GET("/", adminHandler.GetLabsDefinitions)
	gLabsAdminDef.POST("/", adminHandler.CreateLabDefinition)
	gLabsAdminDef.PUT("/:slug", adminHandler.UpdateLabDefinition)
	gLabsAdminDef.DELETE("/:slug", adminHandler.DeleteLabDefinition)
	gLabsAdminDef.GET("/vulnerabilities", adminHandler.GetPossiblesVulnerabilities)
	gLabsAdminDef.GET("/languages", adminHandler.GetPossiblesLanguages)
	gLabsAdminDef.GET("/technologies", adminHandler.GetPossiblesTechnologies)
	gLabsAdminDef.GET("/evaluators", adminHandler.GetPossiblesEvaluators)
	gLabsAdminDef.GET("/images", adminHandler.GetPossiblesImages)

	// Lab Admin Routes - User Management
	gLabsAdminUsers := gAPI.Group("/admin/users")
	gLabsAdminUsers.Use(middleware.Admin)
	gLabsAdminUsers.GET("/", adminHandler.GetAllUsers)
	gLabsAdminUsers.POST("/", adminHandler.CreateUser)
	gLabsAdminUsers.GET("/:id", adminHandler.GetUserByID)
	gLabsAdminUsers.PUT("/:id", adminHandler.UpdateUser)
	gLabsAdminUsers.DELETE("/:id", adminHandler.DeleteUser)

	// Lab Definitions Routes
	gLabDefinitions := gAPI.Group("/lab-definition")
	gLabDefinitions.GET("/", labDefinitionHandler.GetLabsDefinitions)
	gLabDefinitions.GET("/:slug", labDefinitionHandler.GetLabDefinitionBySlug)

	// Webhooks
	gWebhooks := e.Group("/webhook")
	gWebhooks.Use(middleware.WebhookAuth)
	gWebhooks.POST("/finish-evaluation", webhookHandler.FinishEvaluationResult)

	setupTLSServer(e, config)
}

func setupTLSServer(e *echo.Echo, config *config.Config) {
	if config.TLSConfig.Enabled && config.TLSConfig.CertFile != "" && config.TLSConfig.KeyFile != "" {
		slog.Info("Starting server with TLS",
			"cert", config.TLSConfig.CertFile,
			"key", config.TLSConfig.KeyFile,
			"port", config.TLSConfig.Port)

		// HTTP Server para redirect
		go func() {
			redirectServer := &http.Server{
				Addr: ":80",
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					target := "https://" + r.Host + r.URL.Path
					if len(r.URL.RawQuery) > 0 {
						target += "?" + r.URL.RawQuery
					}
					slog.Info("Redirecting to HTTPS", "from", r.URL.String(), "to", target)
					http.Redirect(w, r, target, http.StatusMovedPermanently)
				}),
			}

			slog.Info("Starting HTTP redirect server on :80")
			if err := redirectServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				slog.Error("failed to start HTTP redirect server", "error", err)
			}
		}()

		// HTTPS Server com certificados próprios
		slog.Info("Starting HTTPS server", "port", config.TLSConfig.Port)
		if err := e.StartTLS(":"+config.TLSConfig.Port, config.TLSConfig.CertFile, config.TLSConfig.KeyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start HTTPS server", "error", err)
		}
	} else {
		// Development mode ou TLS disabled
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		slog.Info("Starting HTTP server", "port", port)
		if err := e.Start(":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start server", "error", err)
		}
	}
}
