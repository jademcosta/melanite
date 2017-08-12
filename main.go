package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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
	configFileContent, err := getConfigFileContent()
	if err != nil {
		panic(err)
	}

	configuration, err := config.New(configFileContent)
	if err != nil {
		panic(err)
	}

	http.ListenAndServe(fmt.Sprintf(":%s", defaultPort), GetApp(defaultLogLevel, defaultLogFormatter, configuration))
}

func GetApp(logLevel log.Level, logFormatter log.Formatter,
	configuration config.Config) http.Handler {

	r := http.NewServeMux()
	r.Handle("/", imagecontroller.New(configuration))

	n := negroni.New(negroni.NewRecovery())
	n.Use(negronilogrus.NewMiddlewareFromLogger(getLogger(logLevel, logFormatter),
		"melanite"))

	n.UseHandler(r)
	return n
}

func getConfigFileContent() ([]byte, error) {
	var configFilePath = flag.String("c", "", "The path of the yaml config file")
	flag.Parse()

	if *configFilePath == "" {
		return nil, fmt.Errorf("Config file was not provided. Use -c FILENAME to provide one")
	}

	if _, err := os.Stat(*configFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("File %s does not exist", *configFilePath)
	}

	configContent, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		return nil, err
	}

	return configContent, nil
}

func getLogger(logLevel log.Level, logFormatter log.Formatter) *log.Logger {
	appLog := log.New()
	appLog.SetLevel(logLevel)
	appLog.Formatter = logFormatter
	return appLog
}
