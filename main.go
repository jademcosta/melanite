package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/jademcosta/melanite/config"
	"github.com/jademcosta/melanite/controllers/imagecontroller"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

const defaultLogLevel = log.InfoLevel
const defaultPort = "8080"

var defaultLogFormatter = &log.JSONFormatter{}

func main() {
	logger := buildLogger()

	configuration, err := loadConfig()
	if err != nil {
		logger.Panic(err)
	}

	var port string
	if configuration.Port != "" {
		port = configuration.Port
	} else {
		port = defaultPort
	}

	app := GetApp(*configuration, logger)

	srv := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      app,
		Addr:         fmt.Sprintf(":%s", port),
	}

	logger.Infof("Starting Melanite on port %s", port)
	logger.Fatal(srv.ListenAndServe())
}

func GetApp(configuration config.Config, logger *log.Logger) http.Handler {

	r := http.NewServeMux()
	r.Handle("/", imagecontroller.New(configuration, logger))

	n := negroni.New(negroni.NewRecovery())
	n.Use(negronilogrus.NewMiddlewareFromLogger(logger,
		"melanite"))

	n.UseHandler(r)
	return n
}

func getConfigFileContent(configFilePath string) ([]byte, error) {

	if configFilePath == "" {
		return []byte{}, nil
	}

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("File %s does not exist", configFilePath)
	}

	configContent, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	return configContent, nil
}

func loadConfig() (*config.Config, error) {
	var configFilePath = flag.String("c", "", "The path of the yaml config file")
	flag.Parse()

	configFileContent, err := getConfigFileContent(*configFilePath)
	if err != nil {
		return nil, err
	}

	configuration, err := config.New(configFileContent)
	if err != nil {
		return nil, err
	}
	return &configuration, nil
}

func buildLogger() *log.Logger {
	logger := log.New()
	logger.SetLevel(defaultLogLevel)
	logger.Formatter = defaultLogFormatter
	return logger
}
