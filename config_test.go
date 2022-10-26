package config

import (
	"github.com/sirupsen/logrus"
	"settings/internal"
	"testing"
)

type Config struct {
	Structure Subconfig `json:"structure"`
}

type Subconfig struct {
	Slice []Simple `json:"slice"`
}

type Simple struct {
	//Int      int64         `json:"int" default:"10"`
	//String   string        `json:"string" default:"hello"`
	//Float64  float64       `json:"float64" default:"23.1"`
	//Bool     bool          `json:"bool" default:"true"`
	//Duration time.Duration `json:"dur" default:"1m"`
	Servers []string `json:"servers" validation:"tcp_addr"`
}

func TestConfigurator(t *testing.T) {

	c := internal.New("json", "./")

	conf := &Config{}
	c.LoadConfiguration()

	err := c.Unmarshal(conf)

	logrus.Info(err)
}
