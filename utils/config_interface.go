package utils

type IConfig interface {
	FromKey(key string) IConfig
	GetCfg() map[string]interface{}
	GetValue(key string) interface{}
	GetFloatDefault(key string, dflt float64) float64
	GetIntDefault(key string, dflt int64) int64
	GetBoolDefault(key string, dflt bool) bool
	GetString(key string) *string
	GetStringDefault(key string, defaultVal string) *string
	ReadConfig(filename string) IConfig
	WriteConfigX() error
}
