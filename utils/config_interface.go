package utils

type IConfig interface {

	// FromKey returns section of the config
	FromKey(key string) IConfig

	// GetCfg returns current config as a 1st level map of interfaces
	GetCfg() map[string]interface{}

	// GetValue returns a value as an typeless interface
	GetValue(key string) interface{}

	// GetFloatDefault, etc methods to access section/config variables
	GetFloatDefault(key string, dflt float64) float64
	GetIntDefault(key string, dflt int64) int64
	GetBoolDefault(key string, dflt bool) bool
	GetString(key string) *string
	GetStringDefault(key string, defaultVal string) *string

	GetStringList(key string) []string

	// ReadConfig reads config from a file
	ReadConfig(filename string) IConfig

	// WriteConfigX Writes config back to its file
	WriteConfigX() error

	// ModifiedQ true if modified since last read
	ModifiedQ() bool

	// GetFilename config filename
	GetFileName() string

	// Set a new name
	SetFileName(string)
	// SetRO sets/resets Read-only flag.
	SetRO(bool)

	// GetRO returns Read-only status
	GetRO() bool
}
