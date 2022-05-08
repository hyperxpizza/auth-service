package main

import (
	"flag"

	"github.com/hyperxpizza/auth-service/pkg/config"
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

	cfg, err := config.NewConfig(*configPathOpt)
	if err != nil {
		panic(err)
	}

	logger := logrus.New()
	if level, err := logrus.ParseLevel(*loglevelOpt); err == nil {
		logger.Level = level
	}

	authServiceServer, err := impl.NewAuthServiceServer(cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}

	if cfg.TLS.TlsEnabled {
		authServiceServer.WithTlsEnabled()
	}

	if err := authServiceServer.Run(); err != nil {
		logger.Fatal(err)
	}
}
