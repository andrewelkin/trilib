package utils

import (
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

// Vconfig is just a wrapper around viper Config
type Vconfig struct {
	*viper.Viper
	ro           bool
	lastModified time.Time // last mod timestamp
	fileName     string
}

// SetRO sets/resets Read-only flag.
func (c *Vconfig) SetRO(ro bool) {
	c.ro = ro
}

// GetRO returns Read-only status
func (c *Vconfig) GetRO() bool {
	return c.ro
}

// ModifiedQ returns true if the file modified since last ReadFile
func (c *Vconfig) ModifiedQ() bool {
	fileStats, _ := os.Stat(c.GetFileName())
	return c.lastModified != fileStats.ModTime()
}

func (c *Vconfig) GetFileName() string {
	return c.fileName
}

// FromKey creates a deep copy of the Vconfig section
// this fixes viper bug when calling Sub() for hcl format
func (c *Vconfig) FromKey(key string) IConfig {
	key = strings.ToLower(key)

	subv := viper.New()
	data := c.Get(key)
	if data == nil {
		return nil
	}

	if reflect.TypeOf(data).Kind() == reflect.Map { // never happens in fact, always []map[string]interface{}
		subv.MergeConfigMap(cast.ToStringMap(data))
	} else if reflect.TypeOf(data).Kind() == reflect.Slice {
		sl := cast.ToSlice(data)
		if len(sl) != 1 {
			Throwf("error getting section %v", key)
		}
		subv.MergeConfigMap(cast.ToStringMap(sl[0]))
	} else {
		Throwf("error getting section %v", key)
	}

	return &Vconfig{
		Viper:        subv,
		fileName:     c.GetFileName(),
		lastModified: c.lastModified,
	}
}

func (c *Vconfig) GetStringList(key string) []string {
	rawList := c.GetString(key)
	return strings.Split(*rawList, ",")
}

func (c *Vconfig) GetCfg() map[string]interface{} {

	data := c.AllSettings()
	if reflect.TypeOf(data).Kind() == reflect.Map { // with hcl never happens in fact, always []map[string]interface{}
		return data
	} else if reflect.TypeOf(data).Kind() == reflect.Slice {
		sl := cast.ToSlice(data)
		if len(sl) == 1 {
			return cast.ToStringMap(sl[0])
		}
	}
	Throwf("error getting cfg (casting to map)")
	return nil
}

func (c *Vconfig) GetValue(key string) interface{} {
	key = strings.ToLower(key)
	return c.Get(key)
}

func (c *Vconfig) GetFloatDefault(key string, dflt float64) float64 {
	key = strings.ToLower(key)
	c.SetDefault(key, dflt)
	return c.GetFloat64(key)
}

func (c *Vconfig) GetIntDefault(key string, dflt int64) int64 {
	key = strings.ToLower(key)
	c.SetDefault(key, dflt)
	return c.GetInt64(key)
}

func (c *Vconfig) GetBoolDefault(key string, dflt bool) bool {
	key = strings.ToLower(key)
	c.SetDefault(key, dflt)
	return c.GetBool(key)
}

func (c *Vconfig) GetString(key string) *string {
	key = strings.ToLower(key)
	rs := c.Viper.GetString(key)

	if strings.HasPrefix(rs, "$") {
		envStr := os.Getenv(rs[1:])
		return &envStr
	}
	// replace value with environment variable, windows %xx% style
	if strings.HasPrefix(rs, "%") {
		ndx := strings.Index(rs[1:], "%")
		if ndx < 0 {
			envStr := os.Getenv(rs[1:])
			return &envStr
		}

		envStr := os.Getenv(rs[1 : ndx+1])
		d := envStr + rs[ndx+2:]
		return &d
	}

	return &rs
}

func (c *Vconfig) GetStringDefault(key string, defaultVal string) *string {
	key = strings.ToLower(key)
	c.SetDefault(key, defaultVal)
	return c.GetString(key)
}

// ReadConfig reads json Vconfig. Throws if something is wrong
func (c *Vconfig) ReadConfig(filename string) IConfig {

	if c.Viper == nil {
		c.Viper = viper.New()
	}

	c.fileName = filename
	fileStats, _ := os.Stat(filename)
	c.lastModified = fileStats.ModTime()

	onlyName := strings.TrimRight(strings.Replace(filepath.Base(filename), filepath.Ext(filename), "", 1), ".")
	c.SetConfigName(onlyName) // name of config file (without extension) -- what a f innovation!
	c.SetConfigType(strings.ToLower(strings.TrimLeft(filepath.Ext(filename), ".")))
	path := filepath.Dir(filename)
	if len(path) == 0 {
		path = "./"
	}
	c.AddConfigPath(path)
	err := c.ReadInConfig()
	if err != nil {
		Throwf("fatal error reading config file: %w", err)
	}

	return c
}

// WriteConfigX writes a section of a Vconfig back to file
func (c *Vconfig) WriteConfigX() error {
	Throwf("WriteConfigX not implemented yet")
	return nil
}

func (c *Vconfig) SetFileName(newName string) {
	c.fileName = newName
}
