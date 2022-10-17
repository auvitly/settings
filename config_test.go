package config

import (
	"github.com/sirupsen/logrus"
	"settings/internal"
	"testing"
)

type Config struct {
	Credentials *Subconfig `json:"credentials"`
}

type Subconfig struct {
	Host string             `json:"host" default:"hello"`
	Port string             `json:"port"`
	Map  map[string]float32 `json:"map"`
	SS   []string           `json:"strslice"`
	SI   []int              `json:"intslice"`
}

func TestConfigurator(t *testing.T) {

	c := internal.New("json", "./../../config/")

	conf := &Config{}
	c.LoadConfiguration()

	c.UnmarshalInto(conf)

	logrus.Info()
}
