package controller

import (
	"fmt"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	appMiddleware "github.com/theCompanyDream/id-trials/apps/backend/middleware"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
)

func RunServer(db *gorm.DB) {
	server := NewEchoServer(db)
	// Start the server
	server.Logger.Info("Server is running...")
	port := os.Getenv("BACKEND_PORT")
	if port != "" {
		serverStartCode := fmt.Sprintf(":%s", port)
		server.Logger.Fatal(server.Start(serverStartCode))
	} else {
		server.Logger.Fatal(server.Start(":3000"))
	}
}

func NewEchoServer(db *gorm.DB) *echo.Echo {
	server := echo.New()

	server.HTTPErrorHandler = appMiddleware.HttpErrorHandler
	metricsMiddleware := appMiddleware.NewMetricsMiddleware(db)
	appMiddleware.NewLogger()

	analyticsController := NewAnalyticsController(db)
	ulidController := NewUlidController(db)
	uuid4Controller := NewGormUuidController(db)
	nanoIdController := NewGormNanoController(db)
	ksuidController := NewGormKsuidController(db)
	cuidController := NewGormCuidController(db)
	snowController := NewSnowCuidController(db)

	// Middleware
	server.Use(appMiddleware.LoggingMiddleware)
	server.Use(middleware.Recover())
	server.Use(middleware.RequestID())
	server.Use(middleware.RequestLogger()) // Add request logging for security auditing
	server.Use(middleware.Gzip())
	server.Use(middleware.BodyLimit("20k"))
	server.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(10))))
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Split(os.Getenv("ALLOWED_HOSTS"), ","),
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	server.Use(middleware.Secure())
	server.Use(metricsMiddleware.CaptureMetrics())

	server.GET("/analytics/comparison", analyticsController.GetIDTypeComparison)
	server.GET("/analytics/:type/details", analyticsController.GetIDTypeDetails)
	server.GET("/analytics/:type/percentiles", analyticsController.GetPercentiles)
	server.GET("/analytics/errors", analyticsController.GetErrorRates)
	server.GET("/analytics/:type/timeseries", analyticsController.GetTimeSeries)
	server.GET("/analytics/tableSize", analyticsController.GetTableSizeData)
	// Define main routes
	server.GET("/swagger/*", echoSwagger.WrapHandler)
	server.GET("/", Home)
	server.GET("/ulidIds", ulidController.GetUsers)
	server.GET("/ulidId/:id", ulidController.GetUser)
	server.POST("/ulidId", ulidController.CreateUser)
	server.PUT("/ulidId/:id", ulidController.UpdateUser)
	server.DELETE("/ulidId/:id", ulidController.DeleteUser)
	//uuid
	server.GET("/uuid4s", uuid4Controller.GetUsers)
	server.GET("/uuid4/:id", uuid4Controller.GetUser)
	server.POST("/uuid4", uuid4Controller.CreateUser)
	server.PUT("/uuid4/:id", uuid4Controller.UpdateUser)
	server.DELETE("/uuid4/:id", uuid4Controller.DeleteUser)
	//nanoId
	server.GET("/nanoIds", nanoIdController.GetUsers)
	server.GET("/nanoId/:id", nanoIdController.GetUser)
	server.POST("/nanoId", nanoIdController.CreateUser)
	server.PUT("/nanoId/:id", nanoIdController.UpdateUser)
	server.DELETE("/nanoId/:id", nanoIdController.DeleteUser)
	//ksuidId
	server.GET("/ksuidIds", ksuidController.GetUsers)
	server.GET("/ksuidId/:id", ksuidController.GetUser)
	server.POST("/ksuidId", ksuidController.CreateUser)
	server.PUT("/ksuidId/:id", ksuidController.UpdateUser)
	server.DELETE("/ksuidId/:id", ksuidController.DeleteUser)
	//cuid
	server.GET("/cuidIds", cuidController.GetUsers)
	server.GET("/cuidId/:id", cuidController.GetUser)
	server.POST("/cuidId", cuidController.CreateUser)
	server.PUT("/cuidId/:id", cuidController.UpdateUser)
	server.DELETE("/cuidId/:id", cuidController.DeleteUser)

	server.GET("/snowIds", snowController.GetUsers)
	server.GET("/snowId/:id", snowController.GetUser)
	server.POST("/snowId", snowController.CreateUser)
	server.PUT("/snowId/:id", snowController.UpdateUser)
	server.DELETE("/snowId/:id", snowController.DeleteUser)

	return server
}

