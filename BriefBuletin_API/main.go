package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	// docs "github.com/rest/api/docs"
	"github.com/rest/api/internal/service"
	"github.com/spf13/viper"
)

func main() {
	configFilePath := flag.String("c", "./config.json", "Configuration file")
	verbose := flag.Bool("v", false, "Verbose")
	port := flag.Int("port", 8080, "Port to run the server on")

	flag.Parse()
	configBytes, err := os.ReadFile(*configFilePath)
	if err != nil {
		fmt.Println("Unable to read configuration file ", err.Error())
		os.Exit(1)
	}
	initAppConfigViper(*configFilePath)

	// // ^ Swagger start
	// docs.SwaggerInfo.Host = "" // dynamic: picks request's host:port
	// docs.SwaggerInfo.BasePath = "/"
	// // ^ Swagger end

	srvcInstance := service.NewAPIServer(configBytes, *verbose)
	srvcInstance.Serve(*port)
}

func initAppConfigViper(configPath string) {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Panic("There is no such a config file in path ", configPath)
		} else {
			log.Panic("There is some problem about data in file")
		}
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
