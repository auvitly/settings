package tests

import (
	"net/url"
	"time"
)

type Config struct {
	Base struct {
		Numbers struct {
			Int        int      `json:"int"`
			PtrInt     *int     `json:"int"`
			Int8       int8     `json:"int8"`
			PtrInt8    *int8    `json:"int8"`
			Int16      int16    `json:"int16"`
			PtrInt16   *int16   `json:"int16"`
			Int32      int32    `json:"int32"`
			PtrInt32   *int32   `json:"int32"`
			Int64      int64    `json:"int64"`
			PtrInt64   *int64   `json:"int64"`
			Uint       uint     `json:"uint"`
			PtrUint    *uint    `json:"uint"`
			Uint8      uint8    `json:"uint8"`
			PtrUint8   *uint8   `json:"uint8"`
			Uint16     uint16   `json:"uint16"`
			PtrUint16  *uint16  `json:"uint16"`
			Uint32     uint32   `json:"uint32"`
			PtrUint32  *uint32  `json:"uint32"`
			Uint64     uint64   `json:"uint64"`
			PtrUint64  *uint64  `json:"uint64"`
			Float32    float32  `json:"float32"`
			PtrFloat32 *float32 `json:"float32"`
			Float64    float64  `json:"float64"`
			PtrFloat64 *float64 `json:"float64"`
		} `json:"numbers"`
		Strings struct {
			String    string  `json:"string"`
			PtrString *string `json:"string"`
		} `json:"strings"`
		Slices struct {
			SliceStrings []string         `json:"strings_slice"`
			SliceBytes   []byte           `json:"bytes_slice"`
			SliceMap     []map[string]int `json:"map_slice"`
			SliceStructs []struct {
				Field1 int `json:"field_1"`
			} `json:"struct_slice"`
		} `json:"slices"`
		Maps struct {
			MapInt         map[string]int            `json:"map_int"`
			MapString      map[string]string         `json:"map_string"`
			MapSliceString map[string][]string       `json:"map_slice_string"`
			MapMapInt      map[string]map[string]int `json:"map_map_int"`
			MapStruct      map[string]struct {
				Field1 int `json:"field_1"`
			} `json:"map_struct"`
		} `json:"maps"`
	} `json:"base"`
	Optional struct {
		Time    time.Time      `json:"time"`
		PtrTime *time.Time     `json:"time"`
		Dur     time.Duration  `json:"dur"`
		PtrDur  *time.Duration `json:"dur"`
		Url     url.URL        `json:"url"`
		PtrUrl  *url.URL       `json:"url"`
	} `json:"optional"`
	Logger struct {
		Level          string `json:"level"`
		SyslogAddr     string `json:"syslog_addr"`
		SyslogProtocol string `json:"syslog_protocol"`
		SyslogLevel    string `json:"syslog_level"`
		Colour         bool   `json:"colour"`
		Stdout         bool   `json:"stdout"`
		GraylogLevel   string `json:"graylog_level"`
		Graylog        string `json:"graylog"`
	} `json:"logger"`
	PtrLogger *struct {
		Level          string `json:"level"`
		SyslogAddr     string `json:"syslog_addr"`
		SyslogProtocol string `json:"syslog_protocol"`
		SyslogLevel    string `json:"syslog_level"`
		Colour         bool   `json:"colour"`
		Stdout         bool   `json:"stdout"`
		GraylogLevel   string `json:"graylog_level"`
		Graylog        string `json:"graylog"`
	} `json:"logger"`
}
