package main

import (
	"flag"
	"fmt"
	router "github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"resumpro/internal/app/resumepro/server/handlers"
	"resumpro/internal/app/resumepro/server/middleware"
	"resumpro/internal/app/resumepro/server/routes"
	"syscall"
)

const keyENV = "APP_ENV"

var AppVersion = "dev"
var AppVersionJSON = ""
var AppName = "resumepro"

func readConfig(env string) {
	if len(env) > 0 {
		env = fmt.Sprintf(".%s", env)
	}

	viper.SetConfigFile(fmt.Sprintf("./configs/app%s.yml", env))
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

// For govvv versioning of executable
var GitCommit, GitBranch, GitState, GitSummary, BuildDate, Version string

var logger *log.Logger

func init() {
	env := os.Getenv(keyENV)

	logger = log.New()
	logger.Info(fmt.Sprintf("Starting %s on %s env..", AppName, env))

	readConfig(env)

	v := flag.Bool("v", false, "Prints the version details")
	version := flag.Bool("version", false, "Prints all the version details")
	versionJSON := flag.Bool("json", false, "Prints all the version details in JSON format")
	flag.Parse()

	if *v {
		fmt.Printf("Version: %s\n", Version)
		os.Exit(0)
	}

	AppVersionJSON = fmt.Sprintf("{\"version\":\"%s\",\"build_date\":\"%s\",\"git_commit\":\"%s\",\"git_branch\":\"%s\",\"git_state\":\"%s\",\"git_summary\":\"%s\"}\n", Version, BuildDate, GitCommit, GitBranch, GitState, GitSummary)
	AppVersion = fmt.Sprintf("Version: %s, %s\nBuilt from: %s (%s), %s, %s\n", Version, BuildDate, GitCommit, GitBranch, GitState, GitSummary)

	switch *version {
	case true:
		switch *versionJSON {
		case true:
			fmt.Println(AppVersionJSON)
		default:
			fmt.Println(AppVersion)
		}
		os.Exit(0)
	}
}

func main() {
	host, port := viper.GetString("host"), viper.GetString("port")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT)

	//r := routes.GetRoutesList(logger, router.NewRouter(), middleware.NewLog(logger, true))
	//var r routes.Routes
	//r.Init(logger, router.NewRouter())
	r := routes.GetRoutesList(logger, router.NewRouter(), middleware.NewLog(logger, true))
	r.Add("/ping", "GET", handlers.NewPing(logger, AppVersionJSON).Handler)
	r.Parse()

	httpErr := make(chan error, 1)
	go func() {
		logger.Info(fmt.Sprintf("Started server on %s:%s..", host, port))
		httpErr <- http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), r.Router)
	}()

	select {
	case err := <-httpErr:
		logger.Error(err.Error())
	case <-stop:
		logger.Info("Stopped via signal")
	}

	logger.Info(fmt.Sprintf("Stopping %s..", AppName))
}
