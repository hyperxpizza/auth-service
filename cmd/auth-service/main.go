package main

import (
	"flag"
	"log"

	"github.com/hyperxpizza/auth-service/pkg/impl"
	"github.com/sirupsen/logrus"
)

var configPathOpt = flag.String("config", "", "path to config file")
var loglevelOpt = flag.String("loglevel", "warn", "logger level")

func main() {
	flag.Parse()
	if *configPathOpt == "" {
		log.Fatal("config flag not set")
		return
	}

	logger := logrus.New()
	if level, err := logrus.ParseLevel(*loglevelOpt); err == nil {
		logger.Level = level
	}

	authServiceServer, err := impl.NewAuthServiceServer(*configPathOpt, logger)
	if err != nil {
		log.Fatal(err)
		return
	}

	authServiceServer.Run()
}
