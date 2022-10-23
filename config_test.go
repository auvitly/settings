package config

import (
	"net/url"
	"settings/internal"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Structure *Subconfig `json:"structure"`
}

type Subconfig struct {
	Slice []Simple `json:"slice"`
}

type Simple struct {
	Int      int           `json:"int"`
	String   *string       `json:"string" default:"hello"`
	Duration time.Duration `json:"duration" default:"1m"`
	Url      *url.URL      `json:"url" default:"vk.com"`
}

func TestConfigurator(t *testing.T) {

	c := internal.New("json", "./")

	conf := &Config{}
	c.LoadConfiguration()

	err := c.Unmarshal(conf)

	logrus.Info(err)
}
