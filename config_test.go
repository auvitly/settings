package config

import (
	"github.com/sirupsen/logrus"
	"settings/internal"
	"testing"
	"time"
)

type Config struct {
	Structure Subconfig `json:"structure"`
}

type Subconfig struct {
	Slice []Simple `json:"slice"`
}

type Simple struct {
	Int      int64               `json:"int" default:"10"`
	String   string              `json:"string" default:"hello"`
	Float64  float64             `json:"float64" default:"23.1"`
	Bool     bool                `json:"bool" default:"true"`
	Duration time.Duration       `json:"dur" default:"1m"`
	Servers  []*string           `json:"servers" validation:"tcp_addr"`
	MapType1 map[string]string   `json:"map_type_1"`
	MapType2 map[string][]string `json:"map_type_2"`
	MapType3 map[string][]Simple `json:"map_type_3"`
}

func TestConfigurator(t *testing.T) {

	c := internal.New("json", "./")

	conf := &Config{}
	c.LoadConfiguration()

	err := c.Unmarshal(conf)

	logrus.Info(err)
}
