package service

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rest/api/internal/util"

	"github.com/sirupsen/logrus"
)

var _ServerLog = logrus.New()

type APIServer struct {
	verbose         bool
	dbConnection    *util.DBConnectionWrapper
	authService     *RESTService
	openAPIService  *OpenAPIService
	serverKey       string
	serverCertFile  string
	isTLS           bool
	server          *http.Server
	shutdownChannel chan os.Signal
}

func NewAPIServer(configBytes []byte, verbose bool) *APIServer {
	server := new(APIServer)
	if err := server.Init(configBytes, verbose); err != nil {
		_ServerLog.Errorf("Unable to initialize the APIServer instance")
		return nil
	}
	return server
}

func (s *APIServer) Init(configBytes []byte, verbose bool) error {
	s.verbose = verbose
	var serverConfig struct {
		IsTLS          bool   `json:"isTLS"`
		ServerKeyPath  string `json:"tlsKeyPath"`
		ServerCertPath string `json:"tlsCertPath"`
	}

	if err := json.Unmarshal(configBytes, &serverConfig); err != nil {
		_ServerLog.Errorf("Unable to parse APIServer config file %v", err)
		return err
	}
	s.isTLS = serverConfig.IsTLS
	if s.isTLS {
		if len(serverConfig.ServerKeyPath) == 0 {
			_ServerLog.Errorf("Server key file missing")
			return fmt.Errorf("server key file missing")
		}
		s.serverKey = serverConfig.ServerKeyPath
		if len(serverConfig.ServerCertPath) == 0 {
			_ServerLog.Errorf("Server certificate file missing")
			return fmt.Errorf("server certificate file missing")
		}
		s.serverCertFile = serverConfig.ServerCertPath
	}

	s.dbConnection = util.NewDBConnectionWrapper(configBytes)
	if s.dbConnection == nil {

		return fmt.Errorf("UnableToInitializeDBConnection")
	}
	openAPIService := NewSharedAPIService(configBytes, s.dbConnection, verbose)
	if openAPIService == nil {
		_ServerLog.Errorf("Error in initializing AuthenticationRESTService")
		return fmt.Errorf("Error_Initializing_AuthenticationRESTService")
	}
	s.openAPIService = openAPIService

	s.authService = NewAuthenticationRESTService(configBytes, s.dbConnection, verbose)
	if s.authService == nil {
		return fmt.Errorf("UnableToInitializeAuthService")
	}

	return nil
}
func (s *APIServer) Serve(port int) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	//TODO: Following to be changed for production
	router.MaxMultipartMemory = 8 << 21 //16 MB Max file size
	cnf := cors.Config{
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	cnf.AllowAllOrigins = true
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler)) // ^ Swagger
	router.Use(cors.New(cnf))
	router.Static("/apidoc", "./api/")
	//shared api service
	// s.sharedAPIService.AddRouters(router)
	//JWT Auth handler
	s.openAPIService.AddRouters(router)
	router.Use(func(c *gin.Context) {
		if s.authService.checkAuth(c) {
			c.Next()
			return
		}
		c.JSON(http.StatusMethodNotAllowed, "Unauthorized")
		c.Abort()
	})
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, buildResponse(200, true, "Service Available", nil))
	})

	//Authentication service
	s.authService.AddRouters(router)
	_ServerLog.Infof("Starting server  on port %d....", port)

	go s.runServer(port, router)
	time.Sleep(1 * time.Second)
	s.shutdownChannel = make(chan os.Signal, 2)
	signal.Notify(s.shutdownChannel, syscall.SIGINT, syscall.SIGTERM)
	<-s.shutdownChannel
	_ServerLog.Info("Trying to shutdown the server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		_ServerLog.Error("Server Shutdown:", err)
		os.Exit(1)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx.Done()
	_ServerLog.Info("Shutdown completed...")

}

func (s *APIServer) runServer(port int, router *gin.Engine) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}
	if s.isTLS {
		cer, err := tls.LoadX509KeyPair(s.serverCertFile, s.serverKey)
		if err != nil {
			_ServerLog.Errorf("Sever certificate any key load error %v", err)
			os.Exit(2)
		}
		tlsConfig := &tls.Config{Certificates: []tls.Certificate{cer}}
		server.TLSConfig = tlsConfig
	}
	s.server = server
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		_ServerLog.Errorf("Starting server  on port %d.... failed %v", port, err)
		os.Exit(3)
	}

}
func (s *APIServer) Shutdown() {
	s.shutdownChannel <- syscall.SIGINT
}
