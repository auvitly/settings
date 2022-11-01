## Оглавление
1. [Описание](#Описание)
2. [Файлы конфигурации](#Файлы конфигурации)
3. [Поддерживаемые типы](#Поддерживаемые типы)

---

### 1. Описание

Для парсинга файлов используется библиотека [`viper`](https://github.com/spf13/viper).

Существует два метода использования данного пакета: 
* через глобальную переменную конфигуратора;
* через создание сущности конфигуратора.

<hr>

#### 1.1 Глобальный конфигуратор

Первый вариант довольно прост и имеет в себе всего две функции:
````go
// LoadOptions - функция загрузки настроек из файла
LoadOptions(name string, paths ...string) (*viper.Viper, error)
// LoadSettings - функция установки значений в структуру
LoadSettings(settings interface{}, v *viper.Viper) error
````
Данный способ под капотом также устанавливает хук на syslog и graylog, если используется [стандартная модель](types.go) настроек:

```go
type Logger struct {
	LogLevel       logrus.Level `env:"LOG_LEVEL" toml:"level" json:"level" xml:"level" yaml:"level" default:"debug"`
	Syslog         string       `env:"SYSLOG" toml:"syslog_addr" json:"syslog_addr" xml:"syslog_addr" yaml:"syslog_addr" default:"127.0.0.1:514" validate:"tcp_addr"`
	SyslogProtocol string       `env:"SYSLOG_PROTOCOL" toml:"syslog_protocol" json:"syslog_protocol" xml:"syslog_protocol" yaml:"syslog_protocol" default:"udp" validate:"min=3,max=3"`
	SysLogLevel    SyslogLevel  `env:"SYSLOG_LEVEL" toml:"syslog_level" json:"syslog_level" xml:"syslog_level" yaml:"syslog_level" default:"debug"`
	Colour         bool         `env:"COLOUR" toml:"colour" json:"colour" xml:"colour" yaml:"colour" default:"false"`
	StdOut         bool         `env:"STDOUT" toml:"stdout" json:"stdout" xml:"stdout" yaml:"stdout" default:"true"`
	GraylogLevel   logrus.Level `env:"GRAYLOG_LEVEL" toml:"graylog_level" json:"graylog_level" xml:"graylog_level" yaml:"graylog_level" default:"debug"`
	Graylog        string       `env:"GRAYLOG" toml:"graylog" json:"graylog" xml:"graylog" yaml:"graylog"`
}
```

#### 1.2 Сущность конфигуратора
Данный способ имеет большее количество возможностей, что позволяет выполнять тонкую настройку пакета.
Рассмотрим методику создания сущности:
```go
func main() {
	configurator = config.New("filename", "path")
}
```


### Файлы конфигурации

### Поддерживаемые типы