func NewServerlessEchoServer(db *gorm.DB) *echo.Echo {
	server := echo.New()

	server.HTTPErrorHandler = appMiddleware.HttpErrorHandler
	metricsMiddleware := appMiddleware.NewMetricsMiddleware(db)
	appMiddleware.NewLogger()

	analyticsController := NewAnalyticsController(db)
	ulidController := NewUlidController(db)
	uuid4Controller := NewGormUuidController(db)
	nanoIdController := NewGormNanoController(db)
	ksuidController := NewGormKsuidController(db)
	cuidController := NewGormCuidController(db)
	snowController := NewSnowCuidController(db)

	// Middleware
	server.Use(appMiddleware.LoggingMiddleware)
	server.Use(middleware.Recover())
	server.Use(middleware.RequestID())
	server.Use(middleware.RequestLogger()) // Add request logging for security auditing
	server.Use(middleware.Gzip())
	server.Use(middleware.BodyLimit("20k"))
	server.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(10))))
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Split(os.Getenv("ALLOWED_HOSTS"), ","),
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	server.Use(middleware.Secure())
	server.Use(metricsMiddleware.CaptureMetrics())

	api := server.Group("/api")

	api.GET("/analytics/comparison", analyticsController.GetIDTypeComparison)
	api.GET("/analytics/:type/details", analyticsController.GetIDTypeDetails)
	api.GET("/analytics/:type/percentiles", analyticsController.GetPercentiles)
	api.GET("/analytics/errors", analyticsController.GetErrorRates)
	api.GET("/analytics/:type/timeseries", analyticsController.GetTimeSeries)
	// Define main routes
	api.GET("/swagger/*", echoSwagger.WrapHandler)
	api.GET("/", Home)
	api.GET("/ulidIds", ulidController.GetUsers)
	api.GET("/ulidId/:id", ulidController.GetUser)
	api.POST("/ulidId", ulidController.CreateUser)
	api.PUT("/ulidId/:id", ulidController.UpdateUser)
	api.DELETE("/ulidId/:id", ulidController.DeleteUser)
	//uuid
	api.GET("/uuid4s", uuid4Controller.GetUsers)
	api.GET("/uuid4/:id", uuid4Controller.GetUser)
	api.POST("/uuid4", uuid4Controller.CreateUser)
	api.PUT("/uuid4/:id", uuid4Controller.UpdateUser)
	api.DELETE("/uuid4/:id", uuid4Controller.DeleteUser)
	//nanoId
	api.GET("/nanoIds", nanoIdController.GetUsers)
	api.GET("/nanoId/:id", nanoIdController.GetUser)
	api.POST("/nanoId", nanoIdController.CreateUser)
	api.PUT("/nanoId/:id", nanoIdController.UpdateUser)
	api.DELETE("/nanoId/:id", nanoIdController.DeleteUser)
	//ksuidId
	api.GET("/ksuidIds", ksuidController.GetUsers)
	api.GET("/ksuidId/:id", ksuidController.GetUser)
	api.POST("/ksuidId", ksuidController.CreateUser)
	api.PUT("/ksuidId/:id", ksuidController.UpdateUser)
	api.DELETE("/ksuidId/:id", ksuidController.DeleteUser)
	//cuid
	api.GET("/cuidIds", cuidController.GetUsers)
	api.GET("/cuidId/:id", cuidController.GetUser)
	api.POST("/cuidId", cuidController.CreateUser)
	api.PUT("/cuidId/:id", cuidController.UpdateUser)
	api.DELETE("/cuidId/:id", cuidController.DeleteUser)

	api.GET("/snowIds", snowController.GetUsers)
	api.GET("/snowId/:id", snowController.GetUser)
	api.POST("/snowId", snowController.CreateUser)
	api.PUT("/snowId/:id", snowController.UpdateUser)
	api.DELETE("/snowId/:id", snowController.DeleteUser)

	return server
}
