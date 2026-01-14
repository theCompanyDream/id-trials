package controller

import (
	"fmt"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	appMiddleware "github.com/theCompanyDream/id-trials/apps/backend/middleware"
	"github.com/ziflex/lecho"
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
	logger := lecho.New(
		os.Stdout,
		lecho.WithLevel(log.DEBUG),
		lecho.WithTimestamp(),
		lecho.WithCaller(),
	)

	server.HTTPErrorHandler = appMiddleware.HttpErrorHandler
	metricsMiddleware := appMiddleware.NewMetricsMiddleware(db)
	server.Use(metricsMiddleware.CaptureMetrics())

	analyticsController := NewAnalyticsController(db)
	ulidController := NewUlidController(db)
	uuid4Controller := NewGormUuidController(db)
	nanoIdController := NewGormNanoController(db)
	ksuidController := NewGormKsuidController(db)
	cuidController := NewGormCuidController(db)
	snowController := NewSnowCuidController(db)

	// Middleware
	server.Logger = logger
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
	server.Use(middleware.CORS())
	server.Use(middleware.Secure())

	server.GET("/analytics/comparison", analyticsController.GetIDTypeComparison)
	server.GET("/analytics/:type/details", analyticsController.GetIDTypeDetails)
	server.GET("/analytics/:type/percentiles", analyticsController.GetPercentiles)
	server.GET("/analytics/errors", analyticsController.GetErrorRates)
	server.GET("/analytics/:type/timeseries", analyticsController.GetTimeSeries)
	// Define main routes
	server.GET("/swagger/*", echoSwagger.WrapHandler)
	server.GET("/", Home)
	server.GET("/ulids", ulidController.GetUsers)
	server.GET("/ulid/:id", ulidController.GetUser)
	server.POST("/ulid", ulidController.CreateUser)
	server.PUT("/ulid/:id", ulidController.UpdateUser)
	server.DELETE("/ulid/:id", ulidController.DeleteUser)
	//uuid
	server.GET("/uuid4", uuid4Controller.GetUsers)
	server.GET("/uuid4/:id", uuid4Controller.GetUser)
	server.POST("/uuid4", uuid4Controller.CreateUser)
	server.PUT("/uuid4/:id", uuid4Controller.UpdateUser)
	server.DELETE("/uuid4/:id", uuid4Controller.DeleteUser)
	//nanoId
	server.GET("/nano", nanoIdController.GetUsers)
	server.GET("/nano/:id", nanoIdController.GetUser)
	server.POST("/nano", nanoIdController.CreateUser)
	server.PUT("/nano/:id", nanoIdController.UpdateUser)
	server.DELETE("/nano/:id", nanoIdController.DeleteUser)
	//ksuid
	server.GET("/ksuid", ksuidController.GetUsers)
	server.GET("/ksuid/:id", ksuidController.GetUser)
	server.POST("/ksuid", ksuidController.CreateUser)
	server.PUT("/ksuid/:id", ksuidController.UpdateUser)
	server.DELETE("/ksuid/:id", ksuidController.DeleteUser)
	//cuid
	server.GET("/cuid", cuidController.GetUsers)
	server.GET("/cuid/:id", cuidController.GetUser)
	server.POST("/cuid", cuidController.CreateUser)
	server.PUT("/cuid/:id", cuidController.UpdateUser)
	server.DELETE("/cuid/:id", cuidController.DeleteUser)

	server.GET("/snow", snowController.GetUsers)
	server.GET("/snow/:id", snowController.GetUser)
	server.POST("/snow", snowController.CreateUser)
	server.PUT("/snow/:id", snowController.UpdateUser)
	server.DELETE("/snow/:id", snowController.DeleteUser)

	return server
}
