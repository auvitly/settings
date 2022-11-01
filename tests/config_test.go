package tests

import (
	"github.com/sirupsen/logrus"
	config "settings"
	"testing"
)

func TestBaseTypes(t *testing.T) {

	c := config.New("config_json", "./configs/")

	conf := &Config{}

	c.LoadConfiguration()

	err := c.Unmarshal(conf)

	logrus.Info(err)
}
