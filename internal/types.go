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

// Supported tags
const (
	env          = "env"
	toml         = "toml"
	omit         = "omit"
	yaml         = "yaml"
	json         = "json"
	xml          = "xml"
	defaultValue = "default"
)

var supportedTags = []string{env, toml, yaml, xml, json, omit, defaultValue}

type Tags map[string]string
type LoadValues map[string]interface{}
