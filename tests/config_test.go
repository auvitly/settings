package tests

import (
	"testing"

	config "github.com/Auvitly/settings"
	"github.com/sirupsen/logrus"
)

func TestBaseTypes(t *testing.T) {

	c := config.New("config_json.json", "./configs/")

	conf := &Config{}

	c.LoadOptions()

	err := c.LoadSettings(conf)

	logrus.Info(err)
}
