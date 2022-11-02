## Оглавление
1. [Описание](#desc)
2. [Файлы конфигурации](#files)
3. [Поддерживаемые типы](#types)
4. [Примеры](#examples)

---

<a name="desc"></a>
### 1. Описание

Для парсинга файлов используется библиотека [`viper`](https://github.com/spf13/viper).

Существует два метода использования данного пакета: 
* базовый метод загрузки;
* через сущность конфигуратора.

<hr>

#### 1.1 Базовый метод загрузки настроек

Первый вариант довольно прост и имеет в себе всего две функции:

````go
// LoadOptions - функция загрузки настроек из файла
LoadOptions(name string, paths ...string) (*viper.Viper, error)
// LoadSettings - функция установки значений в структуру
LoadSettings(settings interface{}, v *viper.Viper) error
````

Данный способ также устанавливает хук на syslog, если используется [стандартная модель](../types/logger_unix_syslog.go) настроек логгера:

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

#### 1.2 Конфигуратор
Данный способ имеет сравнительно больше возможностей, что позволяет выполнять тонкую настройку загрузки настроек. 
Конфигуратор определен интерфейсом IConfigurator:

```go
type IConfigurator interface {
	// ReadConfiguration - загрузка viper из внешнего io.Reader
	ReadConfiguration(config io.Reader) error
	// LoadConfiguration - загрузка загрузить viper из файла, который был установлен при создании конфигуратора
	LoadConfiguration() error
	// Unmarshal - установка значения из viper в структуру, указатель на которую передается в качестве аргумента
	Unmarshal(config interface{}) error
	// SetOption - настройка конфигуратора
	SetOption(options types.Options, value interface{}) error
}

```
Конфигуратор может получить viper как из файла, так из иного io.Reader. Рассмотрим методику работы с конфигуратором.

```go
func main() {
	// Указатель на структуру с настройками
	config := new(Config)
	// Инициализация конфигуратора
	configurator = config.New("filename", "path")
	// Загружаем конфигурацию из файла
	if err := configurator.LoadConfiguration(); err != nil {
	    logrus.WithError(err).Panic("Unable to load configuration")	
        }
	// Выгружаем конфигурацию в структуру
	if err := configurator.Unmarshal(config); err != nil {
            logrus.WithError(err).Panic("Unable to unmarshal configuration")
        }
}
```

Тонкая настройка подразумевает следующий набор возможностей:

| Option                 | Type           | Description                                                                       |
|------------------------|----------------|-----------------------------------------------------------------------------------|
| ```TimeFormat```       | string         | установку формата времени для определения даты                                    |
| ```ProcessingMode```   | string         | выбор метода процессинга (overwriting/complement)                                 |
| ```LoggerHook```       | bool           | определение необходимости установки хука на syslog                                |
| ```LoggerInstance```   | *logrus.Logger | работа с установленным ```*logrus.Logger``` (используется при инициализации хука) |
| ```ValidatorEnable```  | bool           | определение необходимости установки валидации                                     |

Пример установки опции конфигуратора:

```go
// Инициализация конфигуратора
configurator = config.New("filename", "path")
if err := configurator.SetOption(types.TimeFormat, time.RFC3339Nano); err != nil {
    logrus.WithError(err).Panic("Unable to set option")
}
```

### Теги и структуры

Если нужно чтобы поле конфигурации было обработано функций ```LoadSettings```, его нужно обеспечить тегами.

Имеется следующий список тегов:
* [env](#env)
* [json, toml, yaml, xml](#general)
* [default](#default)
* [omit](#omit)
* [validate](#validate)

<a name="env"></a>
#### ENV
Тег `env` указывает что значение поля должно загружаться из переменных окружения, тег должен содержать имя переменной
окружения.

Например:

```go
type MoreSettings struct {
    CacheSize   byte `env:"CACHE_SIZE"`
}
```

Есть исключение для полей типа slice, map, struct. Для них тег `env` применять запрещено. Для явного поведения, функция вернёт ошибку
если найдет такой тег у полей типа slice, map, struct.

<a name="general"></a>
#### JSON, TOML, YAML, XML
Теги `json`, `toml`, `yaml`, `xml` указывает что значение поля должно загружаться из файла соответствующего расширения. Тег должен содержать путь к переменной.

Например:

```go
type LogSettings struct {
    LogLevel   byte `toml:"log.level"`
}
```

Также теги могут быть применены для структур, что позволяет объявлять родительский узел настроек для
всех настроек внутри структуры.

Например:
```go
type MainSettings struct {
    LogSettings `toml:"log"`
}

type LogSettings struct {
    LogLevel   byte `toml:"level"`
}
```

Теги исправно работают даже если вложенная структура не имеет тега.

Например:
```go
type MainSettings struct {
    LogSettings `toml:"log"`
}

type LogSettings struct {
    Main
}

type Main struct {
    LogLevel   byte `toml:"level"`
}
```


<a name="default"></a>
#### DEFAULT
Тег `default` определяет дефолтное значение, которое применяется в случае если значение не найдено в переменных
окружения или в toml-файле.

<a name="omit"></a>
#### OMIT
Тег `omit` указывает что поле должно быть проигнорировано при обработке. Это в первую очередь нужно для игнорирования
полей типа указатель и структура, которые по умолчанию создаются и заполняются.

```go
type Settings struct {
    StructMustBeOmitted   *MyStruct `omit:"-"`
}
```
<a name="validate"></a>
#### VALIDATE
Тег `validate` позволяет указывать правила валидации для полей. Тег заполняется в соответствии с описанием пакета
[validator](https://github.com/go-playground/validator).

### Файлы конфигурации
<a name="files"></a>

### Поддерживаемые типы
<a name="types"></a>

### Примеры
<a name="examples"></a>