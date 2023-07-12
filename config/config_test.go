package config_test

import (
	"feastogether/config"
	"fmt"
	"testing"
)

func TestGetConfig(t *testing.T) {
	if cfg, err := config.GetConfig(".."); err != nil {
		t.Error(err)
	} else {
		fmt.Println(cfg)
	}
}
