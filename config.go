package d

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config interface, implement at least the following methods to facilitate internal calls in the devtool library
type InterfaceConfig interface {
	Init()
	GetIntWithDefault(key string, default_value int) int
	GetStringWithDefault(key, default_value string) string
	GetStringMap(key string) map[string]interface{}
	GetBool(key string) bool
	Set(key string, value interface{}) error
}

var (
	config InterfaceConfig // Global variable, stores the initialized interface, if not initialized, it is nil
)

// Config library unified access entry
type Config[T InterfaceConfig] struct{}

// Initialization
func (c Config[T]) Init(conf T) {
	config = conf
}

// Get the initialized interface. If it is not initialized, Viper library is used by default.
func (c Config[T]) Get() T {
	if config == nil {
		LibraryViper{}.Init()
	}
	return config.(T)
}

// Viper library
type LibraryViper struct {
	*viper.Viper
	SetConfigName  string
	AddConfigPath  string
	OnConfigChange func(e fsnotify.Event) // This method is triggered when the configuration file changes
}

// Initialization
func (l LibraryViper) Init() {
	// Create a new viper instance
	var conf = viper.New()

	// Set the default value
	if l.SetConfigName == "" {
		l.SetConfigName = "config"
	}
	if l.AddConfigPath == "" {
		l.AddConfigPath = "."
	}
	if l.OnConfigChange == nil {
		l.OnConfigChange = func(e fsnotify.Event) {
			fmt.Println("Config file changed:", e.Name)
		}
	}

	conf.SetConfigName(l.SetConfigName)
	conf.AddConfigPath(l.AddConfigPath)
	err := conf.ReadInConfig()
	if err != nil {
		panic(err)
	}

	conf.OnConfigChange(l.OnConfigChange)
	conf.WatchConfig()
	Config[LibraryViper]{}.Init(LibraryViper{Viper: conf})
}

// Get the int. If there is no value, get the default value of the setting.
func (l LibraryViper) GetIntWithDefault(key string, default_value int) int {
	Config[LibraryViper]{}.Get().Viper.SetDefault(key, default_value)
	return Config[LibraryViper]{}.Get().Viper.GetInt(key)
}

// Get string
func (l LibraryViper) GetString(key string) string {
	return Config[LibraryViper]{}.Get().Viper.GetString(key)
}

// Get the string. If there is no value, get the default value of the setting.
func (l LibraryViper) GetStringWithDefault(key, default_value string) string {
	Config[LibraryViper]{}.Get().Viper.SetDefault(key, default_value)
	return Config[LibraryViper]{}.Get().Viper.GetString(key)
}

// Get string map
func (l LibraryViper) GetStringMap(key string) map[string]interface{} {
	return Config[LibraryViper]{}.Get().Viper.GetStringMap(key)
}

// Get string map
func (l LibraryViper) GetBool(key string) bool {
	return Config[LibraryViper]{}.Get().Viper.GetBool(key)
}

func (l LibraryViper) Set(key string, value interface{}) error {
	Config[LibraryViper]{}.Get().Viper.Set(key, value)
	return Config[LibraryViper]{}.Get().Viper.WriteConfig()
}
