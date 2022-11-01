package internal

var defaultPaths = []string{".", "/usr/local/etc", "/etc"}

const defaultFileName = "config"

// Struct types
const (
	urlType    = "url.URL"
	timeType   = "time.Time"
	loggerType = "types.Logger"
)

// Based on general types
const (
	durationType    = "time.Duration"
	logrusLevelType = "logrus.Level"
	syslogLevelType = "types.SyslogLevel"
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

const (
	maxInt = 9223372036854775807
	minInt = -9223372036854775808

	maxInt8 = 127
	minInt8 = -128

	maxInt16 = 32767
	minInt16 = -32768

	maxInt32 = 2147483647
	minInt32 = -2147483648

	maxInt64 = 9223372036854775807
	minInt64 = -9223372036854775808

	maxUint = 18446744073709551615
	minUint = 0

	maxUint8 = 255
	minUint8 = 0

	maxUint16 = 65535
	minUint16 = 0

	maxUint32 = 4294967295
	minUint32 = 0

	maxUint64 = 18446744073709551615
	minUint64 = 0
)
