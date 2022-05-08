package main

import (
	"flag"

	"github.com/hyperxpizza/auth-service/pkg/impl"
	"github.com/sirupsen/logrus"
)

var configPathOpt = flag.String("config", "", "path to config file")
var loglevelOpt = flag.String("loglevel", "warn", "logger level")

func main() {
	flag.Parse()
	if *configPathOpt == "" {
		panic("config flag not set")
	}

	logger := logrus.New()
	if level, err := logrus.ParseLevel(*loglevelOpt); err == nil {
		logger.Level = level
	}

	authServiceServer, err := impl.NewAuthServiceServer(*configPathOpt, logger)
	if err != nil {
		logger.Fatal(err)
	}

	authServiceServer.WithTlsEnabled()

	if err := authServiceServer.Run(); err != nil {
		logger.Fatal(err)
	}
}
