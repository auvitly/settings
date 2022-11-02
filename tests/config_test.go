package tests

import (
	"github.com/sirupsen/logrus"
	config "settings"
	"testing"
)

func TestBaseTypes(t *testing.T) {

	c := config.New("config_json.json", "./configs/")

	conf := &Config{}

	c.LoadOptions()

	err := c.LoadSettings(conf)

	logrus.Info(err)
}
