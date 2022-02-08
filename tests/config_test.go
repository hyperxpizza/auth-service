package main

import (
	"testing"

	"github.com/hyperxpizza/auth-service/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	if *configPathOpt == "" {
		t.Fail()
		return
	}

	cfg, err := config.NewConfig(*configPathOpt)
	assert.NoError(t, err)
	cfg.PrettyPrint()
}
