package internal

var defaultPaths = []string{".", "/usr/local/etc", "/etc"}

const defaultFileName = "config"

// Struct types
const (
	urlType  = "url.URL"
	timeType = "time.Time"
)

// Based on general types
const (
	durationType = "<time.Duration Value>"
)

const (
	env          = "env"
	toml         = "toml"
	omit         = "omit"
	yaml         = "yaml"
	json         = "json"
	xml          = "xml"
	defaultValue = "default"
)

var supportedTags = []string{env, toml, omit, yaml, xml, json, defaultValue}

type Tags map[string]string
type LoadValues map[string]interface{}
