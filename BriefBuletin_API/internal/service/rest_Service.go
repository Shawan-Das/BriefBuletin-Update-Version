package service

import (
	"encoding/json"
	"fmt"

	"github.com/rest/api/internal/model"
	"github.com/rest/api/internal/util"

	"github.com/sirupsen/logrus"
)

var _asLogger = logrus.New()
var shareLogger = logrus.New()
var _sysCtlLogger = logrus.New()

// -------------------------- Auth Service ------------------------------------------------------------
// RESTService provides authentication related rest services
type RESTService struct {
	dbConn        *util.DBConnectionWrapper
	jwtSigningKey []byte
	bypassAuth    map[string]bool
}

// NewAuthenticationRESTService returns a new initialized version of the service
func NewAuthenticationRESTService(config []byte, dbConnection *util.DBConnectionWrapper, verbose bool) *RESTService {
	service := new(RESTService)
	if err := service.Init(config, dbConnection, verbose); err != nil {
		_asLogger.Errorf("Unable to initialize service instance %v", err)
		return nil
	}
	return service
}

// Init initializes the service instance
func (s *RESTService) Init(config []byte, dbConnection *util.DBConnectionWrapper, verbose bool) error {
	if verbose {
		_asLogger.SetLevel(logrus.DebugLevel)
	}
	if dbConnection == nil {
		return fmt.Errorf("null DB Util reference passed")
	}
	s.dbConn = dbConnection
	var conf model.AuthServiceConfig
	err := json.Unmarshal(config, &conf)
	if err != nil {
		_asLogger.Error("Unable to parse config json file ", err)
		return err
	}
	if conf.JWTKey != nil && len(*conf.JWTKey) > 0 {
		s.jwtSigningKey = []byte(*conf.JWTKey)
	}
	s.bypassAuth = make(map[string]bool)
	s.bypassAuth["/"] = true
	if len(conf.BypassAuth) > 0 { //conf.BypassAuth != nil &&
		for _, url := range conf.BypassAuth {
			s.bypassAuth[url] = true
		}
	}
	_asLogger.Infof("Successfully initialized AuthenticationRESTService")
	return nil
}

// ------------------------------------ Open Route Service --------------------------------------------------------
type OpenAPIService struct {
	dbConn *util.DBConnectionWrapper
}

func NewSharedAPIService(config []byte, dbConnection *util.DBConnectionWrapper, verbose bool) *OpenAPIService {
	service := new(OpenAPIService)
	if err := service.Init(config, dbConnection, verbose); err != nil {
		shareLogger.Errorf("Error in creating service instance %v", err)
		return nil
	}
	return service
}
func (s *OpenAPIService) Init(config []byte, dbConnection *util.DBConnectionWrapper, verbose bool) error {
	if verbose {
		shareLogger.SetLevel(logrus.DebugLevel)
	}
	if dbConnection == nil {
		return fmt.Errorf("NullDBReference")
	}
	s.dbConn = dbConnection
	return nil
}
